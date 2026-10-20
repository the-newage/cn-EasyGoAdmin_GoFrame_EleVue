// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// Role is the golang structure for table role.
type Role struct {
	Id         int         `json:"id"         orm:"id"          description:"主键ID"`                           // 主键ID
	TenantId   int         `json:"tenantId"   orm:"tenant_id"   description:"租户ID：0=平台级"`                     // 租户ID：0=平台级
	Name       string      `json:"name"       orm:"name"        description:"角色名称"`                           // 角色名称
	Code       string      `json:"code"       orm:"code"        description:"角色编码"`                           // 角色编码
	DataScope  int         `json:"dataScope"  orm:"data_scope"  description:"数据范围：1-全部 2-本部门及以下 3-本部门 4-仅本人"` // 数据范围：1-全部 2-本部门及以下 3-本部门 4-仅本人
	Sort       int         `json:"sort"       orm:"sort"        description:"角色排序"`                           // 角色排序
	Note       string      `json:"note"       orm:"note"        description:"角色备注"`                           // 角色备注
	CreateUser string      `json:"createUser" orm:"create_user" description:"创建人"`                            // 创建人
	CreateTime *gtime.Time `json:"createTime" orm:"create_time" description:"创建时间"`                           // 创建时间
	UpdateUser string      `json:"updateUser" orm:"update_user" description:"更新人"`                            // 更新人
	UpdateTime *gtime.Time `json:"updateTime" orm:"update_time" description:"更新时间"`                           // 更新时间
	DelFlag    int         `json:"delFlag"    orm:"del_flag"    description:"删除标识：0-正常 1-已删除"`                // 删除标识：0-正常 1-已删除
}
