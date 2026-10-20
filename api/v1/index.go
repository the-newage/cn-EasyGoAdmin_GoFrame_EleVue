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

// ============================== 首页/个人中心模块 ==============================

// GetUserReq 获取当前登录用户信息请求
// 返回用户基础信息及其角色与权限集合
type GetUserReq struct {
	gmeta.Meta `method:"get" path:"/api/index/getUser" summary:"获取当前用户信息"`
}

// GetUserRes 当前登录用户信息结果
// 返回用户基础信息、角色与权限集合，敏感字段不返回
type GetUserRes struct {
	Id          int64    `json:"id"`          // 用户ID
	Username    string   `json:"username"`    // 用户名
	Realname    string   `json:"realname"`    // 真实姓名
	Gender      int      `json:"gender"`      // 性别：1-男 2-女 3-保密
	Avatar      string   `json:"avatar"`      // 头像
	Mobile      string   `json:"mobile"`      // 手机号
	Email       string   `json:"email"`       // 邮箱
	DeptId      int      `json:"deptId"`      // 部门ID
	LevelId     int      `json:"levelId"`     // 职级ID
	PositionId  int      `json:"positionId"`  // 岗位ID
	City        []string `json:"city"`        // 省市区街道编码数组
	CityInfo    string   `json:"cityInfo"`    // 省市区信息
	Address     string   `json:"address"`     // 详细地址
	Intro       string   `json:"intro"`       // 个人简介
	Status      int      `json:"status"`      // 状态
	Note        string   `json:"note"`        // 备注
	TenantId    int      `json:"tenantId"`    // 租户ID
	Roles       []string `json:"roles"`       // 角色编码集合
	Permissions []string `json:"permissions"` // 权限标识集合
}

// GetMenusReq 获取当前用户菜单请求
// 无路径参数版本，与带路径参数的菜单路由共用同一处理逻辑
type GetMenusReq struct {
	gmeta.Meta `method:"get" path:"/api/index/getMenus" summary:"获取当前用户菜单"`
	Id         int64 `json:"id" dc:"菜单ID（路由占位，忽略）"` // 菜单ID，路由占位
}

// GetMenusByIdReq 获取当前用户菜单请求（带路径参数）
// 带路径参数版本，与无路径参数的菜单路由共用同一处理逻辑
type GetMenusByIdReq struct {
	gmeta.Meta `method:"get" path:"/api/index/getMenus/:id" summary:"获取当前用户菜单"`
	Id         int64 `json:"id" dc:"菜单ID（路由占位，忽略）"` // 菜单ID，路由占位
}

// MenuTreeItem 菜单树节点
// 字段与菜单表及前端路由生成契约保持一致
type MenuTreeItem struct {
	Id         int64           `json:"id"`                 // 菜单ID
	ParentId   int             `json:"parentId"`           // 上级菜单ID
	Name       string          `json:"name"`               // 菜单名称
	Icon       string          `json:"icon"`               // 图标
	Icon2      string          `json:"icon2"`              // 菜单图标2
	Path       string          `json:"path"`               // 路由路径
	Component  string          `json:"component"`          // 前端组件路径
	Permission string          `json:"permission"`         // 权限标识
	Target     int             `json:"target"`             // 打开方式
	Type       int             `json:"type"`               // 类型：0-菜单 1-节点
	Status     int             `json:"status"`             // 状态
	Hide       int             `json:"hide"`               // 是否隐藏
	Sort       int             `json:"sort"`               // 排序
	Redirect   string          `json:"redirect,omitempty"` // 重定向地址
	Children   []*MenuTreeItem `json:"children"`           // 子节点集合
}

// GetMenusRes 菜单树结果
type GetMenusRes = []*MenuTreeItem

// UpdatePasswordReq 修改个人密码请求
type UpdatePasswordReq struct {
	gmeta.Meta  `method:"put" path:"/api/index/updatePassword" summary:"修改个人密码"`
	Password    string `json:"password" v:"required#请输入原密码" dc:"原密码"`                             // 原密码
	NewPassword string `json:"newPassword" v:"required#请输入新密码|length:6,32#新密码长度需在6~32位" dc:"新密码"` // 新密码
}

// UpdatePasswordRes 修改个人密码结果
type UpdatePasswordRes struct{}

// UpdateProfileReq 更新个人信息请求
type UpdateProfileReq struct {
	gmeta.Meta `method:"put" path:"/api/index/updateProfile" summary:"更新个人信息"`
	Realname   string   `json:"realname" v:"required#请输入姓名" dc:"姓名"` // 真实姓名
	Gender     int      `json:"gender" dc:"性别：1-男 2-女 3-保密"`         // 性别
	Avatar     string   `json:"avatar" dc:"头像"`                      // 头像
	Mobile     string   `json:"mobile" v:"phone#手机号格式不正确" dc:"手机号"`  // 手机号
	Email      string   `json:"email" v:"email#邮箱格式不正确" dc:"邮箱"`     // 邮箱
	City       []string `json:"city" dc:"省市区街道编码数组"`                 // 省市区街道编码数组
	CityInfo   string   `json:"cityInfo" dc:"省市区信息"`                 // 省市区信息
	Address    string   `json:"address" dc:"详细地址"`                   // 详细地址
	Intro      string   `json:"intro" dc:"个人简介"`                     // 个人简介
}

// UpdateProfileRes 更新个人信息结果
type UpdateProfileRes struct{}
