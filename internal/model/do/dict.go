// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// Dict is the golang structure of table goframe_dict for DAO operations like Where/Data.
type Dict struct {
	g.Meta     `orm:"table:goframe_dict, do:true"`
	Id         any         // 主键ID
	Name       any         // 字典名称
	Code       any         // 字典编码
	Sort       any         // 字典排序
	Note       any         // 字典备注
	CreateUser any         // 创建人
	CreateTime *gtime.Time // 创建时间
	UpdateUser any         // 更新人
	UpdateTime *gtime.Time // 更新时间
	DelFlag    any         // 删除标识：0-正常 1-已删除
}
