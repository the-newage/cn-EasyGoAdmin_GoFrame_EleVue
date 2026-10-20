package response

import (
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/net/ghttp"

	"easygoadmin/internal/consts"
)

// ============================== 统一响应 ==============================

// Res 统一响应结构
//
// 前端按此结构解析：
//   - 成功：code=0，前端仅在 code 为 0 时返回 data
//   - 失败：code=1，前端弹出 msg 提示
//   - 登录过期：code=401，前端跳转登录页
type Res struct {
	Code int    `json:"code"` // 业务状态码
	Ok   bool   `json:"ok"`   // code 为 0 时为 true
	Msg  string `json:"msg"`  // 提示信息
	Data any    `json:"data"` // 响应数据
}

// PageData 分页数据结构
//
// 对应前端分页列表契约
type PageData struct {
	List  any `json:"list"`  // 记录列表
	Total int `json:"total"` // 总记录数
}

// Ok 构造成功响应
//
// data 为 nil 时输出空对象，避免前端判空逻辑失败
func Ok(data any) *Res {
	if data == nil {
		data = map[string]any{}
	}
	return &Res{Code: 0, Msg: "操作成功", Data: data}
}

// OkMsg 构造带自定义提示语的成功响应
func OkMsg(msg string, data any) *Res {
	if data == nil {
		data = map[string]any{}
	}
	return &Res{Code: 0, Msg: msg, Data: data}
}

// Page 构造分页响应
func Page(list any, total int) *Res {
	return Ok(PageData{List: list, Total: total})
}

// Fail 构造错误响应
//
// 状态码 401 保持原样，其余错误码统一映射为 1，供前端弹出 msg
func Fail(err error) *Res {
	intCode := gerror.Code(err).Code()
	msg := err.Error()
	if msg == "" {
		msg = consts.CodeFailed.Message()
	}
	if intCode == 401 {
		return &Res{Code: 401, Msg: msg, Data: map[string]any{}}
	}
	return &Res{Code: 1, Msg: msg, Data: map[string]any{}}
}

// Write 输出 JSON 响应
//
// 响应码记入上下文，供操作日志判断正常或异常
func Write(r *ghttp.Request, res *Res) {
	WriteWithCode(r, res, res.Code)
}

// respCodeKey 与 respBodyKey 为请求上下文中的响应码与响应体键
//
// 供操作日志读取
var (
	respCodeKey = struct{ name string }{"respCode"}
	respBodyKey = struct{ name string }{"respBody"}
)

// WriteWithCode 输出 JSON 响应并记录指定状态码
func WriteWithCode(r *ghttp.Request, res *Res, code int) {
	res.Ok = code == 0
	r.SetCtxVar(respCodeKey, code)
	r.SetCtxVar(respBodyKey, res)
	r.Response.WriteJson(res)
}

// GetCode 读取响应码
//
// 第二个返回值表示是否已写出响应
func GetCode(r *ghttp.Request) (int, bool) {
	v := r.Context().Value(respCodeKey)
	if v == nil {
		return 0, false
	}
	code, _ := v.(int)
	return code, true
}

// GetBody 读取已写出的响应信封
//
// 供操作日志记录响应内容使用
func GetBody(r *ghttp.Request) (*Res, bool) {
	v := r.Context().Value(respBodyKey)
	if v == nil {
		return nil, false
	}
	res, ok := v.(*Res)
	return res, ok
}
