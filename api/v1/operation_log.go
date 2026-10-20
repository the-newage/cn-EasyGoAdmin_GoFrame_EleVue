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

// ============================== 操作日志模块 ==============================

// OperLogItem 操作日志行数据
// 用于操作日志分页列表、操作日志列表等场景的日志基础信息
type OperLogItem struct {
	Id            int64       `json:"id"`            // 主键ID
	Title         string      `json:"title"`         // 日志标题
	Type          int         `json:"type"`          // 类型
	TypeText      string      `json:"typeText"`      // 类型文本
	Source        int         `json:"source"`        // 来源
	SourceText    string      `json:"sourceText"`    // 来源文本
	Method        string      `json:"method"`        // 处理方法
	RequestMethod string      `json:"requestMethod"` // 请求方式
	Url           string      `json:"url"`           // 请求地址
	Param         string      `json:"param"`         // 请求参数
	Result        string      `json:"result"`        // 返回结果
	Ip            string      `json:"ip"`            // 客户端IP
	Location      string      `json:"location"`      // 归属地
	Os            string      `json:"os"`            // 操作系统
	Browser       string      `json:"browser"`       // 浏览器
	ConsumeTime   int64       `json:"consumeTime"`   // 耗时（毫秒）
	Status        int         `json:"status"`        // 状态：0-正常 1-异常
	StatusText    string      `json:"statusText"`    // 状态文本
	Error         string      `json:"error"`         // 错误信息
	CreateUser    string      `json:"createUser"`    // 操作人
	CreateTime    *gtime.Time `json:"createTime"`    // 创建时间
	UpdateTime    *gtime.Time `json:"updateTime"`    // 更新时间
}

// OperLogPageReq 操作日志分页请求
// 支持关键词、日志标题、类型、状态等条件组合查询
type OperLogPageReq struct {
	gmeta.Meta `method:"get" path:"/api/oper/log/page" summary:"操作日志分页" permission:"sys:operLog:page" log:"操作日志-查询分页记录" logType:"query"`
	PageReq           // 分页基础参数
	Keyword    string `json:"keyword" dc:"关键词（标题/URL模糊）"` // 关键词，标题或URL模糊匹配
	Title      string `json:"title" dc:"日志标题模糊"`          // 日志标题模糊匹配
	Type       int    `json:"type" dc:"类型（logType）"`      // 类型
	Status     int    `json:"status" dc:"状态：0-正常 1-异常"`   // 状态
}

// OperLogPageRes 操作日志分页结果
type OperLogPageRes = PageData[OperLogItem]

// OperLogListReq 操作日志列表请求
type OperLogListReq struct {
	gmeta.Meta `method:"get" path:"/api/oper/log/list" summary:"操作日志列表" permission:"sys:operLog:list" log:"操作日志-查询列表" logType:"query"`
	Keyword    string `json:"keyword" dc:"关键词"` // 关键词
	Type       int    `json:"type" dc:"类型"`     // 类型
	Status     int    `json:"status" dc:"状态"`   // 状态
	Limit      int    `json:"limit" dc:"条数上限"`  // 条数上限
}

// OperLogListRes 操作日志列表结果
type OperLogListRes = []OperLogItem

// OperLogDetailReq 操作日志详情请求
type OperLogDetailReq struct {
	gmeta.Meta `method:"get" path:"/api/oper/log/detail/:id" summary:"操作日志详情" permission:"sys:operLog:detail" log:"操作日志-查询详情" logType:"query"`
	Id         int64 `json:"id" v:"required#缺少ID" dc:"日志ID"` // 日志ID
}

// OperLogDetailRes 操作日志详情结果
type OperLogDetailRes struct {
	OperLogItem
}

// OperLogDeleteReq 删除操作日志请求
type OperLogDeleteReq struct {
	gmeta.Meta `method:"delete" path:"/api/oper/log/delete/:id" summary:"删除操作日志" permission:"sys:operLog:delete" log:"操作日志-删除记录" logType:"delete"`
	Id         int64 `json:"id" v:"required#缺少ID" dc:"日志ID"` // 日志ID
}

// OperLogDeleteRes 删除操作日志结果
type OperLogDeleteRes struct{}

// OperLogBatchDeleteReq 批量删除操作日志请求
type OperLogBatchDeleteReq struct {
	gmeta.Meta `method:"delete" path:"/api/oper/log/batchDelete" summary:"批量删除操作日志" permission:"sys:operLog:batchDelete" log:"操作日志-批量删除记录" logType:"delete"`
	Ids        []int64 `json:"ids" v:"required#请选择要删除的记录" dc:"日志ID列表"` // 日志ID列表
}

// OperLogBatchDeleteRes 批量删除结果
type OperLogBatchDeleteRes struct{}
