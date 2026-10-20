// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// User is the golang structure of table goframe_user for DAO operations like Where/Data.
type User struct {
	g.Meta       `orm:"table:goframe_user, do:true"`
	Id           any         // 主键ID
	TenantId     any         // 租户ID：0=平台用户
	Code         any         // 用户编码
	Realname     any         // 用户姓名
	Gender       any         // 性别：1-男 2-女 3-保密
	Avatar       any         // 用户头像
	Mobile       any         // 手机号
	Email        any         // 邮箱
	DeptId       any         // 部门ID
	LevelId      any         // 职级ID
	PositionId   any         // 岗位ID
	ProvinceCode any         // 省份编码
	CityCode     any         // 城市编码
	DistrictCode any         // 县区编码
	StreetCode   any         // 街道编码
	CityInfo     any         // 省市区信息
	Address      any         // 详细地址
	Username     any         // 用户名
	Password     any         // 密码
	Salt         any         // 加密盐
	Intro        any         // 个人简介
	Status       any         // 状态：1-正常 2-禁用
	Sort         any         // 排序
	Note         any         // 个人备注
	CreateUser   any         // 创建人
	CreateTime   *gtime.Time // 创建时间
	UpdateUser   any         // 更新人
	UpdateTime   *gtime.Time // 更新时间
	DelFlag      any         // 删除标识：0-正常 1-已删除
}
