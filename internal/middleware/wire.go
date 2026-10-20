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

package middleware

import (
	"context"
	"time"

	"github.com/gogf/gf/v2/os/gtime"
)

// ============================================================================
// 依赖接线：中间件不直接依赖业务层，由 cmd 在启动时注册加载器与落地器。
// ============================================================================

// ============================== 依赖类型 ==============================

// Identity 用户身份扩展信息
// 包含角色编码与权限节点
type Identity struct {
	Roles     []string `json:"roles"`     // 角色编码集合
	Perms     []string `json:"perms"`     // 权限节点集合
	DataScope int      `json:"dataScope"` // 数据范围（多角色取最宽松；0=旧缓存缺字段，按全部处理）
	DeptId    int      `json:"deptId"`    // 所属部门ID（0=未分配）
}

// TenantInfo 租户摘要
// 仅保留状态与过期校验热路径所需的最小字段
type TenantInfo struct {
	ID         int         `json:"id"`         // 租户ID
	Status     int         `json:"status"`     // 状态：1-正常 2-禁用
	ExpireTime *gtime.Time `json:"expireTime"` // 过期时间，nil 表示永不过期
}

// OperationLogEntry 操作日志记录
// 字段覆盖操作日志表全部列
type OperationLogEntry struct {
	Title         string `json:"title"`         // 日志标题
	Type          int    `json:"type"`          // 操作类型
	Source        int    `json:"source"`        // 来源
	Method        string `json:"method"`        // 路由模板
	RequestMethod string `json:"requestMethod"` // HTTP 方法
	URL           string `json:"url"`           // 请求地址
	Param         string `json:"param"`         // 脱敏后的请求参数
	Result        string `json:"result"`        // 响应报文（截断）
	IP            string `json:"ip"`            // 客户端IP
	Location      string `json:"location"`      // 归属地
	OS            string `json:"os"`            // 操作系统
	Browser       string `json:"browser"`       // 浏览器
	ConsumeTime   int64  `json:"consumeTime"`   // 耗时（毫秒）
	Status        int    `json:"status"`        // 状态：0-正常 1-异常
	Error         string `json:"error"`         // 错误信息
	CreateUser    string `json:"createUser"`    // 操作人
}

// ============================== 依赖注册 ==============================

var (
	// IdentityLoader 加载用户角色与权限
	// 为 nil 时跳过加载，仅认证身份
	IdentityLoader func(ctx context.Context, userID int) (*Identity, error)

	// TenantLoader 加载租户状态摘要
	// 为 nil 时跳过租户校验
	TenantLoader func(ctx context.Context, tenantID int) (*TenantInfo, error)

	// LogSink 操作日志异步落库
	// 为 nil 时降级为调试日志
	LogSink func(ctx context.Context, entry *OperationLogEntry) error
)

// SetIdentityLoader 注册身份加载器
func SetIdentityLoader(f func(ctx context.Context, userID int) (*Identity, error)) {
	IdentityLoader = f
}

// SetTenantLoader 注册租户加载器
func SetTenantLoader(f func(ctx context.Context, tenantID int) (*TenantInfo, error)) {
	TenantLoader = f
}

// SetLogSink 注册操作日志落地器
func SetLogSink(f func(ctx context.Context, entry *OperationLogEntry) error) {
	LogSink = f
}

// Now 返回当前时间
// 供各中间件统一使用
func Now() time.Time { return time.Now() }
