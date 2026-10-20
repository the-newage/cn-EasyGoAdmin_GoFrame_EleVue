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

// Package migrate 整库迁移
//
// 处理流程：源库内省 → 中间表示 → 目标方言 DDL → 分批搬运数据 → 序列修复。
// 独立于业务基座（不依赖 dao 与 logic），经 gdb 动态连接源库与目标库进行裸读裸写，
// 全程不触发全局钩子（租户隔离、软删除、审计盖章），保证数据原样落库。
package migrate

// ============================== 中间类型常量 ==============================

// 中间类型常量
//
// 源方言类型归一后的规范形态，各方言渲染器再翻译为目标方言
const (
	TypeBool       = "bool"
	TypeTinyint    = "tinyint"
	TypeSmallint   = "smallint"
	TypeInt        = "int"
	TypeBigint     = "bigint"
	TypeFloat      = "float"
	TypeDouble     = "double"
	TypeDecimal    = "decimal"
	TypeVarchar    = "varchar"
	TypeChar       = "char"
	TypeText       = "text"
	TypeTinytext   = "tinytext"
	TypeMediumText = "mediumtext"
	TypeLongText   = "longtext"
	TypeDate       = "date"
	TypeTime       = "time"
	TypeDateTime   = "datetime"
	TypeBlob       = "blob"
	TypeJson       = "json"
	TypeEnum       = "enum"
	TypeSet        = "set"
	TypeOther      = "other"
)

// ============================== 迁移任务参数 ==============================

// Options 迁移任务参数
type Options struct {
	Src DSN // 源库连接
	Dst DSN // 目标库连接

	DropTargetFirst bool     // 迁移前清空目标库全部表
	Tables          []string // 表白名单（空=全部）
	BatchSize       int      // 数据搬运批大小（默认 500）
}

// DSN 数据库连接要素
//
// SQLite 仅使用 Database 作为文件路径，Host 与 Port 忽略
type DSN struct {
	Type     string // 驱动名：mysql/pgsql/sqlite/mssql/oracle
	Host     string // 主机
	Port     string // 端口
	Database string // 数据库名或文件路径
	Username string // 用户名
	Password string // 密码
}

// ============================== 元数据中间表示 ==============================

// TableMeta 源库表元数据
type TableMeta struct {
	Name    string       // 表名
	Comment string       // 表注释
	Columns []ColumnMeta // 列集合
	Indexes []IndexMeta  // 二级索引集合
}

// ColumnMeta 列元数据
//
// Type 为中间类型常量
type ColumnMeta struct {
	Name     string   // 列名
	Type     string   // 中间类型
	Length   int      // varchar/char 长度；decimal 精度
	Scale    int      // decimal 小数位
	EnumVals []string // 枚举值集合
	Nullable bool     // 是否可空
	Default  string   // 默认值原文（空=无默认）；表达式类原样保留
	AutoInc  bool     // 是否自增
	Comment  string   // 列注释
	IsPK     bool     // 是否主键
}

// IndexMeta 索引元数据
//
// 不含主键，主键经 ColumnMeta.IsPK 表达
type IndexMeta struct {
	Unique  bool     // 是否唯一索引
	Name    string   // 索引名
	Columns []string // 索引列集合
}

// ============================== 迁移汇总 ==============================

// Report 迁移汇总
type Report struct {
	TableTotal  int           // 待迁移表总数
	TableOK     int           // 成功表数
	TableFailed int           // 失败表数
	RowTotal    int64         // 总搬运行数
	Failed      []FailedTable // 失败表记录
}

// FailedTable 失败表记录
type FailedTable struct {
	Table string // 表名
	Err   string // 错误信息
}
