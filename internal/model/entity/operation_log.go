// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// OperationLog is the golang structure for table operation_log.
type OperationLog struct {
	Id            int         `json:"id"            orm:"id"             description:"主键ID"`            // 主键ID
	Title         string      `json:"title"         orm:"title"          description:"请求标题"`            // 请求标题
	Type          int         `json:"type"          orm:"type"           description:"请求类型"`            // 请求类型
	Source        int         `json:"source"        orm:"source"         description:"请求来源：0-后台，1-前台"`  // 请求来源：0-后台，1-前台
	Method        string      `json:"method"        orm:"method"         description:"请求方法"`            // 请求方法
	RequestMethod string      `json:"requestMethod" orm:"request_method" description:"请求方式"`            // 请求方式
	Url           string      `json:"url"           orm:"url"            description:"请求URL"`           // 请求URL
	Param         string      `json:"param"         orm:"param"          description:"请求参数"`            // 请求参数
	Result        string      `json:"result"        orm:"result"         description:"请求结果"`            // 请求结果
	Ip            string      `json:"ip"            orm:"ip"             description:"请求IP"`            // 请求IP
	Location      string      `json:"location"      orm:"location"       description:"请求所在区域"`          // 请求所在区域
	Os            string      `json:"os"            orm:"os"             description:"操作系统"`            // 操作系统
	Browser       string      `json:"browser"       orm:"browser"        description:"浏览器"`             // 浏览器
	ConsumeTime   int64       `json:"consumeTime"   orm:"consume_time"   description:"请求耗时"`            // 请求耗时
	Status        int         `json:"status"        orm:"status"         description:"请求状态：0-正常，1-异常"`  // 请求状态：0-正常，1-异常
	Error         string      `json:"error"         orm:"error"          description:"错误描述"`            // 错误描述
	CreateUser    string      `json:"createUser"    orm:"create_user"    description:"创建人"`             // 创建人
	CreateTime    *gtime.Time `json:"createTime"    orm:"create_time"    description:"创建时间"`            // 创建时间
	UpdateUser    string      `json:"updateUser"    orm:"update_user"    description:"更新人"`             // 更新人
	UpdateTime    *gtime.Time `json:"updateTime"    orm:"update_time"    description:"更新时间"`            // 更新时间
	DelFlag       int         `json:"delFlag"       orm:"del_flag"       description:"删除标识：0-正常 1-已删除"` // 删除标识：0-正常 1-已删除
}
