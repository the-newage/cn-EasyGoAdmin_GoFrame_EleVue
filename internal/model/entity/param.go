// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// Param is the golang structure for table param.
type Param struct {
	Id         int         `json:"id"         orm:"id"          description:"主键ID"`            // 主键ID
	Name       string      `json:"name"       orm:"name"        description:"参数名称"`            // 参数名称
	Code       string      `json:"code"       orm:"code"        description:"参数编码"`            // 参数编码
	Value      string      `json:"value"      orm:"value"       description:"参数值"`             // 参数值
	Type       int         `json:"type"       orm:"type"        description:"参数类型：0-系统 1-业务"`  // 参数类型：0-系统 1-业务
	Status     int         `json:"status"     orm:"status"      description:"参数状态：1-正常 2-禁用"`  // 参数状态：1-正常 2-禁用
	Sort       int         `json:"sort"       orm:"sort"        description:"参数排序"`            // 参数排序
	Note       string      `json:"note"       orm:"note"        description:"参数备注"`            // 参数备注
	CreateUser string      `json:"createUser" orm:"create_user" description:"创建人"`             // 创建人
	CreateTime *gtime.Time `json:"createTime" orm:"create_time" description:"创建时间"`            // 创建时间
	UpdateUser string      `json:"updateUser" orm:"update_user" description:"更新人"`             // 更新人
	UpdateTime *gtime.Time `json:"updateTime" orm:"update_time" description:"更新时间"`            // 更新时间
	DelFlag    int         `json:"delFlag"    orm:"del_flag"    description:"删除标识：0-正常 1-已删除"` // 删除标识：0-正常 1-已删除
}
