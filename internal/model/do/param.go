// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// Param is the golang structure of table goframe_param for DAO operations like Where/Data.
type Param struct {
	g.Meta     `orm:"table:goframe_param, do:true"`
	Id         any         // 主键ID
	Name       any         // 参数名称
	Code       any         // 参数编码
	Value      any         // 参数值
	Type       any         // 参数类型：0-系统 1-业务
	Status     any         // 参数状态：1-正常 2-禁用
	Sort       any         // 参数排序
	Note       any         // 参数备注
	CreateUser any         // 创建人
	CreateTime *gtime.Time // 创建时间
	UpdateUser any         // 更新人
	UpdateTime *gtime.Time // 更新时间
	DelFlag    any         // 删除标识：0-正常 1-已删除
}
