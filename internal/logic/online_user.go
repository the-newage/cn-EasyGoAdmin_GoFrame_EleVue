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
	"math"
	"sort"
	"strings"
	"time"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"

	"easygoadmin/api/v1"
	"easygoadmin/internal/consts"
	"easygoadmin/internal/dao"
	"easygoadmin/internal/service"
	"easygoadmin/utility/cache"
	"easygoadmin/utility/ip2region"
	jwtutil "easygoadmin/utility/jwt"
	"easygoadmin/utility/ua"
	"easygoadmin/utility/upload"
)

// ============================== 在线用户业务实现 ==============================

// sOnlineUser 在线用户业务实现
// 数据源为在线令牌集合（登录登记、登出/强退移除），非数据库表
type sOnlineUser struct{}

func init() {
	service.RegisterOnlineUser(&sOnlineUser{})
}

// onlineInfo 在线令牌详情（JSON 缓存）
type onlineInfo struct {
	Jti             string `json:"jti"`
	UserId          int    `json:"userId"`
	Username        string `json:"username"`
	Realname        string `json:"realname"`
	Ip              string `json:"ip"`
	Location        string `json:"location"`
	Oss             string `json:"oss"`     // 操作系统
	Browser         string `json:"browser"` // 浏览器
	LoginTime       string `json:"loginTime"`
	IssuedAt        int64  `json:"issuedAt"`        // 访问令牌颁发时间戳（秒）
	ExpireAt        int64  `json:"expireAt"`        // 访问令牌过期时间戳（秒）
	RefreshJti      string `json:"refreshJti"`      // 刷新令牌 jti
	RefreshIssuedAt int64  `json:"refreshIssuedAt"` // 刷新令牌颁发时间戳（秒）
	RefreshExpireAt int64  `json:"refreshExpireAt"` // 刷新令牌过期时间戳（秒）
	AccessToken     string `json:"accessToken"`     // 访问令牌原文（详情查看）
	RefreshToken    string `json:"refreshToken"`    // 刷新令牌原文（详情查看）
	TokenType       string `json:"tokenType"`       // 令牌类型（Bearer）
	ExpiresIn       int64  `json:"expiresIn"`       // 访问令牌有效期（秒）
}

// MarkOnline 登录成功后登记在线令牌
// zset score=访问令牌过期时间戳用于区间过滤，详情键携带展示字段；
// 刷新令牌仅登记标识与有效期（不入 zset，列表以访问令牌为准）
func MarkOnline(ctx context.Context, access, refresh *jwtutil.Claims, accessToken, refreshToken string, userId int, username, realname, ip, userAgent string) {
	if access == nil || access.ID == "" {
		return
	}
	now := time.Now()
	expireAt := now.Add(jwtutil.ExpireDuration(ctx, jwtutil.TokenTypeAccess))
	if access.ExpiresAt != nil {
		expireAt = access.ExpiresAt.Time
	}
	info := onlineInfo{
		Jti:       access.ID,
		UserId:    userId,
		Username:  username,
		Realname:  realname,
		Ip:        ip,
		Location:  ip2region.Location(ip),
		Oss:       ua.ParseOS(userAgent),
		Browser:   ua.ParseBrowser(userAgent),
		LoginTime: now.Format("2006-01-02 15:04:05"),
		IssuedAt:  now.Unix(),
		ExpireAt:  expireAt.Unix(),
		RefreshJti: func() string {
			if refresh != nil {
				return refresh.ID
			}
			return ""
		}(),
		RefreshIssuedAt: now.Unix(),
		RefreshExpireAt: func() int64 {
			if refresh != nil && refresh.ExpiresAt != nil {
				return refresh.ExpiresAt.Time.Unix()
			}
			return 0
		}(),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    int64(jwtutil.ExpireDuration(ctx, jwtutil.TokenTypeAccess).Seconds()),
	}
	if b, err := json.Marshal(info); err == nil {
		ttl := int(expireAt.Sub(now).Seconds()) + 60
		_ = cache.Set(ctx, consts.CacheKeyOnlineInfo(access.ID), string(b), ttl)
	}
	_ = cache.ZAdd(ctx, consts.CacheKeyOnlineJtiSet(), access.ID, float64(expireAt.Unix()))
}

// MarkOffline 登出或强退后移除在线登记
func MarkOffline(ctx context.Context, jti string) {
	if jti == "" {
		return
	}
	_ = cache.ZRem(ctx, consts.CacheKeyOnlineJtiSet(), jti)
	_ = cache.Del(ctx, consts.CacheKeyOnlineInfo(jti))
}

// Page 在线用户分页
// 取未过期令牌后按用户名/IP 过滤，内存分页，按登录时间倒序
func (s *sOnlineUser) Page(ctx context.Context, req *v1.OnlineUserPageReq) (res *v1.OnlineUserPageRes, err error) {
	jtis, err := cache.ZRangeByScore(ctx, consts.CacheKeyOnlineJtiSet(), float64(time.Now().Unix()), math.Inf(1))
	if err != nil {
		return nil, err
	}
	items := make([]v1.OnlineUserItem, 0, len(jtis))
	for _, jti := range jtis {
		raw, ok := cache.Get(ctx, consts.CacheKeyOnlineInfo(jti))
		if !ok {
			continue
		}
		var info onlineInfo
		if json.Unmarshal([]byte(raw), &info) != nil {
			continue
		}
		if req.Username != "" && !strings.Contains(info.Username, req.Username) {
			continue
		}
		if req.Ip != "" && !strings.Contains(info.Ip, req.Ip) {
			continue
		}
		items = append(items, v1.OnlineUserItem{
			Id:         info.Jti,
			UserId:     info.UserId,
			Username:   info.Username,
			Realname:   info.Realname,
			Ip:         info.Ip,
			Location:   info.Location,
			Oss:        info.Oss,
			Browser:    info.Browser,
			LoginTime:  info.LoginTime,
			IssuedTime: time.Unix(info.IssuedAt, 0).Format("2006-01-02 15:04:05"),
			ExpireTime: time.Unix(info.ExpireAt, 0).Format("2006-01-02 15:04:05"),
			RefreshIssuedTime: func() string {
				if info.RefreshIssuedAt <= 0 {
					return ""
				}
				return time.Unix(info.RefreshIssuedAt, 0).Format("2006-01-02 15:04:05")
			}(),
			RefreshExpireTime: func() string {
				if info.RefreshExpireAt <= 0 {
					return ""
				}
				return time.Unix(info.RefreshExpireAt, 0).Format("2006-01-02 15:04:05")
			}(),
			TokenType: info.TokenType,
			ExpiresIn: info.ExpiresIn,
		})
	}
	sort.Slice(items, func(i, j int) bool { return items[i].LoginTime > items[j].LoginTime })

	// 同用户多会话标注（序号/总数，按登录时间倒序内的序号）
	if len(items) > 1 {
		totalByUser := map[int]int{}
		seqByUser := map[int]int{}
		for _, it := range items {
			totalByUser[it.UserId]++
		}
		for i := range items {
			seqByUser[items[i].UserId]++
			items[i].SessionNo = seqByUser[items[i].UserId]
			items[i].Sessions = totalByUser[items[i].UserId]
		}
	}

	// 按用户ID批量补齐头像（与用户列表同源，头像变更即时生效）
	if len(items) > 0 {
		if err = fillOnlineAvatars(ctx, items); err != nil {
			return nil, err
		}
	}

	total := int64(len(items))
	pageNo, pageSize := req.GetPageNo(), req.GetPageSize()
	start := int((pageNo - 1) * pageSize)
	if start > len(items) {
		start = len(items)
	}
	end := start + int(pageSize)
	if end > len(items) {
		end = len(items)
	}
	return v1.NewPageData(pageNo, pageSize, int(total), items[start:end]), nil
}

// fillOnlineAvatars 按用户ID批量查询并回填头像
// 查询用户表 avatar 相对路径，经 upload.GetFileUrl 补全访问域名；
// 单个用户多会话只查一次
func fillOnlineAvatars(ctx context.Context, items []v1.OnlineUserItem) error {
	idMap := map[int]bool{}
	ids := make([]int, 0, len(items))
	for _, it := range items {
		if it.UserId > 0 && !idMap[it.UserId] {
			idMap[it.UserId] = true
			ids = append(ids, it.UserId)
		}
	}
	if len(ids) == 0 {
		return nil
	}
	rows, err := dao.User.Ctx(ctx).WhereIn("id", ids).Fields("id, avatar").All()
	if err != nil {
		return err
	}
	avatarMap := make(map[int]string, len(rows))
	for _, r := range rows {
		avatarMap[r["id"].Int()] = upload.GetFileUrl(ctx, r["avatar"].String())
	}
	for i := range items {
		items[i].Avatar = avatarMap[items[i].UserId]
	}
	return nil
}

// Detail 在线会话详情
// 按访问令牌 jti 读取在线登记，返回含令牌原文的完整信息
func (s *sOnlineUser) Detail(ctx context.Context, req *v1.OnlineUserDetailReq) (res *v1.OnlineUserDetailRes, err error) {
	raw, ok := cache.Get(ctx, consts.CacheKeyOnlineInfo(req.Id))
	if !ok {
		return nil, gerror.NewCode(consts.CodeRecordNotFound)
	}
	var info onlineInfo
	if err = json.Unmarshal([]byte(raw), &info); err != nil {
		return nil, gerror.NewCode(consts.CodeRecordNotFound)
	}
	item := v1.OnlineUserItem{
		Id:         info.Jti,
		UserId:     info.UserId,
		Username:   info.Username,
		Realname:   info.Realname,
		Ip:         info.Ip,
		Location:   info.Location,
		Oss:        info.Oss,
		Browser:    info.Browser,
		LoginTime:  info.LoginTime,
		IssuedTime: time.Unix(info.IssuedAt, 0).Format("2006-01-02 15:04:05"),
		ExpireTime: time.Unix(info.ExpireAt, 0).Format("2006-01-02 15:04:05"),
		RefreshIssuedTime: func() string {
			if info.RefreshIssuedAt <= 0 {
				return ""
			}
			return time.Unix(info.RefreshIssuedAt, 0).Format("2006-01-02 15:04:05")
		}(),
		RefreshExpireTime: func() string {
			if info.RefreshExpireAt <= 0 {
				return ""
			}
			return time.Unix(info.RefreshExpireAt, 0).Format("2006-01-02 15:04:05")
		}(),
		TokenType:    info.TokenType,
		ExpiresIn:    info.ExpiresIn,
		AccessToken:  info.AccessToken,
		RefreshToken: info.RefreshToken,
	}
	// 切片字面量为副本，经列表回填后再取回
	list := []v1.OnlineUserItem{item}
	if err = fillOnlineAvatars(ctx, list); err != nil {
		return nil, err
	}
	return &v1.OnlineUserDetailRes{OnlineUserItem: list[0]}, nil
}

// ForceLogout 强制下线
// 将访问令牌拉黑并移除在线登记
func (s *sOnlineUser) ForceLogout(ctx context.Context, req *v1.OnlineUserForceLogoutReq) (res *v1.OnlineUserForceLogoutRes, err error) {
	raw, ok := cache.Get(ctx, consts.CacheKeyOnlineInfo(req.Id))
	if !ok {
		return nil, gerror.NewCode(consts.CodeRecordNotFound)
	}
	var info onlineInfo
	if err = json.Unmarshal([]byte(raw), &info); err != nil {
		return nil, gerror.NewCode(consts.CodeRecordNotFound)
	}
	ttl := int(time.Until(time.Unix(info.ExpireAt, 0)).Seconds()) + 60
	if ttl > 0 {
		_ = cache.Set(ctx, consts.CacheKeyTokenBlacklist(info.Jti), "1", ttl)
	}
	// 刷新令牌一并拉黑，防止被踢会话经 refresh 换新令牌重新上线
	if info.RefreshJti != "" && info.RefreshExpireAt > 0 {
		if rttl := int(time.Until(time.Unix(info.RefreshExpireAt, 0)).Seconds()) + 60; rttl > 0 {
			_ = cache.Set(ctx, consts.CacheKeyTokenBlacklist(info.RefreshJti), "1", rttl)
		}
	}
	MarkOffline(ctx, info.Jti)
	return &v1.OnlineUserForceLogoutRes{}, nil
}

// ForceLogoutUser 踢出该用户全部设备
// 遍历在线集合，将该用户所有会话的访问/刷新令牌拉黑并移除登记
func (s *sOnlineUser) ForceLogoutUser(ctx context.Context, req *v1.OnlineUserForceLogoutUserReq) (res *v1.OnlineUserForceLogoutUserRes, err error) {
	jtis, err := cache.ZRangeByScore(ctx, consts.CacheKeyOnlineJtiSet(), float64(time.Now().Unix()), math.Inf(1))
	if err != nil {
		return nil, err
	}
	kicked := 0
	for _, jti := range jtis {
		raw, ok := cache.Get(ctx, consts.CacheKeyOnlineInfo(jti))
		if !ok {
			continue
		}
		var info onlineInfo
		if json.Unmarshal([]byte(raw), &info) != nil || info.UserId != req.UserId {
			continue
		}
		if ttl := int(time.Until(time.Unix(info.ExpireAt, 0)).Seconds()) + 60; ttl > 0 {
			_ = cache.Set(ctx, consts.CacheKeyTokenBlacklist(info.Jti), "1", ttl)
		}
		if info.RefreshJti != "" && info.RefreshExpireAt > 0 {
			if rttl := int(time.Until(time.Unix(info.RefreshExpireAt, 0)).Seconds()) + 60; rttl > 0 {
				_ = cache.Set(ctx, consts.CacheKeyTokenBlacklist(info.RefreshJti), "1", rttl)
			}
		}
		MarkOffline(ctx, info.Jti)
		kicked++
	}
	g.Log().Infof(ctx, "踢出用户 %d 全部设备，共 %d 个会话", req.UserId, kicked)
	return &v1.OnlineUserForceLogoutUserRes{}, nil
}
