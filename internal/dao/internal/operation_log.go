// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// OperationLogDao is the data access object for the table goframe_operation_log.
type OperationLogDao struct {
	table    string              // table is the underlying table name of the DAO.
	group    string              // group is the database configuration group name of the current DAO.
	columns  OperationLogColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler  // handlers for customized model modification.
}

// OperationLogColumns defines and stores column names for the table goframe_operation_log.
type OperationLogColumns struct {
	Id            string // 主键ID
	Title         string // 请求标题
	Type          string // 请求类型
	Source        string // 请求来源：0-后台，1-前台
	Method        string // 请求方法
	RequestMethod string // 请求方式
	Url           string // 请求URL
	Param         string // 请求参数
	Result        string // 请求结果
	Ip            string // 请求IP
	Location      string // 请求所在区域
	Os            string // 操作系统
	Browser       string // 浏览器
	ConsumeTime   string // 请求耗时
	Status        string // 请求状态：0-正常，1-异常
	Error         string // 错误描述
	CreateUser    string // 创建人
	CreateTime    string // 创建时间
	UpdateUser    string // 更新人
	UpdateTime    string // 更新时间
	DelFlag       string // 删除标识：0-正常 1-已删除
}

// operationLogColumns holds the columns for the table goframe_operation_log.
var operationLogColumns = OperationLogColumns{
	Id:            "id",
	Title:         "title",
	Type:          "type",
	Source:        "source",
	Method:        "method",
	RequestMethod: "request_method",
	Url:           "url",
	Param:         "param",
	Result:        "result",
	Ip:            "ip",
	Location:      "location",
	Os:            "os",
	Browser:       "browser",
	ConsumeTime:   "consume_time",
	Status:        "status",
	Error:         "error",
	CreateUser:    "create_user",
	CreateTime:    "create_time",
	UpdateUser:    "update_user",
	UpdateTime:    "update_time",
	DelFlag:       "del_flag",
}

// NewOperationLogDao creates and returns a new DAO object for table data access.
func NewOperationLogDao(handlers ...gdb.ModelHandler) *OperationLogDao {
	return &OperationLogDao{
		group:    "default",
		table:    "goframe_operation_log",
		columns:  operationLogColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *OperationLogDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *OperationLogDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *OperationLogDao) Columns() OperationLogColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *OperationLogDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *OperationLogDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *OperationLogDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
