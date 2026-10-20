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

	v1 "easygoadmin/api/v1"
	"easygoadmin/internal/dao"
	"easygoadmin/internal/service"
	"easygoadmin/utility/excel"
	"easygoadmin/utility/upload"
)

// ============================== 案例演示业务实现 ==============================

// sExample 案例演示业务实现
// 实现 service.IExample 接口，接口定义见 api/v1/example.go
type sExample struct{}

func init() {
	service.RegisterExample(&sExample{})
}

// exampleSerializeMaps 枚举字段到字典编码的映射
// {field}Text 由基座自动填充
var exampleSerializeMaps = map[string]string{
	"status": "example_status",
}

// exampleSpec 案例演示模块声明
//
// 采用声明式实现：查询、输出与规则均在声明中定义，方法仅做委托
var exampleSpec = ModuleSpec[v1.ExampleItem, v1.ExampleItem, v1.ExampleDetailRes]{
	Model:         dao.Example.Ctx,
	LikeFields:    []string{"name"},
	EqFields:      []string{"type", "status"},
	OrderBy:       "sort, id",
	SerializeMaps: exampleSerializeMaps,
	FileFields: []string{
		"avatar",
	},
	FileSaveDir: "example",
	ToData:      exampleToData,
}

// exampleToData 全量表单语义转换
// 新增与修改均提交全部字段，兼容两种入参类型
func exampleToData(req any) map[string]any {
	switch r := req.(type) {
	case *v1.ExampleAddReq:
		return map[string]any{
			"name":   r.Name,
			"avatar": r.Avatar,
			"type":   r.Type,
			"status": r.Status,
			"sort":   r.Sort,
		}
	case *v1.ExampleUpdateReq:
		return map[string]any{
			"name":   r.Name,
			"avatar": r.Avatar,
			"type":   r.Type,
			"status": r.Status,
			"sort":   r.Sort,
		}
	}
	panic(fmt.Sprintf("exampleToData: 不支持的入参类型 %T", req))
}

// ---------- 标准方法：委托基座 ----------

// Page 案例演示分页
// 按查询条件分页返回案例演示列表
func (s *sExample) Page(ctx context.Context, req *v1.ExamplePageReq) (res *v1.ExamplePageRes, err error) {
	list, total, err := exampleSpec.Page(ctx, req, req.GetPageNo(), req.GetPageSize())
	if err != nil {
		return nil, err
	}
	return v1.NewPageData(req.GetPageNo(), req.GetPageSize(), total, derefSlice(list)), nil
}

// List 案例演示平铺列表
// data 直接为数组，树形结构由前端 buildTree 构建
func (s *sExample) List(ctx context.Context, req *v1.ExampleListReq) (res v1.ExampleListRes, err error) {
	list, err := exampleSpec.List(ctx, req)
	if err != nil {
		return nil, err
	}
	return derefSlice(list), nil
}

// Detail 案例演示详情
// 按主键ID查询详情
func (s *sExample) Detail(ctx context.Context, req *v1.ExampleDetailReq) (res *v1.ExampleDetailRes, err error) {
	return exampleSpec.Detail(ctx, req.Id, req)
}

// Add 新增案例演示
// 新增成功后返回新记录ID
func (s *sExample) Add(ctx context.Context, req *v1.ExampleAddReq) (res *v1.ExampleIdRes, err error) {
	id, err := exampleSpec.Add(ctx, req)
	if err != nil {
		return nil, err
	}
	return &v1.ExampleIdRes{Id: id}, nil
}

// Update 修改案例演示
// 按主键ID更新记录
func (s *sExample) Update(ctx context.Context, req *v1.ExampleUpdateReq) (res *v1.ExampleUpdateRes, err error) {
	if err = exampleSpec.Update(ctx, req.Id, req); err != nil {
		return nil, err
	}
	return &v1.ExampleUpdateRes{}, nil
}

// Status 案例演示状态
// 非标准方法，保留手写实现，设置记录的状态字段
func (s *sExample) Status(ctx context.Context, req *v1.ExampleStatusReq) (res *v1.ExampleStatusRes, err error) {
	if _, err = dao.Example.Ctx(ctx).Where("id", req.Id).Data(map[string]any{"status": req.Status}).Update(); err != nil {
		return nil, err
	}
	return &v1.ExampleStatusRes{}, nil
}

// Delete 删除案例演示
// 按主键ID删除单条记录
func (s *sExample) Delete(ctx context.Context, req *v1.ExampleDeleteReq) (res *v1.ExampleDeleteRes, err error) {
	if err = exampleSpec.Delete(ctx, req.Id); err != nil {
		return nil, err
	}
	return &v1.ExampleDeleteRes{}, nil
}

// BatchDelete 批量删除案例演示
// 按主键ID列表批量删除记录
func (s *sExample) BatchDelete(ctx context.Context, req *v1.ExampleBatchDeleteReq) (res *v1.ExampleBatchDeleteRes, err error) {
	if err = exampleSpec.BatchDelete(ctx, req.Ids); err != nil {
		return nil, err
	}
	return &v1.ExampleBatchDeleteRes{}, nil
}

// Import 导入案例演示
// 表头为表单字段名，按列序写入
func (s *sExample) Import(ctx context.Context, req *v1.ExampleImportReq) (res *v1.ExampleImportRes, err error) {
	_, rows, err := excel.ReadSheet(req.File)
	if err != nil {
		return nil, err
	}
	res = &v1.ExampleImportRes{Total: len(rows)}
	for i, row := range rows {
		if cellAt(row, 0) == "" {
			res.Failed++
			res.Reasons = append(res.Reasons, fmt.Sprintf("第%d行：案例名称为空", i+2))
			continue
		}
		if _, e := dao.Example.Ctx(ctx).Data(map[string]any{
			"name":   cellAt(row, 0),
			"avatar": cellAt(row, 1),
			"type":   cellInt(row, 2, 0),
			"status": cellInt(row, 3, 0),
			"sort":   cellInt(row, 4, 0),
		}).Insert(); e != nil {
			res.Failed++
			res.Reasons = append(res.Reasons, fmt.Sprintf("第%d行：写入失败 %v", i+2, e))
			continue
		}
		res.Success++
	}
	return res, nil
}

// Export 导出案例演示
// 与列表同套筛选条件，枚举列输出字典文本
func (s *sExample) Export(ctx context.Context, req *v1.ExampleExportReq) (res *v1.ExampleExportRes, err error) {
	rows, err := exampleSpec.applyFilters(ctx, exampleSpec.Model(ctx), req).Order("sort, id").All()
	if err != nil {
		return nil, err
	}
	headers := []string{"案例名称", "案例图片", "案例类型", "案例状态", "案例排序"}
	data := make([][]any, 0, len(rows))
	for _, r := range rows {
		data = append(data, []any{r["name"].String(), r["avatar"].String(), DictGetText(ctx, "example_type", r["type"]), DictGetText(ctx, "example_status", r["status"]), r["sort"].String()})
	}
	fileName, _, err := excel.WriteSheet(upload.BaseDir(ctx), "案例演示数据", headers, data)
	if err != nil {
		return nil, err
	}
	url := upload.GetFileUrl(ctx, fileName)
	return &url, nil
}
