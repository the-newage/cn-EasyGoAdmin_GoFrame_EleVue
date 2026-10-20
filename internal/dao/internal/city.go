// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// CityDao is the data access object for the table goframe_city.
type CityDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  CityColumns        // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// CityColumns defines and stores column names for the table goframe_city.
type CityColumns struct {
	Id         string // 主键ID
	CityCode   string // 城市区号
	AreaCode   string // 行政编码
	ParentCode string // 上级行政编码
	ZipCode    string // 邮政编码
	Level      string // 城市级别：0-省份 1-城市 2-县区 3-街道 4-社区
	Pid        string // 上级城市ID
	Name       string // 城市名称
	ShortName  string // 城市简称
	FullName   string // 城市全称
	Pinyin     string // 城市拼音
	Lng        string // 城市经度
	Lat        string // 城市纬度
	CreateUser string // 创建人
	CreateTime string // 创建时间
	UpdateUser string // 更新人
	UpdateTime string // 更新时间
	DelFlag    string // 删除标识：0-正常 1-已删除
}

// cityColumns holds the columns for the table goframe_city.
var cityColumns = CityColumns{
	Id:         "id",
	CityCode:   "city_code",
	AreaCode:   "area_code",
	ParentCode: "parent_code",
	ZipCode:    "zip_code",
	Level:      "level",
	Pid:        "pid",
	Name:       "name",
	ShortName:  "short_name",
	FullName:   "full_name",
	Pinyin:     "pinyin",
	Lng:        "lng",
	Lat:        "lat",
	CreateUser: "create_user",
	CreateTime: "create_time",
	UpdateUser: "update_user",
	UpdateTime: "update_time",
	DelFlag:    "del_flag",
}

// NewCityDao creates and returns a new DAO object for table data access.
func NewCityDao(handlers ...gdb.ModelHandler) *CityDao {
	return &CityDao{
		group:    "default",
		table:    "goframe_city",
		columns:  cityColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *CityDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *CityDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *CityDao) Columns() CityColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *CityDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *CityDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *CityDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
