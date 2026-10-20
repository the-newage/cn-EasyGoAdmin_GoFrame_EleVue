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

// ============================== 菜单模块 ==============================

// MenuItem 菜单行数据
// 用于菜单分页列表、菜单树形列表等场景的菜单基础信息
type MenuItem struct {
	Id         int64  `json:"id"`         // 主键ID
	ParentId   int    `json:"parentId"`   // 上级菜单ID（0=顶级）
	Name       string `json:"name"`       // 菜单名称
	Icon       string `json:"icon"`       // 图标
	Icon2      string `json:"icon2"`      // 菜单图标2
	Path       string `json:"path"`       // 路由路径
	Component  string `json:"component"`  // 前端组件路径
	Permission string `json:"permission"` // 权限节点标识
	Target     int    `json:"target"`     // 打开方式：0-组件 1-内链 2-外链
	Type       int    `json:"type"`       // 类型：0-菜单 1-权限节点
	TypeText   string `json:"typeText"`   // 类型文本
	Status     int    `json:"status"`     // 状态：0-正常 1-停用
	StatusText string `json:"statusText"` // 状态文本
	Hide       int    `json:"hide"`       // 是否可见：0-可见 1-不可见
	Sort       int    `json:"sort"`       // 排序
	Note       string `json:"note"`       // 备注
	CreateUser string `json:"createUser"` // 创建人
}

// MenuPageReq 菜单分页请求（平铺）
// 支持关键词、名称、路径、类型、状态等条件组合查询
type MenuPageReq struct {
	gmeta.Meta `method:"get" path:"/api/menu/page" summary:"菜单分页" permission:"sys:menu:page" log:"菜单管理-查询分页记录" logType:"query"`
	PageReq           // 分页基础参数
	Keyword    string `json:"keyword" dc:"关键词（名称/路径模糊）"` // 关键词，名称或路径模糊匹配
	Name       string `json:"name" dc:"菜单名称模糊"`          // 菜单名称模糊匹配
	Path       string `json:"path" dc:"路由路径模糊"`          // 路由路径模糊匹配
	Type       *int   `json:"type" dc:"类型：0-菜单 1-权限节点"`  // 类型，指针用于区分未传值
	Status     *int   `json:"status" dc:"状态：0-正常 1-停用"`  // 状态，指针用于区分未传值
}

// MenuPageRes 菜单分页结果
type MenuPageRes = PageData[MenuItem]

// MenuListReq 菜单树形列表请求
// 返回平铺列表，由前端自行构建树形结构
type MenuListReq struct {
	gmeta.Meta `method:"get" path:"/api/menu/list" summary:"菜单树形列表" permission:"sys:menu:list" log:"菜单管理-查询树形列表" logType:"query"`
	Keyword    string `json:"keyword" dc:"关键词"` // 关键词
	Type       *int   `json:"type" dc:"类型"`     // 类型
	Status     *int   `json:"status" dc:"状态"`   // 状态
}

// MenuListRes 菜单平铺列表结果
// data 直接为数组，由前端自行构建树形结构
type MenuListRes = []*MenuItem

// MenuDetailReq 菜单详情请求
type MenuDetailReq struct {
	gmeta.Meta `method:"get" path:"/api/menu/detail/:id" summary:"菜单详情" permission:"sys:menu:detail" log:"菜单管理-查询详情" logType:"query"`
	Id         int64 `json:"id" v:"required#缺少ID" dc:"菜单ID"` // 菜单ID
}

// MenuDetailRes 菜单详情结果
type MenuDetailRes struct {
	Id         int64  `json:"id"`         // 主键ID
	ParentId   int    `json:"parentId"`   // 上级菜单ID
	Name       string `json:"name"`       // 菜单名称
	Icon       string `json:"icon"`       // 图标
	Icon2      string `json:"icon2"`      // 菜单图标2
	Path       string `json:"path"`       // 路由路径
	Component  string `json:"component"`  // 前端组件路径
	Permission string `json:"permission"` // 权限节点标识
	Target     int    `json:"target"`     // 打开方式
	Type       int    `json:"type"`       // 类型：0-菜单 1-权限节点
	TypeText   string `json:"typeText"`   // 类型文本
	Status     int    `json:"status"`     // 状态：0-正常 1-停用
	StatusText string `json:"statusText"` // 状态文本
	Hide       int    `json:"hide"`       // 是否可见
	Sort       int    `json:"sort"`       // 排序
	Note       string `json:"note"`       // 备注
}

// MenuAddReq 新增菜单请求
type MenuAddReq struct {
	gmeta.Meta `method:"post" path:"/api/menu/add" summary:"新增菜单" permission:"sys:menu:add" log:"菜单管理-新增记录" logType:"add"`
	ParentId   int    `json:"parentId" dc:"上级ID（0=顶级）"`            // 上级菜单ID
	Name       string `json:"name" v:"required#请输入菜单名称" dc:"菜单名称"` // 菜单名称
	Icon       string `json:"icon" dc:"图标"`                        // 图标
	Icon2      string `json:"icon2" dc:"菜单图标2"`                    // 菜单图标2
	Path       string `json:"path" dc:"路由路径"`                      // 路由路径
	Component  string `json:"component" dc:"前端组件"`                 // 前端组件路径
	Permission string `json:"permission" dc:"权限节点"`                // 权限节点标识
	Target     int    `json:"target" dc:"打开方式：0-组件 1-内链 2-外链"`     // 打开方式
	Type       int    `json:"type" dc:"类型：0-菜单 1-权限节点"`            // 类型
	Status     int    `json:"status" dc:"状态：0-正常 1-停用"`            // 状态
	Hide       int    `json:"hide" dc:"是否可见：0-可见 1-不可见"`           // 是否可见
	Sort       int    `json:"sort" dc:"排序"`                        // 排序
	Note       string `json:"note" dc:"备注"`                        // 备注
}

// MenuUpdateReq 修改菜单请求
// Icon/Icon2/Note 为指针：未传值时该列不更新，保留库中原值；传空串表示清空
type MenuUpdateReq struct {
	gmeta.Meta `method:"put" path:"/api/menu/update" summary:"修改菜单" permission:"sys:menu:update" log:"菜单管理-修改记录" logType:"update"`
	Id         int64   `json:"id" v:"required#缺少ID" dc:"菜单ID"`      // 菜单ID
	ParentId   int     `json:"parentId" dc:"上级ID"`                  // 上级菜单ID
	Name       string  `json:"name" v:"required#请输入菜单名称" dc:"菜单名称"` // 菜单名称
	Icon       *string `json:"icon" dc:"图标"`                        // 图标，未传则保留原值
	Icon2      *string `json:"icon2" dc:"菜单图标2"`                    // 菜单图标2，未传则保留原值
	Path       string  `json:"path" dc:"路由路径"`                      // 路由路径
	Component  string  `json:"component" dc:"前端组件"`                 // 前端组件路径
	Permission string  `json:"permission" dc:"权限节点"`                // 权限节点标识
	Target     int     `json:"target" dc:"打开方式"`                    // 打开方式
	Type       int     `json:"type" dc:"类型"`                        // 类型
	Status     int     `json:"status" dc:"状态"`                      // 状态
	Hide       int     `json:"hide" dc:"是否可见"`                      // 是否可见
	Sort       int     `json:"sort" dc:"排序"`                        // 排序
	Note       *string `json:"note" dc:"备注"`                        // 备注，未传则保留原值
}

// MenuUpdateRes 修改菜单结果
type MenuUpdateRes struct{}

// MenuDeleteReq 删除菜单请求
type MenuDeleteReq struct {
	gmeta.Meta `method:"delete" path:"/api/menu/delete/:id" summary:"删除菜单" permission:"sys:menu:delete" log:"菜单管理-删除记录" logType:"delete"`
	Id         int64 `json:"id" v:"required#缺少ID" dc:"菜单ID"` // 菜单ID
}

// MenuDeleteRes 删除菜单结果
type MenuDeleteRes struct{}

// MenuBatchDeleteReq 批量删除菜单请求
type MenuBatchDeleteReq struct {
	gmeta.Meta `method:"delete" path:"/api/menu/batchDelete" summary:"批量删除菜单" permission:"sys:menu:batchDelete" log:"菜单管理-批量删除记录" logType:"delete"`
	Ids        []int64 `json:"ids" v:"required#请选择要删除的记录" dc:"菜单ID列表"` // 菜单ID列表
}

// MenuBatchDeleteRes 批量删除结果
type MenuBatchDeleteRes struct{}
