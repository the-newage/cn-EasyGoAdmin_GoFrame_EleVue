// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// City is the golang structure for table city.
type City struct {
	Id         int         `json:"id"         orm:"id"          description:"主键ID"`                          // 主键ID
	CityCode   string      `json:"cityCode"   orm:"city_code"   description:"城市区号"`                          // 城市区号
	AreaCode   string      `json:"areaCode"   orm:"area_code"   description:"行政编码"`                          // 行政编码
	ParentCode string      `json:"parentCode" orm:"parent_code" description:"上级行政编码"`                        // 上级行政编码
	ZipCode    string      `json:"zipCode"    orm:"zip_code"    description:"邮政编码"`                          // 邮政编码
	Level      int         `json:"level"      orm:"level"       description:"城市级别：0-省份 1-城市 2-县区 3-街道 4-社区"` // 城市级别：0-省份 1-城市 2-县区 3-街道 4-社区
	Pid        int         `json:"pid"        orm:"pid"         description:"上级城市ID"`                        // 上级城市ID
	Name       string      `json:"name"       orm:"name"        description:"城市名称"`                          // 城市名称
	ShortName  string      `json:"shortName"  orm:"short_name"  description:"城市简称"`                          // 城市简称
	FullName   string      `json:"fullName"   orm:"full_name"   description:"城市全称"`                          // 城市全称
	Pinyin     string      `json:"pinyin"     orm:"pinyin"      description:"城市拼音"`                          // 城市拼音
	Lng        string      `json:"lng"        orm:"lng"         description:"城市经度"`                          // 城市经度
	Lat        string      `json:"lat"        orm:"lat"         description:"城市纬度"`                          // 城市纬度
	CreateUser string      `json:"createUser" orm:"create_user" description:"创建人"`                           // 创建人
	CreateTime *gtime.Time `json:"createTime" orm:"create_time" description:"创建时间"`                          // 创建时间
	UpdateUser string      `json:"updateUser" orm:"update_user" description:"更新人"`                           // 更新人
	UpdateTime *gtime.Time `json:"updateTime" orm:"update_time" description:"更新时间"`                          // 更新时间
	DelFlag    int         `json:"delFlag"    orm:"del_flag"    description:"删除标识：0-正常 1-已删除"`               // 删除标识：0-正常 1-已删除
}
