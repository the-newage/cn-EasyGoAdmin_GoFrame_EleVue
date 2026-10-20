// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// Job is the golang structure of table goframe_job for DAO operations like Where/Data.
type Job struct {
	g.Meta         `orm:"table:goframe_job, do:true"`
	Id             any         // 主键ID
	JobName        any         // 任务名称
	JobAlias       any         // 任务别名
	JobGroup       any         // 任务分组
	JobTrigger     any         // 任务触发器
	Status         any         // 任务状态：0-未发布 1-运行中 2-暂停 3-删除
	CronExpression any         // Cron表达式
	ExecutePolicy  any         // 执行策略：1-立即执行 2-执行一次 3-放弃执行
	IsSync         any         // 是否同步任务：0-否 1-是
	Url            any         // 任务URL
	Params         any         // 任务参数
	Note           any         // 任务备注
	CreateUser     any         // 创建人
	CreateTime     *gtime.Time // 创建时间
	UpdateUser     any         // 更新人
	UpdateTime     *gtime.Time // 更新时间
	DelFlag        any         // 删除标识：0-正常 1-已删除
}
