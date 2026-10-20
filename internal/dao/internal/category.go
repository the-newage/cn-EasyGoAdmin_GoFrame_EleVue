// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// CategoryDao is the data access object for the table goframe_category.
type CategoryDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  CategoryColumns    // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// CategoryColumns defines and stores column names for the table goframe_category.
type CategoryColumns struct {
	Id           string // 主键ID
	TenantId     string // 租户ID：0=平台级
	CategoryName string // 分类名称
	CategoryCode string // 分类编码
	ParentId     string // 上级ID
	Sort         string // 排序
	Note         string // 分类备注
	CreateUser   string // 创建人
	CreateTime   string // 创建时间
	UpdateUser   string // 更新人
	UpdateTime   string // 更新时间
	DelFlag      string // 删除标识：0-正常 1-已删除
}

// categoryColumns holds the columns for the table goframe_category.
var categoryColumns = CategoryColumns{
	Id:           "id",
	TenantId:     "tenant_id",
	CategoryName: "category_name",
	CategoryCode: "category_code",
	ParentId:     "parent_id",
	Sort:         "sort",
	Note:         "note",
	CreateUser:   "create_user",
	CreateTime:   "create_time",
	UpdateUser:   "update_user",
	UpdateTime:   "update_time",
	DelFlag:      "del_flag",
}

// NewCategoryDao creates and returns a new DAO object for table data access.
func NewCategoryDao(handlers ...gdb.ModelHandler) *CategoryDao {
	return &CategoryDao{
		group:    "default",
		table:    "goframe_category",
		columns:  categoryColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *CategoryDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *CategoryDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *CategoryDao) Columns() CategoryColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *CategoryDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *CategoryDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *CategoryDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
