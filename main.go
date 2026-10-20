package main

import (
	_ "easygoadmin/internal/logic"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"

	// 五种数据库驱动注册（按需切换 manifest/config/config.yaml 中 database.default.type/link）
	_ "github.com/gogf/gf/contrib/drivers/mssql/v2"
	_ "github.com/gogf/gf/contrib/drivers/mysql/v2"
	_ "github.com/gogf/gf/contrib/drivers/oracle/v2"
	_ "github.com/gogf/gf/contrib/drivers/pgsql/v2"
	_ "github.com/gogf/gf/contrib/drivers/sqlite/v2"

	// Redis 适配器注册（gredis/gcache 黑名单、验证码、字典缓存）
	_ "github.com/gogf/gf/contrib/nosql/redis/v2"

	"easygoadmin/internal/cmd"
)

// ============================== 服务入口 ==============================

// main 后端服务入口
//
// 完成依赖注册后交由 cmd.Run 启动服务
func main() {
	g.Log().Info(gctx.New(), "EasyGoAdmin 服务启动中...")
	cmd.Run(gctx.New())
}
