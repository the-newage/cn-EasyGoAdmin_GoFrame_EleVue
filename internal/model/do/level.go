// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// Level is the golang structure of table goframe_level for DAO operations like Where/Data.
type Level struct {
	g.Meta     `orm:"table:goframe_level, do:true"`
	Id         any         // 主键ID
	TenantId   any         // 租户ID：0=平台级
	Name       any         // 职级名称
	Status     any         // 职级状态：1-在用 2-停用
	Sort       any         // 职级排序
	CreateUser any         // 创建人
	CreateTime *gtime.Time // 创建时间
	UpdateUser any         // 更新人
	UpdateTime *gtime.Time // 更新时间
	DelFlag    any         // 删除标识：0-正常 1-已删除
}
