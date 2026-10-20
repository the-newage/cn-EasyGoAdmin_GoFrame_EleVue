// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// JobLog is the golang structure of table goframe_job_log for DAO operations like Where/Data.
type JobLog struct {
	g.Meta         `orm:"table:goframe_job_log, do:true"`
	Id             any         // 主键ID
	JobId          any         // 任务ID
	JobName        any         // 任务名称
	JobGroup       any         // 任务分组
	JobTrigger     any         // 任务触发器
	JobMessage     any         // 任务日志信息
	CronExpression any         // Cron表达式
	Status         any         // 执行状态：0-正常 1-失败
	StartTime      *gtime.Time // 执行开始时间
	EndTime        *gtime.Time // 执行结束时间
	ConsumeTime    any         // 执行耗时，单位：毫秒
	ExceptionInfo  any         // 异常信息
	CreateUser     any         // 创建人
	CreateTime     *gtime.Time // 创建时间
	UpdateUser     any         // 更新人
	UpdateTime     *gtime.Time // 更新时间
	DelFlag        any         // 删除标识：0-正常 1-已删除
}
