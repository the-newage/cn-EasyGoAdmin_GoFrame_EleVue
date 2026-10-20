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
	"encoding/json"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"

	"easygoadmin/api/v1"
	"easygoadmin/internal/consts"
	"easygoadmin/internal/dao"
	"easygoadmin/internal/service"
	"easygoadmin/utility/cache"
)

// ============================== 字典业务实现 ==============================

// sDict 字典业务实现
// 实现 service.IDict 接口，承载数据字典模块的业务逻辑
type sDict struct{}

// dictSpec 数据字典模块声明
//
// Detail 含字典项组装，采用手写实现
var dictSpec = ModuleSpec[v1.DictItem, v1.DictItem, v1.DictDetailRes]{
	Model:      dao.Dict.Ctx,
	LikeFields: []string{"name", "code"},
	OrderBy:    "sort, id",
}

func init() {
	service.RegisterDict(&sDict{})
}

// Page 数据字典分页
// 按查询条件分页返回字典列表
func (s *sDict) Page(ctx context.Context, req *v1.DictPageReq) (res *v1.DictPageRes, err error) {
	list, total, err := dictSpec.Page(ctx, req, req.GetPageNo(), req.GetPageSize())
	if err != nil {
		return nil, err
	}
	return v1.NewPageData(req.GetPageNo(), req.GetPageSize(), total, derefSlice(list)), nil
}

// Detail 字典详情（含字典项）
// 按字典ID查询详情及其下的字典项集合
func (s *sDict) Detail(ctx context.Context, req *v1.DictDetailReq) (res *v1.DictDetailRes, err error) {
	one, err := dao.Dict.Ctx(ctx).Where("id", req.Id).One()
	if err != nil {
		return nil, err
	}
	if len(one) == 0 {
		return nil, gerror.NewCode(consts.CodeRecordNotFound)
	}
	items, err := s.itemsByDictId(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &v1.DictDetailRes{DictItem: s.toItem(one), Items: items}, nil
}

// Add 新增字典
// 校验字典编码唯一后新增，返回新记录ID
func (s *sDict) Add(ctx context.Context, req *v1.DictAddReq) (res *v1.IdRes, err error) {
	if err = CheckUnique(ctx, dao.Dict.Ctx(ctx), "code", req.Code, 0, "字典编码"); err != nil {
		return nil, err
	}
	result, err := dao.Dict.Ctx(ctx).Data(map[string]any{
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

// Update 修改字典
// 按字典ID更新记录
func (s *sDict) Update(ctx context.Context, req *v1.DictUpdateReq) (res *v1.DictUpdateRes, err error) {
	if _, err = dao.Dict.Ctx(ctx).Where("id", req.Id).Data(map[string]any{
		"name": req.Name,
		"code": req.Code,
		"sort": req.Sort,
		"note": req.Note,
	}).Update(); err != nil {
		return nil, err
	}
	return &v1.DictUpdateRes{}, nil
}

// Delete 删除数据字典（连带字典项）
// 按字典ID删除记录，并连带删除其下字典项
func (s *sDict) Delete(ctx context.Context, req *v1.DictDeleteReq) (res *v1.DictDeleteRes, err error) {
	if err = dictSpec.Delete(ctx, req.Id); err != nil {
		return nil, err
	}
	if _, err = dao.DictItem.Ctx(ctx).Where("dict_id", req.Id).Delete(); err != nil {
		return nil, err
	}
	return &v1.DictDeleteRes{}, nil
}

// BatchDelete 批量删除数据字典（连带字典项）
// 按字典ID列表批量删除记录，并连带删除其下字典项
func (s *sDict) BatchDelete(ctx context.Context, req *v1.DictBatchDeleteReq) (res *v1.DictBatchDeleteRes, err error) {
	if err = dictSpec.BatchDelete(ctx, req.Ids); err != nil {
		return nil, err
	}
	if _, err = dao.DictItem.Ctx(ctx).WhereIn("dict_id", req.Ids).Delete(); err != nil {
		return nil, err
	}
	return &v1.DictBatchDeleteRes{}, nil
}

// RefreshCache 刷新字典缓存
// 重写 Redis 缓存，编码为空时刷新全部字典
func (s *sDict) RefreshCache(ctx context.Context, req *v1.DictRefreshCacheReq) (res *v1.DictRefreshCacheRes, err error) {
	m := dao.DictItem.Ctx(ctx)
	if req.Code != "" {
		dict, e := dao.Dict.Ctx(ctx).Where("code", req.Code).One()
		if e != nil || len(dict) == 0 {
			return nil, gerror.NewCode(consts.CodeRecordNotFound)
		}
		m = m.Where("dict_id", dict["id"].Int64())
	}
	rows, err := m.Order("sort, id").All()
	if err != nil {
		return nil, err
	}
	byCode := make(map[string][]v1.DictItemRow)
	for _, r := range rows {
		dictId := r["dict_id"].Int64()
		code := s.dictCodeById(ctx, dictId)
		byCode[code] = append(byCode[code], v1.DictItemRow{
			Id:    r["id"].Int64(),
			Name:  r["name"].String(),
			Value: r["value"].String(),
			Sort:  r["sort"].Int(),
		})
	}
	for code, items := range byCode {
		s.setCache(ctx, code, items)
	}
	return &v1.DictRefreshCacheRes{}, nil
}

// AliasList 字典别名列表
// 登录用户即可访问，按关键词返回字典列表
func (s *sDict) AliasList(ctx context.Context, req *v1.DictAliasListReq) (res v1.DictAliasListRes, err error) {
	m := dao.Dict.Ctx(ctx)
	if req.Keyword != "" {
		kw := "%" + req.Keyword + "%"
		m = m.WhereLike("name", kw).WhereOrLike("code", kw)
	}
	rows, err := m.Order("sort, id").All()
	if err != nil {
		return nil, err
	}
	list := make([]v1.DictItem, 0, len(rows))
	for _, r := range rows {
		list = append(list, s.toItem(r))
	}
	return list, nil
}

// AliasInfo 字典别名详情
// 登录用户即可访问，按字典编码查询字典及其字典项
func (s *sDict) AliasInfo(ctx context.Context, req *v1.DictAliasInfoReq) (res *v1.DictAliasInfoRes, err error) {
	one, err := dao.Dict.Ctx(ctx).Where("code", req.Code).One()
	if err != nil {
		return nil, err
	}
	if len(one) == 0 {
		return nil, gerror.NewCode(consts.CodeRecordNotFound)
	}
	items, err := s.itemsByDictId(ctx, one["id"].Int64())
	if err != nil {
		return nil, err
	}
	return &v1.DictAliasInfoRes{DictItem: s.toItem(one), Items: items}, nil
}

// ---------- 私有辅助 ----------

// toItem 数据库记录转换为字典行数据
func (s *sDict) toItem(r gdb.Record) v1.DictItem {
	return v1.DictItem{
		Id:         r["id"].Int64(),
		Name:       r["name"].String(),
		Code:       r["code"].String(),
		Sort:       r["sort"].Int(),
		Note:       r["note"].String(),
		CreateUser: r["create_user"].String(),
	}
}

// itemsByDictId 查询单个字典的字典项
func (s *sDict) itemsByDictId(ctx context.Context, dictId int64) ([]v1.DictItemRow, error) {
	rows, err := dao.DictItem.Ctx(ctx).Where("dict_id", dictId).Order("sort, id").All()
	if err != nil {
		return nil, err
	}
	items := make([]v1.DictItemRow, 0, len(rows))
	for _, r := range rows {
		items = append(items, v1.DictItemRow{
			Id:    r["id"].Int64(),
			Name:  r["name"].String(),
			Value: r["value"].String(),
			Sort:  r["sort"].Int(),
		})
	}
	return items, nil
}

// dictCodeById 查询字典ID对应的编码
func (s *sDict) dictCodeById(ctx context.Context, dictId int64) string {
	one, err := dao.Dict.Ctx(ctx).Where("id", dictId).Fields("code").One()
	if err != nil || len(one) == 0 {
		return ""
	}
	return one["code"].String()
}

// setCache 写入字典项缓存
// 采用尽力而为策略，序列化或写入失败时静默跳过
func (s *sDict) setCache(ctx context.Context, code string, items []v1.DictItemRow) {
	if code == "" {
		return
	}
	// 序列化为简单 JSON（避免额外依赖）
	b, err := json.Marshal(items)
	if err != nil {
		return
	}
	_ = cache.Set(ctx, consts.CacheKeyDictItem(code), string(b), 3600)
}
