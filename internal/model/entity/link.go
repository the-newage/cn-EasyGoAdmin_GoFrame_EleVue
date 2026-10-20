// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// Link is the golang structure for table link.
type Link struct {
	Id         int         `json:"id"         orm:"id"          description:"主键ID"`               // 主键ID
	TenantId   int         `json:"tenantId"   orm:"tenant_id"   description:"租户ID：0=平台级"`         // 租户ID：0=平台级
	Name       string      `json:"name"       orm:"name"        description:"友链名称"`               // 友链名称
	Type       int         `json:"type"       orm:"type"        description:"友链类型：1-友情链接 2-合作伙伴"` // 友链类型：1-友情链接 2-合作伙伴
	Url        string      `json:"url"        orm:"url"         description:"友链URL"`              // 友链URL
	Form       int         `json:"form"       orm:"form"        description:"友链形式：1-文字链接 2-图片链接"` // 友链形式：1-文字链接 2-图片链接
	Image      string      `json:"image"      orm:"image"       description:"友链图片"`               // 友链图片
	Status     int         `json:"status"     orm:"status"      description:"友链状态：1-正常 2-停用"`     // 友链状态：1-正常 2-停用
	Sort       int         `json:"sort"       orm:"sort"        description:"友链顺序"`               // 友链顺序
	Note       string      `json:"note"       orm:"note"        description:"友链备注"`               // 友链备注
	CreateUser string      `json:"createUser" orm:"create_user" description:"创建人"`                // 创建人
	CreateTime *gtime.Time `json:"createTime" orm:"create_time" description:"创建时间"`               // 创建时间
	UpdateUser string      `json:"updateUser" orm:"update_user" description:"更新人"`                // 更新人
	UpdateTime *gtime.Time `json:"updateTime" orm:"update_time" description:"更新时间"`               // 更新时间
	DelFlag    int         `json:"delFlag"    orm:"del_flag"    description:"删除标识：0-正常 1-已删除"`    // 删除标识：0-正常 1-已删除
}
