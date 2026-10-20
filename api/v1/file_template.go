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
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gmeta"
)

// ============================== 文件模板模块 ==============================

// FileTemplateItem 文件模板行数据
// 用于文件模板分页列表、文件模板列表等场景的模板基础信息
type FileTemplateItem struct {
	Id         int64       `json:"id"`         // 主键ID
	Name       string      `json:"name"`       // 模板名称
	Code       string      `json:"code"`       // 模板编码
	FileName   string      `json:"fileName"`   // 文件名
	FilePath   string      `json:"filePath"`   // 文件路径
	Note       string      `json:"note"`       // 备注
	Sort       int         `json:"sort"`       // 排序
	CreateUser string      `json:"createUser"` // 创建人
	CreateTime *gtime.Time `json:"createTime"` // 创建时间
}

// FileTemplatePageReq 文件模板分页请求
// 支持关键词、模板名称等条件组合查询
type FileTemplatePageReq struct {
	gmeta.Meta `method:"get" path:"/api/file/template/page" summary:"文件模板分页" permission:"sys:fileTemplate:page" log:"文件模板-查询分页记录" logType:"query"`
	PageReq           // 分页基础参数
	Keyword    string `json:"keyword" dc:"关键词（名称/编码模糊）"` // 关键词，名称或编码模糊匹配
	Name       string `json:"name" dc:"模板名称模糊"`          // 模板名称模糊匹配
}

// FileTemplatePageRes 文件模板分页结果
type FileTemplatePageRes = PageData[FileTemplateItem]

// FileTemplateListReq 文件模板列表请求
type FileTemplateListReq struct {
	gmeta.Meta `method:"get" path:"/api/file/template/list" summary:"文件模板列表" permission:"sys:fileTemplate:list" log:"文件模板-查询列表" logType:"query"`
	Keyword    string `json:"keyword" dc:"关键词"` // 关键词
}

// FileTemplateListRes 文件模板列表结果
type FileTemplateListRes = []FileTemplateItem

// FileTemplateDetailReq 文件模板详情请求
type FileTemplateDetailReq struct {
	gmeta.Meta `method:"get" path:"/api/file/template/detail/:id" summary:"文件模板详情" permission:"sys:fileTemplate:detail" log:"文件模板-查询详情" logType:"query"`
	Id         int64 `json:"id" v:"required#缺少ID" dc:"模板ID"` // 模板ID
}

// FileTemplateDetailRes 文件模板详情结果
type FileTemplateDetailRes struct {
	Id       int64  `json:"id"`       // 主键ID
	Name     string `json:"name"`     // 模板名称
	Code     string `json:"code"`     // 模板编码
	FileName string `json:"fileName"` // 文件名
	FilePath string `json:"filePath"` // 文件路径
	Note     string `json:"note"`     // 备注
	Sort     int    `json:"sort"`     // 排序
}

// FileTemplateAddReq 新增文件模板请求
type FileTemplateAddReq struct {
	gmeta.Meta `method:"post" path:"/api/file/template/add" summary:"新增文件模板" permission:"sys:fileTemplate:add" log:"文件模板-新增记录" logType:"add"`
	Name       string `json:"name" v:"required#请输入模板名称" dc:"模板名称"` // 模板名称
	Code       string `json:"code" v:"required#请输入模板编码" dc:"模板编码"` // 模板编码
	FileName   string `json:"fileName" dc:"文件名"`                   // 文件名
	FilePath   string `json:"filePath" dc:"文件路径"`                  // 文件路径
	Note       string `json:"note" dc:"备注"`                        // 备注
	Sort       int    `json:"sort" dc:"排序"`                        // 排序
}

// FileTemplateUpdateReq 修改文件模板请求
type FileTemplateUpdateReq struct {
	gmeta.Meta `method:"put" path:"/api/file/template/update" summary:"修改文件模板" permission:"sys:fileTemplate:update" log:"文件模板-修改记录" logType:"update"`
	Id         int64  `json:"id" v:"required#缺少ID" dc:"模板ID"`      // 模板ID
	Name       string `json:"name" v:"required#请输入模板名称" dc:"模板名称"` // 模板名称
	Code       string `json:"code" v:"required#请输入模板编码" dc:"模板编码"` // 模板编码
	FileName   string `json:"fileName" dc:"文件名"`                   // 文件名
	FilePath   string `json:"filePath" dc:"文件路径"`                  // 文件路径
	Note       string `json:"note" dc:"备注"`                        // 备注
	Sort       int    `json:"sort" dc:"排序"`                        // 排序
}

// FileTemplateUpdateRes 修改文件模板结果
type FileTemplateUpdateRes struct{}

// FileTemplateDeleteReq 删除文件模板请求
type FileTemplateDeleteReq struct {
	gmeta.Meta `method:"delete" path:"/api/file/template/delete/:id" summary:"删除文件模板" permission:"sys:fileTemplate:delete" log:"文件模板-删除记录" logType:"delete"`
	Id         int64 `json:"id" v:"required#缺少ID" dc:"模板ID"` // 模板ID
}

// FileTemplateDeleteRes 删除文件模板结果
type FileTemplateDeleteRes struct{}

// FileTemplateBatchDeleteReq 批量删除文件模板请求
type FileTemplateBatchDeleteReq struct {
	gmeta.Meta `method:"delete" path:"/api/file/template/batchDelete" summary:"批量删除文件模板" permission:"sys:fileTemplate:batchDelete" log:"文件模板-批量删除记录" logType:"delete"`
	Ids        []int64 `json:"ids" v:"required#请选择要删除的记录" dc:"模板ID列表"` // 模板ID列表
}

// FileTemplateBatchDeleteRes 批量删除结果
type FileTemplateBatchDeleteRes struct{}

// FileTemplateByCodeReq 按编码查询模板请求
// 登录用户即可访问，按模板编码查询模板详情
type FileTemplateByCodeReq struct {
	gmeta.Meta `method:"get" path:"/api/file/template/getTemplateByCode/:code" summary:"按编码查询模板" log:"文件模板-查询模板详情" logType:"query"`
	Code       string `json:"code" v:"required#请输入模板编码" dc:"模板编码"` // 模板编码
}

// FileTemplateByCodeRes 按编码查询模板结果
type FileTemplateByCodeRes struct {
	FileTemplateItem
}
