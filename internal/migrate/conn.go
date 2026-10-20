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

	"github.com/gogf/gf/v2/database/gdb"
)

// ============================== 数据库连接 ==============================

// openConn 按 DSN 动态建立连接并探活
//
// 通过 gdb.New 注册独立命名节点，不触碰 config.yaml 中的默认连接
func openConn(ctx context.Context, node string, dsn DSN) (gdb.DB, error) {
	link, err := buildLink(dsn)
	if err != nil {
		return nil, err
	}
	db, err := gdb.New(gdb.ConfigNode{
		Name: node,
		Type: dsn.Type,
		Link: link,
	})
	if err != nil {
		return nil, fmt.Errorf("连接 %s 库失败: %v", node, err)
	}
	if _, err = db.Query(ctx, "SELECT 1"); err != nil {
		return nil, fmt.Errorf("%s 库探活失败: %v", node, err)
	}
	return db, nil
}

// ensureDatabase 目标库不存在时先创建
//
// MySQL、PostgreSQL、SQL Server 经维护库建库；Oracle 与 SQLite 无需预建
func ensureDatabase(ctx context.Context, dsn DSN) error {
	var maintainDB, createSQL, checkSQL string
	switch dsn.Type {
	case "mysql":
		maintainDB, checkSQL = "", "SELECT SCHEMA_NAME FROM information_schema.SCHEMATA WHERE SCHEMA_NAME = ?"
		createSQL = "CREATE DATABASE `" + dsn.Database + "` CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci"
	case "pgsql":
		maintainDB, checkSQL = "postgres", "SELECT 1 FROM pg_database WHERE datname = ?"
		createSQL = `CREATE DATABASE "` + dsn.Database + `"`
	case "mssql":
		maintainDB, checkSQL = "master", "SELECT 1 FROM sys.databases WHERE name = ?"
		createSQL = "CREATE DATABASE [" + dsn.Database + "]"
	default:
		return nil // oracle/sqlite 无需预建
	}
	checkDSN := dsn
	checkDSN.Database = maintainDB
	link, err := buildLink(checkDSN)
	if err != nil {
		return err
	}
	maintain, err := gdb.New(gdb.ConfigNode{Name: "migrate_maintain", Type: dsn.Type, Link: link})
	if err != nil {
		return fmt.Errorf("连接 %s 维护库失败: %v", dsn.Type, err)
	}
	defer maintain.Close(ctx)
	row, err := maintain.Query(ctx, checkSQL, dsn.Database)
	if err != nil {
		return fmt.Errorf("检查目标库是否存在失败: %v", err)
	}
	if len(row) > 0 {
		return nil // 已存在
	}
	if _, err = maintain.Exec(ctx, createSQL); err != nil {
		return fmt.Errorf("创建目标库 %s 失败: %v", dsn.Database, err)
	}
	return nil
}

// buildLink 将 DSN 转换为 gdb link 串
//
// 各方言按驱动约定格式组装，对齐配置示例
func buildLink(dsn DSN) (string, error) {
	switch dsn.Type {
	case "mysql":
		return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=true&loc=Local",
			dsn.Username, dsn.Password, dsn.Host, dsn.Port, dsn.Database), nil
	case "pgsql":
		return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?sslmode=disable",
			dsn.Username, dsn.Password, dsn.Host, dsn.Port, dsn.Database), nil
	case "mssql":
		return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?encrypt=disable",
			dsn.Username, dsn.Password, dsn.Host, dsn.Port, dsn.Database), nil
	case "oracle":
		return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
			dsn.Username, dsn.Password, dsn.Host, dsn.Port, dsn.Database), nil
	case "sqlite":
		if dsn.Database == "" {
			return "", fmt.Errorf("sqlite 需以 -src-database/-dst-database 指定数据库文件路径")
		}
		return fmt.Sprintf("sqlite::@file(%s)", dsn.Database), nil
	}
	return "", fmt.Errorf("不支持的数据库类型: %s", dsn.Type)
}

// tableExists 判断目标库中同名表是否已存在
func tableExists(ctx context.Context, db gdb.DB, table string) (bool, error) {
	tables, err := db.Tables(ctx)
	if err != nil {
		return false, err
	}
	for _, t := range tables {
		if t == table {
			return true, nil
		}
	}
	return false, nil
}
