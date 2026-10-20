// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// JobDao is the data access object for the table goframe_job.
type JobDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  JobColumns         // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// JobColumns defines and stores column names for the table goframe_job.
type JobColumns struct {
	Id             string // 主键ID
	JobName        string // 任务名称
	JobAlias       string // 任务别名
	JobGroup       string // 任务分组
	JobTrigger     string // 任务触发器
	Status         string // 任务状态：0-未发布 1-运行中 2-暂停 3-删除
	CronExpression string // Cron表达式
	ExecutePolicy  string // 执行策略：1-立即执行 2-执行一次 3-放弃执行
	IsSync         string // 是否同步任务：0-否 1-是
	Url            string // 任务URL
	Params         string // 任务参数
	Note           string // 任务备注
	CreateUser     string // 创建人
	CreateTime     string // 创建时间
	UpdateUser     string // 更新人
	UpdateTime     string // 更新时间
	DelFlag        string // 删除标识：0-正常 1-已删除
}

// jobColumns holds the columns for the table goframe_job.
var jobColumns = JobColumns{
	Id:             "id",
	JobName:        "job_name",
	JobAlias:       "job_alias",
	JobGroup:       "job_group",
	JobTrigger:     "job_trigger",
	Status:         "status",
	CronExpression: "cron_expression",
	ExecutePolicy:  "execute_policy",
	IsSync:         "is_sync",
	Url:            "url",
	Params:         "params",
	Note:           "note",
	CreateUser:     "create_user",
	CreateTime:     "create_time",
	UpdateUser:     "update_user",
	UpdateTime:     "update_time",
	DelFlag:        "del_flag",
}

// NewJobDao creates and returns a new DAO object for table data access.
func NewJobDao(handlers ...gdb.ModelHandler) *JobDao {
	return &JobDao{
		group:    "default",
		table:    "goframe_job",
		columns:  jobColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *JobDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *JobDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *JobDao) Columns() JobColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *JobDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *JobDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *JobDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
