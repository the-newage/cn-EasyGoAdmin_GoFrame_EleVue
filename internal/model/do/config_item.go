// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// ConfigItem is the golang structure of table goframe_config_item for DAO operations like Where/Data.
type ConfigItem struct {
	g.Meta     `orm:"table:goframe_config_item, do:true"`
	Id         any         // 主键ID
	Name       any         // 配置项名称
	Code       any         // 配置项编码
	Value      any         // 配置项值
	Options    any         // 配置选项
	ConfigId   any         // 配置ID
	Type       any         // 配置类型
	Status     any         // 配置状态：1-正常 2-停用
	Sort       any         // 配置项顺序
	Note       any         // 配置项备注
	CreateUser any         // 创建人
	CreateTime *gtime.Time // 创建时间
	UpdateUser any         // 更新人
	UpdateTime *gtime.Time // 更新时间
	DelFlag    any         // 删除标识：0-正常 1-已删除
}
