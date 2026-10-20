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
	"easygoadmin/internal/dao"
	"easygoadmin/internal/service"
	"easygoadmin/utility/upload"
)

// ============================== 配置项业务实现 ==============================

// sConfigItem 配置项业务实现
// 实现 service.IConfigItem 接口，承载配置项模块的业务逻辑
type sConfigItem struct{}

// configItemSerializeMaps 枚举字段到字典编码的映射
// {field}Text 由 dict_item 解析，字典管理页可修改文案
var configItemSerializeMaps = map[string]string{
	"status": "config_item_status",
}

// configItemSpec 配置项模块声明
var configItemSpec = ModuleSpec[v1.ConfigItemItem, v1.ConfigItemItem, v1.ConfigItemDetailRes]{
	Model:         dao.ConfigItem.Ctx,
	LikeFields:    []string{"name", "code"},
	EqFields:      []string{"status", "config_id"},
	OrderBy:       "sort, id",
	SerializeMaps: configItemSerializeMaps,
	AfterPageList: configItemExpandName,  // 行内所属配置名批量补全
	AfterDetail:   configItemAfterDetail, // 所属配置编码 + 文件类字段拼接访问URL
}

// configItemExpandName 分页行所属配置名批量补全
// 将配置ID映射为配置名称，避免逐行查询
func configItemExpandName(ctx context.Context, list []*v1.ConfigItemItem, req any) error {
	cfgNames, err := (&sConfigItem{}).configNames(ctx)
	if err != nil {
		return err
	}
	for _, item := range list {
		item.ConfigName = cfgNames[item.ConfigId]
	}
	return nil
}

// configItemAfterDetail 详情装配
//
// 补充所属配置编码，并按控件类型拼接文件类字段的访问URL：
//   - image：value 直接拼域名
//   - file：value 为 "原文件名|路径"，路径部分拼域名
//   - images：value 保持相对路径逗号串，valueList 为逐项拼域名后的数组
//   - files：value 保持 "原文件名|路径" 逗号串，valueList 为逐项拼域名后的数组
func configItemAfterDetail(ctx context.Context, item *v1.ConfigItemDetailRes, req any) error {
	code, err := (&sConfigItem{}).configCode(ctx, item.ConfigId)
	if err != nil {
		return err
	}
	item.ConfigCode = code
	switch item.Type {
	case "image":
		item.Value = upload.GetFileUrl(ctx, item.Value)
	case "file":
		item.Value = configItemFileURL(ctx, item.Value)
	case "images":
		item.ValueList = configItemListURLs(ctx, item.Value, false)
	case "files":
		item.ValueList = configItemListURLs(ctx, item.Value, true)
	}
	return nil
}

// configItemSaveValue 保存前按控件类型归一化 value
//
// 文件类字段执行临时转正并剥离域名前缀，入库形态为相对路径：
//   - image：单文件 URL
//   - images：逗号分隔多文件 URL
//   - file：值为 "原文件名|URL"，仅处理 URL 部分
//   - files：逗号分隔 "原文件名|URL" 串，逐项处理 URL 部分
//   - ueditor：富文本，临时转正并将本站URL写为 [IMG_URL] 占位符
//   - 其余类型：原样返回
func configItemSaveValue(ctx context.Context, typ, value string) string {
	switch typ {
	case "image":
		return SaveFileField(ctx, value, "config")
	case "images":
		return SaveMultiFileField(ctx, value, "config")
	case "file":
		return configItemSaveFilePair(ctx, value)
	case "files":
		if value == "" {
			return ""
		}
		parts := strings.Split(value, ",")
		for i, p := range parts {
			parts[i] = configItemSaveFilePair(ctx, strings.TrimSpace(p))
		}
		return strings.Join(parts, ",")
	case "ueditor":
		return LocalizeContentFields(ctx, value, "config")
	}
	return value
}

// configItemSaveFilePair 处理单个 "原文件名|URL" 对
// URL 部分转正并剥离域名，文件名前缀保留
func configItemSaveFilePair(ctx context.Context, value string) string {
	name, url, found := strings.Cut(value, "|")
	if !found {
		return SaveFileField(ctx, value, "config")
	}
	return name + "|" + SaveFileField(ctx, url, "config")
}

// configItemFileURL 单文件值转访问URL
// 将 "原文件名|相对路径" 中的路径部分拼接访问URL
func configItemFileURL(ctx context.Context, value string) string {
	name, path, found := strings.Cut(value, "|")
	if !found {
		return upload.GetFileUrl(ctx, value)
	}
	return name + "|" + upload.GetFileUrl(ctx, path)
}

// configItemListURLs 逗号串逐项拼接访问URL
// keepName 为 true 时保留 "原文件名|" 前缀
func configItemListURLs(ctx context.Context, value string, keepName bool) []string {
	if value == "" {
		return nil
	}
	parts := strings.Split(value, ",")
	list := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		if keepName {
			list = append(list, configItemFileURL(ctx, p))
			continue
		}
		list = append(list, upload.GetFileUrl(ctx, p))
	}
	return list
}

func init() {
	service.RegisterConfigItem(&sConfigItem{})
}

// Page 配置项分页
// 按查询条件分页返回配置项列表
func (s *sConfigItem) Page(ctx context.Context, req *v1.ConfigItemPageReq) (res *v1.ConfigItemPageRes, err error) {
	list, total, err := configItemSpec.Page(ctx, req, req.GetPageNo(), req.GetPageSize())
	if err != nil {
		return nil, err
	}
	return v1.NewPageData(req.GetPageNo(), req.GetPageSize(), total, derefSlice(list)), nil
}

// List 配置项全量列表
// 返回配置项列表并补全所属配置名，精简列表不挂基座钩子，在委托层后处理
func (s *sConfigItem) List(ctx context.Context, req *v1.ConfigItemListReq) (res v1.ConfigItemListRes, err error) {
	list, err := configItemSpec.List(ctx, req)
	if err != nil {
		return nil, err
	}
	cfgNames, err := s.configNames(ctx)
	if err != nil {
		return nil, err
	}
	rows := derefSlice(list)
	for i := range rows {
		rows[i].ConfigName = cfgNames[rows[i].ConfigId]
	}
	return rows, nil
}

// Detail 配置项详情
// 按配置项ID查询详情及其值列表
func (s *sConfigItem) Detail(ctx context.Context, req *v1.ConfigItemDetailReq) (res *v1.ConfigItemDetailRes, err error) {
	return configItemSpec.Detail(ctx, req.Id, req)
}

// Add 新增配置项
// 按控件类型归一化值后落库，返回新记录ID
func (s *sConfigItem) Add(ctx context.Context, req *v1.ConfigItemAddReq) (res *v1.IdRes, err error) {
	result, err := dao.ConfigItem.Ctx(ctx).Data(map[string]any{
		"config_id": req.ConfigId,
		"name":      req.Name,
		"code":      req.Code,
		"value":     configItemSaveValue(ctx, req.Type, req.Value),
		"options":   req.Options,
		"type":      req.Type,
		"status":    req.Status,
		"sort":      req.Sort,
		"note":      req.Note,
	}).Insert()
	if err != nil {
		return nil, err
	}
	idv, _ := result.LastInsertId()
	return &v1.IdRes{Id: idv}, nil
}

// Update 修改配置项
// 按配置项ID更新记录，值按控件类型归一化
func (s *sConfigItem) Update(ctx context.Context, req *v1.ConfigItemUpdateReq) (res *v1.ConfigItemUpdateRes, err error) {
	if _, err = dao.ConfigItem.Ctx(ctx).Where("id", req.Id).Data(map[string]any{
		"config_id": req.ConfigId,
		"name":      req.Name,
		"code":      req.Code,
		"value":     configItemSaveValue(ctx, req.Type, req.Value),
		"options":   req.Options,
		"type":      req.Type,
		"status":    req.Status,
		"sort":      req.Sort,
		"note":      req.Note,
	}).Update(); err != nil {
		return nil, err
	}
	return &v1.ConfigItemUpdateRes{}, nil
}

// Delete 删除配置项
// 按配置项ID删除单条记录
func (s *sConfigItem) Delete(ctx context.Context, req *v1.ConfigItemDeleteReq) (res *v1.ConfigItemDeleteRes, err error) {
	if err = configItemSpec.Delete(ctx, req.Id); err != nil {
		return nil, err
	}
	return &v1.ConfigItemDeleteRes{}, nil
}

// BatchDelete 批量删除配置项
// 按配置项ID列表批量删除记录
func (s *sConfigItem) BatchDelete(ctx context.Context, req *v1.ConfigItemBatchDeleteReq) (res *v1.ConfigItemBatchDeleteRes, err error) {
	if err = configItemSpec.BatchDelete(ctx, req.Ids); err != nil {
		return nil, err
	}
	return &v1.ConfigItemBatchDeleteRes{}, nil
}

// configCode 查询配置ID对应的编码
func (s *sConfigItem) configCode(ctx context.Context, configId int64) (string, error) {
	code, err := dao.Config.Ctx(ctx).Where("id", configId).Fields("code").Value()
	if err != nil {
		return "", err
	}
	return code.String(), nil
}

// configNames 查询配置ID到名称的映射
func (s *sConfigItem) configNames(ctx context.Context) (map[int64]string, error) {
	result := make(map[int64]string)
	rows, err := dao.Config.Ctx(ctx).Fields("id, name").All()
	if err != nil {
		return nil, err
	}
	for _, r := range rows {
		result[r["id"].Int64()] = r["name"].String()
	}
	return result, nil
}
