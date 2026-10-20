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

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"

	"easygoadmin/api/v1"
	"easygoadmin/internal/consts"
	"easygoadmin/internal/dao"
	"easygoadmin/internal/service"
)

// ============================== 城市业务实现 ==============================

// sCity 城市业务实现
// 实现 service.ICity 接口，承载城市模块的业务逻辑
type sCity struct{}

// citySerializeMaps 枚举字段到字典编码的映射
// {field}Text 由 dict_item 解析，字典管理页可修改文案
var citySerializeMaps = map[string]string{
	"level": "city_level",
}

func init() {
	service.RegisterCity(&sCity{})
}

// citySpec 城市模块声明
//
// 仅 Page 采用声明式；List 含 parentCode、limit 等特有查询，Detail 使用嵌入结构，均手写实现
var citySpec = ModuleSpec[v1.CityItem, v1.CityItem, v1.CityDetailRes]{
	Model:         dao.City.Ctx,
	LikeFields:    []string{"name", "pinyin"},
	EqFields:      []string{"level"},
	OrderBy:       "id",
	SerializeMaps: citySerializeMaps,
	ToData:        cityToData,
}

// cityToData 全量表单语义转换
// 兼容新增与修改请求，将入参转换为落库字段
func cityToData(req any) map[string]any {
	switch r := req.(type) {
	case *v1.CityAddReq:
		return map[string]any{
			"city_code": r.CityCode, "area_code": r.AreaCode, "parent_code": r.ParentCode,
			"zip_code": r.ZipCode, "level": r.Level, "pid": r.Pid, "name": r.Name,
			"short_name": r.ShortName, "full_name": r.FullName, "pinyin": r.Pinyin,
			"lng": r.Lng, "lat": r.Lat,
		}
	case *v1.CityUpdateReq:
		return map[string]any{
			"city_code": r.CityCode, "area_code": r.AreaCode, "parent_code": r.ParentCode,
			"zip_code": r.ZipCode, "level": r.Level, "pid": r.Pid, "name": r.Name,
			"short_name": r.ShortName, "full_name": r.FullName, "pinyin": r.Pinyin,
			"lng": r.Lng, "lat": r.Lat,
		}
	}
	panic("cityToData: 不支持的入参类型")
}

// Page 城市分页
// 按查询条件分页返回城市列表，采用声明式实现
func (s *sCity) Page(ctx context.Context, req *v1.CityPageReq) (res *v1.CityPageRes, err error) {
	list, total, err := citySpec.Page(ctx, req, req.GetPageNo(), req.GetPageSize())
	if err != nil {
		return nil, err
	}
	return v1.NewPageData(req.GetPageNo(), req.GetPageSize(), total, derefSlice(list)), nil
}

// List 城市列表
// 按父级编码、层级、关键词等条件返回城市列表，供级联选择器使用
func (s *sCity) List(ctx context.Context, req *v1.CityListReq) (res v1.CityListRes, err error) {
	m := dao.City.Ctx(ctx)
	if req.ParentCode != "" {
		m = m.Where("parent_code", req.ParentCode)
	} else if req.Level > 0 {
		m = m.Where("level", req.Level)
	}
	if req.Keyword != "" {
		kw := "%" + req.Keyword + "%"
		m = m.WhereLike("name", kw).WhereOrLike("pinyin", kw)
	}
	if req.Limit > 0 {
		m = m.Limit(req.Limit)
	}
	rows, err := m.Order("id").All()
	if err != nil {
		return nil, err
	}
	list := make([]v1.CityItem, 0, len(rows))
	for _, r := range rows {
		list = append(list, s.toItem(ctx, r))
	}
	return list, nil
}

// Detail 城市详情
// 按城市ID查询详情，并补充父级名称
func (s *sCity) Detail(ctx context.Context, req *v1.CityDetailReq) (res *v1.CityDetailRes, err error) {
	one, err := dao.City.Ctx(ctx).Where("id", req.Id).One()
	if err != nil {
		return nil, err
	}
	if len(one) == 0 {
		return nil, gerror.NewCode(consts.CodeRecordNotFound)
	}
	res = &v1.CityDetailRes{CityItem: s.toItem(ctx, one)}
	if pid := one["pid"].Int64(); pid > 0 {
		if parent, e := dao.City.Ctx(ctx).Where("id", pid).Fields("name").One(); e == nil && len(parent) > 0 {
			res.ParentName = parent["name"].String()
		}
	}
	return res, nil
}

// Add 新增城市
// 新增成功后返回新记录ID，采用声明式实现
func (s *sCity) Add(ctx context.Context, req *v1.CityAddReq) (res *v1.IdRes, err error) {
	id, err := citySpec.Add(ctx, req)
	if err != nil {
		return nil, err
	}
	return &v1.IdRes{Id: id}, nil
}

// Update 修改城市
// 按城市ID更新记录，采用声明式实现
func (s *sCity) Update(ctx context.Context, req *v1.CityUpdateReq) (res *v1.CityUpdateRes, err error) {
	if err = citySpec.Update(ctx, req.Id, req); err != nil {
		return nil, err
	}
	return &v1.CityUpdateRes{}, nil
}

// Delete 删除城市
// 按城市ID删除单条记录，采用声明式实现
func (s *sCity) Delete(ctx context.Context, req *v1.CityDeleteReq) (res *v1.CityDeleteRes, err error) {
	if err = citySpec.Delete(ctx, req.Id); err != nil {
		return nil, err
	}
	return &v1.CityDeleteRes{}, nil
}

// BatchDelete 批量删除城市
// 按城市ID列表批量删除记录，采用声明式实现
func (s *sCity) BatchDelete(ctx context.Context, req *v1.CityBatchDeleteReq) (res *v1.CityBatchDeleteRes, err error) {
	if err = citySpec.BatchDelete(ctx, req.Ids); err != nil {
		return nil, err
	}
	return &v1.CityBatchDeleteRes{}, nil
}

// ByParent 按父编码查询下级城市
// 登录用户即可访问，返回指定父级下的直属下级城市
func (s *sCity) ByParent(ctx context.Context, req *v1.CityByParentReq) (res *v1.CityByParentRes, err error) {
	rows, err := dao.City.Ctx(ctx).Where("parent_code", req.ParentCode).Order("id").All()
	if err != nil {
		return nil, err
	}
	list := make([]v1.CityItem, 0, len(rows))
	for _, r := range rows {
		list = append(list, s.toItem(ctx, r))
	}
	return &v1.CityByParentRes{List: list}, nil
}

// GetCityList 按父编码获取下级城市
// 登录用户即可访问，供级联选择器懒加载；parentCode 为 0 时返回顶级
func (s *sCity) GetCityList(ctx context.Context, req *v1.CityGetListReq) (res v1.CityGetListRes, err error) {
	rows, err := dao.City.Ctx(ctx).Where("parent_code", req.ParentCode).Order("id").All()
	if err != nil {
		return nil, err
	}
	list := make([]v1.CityItem, 0, len(rows))
	for _, r := range rows {
		list = append(list, s.toItem(ctx, r))
	}
	return list, nil
}

// AreaProvince 省份列表
// 登录用户即可访问，返回层级为省的顶级城市列表
func (s *sCity) AreaProvince(ctx context.Context, req *v1.AreaParentReq) (res *v1.AreaParentRes, err error) {
	rows, err := dao.City.Ctx(ctx).Where("level", 1).Order("id").All()
	if err != nil {
		return nil, err
	}
	list := make([]v1.CityItem, 0, len(rows))
	for _, r := range rows {
		list = append(list, s.toItem(ctx, r))
	}
	return &v1.AreaParentRes{List: list}, nil
}

// ============================================================
// 城市编解码（名称 ↔ 编码互转）
// ============================================================

// ResolveCityNames 城市编码数组转名称串
//
// 按 area_code 查询 name，保留输入顺序，空编码跳过。
// 例：["320000000000","320100000000","320105000000"] 转换为 "江苏省,南京市,建邺区"
func ResolveCityNames(ctx context.Context, codes []string, separator string) string {
	clean := make([]string, 0, len(codes))
	for _, c := range codes {
		if c != "" {
			clean = append(clean, c)
		}
	}
	if len(clean) == 0 {
		return ""
	}
	rows, err := dao.City.Ctx(ctx).WhereIn("area_code", clean).Fields("area_code, name").All()
	if err != nil {
		return ""
	}
	nameByCode := make(map[string]string, len(rows))
	for _, r := range rows {
		nameByCode[r["area_code"].String()] = r["name"].String()
	}
	names := make([]string, 0, len(clean))
	for _, code := range clean {
		if n := nameByCode[code]; n != "" {
			names = append(names, n)
		}
	}
	if separator == "" {
		separator = ","
	}
	return strings.Join(names, separator)
}

// ResolveCityCodes 名称串转四级编码数组
//
// 逐级以 parent 约束避免同名歧义。
// 例："江苏省,南京市,建邺区" 转换为 ["320000000000","320100000000","320105000000",""]
func ResolveCityCodes(ctx context.Context, cityInfo string) [4]string {
	var codes [4]string
	cityInfo = strings.TrimSpace(cityInfo)
	if cityInfo == "" {
		return codes
	}
	names := make([]string, 0, 4)
	for _, part := range strings.Split(cityInfo, ",") {
		if p := strings.TrimSpace(part); p != "" {
			names = append(names, p)
		}
	}
	if len(names) == 0 {
		return codes
	}
	// 第一级：省（pid = 0）
	prov, err := dao.City.Ctx(ctx).Where("pid", 0).Where("name", names[0]).One()
	if err != nil || len(prov) == 0 {
		return codes
	}
	codes[0] = prov["area_code"].String()
	// 第二级：市（parent_code = 省 area_code）
	if len(names) > 1 {
		if city, e := dao.City.Ctx(ctx).Where("parent_code", codes[0]).Where("name", names[1]).One(); e == nil && len(city) > 0 {
			codes[1] = city["area_code"].String()
		}
	}
	// 第三级：区（parent_code = 市 area_code）
	if len(names) > 2 && codes[1] != "" {
		if dist, e := dao.City.Ctx(ctx).Where("parent_code", codes[1]).Where("name", names[2]).One(); e == nil && len(dist) > 0 {
			codes[2] = dist["area_code"].String()
		}
	}
	// 第四级：街道（parent_code = 区 area_code）
	if len(names) > 3 && codes[2] != "" {
		if street, e := dao.City.Ctx(ctx).Where("parent_code", codes[2]).Where("name", names[3]).One(); e == nil && len(street) > 0 {
			codes[3] = street["area_code"].String()
		}
	}
	return codes
}

// toItem 数据库记录转换为城市行数据
// 枚举层级同步补全层级文本
func (s *sCity) toItem(ctx context.Context, r gdb.Record) v1.CityItem {
	return v1.CityItem{
		Id:         r["id"].Int64(),
		CityCode:   r["city_code"].String(),
		AreaCode:   r["area_code"].String(),
		ParentCode: r["parent_code"].String(),
		ZipCode:    r["zip_code"].String(),
		Level:      r["level"].Int(),
		LevelText:  DictGetText(ctx, citySerializeMaps["level"], r["level"].Int()),
		Pid:        r["pid"].Int64(),
		Name:       r["name"].String(),
		ShortName:  r["short_name"].String(),
		FullName:   r["full_name"].String(),
		Pinyin:     r["pinyin"].String(),
		Lng:        r["lng"].String(),
		Lat:        r["lat"].String(),
	}
}
