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
	"easygoadmin/utility/cache"
)

// ============================== 角色业务实现 ==============================

// sRole 角色业务实现
// 实现 service.IRole 接口，接口定义见 api/v1/role.go
type sRole struct{}

func init() {
	service.RegisterRole(&sRole{})
}

// roleSerializeMaps 枚举字段到字典编码的映射
var roleSerializeMaps = map[string]string{
	"data_scope": "role_data_scope",
}

// roleSpec 角色模块声明
//
// 编码唯一 + 删除守卫 + 关联清理
var roleSpec = ModuleSpec[v1.RoleItem, v1.RoleListItem, v1.RoleDetailRes]{
	Model:         dao.Role.Ctx,
	LikeFields:    []string{"name", "code"},
	OrderBy:       "sort, id",
	SerializeMaps: roleSerializeMaps,
	UniqueFields:  []UniqueRule{{Field: "code"}},
	FieldLabels:   map[string]string{"code": "角色编码"}, // 唯一提示文案
	ToData:        roleToData,
	AfterUpdate:   roleClearCache,     // 身份缓存失效
	BeforeDelete:  roleDeleteGuard,    // 存在性 + 关联用户禁止
	AfterDelete:   roleCleanRelations, // 清理 role_menu
	AfterBatchDelete: func(ctx context.Context, ids []int64) error {
		return roleCleanRelationsBatch(ctx, ids)
	},
	AfterPageList: roleFillUserCount, // 行内用户数批量补全
}

// roleToData 全量表单语义转换
//
// 修改请求不包含编码，因编码不允许变更；tenant_id 由 hooks 自动填充；
// 数据范围缺省归一为"全部"（1）
func roleToData(req any) map[string]any {
	switch r := req.(type) {
	case *v1.RoleAddReq:
		scope := r.DataScope
		if scope <= 0 {
			scope = 1
		}
		return map[string]any{"name": r.Name, "code": r.Code, "data_scope": scope, "sort": r.Sort, "note": r.Note}
	case *v1.RoleUpdateReq:
		scope := r.DataScope
		if scope <= 0 {
			scope = 1
		}
		return map[string]any{"name": r.Name, "data_scope": scope, "sort": r.Sort, "note": r.Note}
	}
	panic("roleToData: 不支持的入参类型")
}

// roleFillUserCount 分页行用户数批量补全
// 依据 user_role 关联统计
func roleFillUserCount(ctx context.Context, list []*v1.RoleItem, req any) error {
	counts, err := roleUserCounts(ctx)
	if err != nil {
		return err
	}
	for _, item := range list {
		item.UserCount = counts[item.Id]
	}
	return nil
}

// roleClearCache 角色变更后清理身份缓存
// 失效该角色下所有用户的身份缓存
func roleClearCache(ctx context.Context, id int64, req any) error {
	roleClearIdentityCache(ctx, id)
	return nil
}

// roleDeleteGuard 删除守卫
// 校验记录存在性，存在关联用户时禁止删除
func roleDeleteGuard(ctx context.Context, id int64) error {
	if n, err := dao.Role.Ctx(ctx).Where("id", id).Count(); err != nil {
		return err
	} else if n == 0 {
		return gerror.NewCode(consts.CodeRecordNotFound)
	}
	if n, err := dao.UserRole.Ctx(ctx).Where("role_id", id).Count(); err != nil {
		return err
	} else if n > 0 {
		return gerror.NewCode(consts.CodeInUse)
	}
	return nil
}

// roleCleanRelations 单条删除后清理角色菜单关联
func roleCleanRelations(ctx context.Context, id int64) error {
	_, err := dao.RoleMenu.Ctx(ctx).Where("role_id", id).Delete()
	return err
}

// roleCleanRelationsBatch 批量删除后清理角色菜单关联
func roleCleanRelationsBatch(ctx context.Context, ids []int64) error {
	_, err := dao.RoleMenu.Ctx(ctx).WhereIn("role_id", ids).Delete()
	return err
}

// ---------- 标准方法：3行委托 ----------

// Page 角色分页
// 分页返回角色列表，并补充关联用户数
func (s *sRole) Page(ctx context.Context, req *v1.RolePageReq) (res *v1.RolePageRes, err error) {
	list, total, err := roleSpec.Page(ctx, req, req.GetPageNo(), req.GetPageSize())
	if err != nil {
		return nil, err
	}
	return v1.NewPageData(req.GetPageNo(), req.GetPageSize(), total, derefSlice(list)), nil
}

// List 角色全量列表
// 按关键词返回角色精简列表
func (s *sRole) List(ctx context.Context, req *v1.RoleListReq) (res v1.RoleListRes, err error) {
	list, err := roleSpec.List(ctx, req)
	if err != nil {
		return nil, err
	}
	return derefSlice(list), nil
}

// Detail 角色详情
// 按角色ID查询详情
func (s *sRole) Detail(ctx context.Context, req *v1.RoleDetailReq) (res *v1.RoleDetailRes, err error) {
	return roleSpec.Detail(ctx, req.Id, req)
}

// Add 新增角色
// 校验角色编码唯一后新增，返回新记录ID
func (s *sRole) Add(ctx context.Context, req *v1.RoleAddReq) (res *v1.IdRes, err error) {
	id, err := roleSpec.Add(ctx, req)
	if err != nil {
		return nil, err
	}
	return &v1.IdRes{Id: id}, nil
}

// Update 修改角色
// 按角色ID更新记录，编码不允许变更
func (s *sRole) Update(ctx context.Context, req *v1.RoleUpdateReq) (res *v1.RoleUpdateRes, err error) {
	if err = roleSpec.Update(ctx, req.Id, req); err != nil {
		return nil, err
	}
	return &v1.RoleUpdateRes{}, nil
}

// Delete 删除角色
// 按角色ID删除单条记录，存在关联用户时禁止
func (s *sRole) Delete(ctx context.Context, req *v1.RoleDeleteReq) (res *v1.RoleDeleteRes, err error) {
	if err = roleSpec.Delete(ctx, req.Id); err != nil {
		return nil, err
	}
	return &v1.RoleDeleteRes{}, nil
}

// BatchDelete 批量删除角色
// 按角色ID列表批量删除记录
func (s *sRole) BatchDelete(ctx context.Context, req *v1.RoleBatchDeleteReq) (res *v1.RoleBatchDeleteRes, err error) {
	if err = roleSpec.BatchDelete(ctx, req.Ids); err != nil {
		return nil, err
	}
	return &v1.RoleBatchDeleteRes{}, nil
}

// ---------- 非标准方法（手写实现） ----------

// MenuList 角色菜单勾选列表
// 返回菜单平铺列表并附带勾选状态，供权限配置使用
func (s *sRole) MenuList(ctx context.Context, req *v1.RoleMenuListReq) (res v1.RoleMenuListRes, err error) {
	roleMenus, err := dao.RoleMenu.Ctx(ctx).Where("role_id", req.RoleId).All()
	if err != nil {
		return nil, err
	}
	checkedSet := make(map[int64]bool, len(roleMenus))
	for _, r := range roleMenus {
		checkedSet[r["menu_id"].Int64()] = true
	}
	rows, err := dao.Menu.Ctx(ctx).Order("sort, id").All()
	if err != nil {
		return nil, err
	}
	list := make([]v1.RoleMenuItem, 0, len(rows))
	for _, r := range rows {
		list = append(list, v1.RoleMenuItem{
			Id:         r["id"].Int64(),
			ParentId:   r["parent_id"].Int(),
			Name:       r["name"].String(),
			Icon:       r["icon"].String(),
			Icon2:      r["icon2"].String(),
			Path:       r["path"].String(),
			Component:  r["component"].String(),
			Permission: r["permission"].String(),
			Target:     r["target"].Int(),
			Type:       r["type"].Int(),
			Status:     r["status"].Int(),
			Hide:       r["hide"].Int(),
			Sort:       r["sort"].Int(),
			Checked:    checkedSet[r["id"].Int64()],
		})
	}
	return list, nil
}

// SavePermission 保存角色菜单权限
//
// 全量替换角色菜单关联，并失效相关用户的身份缓存
func (s *sRole) SavePermission(ctx context.Context, req *v1.RoleSavePermissionReq) (res *v1.RoleSavePermissionRes, err error) {
	if _, err = dao.RoleMenu.Ctx(ctx).Where("role_id", req.RoleId).Delete(); err != nil {
		return nil, err
	}
	for _, mid := range req.MenuIds {
		if mid <= 0 {
			continue
		}
		if _, err = dao.RoleMenu.Ctx(ctx).Data(map[string]any{
			"role_id": req.RoleId,
			"menu_id": mid,
		}).Insert(); err != nil {
			return nil, err
		}
	}
	s.clearRoleIdentityCache(ctx, req.RoleId)
	return &v1.RoleSavePermissionRes{}, nil
}

// ---------- 私有辅助 ----------

// roleUserCounts 统计角色ID到用户数的映射
func roleUserCounts(ctx context.Context) (map[int64]int, error) {
	result := make(map[int64]int)
	rows, err := dao.UserRole.Ctx(ctx).All()
	if err != nil {
		return nil, err
	}
	for _, r := range rows {
		result[r["role_id"].Int64()]++
	}
	return result, nil
}

// roleClearIdentityCache 失效角色下所有用户的身份缓存
func roleClearIdentityCache(ctx context.Context, roleId int64) {
	urRows, err := dao.UserRole.Ctx(ctx).Where("role_id", roleId).All()
	if err != nil {
		return
	}
	for _, r := range urRows {
		uid := int(r["user_id"].Int64())
		cache.Del(ctx,
			consts.CacheKeyUserIdentity(uid),
			consts.CacheKeyUserPerms(uid),
			consts.CacheKeyUserRoles(uid),
		)
	}
}

// clearRoleIdentityCache 身份缓存失效（方法版）
// 供 SavePermission 调用
func (s *sRole) clearRoleIdentityCache(ctx context.Context, roleId int64) {
	roleClearIdentityCache(ctx, roleId)
}
