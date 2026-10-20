// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// Article is the golang structure for table article.
type Article struct {
	Id         int         `json:"id"         orm:"id"          description:"主键ID"`            // 主键ID
	TenantId   int         `json:"tenantId"   orm:"tenant_id"   description:"租户ID：0=平台级"`      // 租户ID：0=平台级
	CategoryId int         `json:"categoryId" orm:"category_id" description:"分类ID"`            // 分类ID
	Title      string      `json:"title"      orm:"title"       description:"文章标题"`            // 文章标题
	Cover      string      `json:"cover"      orm:"cover"       description:"文章封面"`            // 文章封面
	Images     string      `json:"images"     orm:"images"      description:"文章图集"`            // 文章图集
	Intro      string      `json:"intro"      orm:"intro"       description:"文章导读"`            // 文章导读
	Content    string      `json:"content"    orm:"content"     description:"文章内容"`            // 文章内容
	Author     string      `json:"author"     orm:"author"      description:"文章作者"`            // 文章作者
	Status     int         `json:"status"     orm:"status"      description:"文章状态：0-正常 1-下架"`  // 文章状态：0-正常 1-下架
	Click      int         `json:"click"      orm:"click"       description:"文章点击率"`           // 文章点击率
	Sort       int         `json:"sort"       orm:"sort"        description:"文章排序"`            // 文章排序
	CreateUser string      `json:"createUser" orm:"create_user" description:"创建人"`             // 创建人
	CreateTime *gtime.Time `json:"createTime" orm:"create_time" description:"创建时间"`            // 创建时间
	UpdateUser string      `json:"updateUser" orm:"update_user" description:"更新人"`             // 更新人
	UpdateTime *gtime.Time `json:"updateTime" orm:"update_time" description:"更新时间"`            // 更新时间
	DelFlag    int         `json:"delFlag"    orm:"del_flag"    description:"删除标识：0-正常 1-已删除"` // 删除标识：0-正常 1-已删除
}
