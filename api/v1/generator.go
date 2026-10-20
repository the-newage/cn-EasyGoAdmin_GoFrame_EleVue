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

// ============================== 代码生成模块（读数据库结构） ==============================

// GeneratorTableItem 数据表行数据
// 用于代码生成模块的数据表分页列表，描述数据库表的元信息
type GeneratorTableItem struct {
	TableName      string `json:"tableName"`      // 表名
	TableComment   string `json:"tableComment"`   // 表注释
	Engine         string `json:"engine"`         // 存储引擎
	TableRows      int64  `json:"tableRows"`      // 表记录数
	DataLength     int64  `json:"dataLength"`     // 数据大小（字节）
	AutoIncrement  int64  `json:"autoIncrement"`  // 自增当前值
	TableCollation string `json:"tableCollation"` // 表排序规则
	CreateTime     string `json:"createTime"`     // 创建时间
}

// GeneratorPageReq 数据表分页请求
// 支持按表名、表注释模糊查询数据库表
type GeneratorPageReq struct {
	gmeta.Meta   `method:"get" path:"/api/generator/page" summary:"查询数据表" permission:"sys:generator:page" log:"代码生成-查询数据表" logType:"query"`
	PageReq             // 分页基础参数
	TableName    string `json:"tableName" dc:"表名模糊"`     // 表名模糊匹配
	TableComment string `json:"tableComment" dc:"表注释模糊"` // 表注释模糊匹配
}

// GeneratorPageRes 数据表分页结果
type GeneratorPageRes = PageData[GeneratorTableItem]

// GeneratorPreviewReq 预览代码请求
// TableNames 元素格式为 "表名|表注释"
type GeneratorPreviewReq struct {
	gmeta.Meta `method:"post" path:"/api/generator/preview" summary:"预览代码" permission:"sys:generator:generator" log:"代码生成-预览代码" logType:"generate"`
	TableNames []string `json:"tableNames" v:"required#请选择数据表" dc:"数据表列表（表名|表注释）"` // 数据表列表
}

// GeneratorFileItem 生成文件
// 描述单个待生成文件的路径与内容
type GeneratorFileItem struct {
	Path    string `json:"path"`    // 文件路径
	Content string `json:"content"` // 文件内容
}

// GeneratorPreviewRes 预览结果
type GeneratorPreviewRes struct {
	List []GeneratorFileItem `json:"list"` // 生成文件列表
}

// GeneratorGenerateReq 一键生成请求
// TableNames 元素格式为 "表名|表注释"
type GeneratorGenerateReq struct {
	gmeta.Meta `method:"post" path:"/api/generator/generate" summary:"一键生成" permission:"sys:generator:generator" log:"代码生成-一键生成" logType:"generate"`
	TableNames []string `json:"tableNames" v:"required#请选择数据表" dc:"数据表列表（表名|表注释）"` // 数据表列表
}

// GeneratorMenuResult 同步创建菜单结果
// 描述代码生成过程中同步创建菜单的执行结果
type GeneratorMenuResult struct {
	Success bool   `json:"success"` // 是否成功
	MenuId  int64  `json:"menuId"`  // 菜单ID
	Message string `json:"message"` // 结果消息
}

// GeneratorDictResult 同步创建字典结果
// 描述代码生成过程中同步创建字典的执行结果
type GeneratorDictResult struct {
	Dict   string `json:"dict"`            // 字典编码
	Name   string `json:"name,omitempty"`  // 字典名称
	Items  int    `json:"items,omitempty"` // 字典项数量
	Status string `json:"status"`          // 执行状态
}

// GeneratorRouteResult 同步注册路由结果
// 描述代码生成过程中同步注册路由的执行结果
type GeneratorRouteResult struct {
	Success bool   `json:"success"` // 是否成功
	Message string `json:"message"` // 结果消息
}

// GeneratorGenerateRes 生成结果
// 汇总代码生成的文件写入情况以及菜单、字典、路由的同步结果
type GeneratorGenerateRes struct {
	Generated   []string              `json:"generated"`   // 新生成的文件
	Skipped     []string              `json:"skipped"`     // 跳过的文件
	Overwritten []string              `json:"overwritten"` // 覆盖的文件
	Failed      []string              `json:"failed"`      // 失败的文件
	Menu        GeneratorMenuResult   `json:"menu"`        // 菜单创建结果
	Dict        []GeneratorDictResult `json:"dict"`        // 字典创建结果
	Route       GeneratorRouteResult  `json:"route"`       // 路由注册结果
}

// GeneratorBatchGenerateReq 批量生成请求
// TableNames 元素格式为 "表名|表注释"
type GeneratorBatchGenerateReq struct {
	gmeta.Meta `method:"post" path:"/api/generator/batchGenerate" summary:"批量生成" permission:"sys:generator:batchGenerator" log:"代码生成-批量生成" logType:"generate"`
	TableNames []string `json:"tableNames" v:"required#请选择数据表" dc:"数据表列表（表名|表注释）"` // 数据表列表
}

// GeneratorBatchItem 单表生成结果
// 描述批量生成中单张表的执行情况
type GeneratorBatchItem struct {
	Table       string                `json:"table"`       // 表名
	Success     bool                  `json:"success"`     // 是否成功
	Generated   int                   `json:"generated"`   // 新生成文件数
	Skipped     int                   `json:"skipped"`     // 跳过文件数
	Overwritten int                   `json:"overwritten"` // 覆盖文件数
	Failed      int                   `json:"failed"`      // 失败文件数
	Message     string                `json:"message"`     // 结果消息
	Menu        GeneratorMenuResult   `json:"menu"`        // 菜单创建结果
	Dict        []GeneratorDictResult `json:"dict"`        // 字典创建结果
	Route       GeneratorRouteResult  `json:"route"`       // 路由注册结果
}

// GeneratorBatchGenerateRes 批量生成结果
type GeneratorBatchGenerateRes struct {
	List []GeneratorBatchItem `json:"list"` // 各表生成结果列表
}
