// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// MenuDao is the data access object for the table goframe_menu.
type MenuDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  MenuColumns        // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// MenuColumns defines and stores column names for the table goframe_menu.
type MenuColumns struct {
	Id         string // 主键ID
	Name       string // 菜单名称
	Icon       string // 菜单图标
	Icon2      string // 菜单图标
	ParentId   string // 上级ID
	Path       string // 菜单路径
	Component  string // 菜单组件
	Permission string // 权限节点
	Target     string // 打开方式：0-组件 1-内链 2-外链
	Type       string // 菜单类型：0-菜单 1-节点
	Status     string // 菜单状态：0-正常 1-停用
	Hide       string // 是否可见：0-可见 1-不可见
	Sort       string // 排序
	Note       string // 备注
	CreateUser string // 创建人
	CreateTime string // 创建时间
	UpdateUser string // 更新人
	UpdateTime string // 更新时间
	DelFlag    string // 删除标识：0-正常 1-已删除
}

// menuColumns holds the columns for the table goframe_menu.
var menuColumns = MenuColumns{
	Id:         "id",
	Name:       "name",
	Icon:       "icon",
	Icon2:      "icon2",
	ParentId:   "parent_id",
	Path:       "path",
	Component:  "component",
	Permission: "permission",
	Target:     "target",
	Type:       "type",
	Status:     "status",
	Hide:       "hide",
	Sort:       "sort",
	Note:       "note",
	CreateUser: "create_user",
	CreateTime: "create_time",
	UpdateUser: "update_user",
	UpdateTime: "update_time",
	DelFlag:    "del_flag",
}

// NewMenuDao creates and returns a new DAO object for table data access.
func NewMenuDao(handlers ...gdb.ModelHandler) *MenuDao {
	return &MenuDao{
		group:    "default",
		table:    "goframe_menu",
		columns:  menuColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *MenuDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *MenuDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *MenuDao) Columns() MenuColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *MenuDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *MenuDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *MenuDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
