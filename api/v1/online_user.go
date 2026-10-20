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

// ============================== 在线用户模块 ==============================

// OnlineUserItem 在线用户行数据
type OnlineUserItem struct {
	Id                string `json:"id"`                // 访问令牌标识（jti）
	UserId            int    `json:"userId"`            // 用户ID
	Username          string `json:"username"`          // 用户名
	Realname          string `json:"realname"`          // 姓名
	Avatar            string `json:"avatar"`            // 头像（已补全访问域名）
	Ip                string `json:"ip"`                // 登录IP
	Location          string `json:"location"`          // 登录地点
	Oss               string `json:"oss"`               // 操作系统
	Browser           string `json:"browser"`           // 浏览器
	LoginTime         string `json:"loginTime"`         // 登录时间
	IssuedTime        string `json:"issuedTime"`        // 令牌颁发时间
	ExpireTime        string `json:"expireTime"`        // 令牌过期时间
	RefreshIssuedTime string `json:"refreshIssuedTime"` // 刷新令牌颁发时间（空表示无）
	RefreshExpireTime string `json:"refreshExpireTime"` // 刷新令牌过期时间（空表示无）
	SessionNo         int    `json:"sessionNo"`         // 同用户会话序号
	Sessions          int    `json:"sessions"`          // 同用户会话总数
	TokenType         string `json:"tokenType"`         // 令牌类型（Bearer）
	ExpiresIn         int64  `json:"expiresIn"`         // 访问令牌有效期（秒）
	AccessToken       string `json:"accessToken"`       // 访问令牌原文（仅详情返回）
	RefreshToken      string `json:"refreshToken"`      // 刷新令牌原文（仅详情返回）
}

// OnlineUserPageReq 在线用户分页请求
type OnlineUserPageReq struct {
	gmeta.Meta `method:"get" path:"/api/onlineUser/page" summary:"在线用户分页" permission:"sys:onlineUser:page" log:"在线用户-查询分页记录" logType:"query"`
	PageReq           // 分页基础参数
	Username   string `json:"username" dc:"用户名模糊"` // 用户名模糊匹配
	Ip         string `json:"ip" dc:"登录IP模糊"`      // 登录IP模糊匹配
}

// OnlineUserPageRes 在线用户分页结果
type OnlineUserPageRes = PageData[OnlineUserItem]

// OnlineUserDetailReq 在线会话详情请求
type OnlineUserDetailReq struct {
	gmeta.Meta `method:"get" path:"/api/onlineUser/detail/:id" summary:"在线会话详情" permission:"sys:onlineUser:detail" log:"在线用户-查询会话详情" logType:"query"`
	Id         string `json:"id" v:"required#缺少令牌标识" dc:"令牌标识"`
}

// OnlineUserDetailRes 在线会话详情结果
// 含访问/刷新令牌原文与登录环境信息
type OnlineUserDetailRes struct {
	OnlineUserItem
}

// OnlineUserForceLogoutReq 强制下线请求
type OnlineUserForceLogoutReq struct {
	gmeta.Meta `method:"get" path:"/api/onlineUser/forceLogout/:id" summary:"强制下线" permission:"sys:onlineUser:forceLogout" log:"在线用户-强制下线" logType:"update"`
	Id         string `json:"id" v:"required#缺少令牌标识" dc:"令牌标识"`
}

// OnlineUserForceLogoutRes 强制下线结果
type OnlineUserForceLogoutRes struct{}

// OnlineUserForceLogoutUserReq 踢出该用户全部设备请求
type OnlineUserForceLogoutUserReq struct {
	gmeta.Meta `method:"get" path:"/api/onlineUser/forceLogoutUser/:userId" summary:"踢出用户全部设备" permission:"sys:onlineUser:forceLogout" log:"在线用户-踢出全部设备" logType:"update"`
	UserId     int `json:"userId" v:"required#缺少用户ID" dc:"用户ID"`
}

// OnlineUserForceLogoutUserRes 踢出该用户全部设备结果
type OnlineUserForceLogoutUserRes struct{}
