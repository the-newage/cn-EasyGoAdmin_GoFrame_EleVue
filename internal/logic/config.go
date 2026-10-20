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
	"strings"

	"easygoadmin/api/v1"
	"easygoadmin/internal/dao"
	"easygoadmin/internal/service"
	"easygoadmin/utility/upload"
)

// ============================== 配置业务实现 ==============================

// sConfig 配置业务实现
// 实现 service.IConfig 接口，承载配置模块的业务逻辑
type sConfig struct{}

// configSpec 配置模块声明
//
// List 含嵌套配置项，采用手写实现
var configSpec = ModuleSpec[v1.ConfigItem, v1.ConfigItem, v1.ConfigDetailRes]{
	Model:      dao.Config.Ctx,
	LikeFields: []string{"name", "code"},
	OrderBy:    "sort, id",
}

func init() {
	service.RegisterConfig(&sConfig{})
}

// Page 配置分页
// 按查询条件分页返回配置列表
func (s *sConfig) Page(ctx context.Context, req *v1.ConfigPageReq) (res *v1.ConfigPageRes, err error) {
	list, total, err := configSpec.Page(ctx, req, req.GetPageNo(), req.GetPageSize())
	if err != nil {
		return nil, err
	}
	return v1.NewPageData(req.GetPageNo(), req.GetPageSize(), total, derefSlice(list)), nil
}

// List 配置列表（含配置项）
// 返回配置列表，每项携带其下的配置项集合
func (s *sConfig) List(ctx context.Context, req *v1.ConfigListReq) (res v1.ConfigListRes, err error) {
	m := dao.Config.Ctx(ctx)
	if req.Keyword != "" {
		kw := "%" + req.Keyword + "%"
		m = m.WhereLike("name", kw).WhereOrLike("code", kw)
	}
	rows, err := m.Order("sort, id").All()
	if err != nil {
		return nil, err
	}
	itemsByConfig, err := s.itemsGrouped(ctx)
	if err != nil {
		return nil, err
	}
	list := make([]v1.ConfigWithItems, 0, len(rows))
	for _, r := range rows {
		cid := r["id"].Int64()
		list = append(list, v1.ConfigWithItems{
			ConfigItem: v1.ConfigItem{
				Id:         cid,
				Name:       r["name"].String(),
				Code:       r["code"].String(),
				Sort:       r["sort"].Int(),
				Note:       r["note"].String(),
				CreateUser: r["create_user"].String(),
			},
			Items: itemsByConfig[cid],
		})
	}
	return list, nil
}

// Detail 配置详情
// 按配置ID查询详情及其下的配置项集合
func (s *sConfig) Detail(ctx context.Context, req *v1.ConfigDetailReq) (res *v1.ConfigDetailRes, err error) {
	return configSpec.Detail(ctx, req.Id, req)
}

// Add 新增配置
// 校验配置编码唯一后新增，返回新记录ID
func (s *sConfig) Add(ctx context.Context, req *v1.ConfigAddReq) (res *v1.IdRes, err error) {
	if err = CheckUnique(ctx, dao.Config.Ctx(ctx), "code", req.Code, 0, "配置编码"); err != nil {
		return nil, err
	}
	result, err := dao.Config.Ctx(ctx).Data(map[string]any{
		"name": req.Name,
		"code": req.Code,
		"sort": req.Sort,
		"note": req.Note,
	}).Insert()
	if err != nil {
		return nil, err
	}
	idv, _ := result.LastInsertId()
	return &v1.IdRes{Id: idv}, nil
}

// Update 修改配置
// 按配置ID更新记录
func (s *sConfig) Update(ctx context.Context, req *v1.ConfigUpdateReq) (res *v1.ConfigUpdateRes, err error) {
	if _, err = dao.Config.Ctx(ctx).Where("id", req.Id).Data(map[string]any{
		"name": req.Name,
		"code": req.Code,
		"sort": req.Sort,
		"note": req.Note,
	}).Update(); err != nil {
		return nil, err
	}
	return &v1.ConfigUpdateRes{}, nil
}

// Delete 删除配置
// 按配置ID删除记录，并连带删除其下配置项
func (s *sConfig) Delete(ctx context.Context, req *v1.ConfigDeleteReq) (res *v1.ConfigDeleteRes, err error) {
	if err = configSpec.Delete(ctx, req.Id); err != nil {
		return nil, err
	}
	if _, err = dao.ConfigItem.Ctx(ctx).Where("config_id", req.Id).Delete(); err != nil {
		return nil, err
	}
	return &v1.ConfigDeleteRes{}, nil
}

// BatchDelete 批量删除配置
// 按配置ID列表批量删除记录，并连带删除其下配置项
func (s *sConfig) BatchDelete(ctx context.Context, req *v1.ConfigBatchDeleteReq) (res *v1.ConfigBatchDeleteRes, err error) {
	if err = configSpec.BatchDelete(ctx, req.Ids); err != nil {
		return nil, err
	}
	if _, err = dao.ConfigItem.Ctx(ctx).WhereIn("config_id", req.Ids).Delete(); err != nil {
		return nil, err
	}
	return &v1.ConfigBatchDeleteRes{}, nil
}

// WebIndex 网站配置聚合
//
// 按配置分组返回，供网站配置页渲染；
// 文件类字段拼接访问URL，options 解析为 param 映射
func (s *sConfig) WebIndex(ctx context.Context, req *v1.ConfigWebIndexReq) (res *v1.ConfigWebIndexRes, err error) {
	configs, err := dao.Config.Ctx(ctx).Order("sort, id").All()
	if err != nil {
		return nil, err
	}
	items, err := dao.ConfigItem.Ctx(ctx).Order("sort, id").All()
	if err != nil {
		return nil, err
	}
	byConfig := make(map[int64][]v1.ConfigWebItem, len(configs))
	for _, r := range items {
		value := r["value"].String()
		typ := r["type"].String()
		options := r["options"].String()
		item := v1.ConfigWebItem{
			Id:      r["id"].Int64(),
			Name:    r["name"].String(),
			Code:    r["code"].String(),
			Value:   value,
			Type:    typ,
			Options: options,
			Param:   parseConfigOptions(options),
		}
		switch typ {
		case "image":
			item.Value = upload.GetFileUrl(ctx, value)
		case "file":
			item.Value = configItemFileURL(ctx, value)
		case "images":
			item.ValueList = configItemListURLs(ctx, value, false)
		case "files":
			item.ValueList = configItemListURLs(ctx, value, true)
		case "ueditor":
			item.Value = ExpandContentFields(ctx, value)
		}
		cid := r["config_id"].Int64()
		byConfig[cid] = append(byConfig[cid], item)
	}
	groups := make([]v1.ConfigWebGroup, 0, len(configs))
	for _, c := range configs {
		cid := c["id"].Int64()
		groups = append(groups, v1.ConfigWebGroup{
			ConfigName: c["name"].String(),
			DataList:   byConfig[cid],
		})
	}
	return &groups, nil
}

// WebSave 批量保存网站配置
//
// 接收配置项编码到值的映射；
// 文件类字段执行转正与剥离域名，富文本执行本地化
func (s *sConfig) WebSave(ctx context.Context, req *v1.ConfigWebSaveReq) (res *v1.ConfigWebSaveRes, err error) {
	for code, raw := range req.Data {
		one, e := dao.ConfigItem.Ctx(ctx).Where("code", code).One()
		if e != nil {
			return nil, e
		}
		if len(one) == 0 {
			continue // 未知编码跳过（前端可能携带非配置项字段）
		}
		value := configItemSaveValue(ctx, one["type"].String(), configValueToString(raw))
		if _, e = dao.ConfigItem.Ctx(ctx).Where("id", one["id"].Int64()).Data("value", value).Update(); e != nil {
			return nil, e
		}
	}
	return &v1.ConfigWebSaveRes{}, nil
}

// parseConfigOptions 解析选项串
// 将 "value=name,value=name" 格式解析为值与显示名的映射
func parseConfigOptions(options string) map[string]string {
	param := make(map[string]string)
	if options == "" {
		return param
	}
	for _, pair := range strings.Split(options, ",") {
		pair = strings.TrimSpace(pair)
		if pair == "" {
			continue
		}
		value, name, found := strings.Cut(pair, "=")
		if !found {
			continue
		}
		param[strings.TrimSpace(value)] = strings.TrimSpace(name)
	}
	return param
}

// configValueToString 前端表单值归一化为字符串
// 多选或多图为数组时按逗号拼接
func configValueToString(v any) string {
	switch t := v.(type) {
	case nil:
		return ""
	case string:
		return t
	case []any:
		parts := make([]string, 0, len(t))
		for _, e := range t {
			parts = append(parts, fmt.Sprintf("%v", e))
		}
		return strings.Join(parts, ",")
	default:
		return fmt.Sprintf("%v", t)
	}
}

// itemsByConfigId 查询单个配置的配置项
func (s *sConfig) itemsByConfigId(ctx context.Context, configId int64) ([]v1.ConfigItemRow, error) {
	rows, err := dao.ConfigItem.Ctx(ctx).Where("config_id", configId).Order("sort, id").All()
	if err != nil {
		return nil, err
	}
	items := make([]v1.ConfigItemRow, 0, len(rows))
	for _, r := range rows {
		items = append(items, v1.ConfigItemRow{
			Id:     r["id"].Int64(),
			Name:   r["name"].String(),
			Code:   r["code"].String(),
			Value:  r["value"].String(),
			Type:   r["type"].String(),
			Status: r["status"].Int(),
			Sort:   r["sort"].Int(),
		})
	}
	return items, nil
}

// itemsGrouped 全量配置项按 config_id 分组
func (s *sConfig) itemsGrouped(ctx context.Context) (map[int64][]v1.ConfigItemRow, error) {
	result := make(map[int64][]v1.ConfigItemRow)
	rows, err := dao.ConfigItem.Ctx(ctx).Order("sort, id").All()
	if err != nil {
		return nil, err
	}
	for _, r := range rows {
		cid := r["config_id"].Int64()
		result[cid] = append(result[cid], v1.ConfigItemRow{
			Id:     r["id"].Int64(),
			Name:   r["name"].String(),
			Code:   r["code"].String(),
			Value:  r["value"].String(),
			Type:   r["type"].String(),
			Status: r["status"].Int(),
			Sort:   r["sort"].Int(),
		})
	}
	return result, nil
}
