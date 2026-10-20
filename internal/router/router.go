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

	"github.com/gogf/gf/v2/net/ghttp"

	"easygoadmin/internal/controller"
	"easygoadmin/internal/middleware"
	"easygoadmin/utility/response"
)

// ============================== 路由注册 ==============================

// errorf 构造状态页错误
//
// 响应层将其映射为 code=1 并附上消息
func errorf(format string, args ...any) error {
	return fmt.Errorf(format, args...)
}

// Register 注册全局中间件链与路由
//
// 中间件顺序：CORS → Response（panic 兜底）→ RateLimit（限流）→ Demo（演示拦截）→ JWT → Tenant → OperationLog → 处理器。
// 权限校验在 Bind 包装层完成，与路由元数据同源；操作日志按元数据异步落库。
func Register(s *ghttp.Server) {
	s.Use(
		middleware.CORS,
		middleware.Response,
		middleware.RateLimit,
		middleware.Demo,
		middleware.JWT,
		middleware.Tenant,
		middleware.OperationLog,
	)

	// 非 200 状态统一输出 JSON（404/500 不返回 HTML）
	s.BindStatusHandler(404, func(r *ghttp.Request) {
		response.Write(r, response.Fail(errorf("接口不存在: %s %s", r.Method, r.URL.Path)))
	})
	s.BindStatusHandler(500, func(r *ghttp.Request) {
		response.Write(r, response.Fail(errorf("服务器内部错误")))
	})

	// 静态文件：将 resource 目录作为静态资源根目录对外访问
	s.SetServerRoot("resource")

	// 健康检查（免登录；公开路由不带 /v1，见 middleware/public.go）
	s.BindHandler("get:/api/health", func(r *ghttp.Request) {
		response.Write(r, response.Ok(map[string]any{"status": "up"}))
	})

	// API 文档页（本地托管 Redoc，静态资源由 SetServerRoot("resource") 提供）
	s.BindHandler("get:/swagger", func(r *ghttp.Request) {
		r.Response.RedirectTo("/swagger/index.html")
	})

	registerControllers(s)
}

// registerControllers 注册各模块控制器
//
// 随实施进度追加
func registerControllers(s *ghttp.Server) {
	// —— auth 认证域 ——
	Bind(s, new(controller.CAuth))

	// —— index 个人中心 ——
	Bind(s, new(controller.CIndex))

	// —— system 域 ——
	Bind(s, new(controller.CUser))
	Bind(s, new(controller.CRole))
	Bind(s, new(controller.CMenu))
	Bind(s, new(controller.CDept))
	Bind(s, new(controller.CLevel))
	Bind(s, new(controller.CPosition))
	Bind(s, new(controller.CTenant))

	// —— data 域 ——
	Bind(s, new(controller.CConfig))
	Bind(s, new(controller.CConfigItem))
	Bind(s, new(controller.CParam))
	Bind(s, new(controller.CDict))
	Bind(s, new(controller.CDictItem))
	Bind(s, new(controller.CNotice))
	Bind(s, new(controller.CCity))

	// —— content 域 ——
	Bind(s, new(controller.CArticle))
	Bind(s, new(controller.CCategory))
	Bind(s, new(controller.CLink))

	// —— monitor 域 ——
	Bind(s, new(controller.COperLog))
	Bind(s, new(controller.CLoginLog))
	Bind(s, new(controller.CJob))
	Bind(s, new(controller.CJobLog))
	Bind(s, new(controller.COnlineUser))
	Bind(s, new(controller.CServerMonitor))
	Bind(s, new(controller.CCacheMonitor))

	// —— tool 域 ——
	Bind(s, new(controller.CFileTemplate))
	Bind(s, new(controller.CGenerator))
	Bind(s, new(controller.CUpload))

	// —— tool 域（代码生成器注册示例） ——
	Bind(s, new(controller.CExample))
}
