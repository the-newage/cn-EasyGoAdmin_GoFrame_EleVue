// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// TenantDao is the data access object for the table goframe_tenant.
type TenantDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  TenantColumns      // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// TenantColumns defines and stores column names for the table goframe_tenant.
type TenantColumns struct {
	Id         string //
	Name       string // 租户名称
	Code       string // 租户编码（系统自动生成：T+4位ID左补零）
	Contact    string // 联系人
	Mobile     string // 联系电话
	Email      string // 邮箱
	Address    string // 地址
	Logo       string // 租户LOGO
	Status     string // 状态：1=正常 2=禁用
	ExpireTime string // 过期时间
	MaxUserNum string // 最大用户数
	Note       string // 备注
	CreateUser string //
	CreateTime string //
	UpdateUser string //
	UpdateTime string //
	DelFlag    string //
}

// tenantColumns holds the columns for the table goframe_tenant.
var tenantColumns = TenantColumns{
	Id:         "id",
	Name:       "name",
	Code:       "code",
	Contact:    "contact",
	Mobile:     "mobile",
	Email:      "email",
	Address:    "address",
	Logo:       "logo",
	Status:     "status",
	ExpireTime: "expire_time",
	MaxUserNum: "max_user_num",
	Note:       "note",
	CreateUser: "create_user",
	CreateTime: "create_time",
	UpdateUser: "update_user",
	UpdateTime: "update_time",
	DelFlag:    "del_flag",
}

// NewTenantDao creates and returns a new DAO object for table data access.
func NewTenantDao(handlers ...gdb.ModelHandler) *TenantDao {
	return &TenantDao{
		group:    "default",
		table:    "goframe_tenant",
		columns:  tenantColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *TenantDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *TenantDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *TenantDao) Columns() TenantColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *TenantDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *TenantDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *TenantDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
