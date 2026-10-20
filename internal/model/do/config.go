// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// Config is the golang structure of table goframe_config for DAO operations like Where/Data.
type Config struct {
	g.Meta     `orm:"table:goframe_config, do:true"`
	Id         any         // 主键ID
	Name       any         // 配置名称
	Code       any         // 配置编码
	Sort       any         // 配置排序
	Note       any         // 配置备注
	CreateUser any         // 创建人
	CreateTime *gtime.Time // 创建时间
	UpdateUser any         // 更新人
	UpdateTime *gtime.Time // 更新时间
	DelFlag    any         // 删除标识：0-正常 1-已删除
}
