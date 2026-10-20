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
	"github.com/gogf/gf/v2/errors/gerror"

	"easygoadmin/api/v1"
	"easygoadmin/internal/consts"
	"easygoadmin/internal/dao"
	"easygoadmin/internal/service"
	"easygoadmin/utility/upload"
)

// ============================== 文件模板业务实现 ==============================

// sFileTemplate 文件模板业务实现
// 实现 service.IFileTemplate 接口，承载文件模板模块的业务逻辑
type sFileTemplate struct{}

// fileTemplateSpec 文件模板模块声明
// ByCode 为特有方法，保留手写实现
var fileTemplateSpec = ModuleSpec[v1.FileTemplateItem, v1.FileTemplateItem, v1.FileTemplateDetailRes]{
	Model:       dao.FileTemplate.Ctx,
	LikeFields:  []string{"name", "code"},
	OrderBy:     "sort, id",
	FileFields:  []string{"filePath"}, // 读侧自动补域名（file.domainUrl）
	FileSaveDir: "file_template",
}

func init() {
	service.RegisterFileTemplate(&sFileTemplate{})
}

// Page 文件模板分页
// 按查询条件分页返回文件模板列表
func (s *sFileTemplate) Page(ctx context.Context, req *v1.FileTemplatePageReq) (res *v1.FileTemplatePageRes, err error) {
	list, total, err := fileTemplateSpec.Page(ctx, req, req.GetPageNo(), req.GetPageSize())
	if err != nil {
		return nil, err
	}
	return v1.NewPageData(req.GetPageNo(), req.GetPageSize(), total, derefSlice(list)), nil
}

// List 文件模板全量列表
// 按关键词返回文件模板列表
func (s *sFileTemplate) List(ctx context.Context, req *v1.FileTemplateListReq) (res v1.FileTemplateListRes, err error) {
	list, err := fileTemplateSpec.List(ctx, req)
	if err != nil {
		return nil, err
	}
	return derefSlice(list), nil
}

// Detail 文件模板详情
// 按模板ID查询详情
func (s *sFileTemplate) Detail(ctx context.Context, req *v1.FileTemplateDetailReq) (res *v1.FileTemplateDetailRes, err error) {
	return fileTemplateSpec.Detail(ctx, req.Id, req)
}

// Add 新增文件模板
// 校验模板编码唯一后新增，文件路径由临时目录转正，返回新记录ID
func (s *sFileTemplate) Add(ctx context.Context, req *v1.FileTemplateAddReq) (res *v1.IdRes, err error) {
	if err = CheckUnique(ctx, dao.FileTemplate.Ctx(ctx), "code", req.Code, 0, "模板编码"); err != nil {
		return nil, err
	}
	result, err := dao.FileTemplate.Ctx(ctx).Data(map[string]any{
		"name":      req.Name,
		"code":      req.Code,
		"file_name": req.FileName,
		"file_path": SaveFileField(ctx, req.FilePath, "file_template"),
		"note":      req.Note,
		"sort":      req.Sort,
	}).Insert()
	if err != nil {
		return nil, err
	}
	idv, _ := result.LastInsertId()
	return &v1.IdRes{Id: idv}, nil
}

// Update 修改文件模板
// 按模板ID更新记录，文件路径由临时目录转正
func (s *sFileTemplate) Update(ctx context.Context, req *v1.FileTemplateUpdateReq) (res *v1.FileTemplateUpdateRes, err error) {
	if _, err = dao.FileTemplate.Ctx(ctx).Where("id", req.Id).Data(map[string]any{
		"name":      req.Name,
		"code":      req.Code,
		"file_name": req.FileName,
		"file_path": SaveFileField(ctx, req.FilePath, "file_template"),
		"note":      req.Note,
		"sort":      req.Sort,
	}).Update(); err != nil {
		return nil, err
	}
	return &v1.FileTemplateUpdateRes{}, nil
}

// Delete 删除文件模板
// 按模板ID删除单条记录
func (s *sFileTemplate) Delete(ctx context.Context, req *v1.FileTemplateDeleteReq) (res *v1.FileTemplateDeleteRes, err error) {
	if err = fileTemplateSpec.Delete(ctx, req.Id); err != nil {
		return nil, err
	}
	return &v1.FileTemplateDeleteRes{}, nil
}

// BatchDelete 批量删除文件模板
// 按模板ID列表批量删除记录
func (s *sFileTemplate) BatchDelete(ctx context.Context, req *v1.FileTemplateBatchDeleteReq) (res *v1.FileTemplateBatchDeleteRes, err error) {
	if err = fileTemplateSpec.BatchDelete(ctx, req.Ids); err != nil {
		return nil, err
	}
	return &v1.FileTemplateBatchDeleteRes{}, nil
}

// ByCode 按编码查询模板
// 登录用户即可访问，按模板编码查询模板详情
func (s *sFileTemplate) ByCode(ctx context.Context, req *v1.FileTemplateByCodeReq) (res *v1.FileTemplateByCodeRes, err error) {
	one, err := dao.FileTemplate.Ctx(ctx).Where("code", req.Code).One()
	if err != nil {
		return nil, err
	}
	if len(one) == 0 {
		return nil, gerror.NewCode(consts.CodeRecordNotFound)
	}
	return &v1.FileTemplateByCodeRes{FileTemplateItem: s.toItem(ctx, one)}, nil
}

// toItem 数据库记录转换为文件模板行数据
// 文件路径补全为可访问的URL
func (s *sFileTemplate) toItem(ctx context.Context, r gdb.Record) v1.FileTemplateItem {
	return v1.FileTemplateItem{
		Id:         r["id"].Int64(),
		Name:       r["name"].String(),
		Code:       r["code"].String(),
		FileName:   r["file_name"].String(),
		FilePath:   upload.GetFileUrl(ctx, r["file_path"].String()),
		Note:       r["note"].String(),
		Sort:       r["sort"].Int(),
		CreateUser: r["create_user"].String(),
		CreateTime: r["create_time"].GTime(),
	}
}
