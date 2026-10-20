// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// ConfigItem is the golang structure for table config_item.
type ConfigItem struct {
	Id         int         `json:"id"         orm:"id"          description:"主键ID"`            // 主键ID
	Name       string      `json:"name"       orm:"name"        description:"配置项名称"`           // 配置项名称
	Code       string      `json:"code"       orm:"code"        description:"配置项编码"`           // 配置项编码
	Value      string      `json:"value"      orm:"value"       description:"配置项值"`            // 配置项值
	Options    string      `json:"options"    orm:"options"     description:"配置选项"`            // 配置选项
	ConfigId   int         `json:"configId"   orm:"config_id"   description:"配置ID"`            // 配置ID
	Type       string      `json:"type"       orm:"type"        description:"配置类型"`            // 配置类型
	Status     int         `json:"status"     orm:"status"      description:"配置状态：1-正常 2-停用"`  // 配置状态：1-正常 2-停用
	Sort       int         `json:"sort"       orm:"sort"        description:"配置项顺序"`           // 配置项顺序
	Note       string      `json:"note"       orm:"note"        description:"配置项备注"`           // 配置项备注
	CreateUser string      `json:"createUser" orm:"create_user" description:"创建人"`             // 创建人
	CreateTime *gtime.Time `json:"createTime" orm:"create_time" description:"创建时间"`            // 创建时间
	UpdateUser string      `json:"updateUser" orm:"update_user" description:"更新人"`             // 更新人
	UpdateTime *gtime.Time `json:"updateTime" orm:"update_time" description:"更新时间"`            // 更新时间
	DelFlag    int         `json:"delFlag"    orm:"del_flag"    description:"删除标识：0-正常 1-已删除"` // 删除标识：0-正常 1-已删除
}
