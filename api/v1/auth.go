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

package v1

import (
	"github.com/gogf/gf/v2/util/gmeta"
)

// ============================== 认证授权模块 ==============================

// LoginReq 登录请求
// 参数支持 query/body 混合提交
type LoginReq struct {
	gmeta.Meta `method:"post" path:"/api/login" summary:"登录" public:"true"`
	Username   string `json:"username" v:"required#请输入用户名" dc:"用户名"` // 用户名
	Password   string `json:"password" v:"required#请输入密码" dc:"密码"`   // 密码
	Code       string `json:"code" dc:"图形验证码"`                       // 图形验证码
	Key        string `json:"key" dc:"验证码标识（captcha 接口返回的 key）"`     // 验证码标识
	TenantCode string `json:"tenantCode" dc:"租户编码，留空按用户名全局定位"`       // 租户编码
}

// LoginRes 登录结果
// 客户端从 data.access_token 中读取并存储令牌
type LoginRes struct {
	AccessToken  string `json:"access_token"`  // 访问令牌
	RefreshToken string `json:"refresh_token"` // 刷新令牌
	ExpiresIn    int64  `json:"expires_in"`    // access 有效期（秒）
	TokenType    string `json:"token_type"`    // 令牌类型，固定为 Bearer
}

// LogoutReq 登出请求
// 拉黑当前 access token，需携带 access token 走 JWT 校验
type LogoutReq struct {
	gmeta.Meta `method:"get" path:"/api/logout" summary:"登出"`
}

// LogoutRes 登出结果
type LogoutRes struct{}

// CaptchaReq 图形验证码请求
type CaptchaReq struct {
	gmeta.Meta `method:"get" path:"/api/captcha" summary:"获取图形验证码" public:"true"`
}

// CaptchaRes 图形验证码结果
type CaptchaRes struct {
	Key     string `json:"key" dc:"验证码标识"`               // 验证码标识
	Captcha string `json:"captcha" dc:"SVG data URL 图片"` // SVG data URL 图片
}

// RefreshReq 刷新令牌请求
// 同时兼容 refresh_token 与 password 两种授权模式
type RefreshReq struct {
	gmeta.Meta   `method:"post" path:"/api/oauth2/token" summary:"刷新令牌/密码模式兑换" public:"true"`
	GrantType    string `json:"grant_type" dc:"refresh_token | password"`              // 授权类型
	RefreshToken string `json:"refresh_token" dc:"刷新令牌（grant_type=refresh_token 时必填）"` // 刷新令牌
	Username     string `json:"username" dc:"用户名（grant_type=password 时必填）"`            // 用户名
	Password     string `json:"password" dc:"密码（grant_type=password 时必填）"`             // 密码
	Code         string `json:"code" dc:"验证码"`                                         // 验证码
	Key          string `json:"key" dc:"验证码标识"`                                        // 验证码标识
}

// RefreshRes 刷新令牌结果
// 结构与登录结果保持一致
type RefreshRes struct {
	AccessToken  string `json:"access_token"`  // 访问令牌
	RefreshToken string `json:"refresh_token"` // 刷新令牌
	ExpiresIn    int64  `json:"expires_in"`    // access 有效期（秒）
	TokenType    string `json:"token_type"`    // 令牌类型
}

// CheckTokenReq 校验令牌请求
type CheckTokenReq struct {
	gmeta.Meta `method:"get" path:"/api/oauth2/check/token" summary:"校验令牌"`
}

// CheckTokenRes 校验令牌结果
type CheckTokenRes struct {
	Valid  bool  `json:"valid"`  // 是否有效
	UserId int64 `json:"userId"` // 用户ID
}

// RemoveTokenReq 指定令牌拉黑请求
type RemoveTokenReq struct {
	gmeta.Meta `method:"delete" path:"/api/oauth2/remove/:token" summary:"撤销指定令牌"`
	Token      string `json:"token" v:"required#缺少token" dc:"待撤销的令牌"` // 待撤销的令牌
}

// RemoveTokenRes 撤销令牌结果
type RemoveTokenRes struct{}
