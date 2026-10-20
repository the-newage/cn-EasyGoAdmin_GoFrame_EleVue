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
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gmeta"
)

// ============================== 用户模块 ==============================

// UserPageReq 用户分页请求
// 支持关键词、姓名、状态、部门、角色、职级、岗位、租户等条件组合查询
type UserPageReq struct {
	gmeta.Meta `method:"get" path:"/api/user/page" summary:"用户分页" permission:"sys:user:page" log:"用户管理-查询分页记录" logType:"query"`
	PageReq           // 分页基础参数
	Keyword    string `json:"keyword" dc:"关键词（用户名/姓名/手机号模糊）"` // 关键词，用户名、姓名或手机号模糊匹配
	Realname   string `json:"realname" dc:"姓名模糊"`             // 姓名模糊匹配
	Status     int    `json:"status" dc:"状态：1-正常 2-禁用"`       // 状态
	DeptId     int    `json:"deptId" dc:"部门ID"`               // 部门ID
	RoleId     int    `json:"roleId" dc:"角色ID"`               // 角色ID
	Role       int    `json:"role" dc:"角色ID（前端参数名兼容）"`        // 角色ID，前端参数名兼容
	LevelId    int    `json:"levelId" dc:"职级ID"`              // 职级ID
	PositionId int    `json:"positionId" dc:"岗位ID"`           // 岗位ID
	TenantId   int    `json:"tenantId" dc:"租户ID"`             // 租户ID
}

// UserRoleRef 用户角色引用
// 仅包含角色ID与名称，用于用户详情中的角色输出
type UserRoleRef struct {
	Id   int64  `json:"id"`   // 角色ID
	Name string `json:"name"` // 角色名称
}

// UserItem 用户行数据
// 用于用户分页列表等场景的用户基础信息
type UserItem struct {
	Id           int64         `json:"id"`           // 主键ID
	Username     string        `json:"username"`     // 用户名
	Code         string        `json:"code"`         // 用户编码
	Realname     string        `json:"realname"`     // 真实姓名
	Avatar       string        `json:"avatar"`       // 头像，已补全访问域名
	Gender       int           `json:"gender"`       // 性别：1-男 2-女 3-保密
	GenderText   string        `json:"genderText"`   // 性别文本
	Mobile       string        `json:"mobile"`       // 手机号
	Email        string        `json:"email"`        // 邮箱
	Roles        []UserRoleRef `json:"roles"`        // 角色引用数组
	City         []string      `json:"city"`         // 省市区街道编码数组
	LevelName    string        `json:"levelName"`    // 职级名称
	PositionName string        `json:"positionName"` // 岗位名称
	DeptName     string        `json:"deptName"`     // 部门名称
	Status       int           `json:"status"`       // 状态：1-正常 2-禁用
	StatusText   string        `json:"statusText"`   // 状态文本
	Sort         int           `json:"sort"`         // 排序
	Note         string        `json:"note"`         // 备注
	CreateUser   string        `json:"createUser"`   // 创建人
	CreateTime   *gtime.Time   `json:"createTime"`   // 创建时间
}

// UserPageRes 用户分页结果
type UserPageRes = PageData[UserItem]

// UserListReq 用户全量列表请求
// 用于下拉选择场景，支持关键词与状态过滤
type UserListReq struct {
	gmeta.Meta `method:"get" path:"/api/user/list" summary:"用户列表" permission:"sys:user:list" log:"用户管理-查询列表" logType:"query"`
	Keyword    string `json:"keyword" dc:"关键词"` // 关键词
	Status     int    `json:"status" dc:"状态"`   // 状态
}

// UserListItem 用户列表行数据
// 精简结构，仅包含列表展示所需字段
type UserListItem struct {
	Id       int64  `json:"id"`       // 用户ID
	Username string `json:"username"` // 用户名
	Realname string `json:"realname"` // 真实姓名
	Role     string `json:"role"`     // 角色
}

// UserListRes 用户列表结果
type UserListRes = []UserListItem

// UserDetailReq 用户详情请求
type UserDetailReq struct {
	gmeta.Meta `method:"get" path:"/api/user/detail/:id" summary:"用户详情" permission:"sys:user:detail" log:"用户管理-查询详情" logType:"query"`
	Id         int64 `json:"id" v:"required#缺少ID" dc:"用户ID"` // 用户ID
}

// UserDetailRes 用户详情结果
// 满足前端编辑表单契约：roles 为角色ID数组，city 为四级区划编码
type UserDetailRes struct {
	Id         int64    `json:"id"`         // 用户ID
	Username   string   `json:"username"`   // 用户名
	Realname   string   `json:"realname"`   // 真实姓名
	Avatar     string   `json:"avatar"`     // 头像
	Gender     int      `json:"gender"`     // 性别：1-男 2-女 3-保密
	GenderText string   `json:"genderText"` // 性别文本
	Mobile     string   `json:"mobile"`     // 手机号
	Email      string   `json:"email"`      // 邮箱
	Roles      []int64  `json:"roles"`      // 角色ID数组
	DeptId     int      `json:"deptId"`     // 部门ID
	LevelId    int      `json:"levelId"`    // 职级ID
	PositionId int      `json:"positionId"` // 岗位ID
	City       []string `json:"city"`       // 省市区街道编码数组
	CityInfo   string   `json:"cityInfo"`   // 省市区信息
	Address    string   `json:"address"`    // 详细地址
	Intro      string   `json:"intro"`      // 个人简介
	Status     int      `json:"status"`     // 状态：1-正常 2-禁用
	StatusText string   `json:"statusText"` // 状态文本
	Sort       int      `json:"sort"`       // 排序
	Note       string   `json:"note"`       // 备注
}

// UserAddReq 新增用户请求
type UserAddReq struct {
	gmeta.Meta `method:"post" path:"/api/user/add" summary:"新增用户" permission:"sys:user:add" log:"用户管理-新增记录" logType:"add"`
	Username   string   `json:"username" v:"required#请输入用户名|length:2,30#用户名长度需在2~30位" dc:"用户名"` // 用户名
	Password   string   `json:"password" v:"required#请输入密码|length:6,32#密码长度需在6~32位" dc:"密码"`    // 密码
	Realname   string   `json:"realname" v:"required#请输入姓名" dc:"姓名"`                            // 真实姓名
	Gender     int      `json:"gender" dc:"性别：1-男 2-女 3-保密"`                                    // 性别
	Mobile     string   `json:"mobile" v:"phone#手机号格式不正确" dc:"手机号"`                             // 手机号
	Email      string   `json:"email" v:"email#邮箱格式不正确" dc:"邮箱"`                                // 邮箱
	Roles      []int64  `json:"roles" dc:"角色ID列表"`                                              // 角色ID列表
	DeptId     int      `json:"deptId" dc:"部门ID"`                                               // 部门ID
	LevelId    int      `json:"levelId" dc:"职级ID"`                                              // 职级ID
	PositionId int      `json:"positionId" dc:"岗位ID"`                                           // 岗位ID
	City       []string `json:"city" dc:"省市区街道编码数组"`                                            // 省市区街道编码数组
	CityInfo   string   `json:"cityInfo" dc:"省市区信息"`                                            // 省市区信息
	Address    string   `json:"address" dc:"详细地址"`                                              // 详细地址
	Avatar     string   `json:"avatar" dc:"头像"`                                                 // 头像
	Intro      string   `json:"intro" dc:"个人简介"`                                                // 个人简介
	Status     int      `json:"status" dc:"状态：1-正常 2-禁用"`                                       // 状态
	Sort       int      `json:"sort" dc:"排序"`                                                   // 排序
	Note       string   `json:"note" dc:"备注"`                                                   // 备注
}

// UserUpdateReq 修改用户请求
type UserUpdateReq struct {
	gmeta.Meta `method:"put" path:"/api/user/update" summary:"修改用户" permission:"sys:user:update" log:"用户管理-修改记录" logType:"update"`
	Id         int64    `json:"id" v:"required#缺少ID" dc:"用户ID"`        // 用户ID
	Username   string   `json:"username" v:"required#请输入用户名" dc:"用户名"` // 用户名
	Password   string   `json:"password" dc:"密码（留空不修改）"`               // 密码，留空不修改
	Realname   string   `json:"realname" v:"required#请输入姓名" dc:"姓名"`   // 真实姓名
	Gender     int      `json:"gender" dc:"性别"`                        // 性别
	Mobile     string   `json:"mobile" v:"phone#手机号格式不正确" dc:"手机号"`    // 手机号
	Email      string   `json:"email" v:"email#邮箱格式不正确" dc:"邮箱"`       // 邮箱
	Roles      []int64  `json:"roles" dc:"角色ID列表"`                     // 角色ID列表
	DeptId     int      `json:"deptId" dc:"部门ID"`                      // 部门ID
	LevelId    int      `json:"levelId" dc:"职级ID"`                     // 职级ID
	PositionId int      `json:"positionId" dc:"岗位ID"`                  // 岗位ID
	City       []string `json:"city" dc:"省市区街道编码数组"`                   // 省市区街道编码数组
	CityInfo   string   `json:"cityInfo" dc:"省市区信息"`                   // 省市区信息
	Address    string   `json:"address" dc:"详细地址"`                     // 详细地址
	Avatar     string   `json:"avatar" dc:"头像"`                        // 头像
	Intro      string   `json:"intro" dc:"个人简介"`                       // 个人简介
	Status     int      `json:"status" dc:"状态"`                        // 状态
	Sort       int      `json:"sort" dc:"排序"`                          // 排序
	Note       string   `json:"note" dc:"备注"`                          // 备注
}

// UserUpdateRes 修改用户结果
type UserUpdateRes struct{}

// UserStatusReq 用户状态请求
type UserStatusReq struct {
	gmeta.Meta `method:"put" path:"/api/user/status" summary:"用户状态" permission:"sys:user:status" log:"用户管理-设置状态" logType:"update"`
	Id         int64 `json:"id" v:"required#缺少ID" dc:"用户ID"`             // 用户ID
	Status     int   `json:"status" v:"required#缺少状态" dc:"状态：1-正常 2-禁用"` // 状态
}

// UserStatusRes 用户状态结果
type UserStatusRes struct{}

// UserResetPwdReq 重置用户密码请求
type UserResetPwdReq struct {
	gmeta.Meta `method:"put" path:"/api/user/resetPwd" summary:"重置用户密码" permission:"sys:user:resetPwd" log:"用户管理-重置密码" logType:"reset"`
	UserId     int64  `json:"userId" v:"required#缺少用户ID" dc:"用户ID"`                                         // 用户ID
	Password   string `json:"password" v:"length:6,32#密码长度需在6~32位" dc:"新密码（为空时重置为系统默认密码 DEFAULT_PASSWORD）"` // 新密码
}

// UserResetPwdRes 重置用户密码结果
type UserResetPwdRes struct{}

// UserDeleteReq 删除用户请求
type UserDeleteReq struct {
	gmeta.Meta `method:"delete" path:"/api/user/delete/:id" summary:"删除用户" permission:"sys:user:delete" log:"用户管理-删除记录" logType:"delete"`
	Id         int64 `json:"id" v:"required#缺少ID" dc:"用户ID"` // 用户ID
}

// UserDeleteRes 删除用户结果
type UserDeleteRes struct{}

// UserBatchDeleteReq 批量删除用户请求
type UserBatchDeleteReq struct {
	gmeta.Meta `method:"delete" path:"/api/user/batchDelete" summary:"批量删除用户" permission:"sys:user:batchDelete" log:"用户管理-批量删除记录" logType:"delete"`
	Ids        []int64 `json:"ids" v:"required#请选择要删除的记录" dc:"用户ID列表"` // 用户ID列表
}

// UserBatchDeleteRes 批量删除结果
type UserBatchDeleteRes struct{}

// UserImportReq 导入用户请求
// 上传 xls/xlsx 文件进行数据导入
type UserImportReq struct {
	gmeta.Meta `method:"post" path:"/api/user/import" summary:"导入用户" permission:"sys:user:import" log:"用户管理-导入数据" logType:"import"`
	File       *ghttp.UploadFile `p:"file" v:"required#请选择文件" dc:"xls/xlsx 文件"` // 上传文件
}

// UserImportRes 导入结果
// 返回导入总数、成功数、失败数及失败原因
type UserImportRes struct {
	Total   int      `json:"total"`   // 导入总数
	Success int      `json:"success"` // 成功数
	Failed  int      `json:"failed"`  // 失败数
	Reasons []string `json:"reasons"` // 失败原因列表
}

// UserExportReq 导出用户请求
// 按查询条件导出符合条件的记录
type UserExportReq struct {
	gmeta.Meta `method:"get" path:"/api/user/export" summary:"导出用户" permission:"sys:user:export" log:"用户管理-导出数据" logType:"export"`
	Keyword    string `json:"keyword" dc:"关键词"`   // 关键词
	Realname   string `json:"realname" dc:"姓名模糊"` // 姓名模糊匹配
	Status     int    `json:"status" dc:"状态"`     // 状态
}

// UserExportRes 导出结果
// data 直接为导出文件的访问地址
type UserExportRes = string

// UserDocumentReq 用户打印信息请求
type UserDocumentReq struct {
	gmeta.Meta `method:"get" path:"/api/user/document/:id" summary:"用户打印信息" permission:"sys:user:detail" log:"用户管理-查询打印信息" logType:"query"`
	Id         int64 `json:"id" v:"required#缺少ID" dc:"用户ID"` // 用户ID
}

// UserDocumentRes 用户打印信息结果
// 在用户详情基础上追加关联的角色、部门、职级、岗位名称
type UserDocumentRes struct {
	UserDetailRes
	RoleNames    []string `json:"roleNames"`    // 角色名称列表
	DeptName     string   `json:"deptName"`     // 部门名称
	LevelName    string   `json:"levelName"`    // 职级名称
	PositionName string   `json:"positionName"` // 岗位名称
}
