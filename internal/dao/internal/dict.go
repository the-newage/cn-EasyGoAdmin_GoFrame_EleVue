// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// DictDao is the data access object for the table goframe_dict.
type DictDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  DictColumns        // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// DictColumns defines and stores column names for the table goframe_dict.
type DictColumns struct {
	Id         string // 主键ID
	Name       string // 字典名称
	Code       string // 字典编码
	Sort       string // 字典排序
	Note       string // 字典备注
	CreateUser string // 创建人
	CreateTime string // 创建时间
	UpdateUser string // 更新人
	UpdateTime string // 更新时间
	DelFlag    string // 删除标识：0-正常 1-已删除
}

// dictColumns holds the columns for the table goframe_dict.
var dictColumns = DictColumns{
	Id:         "id",
	Name:       "name",
	Code:       "code",
	Sort:       "sort",
	Note:       "note",
	CreateUser: "create_user",
	CreateTime: "create_time",
	UpdateUser: "update_user",
	UpdateTime: "update_time",
	DelFlag:    "del_flag",
}

// NewDictDao creates and returns a new DAO object for table data access.
func NewDictDao(handlers ...gdb.ModelHandler) *DictDao {
	return &DictDao{
		group:    "default",
		table:    "goframe_dict",
		columns:  dictColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *DictDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *DictDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *DictDao) Columns() DictColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *DictDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *DictDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *DictDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
