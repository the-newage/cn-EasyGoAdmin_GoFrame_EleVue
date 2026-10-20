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

package hooks

import (
	"context"
	"database/sql"
	"reflect"
	"strings"
	"time"

	"github.com/gogf/gf/v2/database/gdb"

	"easygoadmin/internal/consts"
)

// ============================== 全局模型钩子 ==============================
//
// 全局钩子是架构红线的执行点：
//   - Select、Update、Delete 的条件过滤在 dao.Ctx 覆写中注入（见 ApplyCtx，进入 SQL 构建前最可靠）
//   - Insert：主键交给数据库自增，自动补齐审计字段（create_user/create_time/del_flag）与租户ID
//   - Update：自动为审计字段（update_user/update_time）盖章
//   - Delete：物理删除一律改写为软删除 UPDATE del_flag=1
//
// 所有表均含统一字段 id/create_user/create_time/update_user/update_time/del_flag。

// tenantTables 含 tenant_id 列的租户级表
//
// 与初始化脚本保持一致，共 9 张
var tenantTables = map[string]bool{
	"goframe_user":     true,
	"goframe_role":     true,
	"goframe_dept":     true,
	"goframe_level":    true,
	"goframe_position": true,
	"goframe_article":  true,
	"goframe_category": true,
	"goframe_notice":   true,
	"goframe_link":     true,
}

// IsTenantTable 判断是否为租户级表
func IsTenantTable(table string) bool { return tenantTables[table] }

// Handler 返回全局钩子
// 由 dao 扩展模板挂载
func Handler() gdb.HookHandler {
	return gdb.HookHandler{
		Insert: insertHook,
		Update: updateHook,
		Delete: deleteHook,
	}
}

// ApplyCtx 注入软删除过滤与租户隔离条件
//
// 由 dao 扩展模板的 Ctx 覆写方法调用，在模型层（SQL 构建之前）注入：
//
//	所有表：del_flag = 0
//	租户级表：平台级（超管）不过滤，全量可见；租户用户追加 tenant_id IN (0, 自身租户)
func ApplyCtx(ctx context.Context, m *gdb.Model, table string) *gdb.Model {
	m = m.Where("del_flag", 0)
	if tenantTables[table] && !consts.IsPlatform(ctx) {
		m = m.WhereIn("tenant_id", []int{0, consts.CtxTenantID(ctx)})
	}
	return m
}

// auditUser 获取审计操作人
//
// 优先取用户名，保证租户内稳定可追溯
func auditUser(ctx context.Context) string {
	if u := consts.GetUsername(ctx); u != "" {
		return u
	}
	return consts.GetRealname(ctx)
}

// isBlank 判断插入值是否为空
//
// 为空时由钩子补齐
func isBlank(v any) bool {
	if v == nil {
		return true
	}
	if s, ok := v.(string); ok {
		return s == ""
	}
	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64:
		return rv.IsZero()
	}
	return false
}

// hookTable 归一化钩子回调中的表名
//
// gdb 传入的表名带反引号，如 `goframe_user`
func hookTable(table string) string {
	return strings.Trim(table, "`")
}

// insertHook 插入钩子
//
// 自动补齐主键、审计字段与租户ID
func insertHook(ctx context.Context, in *gdb.HookInsertInput) (sql.Result, error) {
	now := time.Now()
	table := hookTable(in.Table)
	for _, item := range in.Data {
		// 主键交由数据库自增（id 统一 int，不使用雪花算法）
		if isBlank(item["del_flag"]) {
			item["del_flag"] = 0
		}
		if isBlank(item["create_time"]) {
			item["create_time"] = now
		}
		if isBlank(item["create_user"]) {
			item["create_user"] = auditUser(ctx)
		}
		// 租户级表：
		//   - 租户用户未指定租户ID时自动携带自身（零值视为未指定，防跨租户逃逸）
		//   - 超管未指定租户ID时代入默认租户（tenant.defaultId）；
		//     显式传 0 可创建平台共享行，显式传具体租户从显式值
		if tenantTables[table] {
			v, ok := item["tenant_id"]
			if !consts.IsPlatform(ctx) {
				if !ok || isBlank(v) {
					item["tenant_id"] = consts.CtxTenantID(ctx)
				}
			} else if !ok || v == nil {
				item["tenant_id"] = consts.DefaultTenantID(ctx)
			}
		}
	}
	return in.Next(ctx)
}

// updateHook 更新钩子
//
// 自动为审计字段盖章
func updateHook(ctx context.Context, in *gdb.HookUpdateInput) (sql.Result, error) {
	if data, ok := in.Data.(map[string]any); ok {
		data["update_time"] = time.Now()
		data["update_user"] = auditUser(ctx)
	}
	return in.Next(ctx)
}

// deleteHook 删除钩子
//
// 将物理删除改写为软删除，为架构红线；
// 注意：gdb 的 Model.Update 不接收 ctx，模型已携带上下文设置，
// 传入 ctx 会被误当作数据导致匹配字段失败
func deleteHook(ctx context.Context, in *gdb.HookDeleteInput) (sql.Result, error) {
	return in.Model.Update(map[string]any{"del_flag": 1})
}
