// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// Category is the golang structure for table category.
type Category struct {
	Id           int         `json:"id"           orm:"id"            description:"主键ID"`            // 主键ID
	TenantId     int         `json:"tenantId"     orm:"tenant_id"     description:"租户ID：0=平台级"`      // 租户ID：0=平台级
	CategoryName string      `json:"categoryName" orm:"category_name" description:"分类名称"`            // 分类名称
	CategoryCode string      `json:"categoryCode" orm:"category_code" description:"分类编码"`            // 分类编码
	ParentId     int         `json:"parentId"     orm:"parent_id"     description:"上级ID"`            // 上级ID
	Sort         int         `json:"sort"         orm:"sort"          description:"排序"`              // 排序
	Note         string      `json:"note"         orm:"note"          description:"分类备注"`            // 分类备注
	CreateUser   string      `json:"createUser"   orm:"create_user"   description:"创建人"`             // 创建人
	CreateTime   *gtime.Time `json:"createTime"   orm:"create_time"   description:"创建时间"`            // 创建时间
	UpdateUser   string      `json:"updateUser"   orm:"update_user"   description:"更新人"`             // 更新人
	UpdateTime   *gtime.Time `json:"updateTime"   orm:"update_time"   description:"更新时间"`            // 更新时间
	DelFlag      int         `json:"delFlag"      orm:"del_flag"      description:"删除标识：0-正常 1-已删除"` // 删除标识：0-正常 1-已删除
}
