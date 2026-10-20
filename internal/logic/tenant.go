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
	"time"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/gtime"

	"easygoadmin/api/v1"
	"easygoadmin/internal/consts"
	"easygoadmin/internal/dao"
	"easygoadmin/internal/service"
	"easygoadmin/utility/cache"
	"easygoadmin/utility/crypto"
)

// ============================== 租户业务实现 ==============================

// sTenant 租户业务实现
// 实现 service.ITenant 接口，接口定义见 api/v1/tenant.go
type sTenant struct{}

func init() {
	service.RegisterTenant(&sTenant{})
}

// tenantSerializeMaps 枚举字段到字典编码的映射
// {field}Text 由基座自动填充
var tenantSerializeMaps = map[string]string{
	"status": "tenant_status",
}

// tenantSpec 租户模块声明
//
// Add 含临时编码占位回填，采用手写实现
var tenantSpec = ModuleSpec[v1.TenantItem, v1.TenantListItem, v1.TenantDetailRes]{
	Model:         dao.Tenant.Ctx,
	LikeFields:    []string{"name", "code"}, // querySchemas 字符串输入 → 逐列模糊
	EqFields:      []string{"status"},
	OrderBy:       "id DESC",
	ListOrderBy:   "id",
	SerializeMaps: tenantSerializeMaps,
	Filter:        tenantFilter, // keyword 按入参分支（Page 含 contact）
	FileFields:    []string{"logo"},
	FileSaveDir:   "tenant",
	ToData:        tenantUpdateToData,
	AfterUpdate:   tenantClearCache, // 状态等变更 → 租户信息缓存失效
	BeforeDelete:  tenantDeleteGuard,
	AfterDelete: func(ctx context.Context, id int64) error {
		cache.Del(ctx, consts.CacheKeyTenantInfo(int(id)))
		return nil
	},
	AfterBatchDelete: func(ctx context.Context, ids []int64) error {
		for _, id := range ids {
			cache.Del(ctx, consts.CacheKeyTenantInfo(int(id)))
		}
		return nil
	},
	AfterPageList: tenantFillUserCount, // 行内用户数批量补全
}

// tenantFilter 关键词筛选
//
// 分页请求匹配名称、编码与联系人，列表请求匹配名称与编码
func tenantFilter(ctx context.Context, m *gdb.Model, req any) *gdb.Model {
	switch r := req.(type) {
	case *v1.TenantPageReq:
		if r.Keyword != "" {
			kw := "%" + r.Keyword + "%"
			m = m.Where("(name LIKE ? OR code LIKE ? OR contact LIKE ?)", kw, kw, kw)
		}
	case *v1.TenantListReq:
		if r.Keyword != "" {
			kw := "%" + r.Keyword + "%"
			m = m.Where("(name LIKE ? OR code LIKE ?)", kw, kw)
		}
	}
	return m
}

// tenantUpdateToData 修改请求的表单语义转换
// 编码不可变更，过期时间为空时表示清空
func tenantUpdateToData(req any) map[string]any {
	r, ok := req.(*v1.TenantUpdateReq)
	if !ok {
		panic("tenantUpdateToData: 不支持的入参类型")
	}
	data := map[string]any{
		"name":         r.Name,
		"contact":      r.Contact,
		"mobile":       r.Mobile,
		"email":        r.Email,
		"address":      r.Address,
		"logo":         r.Logo,
		"status":       r.Status,
		"max_user_num": r.MaxUserNum,
		"note":         r.Note,
	}
	if r.ExpireTime != "" {
		if t := gtime.New(r.ExpireTime); t != nil {
			data["expire_time"] = t
		}
	} else {
		data["expire_time"] = nil
	}
	return data
}

// tenantFillUserCount 分页行用户数批量补全
func tenantFillUserCount(ctx context.Context, list []*v1.TenantItem, req any) error {
	counts, err := tenantUserCounts(ctx)
	if err != nil {
		return err
	}
	for _, item := range list {
		item.UserCount = counts[item.Id]
	}
	return nil
}

// tenantClearCache 租户变更后清理缓存
// 失效对应租户的信息缓存
func tenantClearCache(ctx context.Context, id int64, req any) error {
	cache.Del(ctx, consts.CacheKeyTenantInfo(int(id)))
	return nil
}

// tenantDeleteGuard 删除守卫
// 校验记录存在性，存在关联用户时禁止删除
func tenantDeleteGuard(ctx context.Context, id int64) error {
	if n, err := dao.Tenant.Ctx(ctx).Where("id", id).Count(); err != nil {
		return err
	} else if n == 0 {
		return gerror.NewCode(consts.CodeRecordNotFound)
	}
	if n, err := dao.User.Ctx(ctx).Where("tenant_id", id).Count(); err != nil {
		return err
	} else if n > 0 {
		return gerror.NewCode(consts.CodeInUse)
	}
	return nil
}

// ---------- 标准方法：3行委托 ----------

// Page 租户分页
// 分页返回租户列表，并补充关联用户数
func (s *sTenant) Page(ctx context.Context, req *v1.TenantPageReq) (res *v1.TenantPageRes, err error) {
	list, total, err := tenantSpec.Page(ctx, req, req.GetPageNo(), req.GetPageSize())
	if err != nil {
		return nil, err
	}
	return v1.NewPageData(req.GetPageNo(), req.GetPageSize(), total, derefSlice(list)), nil
}

// List 租户全量列表
// 按关键词、状态等条件返回租户精简列表
func (s *sTenant) List(ctx context.Context, req *v1.TenantListReq) (res v1.TenantListRes, err error) {
	list, err := tenantSpec.List(ctx, req)
	if err != nil {
		return nil, err
	}
	return derefSlice(list), nil
}

// Detail 租户详情
// 按租户ID查询详情
func (s *sTenant) Detail(ctx context.Context, req *v1.TenantDetailReq) (res *v1.TenantDetailRes, err error) {
	return tenantSpec.Detail(ctx, req.Id, req)
}

// Update 修改租户
// 按租户ID更新记录，编码不允许变更
func (s *sTenant) Update(ctx context.Context, req *v1.TenantUpdateReq) (res *v1.TenantUpdateRes, err error) {
	if err = tenantSpec.Update(ctx, req.Id, req); err != nil {
		return nil, err
	}
	return &v1.TenantUpdateRes{}, nil
}

// Status 租户状态
// 设置租户的正常或禁用状态，并失效租户信息缓存
func (s *sTenant) Status(ctx context.Context, req *v1.TenantStatusReq) (res *v1.TenantStatusRes, err error) {
	if _, err = dao.Tenant.Ctx(ctx).Where("id", req.Id).Data(map[string]any{"status": req.Status}).Update(); err != nil {
		return nil, err
	}
	cache.Del(ctx, consts.CacheKeyTenantInfo(int(req.Id)))
	return &v1.TenantStatusRes{}, nil
}

// Delete 删除租户
// 按租户ID删除单条记录，存在关联用户时禁止
func (s *sTenant) Delete(ctx context.Context, req *v1.TenantDeleteReq) (res *v1.TenantDeleteRes, err error) {
	if err = tenantSpec.Delete(ctx, req.Id); err != nil {
		return nil, err
	}
	return &v1.TenantDeleteRes{}, nil
}

// BatchDelete 批量删除租户
// 按租户ID列表批量删除记录
func (s *sTenant) BatchDelete(ctx context.Context, req *v1.TenantBatchDeleteReq) (res *v1.TenantBatchDeleteRes, err error) {
	if err = tenantSpec.BatchDelete(ctx, req.Ids); err != nil {
		return nil, err
	}
	return &v1.TenantBatchDeleteRes{}, nil
}

// ---------- 非标准方法（手写实现）----------

// Add 新增租户
//
// 编码留空时自动生成，自增回填前先用临时编码占位
func (s *sTenant) Add(ctx context.Context, req *v1.TenantAddReq) (res *v1.IdRes, err error) {
	code := req.Code
	if code != "" {
		count, e := dao.Tenant.Ctx(ctx).Where("code", code).Count()
		if e != nil {
			return nil, e
		}
		if count > 0 {
			return nil, gerror.NewCodef(consts.CodeRecordExists, "租户编码已存在")
		}
	}
	status := req.Status
	if status == 0 {
		status = 1
	}
	// 自增回填前先用临时编码占位（code NOT NULL + 唯一键），插入后修正为 T%04d(id)
	req.Logo = SaveFileField(ctx, req.Logo, "tenant")
	tempCode := code
	if tempCode == "" {
		tempCode = fmt.Sprintf("T%010d", time.Now().UnixNano()%10000000000)
	}
	data := map[string]any{
		"name":         req.Name,
		"code":         tempCode,
		"contact":      req.Contact,
		"mobile":       req.Mobile,
		"email":        req.Email,
		"address":      req.Address,
		"logo":         req.Logo,
		"status":       status,
		"max_user_num": req.MaxUserNum,
		"note":         req.Note,
	}
	if req.ExpireTime != "" {
		if t := gtime.New(req.ExpireTime); t != nil {
			data["expire_time"] = t
		}
	}
	result, err := dao.Tenant.Ctx(ctx).Data(data).Insert()
	if err != nil {
		return nil, err
	}
	idv, _ := result.LastInsertId()
	if code == "" {
		finalCode := fmt.Sprintf("T%04d", idv)
		if _, err = dao.Tenant.Ctx(ctx).Where("id", idv).Data(map[string]any{"code": finalCode}).Update(); err != nil {
			return nil, err
		}
	}
	return &v1.IdRes{Id: idv}, nil
}

// Account 创建租户账号
//
// 同时创建租户及其初始管理员账号；
// bcrypt 属 CPU 密集操作，放在事务外预计算，保持短事务
func (s *sTenant) Account(ctx context.Context, req *v1.TenantAccountReq) (res *v1.TenantAccountRes, err error) {
	// 租户：自增回填（临时编码占位 → T%04d 修正）
	tempCode := fmt.Sprintf("T%010d", time.Now().UnixNano()%10000000000)
	tData := map[string]any{
		"name":         req.Name,
		"code":         tempCode,
		"contact":      req.Contact,
		"mobile":       req.Mobile,
		"email":        req.Email,
		"address":      req.Address,
		"status":       1,
		"max_user_num": req.MaxUserNum,
		"note":         req.Note,
	}
	if req.ExpireTime != "" {
		if t := gtime.New(req.ExpireTime); t != nil {
			tData["expire_time"] = t
		}
	}
	// bcrypt 属 CPU 密集操作，放在事务外预计算（短事务纪律）
	salt := crypto.GenerateSalt()
	hashed, err := crypto.HashPassword(req.Password, salt)
	if err != nil {
		return nil, err
	}
	realname := req.Realname
	if realname == "" {
		realname = req.Username
	}

	// 租户 + 初始管理员三段写入保持原子性（全部经 txCtx）
	var tenantId, uid int64
	var code string
	if err = WithTrans(ctx, func(txCtx context.Context) error {
		tResult, e := dao.Tenant.Ctx(txCtx).Data(tData).Insert()
		if e != nil {
			return e
		}
		tenantId, _ = tResult.LastInsertId()
		code = fmt.Sprintf("T%04d", tenantId)
		if _, e = dao.Tenant.Ctx(txCtx).Where("id", tenantId).Data(map[string]any{"code": code}).Update(); e != nil {
			return e
		}
		// 初始管理员（tenant_id 由 hooks 自动填充）
		tenantCtx := consts.WithTenantID(txCtx, int(tenantId))
		uResult, e := dao.User.Ctx(tenantCtx).Data(map[string]any{
			"username":      req.Username,
			"code":          fmt.Sprintf("%s-000001", code), // 首个用户序号 000001
			"password":      hashed,
			"salt":          salt,
			"realname":      realname,
			"mobile":        req.Mobile,
			"email":         req.Email,
			"province_code": "",
			"city_code":     "",
			"district_code": "",
			"street_code":   "",
			"address":       "",
			"status":        1,
			"tenant_id":     int(tenantId),
		}).Insert()
		if e != nil {
			return e
		}
		uid, _ = uResult.LastInsertId()
		return nil
	}); err != nil {
		return nil, err
	}
	return &v1.TenantAccountRes{TenantId: tenantId, Code: code, UserId: uid}, nil
}

// Users 查询租户用户列表
//
// 平台上下文可跨租户查看，按 tenant_id 精确查询
func (s *sTenant) Users(ctx context.Context, req *v1.TenantUsersReq) (res *v1.TenantUsersRes, err error) {
	// 平台上下文可跨租户查看；直接按 tenant_id 精确查询（绕过 ApplyCtx 的租户收敛，显式指定）
	rows, err := dao.User.Ctx(ctx).Where("tenant_id", req.TenantId).Order("id").All()
	if err != nil {
		return nil, err
	}
	list := make([]v1.TenantUserItem, 0, len(rows))
	for _, r := range rows {
		list = append(list, v1.TenantUserItem{
			Id:       r["id"].Int64(),
			Username: r["username"].String(),
			Realname: r["realname"].String(),
			Mobile:   r["mobile"].String(),
			Email:    r["email"].String(),
			Status:   r["status"].Int(),
		})
	}
	return &v1.TenantUsersRes{List: list}, nil
}

// ---------- 私有辅助 ----------

// tenantUserCounts 统计租户ID到用户数的映射
func tenantUserCounts(ctx context.Context) (map[int64]int, error) {
	result := make(map[int64]int)
	rows, err := dao.User.Ctx(ctx).Fields("tenant_id").All()
	if err != nil {
		return nil, err
	}
	for _, r := range rows {
		if v := r["tenant_id"].Int64(); v > 0 {
			result[v]++
		}
	}
	return result, nil
}
