// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// Category is the golang structure of table goframe_category for DAO operations like Where/Data.
type Category struct {
	g.Meta       `orm:"table:goframe_category, do:true"`
	Id           any         // 主键ID
	TenantId     any         // 租户ID：0=平台级
	CategoryName any         // 分类名称
	CategoryCode any         // 分类编码
	ParentId     any         // 上级ID
	Sort         any         // 排序
	Note         any         // 分类备注
	CreateUser   any         // 创建人
	CreateTime   *gtime.Time // 创建时间
	UpdateUser   any         // 更新人
	UpdateTime   *gtime.Time // 更新时间
	DelFlag      any         // 删除标识：0-正常 1-已删除
}
