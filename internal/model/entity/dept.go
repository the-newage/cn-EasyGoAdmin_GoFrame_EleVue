// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// Dept is the golang structure for table dept.
type Dept struct {
	Id         int         `json:"id"         orm:"id"          description:"主键ID"`                      // 主键ID
	TenantId   int         `json:"tenantId"   orm:"tenant_id"   description:"租户ID：0=平台级"`                // 租户ID：0=平台级
	Name       string      `json:"name"       orm:"name"        description:"部门名称"`                      // 部门名称
	Type       int         `json:"type"       orm:"type"        description:"部门类型：1-公司 2-子公司 3-部门 4-小组"` // 部门类型：1-公司 2-子公司 3-部门 4-小组
	ParentId   int         `json:"parentId"   orm:"parent_id"   description:"上级部门ID"`                    // 上级部门ID
	Sort       int         `json:"sort"       orm:"sort"        description:"部门排序"`                      // 部门排序
	Note       string      `json:"note"       orm:"note"        description:"备注"`                        // 备注
	CreateUser string      `json:"createUser" orm:"create_user" description:"创建人"`                       // 创建人
	CreateTime *gtime.Time `json:"createTime" orm:"create_time" description:"创建时间"`                      // 创建时间
	UpdateUser string      `json:"updateUser" orm:"update_user" description:"更新人"`                       // 更新人
	UpdateTime *gtime.Time `json:"updateTime" orm:"update_time" description:"更新时间"`                      // 更新时间
	DelFlag    int         `json:"delFlag"    orm:"del_flag"    description:"删除标识：0-正常 1-已删除"`           // 删除标识：0-正常 1-已删除
}
