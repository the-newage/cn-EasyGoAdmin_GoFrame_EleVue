// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// LoginLog is the golang structure of table goframe_login_log for DAO operations like Where/Data.
type LoginLog struct {
	g.Meta        `orm:"table:goframe_login_log, do:true"`
	Id            any         // 主键ID
	Username      any         // 登录账号
	LogType       any         // 日志类型：1-登录，2-退出
	Source        any         // 登录来源：0-后台，1-前台
	Method        any         // 请求方法
	RequestMethod any         // 请求方式
	Url           any         // 请求URL
	Param         any         // 请求参数
	Result        any         // 返回结果
	Ip            any         // IP地址
	Location      any         // IP归属地
	Os            any         // 操作系统
	Browser       any         // 浏览器
	ConsumeTime   any         // 请求耗时
	Status        any         // 操作状态：0-正常，1-异常
	Error         any         // 错误信息
	CreateUser    any         // 创建人
	CreateTime    *gtime.Time // 创建时间
	UpdateUser    any         // 更新人
	UpdateTime    *gtime.Time // 更新时间
	DelFlag       any         // 删除标识：0-正常 1-已删除
}
