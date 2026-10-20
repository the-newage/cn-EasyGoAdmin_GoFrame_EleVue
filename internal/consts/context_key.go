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

package consts

import (
	"context"
	"time"

	"github.com/gogf/gf/v2/frame/g"
)

// ============================== 请求上下文键与读写工具 ==============================

// contextKey 上下文键私有类型
// 使用私有类型避免跨包键名冲突
type contextKey string

// 请求上下文键
// 由 JWT、租户中间件写入，供全链路读取
const (
	KeyUserID    contextKey = "userID"      // 当前登录用户ID
	KeyTenantID  contextKey = "tenantID"    // 当前租户ID（0=平台级）
	KeyUsername  contextKey = "username"    // 当前登录用户名
	KeyRealname  contextKey = "realname"    // 当前登录用户姓名
	KeyRoles     contextKey = "roles"       // 角色编码集合 []string
	KeyPerms     contextKey = "perms"       // 权限节点集合 []string，如 sys:user:page
	KeyTokenJTI  contextKey = "tokenJTI"    // 当前令牌 jti
	KeyTokenExp  contextKey = "tokenExpire" // 令牌过期时间 time.Time，用于登出黑名单TTL
	KeyDataScope contextKey = "dataScope"   // 角色数据范围（多角色取最宽松），见 DataScopeXxx 常量
	KeyDeptID    contextKey = "deptID"      // 当前登录用户部门ID（0=未分配）

	// 请求元信息，由 Response 中间件写入，供登录日志、操作日志读取
	KeyClientIP  contextKey = "clientIP"  // 客户端IP
	KeyUserAgent contextKey = "userAgent" // User-Agent
	KeyReqMethod contextKey = "reqMethod" // HTTP方法
	KeyReqURL    contextKey = "reqURL"    // 请求URL
)

// WithUserID 写入当前登录用户ID
func WithUserID(ctx context.Context, id int) context.Context {
	return context.WithValue(ctx, KeyUserID, id)
}

// GetUserID 读取当前登录用户ID
// 未登录时返回 0
func GetUserID(ctx context.Context) int {
	v, _ := ctx.Value(KeyUserID).(int)
	return v
}

// WithTenantID 写入当前租户ID
// 0 表示平台级
func WithTenantID(ctx context.Context, id int) context.Context {
	return context.WithValue(ctx, KeyTenantID, id)
}

// CtxTenantID 读取上下文原始租户ID
//
// 未登录时返回 0（平台级）。仅供租户隔离设施使用（超管判断、查询过滤、中间件校验），
// 业务模块取归属租户请使用 GetTenantID。
func CtxTenantID(ctx context.Context) int {
	v, _ := ctx.Value(KeyTenantID).(int)
	return v
}

// DefaultTenantID 默认归属租户ID
//
// 超管新增租户级数据时使用。读取配置 tenant.defaultId，未配置或非法值（≤0）回退 1。
func DefaultTenantID(ctx context.Context) int {
	v := g.Cfg().MustGet(ctx, "tenant.defaultId").Int()
	if v <= 0 {
		return 1
	}
	return v
}

// GetTenantID 读取当前上下文的归属租户ID
//
// 租户用户返回自身租户ID，超管返回默认租户ID。
// 业务模块新增、导入、唯一校验等取归属租户统一使用本方法。
func GetTenantID(ctx context.Context) int {
	if t := CtxTenantID(ctx); t > 0 {
		return t
	}
	return DefaultTenantID(ctx)
}

// IsPlatform 是否平台级租户
// 上下文租户ID为0时表示平台级（超管）
func IsPlatform(ctx context.Context) bool {
	return CtxTenantID(ctx) == 0
}

// WithUsername 写入当前登录用户名
func WithUsername(ctx context.Context, username string) context.Context {
	return context.WithValue(ctx, KeyUsername, username)
}

// GetUsername 读取当前登录用户名
func GetUsername(ctx context.Context) string {
	v, _ := ctx.Value(KeyUsername).(string)
	return v
}

// WithRealname 写入当前登录用户姓名
func WithRealname(ctx context.Context, realname string) context.Context {
	return context.WithValue(ctx, KeyRealname, realname)
}

// GetRealname 读取当前登录用户姓名
// 审计字段 create_user 优先使用本方法
func GetRealname(ctx context.Context) string {
	v, _ := ctx.Value(KeyRealname).(string)
	return v
}

// WithRoles 写入当前用户角色编码
func WithRoles(ctx context.Context, roles []string) context.Context {
	return context.WithValue(ctx, KeyRoles, roles)
}

// GetRoles 读取当前用户角色编码
func GetRoles(ctx context.Context) []string {
	v, _ := ctx.Value(KeyRoles).([]string)
	return v
}

// WithPerms 写入当前用户权限节点
func WithPerms(ctx context.Context, perms []string) context.Context {
	return context.WithValue(ctx, KeyPerms, perms)
}

// GetPerms 读取当前用户权限节点
func GetPerms(ctx context.Context) []string {
	v, _ := ctx.Value(KeyPerms).([]string)
	return v
}

// HasPerm 校验权限节点
// 权限节点为空视为无需校验；超级管理员拥有全部权限
func HasPerm(ctx context.Context, permission string) bool {
	if permission == "" {
		return true
	}
	for _, r := range GetRoles(ctx) {
		if r == RoleCodeSuperAdmin {
			return true
		}
	}
	for _, p := range GetPerms(ctx) {
		if p == permission {
			return true
		}
	}
	return false
}

// 数据范围档位（角色 data_scope 列值；数值越小可见面越宽）
const (
	DataScopeAll          = 1 // 全部（租户隔离之内）
	DataScopeDeptAndBelow = 2 // 本部门及以下
	DataScopeDept         = 3 // 本部门
	DataScopeSelf         = 4 // 仅本人
)

// WithDataScope 写入当前用户数据范围
func WithDataScope(ctx context.Context, scope int) context.Context {
	return context.WithValue(ctx, KeyDataScope, scope)
}

// GetDataScope 读取当前用户数据范围
// 未注入时返回 DataScopeAll（视为全量可见）
func GetDataScope(ctx context.Context) int {
	v, _ := ctx.Value(KeyDataScope).(int)
	if v <= 0 {
		return DataScopeAll
	}
	return v
}

// WithDeptID 写入当前登录用户部门ID
func WithDeptID(ctx context.Context, id int) context.Context {
	return context.WithValue(ctx, KeyDeptID, id)
}

// GetDeptID 读取当前登录用户部门ID
// 未分配或未登录时返回 0
func GetDeptID(ctx context.Context) int {
	v, _ := ctx.Value(KeyDeptID).(int)
	return v
}

// WithTokenJTI 写入当前令牌 jti
// 登出时依据该值将令牌加入黑名单
func WithTokenJTI(ctx context.Context, jti string) context.Context {
	return context.WithValue(ctx, KeyTokenJTI, jti)
}

// GetTokenJTI 读取当前令牌 jti
func GetTokenJTI(ctx context.Context) string {
	v, _ := ctx.Value(KeyTokenJTI).(string)
	return v
}

// WithTokenExpire 写入当前令牌过期时间
// 黑名单 TTL 取令牌剩余有效期
func WithTokenExpire(ctx context.Context, t time.Time) context.Context {
	return context.WithValue(ctx, KeyTokenExp, t)
}

// GetTokenExpire 读取当前令牌过期时间
func GetTokenExpire(ctx context.Context) time.Time {
	v, _ := ctx.Value(KeyTokenExp).(time.Time)
	return v
}

// WithRequestMeta 写入请求元信息
// 由 Response 中间件调用
func WithRequestMeta(ctx context.Context, ip, userAgent, method, url string) context.Context {
	ctx = context.WithValue(ctx, KeyClientIP, ip)
	ctx = context.WithValue(ctx, KeyUserAgent, userAgent)
	ctx = context.WithValue(ctx, KeyReqMethod, method)
	ctx = context.WithValue(ctx, KeyReqURL, url)
	return ctx
}

// GetClientIP 读取客户端IP
func GetClientIP(ctx context.Context) string {
	v, _ := ctx.Value(KeyClientIP).(string)
	return v
}

// GetUserAgent 读取 User-Agent
func GetUserAgent(ctx context.Context) string {
	v, _ := ctx.Value(KeyUserAgent).(string)
	return v
}

// GetReqMethod 读取HTTP方法
func GetReqMethod(ctx context.Context) string {
	v, _ := ctx.Value(KeyReqMethod).(string)
	return v
}

// GetReqURL 读取请求URL
func GetReqURL(ctx context.Context) string {
	v, _ := ctx.Value(KeyReqURL).(string)
	return v
}
