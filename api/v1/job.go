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

// ============================== 定时任务模块 ==============================

// JobItem 任务行数据
// 用于定时任务分页列表、任务列表等场景的任务基础信息
type JobItem struct {
	Id                int64       `json:"id"`                // 主键ID
	JobName           string      `json:"jobName"`           // 任务名称
	JobAlias          string      `json:"jobAlias"`          // 任务别名
	JobGroup          string      `json:"jobGroup"`          // 任务分组
	JobTrigger        string      `json:"jobTrigger"`        // 触发器类型
	Status            int         `json:"status"`            // 状态：0-未发布 1-运行中 2-暂停
	StatusText        string      `json:"statusText"`        // 状态文本
	CronExpression    string      `json:"cronExpression"`    // Cron表达式
	ExecutePolicy     int         `json:"executePolicy"`     // 执行策略：1-立即执行 2-执行一次 3-放弃执行
	ExecutePolicyText string      `json:"executePolicyText"` // 执行策略文本
	IsSync            int         `json:"isSync"`            // 是否同步：0-否 1-是
	IsSyncText        string      `json:"isSyncText"`        // 是否同步文本
	Url               string      `json:"url"`               // 任务目标：http(s)地址或本地处理方法名
	Params            string      `json:"params"`            // 请求参数（JSON）
	Note              string      `json:"note"`              // 备注
	CreateUser        string      `json:"createUser"`        // 创建人
	CreateTime        *gtime.Time `json:"createTime"`        // 创建时间
}

// JobPageReq 定时任务分页请求
// 支持关键词、任务名称、状态等条件组合查询
type JobPageReq struct {
	gmeta.Meta `method:"get" path:"/api/job/page" summary:"定时任务分页" permission:"sys:job:page" log:"定时任务-查询分页记录" logType:"query"`
	PageReq           // 分页基础参数
	Keyword    string `json:"keyword" dc:"关键词（名称/别名模糊）"`       // 关键词，名称或别名模糊匹配
	JobName    string `json:"jobName" dc:"任务名称模糊"`             // 任务名称模糊匹配
	Status     int    `json:"status" dc:"状态：0-未发布 1-运行中 2-暂停"` // 状态
}

// JobPageRes 定时任务分页结果
type JobPageRes = PageData[JobItem]

// JobListReq 定时任务列表请求
type JobListReq struct {
	gmeta.Meta `method:"get" path:"/api/job/list" summary:"定时任务列表" permission:"sys:job:list" log:"定时任务-查询列表" logType:"query"`
	Keyword    string `json:"keyword" dc:"关键词"` // 关键词
	Status     int    `json:"status" dc:"状态"`   // 状态
}

// JobListRes 定时任务列表结果
type JobListRes = []JobItem

// JobDetailReq 定时任务详情请求
type JobDetailReq struct {
	gmeta.Meta `method:"get" path:"/api/job/detail/:id" summary:"定时任务详情" permission:"sys:job:detail" log:"定时任务-查询详情" logType:"query"`
	Id         int64 `json:"id" v:"required#缺少ID" dc:"任务ID"` // 任务ID
}

// JobDetailRes 定时任务详情结果
type JobDetailRes struct {
	JobItem
}

// JobAddReq 新增定时任务请求
type JobAddReq struct {
	gmeta.Meta     `method:"post" path:"/api/job/add" summary:"新增定时任务" permission:"sys:job:add" log:"定时任务-新增记录" logType:"add"`
	JobName        string `json:"jobName" v:"required#请输入任务名称" dc:"任务名称"`                // 任务名称
	JobAlias       string `json:"jobAlias" dc:"任务别名"`                                    // 任务别名
	JobGroup       string `json:"jobGroup" dc:"任务分组"`                                    // 任务分组
	JobTrigger     string `json:"jobTrigger" dc:"触发器类型"`                                 // 触发器类型
	CronExpression string `json:"cronExpression" v:"required#请输入Cron表达式" dc:"Cron表达式"`   // Cron表达式
	ExecutePolicy  int    `json:"executePolicy" dc:"执行策略：1-立即执行 2-执行一次 3-放弃执行"`          // 执行策略
	IsSync         int    `json:"isSync" dc:"是否同步：0-否 1-是"`                              // 是否同步
	Url            string `json:"url" v:"required#请输入任务目标" dc:"任务目标（http(s)地址或本地处理方法名）"` // 任务目标：http(s)地址或本地处理方法名
	Params         string `json:"params" dc:"请求参数（JSON）"`                                // 请求参数
	Note           string `json:"note" dc:"备注"`                                          // 备注
}

// JobUpdateReq 修改定时任务请求
type JobUpdateReq struct {
	gmeta.Meta     `method:"put" path:"/api/job/update" summary:"修改定时任务" permission:"sys:job:update" log:"定时任务-修改记录" logType:"update"`
	Id             int64  `json:"id" v:"required#缺少ID" dc:"任务ID"`                        // 任务ID
	JobName        string `json:"jobName" v:"required#请输入任务名称" dc:"任务名称"`                // 任务名称
	JobAlias       string `json:"jobAlias" dc:"任务别名"`                                    // 任务别名
	JobGroup       string `json:"jobGroup" dc:"任务分组"`                                    // 任务分组
	JobTrigger     string `json:"jobTrigger" dc:"触发器类型"`                                 // 触发器类型
	CronExpression string `json:"cronExpression" v:"required#请输入Cron表达式" dc:"Cron表达式"`   // Cron表达式
	ExecutePolicy  int    `json:"executePolicy" dc:"执行策略"`                               // 执行策略
	IsSync         int    `json:"isSync" dc:"是否同步"`                                      // 是否同步
	Url            string `json:"url" v:"required#请输入任务目标" dc:"任务目标（http(s)地址或本地处理方法名）"` // 任务目标：http(s)地址或本地处理方法名
	Params         string `json:"params" dc:"请求参数"`                                      // 请求参数
	Note           string `json:"note" dc:"备注"`                                          // 备注
}

// JobUpdateRes 修改定时任务结果
type JobUpdateRes struct{}

// JobDeleteReq 删除定时任务请求
type JobDeleteReq struct {
	gmeta.Meta `method:"delete" path:"/api/job/delete/:id" summary:"删除定时任务" permission:"sys:job:delete" log:"定时任务-删除记录" logType:"delete"`
	Id         int64 `json:"id" v:"required#缺少ID" dc:"任务ID"` // 任务ID
}

// JobDeleteRes 删除定时任务结果
type JobDeleteRes struct{}

// JobBatchDeleteReq 批量删除定时任务请求
type JobBatchDeleteReq struct {
	gmeta.Meta `method:"delete" path:"/api/job/batchDelete" summary:"批量删除定时任务" permission:"sys:job:batchDelete" log:"定时任务-批量删除记录" logType:"delete"`
	Ids        []int64 `json:"ids" v:"required#请选择要删除的记录" dc:"任务ID列表"` // 任务ID列表
}

// JobBatchDeleteRes 批量删除结果
type JobBatchDeleteRes struct{}

// JobRunOnceReq 立即执行一次任务请求
type JobRunOnceReq struct {
	gmeta.Meta `method:"get" path:"/api/job/runOnce/:id" summary:"立即执行任务" permission:"sys:job:execute" log:"定时任务-执行任务" logType:"update"`
	Id         int64 `json:"id" v:"required#缺少ID" dc:"任务ID"` // 任务ID
}

// JobRunOnceRes 立即执行一次任务结果
type JobRunOnceRes struct{}

// JobPauseReq 暂停任务请求
type JobPauseReq struct {
	gmeta.Meta `method:"get" path:"/api/job/pause/:id" summary:"暂停定时任务" permission:"sys:job:pause" log:"定时任务-暂停任务" logType:"update"`
	Id         int64 `json:"id" v:"required#缺少ID" dc:"任务ID"` // 任务ID
}

// JobPauseRes 暂停任务结果
type JobPauseRes struct{}

// JobResumeReq 恢复任务请求
type JobResumeReq struct {
	gmeta.Meta `method:"get" path:"/api/job/resume/:id" summary:"恢复定时任务" permission:"sys:job:resume" log:"定时任务-恢复任务" logType:"update"`
	Id         int64 `json:"id" v:"required#缺少ID" dc:"任务ID"` // 任务ID
}

// JobResumeRes 恢复任务结果
type JobResumeRes struct{}

// JobStatusReq 设置任务状态请求
type JobStatusReq struct {
	gmeta.Meta `method:"post" path:"/api/job/status" summary:"设置定时任务状态" permission:"sys:job:update" log:"定时任务-设置状态" logType:"update"`
	Id         int64 `json:"id" v:"required#缺少ID" dc:"任务ID"`                    // 任务ID
	Status     int   `json:"status" v:"required#缺少状态" dc:"状态：0-未发布 1-运行中 2-暂停"` // 状态
}

// JobStatusRes 设置任务状态结果
type JobStatusRes struct{}

// JobDaemonStatusReq 查询调度器状态请求
type JobDaemonStatusReq struct {
	gmeta.Meta `method:"get" path:"/api/job/daemonStatus" summary:"查询调度器状态" log:"定时任务-查询调度器状态" logType:"query"`
}

// JobDaemonStatusRes 调度器状态结果
type JobDaemonStatusRes struct {
	Running   int      `json:"running"`   // 运行中任务数
	JobIds    []string `json:"jobIds"`    // 调度器中的任务键
	StartTime string   `json:"startTime"` // 调度器启动时间
}
