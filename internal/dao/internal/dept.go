// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// DeptDao is the data access object for the table goframe_dept.
type DeptDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  DeptColumns        // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// DeptColumns defines and stores column names for the table goframe_dept.
type DeptColumns struct {
	Id         string // 主键ID
	TenantId   string // 租户ID：0=平台级
	Name       string // 部门名称
	Type       string // 部门类型：1-公司 2-子公司 3-部门 4-小组
	ParentId   string // 上级部门ID
	Sort       string // 部门排序
	Note       string // 备注
	CreateUser string // 创建人
	CreateTime string // 创建时间
	UpdateUser string // 更新人
	UpdateTime string // 更新时间
	DelFlag    string // 删除标识：0-正常 1-已删除
}

// deptColumns holds the columns for the table goframe_dept.
var deptColumns = DeptColumns{
	Id:         "id",
	TenantId:   "tenant_id",
	Name:       "name",
	Type:       "type",
	ParentId:   "parent_id",
	Sort:       "sort",
	Note:       "note",
	CreateUser: "create_user",
	CreateTime: "create_time",
	UpdateUser: "update_user",
	UpdateTime: "update_time",
	DelFlag:    "del_flag",
}

// NewDeptDao creates and returns a new DAO object for table data access.
func NewDeptDao(handlers ...gdb.ModelHandler) *DeptDao {
	return &DeptDao{
		group:    "default",
		table:    "goframe_dept",
		columns:  deptColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *DeptDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *DeptDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *DeptDao) Columns() DeptColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *DeptDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *DeptDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *DeptDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
