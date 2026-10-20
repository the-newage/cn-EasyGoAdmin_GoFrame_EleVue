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
	"time"

	"github.com/gogf/gf/v2/database/gdb"
)

// ============================== 整库迁移主流程 ==============================

// banner 迁移起止分隔线
const banner = "=============================================================="

// Run 整库迁移主流程
//
// 共五个阶段：建目标库 → 建表结构 → 逐表迁移数据 → 重置自增序列 → 逐表行数校验。
// progress 为进度输出回调，CLI 传入格式化输出函数，为 nil 时静默。
func Run(ctx context.Context, opt Options, progress func(format string, args ...any)) (*Report, error) {
	report := &Report{}
	logf := func(format string, args ...any) {
		if progress != nil {
			progress(format, args...)
		}
	}
	start := time.Now()

	src, err := openConn(ctx, "migrate_src", opt.Src)
	if err != nil {
		return nil, err
	}
	tables, err := listTables(ctx, src, opt.Tables)
	if err != nil {
		return nil, fmt.Errorf("读取源库表清单失败: %v", err)
	}
	if len(tables) == 0 {
		return nil, fmt.Errorf("源库无表（或白名单未命中任何表）")
	}
	report.TableTotal = len(tables)

	renderer, err := NewRenderer(opt.Dst.Type)
	if err != nil {
		return nil, err
	}

	// 抬头
	logf("%s\n", banner)
	logf("[迁移开始] 源: %s:%s → 目标: %s:%s\n", opt.Src.Type, opt.Src.Database, opt.Dst.Type, opt.Dst.Database)
	logf("[迁移开始] 待迁移 %d 张表\n", len(tables))
	logf("%s\n", banner)

	// ── 1/5 建目标库 ──
	phase := time.Now()
	logf("[·] [1/5] 建目标库...\n")
	if err = ensureDatabase(ctx, opt.Dst); err != nil {
		return nil, err
	}
	logf("  · 目标库 \"%s\" 已就绪\n", opt.Dst.Database)
	logf("[✓] [1/5] 建目标库 完成（耗时 %s）\n", dur(time.Since(phase)))

	dst, err := openConn(ctx, "migrate_dst", opt.Dst)
	if err != nil {
		return nil, err
	}

	// ── 2/5 建表结构 ──
	phase = time.Now()
	logf("[·] [2/5] 建表结构（%d 张表）...\n", len(tables))
	if opt.DropTargetFirst {
		if err = renderer.DropAllTables(ctx, dst); err != nil {
			return nil, fmt.Errorf("清空目标库失败: %v", err)
		}
	}
	metas := make([]*TableMeta, 0, len(tables))
	totalIndexes := 0
	for _, table := range tables {
		meta, e := introspectTable(ctx, src, opt.Src.Type, table)
		if e != nil {
			report.TableFailed++
			report.Failed = append(report.Failed, FailedTable{Table: table, Err: "内省失败: " + e.Error()})
			logf("  [✗ %s] 内省失败: %v\n", table, e)
			continue
		}
		// 非清空模式下目标已有同名表则该表跳过（建表用 IF NOT EXISTS，数据会追加，须前置拦截）
		if !opt.DropTargetFirst {
			if exists, e2 := tableExists(ctx, dst, table); e2 == nil && exists {
				report.TableFailed++
				report.Failed = append(report.Failed, FailedTable{Table: table, Err: "目标库已存在同名表"})
				logf("  [✗ %s] 目标库已存在同名表\n", table)
				continue
			}
		}
		if _, e = dst.Exec(ctx, renderer.CreateTable(meta)); e != nil {
			report.TableFailed++
			report.Failed = append(report.Failed, FailedTable{Table: table, Err: "建表失败: " + e.Error()})
			logf("  [✗ %s] 建表失败: %v\n", table, e)
			continue
		}
		for _, idx := range meta.Indexes {
			if _, e = dst.Exec(ctx, renderer.CreateIndex(meta, idx)); e != nil {
				report.TableFailed++
				report.Failed = append(report.Failed, FailedTable{Table: table, Err: "建索引 " + idx.Name + " 失败: " + e.Error()})
				logf("  [✗ %s] 建索引 %s 失败: %v\n", table, idx.Name, e)
				break
			}
		}
		totalIndexes += len(meta.Indexes)
		metas = append(metas, meta)
	}
	logf("[✓] [2/5] 建表结构 完成（耗时 %s），%d 张表，%d 个索引，0 个外键\n",
		dur(time.Since(phase)), len(metas), totalIndexes)

	// ── 3/5 逐表迁移数据 ──
	phase = time.Now()
	logf("[·] [3/5] 逐表迁移数据...\n")
	totalRows := int64(0)
	for i, meta := range metas {
		logf("  · [%d/%d] 正在迁移 %s...\n", i+1, len(metas), meta.Name)
		tableStart := time.Now()
		rows, e := copyTable(ctx, src, dst, renderer, meta, opt.BatchSize, func(batch, total int64) {
			logf("  [%s] +%d（累计 %d）\n", meta.Name, batch, total)
		})
		if e != nil {
			report.TableFailed++
			report.Failed = append(report.Failed, FailedTable{Table: meta.Name, Err: "数据搬运失败: " + e.Error()})
			logf("  [✗ %s] 失败: %v\n", meta.Name, e)
			continue
		}
		report.TableOK++
		totalRows += rows
		logf("  [✓ %s] 共 %d 行（耗时 %s）\n", meta.Name, rows, dur(time.Since(tableStart)))
	}
	report.RowTotal = totalRows
	logf("[✓] [3/5] 逐表迁移数据 完成（耗时 %s），共 %d 行\n", dur(time.Since(phase)), totalRows)

	// ── 4/5 重置自增序列 ──
	phase = time.Now()
	logf("[·] [4/5] 重置自增序列...\n")
	for _, meta := range metas {
		if e := renderer.FixSequence(ctx, dst, meta); e != nil {
			report.Failed = append(report.Failed, FailedTable{Table: meta.Name, Err: "序列修复失败: " + e.Error()})
		}
	}
	logf("[✓] [4/5] 重置自增序列 完成（耗时 %s）\n", dur(time.Since(phase)))

	// ── 5/5 逐表行数校验 ──
	phase = time.Now()
	logf("[·] [5/5] 逐表行数校验...\n")
	allPass := true
	for _, meta := range metas {
		srcCount, e1 := countRows(ctx, src, meta.Name)
		dstCount, e2 := countRows(ctx, dst, meta.Name)
		if e1 != nil || e2 != nil || srcCount != dstCount {
			allPass = false
			logf("  [%s] source=%d target=%d FAIL\n", meta.Name, srcCount, dstCount)
			continue
		}
		logf("  [%s] source=%d target=%d PASS\n", meta.Name, srcCount, dstCount)
	}
	if allPass {
		logf("[✓] [5/5] 逐表行数校验 完成（耗时 %s），%d 张表全部 PASS\n", dur(time.Since(phase)), len(metas))
	} else {
		logf("[✓] [5/5] 逐表行数校验 完成（耗时 %s），存在 FAIL\n", dur(time.Since(phase)))
	}

	// 落款
	logf("%s\n", banner)
	logf("[迁移完成] 成功：%d 行 / %d 张表 / 总耗时 %s\n", totalRows, report.TableOK, dur(time.Since(start)))
	logf("%s\n", banner)

	if report.TableFailed > 0 {
		return report, fmt.Errorf("%d 张表迁移失败", report.TableFailed)
	}
	if !allPass {
		return report, fmt.Errorf("存在行数校验不一致的表")
	}
	return report, nil
}

// countRows 统计表行数
func countRows(ctx context.Context, db gdb.DB, table string) (int64, error) {
	n, err := db.Ctx(ctx).Model(table).Count()
	return int64(n), err
}

// dur 时长格式化
//
// 不足 1 秒输出毫秒，不足 1 分输出秒，否则输出分秒
func dur(d time.Duration) string {
	if d < time.Second {
		return fmt.Sprintf("%dms", d.Milliseconds())
	}
	if d < time.Minute {
		return fmt.Sprintf("%ds", int(d.Seconds()))
	}
	return fmt.Sprintf("%dm%ds", int(d.Minutes()), int(d.Seconds())%60)
}
