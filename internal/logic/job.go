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

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/gcron"

	"easygoadmin/api/v1"
	"easygoadmin/internal/consts"
	"easygoadmin/internal/dao"
	"easygoadmin/internal/service"
)

// ============================== 定时任务业务实现 ==============================

// sJob 定时任务业务实现
// 实现 service.IJob 接口，接口定义见 api/v1/job.go
type sJob struct{}

// jobSerializeMaps 枚举字段到字典编码的映射
// {field}Text 由基座自动填充
var jobSerializeMaps = map[string]string{
	"status":         "job_status",
	"execute_policy": "job_execute_policy",
	"is_sync":        "job_is_sync",
}

func init() {
	service.RegisterJob(&sJob{})
	// 调度器启动加载（cmd 启动流程会调用 DaemonInit，此处保证 init 顺序无关）
}

// jobSpec 定时任务模块声明
// 调度器挂载与摘除通过钩子承载
var jobSpec = ModuleSpec[v1.JobItem, v1.JobItem, v1.JobDetailRes]{
	Model:         dao.Job.Ctx,
	LikeFields:    []string{"job_name", "job_alias"},
	EqFields:      []string{"status"},
	OrderBy:       "id DESC",
	ListOrderBy:   "id",
	SerializeMaps: jobSerializeMaps,
	ToData:        jobToData,
	BeforeAdd:     jobDefaultStatus, // 新增固定 status=0（未发布）
	BeforeUpdate:  jobExistGuard,    // 存在性检查
	AfterUpdate:   jobRemount,       // 运行中任务重挂新配置
	BeforeDelete:  jobDeleteGuard,   // 存在性 + 摘除调度
}

// jobToData 全量表单语义转换
// 新增与修改均提交全部字段，status 不在表单内
func jobToData(req any) map[string]any {
	switch r := req.(type) {
	case *v1.JobAddReq:
		return map[string]any{
			"job_name": r.JobName, "job_alias": r.JobAlias, "job_group": r.JobGroup,
			"job_trigger": r.JobTrigger, "cron_expression": r.CronExpression,
			"execute_policy": r.ExecutePolicy, "is_sync": r.IsSync,
			"url": r.Url, "params": r.Params, "note": r.Note,
		}
	case *v1.JobUpdateReq:
		return map[string]any{
			"job_name": r.JobName, "job_alias": r.JobAlias, "job_group": r.JobGroup,
			"job_trigger": r.JobTrigger, "cron_expression": r.CronExpression,
			"execute_policy": r.ExecutePolicy, "is_sync": r.IsSync,
			"url": r.Url, "params": r.Params, "note": r.Note,
		}
	}
	panic("jobToData: 不支持的入参类型")
}

// jobDefaultStatus 新增时固定状态为未发布
// 发布操作走 Status 或 Resume 方法
func jobDefaultStatus(ctx context.Context, data map[string]any, req any) error {
	data["status"] = 0
	return nil
}

// jobExistGuard 更新前存在性检查
func jobExistGuard(ctx context.Context, data map[string]any, id int64, req any) error {
	if n, err := dao.Job.Ctx(ctx).Where("id", id).Count(); err != nil {
		return err
	} else if n == 0 {
		return gerror.NewCode(consts.CodeRecordNotFound)
	}
	return nil
}

// jobRemount 更新后重挂调度
// 运行中任务先摘除再按新配置挂载，挂载失败时降级为暂停
func jobRemount(ctx context.Context, id int64, req any) error {
	r := mustAs[*v1.JobUpdateReq](req)
	one, err := dao.Job.Ctx(ctx).Where("id", id).One()
	if err != nil {
		return err
	}
	if len(one) == 0 || one["status"].Int() != 1 {
		return nil
	}
	removeJobFromScheduler(id)
	if err = addJobToScheduler(id, r.CronExpression, r.ExecutePolicy, r.IsSync); err != nil {
		// 挂载失败降级为暂停
		_, _ = dao.Job.Ctx(ctx).Where("id", id).Data(map[string]any{"status": 2}).Update()
		return gerror.NewCodef(consts.CodeParamError, "Cron表达式无效: %v", err)
	}
	return nil
}

// jobDeleteGuard 删除守卫
// 校验存在性并摘除调度
func jobDeleteGuard(ctx context.Context, id int64) error {
	if n, err := dao.Job.Ctx(ctx).Where("id", id).Count(); err != nil {
		return err
	} else if n == 0 {
		return gerror.NewCode(consts.CodeRecordNotFound)
	}
	removeJobFromScheduler(id)
	return nil
}

// ---------- 标准方法：3行委托 ----------

// Page 任务分页
// 按查询条件分页返回定时任务列表
func (s *sJob) Page(ctx context.Context, req *v1.JobPageReq) (res *v1.JobPageRes, err error) {
	list, total, err := jobSpec.Page(ctx, req, req.GetPageNo(), req.GetPageSize())
	if err != nil {
		return nil, err
	}
	return v1.NewPageData(req.GetPageNo(), req.GetPageSize(), total, derefSlice(list)), nil
}

// List 任务列表
// 按关键词、状态等条件返回定时任务列表
func (s *sJob) List(ctx context.Context, req *v1.JobListReq) (res v1.JobListRes, err error) {
	list, err := jobSpec.List(ctx, req)
	if err != nil {
		return nil, err
	}
	return derefSlice(list), nil
}

// Detail 任务详情
//
// 因 JobDetailRes 为嵌入结构，基座 FillTextFields 的反射计划要求顶层字段，故手写实现
func (s *sJob) Detail(ctx context.Context, req *v1.JobDetailReq) (res *v1.JobDetailRes, err error) {
	one, err := dao.Job.Ctx(ctx).Where("id", req.Id).One()
	if err != nil {
		return nil, err
	}
	if len(one) == 0 {
		return nil, gerror.NewCode(consts.CodeRecordNotFound)
	}
	var item v1.JobItem
	if err = one.Struct(&item); err != nil {
		return nil, err
	}
	FillTextFields(ctx, &item, jobSerializeMaps)
	return &v1.JobDetailRes{JobItem: item}, nil
}

// Add 新增任务
// 默认未发布状态，发布操作走 Status 或 Resume 方法
func (s *sJob) Add(ctx context.Context, req *v1.JobAddReq) (res *v1.IdRes, err error) {
	id, err := jobSpec.Add(ctx, req)
	if err != nil {
		return nil, err
	}
	return &v1.IdRes{Id: id}, nil
}

// Update 修改任务
// 运行中任务先摘除再重挂
func (s *sJob) Update(ctx context.Context, req *v1.JobUpdateReq) (res *v1.JobUpdateRes, err error) {
	if err = jobSpec.Update(ctx, req.Id, req); err != nil {
		return nil, err
	}
	return &v1.JobUpdateRes{}, nil
}

// Delete 删除任务
// 软删记录并摘除调度
func (s *sJob) Delete(ctx context.Context, req *v1.JobDeleteReq) (res *v1.JobDeleteRes, err error) {
	if err = jobSpec.Delete(ctx, req.Id); err != nil {
		return nil, err
	}
	return &v1.JobDeleteRes{}, nil
}

// BatchDelete 批量删除任务
// 按任务ID列表批量删除记录
func (s *sJob) BatchDelete(ctx context.Context, req *v1.JobBatchDeleteReq) (res *v1.JobBatchDeleteRes, err error) {
	if err = jobSpec.BatchDelete(ctx, req.Ids); err != nil {
		return nil, err
	}
	return &v1.JobBatchDeleteRes{}, nil
}

// ---------- 非标准方法（手写实现） ----------

// RunOnce 立即执行一次
// 异步触发指定任务，未发布任务也可执行
func (s *sJob) RunOnce(ctx context.Context, req *v1.JobRunOnceReq) (res *v1.JobRunOnceRes, err error) {
	one, err := dao.Job.Ctx(ctx).Where("id", req.Id).One()
	if err != nil {
		return nil, err
	}
	if len(one) == 0 {
		return nil, gerror.NewCode(consts.CodeRecordNotFound)
	}
	go executeJob(context.Background(), req.Id, true)
	return &v1.JobRunOnceRes{}, nil
}

// Pause 暂停任务
// 摘除调度并将状态置为暂停
func (s *sJob) Pause(ctx context.Context, req *v1.JobPauseReq) (res *v1.JobPauseRes, err error) {
	one, err := dao.Job.Ctx(ctx).Where("id", req.Id).One()
	if err != nil {
		return nil, err
	}
	if len(one) == 0 {
		return nil, gerror.NewCode(consts.CodeRecordNotFound)
	}
	removeJobFromScheduler(req.Id)
	if _, err = dao.Job.Ctx(ctx).Where("id", req.Id).Data(map[string]any{"status": 2}).Update(); err != nil {
		return nil, err
	}
	return &v1.JobPauseRes{}, nil
}

// Resume 恢复任务
// 重新挂载调度并将状态置为运行中
func (s *sJob) Resume(ctx context.Context, req *v1.JobResumeReq) (res *v1.JobResumeRes, err error) {
	one, err := dao.Job.Ctx(ctx).Where("id", req.Id).One()
	if err != nil {
		return nil, err
	}
	if len(one) == 0 {
		return nil, gerror.NewCode(consts.CodeRecordNotFound)
	}
	if err = addJobToScheduler(req.Id, one["cron_expression"].String(),
		one["execute_policy"].Int(), one["is_sync"].Int()); err != nil {
		return nil, gerror.NewCodef(consts.CodeParamError, "Cron表达式无效: %v", err)
	}
	if _, err = dao.Job.Ctx(ctx).Where("id", req.Id).Data(map[string]any{"status": 1}).Update(); err != nil {
		return nil, err
	}
	return &v1.JobResumeRes{}, nil
}

// Status 设置任务状态
// 运行中状态挂载调度，暂停与未发布状态摘除调度
func (s *sJob) Status(ctx context.Context, req *v1.JobStatusReq) (res *v1.JobStatusRes, err error) {
	one, err := dao.Job.Ctx(ctx).Where("id", req.Id).One()
	if err != nil {
		return nil, err
	}
	if len(one) == 0 {
		return nil, gerror.NewCode(consts.CodeRecordNotFound)
	}
	removeJobFromScheduler(req.Id)
	if req.Status == 1 {
		if err = addJobToScheduler(req.Id, one["cron_expression"].String(),
			one["execute_policy"].Int(), one["is_sync"].Int()); err != nil {
			return nil, gerror.NewCodef(consts.CodeParamError, "Cron表达式无效: %v", err)
		}
	}
	if _, err = dao.Job.Ctx(ctx).Where("id", req.Id).Data(map[string]any{"status": req.Status}).Update(); err != nil {
		return nil, err
	}
	return &v1.JobStatusRes{}, nil
}

// DaemonStatus 查询调度器状态
// 返回调度器中运行的任务数、任务键与启动时间
func (s *sJob) DaemonStatus(ctx context.Context, req *v1.JobDaemonStatusReq) (res *v1.JobDaemonStatusRes, err error) {
	entries := gcron.Entries()
	ids := make([]string, 0, len(entries))
	for _, e := range entries {
		ids = append(ids, e.Name)
	}
	startText := ""
	if !daemonStart.IsZero() {
		startText = daemonStart.Format("2006-01-02 15:04:05")
	}
	return &v1.JobDaemonStatusRes{
		Running:   len(entries),
		JobIds:    ids,
		StartTime: startText,
	}, nil
}
