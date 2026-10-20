// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// ExampleDao is the data access object for the table goframe_example.
type ExampleDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  ExampleColumns     // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// ExampleColumns defines and stores column names for the table goframe_example.
type ExampleColumns struct {
	Id         string // 主键ID
	Name       string // 案例名称
	Avatar     string // 案例图片
	Type       string // 案例类型：1-类型1 2-类型2 3-类型3 4-类型4
	Status     string // 案例状态：1-在用 2-停用
	Sort       string // 案例排序
	CreateUser string // 创建人
	CreateTime string // 创建时间
	UpdateUser string // 更新人
	UpdateTime string // 更新时间
	DelFlag    string // 删除标识：0-正常 1-已删除
}

// exampleColumns holds the columns for the table goframe_example.
var exampleColumns = ExampleColumns{
	Id:         "id",
	Name:       "name",
	Avatar:     "avatar",
	Type:       "type",
	Status:     "status",
	Sort:       "sort",
	CreateUser: "create_user",
	CreateTime: "create_time",
	UpdateUser: "update_user",
	UpdateTime: "update_time",
	DelFlag:    "del_flag",
}

// NewExampleDao creates and returns a new DAO object for table data access.
func NewExampleDao(handlers ...gdb.ModelHandler) *ExampleDao {
	return &ExampleDao{
		group:    "default",
		table:    "goframe_example",
		columns:  exampleColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *ExampleDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *ExampleDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *ExampleDao) Columns() ExampleColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *ExampleDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *ExampleDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *ExampleDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
