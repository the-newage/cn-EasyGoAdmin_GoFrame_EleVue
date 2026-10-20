// Command generator 代码生成命令行
//
// 用法（仓库根目录执行）：
//
//	go run ./cmd/generator -tables "goframe_order|订单表,goframe_sku|商品规格"   # 一键生成（文件+菜单+字典+路由）
//	go run ./cmd/generator -tables "goframe_order|订单表" -preview               # 仅预览渲染结果
//	go run ./cmd/generator -tables "goframe_order|订单表" -batch                 # 批量模式逐表统计
//
// 覆盖模式由配置 generator.overwrite 统一控制（true=生成即覆盖；false/缺省=缺失补生成、存在跳过）。
// 生成后需手动执行：go build ./...
package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"

	_ "github.com/gogf/gf/contrib/drivers/mssql/v2"
	_ "github.com/gogf/gf/contrib/drivers/mysql/v2"
	_ "github.com/gogf/gf/contrib/drivers/oracle/v2"
	_ "github.com/gogf/gf/contrib/drivers/pgsql/v2"
	_ "github.com/gogf/gf/contrib/drivers/sqlite/v2"

	"easygoadmin/api/v1"
	_ "easygoadmin/internal/logic" // 注册全部业务服务（init）
	"easygoadmin/internal/service"
)

// main 代码生成命令行入口
// 解析命令行参数后，按预览、批量或一键生成模式执行代码生成
func main() {
	tables := flag.String("tables", "", "表列表，格式: 表名|表注释,表名|表注释（与前端一键生成入参一致）")
	preview := flag.Bool("preview", false, "仅预览渲染结果，不写文件、不建菜单/字典/路由")
	batch := flag.Bool("batch", false, "批量模式（逐表统计输出）")
	flag.Parse()

	ctx := gctx.New()

	// 解析目标表列表，为空时打印用法并退出
	list := splitTables(*tables)
	if len(list) == 0 {
		flag.Usage()
		fmt.Println("\n示例: go run ./cmd/generator -tables \"goframe_order|订单表\"")
		os.Exit(2)
	}
	g.Log().Infof(ctx, "目标表 %d 张: %v", len(list), list)

	switch {
	// 预览模式：仅渲染并输出文件清单，不落盘
	case *preview:
		res, err := service.Generator().Preview(ctx, &v1.GeneratorPreviewReq{TableNames: list})
		if err != nil {
			fatal("预览失败: %v", err)
		}
		for _, f := range res.List {
			fmt.Printf("── %s (%d bytes)\n", f.Path, len(f.Content))
		}
		fmt.Printf("预览完成：共 %d 个文件\n", len(res.List))

	// 批量模式：逐表生成并输出各表统计结果
	case *batch:
		res, err := service.Generator().BatchGenerate(ctx, &v1.GeneratorBatchGenerateReq{TableNames: list})
		if err != nil {
			fatal("批量生成失败: %v", err)
		}
		for _, item := range res.List {
			mark := "OK"
			if !item.Success {
				mark = "FAIL"
			}
			fmt.Printf("[%s] %s — %s | 菜单:%s | 路由:%s\n",
				mark, item.Table, item.Message, item.Menu.Message, item.Route.Message)
			for _, d := range item.Dict {
				fmt.Printf("       字典 %s: %s\n", d.Dict, d.Status)
			}
		}

	// 默认模式：一键生成文件、菜单、字典与路由
	default:
		res, err := service.Generator().Generate(ctx, &v1.GeneratorGenerateReq{TableNames: list})
		if err != nil {
			fatal("生成失败: %v", err)
		}
		for _, p := range res.Generated {
			fmt.Println("  + 生成", p)
		}
		for _, p := range res.Skipped {
			fmt.Println("  - 跳过", p)
		}
		for _, p := range res.Overwritten {
			fmt.Println("  * 覆盖", p)
		}
		for _, p := range res.Failed {
			fmt.Println("  ! 失败", p)
		}
		fmt.Println("  菜单:", res.Menu.Message)
		for _, d := range res.Dict {
			fmt.Printf("  字典 %s: %s\n", d.Dict, d.Status)
		}
		fmt.Println("  路由:", res.Route.Message)
		fmt.Println("下一步: go build ./... 并重启服务")
	}
}

// fatal 输出错误信息并终止进程
func fatal(format string, args ...any) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}

// splitTables 解析表列表参数
// 按逗号分隔并去除空白项，返回有效的表列表
func splitTables(raw string) []string {
	out := []string{}
	for _, item := range strings.Split(raw, ",") {
		if s := strings.TrimSpace(item); s != "" {
			out = append(out, s)
		}
	}
	return out
}
