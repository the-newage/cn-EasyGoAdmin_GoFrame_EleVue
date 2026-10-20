// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// Position is the golang structure for table position.
type Position struct {
	Id         int         `json:"id"         orm:"id"          description:"主键ID"`            // 主键ID
	TenantId   int         `json:"tenantId"   orm:"tenant_id"   description:"租户ID：0=平台级"`      // 租户ID：0=平台级
	Name       string      `json:"name"       orm:"name"        description:"岗位名称"`            // 岗位名称
	Status     int         `json:"status"     orm:"status"      description:"岗位状态：1-在用 2-停用"`  // 岗位状态：1-在用 2-停用
	Sort       int         `json:"sort"       orm:"sort"        description:"岗位排序"`            // 岗位排序
	CreateUser string      `json:"createUser" orm:"create_user" description:"创建人"`             // 创建人
	CreateTime *gtime.Time `json:"createTime" orm:"create_time" description:"创建时间"`            // 创建时间
	UpdateUser string      `json:"updateUser" orm:"update_user" description:"更新人"`             // 更新人
	UpdateTime *gtime.Time `json:"updateTime" orm:"update_time" description:"更新时间"`            // 更新时间
	DelFlag    int         `json:"delFlag"    orm:"del_flag"    description:"删除标识：0-正常 1-已删除"` // 删除标识：0-正常 1-已删除
}
