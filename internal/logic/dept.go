// +======================================================================
// | EasyGoAdmin 企业数字化开发平台 [ 赋能开发者，助力企业发展 ]
// +----------------------------------------------------------------------
// | Copyright (c) 2019-2026 半城风雨团队 版权所有
// +----------------------------------------------------------------------
// | 开源协议: Apache-2.0
// | 重要声明: 未经书面许可，不得移除或修改本框架相关版权标识。
// +----------------------------------------------------------------------
// | 官方网站: https://www.easygoadmin.vip
// +----------------------------------------------------------------------
// | 研发团队: 小蚂蚁云团队 荣誉出品
// +======================================================================
// | 版权和免责声明:
// | 一、知识产权声明
// | 本团队对该软件框架产品（以下简称“本产品”）拥有完整知识产权，包括但不限于商标权、专利权、
// | 著作权、商业秘密等。前述权利均受中华人民共和国及相关国际法律法规保护。未经本团队书面授权，
// | 任何个人或单位不得基于本产品本身申请任何形式的知识产权。
// | 二、使用限制与禁止条款
// | 您理解并同意，禁止将本产品用于以下任何行为：
// | 1. 任何违反中华人民共和国法律法规的用途；
// | 2. 任何侵害他人合法权益（包括但不限于人身权、财产权、知识产权等）的行为；
// | 3. 任何危害国家安全、公共利益及社会道德风尚的行为。
// | 三、责任豁免
// | 因使用本产品而产生的任何损失或法律纠纷（包括但不限于直接损失、间接损失、附带损失、衍生损失
// | 以及因违约、侵权等引发的责任），本团队均不承担任何法律责任。若因违规使用本产品导致任何第三
// | 方提起索赔或行政机关进行调查，您应独立承担全部责任。
// | 四、违约与法律责任
// | 如发现任何违反本声明的行为，本团队有权：
// | 1. 无条件配合公安机关及其他有权机关的调查取证工作；
// | 2. 依法追究相关方的民事责任，并要求赔偿全部损失；
// | 3. 保留进一步采取刑事控告等法律措施的权利。
// | 五、许可范围说明
// | 1、被授权方应视本产品源代码为核心机密信息，负有最高标准的保密与保管义务。禁止以任何形式（包
// | 括但不限于直接或间接的分享、复制、传播、分发、转售、出租、开源、反向工程、反编译、反汇编或通
// | 过网络服务提供给任何第三方）泄露、披露或处置源代码的全部或部分内容。任何未经授权的复制、使用、
// | 传播行为均构成对本团队知识产权的严重侵权，本团队有权依法追究其全部法律责任，并要求赔偿因此造
// | 成的一切经济损失及商誉损害。
// | 2、本产品仅限于公司或个人在法律允许范围内，用于内部且合规的软件产品研发。具体权利义务以《软件
// | 许可协议》附件为准。
// +======================================================================

package logic

import (
	"context"

	"github.com/gogf/gf/v2/errors/gerror"

	"easygoadmin/api/v1"
	"easygoadmin/internal/consts"
	"easygoadmin/internal/dao"
	"easygoadmin/internal/service"
)

// ============================== 部门业务实现 ==============================

// sDept 部门业务实现
// 实现 service.IDept 接口，接口定义见 api/v1/dept.go
type sDept struct{}

func init() {
	service.RegisterDept(&sDept{})
}

// deptSerializeMaps 枚举字段到字典编码的映射
// {field}Text 由基座自动填充
var deptSerializeMaps = map[string]string{
	"type": "dept_type",
}

// deptSpec 部门模块声明
//
// 平铺列表 + 上级循环校验 + 删除守卫；树形结构由前端 buildTree 构建
var deptSpec = ModuleSpec[v1.DeptItem, v1.DeptItem, v1.DeptDetailRes]{
	Model:         dao.Dept.Ctx,
	LikeFields:    []string{"name"},
	EqFields:      []string{"type"},
	OrderBy:       "sort, id",
	SerializeMaps: deptSerializeMaps,
	ToData:        deptToData,
	BeforeUpdate:  deptParentGuard, // 禁止挂到自身/后代下
	BeforeDelete:  deptDeleteGuard, // 存在性 + 有子部门/用户引用禁止
}

// deptToData 全量表单语义转换
// 将新增或修改请求转换为落库字段，tenant_id 由 hooks.ApplyCtx 自动填充
func deptToData(req any) map[string]any {
	switch r := req.(type) {
	case *v1.DeptAddReq:
		return map[string]any{
			"name": r.Name, "type": r.Type, "parent_id": r.ParentId,
			"sort": r.Sort, "note": r.Note,
		}
	case *v1.DeptUpdateReq:
		return map[string]any{
			"name": r.Name, "type": r.Type, "parent_id": r.ParentId,
			"sort": r.Sort, "note": r.Note,
		}
	}
	panic("deptToData: 不支持的入参类型")
}

// deptParentGuard 上级循环校验
// 禁止将上级设置为自身或自己的后代节点
func deptParentGuard(ctx context.Context, data map[string]any, id int64, req any) error {
	r := mustAs[*v1.DeptUpdateReq](req)
	if r.ParentId <= 0 {
		return nil
	}
	if int64(r.ParentId) == id {
		return gerror.NewCodef(consts.CodeParamError, "上级不能选择自己")
	}
	descIds, err := deptDescendantIds(ctx, id)
	if err != nil {
		return err
	}
	for _, did := range descIds {
		if did == int64(r.ParentId) {
			return gerror.NewCodef(consts.CodeParamError, "上级不能选择自己的下级")
		}
	}
	return nil
}

// deptDeleteGuard 删除守卫
// 校验记录存在性，存在子部门或用户引用时禁止删除
func deptDeleteGuard(ctx context.Context, id int64) error {
	if n, err := dao.Dept.Ctx(ctx).Where("id", id).Count(); err != nil {
		return err
	} else if n == 0 {
		return gerror.NewCode(consts.CodeRecordNotFound)
	}
	if n, err := dao.Dept.Ctx(ctx).Where("parent_id", id).Count(); err != nil {
		return err
	} else if n > 0 {
		return gerror.NewCode(consts.CodeHasChildren)
	}
	if n, err := dao.User.Ctx(ctx).Where("dept_id", id).Count(); err != nil {
		return err
	} else if n > 0 {
		return gerror.NewCode(consts.CodeInUse)
	}
	return nil
}

// deptDescendantIds 递归获取全部后代ID
// 采用广度优先遍历，返回指定部门下所有层级的后代节点ID
func deptDescendantIds(ctx context.Context, deptId int64) ([]int64, error) {
	result := make([]int64, 0)
	queue := []int64{deptId}
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		rows, err := dao.Dept.Ctx(ctx).Where("parent_id", current).Fields("id").All()
		if err != nil {
			return nil, err
		}
		for _, r := range rows {
			idv := r["id"].Int64()
			result = append(result, idv)
			queue = append(queue, idv)
		}
	}
	return result, nil
}

// ---------- 标准方法：3行委托 ----------

// Page 部门分页（平铺）
// 按查询条件分页返回平铺的部门列表
func (s *sDept) Page(ctx context.Context, req *v1.DeptPageReq) (res *v1.DeptPageRes, err error) {
	list, total, err := deptSpec.Page(ctx, req, req.GetPageNo(), req.GetPageSize())
	if err != nil {
		return nil, err
	}
	return v1.NewPageData(req.GetPageNo(), req.GetPageSize(), total, derefSlice(list)), nil
}

// List 部门平铺列表
// data 直接为数组，树形结构由前端 buildTree 构建
func (s *sDept) List(ctx context.Context, req *v1.DeptListReq) (res v1.DeptListRes, err error) {
	list, err := deptSpec.List(ctx, req)
	if err != nil {
		return nil, err
	}
	return list, nil
}

// Detail 部门详情
// 按部门ID查询详情
func (s *sDept) Detail(ctx context.Context, req *v1.DeptDetailReq) (res *v1.DeptDetailRes, err error) {
	return deptSpec.Detail(ctx, req.Id, req)
}

// Add 新增部门
// 新增成功后返回新记录ID
func (s *sDept) Add(ctx context.Context, req *v1.DeptAddReq) (res *v1.IdRes, err error) {
	id, err := deptSpec.Add(ctx, req)
	if err != nil {
		return nil, err
	}
	return &v1.IdRes{Id: id}, nil
}

// Update 修改部门
// 按部门ID更新记录
func (s *sDept) Update(ctx context.Context, req *v1.DeptUpdateReq) (res *v1.DeptUpdateRes, err error) {
	if err = deptSpec.Update(ctx, req.Id, req); err != nil {
		return nil, err
	}
	return &v1.DeptUpdateRes{}, nil
}

// Delete 删除部门
// 按部门ID删除单条记录
func (s *sDept) Delete(ctx context.Context, req *v1.DeptDeleteReq) (res *v1.DeptDeleteRes, err error) {
	if err = deptSpec.Delete(ctx, req.Id); err != nil {
		return nil, err
	}
	return &v1.DeptDeleteRes{}, nil
}

// BatchDelete 批量删除部门
// 按部门ID列表批量删除记录
func (s *sDept) BatchDelete(ctx context.Context, req *v1.DeptBatchDeleteReq) (res *v1.DeptBatchDeleteRes, err error) {
	if err = deptSpec.BatchDelete(ctx, req.Ids); err != nil {
		return nil, err
	}
	return &v1.DeptBatchDeleteRes{}, nil
}
