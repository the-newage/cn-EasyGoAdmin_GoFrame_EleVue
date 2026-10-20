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
	"fmt"

	"github.com/gogf/gf/v2/errors/gerror"

	"easygoadmin/api/v1"
	"easygoadmin/internal/consts"
	"easygoadmin/internal/dao"
	"easygoadmin/internal/service"
	"easygoadmin/utility/excel"
	"easygoadmin/utility/upload"
)

// ============================== 职级业务实现 ==============================

// sLevel 职级业务实现
// 实现 service.ILevel 接口，接口定义见 api/v1/level.go
type sLevel struct{}

func init() {
	service.RegisterLevel(&sLevel{})
}

// levelSerializeMaps 枚举字段到字典编码的映射
// {field}Text 由基座自动填充
var levelSerializeMaps = map[string]string{
	"status": "level_status",
}

// levelSpec 职级模块声明
//
// 采用声明式实现：查询、输出与规则均在声明中定义，方法仅做委托
var levelSpec = ModuleSpec[v1.LevelItem, v1.LevelListItem, v1.LevelDetailRes]{
	Model:         dao.Level.Ctx,
	LikeFields:    []string{"name"},   // 名称模糊（querySchemas 输入）
	EqFields:      []string{"status"}, // status=0 视为未筛选（零值跳过）
	OrderBy:       "sort, id",
	SerializeMaps: levelSerializeMaps,
	ToData:        levelToData,
	BeforeAdd:     levelDefaultStatus, // status 缺省置 1
	// Delete 守卫（存在性 + 用户引用检查）
	BeforeDelete: levelDeleteGuard,
}

// levelToData 全量表单语义转换
// 新增与修改均提交全部字段，兼容两种入参类型
func levelToData(req any) map[string]any {
	switch r := req.(type) {
	case *v1.LevelAddReq:
		return map[string]any{"name": r.Name, "status": r.Status, "sort": r.Sort}
	case *v1.LevelUpdateReq:
		return map[string]any{"name": r.Name, "status": r.Status, "sort": r.Sort}
	}
	panic(fmt.Sprintf("levelToData: 不支持的入参类型 %T", req))
}

// levelDefaultStatus 状态缺省处理
// 状态为零值时置为 1（正常）
func levelDefaultStatus(ctx context.Context, data map[string]any, req any) error {
	if v, ok := data["status"].(int); ok && v == 0 {
		data["status"] = 1
	}
	return nil
}

// levelDeleteGuard 删除守卫
// 校验记录存在性，存在用户引用时禁止删除
func levelDeleteGuard(ctx context.Context, id int64) error {
	if n, err := dao.Level.Ctx(ctx).Where("id", id).Count(); err != nil {
		return err
	} else if n == 0 {
		return gerror.NewCode(consts.CodeRecordNotFound)
	}
	if n, err := dao.User.Ctx(ctx).Where("level_id", id).Count(); err != nil {
		return err
	} else if n > 0 {
		return gerror.NewCode(consts.CodeInUse)
	}
	return nil
}

// ---------- 标准方法：3行委托 ----------

// Page 职级分页
// 分页返回职级列表，并补充关联用户数与状态文本
func (s *sLevel) Page(ctx context.Context, req *v1.LevelPageReq) (res *v1.LevelPageRes, err error) {
	list, total, err := levelSpec.Page(ctx, req, req.GetPageNo(), req.GetPageSize())
	if err != nil {
		return nil, err
	}
	fillUserCount(ctx, "level_id", &list)
	return v1.NewPageData(req.GetPageNo(), req.GetPageSize(), total, derefSlice(list)), nil
}

// List 职级全量列表
// 返回精简结构的职级列表
func (s *sLevel) List(ctx context.Context, req *v1.LevelListReq) (res v1.LevelListRes, err error) {
	list, err := levelSpec.List(ctx, req)
	if err != nil {
		return nil, err
	}
	return derefSlice(list), nil
}

// Detail 职级详情
// 按职级ID查询详情
func (s *sLevel) Detail(ctx context.Context, req *v1.LevelDetailReq) (res *v1.LevelDetailRes, err error) {
	return levelSpec.Detail(ctx, req.Id, req)
}

// Add 新增职级
// 新增成功后返回新记录ID
func (s *sLevel) Add(ctx context.Context, req *v1.LevelAddReq) (res *v1.IdRes, err error) {
	id, err := levelSpec.Add(ctx, req)
	if err != nil {
		return nil, err
	}
	return &v1.IdRes{Id: id}, nil
}

// Update 修改职级
// 按职级ID更新记录
func (s *sLevel) Update(ctx context.Context, req *v1.LevelUpdateReq) (res *v1.LevelUpdateRes, err error) {
	err = levelSpec.Update(ctx, req.Id, req)
	if err != nil {
		return nil, err
	}
	return &v1.LevelUpdateRes{}, nil
}

// Status 职级状态
// 非标准方法，手写实现，设置职级的正常或禁用状态
func (s *sLevel) Status(ctx context.Context, req *v1.LevelStatusReq) (res *v1.LevelStatusRes, err error) {
	if _, err = dao.Level.Ctx(ctx).Where("id", req.Id).Data(map[string]any{"status": req.Status}).Update(); err != nil {
		return nil, err
	}
	return &v1.LevelStatusRes{}, nil
}

// Delete 删除职级
// 按职级ID删除单条记录
func (s *sLevel) Delete(ctx context.Context, req *v1.LevelDeleteReq) (res *v1.LevelDeleteRes, err error) {
	if err = levelSpec.Delete(ctx, req.Id); err != nil {
		return nil, err
	}
	return &v1.LevelDeleteRes{}, nil
}

// BatchDelete 批量删除职级
// 按职级ID列表批量删除记录
func (s *sLevel) BatchDelete(ctx context.Context, req *v1.LevelBatchDeleteReq) (res *v1.LevelBatchDeleteRes, err error) {
	if err = levelSpec.BatchDelete(ctx, req.Ids); err != nil {
		return nil, err
	}
	return &v1.LevelBatchDeleteRes{}, nil
}

// ---------- 非标准方法（手写实现）（Import/Export）----------

// Import 导入职级
//
// 表格列依次为名称、状态、排序，名称重复时跳过并记录原因
func (s *sLevel) Import(ctx context.Context, req *v1.LevelImportReq) (res *v1.LevelImportRes, err error) {
	_, rows, err := excel.ReadSheet(req.File)
	if err != nil {
		return nil, err
	}
	res = &v1.LevelImportRes{Total: len(rows)}
	tenantId := consts.GetTenantID(ctx)
	for i, row := range rows {
		name := cellAt(row, 0)
		if name == "" {
			res.Failed++
			res.Reasons = append(res.Reasons, fmt.Sprintf("第%d行：名称为空", i+2))
			continue
		}
		count, e := dao.Level.Ctx(ctx).Where("name", name).Count()
		if e != nil {
			return nil, e
		}
		if count > 0 {
			res.Failed++
			res.Reasons = append(res.Reasons, fmt.Sprintf("第%d行：名称已存在", i+2))
			continue
		}
		status := cellInt(row, 1, 1)
		if status <= 0 {
			status = 1
		}
		if _, e = dao.Level.Ctx(ctx).Data(map[string]any{
			"name":      name,
			"status":    status,
			"sort":      cellInt(row, 2, 0),
			"tenant_id": tenantId,
		}).Insert(); e != nil {
			res.Failed++
			res.Reasons = append(res.Reasons, fmt.Sprintf("第%d行：写入失败 %v", i+2, e))
			continue
		}
		res.Success++
	}
	return res, nil
}

// Export 导出职级
// 按查询条件导出职级数据并返回文件访问地址
func (s *sLevel) Export(ctx context.Context, req *v1.LevelExportReq) (res *v1.LevelExportRes, err error) {
	m := levelSpec.applyFilters(ctx, levelSpec.Model(ctx), req)
	rows, err := m.Order("sort, id").All()
	if err != nil {
		return nil, err
	}
	headers := []string{"名称", "状态", "排序"}
	data := make([][]any, 0, len(rows))
	for _, r := range rows {
		data = append(data, []any{r["name"].String(), userStatusText(r["status"].Int()), r["sort"].Int()})
	}
	fileName, _, err := excel.WriteSheet(upload.BaseDir(ctx), "职级数据", headers, data)
	if err != nil {
		return nil, err
	}
	url := upload.GetFileUrl(ctx, fileName)
	return &url, nil
}
