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

// ============================== 职级模块 ==============================

// LevelItem 职级行数据
// 用于职级分页列表等场景的职级基础信息
type LevelItem struct {
	Id         int64       `json:"id"`         // 主键ID
	TenantId   int         `json:"tenantId"`   // 租户ID
	Name       string      `json:"name"`       // 职级名称
	Status     int         `json:"status"`     // 状态：1-正常 2-禁用
	StatusText string      `json:"statusText"` // 状态文本
	Sort       int         `json:"sort"`       // 排序
	UserCount  int         `json:"userCount"`  // 关联用户数
	CreateUser string      `json:"createUser"` // 创建人
	CreateTime *gtime.Time `json:"createTime"` // 创建时间
}

// LevelPageReq 职级分页请求
// 支持关键词、名称、状态等条件组合查询
type LevelPageReq struct {
	gmeta.Meta `method:"get" path:"/api/level/page" summary:"职级分页" permission:"sys:level:page" log:"职级管理-查询分页记录" logType:"query"`
	PageReq           // 分页基础参数
	Keyword    string `json:"keyword" dc:"关键词"`         // 关键词
	Name       string `json:"name" dc:"名称模糊"`           // 名称模糊匹配
	Status     int    `json:"status" dc:"状态：1-正常 2-禁用"` // 状态
}

// LevelPageRes 职级分页结果
type LevelPageRes = PageData[LevelItem]

// LevelListReq 职级全量列表请求
type LevelListReq struct {
	gmeta.Meta `method:"get" path:"/api/level/list" summary:"职级列表" permission:"sys:level:list" log:"职级管理-查询列表" logType:"query"`
	Keyword    string `json:"keyword" dc:"关键词"` // 关键词
	Status     int    `json:"status" dc:"状态"`   // 状态
}

// LevelListItem 职级列表行数据
// 精简结构，仅包含列表展示所需字段
type LevelListItem struct {
	Id     int64  `json:"id"`     // 主键ID
	Name   string `json:"name"`   // 职级名称
	Status int    `json:"status"` // 状态：1-正常 2-禁用
	Sort   int    `json:"sort"`   // 排序
}

// LevelListRes 职级列表结果
type LevelListRes = []LevelListItem

// LevelDetailReq 职级详情请求
type LevelDetailReq struct {
	gmeta.Meta `method:"get" path:"/api/level/detail/:id" summary:"职级详情" permission:"sys:level:detail" log:"职级管理-查询详情" logType:"query"`
	Id         int64 `json:"id" v:"required#缺少ID" dc:"职级ID"` // 职级ID
}

// LevelDetailRes 职级详情结果
type LevelDetailRes struct {
	Id         int64  `json:"id"`         // 主键ID
	Name       string `json:"name"`       // 职级名称
	Status     int    `json:"status"`     // 状态：1-正常 2-禁用
	StatusText string `json:"statusText"` // 状态文本
	Sort       int    `json:"sort"`       // 排序
}

// LevelAddReq 新增职级请求
type LevelAddReq struct {
	gmeta.Meta `method:"post" path:"/api/level/add" summary:"新增职级" permission:"sys:level:add" log:"职级管理-新增记录" logType:"add"`
	Name       string `json:"name" v:"required#请输入职级名称" dc:"职级名称"` // 职级名称
	Status     int    `json:"status" dc:"状态：1-正常 2-禁用"`            // 状态
	Sort       int    `json:"sort" dc:"排序"`                        // 排序
}

// LevelUpdateReq 修改职级请求
type LevelUpdateReq struct {
	gmeta.Meta `method:"put" path:"/api/level/update" summary:"修改职级" permission:"sys:level:update" log:"职级管理-修改记录" logType:"update"`
	Id         int64  `json:"id" v:"required#缺少ID" dc:"职级ID"`      // 职级ID
	Name       string `json:"name" v:"required#请输入职级名称" dc:"职级名称"` // 职级名称
	Status     int    `json:"status" dc:"状态"`                      // 状态
	Sort       int    `json:"sort" dc:"排序"`                        // 排序
}

// LevelUpdateRes 修改职级结果
type LevelUpdateRes struct{}

// LevelStatusReq 职级状态请求
type LevelStatusReq struct {
	gmeta.Meta `method:"put" path:"/api/level/status" summary:"职级状态" permission:"sys:level:status" log:"职级管理-设置状态" logType:"update"`
	Id         int64 `json:"id" v:"required#缺少ID" dc:"职级ID"`             // 职级ID
	Status     int   `json:"status" v:"required#缺少状态" dc:"状态：1-正常 2-禁用"` // 状态
}

// LevelStatusRes 职级状态结果
type LevelStatusRes struct{}

// LevelDeleteReq 删除职级请求
type LevelDeleteReq struct {
	gmeta.Meta `method:"delete" path:"/api/level/delete/:id" summary:"删除职级" permission:"sys:level:delete" log:"职级管理-删除记录" logType:"delete"`
	Id         int64 `json:"id" v:"required#缺少ID" dc:"职级ID"` // 职级ID
}

// LevelDeleteRes 删除职级结果
type LevelDeleteRes struct{}

// LevelBatchDeleteReq 批量删除职级请求
type LevelBatchDeleteReq struct {
	gmeta.Meta `method:"delete" path:"/api/level/batchDelete" summary:"批量删除职级" permission:"sys:level:batchDelete" log:"职级管理-批量删除记录" logType:"delete"`
	Ids        []int64 `json:"ids" v:"required#请选择要删除的记录" dc:"职级ID列表"` // 职级ID列表
}

// LevelBatchDeleteRes 批量删除结果
type LevelBatchDeleteRes struct{}

// LevelImportReq 导入职级请求
// 上传 xls/xlsx 文件进行数据导入
type LevelImportReq struct {
	gmeta.Meta `method:"post" path:"/api/level/import" summary:"导入职级" permission:"sys:level:import" log:"职级管理-导入数据" logType:"import"`
	File       *ghttp.UploadFile `p:"file" v:"required#请选择文件" dc:"xls/xlsx 文件"` // 上传文件
}

// LevelImportRes 导入结果
// 返回导入总数、成功数、失败数及失败原因
type LevelImportRes struct {
	Total   int      `json:"total"`   // 导入总数
	Success int      `json:"success"` // 成功数
	Failed  int      `json:"failed"`  // 失败数
	Reasons []string `json:"reasons"` // 失败原因列表
}

// LevelExportReq 导出职级请求
// 按查询条件导出符合条件的记录
type LevelExportReq struct {
	gmeta.Meta `method:"get" path:"/api/level/export" summary:"导出职级" permission:"sys:level:export" log:"职级管理-导出数据" logType:"export"`
	Keyword    string `json:"keyword" dc:"关键词"` // 关键词
	Name       string `json:"name" dc:"名称模糊"`   // 名称模糊匹配
	Status     int    `json:"status" dc:"状态"`   // 状态
}

// LevelExportRes 导出结果
// data 直接为导出文件的访问地址
type LevelExportRes = string
