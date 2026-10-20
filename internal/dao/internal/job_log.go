// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// JobLogDao is the data access object for the table goframe_job_log.
type JobLogDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  JobLogColumns      // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// JobLogColumns defines and stores column names for the table goframe_job_log.
type JobLogColumns struct {
	Id             string // 主键ID
	JobId          string // 任务ID
	JobName        string // 任务名称
	JobGroup       string // 任务分组
	JobTrigger     string // 任务触发器
	JobMessage     string // 任务日志信息
	CronExpression string // Cron表达式
	Status         string // 执行状态：0-正常 1-失败
	StartTime      string // 执行开始时间
	EndTime        string // 执行结束时间
	ConsumeTime    string // 执行耗时，单位：毫秒
	ExceptionInfo  string // 异常信息
	CreateUser     string // 创建人
	CreateTime     string // 创建时间
	UpdateUser     string // 更新人
	UpdateTime     string // 更新时间
	DelFlag        string // 删除标识：0-正常 1-已删除
}

// jobLogColumns holds the columns for the table goframe_job_log.
var jobLogColumns = JobLogColumns{
	Id:             "id",
	JobId:          "job_id",
	JobName:        "job_name",
	JobGroup:       "job_group",
	JobTrigger:     "job_trigger",
	JobMessage:     "job_message",
	CronExpression: "cron_expression",
	Status:         "status",
	StartTime:      "start_time",
	EndTime:        "end_time",
	ConsumeTime:    "consume_time",
	ExceptionInfo:  "exception_info",
	CreateUser:     "create_user",
	CreateTime:     "create_time",
	UpdateUser:     "update_user",
	UpdateTime:     "update_time",
	DelFlag:        "del_flag",
}

// NewJobLogDao creates and returns a new DAO object for table data access.
func NewJobLogDao(handlers ...gdb.ModelHandler) *JobLogDao {
	return &JobLogDao{
		group:    "default",
		table:    "goframe_job_log",
		columns:  jobLogColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *JobLogDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *JobLogDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *JobLogDao) Columns() JobLogColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *JobLogDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *JobLogDao) Ctx(ctx context.Context) *gdb.Model {
	model := dao.DB().Model(dao.table)
	for _, handler := range dao.handlers {
		model = handler(model)
	}
	return model.Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rolls back the transaction and returns the error if function f returns a non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note: Do not commit or roll back the transaction in function f,
// as it is automatically handled by this function.
func (dao *JobLogDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
