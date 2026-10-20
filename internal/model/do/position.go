// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// Position is the golang structure of table goframe_position for DAO operations like Where/Data.
type Position struct {
	g.Meta     `orm:"table:goframe_position, do:true"`
	Id         any         // 主键ID
	TenantId   any         // 租户ID：0=平台级
	Name       any         // 岗位名称
	Status     any         // 岗位状态：1-在用 2-停用
	Sort       any         // 岗位排序
	CreateUser any         // 创建人
	CreateTime *gtime.Time // 创建时间
	UpdateUser any         // 更新人
	UpdateTime *gtime.Time // 更新时间
	DelFlag    any         // 删除标识：0-正常 1-已删除
}
