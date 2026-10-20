// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// Config is the golang structure for table config.
type Config struct {
	Id         int         `json:"id"         orm:"id"          description:"主键ID"`            // 主键ID
	Name       string      `json:"name"       orm:"name"        description:"配置名称"`            // 配置名称
	Code       string      `json:"code"       orm:"code"        description:"配置编码"`            // 配置编码
	Sort       int         `json:"sort"       orm:"sort"        description:"配置排序"`            // 配置排序
	Note       string      `json:"note"       orm:"note"        description:"配置备注"`            // 配置备注
	CreateUser string      `json:"createUser" orm:"create_user" description:"创建人"`             // 创建人
	CreateTime *gtime.Time `json:"createTime" orm:"create_time" description:"创建时间"`            // 创建时间
	UpdateUser string      `json:"updateUser" orm:"update_user" description:"更新人"`             // 更新人
	UpdateTime *gtime.Time `json:"updateTime" orm:"update_time" description:"更新时间"`            // 更新时间
	DelFlag    int         `json:"delFlag"    orm:"del_flag"    description:"删除标识：0-正常 1-已删除"` // 删除标识：0-正常 1-已删除
}
