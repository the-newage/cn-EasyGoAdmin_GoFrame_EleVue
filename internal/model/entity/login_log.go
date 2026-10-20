// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// LoginLog is the golang structure for table login_log.
type LoginLog struct {
	Id            int         `json:"id"            orm:"id"             description:"主键ID"`            // 主键ID
	Username      string      `json:"username"      orm:"username"       description:"登录账号"`            // 登录账号
	LogType       int         `json:"logType"       orm:"log_type"       description:"日志类型：1-登录，2-退出"`  // 日志类型：1-登录，2-退出
	Source        int         `json:"source"        orm:"source"         description:"登录来源：0-后台，1-前台"`  // 登录来源：0-后台，1-前台
	Method        string      `json:"method"        orm:"method"         description:"请求方法"`            // 请求方法
	RequestMethod string      `json:"requestMethod" orm:"request_method" description:"请求方式"`            // 请求方式
	Url           string      `json:"url"           orm:"url"            description:"请求URL"`           // 请求URL
	Param         string      `json:"param"         orm:"param"          description:"请求参数"`            // 请求参数
	Result        string      `json:"result"        orm:"result"         description:"返回结果"`            // 返回结果
	Ip            string      `json:"ip"            orm:"ip"             description:"IP地址"`            // IP地址
	Location      string      `json:"location"      orm:"location"       description:"IP归属地"`           // IP归属地
	Os            string      `json:"os"            orm:"os"             description:"操作系统"`            // 操作系统
	Browser       string      `json:"browser"       orm:"browser"        description:"浏览器"`             // 浏览器
	ConsumeTime   int         `json:"consumeTime"   orm:"consume_time"   description:"请求耗时"`            // 请求耗时
	Status        int         `json:"status"        orm:"status"         description:"操作状态：0-正常，1-异常"`  // 操作状态：0-正常，1-异常
	Error         string      `json:"error"         orm:"error"          description:"错误信息"`            // 错误信息
	CreateUser    string      `json:"createUser"    orm:"create_user"    description:"创建人"`             // 创建人
	CreateTime    *gtime.Time `json:"createTime"    orm:"create_time"    description:"创建时间"`            // 创建时间
	UpdateUser    string      `json:"updateUser"    orm:"update_user"    description:"更新人"`             // 更新人
	UpdateTime    *gtime.Time `json:"updateTime"    orm:"update_time"    description:"更新时间"`            // 更新时间
	DelFlag       int         `json:"delFlag"       orm:"del_flag"       description:"删除标识：0-正常 1-已删除"` // 删除标识：0-正常 1-已删除
}
