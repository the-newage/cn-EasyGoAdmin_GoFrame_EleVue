// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// Tenant is the golang structure of table goframe_tenant for DAO operations like Where/Data.
type Tenant struct {
	g.Meta     `orm:"table:goframe_tenant, do:true"`
	Id         any         //
	Name       any         // 租户名称
	Code       any         // 租户编码（系统自动生成：T+4位ID左补零）
	Contact    any         // 联系人
	Mobile     any         // 联系电话
	Email      any         // 邮箱
	Address    any         // 地址
	Logo       any         // 租户LOGO
	Status     any         // 状态：1=正常 2=禁用
	ExpireTime *gtime.Time // 过期时间
	MaxUserNum any         // 最大用户数
	Note       any         // 备注
	CreateUser any         //
	CreateTime *gtime.Time //
	UpdateUser any         //
	UpdateTime *gtime.Time //
	DelFlag    any         //
}
