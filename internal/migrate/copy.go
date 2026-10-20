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
	"time"

	"github.com/gogf/gf/v2/database/gdb"
)

// ============================== 数据搬运 ==============================

// copyTable 分批搬运单表数据
//
// 源端按主键排序做快照分页，目标端使用多行 INSERT。
// onBatch 为每批完成回调（本批行数、累计行数），为 nil 时静默。
// 返回搬运行数；无主键表直接报错。
// SQL Server 目标使用内联字面量 INSERT 绕过参数上限，含自增列时以事务包裹切换 IDENTITY_INSERT。
func copyTable(ctx context.Context, src, dst gdb.DB, r Renderer, t *TableMeta, batchSize int, onBatch func(batch, total int64)) (int64, error) {
	pk := pkColumns(t)
	if len(pk) == 0 {
		return 0, errNoPK(t.Name)
	}
	pkCol := pk[0]
	if batchSize <= 0 {
		batchSize = 1000
	}

	// 目标端列清单（以源表列序组装 INSERT，map 键名驱动目标方言语句生成）
	colNames := make([]string, 0, len(t.Columns))
	for _, c := range t.Columns {
		colNames = append(colNames, c.Name)
	}

	// mssql：内联字面量 INSERT（无参数占位符，不受 2100 上限约束）；自增列事务钉连接
	if dst.GetConfig().Type == "mssql" {
		return copyMssql(ctx, src, dst, r, t, pkCol, colNames, batchSize, onBatch)
	}

	// 其余方言：参数化多行 INSERT
	return doCopy(ctx, src, dst.Ctx(ctx), t.Name, pkCol, colNames, batchSize, onBatch,
		func(ins inserter, batch []map[string]any) error {
			_, err := ins.Model(t.Name).Data(batch).Insert()
			return err
		})
}

// copyMssql SQL Server 数据搬运
//
// 使用内联字面量 INSERT 绕过参数上限；
// 自增列每批独立小事务切换 IDENTITY_INSERT，
// 既钉住连接，又避免整表单事务撑爆事务日志
func copyMssql(ctx context.Context, src, dst gdb.DB, r Renderer, t *TableMeta, pkCol string, colNames []string, batchSize int, onBatch func(batch, total int64)) (int64, error) {
	quoted := r.Quote(t.Name)

	// 自增列：每批 SET ON → INSERT → SET OFF → COMMIT（事务粒度=单批，日志占用受批大小约束）
	if hasAutoIncCol(t) {
		return doCopy(ctx, src, dst.Ctx(ctx), t.Name, pkCol, colNames, batchSize, onBatch,
			func(ins inserter, batch []map[string]any) error {
				tx, err := dst.Begin(ctx)
				if err != nil {
					return err
				}
				defer tx.Rollback()
				if _, err = tx.Exec("SET IDENTITY_INSERT " + quoted + " ON"); err != nil {
					return err
				}
				if _, err = tx.Exec(buildInsertSQL(r, t.Name, colNames, batch)); err != nil {
					return err
				}
				if _, err = tx.Exec("SET IDENTITY_INSERT " + quoted + " OFF"); err != nil {
					return err
				}
				return tx.Commit()
			})
	}

	return doCopy(ctx, src, dst.Ctx(ctx), t.Name, pkCol, colNames, batchSize, onBatch,
		func(ins inserter, batch []map[string]any) error {
			_, err := dst.Exec(ctx, buildInsertSQL(r, t.Name, colNames, batch))
			return err
		})
}

// inserter 目标端写入器
// gdb.DB 与 gdb.TX 均满足该接口
type inserter interface {
	Model(tableNameOrStruct ...any) *gdb.Model
}

// doCopy 核心搬运循环
//
// insertBatch 为每批写入实现，方言差异由此注入
func doCopy(ctx context.Context, src gdb.DB, ins inserter, table, pkCol string, colNames []string, batchSize int, onBatch func(batch, total int64), insertBatch func(ins inserter, batch []map[string]any) error) (int64, error) {
	var (
		total  int64
		lastPK any
	)
	for {
		m := src.Ctx(ctx).Model(table).Fields(strings.Join(quoteAll(src, colNames), ",")).Order(pkCol).Limit(batchSize)
		if lastPK != nil {
			m = m.Where(pkCol+" > ?", lastPK)
		}
		rows, err := m.All()
		if err != nil {
			return total, err
		}
		if len(rows) == 0 {
			return total, nil
		}
		batch := make([]map[string]any, 0, len(rows))
		for _, row := range rows {
			item := make(map[string]any, len(colNames))
			for _, c := range colNames {
				item[c] = convertValue(row[c])
			}
			batch = append(batch, item)
			lastPK = row[pkCol].Val()
		}
		if err = insertBatch(ins, batch); err != nil {
			return total, err
		}
		total += int64(len(batch))
		if onBatch != nil {
			onBatch(int64(len(batch)), total)
		}
		if len(rows) < batchSize {
			return total, nil
		}
	}
}

// buildInsertSQL 构建带内联字面量的多行 INSERT
//
// SQL Server 专用，绕过 2100 参数占位符上限
func buildInsertSQL(r Renderer, table string, colNames []string, batch []map[string]any) string {
	cols := make([]string, len(colNames))
	for i, c := range colNames {
		cols[i] = r.Quote(c)
	}
	valueRows := make([]string, 0, len(batch))
	for _, row := range batch {
		vals := make([]string, len(colNames))
		for i, c := range colNames {
			vals[i] = inlineValue(row[c])
		}
		valueRows = append(valueRows, "("+strings.Join(vals, ",")+")")
	}
	return "INSERT INTO " + r.Quote(table) + " (" + strings.Join(cols, ",") + ") VALUES " + strings.Join(valueRows, ",")
}

// inlineValue 将值转换为内联字面量
//
// nil 转 NULL，字符串与字节切片转义单引号，
// 数值与布尔原样，时间格式化输出
func inlineValue(v any) string {
	if v == nil {
		return "NULL"
	}
	switch t := v.(type) {
	case string:
		return "'" + strings.ReplaceAll(t, "'", "''") + "'"
	case []byte:
		if len(t) == 0 {
			return "NULL"
		}
		return "'" + strings.ReplaceAll(string(t), "'", "''") + "'"
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
		return fmt.Sprintf("%v", t)
	case bool:
		if t {
			return "1"
		}
		return "0"
	case time.Time:
		return "'" + t.Format("2006-01-02 15:04:05") + "'"
	default:
		return "'" + strings.ReplaceAll(fmt.Sprintf("%v", t), "'", "''") + "'"
	}
}

// quoteAll 源端标识符引用
//
// 源方言与目标方言可能不同，按源方言的引号规则处理
func quoteAll(db gdb.DB, cols []string) []string {
	out := make([]string, len(cols))
	for i, c := range cols {
		switch db.GetConfig().Type {
		case "mysql":
			out[i] = "`" + c + "`"
		case "mssql":
			out[i] = "[" + c + "]"
		default:
			out[i] = `"` + c + `"`
		}
	}
	return out
}

// convertValue 值转换
//
// 字节切片文本转为字符串，便于目标驱动绑定，其余原样返回
func convertValue(v gdb.Value) any {
	val := v.Val()
	if b, ok := val.([]byte); ok {
		return string(b)
	}
	return val
}

// errNoPK 无主键表错误
func errNoPK(table string) error {
	return &migrateError{"表 " + table + " 无主键，无法分页搬运（请先补主键）"}
}

// migrateError 迁移错误类型
type migrateError struct{ msg string }

// Error 返回错误信息
func (e *migrateError) Error() string { return e.msg }
