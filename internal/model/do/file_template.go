// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// FileTemplate is the golang structure of table goframe_file_template for DAO operations like Where/Data.
type FileTemplate struct {
	g.Meta     `orm:"table:goframe_file_template, do:true"`
	Id         any         // 主键ID
	Name       any         // 模板名称
	Code       any         // 模板编码
	FileName   any         // 文件名称
	FilePath   any         // 文件路径
	Note       any         // 文件模板备注
	Sort       any         // 文件模板顺序
	CreateUser any         // 创建人
	CreateTime *gtime.Time // 创建时间
	UpdateUser any         // 更新人
	UpdateTime *gtime.Time // 更新时间
	DelFlag    any         // 删除标识：0-正常 1-已删除
}
