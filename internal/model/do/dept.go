// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// Dept is the golang structure of table goframe_dept for DAO operations like Where/Data.
type Dept struct {
	g.Meta     `orm:"table:goframe_dept, do:true"`
	Id         any         // 主键ID
	TenantId   any         // 租户ID：0=平台级
	Name       any         // 部门名称
	Type       any         // 部门类型：1-公司 2-子公司 3-部门 4-小组
	ParentId   any         // 上级部门ID
	Sort       any         // 部门排序
	Note       any         // 备注
	CreateUser any         // 创建人
	CreateTime *gtime.Time // 创建时间
	UpdateUser any         // 更新人
	UpdateTime *gtime.Time // 更新时间
	DelFlag    any         // 删除标识：0-正常 1-已删除
}
