// +======================================================================
// | EasyGoAdmin 企业数字化开发平台 [ 赋能开发者，助力企业发展 ]
// +----------------------------------------------------------------------
// | Copyright (c) 2019-2026 半城风雨团队 版权所有
// +----------------------------------------------------------------------
// | 开源协议: Apache-2.0
// | 重要声明: 未经书面许可，不得移除或修改本框架相关版权标识。
// +----------------------------------------------------------------------
// | 官方网站: https://www.easygoadmin.vip
// +----------------------------------------------------------------------
// | 研发团队: 小蚂蚁云团队 荣誉出品
// +======================================================================
// | 版权和免责声明:
// | 一、知识产权声明
// | 本团队对该软件框架产品（以下简称“本产品”）拥有完整知识产权，包括但不限于商标权、专利权、
// | 著作权、商业秘密等。前述权利均受中华人民共和国及相关国际法律法规保护。未经本团队书面授权，
// | 任何个人或单位不得基于本产品本身申请任何形式的知识产权。
// | 二、使用限制与禁止条款
// | 您理解并同意，禁止将本产品用于以下任何行为：
// | 1. 任何违反中华人民共和国法律法规的用途；
// | 2. 任何侵害他人合法权益（包括但不限于人身权、财产权、知识产权等）的行为；
// | 3. 任何危害国家安全、公共利益及社会道德风尚的行为。
// | 三、责任豁免
// | 因使用本产品而产生的任何损失或法律纠纷（包括但不限于直接损失、间接损失、附带损失、衍生损失
// | 以及因违约、侵权等引发的责任），本团队均不承担任何法律责任。若因违规使用本产品导致任何第三
// | 方提起索赔或行政机关进行调查，您应独立承担全部责任。
// | 四、违约与法律责任
// | 如发现任何违反本声明的行为，本团队有权：
// | 1. 无条件配合公安机关及其他有权机关的调查取证工作；
// | 2. 依法追究相关方的民事责任，并要求赔偿全部损失；
// | 3. 保留进一步采取刑事控告等法律措施的权利。
// | 五、许可范围说明
// | 1、被授权方应视本产品源代码为核心机密信息，负有最高标准的保密与保管义务。禁止以任何形式（包
// | 括但不限于直接或间接的分享、复制、传播、分发、转售、出租、开源、反向工程、反编译、反汇编或通
// | 过网络服务提供给任何第三方）泄露、披露或处置源代码的全部或部分内容。任何未经授权的复制、使用、
// | 传播行为均构成对本团队知识产权的严重侵权，本团队有权依法追究其全部法律责任，并要求赔偿因此造
// | 成的一切经济损失及商誉损害。
// | 2、本产品仅限于公司或个人在法律允许范围内，用于内部且合规的软件产品研发。具体权利义务以《软件
// | 许可协议》附件为准。
// +======================================================================

package v1

import (
	"github.com/gogf/gf/v2/util/gmeta"
)

// ============================== 城市模块 ==============================

// CityItem 城市行数据
// 用于分页列表、下拉列表等场景的城市基础信息
type CityItem struct {
	Id         int64  `json:"id"`         // 主键ID
	CityCode   string `json:"cityCode"`   // 城市编码
	AreaCode   string `json:"areaCode"`   // 区域编码
	ParentCode string `json:"parentCode"` // 父级编码
	ZipCode    string `json:"zipCode"`    // 邮政编码
	Level      int    `json:"level"`      // 层级：1-省 2-市 3-区 4-街道
	LevelText  string `json:"levelText"`  // 层级文本
	Pid        int64  `json:"pid"`        // 父级ID
	Name       string `json:"name"`       // 城市名称
	ShortName  string `json:"shortName"`  // 城市简称
	FullName   string `json:"fullName"`   // 城市全称
	Pinyin     string `json:"pinyin"`     // 拼音
	Lng        string `json:"lng"`        // 经度
	Lat        string `json:"lat"`        // 纬度
}

// CityPageReq 城市分页请求
// 支持关键词、名称、拼音及层级等条件组合查询
type CityPageReq struct {
	gmeta.Meta `method:"get" path:"/api/city/page" summary:"城市分页" permission:"sys:city:page" log:"城市管理-查询分页记录" logType:"query"`
	PageReq           // 分页基础参数
	Keyword    string `json:"keyword" dc:"关键词（名称/拼音模糊）"`      // 关键词，名称或拼音模糊匹配
	Name       string `json:"name" dc:"城市名称模糊"`               // 城市名称模糊匹配
	Pinyin     string `json:"pinyin" dc:"拼音模糊"`               // 拼音模糊匹配
	Level      int    `json:"level" dc:"层级：1-省 2-市 3-区 4-街道"` // 层级
}

// CityPageRes 城市分页结果
type CityPageRes = PageData[CityItem]

// CityListReq 城市列表请求
// 按父级编码、层级等条件查询，供级联选择场景使用
type CityListReq struct {
	gmeta.Meta `method:"get" path:"/api/city/list" summary:"城市列表" permission:"sys:city:list" log:"城市管理-查询列表" logType:"query"`
	Keyword    string `json:"keyword" dc:"关键词"`           // 关键词
	Level      int    `json:"level" dc:"层级"`              // 层级
	ParentCode string `json:"parentCode" dc:"父级编码（顶级为空）"` // 父级编码
	Limit      int    `json:"limit" dc:"条数上限"`            // 条数上限
}

// CityListRes 城市列表结果
type CityListRes = []CityItem

// CityDetailReq 城市详情请求
type CityDetailReq struct {
	gmeta.Meta `method:"get" path:"/api/city/detail/:id" summary:"城市详情" permission:"sys:city:detail" log:"城市管理-查询详情" logType:"query"`
	Id         int64 `json:"id" v:"required#缺少ID" dc:"城市ID"` // 城市ID
}

// CityDetailRes 城市详情结果
type CityDetailRes struct {
	CityItem
	ParentName string `json:"parentName" dc:"父级名称"` // 父级名称
}

// CityAddReq 新增城市请求
type CityAddReq struct {
	gmeta.Meta `method:"post" path:"/api/city/add" summary:"新增城市" permission:"sys:city:add" log:"城市管理-新增记录" logType:"add"`
	CityCode   string `json:"cityCode" v:"required#请输入城市编码" dc:"城市编码"` // 城市编码
	AreaCode   string `json:"areaCode" dc:"区域编码"`                      // 区域编码
	ParentCode string `json:"parentCode" dc:"父级编码"`                    // 父级编码
	ZipCode    string `json:"zipCode" dc:"邮编"`                         // 邮政编码
	Level      int    `json:"level" dc:"层级"`                           // 层级
	Pid        int64  `json:"pid" dc:"父级ID"`                           // 父级ID
	Name       string `json:"name" v:"required#请输入名称" dc:"名称"`         // 城市名称
	ShortName  string `json:"shortName" dc:"简称"`                       // 城市简称
	FullName   string `json:"fullName" dc:"全称"`                        // 城市全称
	Pinyin     string `json:"pinyin" dc:"拼音"`                          // 拼音
	Lng        string `json:"lng" dc:"经度"`                             // 经度
	Lat        string `json:"lat" dc:"纬度"`                             // 纬度
}

// CityUpdateReq 修改城市请求
type CityUpdateReq struct {
	gmeta.Meta `method:"put" path:"/api/city/update" summary:"修改城市" permission:"sys:city:update" log:"城市管理-修改记录" logType:"update"`
	Id         int64  `json:"id" v:"required#缺少ID" dc:"城市ID"` // 城市ID
	CityCode   string `json:"cityCode" dc:"城市编码"`             // 城市编码
	AreaCode   string `json:"areaCode" dc:"区域编码"`             // 区域编码
	ParentCode string `json:"parentCode" dc:"父级编码"`           // 父级编码
	ZipCode    string `json:"zipCode" dc:"邮编"`                // 邮政编码
	Level      int    `json:"level" dc:"层级"`                  // 层级
	Pid        int64  `json:"pid" dc:"父级ID"`                  // 父级ID
	Name       string `json:"name" dc:"名称"`                   // 城市名称
	ShortName  string `json:"shortName" dc:"简称"`              // 城市简称
	FullName   string `json:"fullName" dc:"全称"`               // 城市全称
	Pinyin     string `json:"pinyin" dc:"拼音"`                 // 拼音
	Lng        string `json:"lng" dc:"经度"`                    // 经度
	Lat        string `json:"lat" dc:"纬度"`                    // 纬度
}

// CityUpdateRes 修改城市结果
type CityUpdateRes struct{}

// CityDeleteReq 删除城市请求
type CityDeleteReq struct {
	gmeta.Meta `method:"delete" path:"/api/city/delete/:id" summary:"删除城市" permission:"sys:city:delete" log:"城市管理-删除记录" logType:"delete"`
	Id         int64 `json:"id" v:"required#缺少ID" dc:"城市ID"` // 城市ID
}

// CityDeleteRes 删除城市结果
type CityDeleteRes struct{}

// CityBatchDeleteReq 批量删除城市请求
type CityBatchDeleteReq struct {
	gmeta.Meta `method:"delete" path:"/api/city/batchDelete" summary:"批量删除城市" permission:"sys:city:batchDelete" log:"城市管理-批量删除记录" logType:"delete"`
	Ids        []int64 `json:"ids" v:"required#请选择要删除的记录" dc:"城市ID列表"` // 城市ID列表
}

// CityBatchDeleteRes 批量删除结果
type CityBatchDeleteRes struct{}

// CityByParentReq 按父编码查询下级城市请求
// 登录用户即可访问，用于查询指定父级下的直属下级城市
type CityByParentReq struct {
	gmeta.Meta `method:"get" path:"/api/area/findByParentCode" summary:"按父编码查询城市" log:"城市查询-查询下级城市" logType:"query"`
	ParentCode string `json:"parentCode" v:"required#请输入父级编码" dc:"父级编码"` // 父级编码
}

// CityByParentRes 按父编码查询结果
type CityByParentRes struct {
	List []CityItem `json:"list"` // 城市列表
}

// AreaParentReq 省份列表请求
// 登录用户即可访问，用于查询顶级省份列表
type AreaParentReq struct {
	gmeta.Meta `method:"get" path:"/api/area/getParent" summary:"查询省份列表" log:"城市查询-查询省份列表" logType:"query"`
}

// AreaParentRes 省份列表结果
type AreaParentRes struct {
	List []CityItem `json:"list"` // 省份列表
}

// CityGetListReq 按父编码获取下级城市请求
// 登录用户即可访问，用于级联选择器懒加载下级城市
type CityGetListReq struct {
	gmeta.Meta `method:"get" path:"/api/city/getCityList/:parentCode" summary:"按父编码获取城市列表" log:"城市查询-查询城市列表" logType:"query"`
	ParentCode string `json:"parentCode" v:"required#缺少父级编码" dc:"父级行政编码（顶级为0）"` // 父级行政编码，顶级为0
}

// CityGetListRes 城市列表结果
// data 直接为数组
type CityGetListRes = []CityItem
