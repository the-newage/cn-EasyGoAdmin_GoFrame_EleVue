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

// ============================== 公告模块 ==============================

// NoticeItem 公告行数据
// 用于公告分页列表等场景的公告基础信息
type NoticeItem struct {
	Id         int64       `json:"id"`         // 主键ID
	Title      string      `json:"title"`      // 公告标题
	Cover      string      `json:"cover"`      // 封面图
	Type       int         `json:"type"`       // 类型
	TypeText   string      `json:"typeText"`   // 类型文本
	ClickNum   int         `json:"clickNum"`   // 点击量
	Status     int         `json:"status"`     // 状态：0-草稿 1-发布
	StatusText string      `json:"statusText"` // 状态文本
	CreateUser string      `json:"createUser"` // 创建人
	CreateTime *gtime.Time `json:"createTime"` // 创建时间
}

// NoticePageReq 公告分页请求
// 支持关键词、标题、类型、状态等条件组合查询
type NoticePageReq struct {
	gmeta.Meta `method:"get" path:"/api/notice/page" summary:"公告分页" permission:"sys:notice:page" log:"公告管理-查询分页记录" logType:"query"`
	PageReq           // 分页基础参数
	Keyword    string `json:"keyword" dc:"关键词（标题模糊）"`   // 关键词，标题模糊匹配
	Title      string `json:"title" dc:"标题模糊"`          // 标题模糊匹配
	Type       int    `json:"type" dc:"类型"`             // 类型
	Status     int    `json:"status" dc:"状态：0-草稿 1-发布"` // 状态
}

// NoticePageRes 公告分页结果
type NoticePageRes = PageData[NoticeItem]

// NoticeListReq 公告列表请求
type NoticeListReq struct {
	gmeta.Meta `method:"get" path:"/api/notice/list" summary:"公告列表" permission:"sys:notice:list" log:"公告管理-查询列表" logType:"query"`
	Keyword    string `json:"keyword" dc:"关键词"` // 关键词
	Type       int    `json:"type" dc:"类型"`     // 类型
	Status     int    `json:"status" dc:"状态"`   // 状态
	Limit      int    `json:"limit" dc:"条数上限"`  // 条数上限
}

// NoticeListItem 公告列表行数据
// 精简结构，仅包含列表展示所需字段
type NoticeListItem struct {
	Id     int64  `json:"id"`     // 主键ID
	Title  string `json:"title"`  // 公告标题
	Type   int    `json:"type"`   // 类型
	Status int    `json:"status"` // 状态：0-草稿 1-发布
}

// NoticeListRes 公告列表结果
type NoticeListRes = []NoticeListItem

// NoticeDetailReq 公告详情请求
type NoticeDetailReq struct {
	gmeta.Meta `method:"get" path:"/api/notice/detail/:id" summary:"公告详情" permission:"sys:notice:detail" log:"公告管理-查询详情" logType:"query"`
	Id         int64 `json:"id" v:"required#缺少ID" dc:"公告ID"` // 公告ID
}

// NoticeDetailRes 公告详情结果
type NoticeDetailRes struct {
	Id         int64  `json:"id"`         // 主键ID
	Title      string `json:"title"`      // 公告标题
	Cover      string `json:"cover"`      // 封面图
	Type       int    `json:"type"`       // 类型
	TypeText   string `json:"typeText"`   // 类型文本
	Status     int    `json:"status"`     // 状态：0-草稿 1-发布
	StatusText string `json:"statusText"` // 状态文本
	Content    string `json:"content"`    // 正文
}

// NoticeAddReq 新增公告请求
type NoticeAddReq struct {
	gmeta.Meta `method:"post" path:"/api/notice/add" summary:"新增公告" permission:"sys:notice:add" log:"公告管理-新增记录" logType:"add"`
	Title      string `json:"title" v:"required#请输入标题" dc:"标题"` // 公告标题
	Cover      string `json:"cover" dc:"封面图"`                   // 封面图
	Type       int    `json:"type" dc:"类型"`                     // 类型
	Status     int    `json:"status" dc:"状态：0-草稿 1-发布"`         // 状态
	Content    string `json:"content" dc:"正文"`                  // 正文
}

// NoticeUpdateReq 修改公告请求
type NoticeUpdateReq struct {
	gmeta.Meta `method:"put" path:"/api/notice/update" summary:"修改公告" permission:"sys:notice:update" log:"公告管理-修改记录" logType:"update"`
	Id         int64  `json:"id" v:"required#缺少ID" dc:"公告ID"`   // 公告ID
	Title      string `json:"title" v:"required#请输入标题" dc:"标题"` // 公告标题
	Cover      string `json:"cover" dc:"封面图"`                   // 封面图
	Type       int    `json:"type" dc:"类型"`                     // 类型
	Status     int    `json:"status" dc:"状态"`                   // 状态
	Content    string `json:"content" dc:"正文"`                  // 正文
}

// NoticeUpdateRes 修改公告结果
type NoticeUpdateRes struct{}

// NoticeStatusReq 公告状态请求
type NoticeStatusReq struct {
	gmeta.Meta `method:"put" path:"/api/notice/status" summary:"公告状态" permission:"sys:notice:status" log:"公告管理-设置状态" logType:"update"`
	Id         int64 `json:"id" v:"required#缺少ID" dc:"公告ID"`             // 公告ID
	Status     int   `json:"status" v:"required#缺少状态" dc:"状态：0-草稿 1-发布"` // 状态
}

// NoticeStatusRes 公告状态结果
type NoticeStatusRes struct{}

// NoticeDeleteReq 删除公告请求
type NoticeDeleteReq struct {
	gmeta.Meta `method:"delete" path:"/api/notice/delete/:id" summary:"删除公告" permission:"sys:notice:delete" log:"公告管理-删除记录" logType:"delete"`
	Id         int64 `json:"id" v:"required#缺少ID" dc:"公告ID"` // 公告ID
}

// NoticeDeleteRes 删除公告结果
type NoticeDeleteRes struct{}

// NoticeBatchDeleteReq 批量删除公告请求
type NoticeBatchDeleteReq struct {
	gmeta.Meta `method:"delete" path:"/api/notice/batchDelete" summary:"批量删除公告" permission:"sys:notice:batchDelete" log:"公告管理-批量删除记录" logType:"delete"`
	Ids        []int64 `json:"ids" v:"required#请选择要删除的记录" dc:"公告ID列表"` // 公告ID列表
}

// NoticeBatchDeleteRes 批量删除结果
type NoticeBatchDeleteRes struct{}
