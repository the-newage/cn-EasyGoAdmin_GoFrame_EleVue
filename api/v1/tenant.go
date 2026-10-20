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
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gmeta"
)

// ============================== 租户模块 ==============================

// TenantItem 租户行数据
// 用于租户分页列表等场景的租户基础信息
type TenantItem struct {
	Id         int64       `json:"id"`         // 主键ID
	Name       string      `json:"name"`       // 租户名称
	Code       string      `json:"code"`       // 租户编码
	Contact    string      `json:"contact"`    // 联系人
	Mobile     string      `json:"mobile"`     // 联系电话
	Email      string      `json:"email"`      // 邮箱
	Address    string      `json:"address"`    // 地址
	Logo       string      `json:"logo"`       // Logo
	Status     int         `json:"status"`     // 状态：1-正常 2-禁用
	StatusText string      `json:"statusText"` // 状态文本
	ExpireTime *gtime.Time `json:"expireTime"` // 过期时间
	MaxUserNum int         `json:"maxUserNum"` // 用户数上限（0=不限）
	UserCount  int         `json:"userCount"`  // 已用用户数
	Note       string      `json:"note"`       // 备注
	CreateUser string      `json:"createUser"` // 创建人
	CreateTime *gtime.Time `json:"createTime"` // 创建时间
}

// TenantPageReq 租户分页请求
// 支持关键词、租户名称、租户编码、状态等条件组合查询
type TenantPageReq struct {
	gmeta.Meta `method:"get" path:"/api/tenant/page" summary:"租户分页" permission:"sys:tenant:page" log:"租户管理-查询分页记录" logType:"query"`
	PageReq           // 分页基础参数
	Keyword    string `json:"keyword" dc:"关键词（名称/编码/联系人模糊）"` // 关键词，名称、编码或联系人模糊匹配
	Name       string `json:"name" dc:"租户名称模糊"`              // 租户名称模糊匹配
	Code       string `json:"code" dc:"租户编码模糊"`              // 租户编码模糊匹配
	Status     int    `json:"status" dc:"状态：1-正常 2-禁用"`      // 状态
}

// TenantPageRes 租户分页结果
type TenantPageRes = PageData[TenantItem]

// TenantListReq 租户全量列表请求
type TenantListReq struct {
	gmeta.Meta `method:"get" path:"/api/tenant/list" summary:"租户列表" permission:"sys:tenant:list" log:"租户管理-查询列表" logType:"query"`
	Keyword    string `json:"keyword" dc:"关键词"` // 关键词
	Status     int    `json:"status" dc:"状态"`   // 状态
}

// TenantListItem 租户列表行数据
// 精简结构，仅包含列表展示所需字段
type TenantListItem struct {
	Id     int64  `json:"id"`     // 主键ID
	Name   string `json:"name"`   // 租户名称
	Code   string `json:"code"`   // 租户编码
	Status int    `json:"status"` // 状态：1-正常 2-禁用
}

// TenantListRes 租户列表结果
type TenantListRes = []TenantListItem

// TenantDetailReq 租户详情请求
type TenantDetailReq struct {
	gmeta.Meta `method:"get" path:"/api/tenant/detail/:id" summary:"租户详情" permission:"sys:tenant:detail" log:"租户管理-查询详情" logType:"query"`
	Id         int64 `json:"id" v:"required#缺少ID" dc:"租户ID"` // 租户ID
}

// TenantDetailRes 租户详情结果
type TenantDetailRes struct {
	Id         int64       `json:"id"`         // 主键ID
	Name       string      `json:"name"`       // 租户名称
	Code       string      `json:"code"`       // 租户编码
	Contact    string      `json:"contact"`    // 联系人
	Mobile     string      `json:"mobile"`     // 联系电话
	Email      string      `json:"email"`      // 邮箱
	Address    string      `json:"address"`    // 地址
	Logo       string      `json:"logo"`       // Logo
	Status     int         `json:"status"`     // 状态：1-正常 2-禁用
	StatusText string      `json:"statusText"` // 状态文本
	ExpireTime *gtime.Time `json:"expireTime"` // 过期时间
	MaxUserNum int         `json:"maxUserNum"` // 用户数上限
	Note       string      `json:"note"`       // 备注
}

// TenantAddReq 新增租户请求
type TenantAddReq struct {
	gmeta.Meta `method:"post" path:"/api/tenant/add" summary:"新增租户" permission:"sys:tenant:add" log:"租户管理-新增记录" logType:"add"`
	Name       string `json:"name" v:"required#请输入租户名称" dc:"租户名称"`    // 租户名称
	Code       string `json:"code" dc:"租户编码（留空自动生成）"`                 // 租户编码
	Contact    string `json:"contact" dc:"联系人"`                       // 联系人
	Mobile     string `json:"mobile" v:"phone#手机号格式不正确" dc:"联系电话"`    // 联系电话
	Email      string `json:"email" v:"email#邮箱格式不正确" dc:"邮箱"`        // 邮箱
	Address    string `json:"address" dc:"地址"`                        // 地址
	Logo       string `json:"logo" dc:"Logo"`                         // Logo
	Status     int    `json:"status" dc:"状态：1-正常 2-禁用"`               // 状态
	ExpireTime string `json:"expireTime" dc:"过期时间（Y-m-d H:i:s，留空永久）"` // 过期时间
	MaxUserNum int    `json:"maxUserNum" dc:"用户数上限（0=不限）"`            // 用户数上限
	Note       string `json:"note" dc:"备注"`                           // 备注
}

// TenantUpdateReq 修改租户请求
type TenantUpdateReq struct {
	gmeta.Meta `method:"put" path:"/api/tenant/update" summary:"修改租户" permission:"sys:tenant:update" log:"租户管理-修改记录" logType:"update"`
	Id         int64  `json:"id" v:"required#缺少ID" dc:"租户ID"`      // 租户ID
	Name       string `json:"name" v:"required#请输入租户名称" dc:"租户名称"` // 租户名称
	Code       string `json:"code" dc:"租户编码（不允许修改）"`               // 租户编码，不允许修改
	Contact    string `json:"contact" dc:"联系人"`                    // 联系人
	Mobile     string `json:"mobile" v:"phone#手机号格式不正确" dc:"联系电话"` // 联系电话
	Email      string `json:"email" v:"email#邮箱格式不正确" dc:"邮箱"`     // 邮箱
	Address    string `json:"address" dc:"地址"`                     // 地址
	Logo       string `json:"logo" dc:"Logo"`                      // Logo
	Status     int    `json:"status" dc:"状态"`                      // 状态
	ExpireTime string `json:"expireTime" dc:"过期时间"`                // 过期时间
	MaxUserNum int    `json:"maxUserNum" dc:"用户数上限"`               // 用户数上限
	Note       string `json:"note" dc:"备注"`                        // 备注
}

// TenantUpdateRes 修改租户结果
type TenantUpdateRes struct{}

// TenantStatusReq 租户状态请求
type TenantStatusReq struct {
	gmeta.Meta `method:"put" path:"/api/tenant/status" summary:"租户状态" permission:"sys:tenant:status" log:"租户管理-设置状态" logType:"update"`
	Id         int64 `json:"id" v:"required#缺少ID" dc:"租户ID"`             // 租户ID
	Status     int   `json:"status" v:"required#缺少状态" dc:"状态：1-正常 2-禁用"` // 状态
}

// TenantStatusRes 租户状态结果
type TenantStatusRes struct{}

// TenantDeleteReq 删除租户请求
type TenantDeleteReq struct {
	gmeta.Meta `method:"delete" path:"/api/tenant/delete/:id" summary:"删除租户" permission:"sys:tenant:delete" log:"租户管理-删除记录" logType:"delete"`
	Id         int64 `json:"id" v:"required#缺少ID" dc:"租户ID"` // 租户ID
}

// TenantDeleteRes 删除租户结果
type TenantDeleteRes struct{}

// TenantBatchDeleteReq 批量删除租户请求
type TenantBatchDeleteReq struct {
	gmeta.Meta `method:"delete" path:"/api/tenant/batchDelete" summary:"批量删除租户" permission:"sys:tenant:batchDelete" log:"租户管理-批量删除记录" logType:"delete"`
	Ids        []int64 `json:"ids" v:"required#请选择要删除的记录" dc:"租户ID列表"` // 租户ID列表
}

// TenantBatchDeleteRes 批量删除结果
type TenantBatchDeleteRes struct{}

// TenantAccountReq 创建租户账号请求
// 同时创建租户及其初始管理员账号
type TenantAccountReq struct {
	gmeta.Meta `method:"post" path:"/api/tenant/account" summary:"创建租户账号" permission:"sys:tenant:add" log:"租户管理-创建租户账号" logType:"add"`
	Name       string `json:"name" v:"required#请输入租户名称" dc:"租户名称"`                                 // 租户名称
	Contact    string `json:"contact" dc:"联系人"`                                                    // 联系人
	Mobile     string `json:"mobile" v:"phone#手机号格式不正确" dc:"联系电话"`                                 // 联系电话
	Email      string `json:"email" v:"email#邮箱格式不正确" dc:"邮箱"`                                     // 邮箱
	Address    string `json:"address" dc:"地址"`                                                     // 地址
	ExpireTime string `json:"expireTime" dc:"过期时间"`                                                // 过期时间
	MaxUserNum int    `json:"maxUserNum" dc:"用户数上限"`                                               // 用户数上限
	Note       string `json:"note" dc:"备注"`                                                        // 备注
	Username   string `json:"username" v:"required#请输入管理员账号|length:2,30#账号长度需在2~30位" dc:"初始管理员账号"` // 初始管理员账号
	Password   string `json:"password" v:"required#请输入密码|length:6,32#密码长度需在6~32位" dc:"初始管理员密码"`    // 初始管理员密码
	Realname   string `json:"realname" dc:"管理员姓名"`                                                 // 管理员姓名
}

// TenantAccountRes 创建租户账号结果
type TenantAccountRes struct {
	TenantId int64  `json:"tenantId"` // 租户ID
	Code     string `json:"code"`     // 租户编码
	UserId   int64  `json:"userId"`   // 管理员用户ID
}

// TenantUsersReq 查询租户用户请求
type TenantUsersReq struct {
	gmeta.Meta `method:"get" path:"/api/tenant/users/:tenantId" summary:"查询租户用户" permission:"sys:tenant:users" log:"租户管理-查询租户用户" logType:"query"`
	TenantId   int64 `json:"tenantId" v:"required#缺少租户ID" dc:"租户ID"` // 租户ID
}

// TenantUserItem 租户用户行数据
// 精简结构，仅包含租户用户列表展示所需字段
type TenantUserItem struct {
	Id       int64  `json:"id"`       // 用户ID
	Username string `json:"username"` // 用户名
	Realname string `json:"realname"` // 真实姓名
	Mobile   string `json:"mobile"`   // 手机号
	Email    string `json:"email"`    // 邮箱
	Status   int    `json:"status"`   // 状态
}

// TenantUsersRes 查询租户用户结果
type TenantUsersRes struct {
	List []TenantUserItem `json:"list"` // 租户用户列表
}
