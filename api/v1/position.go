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

// ============================== 岗位模块 ==============================

// PositionItem 岗位行数据
// 用于岗位分页列表等场景的岗位基础信息
type PositionItem struct {
	Id         int64       `json:"id"`         // 主键ID
	TenantId   int         `json:"tenantId"`   // 租户ID
	Name       string      `json:"name"`       // 岗位名称
	Status     int         `json:"status"`     // 状态：1-正常 2-禁用
	StatusText string      `json:"statusText"` // 状态文本
	Sort       int         `json:"sort"`       // 排序
	UserCount  int         `json:"userCount"`  // 关联用户数
	CreateUser string      `json:"createUser"` // 创建人
	CreateTime *gtime.Time `json:"createTime"` // 创建时间
}

// PositionPageReq 岗位分页请求
// 支持关键词、岗位名称、状态等条件组合查询
type PositionPageReq struct {
	gmeta.Meta `method:"get" path:"/api/position/page" summary:"岗位分页" permission:"sys:position:page" log:"岗位管理-查询分页记录" logType:"query"`
	PageReq           // 分页基础参数
	Keyword    string `json:"keyword" dc:"关键词"`         // 关键词
	Name       string `json:"name" dc:"岗位名称模糊"`         // 岗位名称模糊匹配
	Status     int    `json:"status" dc:"状态：1-正常 2-禁用"` // 状态
}

// PositionPageRes 岗位分页结果
type PositionPageRes = PageData[PositionItem]

// PositionListReq 岗位全量列表请求
type PositionListReq struct {
	gmeta.Meta `method:"get" path:"/api/position/list" summary:"岗位列表" permission:"sys:position:list" log:"岗位管理-查询列表" logType:"query"`
	Keyword    string `json:"keyword" dc:"关键词"` // 关键词
	Status     int    `json:"status" dc:"状态"`   // 状态
}

// PositionListItem 岗位列表行数据
// 精简结构，仅包含列表展示所需字段
type PositionListItem struct {
	Id     int64  `json:"id"`     // 主键ID
	Name   string `json:"name"`   // 岗位名称
	Status int    `json:"status"` // 状态：1-正常 2-禁用
	Sort   int    `json:"sort"`   // 排序
}

// PositionListRes 岗位列表结果
type PositionListRes = []PositionListItem

// PositionDetailReq 岗位详情请求
type PositionDetailReq struct {
	gmeta.Meta `method:"get" path:"/api/position/detail/:id" summary:"岗位详情" permission:"sys:position:detail" log:"岗位管理-查询详情" logType:"query"`
	Id         int64 `json:"id" v:"required#缺少ID" dc:"岗位ID"` // 岗位ID
}

// PositionDetailRes 岗位详情结果
type PositionDetailRes struct {
	Id         int64  `json:"id"`         // 主键ID
	Name       string `json:"name"`       // 岗位名称
	Status     int    `json:"status"`     // 状态：1-正常 2-禁用
	StatusText string `json:"statusText"` // 状态文本
	Sort       int    `json:"sort"`       // 排序
}

// PositionAddReq 新增岗位请求
type PositionAddReq struct {
	gmeta.Meta `method:"post" path:"/api/position/add" summary:"新增岗位" permission:"sys:position:add" log:"岗位管理-新增记录" logType:"add"`
	Name       string `json:"name" v:"required#请输入岗位名称" dc:"岗位名称"` // 岗位名称
	Status     int    `json:"status" dc:"状态：1-正常 2-禁用"`            // 状态
	Sort       int    `json:"sort" dc:"排序"`                        // 排序
}

// PositionUpdateReq 修改岗位请求
type PositionUpdateReq struct {
	gmeta.Meta `method:"put" path:"/api/position/update" summary:"修改岗位" permission:"sys:position:update" log:"岗位管理-修改记录" logType:"update"`
	Id         int64  `json:"id" v:"required#缺少ID" dc:"岗位ID"`      // 岗位ID
	Name       string `json:"name" v:"required#请输入岗位名称" dc:"岗位名称"` // 岗位名称
	Status     int    `json:"status" dc:"状态"`                      // 状态
	Sort       int    `json:"sort" dc:"排序"`                        // 排序
}

// PositionUpdateRes 修改岗位结果
type PositionUpdateRes struct{}

// PositionStatusReq 岗位状态请求
type PositionStatusReq struct {
	gmeta.Meta `method:"put" path:"/api/position/status" summary:"岗位状态" permission:"sys:position:status" log:"岗位管理-设置状态" logType:"update"`
	Id         int64 `json:"id" v:"required#缺少ID" dc:"岗位ID"`             // 岗位ID
	Status     int   `json:"status" v:"required#缺少状态" dc:"状态：1-正常 2-禁用"` // 状态
}

// PositionStatusRes 岗位状态结果
type PositionStatusRes struct{}

// PositionDeleteReq 删除岗位请求
type PositionDeleteReq struct {
	gmeta.Meta `method:"delete" path:"/api/position/delete/:id" summary:"删除岗位" permission:"sys:position:delete" log:"岗位管理-删除记录" logType:"delete"`
	Id         int64 `json:"id" v:"required#缺少ID" dc:"岗位ID"` // 岗位ID
}

// PositionDeleteRes 删除岗位结果
type PositionDeleteRes struct{}

// PositionBatchDeleteReq 批量删除岗位请求
type PositionBatchDeleteReq struct {
	gmeta.Meta `method:"delete" path:"/api/position/batchDelete" summary:"批量删除岗位" permission:"sys:position:batchDelete" log:"岗位管理-批量删除记录" logType:"delete"`
	Ids        []int64 `json:"ids" v:"required#请选择要删除的记录" dc:"岗位ID列表"` // 岗位ID列表
}

// PositionBatchDeleteRes 批量删除结果
type PositionBatchDeleteRes struct{}
