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

// ============================== 角色模块 ==============================

// RolePageReq 角色分页请求
// 支持关键词、角色名称等条件组合查询
type RolePageReq struct {
	gmeta.Meta `method:"get" path:"/api/role/page" summary:"角色分页" permission:"sys:role:page" log:"角色管理-查询分页记录" logType:"query"`
	PageReq           // 分页基础参数
	Keyword    string `json:"keyword" dc:"关键词（名称/编码模糊）"` // 关键词，名称或编码模糊匹配
	Name       string `json:"name" dc:"角色名称模糊"`          // 角色名称模糊匹配
}

// RoleItem 角色行数据
// 用于角色分页列表等场景的角色基础信息
type RoleItem struct {
	Id            int64       `json:"id"`            // 主键ID
	TenantId      int         `json:"tenantId"`      // 租户ID
	Name          string      `json:"name"`          // 角色名称
	Code          string      `json:"code"`          // 角色编码
	DataScope     int         `json:"dataScope"`     // 数据范围：1-全部 2-本部门及以下 3-本部门 4-仅本人
	DataScopeText string      `json:"dataScopeText"` // 数据范围文案（字典驱动）
	Sort          int         `json:"sort"`          // 排序
	Note          string      `json:"note"`          // 备注
	UserCount     int         `json:"userCount"`     // 关联用户数
	CreateUser    string      `json:"createUser"`    // 创建人
	CreateTime    *gtime.Time `json:"createTime"`    // 创建时间
}

// RolePageRes 角色分页结果
type RolePageRes = PageData[RoleItem]

// RoleListReq 角色全量列表请求
type RoleListReq struct {
	gmeta.Meta `method:"get" path:"/api/role/list" summary:"角色列表" permission:"sys:role:list" log:"角色管理-查询列表" logType:"query"`
	Keyword    string `json:"keyword" dc:"关键词"` // 关键词
}

// RoleListItem 角色列表行数据
// 精简结构，仅包含列表展示所需字段
type RoleListItem struct {
	Id   int64  `json:"id"`   // 主键ID
	Name string `json:"name"` // 角色名称
	Code string `json:"code"` // 角色编码
	Sort int    `json:"sort"` // 排序
}

// RoleListRes 角色列表结果
type RoleListRes = []RoleListItem

// RoleDetailReq 角色详情请求
type RoleDetailReq struct {
	gmeta.Meta `method:"get" path:"/api/role/detail/:id" summary:"角色详情" permission:"sys:role:detail" log:"角色管理-查询详情" logType:"query"`
	Id         int64 `json:"id" v:"required#缺少ID" dc:"角色ID"` // 角色ID
}

// RoleDetailRes 角色详情结果
type RoleDetailRes struct {
	Id        int64  `json:"id"`        // 主键ID
	TenantId  int    `json:"tenantId"`  // 租户ID
	Name      string `json:"name"`      // 角色名称
	Code      string `json:"code"`      // 角色编码
	DataScope int    `json:"dataScope"` // 数据范围：1-全部 2-本部门及以下 3-本部门 4-仅本人
	Sort      int    `json:"sort"`      // 排序
	Note      string `json:"note"`      // 备注
}

// RoleAddReq 新增角色请求
type RoleAddReq struct {
	gmeta.Meta `method:"post" path:"/api/role/add" summary:"新增角色" permission:"sys:role:add" log:"角色管理-新增记录" logType:"add"`
	Name       string `json:"name" v:"required#请输入角色名称" dc:"角色名称"`                                        // 角色名称
	Code       string `json:"code" v:"required#请输入角色编码|length:2,30#编码长度需在2~30位" dc:"角色编码"`                // 角色编码
	DataScope  int    `json:"dataScope" v:"in:1,2,3,4#数据范围不合法" d:"1" dc:"数据范围：1-全部 2-本部门及以下 3-本部门 4-仅本人"` // 数据范围
	Sort       int    `json:"sort" dc:"排序"`                                                               // 排序
	Note       string `json:"note" dc:"备注"`                                                               // 备注
}

// RoleUpdateReq 修改角色请求
type RoleUpdateReq struct {
	gmeta.Meta `method:"put" path:"/api/role/update" summary:"修改角色" permission:"sys:role:update" log:"角色管理-修改记录" logType:"update"`
	Id         int64  `json:"id" v:"required#缺少ID" dc:"角色ID"`                                       // 角色ID
	Name       string `json:"name" v:"required#请输入角色名称" dc:"角色名称"`                                  // 角色名称
	Code       string `json:"code" dc:"角色编码（不允许修改）"`                                                // 角色编码，不允许修改
	DataScope  int    `json:"dataScope" v:"in:1,2,3,4#数据范围不合法" dc:"数据范围：1-全部 2-本部门及以下 3-本部门 4-仅本人"` // 数据范围
	Sort       int    `json:"sort" dc:"排序"`                                                         // 排序
	Note       string `json:"note" dc:"备注"`                                                         // 备注
}

// RoleUpdateRes 修改角色结果
type RoleUpdateRes struct{}

// RoleDeleteReq 删除角色请求
type RoleDeleteReq struct {
	gmeta.Meta `method:"delete" path:"/api/role/delete/:id" summary:"删除角色" permission:"sys:role:delete" log:"角色管理-删除记录" logType:"delete"`
	Id         int64 `json:"id" v:"required#缺少ID" dc:"角色ID"` // 角色ID
}

// RoleDeleteRes 删除角色结果
type RoleDeleteRes struct{}

// RoleBatchDeleteReq 批量删除角色请求
type RoleBatchDeleteReq struct {
	gmeta.Meta `method:"delete" path:"/api/role/batchDelete" summary:"批量删除角色" permission:"sys:role:batchDelete" log:"角色管理-批量删除记录" logType:"delete"`
	Ids        []int64 `json:"ids" v:"required#请选择要删除的记录" dc:"角色ID列表"` // 角色ID列表
}

// RoleBatchDeleteRes 批量删除结果
type RoleBatchDeleteRes struct{}

// RoleMenuListReq 角色菜单权限请求
// 返回菜单平铺列表并附带勾选状态
type RoleMenuListReq struct {
	gmeta.Meta `method:"get" path:"/api/role/menu/list/:roleId" summary:"查询角色菜单权限" permission:"sys:role:getPermission" log:"角色管理-查询角色菜单" logType:"query"`
	RoleId     int64 `json:"roleId" v:"required#缺少角色ID" dc:"角色ID"` // 角色ID
}

// RoleMenuItem 菜单勾选行数据
// 供前端构建树形结构并依据 checked 标记勾选状态
type RoleMenuItem struct {
	Id         int64  `json:"id"`         // 菜单ID
	ParentId   int    `json:"parentId"`   // 上级菜单ID
	Name       string `json:"name"`       // 菜单名称
	Icon       string `json:"icon"`       // 图标
	Icon2      string `json:"icon2"`      // 菜单图标2
	Path       string `json:"path"`       // 路由路径
	Component  string `json:"component"`  // 前端组件路径
	Permission string `json:"permission"` // 权限节点标识
	Target     int    `json:"target"`     // 打开方式
	Type       int    `json:"type"`       // 类型：0-菜单 1-权限节点
	Status     int    `json:"status"`     // 状态
	Hide       int    `json:"hide"`       // 是否可见
	Sort       int    `json:"sort"`       // 排序
	Checked    bool   `json:"checked"`    // 是否已勾选
}

// RoleMenuListRes 角色菜单勾选列表结果（平铺）
type RoleMenuListRes = []RoleMenuItem

// RoleSavePermissionReq 保存角色菜单权限请求
type RoleSavePermissionReq struct {
	gmeta.Meta `method:"post" path:"/api/role/menu/save" summary:"保存角色菜单权限" permission:"sys:role:savePermission" log:"角色管理-保存权限" logType:"update"`
	RoleId     int64   `json:"roleId" v:"required#缺少角色ID" dc:"角色ID"` // 角色ID
	MenuIds    []int64 `json:"menuIds" dc:"菜单ID列表（含半选父级）"`           // 菜单ID列表
}

// RoleSavePermissionRes 保存角色菜单权限结果
type RoleSavePermissionRes struct{}
