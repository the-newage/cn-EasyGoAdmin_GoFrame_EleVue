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

package consts

// ============================== 日志类型与来源 ==============================

// 操作日志类型
//
// 对应操作日志表 type 列以及 gmeta logType 标签。
// 数值定义勿随意改动（前端操作日志类型按此映射显示）。
const (
	LogTypeAdd       = 1  // 新增
	LogTypeUpdate    = 2  // 修改
	LogTypeDelete    = 3  // 删除
	LogTypeQuery     = 4  // 查询
	LogTypeImport    = 5  // 导入
	LogTypeExport    = 6  // 导出
	LogTypeApprove   = 7  // 审批
	LogTypeReject    = 8  // 驳回
	LogTypeSubmit    = 9  // 提交
	LogTypeWithdraw  = 10 // 撤回
	LogTypeClose     = 11 // 关闭
	LogTypeOpen      = 12 // 开启
	LogTypeEnable    = 13 // 启用
	LogTypeDisable   = 14 // 禁用
	LogTypeUpload    = 15 // 上传
	LogTypeDownload  = 16 // 下载
	LogTypeSend      = 17 // 发送
	LogTypeSync      = 18 // 同步
	LogTypeGenerate  = 19 // 代码生成
	LogTypeClear     = 20 // 清空
	LogTypeRestore   = 21 // 还原
	LogTypeLogin     = 22 // 登录
	LogTypeLogout    = 23 // 登出
	LogTypeReset     = 24 // 重置
	LogTypeAssign    = 25 // 分配
	LogTypeAuthorize = 26 // 授权
	LogTypeExecute   = 27 // 执行
	LogTypePause     = 28 // 暂停
	LogTypeResume    = 29 // 恢复
	LogTypePublish   = 30 // 发布
	LogTypeUnpublish = 31 // 下架
	LogTypeCopy      = 32 // 复制
	LogTypeMove      = 33 // 移动
	LogTypeSort      = 34 // 排序
	LogTypeLock      = 35 // 锁定
	LogTypeUnlock    = 36 // 解锁
	LogTypeArchive   = 37 // 归档
	LogTypeUnarchive = 38 // 取消归档
	LogTypeBind      = 39 // 绑定
	LogTypeUnbind    = 40 // 解绑
	LogTypeOther     = 99 // 其他
	LogTypeAuto      = 0  // 自动推断：由中间件按 HTTP 方法判定（GET→查询，其余→修改）
)

// LogTypeNames 操作日志类型值到中文名的映射
// 供日志展示使用
var LogTypeNames = map[int]string{
	LogTypeAdd: "新增", LogTypeUpdate: "修改", LogTypeDelete: "删除", LogTypeQuery: "查询",
	LogTypeImport: "导入", LogTypeExport: "导出", LogTypeApprove: "审批", LogTypeReject: "驳回",
	LogTypeSubmit: "提交", LogTypeWithdraw: "撤回", LogTypeClose: "关闭", LogTypeOpen: "开启",
	LogTypeEnable: "启用", LogTypeDisable: "禁用", LogTypeUpload: "上传", LogTypeDownload: "下载",
	LogTypeSend: "发送", LogTypeSync: "同步", LogTypeGenerate: "代码生成", LogTypeClear: "清空",
	LogTypeRestore: "还原", LogTypeLogin: "登录", LogTypeLogout: "登出", LogTypeReset: "重置",
	LogTypeAssign: "分配", LogTypeAuthorize: "授权", LogTypeExecute: "执行", LogTypePause: "暂停",
	LogTypeResume: "恢复", LogTypePublish: "发布", LogTypeUnpublish: "下架", LogTypeCopy: "复制",
	LogTypeMove: "移动", LogTypeSort: "排序", LogTypeLock: "锁定", LogTypeUnlock: "解锁",
	LogTypeArchive: "归档", LogTypeUnarchive: "取消归档", LogTypeBind: "绑定", LogTypeUnbind: "解绑",
	LogTypeOther: "其他",
}

// LogTypeMap gmeta logType 标签字符串到枚举值的映射
// 空字符串映射为自动推断
var LogTypeMap = map[string]int{
	"add":       LogTypeAdd,
	"update":    LogTypeUpdate,
	"delete":    LogTypeDelete,
	"query":     LogTypeQuery,
	"import":    LogTypeImport,
	"export":    LogTypeExport,
	"approve":   LogTypeApprove,
	"reject":    LogTypeReject,
	"submit":    LogTypeSubmit,
	"withdraw":  LogTypeWithdraw,
	"enable":    LogTypeEnable,
	"disable":   LogTypeDisable,
	"upload":    LogTypeUpload,
	"download":  LogTypeDownload,
	"generate":  LogTypeGenerate,
	"login":     LogTypeLogin,
	"logout":    LogTypeLogout,
	"reset":     LogTypeReset,
	"assign":    LogTypeAssign,
	"authorize": LogTypeAuthorize,
	"execute":   LogTypeExecute,
	"pause":     LogTypePause,
	"resume":    LogTypeResume,
	"publish":   LogTypePublish,
	"sort":      LogTypeSort,
	"other":     LogTypeOther,
	"":          LogTypeAuto,
}

// 登录日志类型
// 对应登录日志表 log_type 列
const (
	LoginLogTypeLogin  = 1 // 登录
	LoginLogTypeLogout = 2 // 退出
)

// 登录来源
// 对应登录日志表 source 列
const (
	LoginSourceAdmin = 0 // 后台
	LoginSourceFront = 1 // 前台
)
