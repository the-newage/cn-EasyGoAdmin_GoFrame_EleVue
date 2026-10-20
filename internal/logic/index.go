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
	"sort"

	"github.com/gogf/gf/v2/errors/gerror"

	"easygoadmin/api/v1"
	"easygoadmin/internal/consts"
	"easygoadmin/internal/dao"
	"easygoadmin/internal/service"
	"easygoadmin/utility/crypto"
	"easygoadmin/utility/upload"
)

// ============================== 个人中心与菜单聚合业务实现 ==============================

// sIndex 个人中心与菜单聚合业务实现
// 实现 service.IIndex 接口，承载当前用户信息、菜单树与个人资料相关业务
type sIndex struct{}

func init() {
	service.RegisterIndex(&sIndex{})
}

// GetUser 获取当前登录用户信息
//
// 敏感字段永不返回；角色与权限来自身份加载器，
// 超管返回通配权限码，供前端放开全部按钮权限
func (s *sIndex) GetUser(ctx context.Context, req *v1.GetUserReq) (res *v1.GetUserRes, err error) {
	uid := consts.GetUserID(ctx)
	if uid <= 0 {
		return nil, gerror.NewCode(consts.CodeUnauthorized)
	}
	one, err := dao.User.Ctx(ctx).Where("id", uid).One()
	if err != nil {
		return nil, err
	}
	if len(one) == 0 {
		return nil, gerror.NewCode(consts.CodeUserNotFound)
	}

	// 角色编码与权限节点（与 JWT 中间件同源：user_role → role / role_menu → menu.permission）
	roles := consts.GetRoles(ctx)
	perms := consts.GetPerms(ctx)
	if len(roles) == 0 {
		if idn, lerr := LoadIdentity(ctx, uid); lerr == nil && idn != nil {
			roles, perms = idn.Roles, idn.Perms
		}
	}
	if roles == nil {
		roles = []string{}
	}
	// 超管返回通配权限码 *:*:*（前端据此放开全部按钮权限）
	if s.isSuperAdmin(ctx, uid) {
		perms = []string{consts.PermWildcard}
	} else if perms == nil {
		perms = []string{}
	}

	// 城市编码数组（省市区街道，空编码跳过）；city_info 为空时按编码回填名称
	city := make([]string, 0, 4)
	for _, col := range []string{"province_code", "city_code", "district_code", "street_code"} {
		if v := one[col].String(); v != "" {
			city = append(city, v)
		}
	}
	cityInfo := one["city_info"].String()
	if cityInfo == "" && len(city) > 0 {
		cityInfo = ResolveCityNames(ctx, city, ",")
	}
	return &v1.GetUserRes{
		Id:          one["id"].Int64(),
		Username:    one["username"].String(),
		Realname:    one["realname"].String(),
		Gender:      one["gender"].Int(),
		Avatar:      upload.GetFileUrl(ctx, one["avatar"].String()),
		Mobile:      one["mobile"].String(),
		Email:       one["email"].String(),
		DeptId:      one["dept_id"].Int(),
		LevelId:     one["level_id"].Int(),
		PositionId:  one["position_id"].Int(),
		City:        city,
		CityInfo:    cityInfo,
		Address:     one["address"].String(),
		Intro:       one["intro"].String(),
		Status:      one["status"].Int(),
		Note:        one["note"].String(),
		TenantId:    one["tenant_id"].Int(),
		Roles:       roles,
		Permissions: perms,
	}, nil
}

// GetMenus 获取当前用户菜单树
//
// 关联路径：user → user_role → role → role_menu → menu，超管返回全量
func (s *sIndex) GetMenus(ctx context.Context, req *v1.GetMenusReq) (res *v1.GetMenusRes, err error) {
	return s.getMenuTree(ctx)
}

// GetMenusById 获取当前用户菜单树（带路径参数）
//
// 与无路径参数的菜单路由共用同一处理逻辑
func (s *sIndex) GetMenusById(ctx context.Context, req *v1.GetMenusByIdReq) (res *v1.GetMenusRes, err error) {
	return s.getMenuTree(ctx)
}

// getMenuTree 组装菜单树
//
// 仅取目录与菜单节点，按钮权限节点不进菜单树
func (s *sIndex) getMenuTree(ctx context.Context) (*v1.GetMenusRes, error) {
	uid := consts.GetUserID(ctx)
	if uid <= 0 {
		return nil, gerror.NewCode(consts.CodeUnauthorized)
	}

	// 1. 可见菜单集合：超管全量，其余按角色关联
	//    仅取目录/菜单（type=0）、启用（status=0）、非隐藏（hide=0）；按钮权限节点不进菜单树（由 GetUser.permissions 下发）
	query := dao.Menu.Ctx(ctx).Where("type", 0).Where("status", 0).Where("hide", 0)
	if !s.isSuperAdmin(ctx, uid) {
		roleRows, err := dao.UserRole.Ctx(ctx).Where("user_id", uid).All()
		if err != nil {
			return nil, err
		}
		if len(roleRows) == 0 {
			empty := []*v1.MenuTreeItem{}
			return &empty, nil
		}
		roleIds := make([]int64, 0, len(roleRows))
		for _, r := range roleRows {
			roleIds = append(roleIds, r["role_id"].Int64())
		}
		menuRows, err := dao.RoleMenu.Ctx(ctx).WhereIn("role_id", roleIds).All()
		if err != nil {
			return nil, err
		}
		if len(menuRows) == 0 {
			empty := []*v1.MenuTreeItem{}
			return &empty, nil
		}
		menuIds := make([]int64, 0, len(menuRows))
		for _, r := range menuRows {
			menuIds = append(menuIds, r["menu_id"].Int64())
		}
		query = query.WhereIn("id", menuIds)
	}
	rows, err := query.Order("sort, id").All()
	if err != nil {
		return nil, err
	}

	// 2. 平铺节点
	nodes := make(map[int64]*v1.MenuTreeItem, len(rows))
	ids := make([]int64, 0, len(rows))
	for _, row := range rows {
		item := &v1.MenuTreeItem{
			Id:         row["id"].Int64(),
			ParentId:   row["parent_id"].Int(),
			Name:       row["name"].String(),
			Icon:       row["icon"].String(),
			Icon2:      row["icon2"].String(),
			Path:       row["path"].String(),
			Component:  row["component"].String(),
			Permission: row["permission"].String(),
			Target:     row["target"].Int(),
			Type:       row["type"].Int(),
			Status:     row["status"].Int(),
			Hide:       row["hide"].Int(),
			Sort:       row["sort"].Int(),
			Children:   []*v1.MenuTreeItem{},
		}
		nodes[item.Id] = item
		ids = append(ids, item.Id)
	}

	// 3. 挂树：父节点不在集合内（未授权/顶级）一律挂根，避免菜单丢失
	roots := make([]*v1.MenuTreeItem, 0, len(ids))
	for _, id := range ids {
		node := nodes[id]
		if node.ParentId != 0 {
			if parent, ok := nodes[int64(node.ParentId)]; ok {
				parent.Children = append(parent.Children, node)
				continue
			}
		}
		roots = append(roots, node)
	}
	// 父子均已按 sort,id 查询，子级顺序天然有序；根级再兜底一次
	sort.SliceStable(roots, func(i, j int) bool {
		if roots[i].Sort == roots[j].Sort {
			return roots[i].Id < roots[j].Id
		}
		return roots[i].Sort < roots[j].Sort
	})
	return &roots, nil
}

// isSuperAdmin 判断是否超级管理员
// 优先读取上下文角色，未命中时回源用户角色表
func (s *sIndex) isSuperAdmin(ctx context.Context, uid int) bool {
	for _, r := range consts.GetRoles(ctx) {
		if r == consts.RoleCodeSuperAdmin {
			return true
		}
	}
	// 角色缓存未命中时回源
	roleRows, err := dao.UserRole.Ctx(ctx).Where("user_id", uid).All()
	if err != nil || len(roleRows) == 0 {
		return false
	}
	roleIds := make([]int64, 0, len(roleRows))
	for _, r := range roleRows {
		roleIds = append(roleIds, r["role_id"].Int64())
	}
	count, err := dao.Role.Ctx(ctx).WhereIn("id", roleIds).Where("code", consts.RoleCodeSuperAdmin).Count()
	return err == nil && count > 0
}

// UpdatePassword 修改个人密码
//
// 校验旧密码后使用 bcrypt 重写密码与盐
func (s *sIndex) UpdatePassword(ctx context.Context, req *v1.UpdatePasswordReq) (res *v1.UpdatePasswordRes, err error) {
	uid := consts.GetUserID(ctx)
	if uid <= 0 {
		return nil, gerror.NewCode(consts.CodeUnauthorized)
	}
	one, err := dao.User.Ctx(ctx).Where("id", uid).Fields("id, password, salt").One()
	if err != nil {
		return nil, err
	}
	if len(one) == 0 {
		return nil, gerror.NewCode(consts.CodeUserNotFound)
	}
	// 旧密码校验（外部 salt + bcrypt(password+salt)；salt 为空兼容旧数据）
	if !crypto.IsBcrypt(one["password"].String()) || !crypto.VerifyPassword(one["password"].String(), req.Password, one["salt"].String()) {
		return nil, gerror.NewCode(consts.CodeOldPwdWrong)
	}
	salt := crypto.GenerateSalt()
	hashed, err := crypto.HashPassword(req.NewPassword, salt)
	if err != nil {
		return nil, err
	}
	if _, err = dao.User.Ctx(ctx).Where("id", uid).Data(map[string]any{
		"password": hashed,
		"salt":     salt,
	}).Update(); err != nil {
		return nil, err
	}
	return &v1.UpdatePasswordRes{}, nil
}

// UpdateProfile 更新个人信息
//
// 头像执行临时转正，城市信息优先由编码数组解析名称，审计字段由全局钩子盖章
func (s *sIndex) UpdateProfile(ctx context.Context, req *v1.UpdateProfileReq) (res *v1.UpdateProfileRes, err error) {
	uid := consts.GetUserID(ctx)
	if uid <= 0 {
		return nil, gerror.NewCode(consts.CodeUnauthorized)
	}
	data := map[string]any{
		"realname":      req.Realname,
		"gender":        req.Gender,
		"avatar":        SaveFileField(ctx, req.Avatar, "user"), // 头像 temp 转正 + 剥离域名
		"mobile":        req.Mobile,
		"email":         req.Email,
		"city_info":     req.CityInfo,
		"province_code": cityCodeAt(req.City, 0),
		"city_code":     cityCodeAt(req.City, 1),
		"district_code": cityCodeAt(req.City, 2),
		"street_code":   cityCodeAt(req.City, 3),
		"address":       req.Address,
		"intro":         req.Intro,
	}
	// 城市信息：优先由编码数组解析名称
	if len(req.City) > 0 {
		if info := ResolveCityNames(ctx, req.City, ","); info != "" {
			data["city_info"] = info
		}
	}
	if _, err = dao.User.Ctx(ctx).Where("id", uid).Data(data).Update(); err != nil {
		return nil, err
	}
	return &v1.UpdateProfileRes{}, nil
}
