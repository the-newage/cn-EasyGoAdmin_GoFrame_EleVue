// Command migrator 整库迁移命令行
//
// 用法（仓库根目录执行）：
//
//	go run ./cmd/migrator \
//	  -src-type=mysql -src-host=127.0.0.1 -src-hostport=3306 \
//	  -src-database=easygoadmin.goframe.elevue -src-username=root -src-password=root \
//	  -dst-type=pgsql -dst-host=127.0.0.1 -dst-hostport=5432 \
//	  -dst-database=easygoadmin.goframe.elevue -dst-username=postgres -dst-password=123456 \
//	  -drop-target-first
//
// 可选：-tables "goframe_user,goframe_level"（只迁指定表）、-batch-size=1000（数据批大小）
package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/gogf/gf/v2/os/gctx"

	_ "github.com/gogf/gf/contrib/drivers/mssql/v2"
	_ "github.com/gogf/gf/contrib/drivers/mysql/v2"
	_ "github.com/gogf/gf/contrib/drivers/oracle/v2"
	_ "github.com/gogf/gf/contrib/drivers/pgsql/v2"
	_ "github.com/gogf/gf/contrib/drivers/sqlite/v2"

	"easygoadmin/internal/migrate"
)

func main() {
	// 源库连接要素
	srcType := flag.String("src-type", "mysql", "源库类型：mysql/pgsql/sqlite/mssql/oracle")
	srcHost := flag.String("src-host", "127.0.0.1", "源库地址")
	srcPort := flag.String("src-hostport", "3306", "源库端口")
	srcDB := flag.String("src-database", "", "源库库名（sqlite 填文件路径）")
	srcUser := flag.String("src-username", "root", "源库用户名")
	srcPass := flag.String("src-password", "", "源库密码")
	// 目标库连接要素
	dstType := flag.String("dst-type", "pgsql", "目标库类型：mysql/pgsql/sqlite/mssql/oracle")
	dstHost := flag.String("dst-host", "127.0.0.1", "目标库地址")
	dstPort := flag.String("dst-hostport", "5432", "目标库端口")
	dstDB := flag.String("dst-database", "", "目标库库名（sqlite 填文件路径）")
	dstUser := flag.String("dst-username", "postgres", "目标库用户名")
	dstPass := flag.String("dst-password", "", "目标库密码")
	// 迁移行为
	dropFirst := flag.Bool("drop-target-first", false, "迁移前清空目标库全部表")
	tables := flag.String("tables", "", "只迁指定表（逗号分隔，留空迁全部）")
	batchSize := flag.Int("batch-size", 1000, "数据搬运批大小")
	flag.Parse()

	if *srcDB == "" || *dstDB == "" {
		flag.Usage()
		fmt.Fprintln(os.Stderr, "\n示例: go run ./cmd/migrator -src-database=easygoadmin -dst-database=easygoadmin -drop-target-first")
		os.Exit(2)
	}

	ctx := gctx.New()
	opt := migrate.Options{
		Src: migrate.DSN{
			Type: *srcType, Host: *srcHost, Port: *srcPort,
			Database: *srcDB, Username: *srcUser, Password: *srcPass,
		},
		Dst: migrate.DSN{
			Type: *dstType, Host: *dstHost, Port: *dstPort,
			Database: *dstDB, Username: *dstUser, Password: *dstPass,
		},
		DropTargetFirst: *dropFirst,
		Tables:          splitList(*tables),
		BatchSize:       *batchSize,
	}

	_, err := migrate.Run(ctx, opt, func(format string, args ...any) {
		fmt.Printf(format, args...)
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "\n迁移失败: %v\n", err)
		os.Exit(1)
	}
}

// splitList 逗号分隔列表拆分（去空白）
func splitList(raw string) []string {
	if strings.TrimSpace(raw) == "" {
		return nil
	}
	out := make([]string, 0)
	for _, item := range strings.Split(raw, ",") {
		if s := strings.TrimSpace(item); s != "" {
			out = append(out, s)
		}
	}
	return out
}
