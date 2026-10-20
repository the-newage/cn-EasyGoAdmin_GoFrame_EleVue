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

package migrate

import (
	"context"
	"fmt"
	"strings"

	"github.com/gogf/gf/v2/database/gdb"
)

// ============================== 目标方言 DDL 渲染 ==============================

// Renderer 目标方言 DDL 渲染器
// 负责标识符引用、列定义、建表、建索引、清表与自增序列修正
type Renderer interface {
	Quote(name string) string
	Column(col ColumnMeta) string
	CreateTable(t *TableMeta) string
	CreateIndex(t *TableMeta, idx IndexMeta) string
	DropAllTables(ctx context.Context, db gdb.DB) error
	FixSequence(ctx context.Context, db gdb.DB, t *TableMeta) error
}

// NewRenderer 按目标方言创建渲染器
func NewRenderer(dialect string) (Renderer, error) {
	switch dialect {
	case "mysql":
		return &mysqlRenderer{dialect: "mysql"}, nil
	case "pgsql":
		return &pgsqlRenderer{dialect: "pgsql"}, nil
	case "sqlite":
		return &sqliteRenderer{dialect: "sqlite"}, nil
	case "mssql":
		return &mssqlRenderer{dialect: "mssql"}, nil
	case "oracle":
		return &oracleRenderer{dialect: "oracle"}, nil
	}
	return nil, fmt.Errorf("不支持的目标数据库类型: %s", dialect)
}

// quoteStd 双引号引用
// 适用于 PostgreSQL、Oracle、SQLite
func quoteStd(name string) string {
	return `"` + strings.ReplaceAll(name, `"`, `""`) + `"`
}

// quoteTick 反引号引用
// 适用于 MySQL
func quoteTick(name string) string {
	return "`" + strings.ReplaceAll(name, "`", "``") + "`"
}

// quoteBracket 方括号引用
// 适用于 SQL Server
func quoteBracket(name string) string {
	return "[" + strings.ReplaceAll(name, "]", "]]") + "]"
}

// defaultLiteral 渲染默认值字面量
//
// 数值原样输出，其余加引号；
// 时间函数按目标方言归一转换
func defaultLiteral(dialect, v string) string {
	if v == "" {
		return ""
	}
	u := strings.ToUpper(v)
	// 时间函数按目标方言归一（NOW/CURRENT_TIMESTAMP → 各库原生）
	if u == "CURRENT_TIMESTAMP" || u == "NOW()" || u == "NOW" {
		switch dialect {
		case "mssql":
			return "GETDATE()"
		case "oracle":
			return "SYSDATE"
		default:
			return "CURRENT_TIMESTAMP"
		}
	}
	if u == "NULL" {
		return "NULL"
	}
	// 纯数值/负数原样
	if _, err := fmt.Sscanf(v, "%f", new(float64)); err == nil {
		return v
	}
	return "'" + strings.ReplaceAll(v, "'", "''") + "'"
}

// hasAutoIncCol 判断表中是否存在自增列
func hasAutoIncCol(t *TableMeta) bool {
	for _, c := range t.Columns {
		if c.AutoInc {
			return true
		}
	}
	return false
}

// pkColumns 返回表的主键列清单
func pkColumns(t *TableMeta) []string {
	var cols []string
	for _, c := range t.Columns {
		if c.IsPK {
			cols = append(cols, c.Name)
		}
	}
	return cols
}

// ============================================================
// MySQL
// ============================================================

// mysqlRenderer MySQL 方言渲染器
type mysqlRenderer struct{ dialect string }

// Quote 引用标识符
func (r *mysqlRenderer) Quote(name string) string { return quoteTick(name) }

// Column 渲染列定义
func (r *mysqlRenderer) Column(col ColumnMeta) string {
	var b strings.Builder
	b.WriteString(quoteTick(col.Name) + " " + r.colType(col))
	if col.AutoInc {
		b.WriteString(" NOT NULL AUTO_INCREMENT")
	} else {
		if !col.Nullable {
			b.WriteString(" NOT NULL")
		}
		if d := defaultLiteral(r.dialect, col.Default); d != "" {
			b.WriteString(" DEFAULT " + d)
		}
	}
	if col.Comment != "" {
		b.WriteString(" COMMENT '" + strings.ReplaceAll(col.Comment, "'", "''") + "'")
	}
	return b.String()
}

// colType 列类型映射
func (r *mysqlRenderer) colType(col ColumnMeta) string {
	switch col.Type {
	case TypeBool:
		return "tinyint(1)"
	case TypeTinyint:
		return "tinyint"
	case TypeSmallint:
		return "smallint"
	case TypeInt:
		return "int"
	case TypeBigint:
		return "bigint"
	case TypeFloat:
		return "float"
	case TypeDouble:
		return "double"
	case TypeDecimal:
		if col.Length > 0 {
			return fmt.Sprintf("decimal(%d,%d)", col.Length, col.Scale)
		}
		return "decimal"
	case TypeVarchar:
		if col.Length > 0 {
			return fmt.Sprintf("varchar(%d)", col.Length)
		}
		return "varchar(255)"
	case TypeChar:
		if col.Length > 0 {
			return fmt.Sprintf("char(%d)", col.Length)
		}
		return "char(36)"
	case TypeTinytext:
		return "tinytext"
	case TypeMediumText:
		return "mediumtext"
	case TypeLongText:
		return "longtext"
	case TypeText:
		return "text"
	case TypeDate:
		return "date"
	case TypeTime:
		return "time"
	case TypeDateTime:
		return "datetime"
	case TypeBlob:
		return "longblob"
	case TypeJson:
		return "json"
	case TypeEnum:
		return "enum('" + strings.Join(col.EnumVals, "','") + "')"
	case TypeSet:
		return "set('" + strings.Join(col.EnumVals, "','") + "')"
	}
	return "text"
}

// CreateTable 渲染建表语句
func (r *mysqlRenderer) CreateTable(t *TableMeta) string {
	var lines []string
	for _, c := range t.Columns {
		lines = append(lines, "  "+r.Column(c))
	}
	if pk := pkColumns(t); len(pk) > 0 {
		quoted := make([]string, len(pk))
		for i, c := range pk {
			quoted[i] = quoteTick(c)
		}
		lines = append(lines, "  PRIMARY KEY ("+strings.Join(quoted, ",")+")")
	}
	sql := "CREATE TABLE IF NOT EXISTS " + quoteTick(t.Name) + " (\n" + strings.Join(lines, ",\n") + "\n)"
	if t.Comment != "" {
		sql += " COMMENT='" + strings.ReplaceAll(t.Comment, "'", "''") + "'"
	}
	return sql + ";"
}

// CreateIndex 渲染索引语句
func (r *mysqlRenderer) CreateIndex(t *TableMeta, idx IndexMeta) string {
	unique := ""
	if idx.Unique {
		unique = "UNIQUE "
	}
	cols := make([]string, len(idx.Columns))
	for i, c := range idx.Columns {
		cols[i] = quoteTick(c)
	}
	return fmt.Sprintf("CREATE %sINDEX %s ON %s (%s);",
		unique, quoteTick(idx.Name), quoteTick(t.Name), strings.Join(cols, ","))
}

// DropAllTables 清空全部表
func (r *mysqlRenderer) DropAllTables(ctx context.Context, db gdb.DB) error {
	if _, err := db.Exec(ctx, "SET FOREIGN_KEY_CHECKS=0"); err != nil {
		return err
	}
	defer db.Exec(ctx, "SET FOREIGN_KEY_CHECKS=1")
	tables, err := db.Tables(ctx)
	if err != nil {
		return err
	}
	for _, t := range tables {
		if _, err = db.Exec(ctx, "DROP TABLE IF EXISTS "+quoteTick(t)); err != nil {
			return err
		}
	}
	return nil
}

// FixSequence 修正自增序列
// MySQL 自增自动续接，无需处理
func (r *mysqlRenderer) FixSequence(ctx context.Context, db gdb.DB, t *TableMeta) error { return nil }

// ============================================================
// PostgreSQL
// ============================================================

// pgsqlRenderer PostgreSQL 方言渲染器
type pgsqlRenderer struct{ dialect string }

// Quote 引用标识符
func (r *pgsqlRenderer) Quote(name string) string { return quoteStd(name) }

// Column 渲染列定义
func (r *pgsqlRenderer) Column(col ColumnMeta) string {
	var b strings.Builder
	b.WriteString(quoteStd(col.Name) + " " + r.colType(col))
	if col.AutoInc {
		b.WriteString(" GENERATED BY DEFAULT AS IDENTITY")
	}
	if !col.Nullable && !col.AutoInc {
		b.WriteString(" NOT NULL")
	}
	if !col.AutoInc {
		if d := defaultLiteral(r.dialect, col.Default); d != "" {
			b.WriteString(" DEFAULT " + d)
		}
	}
	return b.String()
}

// colType 列类型映射
func (r *pgsqlRenderer) colType(col ColumnMeta) string {
	switch col.Type {
	case TypeBool:
		return "BOOLEAN"
	case TypeTinyint, TypeSmallint:
		return "SMALLINT"
	case TypeInt:
		return "INTEGER"
	case TypeBigint:
		return "BIGINT"
	case TypeFloat:
		return "REAL"
	case TypeDouble:
		return "DOUBLE PRECISION"
	case TypeDecimal:
		if col.Length > 0 {
			return fmt.Sprintf("NUMERIC(%d,%d)", col.Length, col.Scale)
		}
		return "NUMERIC"
	case TypeVarchar:
		if col.Length > 0 {
			return fmt.Sprintf("VARCHAR(%d)", col.Length)
		}
		return "VARCHAR"
	case TypeChar:
		if col.Length > 0 {
			return fmt.Sprintf("CHAR(%d)", col.Length)
		}
		return "CHAR(36)"
	case TypeText, TypeTinytext, TypeMediumText, TypeLongText:
		return "TEXT"
	case TypeDate:
		return "DATE"
	case TypeTime:
		return "TIME"
	case TypeDateTime:
		return "TIMESTAMP"
	case TypeBlob:
		return "BYTEA"
	case TypeJson:
		return "JSONB"
	case TypeEnum, TypeSet:
		return "VARCHAR(255)" // pgsql 无内建 enum/set，降级为 varchar（值域由应用层保证）
	}
	return "TEXT"
}

// CreateTable 渲染建表语句
func (r *pgsqlRenderer) CreateTable(t *TableMeta) string {
	var lines []string
	for _, c := range t.Columns {
		lines = append(lines, "  "+r.Column(c))
	}
	if pk := pkColumns(t); len(pk) > 0 {
		quoted := make([]string, len(pk))
		for i, c := range pk {
			quoted[i] = quoteStd(c)
		}
		lines = append(lines, "  PRIMARY KEY ("+strings.Join(quoted, ",")+")")
	}
	sql := "CREATE TABLE IF NOT EXISTS " + quoteStd(t.Name) + " (\n" + strings.Join(lines, ",\n") + "\n);"
	if t.Comment != "" {
		sql += "\nCOMMENT ON TABLE " + quoteStd(t.Name) + " IS '" + strings.ReplaceAll(t.Comment, "'", "''") + "';"
	}
	// 列注释随建表语句一并输出
	for _, c := range t.Columns {
		if c.Comment != "" {
			sql += "\nCOMMENT ON COLUMN " + quoteStd(t.Name) + "." + quoteStd(c.Name) +
				" IS '" + strings.ReplaceAll(c.Comment, "'", "''") + "';"
		}
	}
	return sql
}

// CreateIndex 渲染索引语句
func (r *pgsqlRenderer) CreateIndex(t *TableMeta, idx IndexMeta) string {
	unique := ""
	if idx.Unique {
		unique = "UNIQUE "
	}
	cols := make([]string, len(idx.Columns))
	for i, c := range idx.Columns {
		cols[i] = quoteStd(c)
	}
	return fmt.Sprintf("CREATE %sINDEX IF NOT EXISTS %s ON %s (%s);",
		unique, quoteStd(idx.Name), quoteStd(t.Name), strings.Join(cols, ","))
}

// DropAllTables 清空全部表
func (r *pgsqlRenderer) DropAllTables(ctx context.Context, db gdb.DB) error {
	rows, err := db.Ctx(ctx).Raw(
		"SELECT tablename FROM pg_tables WHERE schemaname = current_schema()").All()
	if err != nil {
		return err
	}
	for _, row := range rows {
		if _, err = db.Exec(ctx, "DROP TABLE IF EXISTS "+quoteStd(row["tablename"].String())+" CASCADE"); err != nil {
			return err
		}
	}
	return nil
}

// FixSequence 修正自增序列
//
// 将序列设置到当前 max(id)，下次插入取 max+1；
// 空表不处理，IDENTITY 从 1 开始
func (r *pgsqlRenderer) FixSequence(ctx context.Context, db gdb.DB, t *TableMeta) error {
	for _, c := range t.Columns {
		if !c.AutoInc {
			continue
		}
		// setval 到当前 max(id)，下次插入取 max+1；空表不动（IDENTITY 从1开始）
		_, err := db.Exec(ctx, fmt.Sprintf(
			`SELECT setval(pg_get_serial_sequence('%s', '%s'), COALESCE((SELECT MAX(%s) FROM %s), 0) + 1, false)`,
			t.Name, c.Name, quoteStd(c.Name), quoteStd(t.Name)))
		if err != nil {
			return err
		}
	}
	return nil
}

// ============================================================
// SQLite
// ============================================================

// sqliteRenderer SQLite 方言渲染器
type sqliteRenderer struct{ dialect string }

// Quote 引用标识符
func (r *sqliteRenderer) Quote(name string) string { return quoteStd(name) }

// Column 渲染列定义
func (r *sqliteRenderer) Column(col ColumnMeta) string {
	var b strings.Builder
	b.WriteString(quoteStd(col.Name) + " " + r.colType(col))
	if col.IsPK && col.AutoInc {
		b.WriteString(" PRIMARY KEY AUTOINCREMENT")
	} else {
		if col.IsPK {
			b.WriteString(" PRIMARY KEY")
		}
		if !col.Nullable {
			b.WriteString(" NOT NULL")
		}
		if d := defaultLiteral(r.dialect, col.Default); d != "" {
			b.WriteString(" DEFAULT " + d)
		}
	}
	return b.String()
}

// colType 列类型映射
func (r *sqliteRenderer) colType(col ColumnMeta) string {
	switch col.Type {
	case TypeBool, TypeTinyint, TypeSmallint, TypeInt, TypeBigint, TypeDate, TypeTime, TypeDateTime:
		return "INTEGER"
	case TypeFloat, TypeDouble, TypeDecimal:
		return "REAL"
	case TypeBlob:
		return "BLOB"
	default:
		return "TEXT"
	}
}

// CreateTable 渲染建表语句
func (r *sqliteRenderer) CreateTable(t *TableMeta) string {
	var lines []string
	for _, c := range t.Columns {
		lines = append(lines, "  "+r.Column(c))
	}
	// 复合主键走表级约束（列级 PRIMARY KEY 仅单列生效）
	if pk := pkColumns(t); len(pk) > 1 || !hasAutoIncCol(t) {
		if pk := pkColumns(t); len(pk) > 0 {
			quoted := make([]string, len(pk))
			for i, c := range pk {
				quoted[i] = quoteStd(c)
			}
			lines = append(lines, "  PRIMARY KEY ("+strings.Join(quoted, ",")+")")
		}
	}
	// sqlite 无表/列注释
	return "CREATE TABLE IF NOT EXISTS " + quoteStd(t.Name) + " (\n" + strings.Join(lines, ",\n") + "\n);"
}

// CreateIndex 渲染索引语句
func (r *sqliteRenderer) CreateIndex(t *TableMeta, idx IndexMeta) string {
	unique := ""
	if idx.Unique {
		unique = "UNIQUE "
	}
	cols := make([]string, len(idx.Columns))
	for i, c := range idx.Columns {
		cols[i] = quoteStd(c)
	}
	return fmt.Sprintf("CREATE %sINDEX IF NOT EXISTS %s ON %s (%s);",
		unique, quoteStd(idx.Name), quoteStd(t.Name), strings.Join(cols, ","))
}

// DropAllTables 清空全部表
func (r *sqliteRenderer) DropAllTables(ctx context.Context, db gdb.DB) error {
	rows, err := db.Ctx(ctx).Raw(
		"SELECT name FROM sqlite_master WHERE type='table' AND name NOT LIKE 'sqlite_%'").All()
	if err != nil {
		return err
	}
	if _, err = db.Exec(ctx, "PRAGMA foreign_keys = OFF"); err != nil {
		return err
	}
	defer db.Exec(ctx, "PRAGMA foreign_keys = ON")
	for _, row := range rows {
		if _, err = db.Exec(ctx, "DROP TABLE IF EXISTS "+quoteStd(row["name"].String())); err != nil {
			return err
		}
	}
	return nil
}

// FixSequence 修正自增序列
// AUTOINCREMENT 自动续接，无需处理
func (r *sqliteRenderer) FixSequence(ctx context.Context, db gdb.DB, t *TableMeta) error { return nil }

// ============================================================
// SQL Server（首期实测链路为 mysql→pgsql，mssql/oracle 按定稿方言约定实现，待人工校对）
// ============================================================

// mssqlRenderer SQL Server 方言渲染器
type mssqlRenderer struct{ dialect string }

// Quote 引用标识符
func (r *mssqlRenderer) Quote(name string) string { return quoteBracket(name) }

// Column 渲染列定义
func (r *mssqlRenderer) Column(col ColumnMeta) string {
	var b strings.Builder
	b.WriteString(quoteBracket(col.Name) + " " + r.colType(col))
	if col.AutoInc {
		b.WriteString(" IDENTITY(1,1)")
	}
	if !col.Nullable {
		b.WriteString(" NOT NULL")
	}
	if !col.AutoInc {
		if d := defaultLiteral(r.dialect, col.Default); d != "" {
			b.WriteString(" DEFAULT " + d)
		}
	}
	return b.String()
}

// colType 列类型映射
func (r *mssqlRenderer) colType(col ColumnMeta) string {
	switch col.Type {
	case TypeBool:
		return "BIT"
	case TypeTinyint:
		return "TINYINT"
	case TypeSmallint:
		return "SMALLINT"
	case TypeInt:
		return "INT"
	case TypeBigint:
		return "BIGINT"
	case TypeFloat:
		return "FLOAT"
	case TypeDouble:
		return "FLOAT"
	case TypeDecimal:
		if col.Length > 0 {
			return fmt.Sprintf("DECIMAL(%d,%d)", col.Length, col.Scale)
		}
		return "DECIMAL"
	case TypeVarchar:
		if col.Length > 0 {
			return fmt.Sprintf("NVARCHAR(%d)", col.Length)
		}
		return "NVARCHAR(255)"
	case TypeChar:
		if col.Length > 0 {
			return fmt.Sprintf("NCHAR(%d)", col.Length)
		}
		return "NCHAR(36)"
	case TypeText, TypeTinytext, TypeMediumText, TypeLongText, TypeJson, TypeEnum, TypeSet:
		return "NVARCHAR(MAX)"
	case TypeDate:
		return "DATE"
	case TypeTime:
		return "TIME"
	case TypeDateTime:
		return "DATETIME2"
	case TypeBlob:
		return "VARBINARY(MAX)"
	}
	return "NVARCHAR(MAX)"
}

// CreateTable 渲染建表语句
func (r *mssqlRenderer) CreateTable(t *TableMeta) string {
	var lines []string
	for _, c := range t.Columns {
		lines = append(lines, "  "+r.Column(c))
	}
	if pk := pkColumns(t); len(pk) > 0 {
		quoted := make([]string, len(pk))
		for i, c := range pk {
			quoted[i] = quoteBracket(c)
		}
		lines = append(lines, "  PRIMARY KEY ("+strings.Join(quoted, ",")+")")
	}
	sql := "IF OBJECT_ID(N'" + t.Name + "', N'U') IS NULL CREATE TABLE " + quoteBracket(t.Name) +
		" (\n" + strings.Join(lines, ",\n") + "\n);"
	// 表注释（扩展属性 MS_Description）
	if t.Comment != "" {
		sql += "\nEXEC sp_addextendedproperty N'MS_Description', N'" +
			strings.ReplaceAll(t.Comment, "'", "''") +
			"', N'SCHEMA', N'dbo', N'TABLE', N'" + t.Name + "';"
	}
	// 列注释（扩展属性 MS_Description）
	for _, c := range t.Columns {
		if c.Comment != "" {
			sql += "\nEXEC sp_addextendedproperty N'MS_Description', N'" +
				strings.ReplaceAll(c.Comment, "'", "''") +
				"', N'SCHEMA', N'dbo', N'TABLE', N'" + t.Name + "', N'COLUMN', N'" + c.Name + "';"
		}
	}
	return sql
}

// CreateIndex 渲染索引语句
func (r *mssqlRenderer) CreateIndex(t *TableMeta, idx IndexMeta) string {
	unique := ""
	if idx.Unique {
		unique = "UNIQUE "
	}
	cols := make([]string, len(idx.Columns))
	for i, c := range idx.Columns {
		cols[i] = quoteBracket(c)
	}
	// IF NOT EXISTS 保护：重复建索引跳过
	return fmt.Sprintf("IF NOT EXISTS (SELECT * FROM sys.indexes WHERE name = N'%s' AND object_id = OBJECT_ID(N'%s')) "+
		"CREATE %sINDEX %s ON %s (%s);",
		idx.Name, t.Name, unique, quoteBracket(idx.Name), quoteBracket(t.Name), strings.Join(cols, ","))
}

// DropAllTables 清空全部表
func (r *mssqlRenderer) DropAllTables(ctx context.Context, db gdb.DB) error {
	// 先禁用全部外键约束，再逐表 DROP
	fks, err := db.Ctx(ctx).Raw(
		`SELECT OBJECT_NAME(parent_object_id) AS tbl, name AS fk FROM sys.foreign_keys`).All()
	if err != nil {
		return err
	}
	for _, fk := range fks {
		if _, err = db.Exec(ctx, fmt.Sprintf("ALTER TABLE %s DROP CONSTRAINT %s",
			quoteBracket(fk["tbl"].String()), quoteBracket(fk["fk"].String()))); err != nil {
			return err
		}
	}
	tables, err := db.Tables(ctx)
	if err != nil {
		return err
	}
	for _, t := range tables {
		if _, err = db.Exec(ctx, "DROP TABLE "+quoteBracket(t)); err != nil {
			return err
		}
	}
	return nil
}

// FixSequence 修正自增序列
func (r *mssqlRenderer) FixSequence(ctx context.Context, db gdb.DB, t *TableMeta) error {
	for _, c := range t.Columns {
		if !c.AutoInc {
			continue
		}
		if _, err := db.Exec(ctx, fmt.Sprintf(
			"DBCC CHECKIDENT ('%s', RESEED)", t.Name)); err != nil {
			return err
		}
	}
	return nil
}

// ============================================================
// Oracle
// ============================================================

// oracleRenderer Oracle 方言渲染器
type oracleRenderer struct{ dialect string }

// Quote 引用标识符
func (r *oracleRenderer) Quote(name string) string { return quoteStd(name) }

// Column 渲染列定义
func (r *oracleRenderer) Column(col ColumnMeta) string {
	var b strings.Builder
	b.WriteString(quoteStd(col.Name) + " " + r.colType(col))
	if col.AutoInc {
		b.WriteString(" GENERATED BY DEFAULT AS IDENTITY")
	}
	if !col.Nullable && !col.AutoInc {
		b.WriteString(" NOT NULL")
	}
	if !col.AutoInc {
		if d := defaultLiteral(r.dialect, col.Default); d != "" {
			b.WriteString(" DEFAULT " + d)
		}
	}
	return b.String()
}

// colType 列类型映射
func (r *oracleRenderer) colType(col ColumnMeta) string {
	switch col.Type {
	case TypeBool, TypeTinyint, TypeSmallint:
		return "NUMBER(3)"
	case TypeInt:
		return "NUMBER(10)"
	case TypeBigint:
		return "NUMBER(19)"
	case TypeFloat:
		return "BINARY_FLOAT"
	case TypeDouble:
		return "BINARY_DOUBLE"
	case TypeDecimal:
		if col.Length > 0 {
			return fmt.Sprintf("NUMBER(%d,%d)", col.Length, col.Scale)
		}
		return "NUMBER"
	case TypeVarchar:
		if col.Length > 0 {
			return fmt.Sprintf("NVARCHAR2(%d)", col.Length)
		}
		return "NVARCHAR2(255)"
	case TypeChar:
		if col.Length > 0 {
			return fmt.Sprintf("NCHAR(%d)", col.Length)
		}
		return "NCHAR(36)"
	case TypeText, TypeTinytext, TypeMediumText, TypeLongText, TypeJson, TypeEnum, TypeSet:
		return "NCLOB"
	case TypeDate:
		return "DATE"
	case TypeTime, TypeDateTime:
		return "TIMESTAMP"
	case TypeBlob:
		return "BLOB"
	}
	return "NCLOB"
}

// CreateTable 渲染建表语句
func (r *oracleRenderer) CreateTable(t *TableMeta) string {
	var lines []string
	for _, c := range t.Columns {
		lines = append(lines, "  "+r.Column(c))
	}
	if pk := pkColumns(t); len(pk) > 0 {
		quoted := make([]string, len(pk))
		for i, c := range pk {
			quoted[i] = quoteStd(c)
		}
		lines = append(lines, "  PRIMARY KEY ("+strings.Join(quoted, ",")+")")
	}
	sql := "CREATE TABLE " + quoteStd(t.Name) + " (\n" + strings.Join(lines, ",\n") + "\n)"
	if t.Comment != "" {
		sql += "\nCOMMENT ON TABLE " + quoteStd(t.Name) + " IS '" + strings.ReplaceAll(t.Comment, "'", "''") + "'"
	}
	for _, c := range t.Columns {
		if c.Comment != "" {
			sql += "\nCOMMENT ON COLUMN " + quoteStd(t.Name) + "." + quoteStd(c.Name) +
				" IS '" + strings.ReplaceAll(c.Comment, "'", "''") + "'"
		}
	}
	return sql + ";"
}

// CreateIndex 渲染索引语句
//
// Oracle 无 IF NOT EXISTS，由编排层预检索引存在性后决定是否执行
func (r *oracleRenderer) CreateIndex(t *TableMeta, idx IndexMeta) string {
	unique := ""
	if idx.Unique {
		unique = "UNIQUE "
	}
	cols := make([]string, len(idx.Columns))
	for i, c := range idx.Columns {
		cols[i] = quoteStd(c)
	}
	// oracle 无 IF NOT EXISTS，由编排层预检索引存在性后决定是否执行
	return fmt.Sprintf("CREATE %sINDEX %s ON %s (%s);",
		unique, quoteStd(idx.Name), quoteStd(t.Name), strings.Join(cols, ","))
}

// DropAllTables 清空全部表
func (r *oracleRenderer) DropAllTables(ctx context.Context, db gdb.DB) error {
	rows, err := db.Ctx(ctx).Raw("SELECT table_name FROM user_tables").All()
	if err != nil {
		return err
	}
	for _, row := range rows {
		if _, err = db.Exec(ctx, "DROP TABLE "+quoteStd(row["table_name"].String())+" CASCADE CONSTRAINTS PURGE"); err != nil {
			return err
		}
	}
	return nil
}

// FixSequence 修正自增序列
//
// Oracle 的 IDENTITY 由缓存失效自动续接，无需显式修序
func (r *oracleRenderer) FixSequence(ctx context.Context, db gdb.DB, t *TableMeta) error {
	for _, c := range t.Columns {
		if !c.AutoInc {
			continue
		}
		// oracle identity 由缓存失效自动续接，无需显式修序
	}
	return nil
}
