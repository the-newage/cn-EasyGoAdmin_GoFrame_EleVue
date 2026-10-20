// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// JobLog is the golang structure for table job_log.
type JobLog struct {
	Id             int         `json:"id"             orm:"id"              description:"主键ID"`            // 主键ID
	JobId          int         `json:"jobId"          orm:"job_id"          description:"任务ID"`            // 任务ID
	JobName        string      `json:"jobName"        orm:"job_name"        description:"任务名称"`            // 任务名称
	JobGroup       string      `json:"jobGroup"       orm:"job_group"       description:"任务分组"`            // 任务分组
	JobTrigger     string      `json:"jobTrigger"     orm:"job_trigger"     description:"任务触发器"`           // 任务触发器
	JobMessage     string      `json:"jobMessage"     orm:"job_message"     description:"任务日志信息"`          // 任务日志信息
	CronExpression string      `json:"cronExpression" orm:"cron_expression" description:"Cron表达式"`         // Cron表达式
	Status         int         `json:"status"         orm:"status"          description:"执行状态：0-正常 1-失败"`  // 执行状态：0-正常 1-失败
	StartTime      *gtime.Time `json:"startTime"      orm:"start_time"      description:"执行开始时间"`          // 执行开始时间
	EndTime        *gtime.Time `json:"endTime"        orm:"end_time"        description:"执行结束时间"`          // 执行结束时间
	ConsumeTime    int64       `json:"consumeTime"    orm:"consume_time"    description:"执行耗时，单位：毫秒"`      // 执行耗时，单位：毫秒
	ExceptionInfo  string      `json:"exceptionInfo"  orm:"exception_info"  description:"异常信息"`            // 异常信息
	CreateUser     string      `json:"createUser"     orm:"create_user"     description:"创建人"`             // 创建人
	CreateTime     *gtime.Time `json:"createTime"     orm:"create_time"     description:"创建时间"`            // 创建时间
	UpdateUser     string      `json:"updateUser"     orm:"update_user"     description:"更新人"`             // 更新人
	UpdateTime     *gtime.Time `json:"updateTime"     orm:"update_time"     description:"更新时间"`            // 更新时间
	DelFlag        int         `json:"delFlag"        orm:"del_flag"        description:"删除标识：0-正常 1-已删除"` // 删除标识：0-正常 1-已删除
}
