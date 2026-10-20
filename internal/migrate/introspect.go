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
	"sort"
	"strings"

	"github.com/gogf/gf/v2/database/gdb"
)

// ============================== 源库结构内省 ==============================

// listTables 查询源库表清单
//
// 受白名单过滤，按表名排序以保证结果可复现
func listTables(ctx context.Context, db gdb.DB, whitelist []string) ([]string, error) {
	tables, err := db.Tables(ctx)
	if err != nil {
		return nil, err
	}
	if len(whitelist) > 0 {
		allow := make(map[string]bool, len(whitelist))
		for _, t := range whitelist {
			allow[strings.TrimSpace(t)] = true
		}
		filtered := make([]string, 0, len(tables))
		for _, t := range tables {
			if allow[t] {
				filtered = append(filtered, t)
			}
		}
		tables = filtered
	}
	sort.Strings(tables)
	return tables, nil
}

// introspectTable 单表全量内省
//
// 采集列信息（含注释、默认值、自增、主键）、二级索引与表注释
func introspectTable(ctx context.Context, db gdb.DB, dialect, table string) (*TableMeta, error) {
	fields, err := db.TableFields(ctx, table)
	if err != nil {
		return nil, err
	}
	meta := &TableMeta{Name: table}
	// TableFields 返回 map，按 Index 排序还原物理列序
	sorted := make([]*gdb.TableField, 0, len(fields))
	for _, f := range fields {
		sorted = append(sorted, f)
	}
	sort.Slice(sorted, func(i, j int) bool { return sorted[i].Index < sorted[j].Index })
	for _, f := range sorted {
		typ, length, scale, enumVals := normalizeType(dialect, f.Type)
		meta.Columns = append(meta.Columns, ColumnMeta{
			Name:     f.Name,
			Type:     typ,
			Length:   length,
			Scale:    scale,
			EnumVals: enumVals,
			Nullable: f.Null,
			Default:  cleanDefault(dialect, f.Default),
			AutoInc:  isAutoInc(f.Extra),
			Comment:  f.Comment,
			IsPK:     strings.EqualFold(f.Key, "PRI") || strings.EqualFold(f.Key, "pri"),
		})
	}
	// 索引与表注释（方言相关）
	switch dialect {
	case "mysql":
		if err = introspectIndexesMySQL(ctx, db, table, meta); err != nil {
			return nil, err
		}
		meta.Comment = tableCommentMySQL(ctx, db, table)
	case "pgsql":
		if err = introspectIndexesPgsql(ctx, db, table, meta); err != nil {
			return nil, err
		}
		meta.Comment = tableCommentPgsql(ctx, db, table)
	default:
		// 其余方言：主键经列元数据已识别，二级索引降级跳过（首期实测链路 mysql→pgsql）
	}
	// 统一索引命名（迁移后全库一种规则，见 normalizeIndexName）
	normalizeIndexes(meta)
	return meta, nil
}

// normalizeIndexes 归一化二级索引名
//
// 主键随建表落地，不包含在索引集合中
func normalizeIndexes(meta *TableMeta) {
	for i := range meta.Indexes {
		idx := &meta.Indexes[i]
		idx.Name = normalizeIndexName(meta.Name, idx.Unique, idx.Columns)
	}
}

// normalizeIndexName 统一索引命名
//
// 唯一索引使用 uk_ 前缀，普通索引使用 idx_ 前缀，后接表名与列名；
// 名称超过 63 字符时截断，兼容多数数据库的索引名长度上限
func normalizeIndexName(table string, unique bool, cols []string) string {
	prefix := "idx"
	if unique {
		prefix = "uk"
	}
	name := prefix + "_" + table + "_" + strings.Join(cols, "_")
	if len(name) > 63 {
		name = name[:63]
	}
	return name
}

// isAutoInc 判断是否自增列
//
// 兼容 MySQL 的 auto_increment 以及 PostgreSQL、SQL Server 的 identity
func isAutoInc(extra string) bool {
	e := strings.ToLower(extra)
	return strings.Contains(e, "auto_increment") || strings.Contains(e, "identity")
}

// introspectIndexesMySQL 采集 MySQL 索引
//
// 主键写入列元数据，其余二级索引写入索引集合
func introspectIndexesMySQL(ctx context.Context, db gdb.DB, table string, meta *TableMeta) error {
	rows, err := db.Ctx(ctx).Raw("SHOW INDEX FROM `" + table + "`").All()
	if err != nil {
		return err
	}
	type idxHolder struct {
		unique bool
		cols   []string
	}
	holders := map[string]*idxHolder{}
	order := make([]string, 0)
	for _, r := range rows {
		name := r["Key_name"].String()
		h, ok := holders[name]
		if !ok {
			h = &idxHolder{unique: r["Non_unique"].Int() == 0}
			holders[name] = h
			order = append(order, name)
		}
		h.cols = append(h.cols, r["Column_name"].String())
	}
	for _, name := range order {
		h := holders[name]
		if name == "PRIMARY" {
			for i := range meta.Columns {
				for _, c := range h.cols {
					if meta.Columns[i].Name == c {
						meta.Columns[i].IsPK = true
					}
				}
			}
			continue
		}
		meta.Indexes = append(meta.Indexes, IndexMeta{Unique: h.unique, Name: name, Columns: h.cols})
	}
	return nil
}

// tableCommentMySQL 查询 MySQL 表注释
func tableCommentMySQL(ctx context.Context, db gdb.DB, table string) string {
	row, err := db.Ctx(ctx).Raw(
		"SELECT TABLE_COMMENT FROM information_schema.TABLES WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = ?", table,
	).One()
	if err != nil || len(row) == 0 {
		return ""
	}
	return row["TABLE_COMMENT"].String()
}

// introspectIndexesPgsql 采集 PostgreSQL 索引
//
// 主键写入列元数据，其余二级索引写入索引集合
func introspectIndexesPgsql(ctx context.Context, db gdb.DB, table string, meta *TableMeta) error {
	rows, err := db.Ctx(ctx).Raw(
		"SELECT indexname, indexdef FROM pg_indexes WHERE schemaname = current_schema() AND tablename = ?", table,
	).All()
	if err != nil {
		return err
	}
	for _, r := range rows {
		name := r["indexname"].String()
		def := r["indexdef"].String()
		cols := parsePgIndexCols(def)
		if len(cols) == 0 {
			continue
		}
		if name == table+"_pkey" {
			for i := range meta.Columns {
				for _, c := range cols {
					if meta.Columns[i].Name == c {
						meta.Columns[i].IsPK = true
					}
				}
			}
			continue
		}
		meta.Indexes = append(meta.Indexes, IndexMeta{
			Unique:  strings.Contains(strings.ToUpper(def), "UNIQUE"),
			Name:    name,
			Columns: cols,
		})
	}
	return nil
}

// parsePgIndexCols 从建索引语句中提取列清单
//
// 取括号内内容按逗号切分，跳过表达式部分
func parsePgIndexCols(def string) []string {
	start := strings.Index(def, "(")
	end := strings.LastIndex(def, ")")
	if start < 0 || end <= start {
		return nil
	}
	parts := strings.Split(def[start+1:end], ",")
	cols := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		// 表达式索引（如 lower(name)）取首标识符，跳过复杂表达式
		if i := strings.IndexAny(p, " ("); i > 0 {
			p = p[:i]
		}
		p = strings.Trim(p, `"`)
		if p != "" {
			cols = append(cols, p)
		}
	}
	return cols
}

// tableCommentPgsql 查询 PostgreSQL 表注释
func tableCommentPgsql(ctx context.Context, db gdb.DB, table string) string {
	row, err := db.Ctx(ctx).Raw(
		"SELECT obj_description(c.oid) AS comment FROM pg_class c JOIN pg_namespace n ON n.oid = c.relnamespace WHERE n.nspname = current_schema() AND c.relname = ?", table,
	).One()
	if err != nil || len(row) == 0 {
		return ""
	}
	return row["comment"].String()
}
