// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// Notice is the golang structure of table goframe_notice for DAO operations like Where/Data.
type Notice struct {
	g.Meta     `orm:"table:goframe_notice, do:true"`
	Id         any         // 主键ID
	TenantId   any         // 租户ID：0=平台级
	Title      any         // 通知标题
	Cover      any         // 通知封面
	Type       any         // 通知类型：1-通知 2-公告
	ClickNum   any         // 点击率
	Status     any         // 通知状态：1-正常 2-停用
	Content    any         // 通知公告内容
	CreateUser any         // 创建人
	CreateTime *gtime.Time // 创建时间
	UpdateUser any         // 更新人
	UpdateTime *gtime.Time // 更新时间
	DelFlag    any         // 删除标识：0-正常 1-已删除
}
