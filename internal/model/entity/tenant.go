// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// Tenant is the golang structure for table tenant.
type Tenant struct {
	Id         int         `json:"id"         orm:"id"           description:""`                       //
	Name       string      `json:"name"       orm:"name"         description:"租户名称"`                   // 租户名称
	Code       string      `json:"code"       orm:"code"         description:"租户编码（系统自动生成：T+4位ID左补零）"` // 租户编码（系统自动生成：T+4位ID左补零）
	Contact    string      `json:"contact"    orm:"contact"      description:"联系人"`                    // 联系人
	Mobile     string      `json:"mobile"     orm:"mobile"       description:"联系电话"`                   // 联系电话
	Email      string      `json:"email"      orm:"email"        description:"邮箱"`                     // 邮箱
	Address    string      `json:"address"    orm:"address"      description:"地址"`                     // 地址
	Logo       string      `json:"logo"       orm:"logo"         description:"租户LOGO"`                 // 租户LOGO
	Status     int         `json:"status"     orm:"status"       description:"状态：1=正常 2=禁用"`           // 状态：1=正常 2=禁用
	ExpireTime *gtime.Time `json:"expireTime" orm:"expire_time"  description:"过期时间"`                   // 过期时间
	MaxUserNum int         `json:"maxUserNum" orm:"max_user_num" description:"最大用户数"`                  // 最大用户数
	Note       string      `json:"note"       orm:"note"         description:"备注"`                     // 备注
	CreateUser string      `json:"createUser" orm:"create_user"  description:""`                       //
	CreateTime *gtime.Time `json:"createTime" orm:"create_time"  description:""`                       //
	UpdateUser string      `json:"updateUser" orm:"update_user"  description:""`                       //
	UpdateTime *gtime.Time `json:"updateTime" orm:"update_time"  description:""`                       //
	DelFlag    int         `json:"delFlag"    orm:"del_flag"     description:""`                       //
}
