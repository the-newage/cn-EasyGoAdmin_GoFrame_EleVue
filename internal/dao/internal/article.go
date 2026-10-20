// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// ArticleDao is the data access object for the table goframe_article.
type ArticleDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  ArticleColumns     // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// ArticleColumns defines and stores column names for the table goframe_article.
type ArticleColumns struct {
	Id         string // 主键ID
	TenantId   string // 租户ID：0=平台级
	CategoryId string // 分类ID
	Title      string // 文章标题
	Cover      string // 文章封面
	Images     string // 文章图集
	Intro      string // 文章导读
	Content    string // 文章内容
	Author     string // 文章作者
	Status     string // 文章状态：0-正常 1-下架
	Click      string // 文章点击率
	Sort       string // 文章排序
	CreateUser string // 创建人
	CreateTime string // 创建时间
	UpdateUser string // 更新人
	UpdateTime string // 更新时间
	DelFlag    string // 删除标识：0-正常 1-已删除
}

// articleColumns holds the columns for the table goframe_article.
var articleColumns = ArticleColumns{
	Id:         "id",
	TenantId:   "tenant_id",
	CategoryId: "category_id",
	Title:      "title",
	Cover:      "cover",
	Images:     "images",
	Intro:      "intro",
	Content:    "content",
	Author:     "author",
	Status:     "status",
	Click:      "click",
	Sort:       "sort",
	CreateUser: "create_user",
	CreateTime: "create_time",
	UpdateUser: "update_user",
	UpdateTime: "update_time",
	DelFlag:    "del_flag",
}

// NewArticleDao creates and returns a new DAO object for table data access.
func NewArticleDao(handlers ...gdb.ModelHandler) *ArticleDao {
	return &ArticleDao{
		group:    "default",
		table:    "goframe_article",
		columns:  articleColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *ArticleDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *ArticleDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *ArticleDao) Columns() ArticleColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *ArticleDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *ArticleDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *ArticleDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
