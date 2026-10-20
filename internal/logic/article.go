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
	"strings"

	"easygoadmin/api/v1"
	"easygoadmin/internal/consts"
	"easygoadmin/internal/dao"
	"easygoadmin/internal/service"
)

// ============================== 文章业务实现 ==============================

// sArticle 文章业务实现
// 实现 service.IArticle 接口，承载文章模块的业务逻辑
type sArticle struct{}

// articleSerializeMaps 枚举字段到字典编码的映射
// {field}Text 由 dict_item 解析，字典管理页可修改文案
var articleSerializeMaps = map[string]string{
	"status": "article_status",
}

// articleSpec 文章模块声明
//
// Status 特有方法与文件字段保存在 Add/Update 手写体内处理；
// 文件与富文本读侧还原由基座 readTransform 按下列声明自动完成
var articleSpec = ModuleSpec[v1.ArticleItem, v1.ArticleListItem, v1.ArticleDetailRes]{
	Model:           dao.Article.Ctx,
	LikeFields:      []string{"title", "author"},
	EqFields:        []string{"status", "category_id"}, // 零值不筛（列表初始显示全部）
	OrderBy:         "sort, id DESC",
	ListOrderBy:     "sort, id",
	SerializeMaps:   articleSerializeMaps,
	FileFields:      []string{"cover"},   // 读侧补域名
	MultiFileFields: []string{"images"},  // 逗号多图：读侧逐个补域名
	ContentFields:   []string{"content"}, // 富文本：读侧 [IMG_URL] 还原
	FileSaveDir:     "article",
	AfterPageList:   articleExpandCategory,  // 行内分类名批量补全
	AfterDetail:     articleSplitImagesList, // 多图拆分为 imagesList 数组
}

// articleExpandCategory 分页行分类名批量补全
// 将分类ID映射为分类名称，避免逐行查询
func articleExpandCategory(ctx context.Context, list []*v1.ArticleItem, req any) error {
	catNames := (&sArticle{}).categoryNames(ctx)
	for _, item := range list {
		item.CategoryName = catNames[int64(item.CategoryId)]
	}
	return nil
}

func init() {
	service.RegisterArticle(&sArticle{})
}

// Page 文章分页
// 按查询条件分页返回文章列表
func (s *sArticle) Page(ctx context.Context, req *v1.ArticlePageReq) (res *v1.ArticlePageRes, err error) {
	list, total, err := articleSpec.Page(ctx, req, req.GetPageNo(), req.GetPageSize())
	if err != nil {
		return nil, err
	}
	return v1.NewPageData(req.GetPageNo(), req.GetPageSize(), total, derefSlice(list)), nil
}

// List 文章全量列表
// 按查询条件返回文章精简列表
func (s *sArticle) List(ctx context.Context, req *v1.ArticleListReq) (res v1.ArticleListRes, err error) {
	list, err := articleSpec.List(ctx, req)
	if err != nil {
		return nil, err
	}
	return derefSlice(list), nil
}

// Detail 文章详情
// 按文章ID查询详情
func (s *sArticle) Detail(ctx context.Context, req *v1.ArticleDetailReq) (res *v1.ArticleDetailRes, err error) {
	return articleSpec.Detail(ctx, req.Id, req)
}

// articleSplitImagesList 详情装配
// 将多图逗号串（读转换已补域名）拆分为 imagesList 数组
func articleSplitImagesList(ctx context.Context, item *v1.ArticleDetailRes, req any) error {
	if item.Images == "" {
		return nil
	}
	for _, u := range strings.Split(item.Images, ",") {
		if u = strings.TrimSpace(u); u != "" {
			item.ImagesList = append(item.ImagesList, u)
		}
	}
	return nil
}

// Add 新增文章
// 封面与多图由临时目录迁移至正式目录，正文富文本本地化后落库
func (s *sArticle) Add(ctx context.Context, req *v1.ArticleAddReq) (res *v1.IdRes, err error) {
	// 封面/多图 temp → 正式目录；正文富文本本地化（[IMG_URL] 占位符契约）
	req.Cover = SaveFileField(ctx, req.Cover, "article")
	req.Images = SaveMultiFileField(ctx, req.Images, "article")
	req.Content = LocalizeContentFields(ctx, req.Content, "article")
	result, err := dao.Article.Ctx(ctx).Data(map[string]any{
		"category_id": req.CategoryId,
		"title":       req.Title,
		"cover":       req.Cover,
		"images":      req.Images,
		"intro":       req.Intro,
		"content":     req.Content,
		"author":      req.Author,
		"status":      req.Status,
		"sort":        req.Sort,
		"tenant_id":   consts.GetTenantID(ctx),
	}).Insert()
	if err != nil {
		return nil, err
	}
	idv, _ := result.LastInsertId()
	return &v1.IdRes{Id: idv}, nil
}

// Update 修改文章
// 封面与多图由临时目录迁移至正式目录，正文富文本本地化后更新
func (s *sArticle) Update(ctx context.Context, req *v1.ArticleUpdateReq) (res *v1.ArticleUpdateRes, err error) {
	req.Cover = SaveFileField(ctx, req.Cover, "article")
	req.Images = SaveMultiFileField(ctx, req.Images, "article")
	req.Content = LocalizeContentFields(ctx, req.Content, "article")
	if _, err = dao.Article.Ctx(ctx).Where("id", req.Id).Data(map[string]any{
		"category_id": req.CategoryId,
		"title":       req.Title,
		"cover":       req.Cover,
		"images":      req.Images,
		"intro":       req.Intro,
		"content":     req.Content,
		"author":      req.Author,
		"status":      req.Status,
		"click":       req.Click,
		"sort":        req.Sort,
	}).Update(); err != nil {
		return nil, err
	}
	return &v1.ArticleUpdateRes{}, nil
}

// Status 文章状态
// 设置文章的草稿或发布状态
func (s *sArticle) Status(ctx context.Context, req *v1.ArticleStatusReq) (res *v1.ArticleStatusRes, err error) {
	if _, err = dao.Article.Ctx(ctx).Where("id", req.Id).Data(map[string]any{"status": req.Status}).Update(); err != nil {
		return nil, err
	}
	return &v1.ArticleStatusRes{}, nil
}

// Delete 删除文章
// 按文章ID删除单条记录
func (s *sArticle) Delete(ctx context.Context, req *v1.ArticleDeleteReq) (res *v1.ArticleDeleteRes, err error) {
	if err = articleSpec.Delete(ctx, req.Id); err != nil {
		return nil, err
	}
	return &v1.ArticleDeleteRes{}, nil
}

// BatchDelete 批量删除文章
// 按文章ID列表批量删除记录
func (s *sArticle) BatchDelete(ctx context.Context, req *v1.ArticleBatchDeleteReq) (res *v1.ArticleBatchDeleteRes, err error) {
	if err = articleSpec.BatchDelete(ctx, req.Ids); err != nil {
		return nil, err
	}
	return &v1.ArticleBatchDeleteRes{}, nil
}

// categoryNames 查询分类ID到名称的映射
func (s *sArticle) categoryNames(ctx context.Context) map[int64]string {
	result := make(map[int64]string)
	rows, err := dao.Category.Ctx(ctx).Fields("id, category_name").All()
	if err != nil {
		return result
	}
	for _, r := range rows {
		result[r["id"].Int64()] = r["category_name"].String()
	}
	return result
}
