// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// LoginLogDao is the data access object for the table goframe_login_log.
type LoginLogDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  LoginLogColumns    // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// LoginLogColumns defines and stores column names for the table goframe_login_log.
type LoginLogColumns struct {
	Id            string // 主键ID
	Username      string // 登录账号
	LogType       string // 日志类型：1-登录，2-退出
	Source        string // 登录来源：0-后台，1-前台
	Method        string // 请求方法
	RequestMethod string // 请求方式
	Url           string // 请求URL
	Param         string // 请求参数
	Result        string // 返回结果
	Ip            string // IP地址
	Location      string // IP归属地
	Os            string // 操作系统
	Browser       string // 浏览器
	ConsumeTime   string // 请求耗时
	Status        string // 操作状态：0-正常，1-异常
	Error         string // 错误信息
	CreateUser    string // 创建人
	CreateTime    string // 创建时间
	UpdateUser    string // 更新人
	UpdateTime    string // 更新时间
	DelFlag       string // 删除标识：0-正常 1-已删除
}

// loginLogColumns holds the columns for the table goframe_login_log.
var loginLogColumns = LoginLogColumns{
	Id:            "id",
	Username:      "username",
	LogType:       "log_type",
	Source:        "source",
	Method:        "method",
	RequestMethod: "request_method",
	Url:           "url",
	Param:         "param",
	Result:        "result",
	Ip:            "ip",
	Location:      "location",
	Os:            "os",
	Browser:       "browser",
	ConsumeTime:   "consume_time",
	Status:        "status",
	Error:         "error",
	CreateUser:    "create_user",
	CreateTime:    "create_time",
	UpdateUser:    "update_user",
	UpdateTime:    "update_time",
	DelFlag:       "del_flag",
}

// NewLoginLogDao creates and returns a new DAO object for table data access.
func NewLoginLogDao(handlers ...gdb.ModelHandler) *LoginLogDao {
	return &LoginLogDao{
		group:    "default",
		table:    "goframe_login_log",
		columns:  loginLogColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *LoginLogDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *LoginLogDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *LoginLogDao) Columns() LoginLogColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *LoginLogDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *LoginLogDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *LoginLogDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
