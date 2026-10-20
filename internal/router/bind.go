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

package router

import (
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/net/goai"

	"easygoadmin/internal/consts"
	"easygoadmin/internal/router/meta"
	"easygoadmin/utility/response"
)

// ============================== 控制器反射绑定 ==============================

// BodyUnmarshaler 自定义请求体解析逃生舱
//
// 当请求体顶层形态无法用结构体标签表达时实现该接口；
// 实现后跳过默认的请求参数解析，gvalid 校验需在实现内自行完成
type BodyUnmarshaler interface {
	UnmarshalBody(r *ghttp.Request) error
}

// Bind 反射注册控制器
//
// 基于 gmeta 扩展标签完成路由注册与统一包装：
//
//  1. 扫描控制器导出方法，校验签名：
//     func (c *Xxx) Page(ctx context.Context, req *v1.XxxReq) (*v1.XxxRes, error)
//  2. 解析请求结构体内嵌 gmeta.Meta 上的扩展标签：
//     method、path、summary、permission、log、logType、public
//  3. 按 method:path 注册路由并写入路由元数据表
//  4. 统一包装：权限校验 → 参数解析（含 gvalid 校验）→ 调用 → 信封写出
func Bind(s *ghttp.Server, ctrl any) {
	ctrlVal := reflect.ValueOf(ctrl)
	ctrlType := reflect.TypeOf(ctrl)
	for i := 0; i < ctrlType.NumMethod(); i++ {
		m := ctrlType.Method(i)
		// 目标签名：(ctx, *Req) (*Res, error)；In(0) 为接收者
		if m.Type.NumIn() != 3 || m.Type.NumOut() != 2 {
			continue
		}
		reqType := m.Type.In(2)
		if reqType.Kind() != reflect.Ptr {
			continue
		}
		metaField, ok := reqType.Elem().FieldByName("Meta")
		if !ok {
			continue
		}
		tag := metaField.Tag
		method := strings.ToLower(tag.Get("method"))
		path := tag.Get("path")
		if method == "" || path == "" {
			continue
		}
		rm := meta.RouteMeta{
			Method:     strings.ToUpper(method),
			Path:       path,
			Summary:    tag.Get("summary"),
			Permission: tag.Get("permission"),
			Log:        tag.Get("log"),
			LogType:    tag.Get("logType"),
			Public:     tag.Get("public") == "true",
		}
		meta.Register(rm)
		// 补注册进 OpenAPI 文档：闭包路由非严格路由，gf 不会自动收集，
		// 此处将带类型的控制器方法（签名 func(ctx,*Req)(*Res,error)）注册进 goai，
		// gmeta 的 path/method/summary 标签由 goai 解析为文档条目
		if err := s.GetOpenApi().Add(goai.AddInput{Object: ctrlVal.Method(i).Interface()}); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "openapi 注册失败 %s: %v\n", path, err)
		}
		handler := makeHandler(ctrlVal.Method(i), reqType, rm)
		s.BindHandler(method+":"+path, handler)
	}
}

// makeHandler 构造统一请求处理包装
//
// 依次执行权限校验、参数解析、业务调用与响应写出
func makeHandler(method reflect.Value, reqType reflect.Type, rm meta.RouteMeta) func(r *ghttp.Request) {
	return func(r *ghttp.Request) {
		// 路由元数据入上下文（操作日志读取）
		r.SetCtxVar(meta.CtxKey, rm)

		// 权限校验（与日志元数据同源；permission 为空 = 登录即可）
		if rm.Permission != "" && !consts.HasPerm(r.Context(), rm.Permission) {
			response.Write(r, response.Fail(gerror.NewCode(consts.CodeForbidden)))
			return
		}

		// 参数解析：路由参数 + query + body，并执行 gvalid 校验
		// 逃生舱：请求对象实现 BodyUnmarshaler 时走自定义解析（body 形态不规则的接口，如平铺 map）
		req := reflect.New(reqType.Elem())
		if u, ok := req.Interface().(BodyUnmarshaler); ok {
			if err := u.UnmarshalBody(r); err != nil {
				response.Write(r, response.Fail(err))
				return
			}
		} else if err := r.Parse(req.Interface()); err != nil {
			response.Write(r, response.Fail(err))
			return
		}

		// 调用业务方法
		out := method.Call([]reflect.Value{
			reflect.ValueOf(r.Context()),
			req,
		})
		var err error
		if !out[1].IsNil() {
			err, _ = out[1].Interface().(error)
		}
		if err != nil {
			response.Write(r, response.Fail(err))
			return
		}
		if out[0].IsNil() {
			response.Write(r, response.Ok(nil))
			return
		}
		response.Write(r, response.Ok(out[0].Interface()))
	}
}
