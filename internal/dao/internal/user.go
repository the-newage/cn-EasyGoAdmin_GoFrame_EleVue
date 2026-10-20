// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// UserDao is the data access object for the table goframe_user.
type UserDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  UserColumns        // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// UserColumns defines and stores column names for the table goframe_user.
type UserColumns struct {
	Id           string // 主键ID
	TenantId     string // 租户ID：0=平台用户
	Code         string // 用户编码
	Realname     string // 用户姓名
	Gender       string // 性别：1-男 2-女 3-保密
	Avatar       string // 用户头像
	Mobile       string // 手机号
	Email        string // 邮箱
	DeptId       string // 部门ID
	LevelId      string // 职级ID
	PositionId   string // 岗位ID
	ProvinceCode string // 省份编码
	CityCode     string // 城市编码
	DistrictCode string // 县区编码
	StreetCode   string // 街道编码
	CityInfo     string // 省市区信息
	Address      string // 详细地址
	Username     string // 用户名
	Password     string // 密码
	Salt         string // 加密盐
	Intro        string // 个人简介
	Status       string // 状态：1-正常 2-禁用
	Sort         string // 排序
	Note         string // 个人备注
	CreateUser   string // 创建人
	CreateTime   string // 创建时间
	UpdateUser   string // 更新人
	UpdateTime   string // 更新时间
	DelFlag      string // 删除标识：0-正常 1-已删除
}

// userColumns holds the columns for the table goframe_user.
var userColumns = UserColumns{
	Id:           "id",
	TenantId:     "tenant_id",
	Code:         "code",
	Realname:     "realname",
	Gender:       "gender",
	Avatar:       "avatar",
	Mobile:       "mobile",
	Email:        "email",
	DeptId:       "dept_id",
	LevelId:      "level_id",
	PositionId:   "position_id",
	ProvinceCode: "province_code",
	CityCode:     "city_code",
	DistrictCode: "district_code",
	StreetCode:   "street_code",
	CityInfo:     "city_info",
	Address:      "address",
	Username:     "username",
	Password:     "password",
	Salt:         "salt",
	Intro:        "intro",
	Status:       "status",
	Sort:         "sort",
	Note:         "note",
	CreateUser:   "create_user",
	CreateTime:   "create_time",
	UpdateUser:   "update_user",
	UpdateTime:   "update_time",
	DelFlag:      "del_flag",
}

// NewUserDao creates and returns a new DAO object for table data access.
func NewUserDao(handlers ...gdb.ModelHandler) *UserDao {
	return &UserDao{
		group:    "default",
		table:    "goframe_user",
		columns:  userColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *UserDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *UserDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *UserDao) Columns() UserColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *UserDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *UserDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *UserDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
