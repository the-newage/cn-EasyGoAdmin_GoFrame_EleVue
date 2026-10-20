// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// OperationLog is the golang structure of table goframe_operation_log for DAO operations like Where/Data.
type OperationLog struct {
	g.Meta        `orm:"table:goframe_operation_log, do:true"`
	Id            any         // 主键ID
	Title         any         // 请求标题
	Type          any         // 请求类型
	Source        any         // 请求来源：0-后台，1-前台
	Method        any         // 请求方法
	RequestMethod any         // 请求方式
	Url           any         // 请求URL
	Param         any         // 请求参数
	Result        any         // 请求结果
	Ip            any         // 请求IP
	Location      any         // 请求所在区域
	Os            any         // 操作系统
	Browser       any         // 浏览器
	ConsumeTime   any         // 请求耗时
	Status        any         // 请求状态：0-正常，1-异常
	Error         any         // 错误描述
	CreateUser    any         // 创建人
	CreateTime    *gtime.Time // 创建时间
	UpdateUser    any         // 更新人
	UpdateTime    *gtime.Time // 更新时间
	DelFlag       any         // 删除标识：0-正常 1-已删除
}
