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
	"encoding/json"
	"strings"
	"time"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"

	"easygoadmin/api/v1"
	"easygoadmin/internal/consts"
	"easygoadmin/internal/dao"
	"easygoadmin/internal/middleware"
	"easygoadmin/internal/service"
	"easygoadmin/utility/cache"
	"easygoadmin/utility/captcha"
	"easygoadmin/utility/crypto"
	"easygoadmin/utility/ip2region"
	jwtutil "easygoadmin/utility/jwt"
	"easygoadmin/utility/ua"
)

// ============================== 认证业务实现 ==============================

// sAuth 认证业务实现
// 实现 service.IAuth 接口，承载登录、登出、验证码、令牌等业务逻辑
type sAuth struct{}

func init() {
	service.RegisterAuth(&sAuth{})
	// 中间件依赖接线：身份扩展加载 / 租户状态校验
	middleware.SetIdentityLoader(LoadIdentity)
	middleware.SetTenantLoader(LoadTenant)
}

// Login 登录
// 流程：验证码校验 → 定位用户 → 账号状态与租户校验 → 密码校验 → 签发双令牌 → 写入登录日志
func (s *sAuth) Login(ctx context.Context, req *v1.LoginReq) (res *v1.LoginRes, err error) {
	start := time.Now()
	username := strings.TrimSpace(req.Username)

	// 1. 图形验证码（提供 key 时强制校验）
	if !captcha.Verify(ctx, req.Key, req.Code) {
		s.writeLoginLog(ctx, username, consts.LoginLogTypeLogin, 1, "验证码错误或已过期", start, req)
		return nil, gerror.NewCode(consts.CodeCaptchaError)
	}

	// 2. 定位用户（用户名在租户内唯一；未指定租户时全局匹配，优先平台租户）
	user, err := s.findUser(ctx, username, req.TenantCode)
	if err != nil {
		s.writeLoginLog(ctx, username, consts.LoginLogTypeLogin, 1, err.Error(), start, req)
		return nil, err
	}

	// 3. 账号状态
	if user.Status != 1 {
		s.writeLoginLog(ctx, username, consts.LoginLogTypeLogin, 1, "账号已被禁用", start, req)
		return nil, gerror.NewCode(consts.CodeUserDisabled)
	}

	// 4. 租户状态与过期校验
	if info := mustLoadTenant(ctx, user.TenantId); info != nil {
		if info.Status != 1 {
			s.writeLoginLog(ctx, username, consts.LoginLogTypeLogin, 1, "租户已禁用", start, req)
			return nil, gerror.NewCode(consts.CodeTenantDisabled)
		}
		if info.ExpireTime != nil && info.ExpireTime.Time.Before(time.Now()) {
			s.writeLoginLog(ctx, username, consts.LoginLogTypeLogin, 1, "租户已过期", start, req)
			return nil, gerror.NewCode(consts.CodeTenantExpired)
		}
	}

	// 5. 密码校验（外部 salt + bcrypt(password+salt)；salt 为空兼容旧数据）
	if !crypto.IsBcrypt(user.Password) || !crypto.VerifyPassword(user.Password, req.Password, user.Salt) {
		s.writeLoginLog(ctx, username, consts.LoginLogTypeLogin, 1, "密码错误", start, req)
		return nil, gerror.NewCode(consts.CodePasswordWrong)
	}

	// 6. 签发双令牌
	loginCtx := consts.WithUserID(consts.WithUsername(ctx, username), int(user.Id))
	loginCtx = consts.WithTenantID(loginCtx, user.TenantId)
	access, accessClaims, err := jwtutil.Sign(loginCtx, jwtutil.TokenTypeAccess,
		int(user.Id), user.TenantId, user.Username, user.Realname)
	if err != nil {
		return nil, err
	}
	refresh, refreshClaims, err := jwtutil.Sign(loginCtx, jwtutil.TokenTypeRefresh,
		int(user.Id), user.TenantId, user.Username, user.Realname)
	if err != nil {
		return nil, err
	}

	// 7. 登录日志（成功）
	s.writeLoginLog(ctx, username, consts.LoginLogTypeLogin, 0, "", start, req)

	// 在线登记（访问令牌 + 刷新令牌信息）
	MarkOnline(ctx, accessClaims, refreshClaims, access, refresh, int(user.Id), user.Username, user.Realname,
		consts.GetClientIP(ctx), consts.GetUserAgent(ctx))

	expireSec := int64(jwtutil.ExpireDuration(ctx, jwtutil.TokenTypeAccess).Seconds())
	return &v1.LoginRes{
		AccessToken:  access,
		RefreshToken: refresh,
		ExpiresIn:    expireSec,
		TokenType:    "Bearer",
	}, nil
}

// Logout 登出
// 将当前 access token 的 jti 加入黑名单，并写入登录日志
func (s *sAuth) Logout(ctx context.Context, req *v1.LogoutReq) (res *v1.LogoutRes, err error) {
	start := time.Now()
	jti := consts.GetTokenJTI(ctx)
	if jti != "" {
		ttl := int(time.Until(consts.GetTokenExpire(ctx)).Seconds())
		if ttl > 0 {
			_ = cache.Set(ctx, consts.CacheKeyTokenBlacklist(jti), "1", ttl)
		}
		MarkOffline(ctx, jti)
	}
	s.writeLoginLog(ctx, consts.GetUsername(ctx), consts.LoginLogTypeLogout, 0, "", start, nil)
	return &v1.LogoutRes{}, nil
}

// Captcha 生成图形验证码
// 返回验证码标识与图片内容
func (s *sAuth) Captcha(ctx context.Context, req *v1.CaptchaReq) (res *v1.CaptchaRes, err error) {
	key, dataURL, err := captcha.Generate(ctx)
	if err != nil {
		return nil, err
	}
	return &v1.CaptchaRes{Key: key, Captcha: dataURL}, nil
}

// Refresh 刷新令牌
// 支持 refresh_token 轮换，并兼容密码模式兑换
func (s *sAuth) Refresh(ctx context.Context, req *v1.RefreshReq) (res *v1.RefreshRes, err error) {
	// 密码模式（前端 login2 变体）：复用 Login
	if req.GrantType == "password" || req.RefreshToken == "" {
		loginRes, lerr := s.Login(ctx, &v1.LoginReq{
			Username: req.Username,
			Password: req.Password,
			Code:     req.Code,
			Key:      req.Key,
		})
		if lerr != nil {
			return nil, lerr
		}
		return &v1.RefreshRes{
			AccessToken:  loginRes.AccessToken,
			RefreshToken: loginRes.RefreshToken,
			ExpiresIn:    loginRes.ExpiresIn,
			TokenType:    loginRes.TokenType,
		}, nil
	}

	// refresh_token 模式
	claims, err := jwtutil.Parse(ctx, req.RefreshToken)
	if err != nil {
		return nil, err
	}
	if claims.Type != jwtutil.TokenTypeRefresh {
		return nil, gerror.NewCode(consts.CodeTokenInvalid)
	}
	// 旧 refresh 已拉黑则失效
	if cache.Exists(ctx, consts.CacheKeyTokenBlacklist(claims.ID)) {
		return nil, gerror.NewCode(consts.CodeTokenInvalid)
	}
	// 轮换：旧 refresh 入黑名单
	if claims.ExpiresAt != nil {
		ttl := int(time.Until(claims.ExpiresAt.Time).Seconds())
		if ttl > 0 {
			_ = cache.Set(ctx, consts.CacheKeyTokenBlacklist(claims.ID), "1", ttl)
		}
	}
	access, accessClaims, err := jwtutil.Sign(ctx, jwtutil.TokenTypeAccess,
		claims.UserID, claims.TenantID, claims.Username, claims.Realname)
	if err != nil {
		return nil, err
	}
	// 令牌轮换后刷新在线登记
	refresh, refreshClaims, err := jwtutil.Sign(ctx, jwtutil.TokenTypeRefresh,
		claims.UserID, claims.TenantID, claims.Username, claims.Realname)
	if err != nil {
		return nil, err
	}
	// 令牌轮换后刷新在线登记
	MarkOnline(ctx, accessClaims, refreshClaims, access, refresh, claims.UserID, claims.Username, claims.Realname,
		consts.GetClientIP(ctx), consts.GetUserAgent(ctx))
	return &v1.RefreshRes{
		AccessToken:  access,
		RefreshToken: refresh,
		ExpiresIn:    int64(jwtutil.ExpireDuration(ctx, jwtutil.TokenTypeAccess).Seconds()),
		TokenType:    "Bearer",
	}, nil
}

// CheckToken 校验令牌有效性
// JWT 中间件已通过即视为有效，返回当前登录用户ID
func (s *sAuth) CheckToken(ctx context.Context, req *v1.CheckTokenReq) (res *v1.CheckTokenRes, err error) {
	return &v1.CheckTokenRes{
		Valid:  true,
		UserId: int64(consts.GetUserID(ctx)),
	}, nil
}

// RemoveToken 指定令牌撤销（踢下线）
// 将指定令牌加入黑名单
func (s *sAuth) RemoveToken(ctx context.Context, req *v1.RemoveTokenReq) (res *v1.RemoveTokenRes, err error) {
	claims, err := jwtutil.Parse(ctx, req.Token)
	if err != nil {
		return nil, err
	}
	if claims.ExpiresAt != nil {
		ttl := int(time.Until(claims.ExpiresAt.Time).Seconds())
		if ttl > 0 {
			_ = cache.Set(ctx, consts.CacheKeyTokenBlacklist(claims.ID), "1", ttl)
		}
	}
	return &v1.RemoveTokenRes{}, nil
}

// ---------- 定位与校验 ----------

// findUser 按用户名与可选租户编码定位用户
// 未指定租户时全局匹配，多租户同名优先平台租户，仍歧义则提示携带租户编码
func (s *sAuth) findUser(ctx context.Context, username, tenantCode string) (*entityUser, error) {
	if username == "" {
		return nil, gerror.NewCode(consts.CodeUserNotFound)
	}
	query := dao.User.Ctx(ctx).Where("username", username)
	tenantID := 0
	if tenantCode != "" {
		tOne, terr := dao.Tenant.Ctx(ctx).Where("code", tenantCode).One()
		if terr != nil || len(tOne) == 0 {
			return nil, gerror.NewCode(consts.CodeTenantNotFound)
		}
		tenantID = tOne["id"].Int()
		query = query.Where("tenant_id", tenantID)
	}
	list, err := query.All()
	if err != nil {
		return nil, err
	}
	if len(list) == 0 {
		return nil, gerror.NewCode(consts.CodeUserNotFound)
	}
	// 多租户同名：优先平台租户(tenant_id=0)，仍歧义则要求携带租户编码
	if len(list) == 1 {
		return rowToUser(list[0]), nil
	}
	for _, row := range list {
		if row["tenant_id"].Int() == 0 {
			return rowToUser(row), nil
		}
	}
	return nil, gerror.NewCodef(consts.CodeUserNotFound, "该用户名存在于多个租户，请携带租户编码登录")
}

// entityUser 用户摘要
// 登录所需的最小字段集合
type entityUser struct {
	Id       int64  // 用户ID
	TenantId int    // 租户ID
	Username string // 用户名
	Password string // 密码
	Salt     string // 密码盐
	Realname string // 真实姓名
	Status   int    // 状态：1-正常 2-禁用
}

// rowToUser 数据库记录转换为用户摘要
func rowToUser(row gdb.Record) *entityUser {
	return &entityUser{
		Id:       row["id"].Int64(),
		TenantId: row["tenant_id"].Int(),
		Username: row["username"].String(),
		Password: row["password"].String(),
		Salt:     row["salt"].String(),
		Realname: row["realname"].String(),
		Status:   row["status"].Int(),
	}
}

// ---------- 中间件加载器 ----------

// LoadIdentity 加载用户角色与权限
// 由 JWT 中间件调用，结果由中间件缓存 300 秒
func LoadIdentity(ctx context.Context, userID int) (*middleware.Identity, error) {
	// DataScope 以 0 为"未取到"哨兵，角色循环内取最小（最宽松），结束时回退"全部"
	idn := &middleware.Identity{Roles: []string{}, Perms: []string{}}
	if userID <= 0 {
		return idn, nil
	}
	// 所属部门（数据权限过滤用）
	if one, e := dao.User.Ctx(ctx).Where("id", userID).Fields("dept_id").One(); e == nil && len(one) > 0 {
		idn.DeptId = one["dept_id"].Int()
	}
	// 角色编码
	roleRows, err := dao.UserRole.Ctx(ctx).Where("user_id", userID).All()
	if err != nil || len(roleRows) == 0 {
		return idn, err
	}
	roleIds := make([]int64, 0, len(roleRows))
	for _, r := range roleRows {
		roleIds = append(roleIds, r["role_id"].Int64())
	}
	roleRows2, err := dao.Role.Ctx(ctx).WhereIn("id", roleIds).All()
	if err != nil {
		return idn, err
	}
	for _, r := range roleRows2 {
		idn.Roles = append(idn.Roles, r["code"].String())
		// 数据范围多角色取最宽松（数值越小越宽；缺列值按"全部"参与比较）
		sc := r["data_scope"].Int()
		if sc <= 0 {
			sc = consts.DataScopeAll
		}
		if idn.DataScope == 0 || sc < idn.DataScope {
			idn.DataScope = sc
		}
	}
	if idn.DataScope == 0 {
		idn.DataScope = consts.DataScopeAll // 无角色或角色均无该列时回退
	}
	// 权限节点（role → menu.permission）
	menuRows, err := dao.RoleMenu.Ctx(ctx).WhereIn("role_id", roleIds).All()
	if err != nil || len(menuRows) == 0 {
		return idn, err
	}
	menuIds := make([]int64, 0, len(menuRows))
	for _, r := range menuRows {
		menuIds = append(menuIds, r["menu_id"].Int64())
	}
	permRows, err := dao.Menu.Ctx(ctx).WhereIn("id", menuIds).
		Where("type", 1).Where("status", 0).Where("hide", 0).All() // 仅取权限节点（type=1），启用且非隐藏
	if err != nil {
		return idn, err
	}
	seen := make(map[string]bool)
	for _, r := range permRows {
		p := r["permission"].String()
		if p != "" && !seen[p] {
			seen[p] = true
			idn.Perms = append(idn.Perms, p)
		}
	}
	return idn, nil
}

// LoadTenant 加载租户状态摘要
// 由 Tenant 中间件调用，缓存 60 秒；租户不存在按禁用处理，实现失效关闭
func LoadTenant(ctx context.Context, tenantID int) (*middleware.TenantInfo, error) {
	if tenantID == 0 {
		return &middleware.TenantInfo{ID: 0, Status: 1}, nil
	}
	one, err := dao.Tenant.Ctx(ctx).Where("id", tenantID).One()
	if err != nil {
		return nil, err // 基础设施错误：中间件放行
	}
	if len(one) == 0 {
		return &middleware.TenantInfo{ID: tenantID, Status: 2}, nil // 租户不存在 → 禁用
	}
	info := &middleware.TenantInfo{
		ID:     tenantID,
		Status: one["status"].Int(),
	}
	if v := one["expire_time"]; !v.IsNil() {
		t := v.GTime()
		if t != nil {
			info.ExpireTime = t
		}
	}
	return info, nil
}

// mustLoadTenant 登录路径上的租户校验
// 带 60 秒缓存，避免登录过程中重复查询租户信息
func mustLoadTenant(ctx context.Context, tenantID int) *middleware.TenantInfo {
	ck := consts.CacheKeyTenantInfo(tenantID)
	if v, ok := cache.Get(ctx, ck); ok {
		info := &middleware.TenantInfo{}
		if json.Unmarshal([]byte(v), info) == nil {
			return info
		}
	}
	info, err := LoadTenant(ctx, tenantID)
	if err != nil || info == nil {
		return nil
	}
	if b, berr := json.Marshal(info); berr == nil {
		_ = cache.Set(ctx, ck, string(b), 60)
	}
	return info
}

// writeLoginLog 写入登录日志
// 覆盖登录日志表所有 NOT NULL 字段，写入失败仅记录警告不影响主流程
func (s *sAuth) writeLoginLog(ctx context.Context, username string, logType, status int, errMsg string, start time.Time, req *v1.LoginReq) {
	param := ""
	if req != nil {
		if b, err := json.Marshal(map[string]any{
			"username": req.Username,
			"code":     maskValue(req.Code),
			"key":      req.Key,
		}); err == nil {
			param = string(b)
		}
	}
	result := "成功"
	if status != 0 {
		result = "失败"
	}
	// 登录/退出日志的 create_user 显式取本次登录账号：
	// 登录走免登录公开路由，ctx 中无登录身份，审计钩子取不到操作人
	createUser := username
	if createUser == "" {
		createUser = consts.GetUsername(ctx)
	}
	_, err := dao.LoginLog.Ctx(ctx).Data(map[string]any{
		"username":       username,
		"log_type":       logType,
		"source":         consts.LoginSourceAdmin,
		"method":         consts.GetReqMethod(ctx),
		"request_method": consts.GetReqMethod(ctx),
		"url":            consts.GetReqURL(ctx),
		"param":          param,
		"result":         result,
		"ip":             consts.GetClientIP(ctx),
		"location":       ip2region.Location(consts.GetClientIP(ctx)),
		"os":             ua.ParseOS(consts.GetUserAgent(ctx)),
		"browser":        ua.ParseBrowser(consts.GetUserAgent(ctx)),
		"consume_time":   int(time.Since(start).Milliseconds()),
		"status":         status,
		"error":          errMsg,
		"create_user":    createUser,
	}).Insert()
	if err != nil {
		g.Log().Warningf(ctx, "登录日志写入失败: %v", err)
	}
}

// maskValue 敏感值脱敏
// 空值返回空串，非空统一返回掩码
func maskValue(s string) string {
	if s == "" {
		return ""
	}
	return "***"
}
