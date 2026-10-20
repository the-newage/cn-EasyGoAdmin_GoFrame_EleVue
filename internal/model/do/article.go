// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// Article is the golang structure of table goframe_article for DAO operations like Where/Data.
type Article struct {
	g.Meta     `orm:"table:goframe_article, do:true"`
	Id         any         // 主键ID
	TenantId   any         // 租户ID：0=平台级
	CategoryId any         // 分类ID
	Title      any         // 文章标题
	Cover      any         // 文章封面
	Images     any         // 文章图集
	Intro      any         // 文章导读
	Content    any         // 文章内容
	Author     any         // 文章作者
	Status     any         // 文章状态：0-正常 1-下架
	Click      any         // 文章点击率
	Sort       any         // 文章排序
	CreateUser any         // 创建人
	CreateTime *gtime.Time // 创建时间
	UpdateUser any         // 更新人
	UpdateTime *gtime.Time // 更新时间
	DelFlag    any         // 删除标识：0-正常 1-已删除
}
