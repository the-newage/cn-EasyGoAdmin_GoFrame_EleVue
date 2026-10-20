// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// Job is the golang structure for table job.
type Job struct {
	Id             int         `json:"id"             orm:"id"              description:"主键ID"`                       // 主键ID
	JobName        string      `json:"jobName"        orm:"job_name"        description:"任务名称"`                       // 任务名称
	JobAlias       string      `json:"jobAlias"       orm:"job_alias"       description:"任务别名"`                       // 任务别名
	JobGroup       string      `json:"jobGroup"       orm:"job_group"       description:"任务分组"`                       // 任务分组
	JobTrigger     string      `json:"jobTrigger"     orm:"job_trigger"     description:"任务触发器"`                      // 任务触发器
	Status         int         `json:"status"         orm:"status"          description:"任务状态：0-未发布 1-运行中 2-暂停 3-删除"` // 任务状态：0-未发布 1-运行中 2-暂停 3-删除
	CronExpression string      `json:"cronExpression" orm:"cron_expression" description:"Cron表达式"`                    // Cron表达式
	ExecutePolicy  int         `json:"executePolicy"  orm:"execute_policy"  description:"执行策略：1-立即执行 2-执行一次 3-放弃执行"`  // 执行策略：1-立即执行 2-执行一次 3-放弃执行
	IsSync         int         `json:"isSync"         orm:"is_sync"         description:"是否同步任务：0-否 1-是"`             // 是否同步任务：0-否 1-是
	Url            string      `json:"url"            orm:"url"             description:"任务URL"`                      // 任务URL
	Params         string      `json:"params"         orm:"params"          description:"任务参数"`                       // 任务参数
	Note           string      `json:"note"           orm:"note"            description:"任务备注"`                       // 任务备注
	CreateUser     string      `json:"createUser"     orm:"create_user"     description:"创建人"`                        // 创建人
	CreateTime     *gtime.Time `json:"createTime"     orm:"create_time"     description:"创建时间"`                       // 创建时间
	UpdateUser     string      `json:"updateUser"     orm:"update_user"     description:"更新人"`                        // 更新人
	UpdateTime     *gtime.Time `json:"updateTime"     orm:"update_time"     description:"更新时间"`                       // 更新时间
	DelFlag        int         `json:"delFlag"        orm:"del_flag"        description:"删除标识：0-正常 1-已删除"`            // 删除标识：0-正常 1-已删除
}
