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

// ============================== 通知公告业务实现 ==============================

// sNotice 通知公告业务实现
// 实现 service.INotice 接口，接口定义见 api/v1/notice.go
type sNotice struct{}

func init() {
	service.RegisterNotice(&sNotice{})
}

// noticeSerializeMaps 枚举字段到字典编码的映射
var noticeSerializeMaps = map[string]string{
	"type":   "notice_type",
	"status": "notice_status",
}

// noticeSpec 通知公告模块声明
var noticeSpec = ModuleSpec[v1.NoticeItem, v1.NoticeListItem, v1.NoticeDetailRes]{
	Model:         dao.Notice.Ctx,
	LikeFields:    []string{"title"},
	EqFields:      []string{"type", "status"},
	OrderBy:       "id DESC",
	SerializeMaps: noticeSerializeMaps,
	FileFields:    []string{"cover"},   // 保存 temp→正式；读侧补域名
	ContentFields: []string{"content"}, // 富文本：存 [IMG_URL] 占位，读还原
	FileSaveDir:   "notice",
	ToData:        noticeToData,
	// 详情读取后点击数+1（读后更新，返回原值）
	AfterDetail: noticeCountClick,
}

// noticeToData 全量表单语义转换
// 兼容新增与修改请求，将入参转换为落库字段
func noticeToData(req any) map[string]any {
	switch r := req.(type) {
	case *v1.NoticeAddReq:
		return map[string]any{"title": r.Title, "cover": r.Cover, "type": r.Type,
			"status": r.Status, "content": r.Content}
	case *v1.NoticeUpdateReq:
		return map[string]any{"title": r.Title, "cover": r.Cover, "type": r.Type,
			"status": r.Status, "content": r.Content}
	}
	panic("noticeToData: 不支持的入参类型")
}

// noticeCountClick 详情点击数自增
//
// 先读取原值再更新，响应保持查询时的数值
func noticeCountClick(ctx context.Context, item *v1.NoticeDetailRes, req any) error {
	one, err := dao.Notice.Ctx(ctx).Where("id", item.Id).Fields("click_num").One()
	if err != nil || len(one) == 0 {
		return err
	}
	_, err = dao.Notice.Ctx(ctx).Where("id", item.Id).Data(map[string]any{
		"click_num": one["click_num"].Int() + 1,
	}).Update()
	return err
}

// Page 通知公告分页
// 按查询条件分页返回通知公告列表
func (s *sNotice) Page(ctx context.Context, req *v1.NoticePageReq) (res *v1.NoticePageRes, err error) {
	list, total, err := noticeSpec.Page(ctx, req, req.GetPageNo(), req.GetPageSize())
	if err != nil {
		return nil, err
	}
	return v1.NewPageData(req.GetPageNo(), req.GetPageSize(), total, derefSlice(list)), nil
}

// List 通知公告全量列表
// 按关键词、类型、状态等条件返回通知公告列表
func (s *sNotice) List(ctx context.Context, req *v1.NoticeListReq) (res v1.NoticeListRes, err error) {
	list, err := noticeSpec.List(ctx, req)
	if err != nil {
		return nil, err
	}
	return derefSlice(list), nil
}

// Detail 通知公告详情
// 按公告ID查询详情，并触发点击数自增
func (s *sNotice) Detail(ctx context.Context, req *v1.NoticeDetailReq) (res *v1.NoticeDetailRes, err error) {
	return noticeSpec.Detail(ctx, req.Id, req)
}

// Add 新增通知公告
// 新增成功后返回新记录ID
func (s *sNotice) Add(ctx context.Context, req *v1.NoticeAddReq) (res *v1.IdRes, err error) {
	id, err := noticeSpec.Add(ctx, req)
	if err != nil {
		return nil, err
	}
	return &v1.IdRes{Id: id}, nil
}

// Update 修改通知公告
// 按公告ID更新记录
func (s *sNotice) Update(ctx context.Context, req *v1.NoticeUpdateReq) (res *v1.NoticeUpdateRes, err error) {
	if err = noticeSpec.Update(ctx, req.Id, req); err != nil {
		return nil, err
	}
	return &v1.NoticeUpdateRes{}, nil
}

// Status 通知公告状态
// 非标准方法，手写实现，设置公告的草稿或发布状态
func (s *sNotice) Status(ctx context.Context, req *v1.NoticeStatusReq) (res *v1.NoticeStatusRes, err error) {
	if _, err = dao.Notice.Ctx(ctx).Where("id", req.Id).Data(map[string]any{"status": req.Status}).Update(); err != nil {
		return nil, err
	}
	return &v1.NoticeStatusRes{}, nil
}

// Delete 删除通知公告
// 按公告ID删除单条记录
func (s *sNotice) Delete(ctx context.Context, req *v1.NoticeDeleteReq) (res *v1.NoticeDeleteRes, err error) {
	if err = noticeSpec.Delete(ctx, req.Id); err != nil {
		return nil, err
	}
	return &v1.NoticeDeleteRes{}, nil
}

// BatchDelete 批量删除通知公告
// 按公告ID列表批量删除记录
func (s *sNotice) BatchDelete(ctx context.Context, req *v1.NoticeBatchDeleteReq) (res *v1.NoticeBatchDeleteRes, err error) {
	if err = noticeSpec.BatchDelete(ctx, req.Ids); err != nil {
		return nil, err
	}
	return &v1.NoticeBatchDeleteRes{}, nil
}
