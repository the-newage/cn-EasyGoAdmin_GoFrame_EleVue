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

	"github.com/gogf/gf/v2/database/gdb"

	"easygoadmin/api/v1"
	"easygoadmin/internal/dao"
	"easygoadmin/internal/service"
)

// ============================== 任务日志业务实现 ==============================

// sJobLog 任务日志业务实现
// 实现 service.IJobLog 接口，承载任务日志模块的业务逻辑
type sJobLog struct{}

// jobLogSerializeMaps 枚举字段到字典编码的映射
// {field}Text 由 dict_item 解析，字典管理页可修改文案
var jobLogSerializeMaps = map[string]string{
	"status": "job_log_status",
}

// jobLogSpec 任务日志模块声明
var jobLogSpec = ModuleSpec[v1.JobLogItem, v1.JobLogItem, v1.JobLogDetailRes]{
	Model:         dao.JobLog.Ctx,
	LikeFields:    []string{"job_name"},
	EqFields:      []string{"status", "job_id"},
	OrderBy:       "id DESC",
	SerializeMaps: jobLogSerializeMaps,
}

func init() {
	service.RegisterJobLog(&sJobLog{})
}

// Page 任务日志分页
// 按查询条件分页返回任务日志列表
func (s *sJobLog) Page(ctx context.Context, req *v1.JobLogPageReq) (res *v1.JobLogPageRes, err error) {
	list, total, err := jobLogSpec.Page(ctx, req, req.GetPageNo(), req.GetPageSize())
	if err != nil {
		return nil, err
	}
	return v1.NewPageData(req.GetPageNo(), req.GetPageSize(), total, derefSlice(list)), nil
}

// List 任务日志全量列表
// 按关键词、任务ID、状态等条件返回任务日志列表
func (s *sJobLog) List(ctx context.Context, req *v1.JobLogListReq) (res v1.JobLogListRes, err error) {
	list, err := jobLogSpec.List(ctx, req)
	if err != nil {
		return nil, err
	}
	return derefSlice(list), nil
}

// Detail 任务日志详情
// 按日志ID查询详情
func (s *sJobLog) Detail(ctx context.Context, req *v1.JobLogDetailReq) (res *v1.JobLogDetailRes, err error) {
	return jobLogSpec.Detail(ctx, req.Id, req)
}

// Delete 删除任务日志
// 按日志ID删除单条记录
func (s *sJobLog) Delete(ctx context.Context, req *v1.JobLogDeleteReq) (res *v1.JobLogDeleteRes, err error) {
	if err = jobLogSpec.Delete(ctx, req.Id); err != nil {
		return nil, err
	}
	return &v1.JobLogDeleteRes{}, nil
}

// BatchDelete 批量删除任务日志
// 按日志ID列表批量删除记录
func (s *sJobLog) BatchDelete(ctx context.Context, req *v1.JobLogBatchDeleteReq) (res *v1.JobLogBatchDeleteRes, err error) {
	if err = jobLogSpec.BatchDelete(ctx, req.Ids); err != nil {
		return nil, err
	}
	return &v1.JobLogBatchDeleteRes{}, nil
}

// toItem 数据库记录转换为任务日志行数据
// 时间字段按可空处理，枚举状态同步补全状态文本
func (s *sJobLog) toItem(ctx context.Context, r gdb.Record) v1.JobLogItem {
	item := v1.JobLogItem{
		Id:             r["id"].Int64(),
		JobId:          r["job_id"].Int64(),
		JobName:        r["job_name"].String(),
		JobGroup:       r["job_group"].String(),
		JobTrigger:     r["job_trigger"].String(),
		JobMessage:     r["job_message"].String(),
		CronExpression: r["cron_expression"].String(),
		Status:         r["status"].Int(),
		StatusText:     DictGetText(ctx, jobLogSerializeMaps["status"], r["status"].Int()),
		ConsumeTime:    r["consume_time"].Int64(),
		ExceptionInfo:  r["exception_info"].String(),
		CreateTime:     r["create_time"].GTime(),
	}
	if v := r["start_time"]; !v.IsNil() {
		item.StartTime = v.GTime()
	}
	if v := r["end_time"]; !v.IsNil() {
		item.EndTime = v.GTime()
	}
	return item
}
