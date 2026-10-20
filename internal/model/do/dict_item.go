// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// DictItem is the golang structure of table goframe_dict_item for DAO operations like Where/Data.
type DictItem struct {
	g.Meta     `orm:"table:goframe_dict_item, do:true"`
	Id         any         // 主键ID
	Name       any         // 字典项名称
	Value      any         // 字典项值
	DictId     any         // 字典ID
	Sort       any         // 字典项顺序
	Note       any         // 字典项备注
	CreateUser any         // 创建人
	CreateTime *gtime.Time // 创建时间
	UpdateUser any         // 更新人
	UpdateTime *gtime.Time // 更新时间
	DelFlag    any         // 删除标识：0-正常 1-已删除
}
