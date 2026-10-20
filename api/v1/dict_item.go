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

// ============================== 字典项模块 ==============================

// DictItemItemRow 字典项行数据
// 用于字典项分页列表等场景的字典项基础信息
type DictItemItemRow struct {
	Id         int64  `json:"id"`         // 主键ID
	DictId     int64  `json:"dictId"`     // 所属字典ID
	DictName   string `json:"dictName"`   // 所属字典名称
	DictCode   string `json:"dictCode"`   // 所属字典编码
	Name       string `json:"name"`       // 字典项名称
	Value      string `json:"value"`      // 字典项值
	Sort       int    `json:"sort"`       // 排序
	Note       string `json:"note"`       // 备注
	CreateUser string `json:"createUser"` // 创建人
}

// DictItemPageReq 字典项分页请求
// 支持关键词、字典项名称、所属字典等条件组合查询
type DictItemPageReq struct {
	gmeta.Meta `method:"get" path:"/api/dict/item/page" summary:"字典项分页" permission:"sys:dictitem:page" log:"字典项管理-查询分页记录" logType:"query"`
	PageReq           // 分页基础参数
	Keyword    string `json:"keyword" dc:"关键词（名称/值模糊）"` // 关键词，名称或值模糊匹配
	Name       string `json:"name" dc:"字典项名称模糊"`        // 字典项名称模糊匹配
	DictId     int64  `json:"dictId" dc:"所属字典ID"`       // 所属字典ID
}

// DictItemPageRes 字典项分页结果
type DictItemPageRes = PageData[DictItemItemRow]

// DictItemDetailReq 字典项详情请求
type DictItemDetailReq struct {
	gmeta.Meta `method:"get" path:"/api/dict/item/detail/:id" summary:"字典项详情" permission:"sys:dictitem:detail" log:"字典项管理-查询详情" logType:"query"`
	Id         int64 `json:"id" v:"required#缺少ID" dc:"字典项ID"` // 字典项ID
}

// DictItemDetailRes 字典项详情结果
type DictItemDetailRes struct {
	Id     int64  `json:"id"`     // 主键ID
	DictId int64  `json:"dictId"` // 所属字典ID
	Name   string `json:"name"`   // 字典项名称
	Value  string `json:"value"`  // 字典项值
	Sort   int    `json:"sort"`   // 排序
	Note   string `json:"note"`   // 备注
}

// DictItemAddReq 新增字典项请求
type DictItemAddReq struct {
	gmeta.Meta `method:"post" path:"/api/dict/item/add" summary:"新增字典项" permission:"sys:dictitem:add" log:"字典项管理-新增记录" logType:"add"`
	DictId     int64  `json:"dictId" v:"required#请选择所属字典" dc:"所属字典ID"` // 所属字典ID
	Name       string `json:"name" v:"required#请输入字典项名称" dc:"字典项名称"`   // 字典项名称
	Value      string `json:"value" v:"required#请输入字典项值" dc:"字典项值"`    // 字典项值
	Sort       int    `json:"sort" dc:"排序"`                            // 排序
	Note       string `json:"note" dc:"备注"`                            // 备注
}

// DictItemUpdateReq 修改字典项请求
type DictItemUpdateReq struct {
	gmeta.Meta `method:"put" path:"/api/dict/item/update" summary:"修改字典项" permission:"sys:dictitem:update" log:"字典项管理-修改记录" logType:"update"`
	Id         int64  `json:"id" v:"required#缺少ID" dc:"字典项ID"`       // 字典项ID
	DictId     int64  `json:"dictId" dc:"所属字典ID"`                    // 所属字典ID
	Name       string `json:"name" v:"required#请输入字典项名称" dc:"字典项名称"` // 字典项名称
	Value      string `json:"value" v:"required#请输入字典项值" dc:"字典项值"`  // 字典项值
	Sort       int    `json:"sort" dc:"排序"`                          // 排序
	Note       string `json:"note" dc:"备注"`                          // 备注
}

// DictItemUpdateRes 修改字典项结果
type DictItemUpdateRes struct{}

// DictItemDeleteReq 删除字典项请求
type DictItemDeleteReq struct {
	gmeta.Meta `method:"delete" path:"/api/dict/item/delete/:id" summary:"删除字典项" permission:"sys:dictitem:delete" log:"字典项管理-删除记录" logType:"delete"`
	Id         int64 `json:"id" v:"required#缺少ID" dc:"字典项ID"` // 字典项ID
}

// DictItemDeleteRes 删除字典项结果
type DictItemDeleteRes struct{}

// DictItemBatchDeleteReq 批量删除字典项请求
type DictItemBatchDeleteReq struct {
	gmeta.Meta `method:"delete" path:"/api/dict/item/batchDelete" summary:"批量删除字典项" permission:"sys:dictitem:batchDelete" log:"字典项管理-批量删除记录" logType:"delete"`
	Ids        []int64 `json:"ids" v:"required#请选择要删除的记录" dc:"字典项ID列表"` // 字典项ID列表
}

// DictItemBatchDeleteRes 批量删除结果
type DictItemBatchDeleteRes struct{}

// DictItemQuickListReq 快捷字典项列表请求
// 登录用户即可访问，按字典编码查询并走缓存
type DictItemQuickListReq struct {
	gmeta.Meta `method:"get" path:"/api/dict/item/getDictItemList/:code" summary:"按编码查询字典项" log:"字典项管理-查询列表" logType:"query"`
	Code       string `json:"code" v:"required#请输入字典编码" dc:"字典编码"` // 字典编码
}

// DictItemQuickRow 快捷字典项行数据
// 精简结构，仅包含前端下拉等场景所需字段
type DictItemQuickRow struct {
	Id    int64  `json:"id"`    // 主键ID
	Name  string `json:"name"`  // 字典项名称
	Value string `json:"value"` // 字典项值
}

// DictItemQuickListRes 快捷字典项列表结果
type DictItemQuickListRes = []DictItemQuickRow
