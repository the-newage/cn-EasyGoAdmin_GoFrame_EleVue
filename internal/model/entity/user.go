// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// User is the golang structure for table user.
type User struct {
	Id           int         `json:"id"           orm:"id"            description:"主键ID"`            // 主键ID
	TenantId     int         `json:"tenantId"     orm:"tenant_id"     description:"租户ID：0=平台用户"`     // 租户ID：0=平台用户
	Code         string      `json:"code"         orm:"code"          description:"用户编码"`            // 用户编码
	Realname     string      `json:"realname"     orm:"realname"      description:"用户姓名"`            // 用户姓名
	Gender       int         `json:"gender"       orm:"gender"        description:"性别：1-男 2-女 3-保密"` // 性别：1-男 2-女 3-保密
	Avatar       string      `json:"avatar"       orm:"avatar"        description:"用户头像"`            // 用户头像
	Mobile       string      `json:"mobile"       orm:"mobile"        description:"手机号"`             // 手机号
	Email        string      `json:"email"        orm:"email"         description:"邮箱"`              // 邮箱
	DeptId       int         `json:"deptId"       orm:"dept_id"       description:"部门ID"`            // 部门ID
	LevelId      int         `json:"levelId"      orm:"level_id"      description:"职级ID"`            // 职级ID
	PositionId   int         `json:"positionId"   orm:"position_id"   description:"岗位ID"`            // 岗位ID
	ProvinceCode string      `json:"provinceCode" orm:"province_code" description:"省份编码"`            // 省份编码
	CityCode     string      `json:"cityCode"     orm:"city_code"     description:"城市编码"`            // 城市编码
	DistrictCode string      `json:"districtCode" orm:"district_code" description:"县区编码"`            // 县区编码
	StreetCode   string      `json:"streetCode"   orm:"street_code"   description:"街道编码"`            // 街道编码
	CityInfo     string      `json:"cityInfo"     orm:"city_info"     description:"省市区信息"`           // 省市区信息
	Address      string      `json:"address"      orm:"address"       description:"详细地址"`            // 详细地址
	Username     string      `json:"username"     orm:"username"      description:"用户名"`             // 用户名
	Password     string      `json:"password"     orm:"password"      description:"密码"`              // 密码
	Salt         string      `json:"salt"         orm:"salt"          description:"加密盐"`             // 加密盐
	Intro        string      `json:"intro"        orm:"intro"         description:"个人简介"`            // 个人简介
	Status       int         `json:"status"       orm:"status"        description:"状态：1-正常 2-禁用"`    // 状态：1-正常 2-禁用
	Sort         int         `json:"sort"         orm:"sort"          description:"排序"`              // 排序
	Note         string      `json:"note"         orm:"note"          description:"个人备注"`            // 个人备注
	CreateUser   string      `json:"createUser"   orm:"create_user"   description:"创建人"`             // 创建人
	CreateTime   *gtime.Time `json:"createTime"   orm:"create_time"   description:"创建时间"`            // 创建时间
	UpdateUser   string      `json:"updateUser"   orm:"update_user"   description:"更新人"`             // 更新人
	UpdateTime   *gtime.Time `json:"updateTime"   orm:"update_time"   description:"更新时间"`            // 更新时间
	DelFlag      int         `json:"delFlag"      orm:"del_flag"      description:"删除标识：0-正常 1-已删除"` // 删除标识：0-正常 1-已删除
}
