// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// City is the golang structure of table goframe_city for DAO operations like Where/Data.
type City struct {
	g.Meta     `orm:"table:goframe_city, do:true"`
	Id         any         // 主键ID
	CityCode   any         // 城市区号
	AreaCode   any         // 行政编码
	ParentCode any         // 上级行政编码
	ZipCode    any         // 邮政编码
	Level      any         // 城市级别：0-省份 1-城市 2-县区 3-街道 4-社区
	Pid        any         // 上级城市ID
	Name       any         // 城市名称
	ShortName  any         // 城市简称
	FullName   any         // 城市全称
	Pinyin     any         // 城市拼音
	Lng        any         // 城市经度
	Lat        any         // 城市纬度
	CreateUser any         // 创建人
	CreateTime *gtime.Time // 创建时间
	UpdateUser any         // 更新人
	UpdateTime *gtime.Time // 更新时间
	DelFlag    any         // 删除标识：0-正常 1-已删除
}
