// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// RoleMenu is the golang structure for table role_menu.
type RoleMenu struct {
	Id         int         `json:"id"         orm:"id"          description:"主键ID"`            // 主键ID
	RoleId     int         `json:"roleId"     orm:"role_id"     description:"角色ID"`            // 角色ID
	MenuId     int         `json:"menuId"     orm:"menu_id"     description:"菜单ID"`            // 菜单ID
	CreateUser string      `json:"createUser" orm:"create_user" description:"创建人"`             // 创建人
	CreateTime *gtime.Time `json:"createTime" orm:"create_time" description:"创建时间"`            // 创建时间
	UpdateUser string      `json:"updateUser" orm:"update_user" description:"更新人"`             // 更新人
	UpdateTime *gtime.Time `json:"updateTime" orm:"update_time" description:"更新时间"`            // 更新时间
	DelFlag    int         `json:"delFlag"    orm:"del_flag"    description:"删除标识：0-正常 1-已删除"` // 删除标识：0-正常 1-已删除
}
