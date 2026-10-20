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

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"

	"easygoadmin/api/v1"
	"easygoadmin/internal/consts"
	"easygoadmin/internal/dao"
	"easygoadmin/internal/service"
	"easygoadmin/utility/cache"
)

// ============================== 菜单业务实现 ==============================

// sMenu 菜单业务实现
// 实现 service.IMenu 接口，接口定义见 api/v1/menu.go
type sMenu struct{}

func init() {
	service.RegisterMenu(&sMenu{})
}

// menuSerializeMaps 枚举字段到字典编码的映射
// {field}Text 由基座自动填充
var menuSerializeMaps = map[string]string{
	"type":   "menu_type",
	"status": "menu_status",
}

// menuSpec 菜单模块声明
//
// 平铺列表 + 上级循环校验 + 删除守卫，并连带清理 role_menu 与缓存；
// 树形结构由前端 buildTree 构建
var menuSpec = ModuleSpec[v1.MenuItem, v1.MenuItem, v1.MenuDetailRes]{
	Model:         dao.Menu.Ctx,
	LikeFields:    []string{"name", "path"},
	OrderBy:       "sort, id",
	SerializeMaps: menuSerializeMaps,
	Filter:        menuFilter, // Type/Status 为 *int 指针（nil=不筛，指向0也筛），EqFields 表达不了
	ToData:        menuToData,
	BeforeUpdate:  menuParentGuard,    // 禁止挂到自身/后代下
	BeforeDelete:  menuDeleteGuard,    // 存在性 + 有子菜单禁止
	AfterDelete:   menuCleanRelations, // 清理 role_menu + 失效受影响用户缓存
	AfterBatchDelete: func(ctx context.Context, ids []int64) error {
		return menuCleanRelationsBatch(ctx, ids)
	},
}

// menuFilter 复杂筛选
// Type 与 Status 为指针字段，nil 时跳过，指向 0 时同样生效
func menuFilter(ctx context.Context, m *gdb.Model, req any) *gdb.Model {
	switch r := req.(type) {
	case *v1.MenuPageReq:
		if r.Type != nil {
			m = m.Where("type", *r.Type)
		}
		if r.Status != nil {
			m = m.Where("status", *r.Status)
		}
	case *v1.MenuListReq:
		if r.Type != nil {
			m = m.Where("type", *r.Type)
		}
		if r.Status != nil {
			m = m.Where("status", *r.Status)
		}
	}
	return m
}

// menuToData 全量表单语义转换
// 新增与修改均提交全部字段；修改时 icon/icon2/note 未传值则不落键，保留库中原值
func menuToData(req any) map[string]any {
	switch r := req.(type) {
	case *v1.MenuAddReq:
		return map[string]any{
			"parent_id": r.ParentId, "name": r.Name, "icon": r.Icon, "icon2": r.Icon2, "path": r.Path,
			"component": r.Component, "permission": r.Permission, "target": r.Target,
			"type": r.Type, "status": r.Status, "hide": r.Hide, "sort": r.Sort, "note": r.Note,
		}
	case *v1.MenuUpdateReq:
		data := map[string]any{
			"parent_id": r.ParentId, "name": r.Name, "path": r.Path,
			"component": r.Component, "permission": r.Permission, "target": r.Target,
			"type": r.Type, "status": r.Status, "hide": r.Hide, "sort": r.Sort,
		}
		if r.Icon != nil {
			data["icon"] = *r.Icon
		}
		if r.Icon2 != nil {
			data["icon2"] = *r.Icon2
		}
		if r.Note != nil {
			data["note"] = *r.Note
		}
		return data
	}
	panic("menuToData: 不支持的入参类型")
}

// menuParentGuard 上级循环校验
// 禁止将上级设置为自身或自己的后代节点
func menuParentGuard(ctx context.Context, data map[string]any, id int64, req any) error {
	r := mustAs[*v1.MenuUpdateReq](req)
	if r.ParentId <= 0 {
		return nil
	}
	if int64(r.ParentId) == id {
		return gerror.NewCodef(consts.CodeParamError, "上级不能选择自己")
	}
	descIds, err := menuDescendantIds(ctx, id)
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

// menuDeleteGuard 删除守卫
// 校验记录存在性，存在子菜单时禁止删除
func menuDeleteGuard(ctx context.Context, id int64) error {
	if n, err := dao.Menu.Ctx(ctx).Where("id", id).Count(); err != nil {
		return err
	} else if n == 0 {
		return gerror.NewCode(consts.CodeRecordNotFound)
	}
	if n, err := dao.Menu.Ctx(ctx).Where("parent_id", id).Count(); err != nil {
		return err
	} else if n > 0 {
		return gerror.NewCode(consts.CodeHasChildren)
	}
	return nil
}

// menuCleanRelations 单条删除后清理
// 清理 role_menu 关联并失效受影响用户的缓存
func menuCleanRelations(ctx context.Context, id int64) error {
	return menuCleanRelationsBatch(ctx, []int64{id})
}

// menuCleanRelationsBatch 批量删除后清理
// 清理 role_menu 关联并失效受影响用户的缓存
func menuCleanRelationsBatch(ctx context.Context, ids []int64) error {
	roleMenus, err := dao.RoleMenu.Ctx(ctx).WhereIn("menu_id", ids).All()
	if err != nil {
		return err
	}
	if _, err = dao.RoleMenu.Ctx(ctx).WhereIn("menu_id", ids).Delete(); err != nil {
		return err
	}
	if len(roleMenus) == 0 {
		return nil
	}
	roleIds := make([]int64, 0, len(roleMenus))
	for _, rm := range roleMenus {
		roleIds = append(roleIds, rm["role_id"].Int64())
	}
	urRows, err := dao.UserRole.Ctx(ctx).WhereIn("role_id", roleIds).All()
	if err != nil {
		return err
	}
	for _, ur := range urRows {
		uid := int(ur["user_id"].Int64())
		cache.Del(ctx,
			consts.CacheKeyUserIdentity(uid),
			consts.CacheKeyUserPerms(uid),
			consts.CacheKeyUserRoles(uid),
		)
	}
	return nil
}

// menuDescendantIds 递归获取全部后代ID
// 采用广度优先遍历，返回指定菜单下所有层级的后代节点ID
func menuDescendantIds(ctx context.Context, menuId int64) ([]int64, error) {
	result := make([]int64, 0)
	queue := []int64{menuId}
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		rows, err := dao.Menu.Ctx(ctx).Where("parent_id", current).Fields("id").All()
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

// Page 菜单分页（平铺）
// 按查询条件分页返回平铺的菜单列表
func (s *sMenu) Page(ctx context.Context, req *v1.MenuPageReq) (res *v1.MenuPageRes, err error) {
	list, total, err := menuSpec.Page(ctx, req, req.GetPageNo(), req.GetPageSize())
	if err != nil {
		return nil, err
	}
	return v1.NewPageData(req.GetPageNo(), req.GetPageSize(), total, derefSlice(list)), nil
}

// List 菜单平铺列表
// data 直接为数组，树形结构由前端 buildTree 构建
func (s *sMenu) List(ctx context.Context, req *v1.MenuListReq) (res v1.MenuListRes, err error) {
	list, err := menuSpec.List(ctx, req)
	if err != nil {
		return nil, err
	}
	return list, nil
}

// Detail 菜单详情
// 按菜单ID查询详情
func (s *sMenu) Detail(ctx context.Context, req *v1.MenuDetailReq) (res *v1.MenuDetailRes, err error) {
	return menuSpec.Detail(ctx, req.Id, req)
}

// Add 新增菜单
// 新增成功后返回新记录ID
func (s *sMenu) Add(ctx context.Context, req *v1.MenuAddReq) (res *v1.IdRes, err error) {
	id, err := menuSpec.Add(ctx, req)
	if err != nil {
		return nil, err
	}
	return &v1.IdRes{Id: id}, nil
}

// Update 修改菜单
// 按菜单ID更新记录
func (s *sMenu) Update(ctx context.Context, req *v1.MenuUpdateReq) (res *v1.MenuUpdateRes, err error) {
	if err = menuSpec.Update(ctx, req.Id, req); err != nil {
		return nil, err
	}
	return &v1.MenuUpdateRes{}, nil
}

// Delete 删除菜单
// 按菜单ID删除单条记录，并清理关联关系与缓存
func (s *sMenu) Delete(ctx context.Context, req *v1.MenuDeleteReq) (res *v1.MenuDeleteRes, err error) {
	if err = menuSpec.Delete(ctx, req.Id); err != nil {
		return nil, err
	}
	return &v1.MenuDeleteRes{}, nil
}

// BatchDelete 批量删除菜单
// 按菜单ID列表批量删除记录，并清理关联关系与缓存
func (s *sMenu) BatchDelete(ctx context.Context, req *v1.MenuBatchDeleteReq) (res *v1.MenuBatchDeleteRes, err error) {
	if err = menuSpec.BatchDelete(ctx, req.Ids); err != nil {
		return nil, err
	}
	return &v1.MenuBatchDeleteRes{}, nil
}
