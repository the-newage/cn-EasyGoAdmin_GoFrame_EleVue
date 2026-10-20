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
	"fmt"
	"strconv"
	"strings"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"

	"easygoadmin/api/v1"
	"easygoadmin/internal/consts"
	"easygoadmin/internal/dao"
	"easygoadmin/internal/service"
	"easygoadmin/utility/cache"
	"easygoadmin/utility/crypto"
	"easygoadmin/utility/excel"
	"easygoadmin/utility/upload"
)

// ============================== 用户业务实现 ==============================

// sUser 用户业务实现
// 实现 service.IUser 接口，承载用户模块的业务逻辑
type sUser struct{}

func init() {
	service.RegisterUser(&sUser{})
}

// userSerializeMaps 枚举字段到字典编码的映射
// {field}Text 由 dict_item 解析，字典管理页可修改文案
var userSerializeMaps = map[string]string{
	"gender": "user_gender",
	"status": "user_status",
}

// userSpec 用户模块声明
//
// Add、Update、Detail 涉及密码、城市、头像与角色装配，采用手写实现
var userSpec = ModuleSpec[v1.UserItem, v1.UserListItem, v1.UserDetailRes]{
	Model:         dao.User.Ctx,
	LikeFields:    []string{"realname"}, // querySchemas 姓名输入 → 模糊
	EqFields:      []string{"status", "dept_id", "level_id", "position_id", "tenant_id"},
	OrderBy:       "id DESC",
	ListOrderBy:   "sort, id",
	SerializeMaps: userSerializeMaps,
	FileFields:    []string{"avatar"}, // 读侧自动补域名（file.domainUrl）
	FileSaveDir:   "user",
	Filter:        userFilter,                // keyword 列 Page/List 有别 + role_id 子查询
	DataScope:     ApplyDataScope("dept_id"), // 角色数据权限（Page/List/导出生效）
	BeforeDelete:  userDeleteGuard,
	AfterDelete:   userCleanRelations, // 清理 user_role + 身份缓存
	AfterBatchDelete: func(ctx context.Context, ids []int64) error {
		return userCleanRelationsBatch(ctx, ids)
	},
	AfterPageList: userExpandNames, // 行内角色/部门/职级/岗位名批量补全
}

// userFilter 复杂筛选
//
// 关键词按入参分支，分页请求额外匹配手机号，并支持按角色ID过滤
func userFilter(ctx context.Context, m *gdb.Model, req any) *gdb.Model {
	switch r := req.(type) {
	case *v1.UserPageReq:
		if r.Keyword != "" {
			kw := "%" + r.Keyword + "%"
			m = m.Where("(username LIKE ? OR realname LIKE ? OR mobile LIKE ?)", kw, kw, kw)
		}
		roleId := r.RoleId
		if roleId == 0 {
			roleId = r.Role // 前端传参名兼容（role）
		}
		if roleId > 0 {
			m = m.WhereIn("id", dao.UserRole.Ctx(ctx).Where("role_id", roleId).Fields("user_id"))
		}
	case *v1.UserListReq:
		if r.Keyword != "" {
			kw := "%" + r.Keyword + "%"
			m = m.Where("(username LIKE ? OR realname LIKE ?)", kw, kw)
		}
	}
	return m
}

// userDeleteGuard 删除守卫
// 校验记录存在性
func userDeleteGuard(ctx context.Context, id int64) error {
	if n, err := dao.User.Ctx(ctx).Where("id", id).Count(); err != nil {
		return err
	} else if n == 0 {
		return gerror.NewCode(consts.CodeRecordNotFound)
	}
	return nil
}

// userCleanRelations 单条删除后清理
// 清理 user_role 关联并失效身份缓存
func userCleanRelations(ctx context.Context, id int64) error {
	return userCleanRelationsBatch(ctx, []int64{id})
}

// userCleanRelationsBatch 批量删除后清理
// 清理 user_role 关联并失效身份缓存
func userCleanRelationsBatch(ctx context.Context, ids []int64) error {
	if _, err := dao.UserRole.Ctx(ctx).WhereIn("user_id", ids).Delete(); err != nil {
		return err
	}
	for _, id := range ids {
		cache.Del(ctx,
			consts.CacheKeyUserIdentity(int(id)),
			consts.CacheKeyUserPerms(int(id)),
			consts.CacheKeyUserRoles(int(id)),
		)
	}
	return nil
}

// userExpandNames 分页行冗余装配
//
// 补全角色引用数组、城市编码数组、部门、职级、岗位名称及头像域名
func userExpandNames(ctx context.Context, list []*v1.UserItem, req any) error {
	if len(list) == 0 {
		return nil
	}
	ids := make([]int64, 0, len(list))
	for _, item := range list {
		ids = append(ids, item.Id)
	}
	// UserItem 无 dept/level/position/city id 列 → 按主键回查装配所需外键
	rows, err := dao.User.Ctx(ctx).WhereIn("id", ids).
		Fields("id, dept_id, level_id, position_id, province_code, city_code, district_code, street_code").All()
	if err != nil {
		return err
	}
	u := &sUser{}
	levelIds := make([]int64, 0)
	positionIds := make([]int64, 0)
	deptIds := make([]int64, 0)
	fkOf := make(map[int64]gdb.Record, len(rows))
	for _, r := range rows {
		fkOf[r["id"].Int64()] = r
		if v := r["level_id"].Int64(); v > 0 {
			levelIds = append(levelIds, v)
		}
		if v := r["position_id"].Int64(); v > 0 {
			positionIds = append(positionIds, v)
		}
		if v := r["dept_id"].Int64(); v > 0 {
			deptIds = append(deptIds, v)
		}
	}
	nameOfLevel := u.namesById(ctx, dao.Level.Ctx(ctx), levelIds)
	nameOfPosition := u.namesById(ctx, dao.Position.Ctx(ctx), positionIds)
	nameOfDept := u.namesById(ctx, dao.Dept.Ctx(ctx), deptIds)
	roleRefMap, err := u.userRoleRefs(ctx, ids)
	if err != nil {
		return err
	}
	for _, item := range list {
		item.Roles = roleRefMap[item.Id]
		if item.Roles == nil {
			item.Roles = []v1.UserRoleRef{}
		}
		if fk, ok := fkOf[item.Id]; ok {
			item.LevelName = nameOfLevel[fk["level_id"].Int64()]
			item.PositionName = nameOfPosition[fk["position_id"].Int64()]
			item.DeptName = nameOfDept[fk["dept_id"].Int64()]
			city := make([]string, 0, 4)
			for _, col := range []string{"province_code", "city_code", "district_code", "street_code"} {
				if v := fk[col].String(); v != "" {
					city = append(city, v)
				}
			}
			item.City = city
		}
		if item.City == nil {
			item.City = []string{}
		}
	}
	return nil
}

// Page 用户分页
//
// 支持关键词、状态、部门、角色、职级、岗位等条件，行内冗余名称批量补全
func (s *sUser) Page(ctx context.Context, req *v1.UserPageReq) (res *v1.UserPageRes, err error) {
	list, total, err := userSpec.Page(ctx, req, req.GetPageNo(), req.GetPageSize())
	if err != nil {
		return nil, err
	}
	return v1.NewPageData(req.GetPageNo(), req.GetPageSize(), total, derefSlice(list)), nil
}

// List 用户全量列表
//
// 供下拉选择使用，角色名在委托层后处理
func (s *sUser) List(ctx context.Context, req *v1.UserListReq) (res v1.UserListRes, err error) {
	list, err := userSpec.List(ctx, req)
	if err != nil {
		return nil, err
	}
	// 精简列表不挂基座钩子 → 角色名在委托层后处理
	roleMap, err := s.userRoleNames(ctx)
	if err != nil {
		return nil, err
	}
	for _, item := range list {
		item.Role = roleMap[item.Id]
	}
	return derefSlice(list), nil
}

// Detail 用户详情
//
// 含角色ID数组与四级区划编码，城市信息为空时按编码回填名称
func (s *sUser) Detail(ctx context.Context, req *v1.UserDetailReq) (res *v1.UserDetailRes, err error) {
	one, err := dao.User.Ctx(ctx).Where("id", req.Id).One()
	if err != nil {
		return nil, err
	}
	if len(one) == 0 {
		return nil, gerror.NewCode(consts.CodeRecordNotFound)
	}
	roleIds, err := s.userRoleIds(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	detail := s.toDetail(ctx, one, roleIds)
	if detail.CityInfo == "" && len(detail.City) > 0 {
		detail.CityInfo = ResolveCityNames(ctx, detail.City, ",")
	}
	return detail, nil
}

// Add 新增用户
//
// 用户名与手机号租户内唯一，校验租户用户上限，自动生成用户编码并写入 user_role
func (s *sUser) Add(ctx context.Context, req *v1.UserAddReq) (res *v1.IdRes, err error) {
	tenantId := consts.GetTenantID(ctx)
	if exists, e := s.usernameExists(ctx, req.Username, tenantId, 0); e != nil {
		return nil, e
	} else if exists {
		return nil, gerror.NewCodef(consts.CodeRecordExists, "用户名已存在")
	}
	if req.Mobile != "" {
		if exists, e := s.mobileExists(ctx, req.Mobile, tenantId, 0); e != nil {
			return nil, e
		} else if exists {
			return nil, gerror.NewCodef(consts.CodeRecordExists, "手机号已被使用")
		}
	}
	if err = s.checkTenantUserMax(ctx, tenantId); err != nil {
		return nil, err
	}
	// 自动生成用户编码：{租户编码}-{租户内6位序号}（如 T0001-000020；平台租户用 SYS）
	code, err := genUserCode(ctx, tenantId)
	if err != nil {
		return nil, err
	}
	// 双重加盐：外部 salt + bcrypt(password+salt)，两值分别入 password/salt 列
	salt := crypto.GenerateSalt()
	hashed, err := crypto.HashPassword(req.Password, salt)
	if err != nil {
		return nil, err
	}
	status := req.Status
	if status == 0 {
		status = 1
	}
	data := map[string]any{
		"tenant_id":     tenantId,
		"code":          code,
		"username":      req.Username,
		"password":      hashed,
		"salt":          salt,
		"realname":      req.Realname,
		"gender":        req.Gender,
		"mobile":        req.Mobile,
		"email":         req.Email,
		"dept_id":       req.DeptId,
		"level_id":      req.LevelId,
		"position_id":   req.PositionId,
		"city_info":     req.CityInfo,
		"address":       req.Address,
		"avatar":        req.Avatar,
		"intro":         req.Intro,
		"status":        status,
		"sort":          req.Sort,
		"note":          req.Note,
		"province_code": cityCodeAt(req.City, 0),
		"city_code":     cityCodeAt(req.City, 1),
		"district_code": cityCodeAt(req.City, 2),
		"street_code":   cityCodeAt(req.City, 3),
	}
	// 头像 temp → 正式目录（入库形态为相对 uploadDir 的路径，剥离域名）
	data["avatar"] = SaveFileField(ctx, req.Avatar, "user")
	// 城市信息：优先由编码数组解析名称
	if len(req.City) > 0 {
		if info := ResolveCityNames(ctx, req.City, ","); info != "" {
			data["city_info"] = info
		}
	}
	result, err := dao.User.Ctx(ctx).Data(data).Insert()
	if err != nil {
		return nil, err
	}
	idv, _ := result.LastInsertId()
	if err = s.AssignRoles(ctx, idv, req.Roles); err != nil {
		return nil, err
	}
	return &v1.IdRes{Id: idv}, nil
}

// Update 修改用户
//
// 密码留空时不做修改，角色关联全量重写
func (s *sUser) Update(ctx context.Context, req *v1.UserUpdateReq) (res *v1.UserUpdateRes, err error) {
	if exists, e := s.usernameExists(ctx, req.Username, consts.GetTenantID(ctx), req.Id); e != nil {
		return nil, e
	} else if exists {
		return nil, gerror.NewCodef(consts.CodeRecordExists, "用户名已存在")
	}
	if req.Mobile != "" {
		if exists, e := s.mobileExists(ctx, req.Mobile, consts.GetTenantID(ctx), req.Id); e != nil {
			return nil, e
		} else if exists {
			return nil, gerror.NewCodef(consts.CodeRecordExists, "手机号已被使用")
		}
	}
	data := map[string]any{
		"username":      req.Username,
		"realname":      req.Realname,
		"gender":        req.Gender,
		"mobile":        req.Mobile,
		"email":         req.Email,
		"dept_id":       req.DeptId,
		"level_id":      req.LevelId,
		"position_id":   req.PositionId,
		"city_info":     req.CityInfo,
		"address":       req.Address,
		"avatar":        req.Avatar,
		"intro":         req.Intro,
		"status":        req.Status,
		"sort":          req.Sort,
		"note":          req.Note,
		"province_code": cityCodeAt(req.City, 0),
		"city_code":     cityCodeAt(req.City, 1),
		"district_code": cityCodeAt(req.City, 2),
		"street_code":   cityCodeAt(req.City, 3),
	}
	if len(req.City) > 0 {
		if info := ResolveCityNames(ctx, req.City, ","); info != "" {
			data["city_info"] = info
		}
	}
	data["avatar"] = SaveFileField(ctx, req.Avatar, "user")
	if req.Password != "" {
		salt := crypto.GenerateSalt()
		hashed, e := crypto.HashPassword(req.Password, salt)
		if e != nil {
			return nil, e
		}
		data["password"] = hashed
		data["salt"] = salt
	}
	if _, err = dao.User.Ctx(ctx).Where("id", req.Id).Data(data).Update(); err != nil {
		return nil, err
	}
	if err = s.AssignRoles(ctx, req.Id, req.Roles); err != nil {
		return nil, err
	}
	return &v1.UserUpdateRes{}, nil
}

// Status 用户状态
// 设置用户的正常或禁用状态
func (s *sUser) Status(ctx context.Context, req *v1.UserStatusReq) (res *v1.UserStatusRes, err error) {
	if _, err = dao.User.Ctx(ctx).Where("id", req.Id).Data(map[string]any{"status": req.Status}).Update(); err != nil {
		return nil, err
	}
	return &v1.UserStatusRes{}, nil
}

// ResetPwd 重置用户密码
// 未传新密码时重置为系统默认密码
func (s *sUser) ResetPwd(ctx context.Context, req *v1.UserResetPwdReq) (res *v1.UserResetPwdRes, err error) {
	if n, e := dao.User.Ctx(ctx).Where("id", req.UserId).Count(); e != nil {
		return nil, e
	} else if n == 0 {
		return nil, gerror.New("用户不存在")
	}
	password := req.Password
	if password == "" {
		password = getDefaultPassword(ctx)
	}
	salt := crypto.GenerateSalt()
	hashed, err := crypto.HashPassword(password, salt)
	if err != nil {
		return nil, err
	}
	if _, err = dao.User.Ctx(ctx).Where("id", req.UserId).Data(map[string]any{
		"password": hashed,
		"salt":     salt,
	}).Update(); err != nil {
		return nil, err
	}
	return &v1.UserResetPwdRes{}, nil
}

// Delete 删除用户
// 软删记录并清理 user_role 关联
func (s *sUser) Delete(ctx context.Context, req *v1.UserDeleteReq) (res *v1.UserDeleteRes, err error) {
	if err = userSpec.Delete(ctx, req.Id); err != nil {
		return nil, err
	}
	return &v1.UserDeleteRes{}, nil
}

// BatchDelete 批量删除用户
// 按用户ID列表批量删除记录
func (s *sUser) BatchDelete(ctx context.Context, req *v1.UserBatchDeleteReq) (res *v1.UserBatchDeleteRes, err error) {
	if err = userSpec.BatchDelete(ctx, req.Ids); err != nil {
		return nil, err
	}
	return &v1.UserBatchDeleteRes{}, nil
}

// Import 导入用户
//
// 表格列序：姓名、性别、手机号、邮箱、部门、职级、岗位、城市、地址、用户名、简介、状态、备注；
// 密码统一取系统默认密码，用户编码自动生成；
// 性别与状态支持文本或数字，部门、职级、岗位支持名称或ID，城市为逗号分隔的名称串
func (s *sUser) Import(ctx context.Context, req *v1.UserImportReq) (res *v1.UserImportRes, err error) {
	_, rows, err := excel.ReadSheet(req.File)
	if err != nil {
		return nil, err
	}
	res = &v1.UserImportRes{Total: len(rows)}
	tenantId := consts.GetTenantID(ctx)
	for i, row := range rows {
		username := cellAt(row, 9)
		if username == "" {
			res.Failed++
			res.Reasons = append(res.Reasons, fmt.Sprintf("第%d行：用户名为空", i+2))
			continue
		}
		if exists, e := s.usernameExists(ctx, username, tenantId, 0); e != nil {
			return nil, e
		} else if exists {
			res.Failed++
			res.Reasons = append(res.Reasons, fmt.Sprintf("第%d行：用户名已存在", i+2))
			continue
		}
		salt := crypto.GenerateSalt()
		hashed, e := crypto.HashPassword(getDefaultPassword(ctx), salt)
		if e != nil {
			res.Failed++
			res.Reasons = append(res.Reasons, fmt.Sprintf("第%d行：密码加密失败", i+2))
			continue
		}
		code, e := genUserCode(ctx, tenantId)
		if e != nil {
			return nil, e
		}
		cityInfo := cellAt(row, 7)
		codes := ResolveCityCodes(ctx, cityInfo)
		if _, e = dao.User.Ctx(ctx).Data(map[string]any{
			"username":      username,
			"realname":      cellAt(row, 0),
			"code":          code,
			"password":      hashed,
			"salt":          salt,
			"gender":        parseGenderCell(cellAt(row, 1)),
			"mobile":        cellAt(row, 2),
			"email":         cellAt(row, 3),
			"dept_id":       resolveRefID(ctx, dao.Dept.Ctx(ctx), cellAt(row, 4)),
			"level_id":      resolveRefID(ctx, dao.Level.Ctx(ctx), cellAt(row, 5)),
			"position_id":   resolveRefID(ctx, dao.Position.Ctx(ctx), cellAt(row, 6)),
			"city_info":     cityInfo,
			"address":       cellAt(row, 8),
			"intro":         cellAt(row, 10),
			"status":        parseStatusCell(cellAt(row, 11)),
			"note":          cellAt(row, 12),
			"tenant_id":     tenantId,
			"province_code": codes[0],
			"city_code":     codes[1],
			"district_code": codes[2],
			"street_code":   codes[3],
		}).Insert(); e != nil {
			res.Failed++
			res.Reasons = append(res.Reasons, fmt.Sprintf("第%d行：写入失败 %v", i+2, e))
			continue
		}
		res.Success++
	}
	return res, nil
}

// Export 导出用户
//
// 列序对齐导出模板，每批 500 条防 OOM，写出 xlsx 并返回访问地址
func (s *sUser) Export(ctx context.Context, req *v1.UserExportReq) (res *v1.UserExportRes, err error) {
	// 部门/职级/岗位 ID→名称（小维表全量加载，行内引用）
	deptNames, err := idNameMap(ctx, dao.Dept.Ctx(ctx))
	if err != nil {
		return nil, err
	}
	levelNames, err := idNameMap(ctx, dao.Level.Ctx(ctx))
	if err != nil {
		return nil, err
	}
	positionNames, err := idNameMap(ctx, dao.Position.Ctx(ctx))
	if err != nil {
		return nil, err
	}
	headers := []string{"用户编码", "姓名", "性别", "手机号", "邮箱", "部门", "职级", "岗位", "城市", "地址", "用户名", "简介", "状态", "备注"}
	var rows [][]any
	// 分块遍历导出（每页重建查询防 OOM；并发写入下可能有翻页漂移，导出场景可接受）
	if err = ChunkQuery(ctx, func(page, size int) *gdb.Model {
		m := userSpec.applyFilters(ctx, dao.User.Ctx(ctx), req)
		return m.Page(page, size).Order("id DESC")
	}, 500, func(part gdb.Result) error {
		for _, r := range part {
			codes := make([]string, 0, 4)
			for _, col := range []string{"province_code", "city_code", "district_code", "street_code"} {
				if v := r[col].String(); v != "" {
					codes = append(codes, v)
				}
			}
			// 城市优先由四级编码解析名称，编码缺失时回退 city_info 存量值
			city := ResolveCityNames(ctx, codes, ",")
			if city == "" {
				city = r["city_info"].String()
			}
			rows = append(rows, []any{
				r["code"].String(), r["realname"].String(), genderText(r["gender"].Int()),
				r["mobile"].String(), r["email"].String(),
				deptNames[r["dept_id"].Int64()], levelNames[r["level_id"].Int64()],
				positionNames[r["position_id"].Int64()], city, r["address"].String(),
				r["username"].String(), r["intro"].String(),
				userStatusText(r["status"].Int()), r["note"].String(),
			})
		}
		return nil
	}); err != nil {
		return nil, err
	}
	fileName, _, err := excel.WriteSheet(upload.BaseDir(ctx), "用户数据", headers, rows)
	if err != nil {
		return nil, err
	}
	url := upload.GetFileUrl(ctx, fileName)
	return &url, nil
}

// Document 用户打印信息
//
// 在用户详情基础上追加角色、部门、职级、岗位名称
func (s *sUser) Document(ctx context.Context, req *v1.UserDocumentReq) (res *v1.UserDocumentRes, err error) {
	one, err := dao.User.Ctx(ctx).Where("id", req.Id).One()
	if err != nil {
		return nil, err
	}
	if len(one) == 0 {
		return nil, gerror.NewCode(consts.CodeRecordNotFound)
	}
	roleIds, err := s.userRoleIds(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	detail := s.toDetail(ctx, one, roleIds)
	res = &v1.UserDocumentRes{UserDetailRes: *detail, RoleNames: []string{}}
	if len(roleIds) > 0 {
		roleRows, e := dao.Role.Ctx(ctx).WhereIn("id", roleIds).All()
		if e == nil {
			for _, r := range roleRows {
				res.RoleNames = append(res.RoleNames, r["name"].String())
			}
		}
	}
	res.DeptName = s.dictName(ctx, dao.Dept.Ctx(ctx), one["dept_id"].Int())
	res.LevelName = s.dictName(ctx, dao.Level.Ctx(ctx), one["level_id"].Int())
	res.PositionName = s.dictName(ctx, dao.Position.Ctx(ctx), one["position_id"].Int())
	return res, nil
}

// AssignRoles 重写用户角色关联
//
// 供新增与修改调用，完成后失效该用户的身份缓存
func (s *sUser) AssignRoles(ctx context.Context, userId int64, roleIds []int64) error {
	if _, err := dao.UserRole.Ctx(ctx).Where("user_id", userId).Delete(); err != nil {
		return err
	}
	for _, rid := range roleIds {
		if rid <= 0 {
			continue
		}
		if _, err := dao.UserRole.Ctx(ctx).Data(map[string]any{
			"user_id": userId,
			"role_id": rid,
		}).Insert(); err != nil {
			return err
		}
	}
	s.clearIdentityCache(ctx, userId)
	return nil
}

// ---------- 私有辅助 ----------

// clearIdentityCache 权限相关缓存失效
// 采用尽力而为策略
func (s *sUser) clearIdentityCache(ctx context.Context, userId int64) {
	cache.Del(ctx,
		consts.CacheKeyUserIdentity(int(userId)),
		consts.CacheKeyUserPerms(int(userId)),
		consts.CacheKeyUserRoles(int(userId)),
	)
}

// usernameExists 判断用户名是否已存在
// tenantId 限定租户，excludeId 排除自身
func (s *sUser) usernameExists(ctx context.Context, username string, tenantId int, excludeId int64) (bool, error) {
	m := dao.User.Ctx(ctx).Where("username", username).Where("tenant_id", tenantId)
	if excludeId > 0 {
		m = m.WhereNot("id", excludeId)
	}
	count, err := m.Count()
	return count > 0, err
}

// mobileExists 判断手机号是否已存在
func (s *sUser) mobileExists(ctx context.Context, mobile string, tenantId int, excludeId int64) (bool, error) {
	m := dao.User.Ctx(ctx).Where("mobile", mobile).Where("tenant_id", tenantId)
	if excludeId > 0 {
		m = m.WhereNot("id", excludeId)
	}
	count, err := m.Count()
	return count > 0, err
}

// checkTenantUserMax 校验租户用户数上限
func (s *sUser) checkTenantUserMax(ctx context.Context, tenantId int) error {
	if tenantId == 0 {
		return nil
	}
	one, err := dao.Tenant.Ctx(ctx).Where("id", tenantId).One()
	if err != nil || len(one) == 0 {
		return nil
	}
	maxNum := one["max_user_num"].Int()
	if maxNum <= 0 {
		return nil
	}
	count, err := dao.User.Ctx(ctx).Where("tenant_id", tenantId).Count()
	if err != nil {
		return err
	}
	if count >= maxNum {
		return gerror.NewCode(consts.CodeTenantUserLimit)
	}
	return nil
}

// genUserCode 生成用户编码
//
// 格式为 {租户编码}-{租户内6位序号}，如 T0001-000020；平台租户前缀为 SYS。
// 序号取租户内未删除用户数加一，与现有编码顺延
func genUserCode(ctx context.Context, tenantId int) (string, error) {
	prefix := "SYS"
	if tenantId > 0 {
		one, err := dao.Tenant.Ctx(ctx).Where("id", tenantId).Fields("code").One()
		if err != nil {
			return "", err
		}
		if len(one) == 0 {
			return "", gerror.NewCode(consts.CodeTenantNotFound)
		}
		prefix = one["code"].String()
	}
	n, err := dao.User.Ctx(ctx).Where("tenant_id", tenantId).Count()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s-%06d", prefix, n+1), nil
}

// userRoleIds 查询用户的角色ID列表
func (s *sUser) userRoleIds(ctx context.Context, userId int64) ([]int64, error) {
	rows, err := dao.UserRole.Ctx(ctx).Where("user_id", userId).All()
	if err != nil {
		return nil, err
	}
	ids := make([]int64, 0, len(rows))
	for _, r := range rows {
		ids = append(ids, r["role_id"].Int64())
	}
	return ids, nil
}

// userRoleRefs 批量查询用户角色引用数组
//
// 返回用户ID到角色引用的映射，一次查询避免 N+1
func (s *sUser) userRoleRefs(ctx context.Context, userIds []int64) (map[int64][]v1.UserRoleRef, error) {
	result := make(map[int64][]v1.UserRoleRef)
	if len(userIds) == 0 {
		return result, nil
	}
	// 注意：不可给主表起别名——dao.Ctx 注入的软删条件使用完整表名，带别名时 MySQL 无法解析
	urTable := dao.UserRole.Table()
	roleData, err := dao.UserRole.Ctx(ctx).
		LeftJoin(dao.Role.Table()+" r", urTable+".role_id = r.id AND r.del_flag=0").
		WhereIn(urTable+".user_id", userIds).
		Fields(urTable + ".user_id, r.id, r.name").
		Order(urTable + ".role_id").
		All()
	if err != nil {
		return nil, err
	}
	for _, rd := range roleData {
		uid := rd["user_id"].Int64()
		result[uid] = append(result[uid], v1.UserRoleRef{
			Id:   rd["id"].Int64(),
			Name: rd["name"].String(),
		})
	}
	return result, nil
}

// userRoleNames 查询用户ID到角色名的映射
// 多个角色名以逗号分隔
func (s *sUser) userRoleNames(ctx context.Context) (map[int64]string, error) {
	result := make(map[int64]string)
	urRows, err := dao.UserRole.Ctx(ctx).All()
	if err != nil {
		return nil, err
	}
	if len(urRows) == 0 {
		return result, nil
	}
	byUser := make(map[int64][]int64)
	roleIds := make([]int64, 0, len(urRows))
	for _, r := range urRows {
		uid, rid := r["user_id"].Int64(), r["role_id"].Int64()
		byUser[uid] = append(byUser[uid], rid)
		roleIds = append(roleIds, rid)
	}
	nameById := make(map[int64]string)
	roleRows, err := dao.Role.Ctx(ctx).WhereIn("id", roleIds).All()
	if err != nil {
		return nil, err
	}
	for _, r := range roleRows {
		nameById[r["id"].Int64()] = r["name"].String()
	}
	for uid, ids := range byUser {
		names := ""
		for i, rid := range ids {
			if i > 0 {
				names += ","
			}
			names += nameById[rid]
		}
		result[uid] = names
	}
	return result, nil
}

// toDetail 数据库记录转换为用户详情
func (s *sUser) toDetail(ctx context.Context, row gdb.Record, roleIds []int64) *v1.UserDetailRes {
	if roleIds == nil {
		roleIds = []int64{}
	}
	city := make([]string, 0, 4)
	for _, col := range []string{"province_code", "city_code", "district_code", "street_code"} {
		if v := row[col].String(); v != "" {
			city = append(city, v)
		}
	}
	return &v1.UserDetailRes{
		Id:         row["id"].Int64(),
		Username:   row["username"].String(),
		Realname:   row["realname"].String(),
		Avatar:     upload.GetFileUrl(ctx, row["avatar"].String()), // 逃生舱手拼 DTO，直接调底层公共助手补域名
		Gender:     row["gender"].Int(),
		GenderText: DictGetText(ctx, userSerializeMaps["gender"], row["gender"].Int()),
		Mobile:     row["mobile"].String(),
		Email:      row["email"].String(),
		Roles:      roleIds,
		DeptId:     row["dept_id"].Int(),
		LevelId:    row["level_id"].Int(),
		PositionId: row["position_id"].Int(),
		City:       city,
		CityInfo:   row["city_info"].String(),
		Address:    row["address"].String(),
		Intro:      row["intro"].String(),
		Status:     row["status"].Int(),
		StatusText: DictGetText(ctx, userSerializeMaps["status"], row["status"].Int()),
		Sort:       row["sort"].Int(),
		Note:       row["note"].String(),
	}
}

// namesById 查询ID到名称的映射
// 传入的模型需已构造好上下文
func (s *sUser) namesById(ctx context.Context, m *gdb.Model, ids []int64) map[int64]string {
	result := make(map[int64]string)
	if len(ids) == 0 {
		return result
	}
	rows, err := m.WhereIn("id", ids).Fields("id, name").All()
	if err != nil {
		return result
	}
	for _, r := range rows {
		result[r["id"].Int64()] = r["name"].String()
	}
	return result
}

// dictName 按ID查询名称
// 适用于部门、职级、岗位等通用维表
func (s *sUser) dictName(ctx context.Context, m *gdb.Model, idv int) string {
	if idv <= 0 {
		return ""
	}
	one, err := m.Where("id", idv).Fields("name").One()
	if err != nil || len(one) == 0 {
		return ""
	}
	return one["name"].String()
}

// cityCodeAt 按下标取城市编码
// 下标越界时返回空串
func cityCodeAt(cities []string, i int) string {
	if i < len(cities) {
		return cities[i]
	}
	return ""
}

// cellAt 按下标取单元格文本
// 下标越界时返回空串
func cellAt(row []string, i int) string {
	if i < len(row) {
		return row[i]
	}
	return ""
}

// cellInt 按下标取单元格整数
// 解析失败时返回默认值
func cellInt(row []string, i, def int) int {
	if v := cellAt(row, i); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return def
}

// idNameMap 加载维表的ID到名称映射
// 供部门、职级、岗位导出人读化使用
func idNameMap(ctx context.Context, m *gdb.Model) (map[int64]string, error) {
	rows, err := m.Fields("id, name").All()
	if err != nil {
		return nil, err
	}
	result := make(map[int64]string, len(rows))
	for _, r := range rows {
		result[r["id"].Int64()] = r["name"].String()
	}
	return result, nil
}

// resolveRefID 解析引用列
// 数字按ID处理，非数字按名称反查，未命中时返回 0
func resolveRefID(ctx context.Context, m *gdb.Model, v string) int64 {
	v = strings.TrimSpace(v)
	if v == "" {
		return 0
	}
	if n, err := strconv.ParseInt(v, 10, 64); err == nil {
		return n
	}
	id, err := m.Where("name", v).Fields("id").Value()
	if err != nil {
		return 0
	}
	return id.Int64()
}

// parseGenderCell 解析性别单元格
// 支持男、女、保密文本或 1、2、3 数字，其余取保密
func parseGenderCell(v string) int {
	switch strings.TrimSpace(v) {
	case "男":
		return 1
	case "女":
		return 2
	case "保密":
		return 3
	}
	if n, err := strconv.Atoi(strings.TrimSpace(v)); err == nil && n >= 1 && n <= 3 {
		return n
	}
	return 3
}

// parseStatusCell 解析状态单元格
// 支持正常、禁用文本或 1、2 数字，其余取正常
func parseStatusCell(v string) int {
	switch strings.TrimSpace(v) {
	case "正常":
		return 1
	case "禁用":
		return 2
	}
	if n, err := strconv.Atoi(strings.TrimSpace(v)); err == nil && n > 0 {
		return n
	}
	return 1
}

// genderText 性别值转文本
func genderText(g int) string {
	switch g {
	case 1:
		return "男"
	case 2:
		return "女"
	default:
		return "保密"
	}
}

// userStatusText 用户状态值转文本
func userStatusText(st int) string {
	if st == 1 {
		return "正常"
	}
	return "禁用"
}
