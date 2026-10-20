// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// Menu is the golang structure of table goframe_menu for DAO operations like Where/Data.
type Menu struct {
	g.Meta     `orm:"table:goframe_menu, do:true"`
	Id         any         // 主键ID
	Name       any         // 菜单名称
	Icon       any         // 菜单图标
	Icon2      any         // 菜单图标
	ParentId   any         // 上级ID
	Path       any         // 菜单路径
	Component  any         // 菜单组件
	Permission any         // 权限节点
	Target     any         // 打开方式：0-组件 1-内链 2-外链
	Type       any         // 菜单类型：0-菜单 1-节点
	Status     any         // 菜单状态：0-正常 1-停用
	Hide       any         // 是否可见：0-可见 1-不可见
	Sort       any         // 排序
	Note       any         // 备注
	CreateUser any         // 创建人
	CreateTime *gtime.Time // 创建时间
	UpdateUser any         // 更新人
	UpdateTime *gtime.Time // 更新时间
	DelFlag    any         // 删除标识：0-正常 1-已删除
}
