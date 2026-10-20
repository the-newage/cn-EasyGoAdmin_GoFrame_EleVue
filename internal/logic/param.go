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

package logic

import (
	"context"

	"easygoadmin/api/v1"
	"easygoadmin/internal/dao"
	"easygoadmin/internal/service"
)

// ============================== 系统参数业务实现 ==============================

// sParam 系统参数业务实现
// 实现 service.IParam 接口，接口定义见 api/v1/param.go
type sParam struct{}

func init() {
	service.RegisterParam(&sParam{})
}

// paramSerializeMaps 枚举字段到字典编码的映射
var paramSerializeMaps = map[string]string{
	"type":   "param_type",
	"status": "param_status",
}

// paramSpec 系统参数模块声明
var paramSpec = ModuleSpec[v1.ParamItem, v1.ParamItem, v1.ParamDetailRes]{
	Model:         dao.Param.Ctx,
	LikeFields:    []string{"name", "code"},
	EqFields:      []string{"type", "status"},
	OrderBy:       "sort, id",
	SerializeMaps: paramSerializeMaps,
	UniqueFields:  []UniqueRule{{Field: "code"}},     // 参数编码唯一
	FieldLabels:   map[string]string{"code": "参数编码"}, // 唯一提示显式声明（原取自列注释）
	ToData:        paramToData,
}

// paramToData 全量表单语义转换
// 兼容新增与修改请求，将入参转换为落库字段
func paramToData(req any) map[string]any {
	switch r := req.(type) {
	case *v1.ParamAddReq:
		return map[string]any{"name": r.Name, "code": r.Code, "value": r.Value,
			"type": r.Type, "status": r.Status, "sort": r.Sort, "note": r.Note}
	case *v1.ParamUpdateReq:
		return map[string]any{"name": r.Name, "code": r.Code, "value": r.Value,
			"type": r.Type, "status": r.Status, "sort": r.Sort, "note": r.Note}
	}
	panic("paramToData: 不支持的入参类型")
}

// Page 系统参数分页
// 按查询条件分页返回系统参数列表
func (s *sParam) Page(ctx context.Context, req *v1.ParamPageReq) (res *v1.ParamPageRes, err error) {
	list, total, err := paramSpec.Page(ctx, req, req.GetPageNo(), req.GetPageSize())
	if err != nil {
		return nil, err
	}
	return v1.NewPageData(req.GetPageNo(), req.GetPageSize(), total, derefSlice(list)), nil
}

// List 系统参数全量列表
// 按关键词、类型、状态等条件返回系统参数列表
func (s *sParam) List(ctx context.Context, req *v1.ParamListReq) (res v1.ParamListRes, err error) {
	list, err := paramSpec.List(ctx, req)
	if err != nil {
		return nil, err
	}
	return derefSlice(list), nil
}

// Detail 系统参数详情
// 按参数ID查询详情
func (s *sParam) Detail(ctx context.Context, req *v1.ParamDetailReq) (res *v1.ParamDetailRes, err error) {
	return paramSpec.Detail(ctx, req.Id, req)
}

// Add 新增系统参数
// 校验参数编码唯一后新增，返回新记录ID
func (s *sParam) Add(ctx context.Context, req *v1.ParamAddReq) (res *v1.IdRes, err error) {
	id, err := paramSpec.Add(ctx, req)
	if err != nil {
		return nil, err
	}
	return &v1.IdRes{Id: id}, nil
}

// Update 修改系统参数
// 按参数ID更新记录
func (s *sParam) Update(ctx context.Context, req *v1.ParamUpdateReq) (res *v1.ParamUpdateRes, err error) {
	if err = paramSpec.Update(ctx, req.Id, req); err != nil {
		return nil, err
	}
	return &v1.ParamUpdateRes{}, nil
}

// Delete 删除系统参数
// 按参数ID删除单条记录
func (s *sParam) Delete(ctx context.Context, req *v1.ParamDeleteReq) (res *v1.ParamDeleteRes, err error) {
	if err = paramSpec.Delete(ctx, req.Id); err != nil {
		return nil, err
	}
	return &v1.ParamDeleteRes{}, nil
}

// BatchDelete 批量删除系统参数
// 按参数ID列表批量删除记录
func (s *sParam) BatchDelete(ctx context.Context, req *v1.ParamBatchDeleteReq) (res *v1.ParamBatchDeleteRes, err error) {
	if err = paramSpec.BatchDelete(ctx, req.Ids); err != nil {
		return nil, err
	}
	return &v1.ParamBatchDeleteRes{}, nil
}

// paramGetString 按编码读取启用状态的系统参数值
//
// 软删除由全局钩子过滤，未命中时返回 fallback
func paramGetString(ctx context.Context, code, fallback string) string {
	v, err := dao.Param.Ctx(ctx).Where("code", code).Where("status", 1).Fields("value").Value()
	if err != nil || v.IsNil() || v.String() == "" {
		return fallback
	}
	return v.String()
}

// getDefaultPassword 获取默认密码明文
//
// 取系统参数 DEFAULT_PASSWORD，未配置时回退 123456；
// 仅用于初始化密码的中间环节，禁止写入日志
func getDefaultPassword(ctx context.Context) string {
	return paramGetString(ctx, "DEFAULT_PASSWORD", "123456")
}
