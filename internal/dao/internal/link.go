// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// LinkDao is the data access object for the table goframe_link.
type LinkDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  LinkColumns        // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// LinkColumns defines and stores column names for the table goframe_link.
type LinkColumns struct {
	Id         string // 主键ID
	TenantId   string // 租户ID：0=平台级
	Name       string // 友链名称
	Type       string // 友链类型：1-友情链接 2-合作伙伴
	Url        string // 友链URL
	Form       string // 友链形式：1-文字链接 2-图片链接
	Image      string // 友链图片
	Status     string // 友链状态：1-正常 2-停用
	Sort       string // 友链顺序
	Note       string // 友链备注
	CreateUser string // 创建人
	CreateTime string // 创建时间
	UpdateUser string // 更新人
	UpdateTime string // 更新时间
	DelFlag    string // 删除标识：0-正常 1-已删除
}

// linkColumns holds the columns for the table goframe_link.
var linkColumns = LinkColumns{
	Id:         "id",
	TenantId:   "tenant_id",
	Name:       "name",
	Type:       "type",
	Url:        "url",
	Form:       "form",
	Image:      "image",
	Status:     "status",
	Sort:       "sort",
	Note:       "note",
	CreateUser: "create_user",
	CreateTime: "create_time",
	UpdateUser: "update_user",
	UpdateTime: "update_time",
	DelFlag:    "del_flag",
}

// NewLinkDao creates and returns a new DAO object for table data access.
func NewLinkDao(handlers ...gdb.ModelHandler) *LinkDao {
	return &LinkDao{
		group:    "default",
		table:    "goframe_link",
		columns:  linkColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *LinkDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *LinkDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *LinkDao) Columns() LinkColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *LinkDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *LinkDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *LinkDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
