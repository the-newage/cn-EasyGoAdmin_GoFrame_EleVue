// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// Role is the golang structure of table goframe_role for DAO operations like Where/Data.
type Role struct {
	g.Meta     `orm:"table:goframe_role, do:true"`
	Id         any         // 主键ID
	TenantId   any         // 租户ID：0=平台级
	Name       any         // 角色名称
	Code       any         // 角色编码
	DataScope  any         // 数据范围：1-全部 2-本部门及以下 3-本部门 4-仅本人
	Sort       any         // 角色排序
	Note       any         // 角色备注
	CreateUser any         // 创建人
	CreateTime *gtime.Time // 创建时间
	UpdateUser any         // 更新人
	UpdateTime *gtime.Time // 更新时间
	DelFlag    any         // 删除标识：0-正常 1-已删除
}
