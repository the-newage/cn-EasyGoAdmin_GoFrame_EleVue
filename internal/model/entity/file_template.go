// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// FileTemplate is the golang structure for table file_template.
type FileTemplate struct {
	Id         int         `json:"id"         orm:"id"          description:"主键ID"`            // 主键ID
	Name       string      `json:"name"       orm:"name"        description:"模板名称"`            // 模板名称
	Code       string      `json:"code"       orm:"code"        description:"模板编码"`            // 模板编码
	FileName   string      `json:"fileName"   orm:"file_name"   description:"文件名称"`            // 文件名称
	FilePath   string      `json:"filePath"   orm:"file_path"   description:"文件路径"`            // 文件路径
	Note       string      `json:"note"       orm:"note"        description:"文件模板备注"`          // 文件模板备注
	Sort       int         `json:"sort"       orm:"sort"        description:"文件模板顺序"`          // 文件模板顺序
	CreateUser string      `json:"createUser" orm:"create_user" description:"创建人"`             // 创建人
	CreateTime *gtime.Time `json:"createTime" orm:"create_time" description:"创建时间"`            // 创建时间
	UpdateUser string      `json:"updateUser" orm:"update_user" description:"更新人"`             // 更新人
	UpdateTime *gtime.Time `json:"updateTime" orm:"update_time" description:"更新时间"`            // 更新时间
	DelFlag    int         `json:"delFlag"    orm:"del_flag"    description:"删除标识：0-正常 1-已删除"` // 删除标识：0-正常 1-已删除
}
