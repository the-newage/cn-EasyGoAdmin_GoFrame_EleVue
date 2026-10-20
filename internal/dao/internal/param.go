// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// ParamDao is the data access object for the table goframe_param.
type ParamDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  ParamColumns       // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// ParamColumns defines and stores column names for the table goframe_param.
type ParamColumns struct {
	Id         string // 主键ID
	Name       string // 参数名称
	Code       string // 参数编码
	Value      string // 参数值
	Type       string // 参数类型：0-系统 1-业务
	Status     string // 参数状态：1-正常 2-禁用
	Sort       string // 参数排序
	Note       string // 参数备注
	CreateUser string // 创建人
	CreateTime string // 创建时间
	UpdateUser string // 更新人
	UpdateTime string // 更新时间
	DelFlag    string // 删除标识：0-正常 1-已删除
}

// paramColumns holds the columns for the table goframe_param.
var paramColumns = ParamColumns{
	Id:         "id",
	Name:       "name",
	Code:       "code",
	Value:      "value",
	Type:       "type",
	Status:     "status",
	Sort:       "sort",
	Note:       "note",
	CreateUser: "create_user",
	CreateTime: "create_time",
	UpdateUser: "update_user",
	UpdateTime: "update_time",
	DelFlag:    "del_flag",
}

// NewParamDao creates and returns a new DAO object for table data access.
func NewParamDao(handlers ...gdb.ModelHandler) *ParamDao {
	return &ParamDao{
		group:    "default",
		table:    "goframe_param",
		columns:  paramColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *ParamDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *ParamDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *ParamDao) Columns() ParamColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *ParamDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *ParamDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *ParamDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
