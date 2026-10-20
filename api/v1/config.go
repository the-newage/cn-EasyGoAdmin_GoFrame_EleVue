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
	"encoding/json"
	"errors"

	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/util/gmeta"
)

// ============================== 配置模块 ==============================

// ConfigItemRow 配置项行数据
// 内嵌使用，供配置详情、配置列表展示配置项明细
type ConfigItemRow struct {
	Id     int64  `json:"id"`     // 主键ID
	Name   string `json:"name"`   // 配置项名称
	Code   string `json:"code"`   // 配置项编码
	Value  string `json:"value"`  // 配置项值
	Type   string `json:"type"`   // 控件类型
	Status int    `json:"status"` // 状态
	Sort   int    `json:"sort"`   // 排序
}

// ConfigItem 配置行数据
// 用于分页列表等场景的配置基础信息
type ConfigItem struct {
	Id         int64  `json:"id"`         // 主键ID
	Name       string `json:"name"`       // 配置名称
	Code       string `json:"code"`       // 配置编码
	Sort       int    `json:"sort"`       // 排序
	Note       string `json:"note"`       // 备注
	CreateUser string `json:"createUser"` // 创建人
}

// ConfigPageReq 配置分页请求
// 支持关键词、配置名称等条件组合查询
type ConfigPageReq struct {
	gmeta.Meta `method:"get" path:"/api/config/page" summary:"配置分页" permission:"sys:config:page" log:"配置管理-查询分页记录" logType:"query"`
	PageReq           // 分页基础参数
	Keyword    string `json:"keyword" dc:"关键词（名称/编码模糊）"` // 关键词，名称或编码模糊匹配
	Name       string `json:"name" dc:"配置名称模糊"`          // 配置名称模糊匹配
}

// ConfigPageRes 配置分页结果
type ConfigPageRes = PageData[ConfigItem]

// ConfigListReq 配置列表请求
// 返回的每一项均携带其下的配置项集合
type ConfigListReq struct {
	gmeta.Meta `method:"get" path:"/api/config/list" summary:"配置列表" permission:"sys:config:list" log:"配置管理-查询列表" logType:"query"`
	Keyword    string `json:"keyword" dc:"关键词"` // 关键词
}

// ConfigListRes 配置列表结果
// 列表每一项携带其下的配置项集合
type ConfigListRes = []ConfigWithItems

// ConfigWithItems 配置及其配置项集合
// 用于列表、详情等需要同时返回配置与配置项的场景
type ConfigWithItems struct {
	ConfigItem
	Items []ConfigItemRow `json:"items"` // 配置项集合
}

// ConfigDetailReq 配置详情请求
type ConfigDetailReq struct {
	gmeta.Meta `method:"get" path:"/api/config/detail/:id" summary:"配置详情" permission:"sys:config:detail" log:"配置管理-查询详情" logType:"query"`
	Id         int64 `json:"id" v:"required#缺少ID" dc:"配置ID"` // 配置ID
}

// ConfigDetailRes 配置详情结果
// 包含配置基础信息及其下的配置项集合
type ConfigDetailRes struct {
	ConfigItem
	Items []ConfigItemRow `json:"items"` // 配置项集合
}

// ConfigAddReq 新增配置请求
type ConfigAddReq struct {
	gmeta.Meta `method:"post" path:"/api/config/add" summary:"新增配置" permission:"sys:config:add" log:"配置管理-新增记录" logType:"add"`
	Name       string `json:"name" v:"required#请输入配置名称" dc:"配置名称"` // 配置名称
	Code       string `json:"code" v:"required#请输入配置编码" dc:"配置编码"` // 配置编码
	Sort       int    `json:"sort" dc:"排序"`                        // 排序
	Note       string `json:"note" dc:"备注"`                        // 备注
}

// ConfigUpdateReq 修改配置请求
type ConfigUpdateReq struct {
	gmeta.Meta `method:"put" path:"/api/config/update" summary:"修改配置" permission:"sys:config:update" log:"配置管理-修改记录" logType:"update"`
	Id         int64  `json:"id" v:"required#缺少ID" dc:"配置ID"`      // 配置ID
	Name       string `json:"name" v:"required#请输入配置名称" dc:"配置名称"` // 配置名称
	Code       string `json:"code" v:"required#请输入配置编码" dc:"配置编码"` // 配置编码
	Sort       int    `json:"sort" dc:"排序"`                        // 排序
	Note       string `json:"note" dc:"备注"`                        // 备注
}

// ConfigUpdateRes 修改配置结果
type ConfigUpdateRes struct{}

// ConfigDeleteReq 删除配置请求
type ConfigDeleteReq struct {
	gmeta.Meta `method:"delete" path:"/api/config/delete/:id" summary:"删除配置" permission:"sys:config:delete" log:"配置管理-删除记录" logType:"delete"`
	Id         int64 `json:"id" v:"required#缺少ID" dc:"配置ID"` // 配置ID
}

// ConfigDeleteRes 删除配置结果
type ConfigDeleteRes struct{}

// ConfigBatchDeleteReq 批量删除配置请求
type ConfigBatchDeleteReq struct {
	gmeta.Meta `method:"delete" path:"/api/config/batchDelete" summary:"批量删除配置" permission:"sys:config:batchDelete" log:"配置管理-批量删除记录" logType:"delete"`
	Ids        []int64 `json:"ids" v:"required#请选择要删除的记录" dc:"配置ID列表"` // 配置ID列表
}

// ConfigBatchDeleteRes 批量删除结果
type ConfigBatchDeleteRes struct{}

// ConfigWebIndexReq 获取网站配置请求
// 返回配置分组及其下的配置项，供网站配置页渲染
type ConfigWebIndexReq struct {
	gmeta.Meta `method:"get" path:"/api/config/web/index" summary:"获取网站配置" permission:"sys:configWeb:data" log:"网站配置-获取配置" logType:"query"`
}

// ConfigWebItem 网站配置项
// value 按控件类型拼接访问URL；images/files 类型拆分到 valueList；options 解析为 param 映射
type ConfigWebItem struct {
	Id        int64             `json:"id"`        // 主键ID
	Name      string            `json:"name"`      // 配置项名称
	Code      string            `json:"code"`      // 配置项编码
	Value     string            `json:"value"`     // 配置项值
	ValueList []string          `json:"valueList"` // 值列表，多值类型使用
	Type      string            `json:"type"`      // 控件类型
	Options   string            `json:"options"`   // 可选项原始串
	Param     map[string]string `json:"param"`     // 可选项解析后的映射
}

// ConfigWebGroup 配置分组
// 表示一个配置标签页下的配置项集合
type ConfigWebGroup struct {
	ConfigName string          `json:"configName"` // 配置分组名称
	DataList   []ConfigWebItem `json:"dataList"`   // 分组下的配置项集合
}

// ConfigWebIndexRes 网站配置聚合结果
type ConfigWebIndexRes = []ConfigWebGroup

// ConfigWebSaveReq 保存网站配置请求
// body 顶层为配置项编码到值的平铺映射，采用自定义解析
type ConfigWebSaveReq struct {
	gmeta.Meta `method:"post" path:"/api/config/web/save" summary:"保存网站配置" permission:"sys:configWeb:save" log:"网站配置-保存配置" logType:"update"`
	Data       map[string]any `json:"data" dc:"配置项编码 → 值"` // 配置项编码到值的映射
}

// UnmarshalBody 自定义请求体解析
// body 顶层为配置项编码到值的平铺映射，跳过默认的请求参数解析流程
func (r *ConfigWebSaveReq) UnmarshalBody(req *ghttp.Request) error {
	var data map[string]any
	if err := json.Unmarshal(req.GetBody(), &data); err != nil {
		return err
	}
	if len(data) == 0 {
		return errors.New("请提交配置数据")
	}
	r.Data = data
	return nil
}

// ConfigWebSaveRes 保存网站配置结果
type ConfigWebSaveRes struct{}
