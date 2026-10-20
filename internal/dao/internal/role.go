// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// RoleDao is the data access object for the table goframe_role.
type RoleDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  RoleColumns        // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// RoleColumns defines and stores column names for the table goframe_role.
type RoleColumns struct {
	Id         string // 主键ID
	TenantId   string // 租户ID：0=平台级
	Name       string // 角色名称
	Code       string // 角色编码
	DataScope  string // 数据范围：1-全部 2-本部门及以下 3-本部门 4-仅本人
	Sort       string // 角色排序
	Note       string // 角色备注
	CreateUser string // 创建人
	CreateTime string // 创建时间
	UpdateUser string // 更新人
	UpdateTime string // 更新时间
	DelFlag    string // 删除标识：0-正常 1-已删除
}

// roleColumns holds the columns for the table goframe_role.
var roleColumns = RoleColumns{
	Id:         "id",
	TenantId:   "tenant_id",
	Name:       "name",
	Code:       "code",
	DataScope:  "data_scope",
	Sort:       "sort",
	Note:       "note",
	CreateUser: "create_user",
	CreateTime: "create_time",
	UpdateUser: "update_user",
	UpdateTime: "update_time",
	DelFlag:    "del_flag",
}

// NewRoleDao creates and returns a new DAO object for table data access.
func NewRoleDao(handlers ...gdb.ModelHandler) *RoleDao {
	return &RoleDao{
		group:    "default",
		table:    "goframe_role",
		columns:  roleColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *RoleDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *RoleDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *RoleDao) Columns() RoleColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *RoleDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *RoleDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *RoleDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
