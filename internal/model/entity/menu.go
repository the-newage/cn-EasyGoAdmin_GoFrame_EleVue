// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// Menu is the golang structure for table menu.
type Menu struct {
	Id         int         `json:"id"         orm:"id"          description:"主键ID"`                // 主键ID
	Name       string      `json:"name"       orm:"name"        description:"菜单名称"`                // 菜单名称
	Icon       string      `json:"icon"       orm:"icon"        description:"菜单图标"`                // 菜单图标
	Icon2      string      `json:"icon2"      orm:"icon2"       description:"菜单图标"`                // 菜单图标
	ParentId   int         `json:"parentId"   orm:"parent_id"   description:"上级ID"`                // 上级ID
	Path       string      `json:"path"       orm:"path"        description:"菜单路径"`                // 菜单路径
	Component  string      `json:"component"  orm:"component"   description:"菜单组件"`                // 菜单组件
	Permission string      `json:"permission" orm:"permission"  description:"权限节点"`                // 权限节点
	Target     int         `json:"target"     orm:"target"      description:"打开方式：0-组件 1-内链 2-外链"` // 打开方式：0-组件 1-内链 2-外链
	Type       int         `json:"type"       orm:"type"        description:"菜单类型：0-菜单 1-节点"`      // 菜单类型：0-菜单 1-节点
	Status     int         `json:"status"     orm:"status"      description:"菜单状态：0-正常 1-停用"`      // 菜单状态：0-正常 1-停用
	Hide       int         `json:"hide"       orm:"hide"        description:"是否可见：0-可见 1-不可见"`     // 是否可见：0-可见 1-不可见
	Sort       int         `json:"sort"       orm:"sort"        description:"排序"`                  // 排序
	Note       string      `json:"note"       orm:"note"        description:"备注"`                  // 备注
	CreateUser string      `json:"createUser" orm:"create_user" description:"创建人"`                 // 创建人
	CreateTime *gtime.Time `json:"createTime" orm:"create_time" description:"创建时间"`                // 创建时间
	UpdateUser string      `json:"updateUser" orm:"update_user" description:"更新人"`                 // 更新人
	UpdateTime *gtime.Time `json:"updateTime" orm:"update_time" description:"更新时间"`                // 更新时间
	DelFlag    int         `json:"delFlag"    orm:"del_flag"    description:"删除标识：0-正常 1-已删除"`     // 删除标识：0-正常 1-已删除
}
