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

// ============================== 案例演示模块（goframe_example） ==============================

// ExampleItem 案例演示行数据
// 用于分页列表等场景的基础信息
type ExampleItem struct {
	Id         int         `json:"id"`         // 主键ID
	Name       string      `json:"name"`       // 案例名称
	Avatar     string      `json:"avatar"`     // 案例图片
	Type       int         `json:"type"`       // 案例类型
	Status     int         `json:"status"`     // 案例状态
	Sort       int         `json:"sort"`       // 案例排序
	CreateUser string      `json:"createUser"` // 创建人
	CreateTime *gtime.Time `json:"createTime"` // 创建时间
	StatusText string      `json:"statusText"` // 状态显示名（dict_item 解析）
}

// ExamplePageReq 案例演示分页请求
type ExamplePageReq struct {
	gmeta.Meta `method:"get" path:"/api/example/page" summary:"案例演示分页" permission:"sys:example:page" log:"案例演示-查询分页记录" logType:"query"`
	PageReq
	Name    string `json:"name" dc:"案例名称"`   // 案例名称
	Type    int    `json:"type" dc:"案例类型"`   // 案例类型
	Status  int    `json:"status" dc:"案例状态"` // 案例状态
	Keyword string `json:"keyword" dc:"关键词"` // 关键词
}

// ExamplePageRes 案例演示分页结果
// 统一壳：records/total/size/current/pages
type ExamplePageRes = PageData[ExampleItem]

// ExampleListReq 案例演示列表请求
type ExampleListReq struct {
	gmeta.Meta `method:"get" path:"/api/example/list" summary:"案例演示列表" permission:"sys:example:list" log:"案例演示-查询列表" logType:"query"`
	Keyword    string `json:"keyword" dc:"关键词"` // 关键词
	Limit      int    `json:"limit" dc:"条数上限"`  // 条数上限
}

// ExampleListRes 案例演示平铺列表结果
// data 直接为数组，树形结构由前端 buildTree 构建
type ExampleListRes = []ExampleItem

// ExampleDetailReq 案例演示详情请求
type ExampleDetailReq struct {
	gmeta.Meta `method:"get" path:"/api/example/detail/:id" summary:"案例演示详情" permission:"sys:example:detail" log:"案例演示-查询详情" logType:"query"`
	Id         int64 `json:"id" v:"required#缺少ID" dc:"记录ID"` // 记录ID
}

// ExampleDetailRes 案例演示详情结果
type ExampleDetailRes struct {
	Id         int64  `json:"id"`         // 主键ID
	Name       string `json:"name"`       // 案例名称
	Avatar     string `json:"avatar"`     // 案例图片
	Type       int    `json:"type"`       // 案例类型
	Status     int    `json:"status"`     // 案例状态
	Sort       int    `json:"sort"`       // 案例排序
	StatusText string `json:"statusText"` // 状态显示名（dict_item 解析）
}

// ExampleAddReq 新增案例演示请求
type ExampleAddReq struct {
	gmeta.Meta `method:"post" path:"/api/example/add" summary:"新增案例演示" permission:"sys:example:add" log:"案例演示-新增记录" logType:"add"`
	Name       string `json:"name" v:"required#请输入案例名称" dc:"案例名称"` // 案例名称
	Avatar     string `json:"avatar" dc:"案例图片"`                    // 案例图片
	Type       int    `json:"type" dc:"案例类型"`                      // 案例类型
	Status     int    `json:"status" dc:"案例状态"`                    // 案例状态
	Sort       int    `json:"sort" dc:"案例排序"`                      // 案例排序
}

// ExampleIdRes 新增返回ID结构
type ExampleIdRes struct {
	Id int64 `json:"id"` // 主键ID
}

// ExampleUpdateReq 修改案例演示请求
type ExampleUpdateReq struct {
	gmeta.Meta `method:"put" path:"/api/example/update" summary:"修改案例演示" permission:"sys:example:update" log:"案例演示-修改记录" logType:"update"`
	Id         int64  `json:"id" v:"required#缺少ID" dc:"记录ID"`      // 记录ID
	Name       string `json:"name" v:"required#请输入案例名称" dc:"案例名称"` // 案例名称
	Avatar     string `json:"avatar" dc:"案例图片"`                    // 案例图片
	Type       int    `json:"type" dc:"案例类型"`                      // 案例类型
	Status     int    `json:"status" dc:"案例状态"`                    // 案例状态
	Sort       int    `json:"sort" dc:"案例排序"`                      // 案例排序
}

// ExampleUpdateRes 修改案例演示结果
type ExampleUpdateRes struct{}

// ExampleStatusReq 案例演示状态请求
type ExampleStatusReq struct {
	gmeta.Meta `method:"put" path:"/api/example/status" summary:"案例演示状态" permission:"sys:example:status" log:"案例演示-设置状态" logType:"update"`
	Id         int64 `json:"id" v:"required#缺少ID" dc:"记录ID"`   // 记录ID
	Status     int   `json:"status" v:"required#缺少状态" dc:"状态"` // 状态
}

// ExampleStatusRes 案例演示状态结果
type ExampleStatusRes struct{}

// ExampleDeleteReq 删除案例演示请求
type ExampleDeleteReq struct {
	gmeta.Meta `method:"delete" path:"/api/example/delete/:id" summary:"删除案例演示" permission:"sys:example:delete" log:"案例演示-删除记录" logType:"delete"`
	Id         int64 `json:"id" v:"required#缺少ID" dc:"记录ID"` // 记录ID
}

// ExampleDeleteRes 删除案例演示结果
type ExampleDeleteRes struct{}

// ExampleBatchDeleteReq 批量删除案例演示请求
type ExampleBatchDeleteReq struct {
	gmeta.Meta `method:"delete" path:"/api/example/batchDelete" summary:"批量删除案例演示" permission:"sys:example:batchDelete" log:"案例演示-批量删除记录" logType:"delete"`
	Ids        []int64 `json:"ids" v:"required#请选择要删除的记录" dc:"ID列表"` // ID列表
}

// ExampleBatchDeleteRes 批量删除结果
type ExampleBatchDeleteRes struct{}

// ExampleImportReq 导入案例演示请求
type ExampleImportReq struct {
	gmeta.Meta `method:"post" path:"/api/example/import" summary:"导入案例演示" permission:"sys:example:import" log:"案例演示-导入数据" logType:"import"`
	File       *ghttp.UploadFile `p:"file" v:"required#请选择文件" dc:"xls/xlsx 文件"` // 上传文件
}

// ExampleImportRes 导入结果
// 返回导入总数、成功数、失败数及失败原因
type ExampleImportRes struct {
	Total   int      `json:"total"`   // 导入总数
	Success int      `json:"success"` // 成功数
	Failed  int      `json:"failed"`  // 失败数
	Reasons []string `json:"reasons"` // 失败原因列表
}

// ExampleExportReq 导出案例演示请求
type ExampleExportReq struct {
	gmeta.Meta `method:"get" path:"/api/example/export" summary:"导出案例演示" permission:"sys:example:export" log:"案例演示-导出数据" logType:"export"`
	Keyword    string `json:"keyword" dc:"关键词"` // 关键词
	Name       string `json:"name" dc:"案例名称"`   // 案例名称
	Type       int    `json:"type" dc:"案例类型"`   // 案例类型
	Status     int    `json:"status" dc:"案例状态"` // 案例状态
}

// ExampleExportRes 导出结果
// data 直接为文件访问 URL 字符串，前端可 window.open(data)
type ExampleExportRes = string
