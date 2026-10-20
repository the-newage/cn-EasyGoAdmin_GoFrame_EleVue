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

// ============================== 字典模块 ==============================

// DictItemRow 字典项行数据
// 内嵌使用，供字典详情展示字典项明细
type DictItemRow struct {
	Id    int64  `json:"id"`    // 主键ID
	Name  string `json:"name"`  // 字典项名称
	Value string `json:"value"` // 字典项值
	Sort  int    `json:"sort"`  // 排序
}

// DictItem 字典行数据
// 用于字典分页列表、字典列表等场景的字典基础信息
type DictItem struct {
	Id         int64  `json:"id"`         // 主键ID
	Name       string `json:"name"`       // 字典名称
	Code       string `json:"code"`       // 字典编码
	Sort       int    `json:"sort"`       // 排序
	Note       string `json:"note"`       // 备注
	CreateUser string `json:"createUser"` // 创建人
}

// DictPageReq 字典分页请求
// 支持关键词、字典名称等条件组合查询
type DictPageReq struct {
	gmeta.Meta `method:"get" path:"/api/dict/page" summary:"字典分页" permission:"sys:dict:page" log:"字典管理-查询分页记录" logType:"query"`
	PageReq           // 分页基础参数
	Keyword    string `json:"keyword" dc:"关键词（名称/编码模糊）"` // 关键词，名称或编码模糊匹配
	Name       string `json:"name" dc:"字典名称模糊"`          // 字典名称模糊匹配
}

// DictPageRes 字典分页结果
type DictPageRes = PageData[DictItem]

// DictDetailReq 字典详情请求
// 返回字典基础信息及其下的字典项集合
type DictDetailReq struct {
	gmeta.Meta `method:"get" path:"/api/dict/detail/:id" summary:"字典详情" permission:"sys:dict:detail" log:"字典管理-查询详情" logType:"query"`
	Id         int64 `json:"id" v:"required#缺少ID" dc:"字典ID"` // 字典ID
}

// DictDetailRes 字典详情结果
// 包含字典基础信息及其下的字典项集合
type DictDetailRes struct {
	DictItem
	Items []DictItemRow `json:"items"` // 字典项集合
}

// DictAddReq 新增字典请求
type DictAddReq struct {
	gmeta.Meta `method:"post" path:"/api/dict/add" summary:"新增字典" permission:"sys:dict:add" log:"字典管理-新增记录" logType:"add"`
	Name       string `json:"name" v:"required#请输入字典名称" dc:"字典名称"` // 字典名称
	Code       string `json:"code" v:"required#请输入字典编码" dc:"字典编码"` // 字典编码
	Sort       int    `json:"sort" dc:"排序"`                        // 排序
	Note       string `json:"note" dc:"备注"`                        // 备注
}

// DictUpdateReq 修改字典请求
type DictUpdateReq struct {
	gmeta.Meta `method:"put" path:"/api/dict/update" summary:"修改字典" permission:"sys:dict:update" log:"字典管理-修改记录" logType:"update"`
	Id         int64  `json:"id" v:"required#缺少ID" dc:"字典ID"`      // 字典ID
	Name       string `json:"name" v:"required#请输入字典名称" dc:"字典名称"` // 字典名称
	Code       string `json:"code" v:"required#请输入字典编码" dc:"字典编码"` // 字典编码
	Sort       int    `json:"sort" dc:"排序"`                        // 排序
	Note       string `json:"note" dc:"备注"`                        // 备注
}

// DictUpdateRes 修改字典结果
type DictUpdateRes struct{}

// DictDeleteReq 删除字典请求
type DictDeleteReq struct {
	gmeta.Meta `method:"delete" path:"/api/dict/delete/:id" summary:"删除字典" permission:"sys:dict:delete" log:"字典管理-删除记录" logType:"delete"`
	Id         int64 `json:"id" v:"required#缺少ID" dc:"字典ID"` // 字典ID
}

// DictDeleteRes 删除字典结果
type DictDeleteRes struct{}

// DictBatchDeleteReq 批量删除字典请求
type DictBatchDeleteReq struct {
	gmeta.Meta `method:"delete" path:"/api/dict/batchDelete" summary:"批量删除字典" permission:"sys:dict:batchDelete" log:"字典管理-批量删除记录" logType:"delete"`
	Ids        []int64 `json:"ids" v:"required#请选择要删除的记录" dc:"字典ID列表"` // 字典ID列表
}

// DictBatchDeleteRes 批量删除结果
type DictBatchDeleteRes struct{}

// DictRefreshCacheReq 刷新字典缓存请求
type DictRefreshCacheReq struct {
	gmeta.Meta `method:"get" path:"/api/dict/refreshCache" summary:"刷新字典缓存" permission:"sys:dict:cache" log:"字典管理-刷新缓存" logType:"update"`
	Code       string `json:"code" dc:"字典编码（留空刷新全部）"` // 字典编码，留空刷新全部
}

// DictRefreshCacheRes 刷新缓存结果
type DictRefreshCacheRes struct{}

// DictAliasListReq 字典别名列表请求
// 登录用户即可访问，供前端字典查询场景使用
type DictAliasListReq struct {
	gmeta.Meta `method:"get" path:"/api/dictionary/list" summary:"字典列表" log:"字典管理-查询列表" logType:"query"`
	Keyword    string `json:"keyword" dc:"关键词"` // 关键词
}

// DictAliasListRes 字典别名列表结果
type DictAliasListRes = []DictItem

// DictAliasInfoReq 字典别名详情请求
// 登录用户即可访问，按字典编码查询字典及其字典项
type DictAliasInfoReq struct {
	gmeta.Meta `method:"get" path:"/api/dictionary/info" summary:"字典详情" log:"字典管理-查询详情" logType:"query"`
	Code       string `json:"code" v:"required#请输入字典编码" dc:"字典编码"` // 字典编码
}

// DictAliasInfoRes 字典别名详情结果
// 包含字典基础信息及其下的字典项集合
type DictAliasInfoRes struct {
	DictItem
	Items []DictItemRow `json:"items"` // 字典项集合
}
