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

// ============================== 岗位业务实现 ==============================

// sPosition 岗位业务实现
// 实现 service.IPosition 接口，接口定义见 api/v1/position.go
type sPosition struct{}

func init() {
	service.RegisterPosition(&sPosition{})
}

// positionSerializeMaps 枚举字段到字典编码的映射
// {field}Text 由基座自动填充
var positionSerializeMaps = map[string]string{
	"status": "position_status",
}

// positionSpec 岗位模块声明
var positionSpec = ModuleSpec[v1.PositionItem, v1.PositionListItem, v1.PositionDetailRes]{
	Model:         dao.Position.Ctx,
	LikeFields:    []string{"name"},
	EqFields:      []string{"status"},
	OrderBy:       "sort, id",
	SerializeMaps: positionSerializeMaps,
	ToData:        positionToData,
	BeforeAdd:     positionDefaultStatus, // status 缺省置 1
	BeforeDelete:  positionDeleteGuard,   // 存在性 + 用户引用检查
}

// positionToData 全量表单语义转换
// 兼容新增与修改请求，将入参转换为落库字段
func positionToData(req any) map[string]any {
	switch r := req.(type) {
	case *v1.PositionAddReq:
		return map[string]any{"name": r.Name, "status": r.Status, "sort": r.Sort}
	case *v1.PositionUpdateReq:
		return map[string]any{"name": r.Name, "status": r.Status, "sort": r.Sort}
	}
	panic("positionToData: 不支持的入参类型")
}

// positionDefaultStatus 状态缺省处理
// 状态为零值时置为 1（正常）
func positionDefaultStatus(ctx context.Context, data map[string]any, req any) error {
	if v, ok := data["status"].(int); ok && v == 0 {
		data["status"] = 1
	}
	return nil
}

// positionDeleteGuard 删除守卫
// 校验记录存在性，存在用户引用时禁止删除
func positionDeleteGuard(ctx context.Context, id int64) error {
	if n, err := dao.Position.Ctx(ctx).Where("id", id).Count(); err != nil {
		return err
	} else if n == 0 {
		return gerror.NewCode(consts.CodeRecordNotFound)
	}
	if n, err := dao.User.Ctx(ctx).Where("position_id", id).Count(); err != nil {
		return err
	} else if n > 0 {
		return gerror.NewCode(consts.CodeInUse)
	}
	return nil
}

// Page 岗位分页
// 分页返回岗位列表，并补充关联用户数与状态文本
func (s *sPosition) Page(ctx context.Context, req *v1.PositionPageReq) (res *v1.PositionPageRes, err error) {
	list, total, err := positionSpec.Page(ctx, req, req.GetPageNo(), req.GetPageSize())
	if err != nil {
		return nil, err
	}
	fillUserCount(ctx, "position_id", &list)
	return v1.NewPageData(req.GetPageNo(), req.GetPageSize(), total, derefSlice(list)), nil
}

// List 岗位全量列表
// 按关键词、状态等条件返回岗位列表
func (s *sPosition) List(ctx context.Context, req *v1.PositionListReq) (res v1.PositionListRes, err error) {
	list, err := positionSpec.List(ctx, req)
	if err != nil {
		return nil, err
	}
	return derefSlice(list), nil
}

// Detail 岗位详情
// 按岗位ID查询详情
func (s *sPosition) Detail(ctx context.Context, req *v1.PositionDetailReq) (res *v1.PositionDetailRes, err error) {
	return positionSpec.Detail(ctx, req.Id, req)
}

// Add 新增岗位
// 新增成功后返回新记录ID
func (s *sPosition) Add(ctx context.Context, req *v1.PositionAddReq) (res *v1.IdRes, err error) {
	id, err := positionSpec.Add(ctx, req)
	if err != nil {
		return nil, err
	}
	return &v1.IdRes{Id: id}, nil
}

// Update 修改岗位
// 按岗位ID更新记录
func (s *sPosition) Update(ctx context.Context, req *v1.PositionUpdateReq) (res *v1.PositionUpdateRes, err error) {
	if err = positionSpec.Update(ctx, req.Id, req); err != nil {
		return nil, err
	}
	return &v1.PositionUpdateRes{}, nil
}

// Status 岗位状态
// 非标准方法，手写实现，设置岗位的正常或禁用状态
func (s *sPosition) Status(ctx context.Context, req *v1.PositionStatusReq) (res *v1.PositionStatusRes, err error) {
	if _, err = dao.Position.Ctx(ctx).Where("id", req.Id).Data(map[string]any{"status": req.Status}).Update(); err != nil {
		return nil, err
	}
	return &v1.PositionStatusRes{}, nil
}

// Delete 删除岗位
// 按岗位ID删除单条记录
func (s *sPosition) Delete(ctx context.Context, req *v1.PositionDeleteReq) (res *v1.PositionDeleteRes, err error) {
	if err = positionSpec.Delete(ctx, req.Id); err != nil {
		return nil, err
	}
	return &v1.PositionDeleteRes{}, nil
}

// BatchDelete 批量删除岗位
// 按岗位ID列表批量删除记录
func (s *sPosition) BatchDelete(ctx context.Context, req *v1.PositionBatchDeleteReq) (res *v1.PositionBatchDeleteRes, err error) {
	if err = positionSpec.BatchDelete(ctx, req.Ids); err != nil {
		return nil, err
	}
	return &v1.PositionBatchDeleteRes{}, nil
}
