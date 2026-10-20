// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// Example is the golang structure of table goframe_example for DAO operations like Where/Data.
type Example struct {
	g.Meta     `orm:"table:goframe_example, do:true"`
	Id         any         // 主键ID
	Name       any         // 案例名称
	Avatar     any         // 案例图片
	Type       any         // 案例类型：1-类型1 2-类型2 3-类型3 4-类型4
	Status     any         // 案例状态：1-在用 2-停用
	Sort       any         // 案例排序
	CreateUser any         // 创建人
	CreateTime *gtime.Time // 创建时间
	UpdateUser any         // 更新人
	UpdateTime *gtime.Time // 更新时间
	DelFlag    any         // 删除标识：0-正常 1-已删除
}
