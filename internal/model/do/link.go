// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// Link is the golang structure of table goframe_link for DAO operations like Where/Data.
type Link struct {
	g.Meta     `orm:"table:goframe_link, do:true"`
	Id         any         // 主键ID
	TenantId   any         // 租户ID：0=平台级
	Name       any         // 友链名称
	Type       any         // 友链类型：1-友情链接 2-合作伙伴
	Url        any         // 友链URL
	Form       any         // 友链形式：1-文字链接 2-图片链接
	Image      any         // 友链图片
	Status     any         // 友链状态：1-正常 2-停用
	Sort       any         // 友链顺序
	Note       any         // 友链备注
	CreateUser any         // 创建人
	CreateTime *gtime.Time // 创建时间
	UpdateUser any         // 更新人
	UpdateTime *gtime.Time // 更新时间
	DelFlag    any         // 删除标识：0-正常 1-已删除
}
