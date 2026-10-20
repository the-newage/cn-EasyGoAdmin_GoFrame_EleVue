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

// ================================================================================
// 手维护服务接口：监控域无表模块（在线用户/服务器监控/缓存监控）
// ================================================================================

package service

import (
	"context"

	"easygoadmin/api/v1"
)

// IOnlineUser 在线用户服务接口
type IOnlineUser interface {
	// Page 在线用户分页
	Page(ctx context.Context, req *v1.OnlineUserPageReq) (res *v1.OnlineUserPageRes, err error)
	// Detail 在线会话详情（含令牌原文）
	Detail(ctx context.Context, req *v1.OnlineUserDetailReq) (res *v1.OnlineUserDetailRes, err error)
	// ForceLogout 强制下线
	ForceLogout(ctx context.Context, req *v1.OnlineUserForceLogoutReq) (res *v1.OnlineUserForceLogoutRes, err error)
	// ForceLogoutUser 踢出该用户全部设备
	ForceLogoutUser(ctx context.Context, req *v1.OnlineUserForceLogoutUserReq) (res *v1.OnlineUserForceLogoutUserRes, err error)
}

// IServerMonitor 服务器监控服务接口
type IServerMonitor interface {
	// Info 服务器监控信息
	Info(ctx context.Context, req *v1.ServerInfoReq) (res *v1.ServerInfoRes, err error)
}

// ICacheMonitor 缓存监控服务接口
type ICacheMonitor interface {
	// Info 缓存监控信息
	Info(ctx context.Context, req *v1.CacheInfoReq) (res *v1.CacheInfoRes, err error)
}

var (
	localOnlineUser    IOnlineUser
	localServerMonitor IServerMonitor
	localCacheMonitor  ICacheMonitor
)

// OnlineUser 获取在线用户服务实例
func OnlineUser() IOnlineUser {
	if localOnlineUser == nil {
		panic("implement not found for interface IOnlineUser, forgot register?")
	}
	return localOnlineUser
}

// RegisterOnlineUser 注册在线用户服务实现
func RegisterOnlineUser(i IOnlineUser) { localOnlineUser = i }

// ServerMonitor 获取服务器监控服务实例
func ServerMonitor() IServerMonitor {
	if localServerMonitor == nil {
		panic("implement not found for interface IServerMonitor, forgot register?")
	}
	return localServerMonitor
}

// RegisterServerMonitor 注册服务器监控服务实现
func RegisterServerMonitor(i IServerMonitor) { localServerMonitor = i }

// CacheMonitor 获取缓存监控服务实例
func CacheMonitor() ICacheMonitor {
	if localCacheMonitor == nil {
		panic("implement not found for interface ICacheMonitor, forgot register?")
	}
	return localCacheMonitor
}

// RegisterCacheMonitor 注册缓存监控服务实现
func RegisterCacheMonitor(i ICacheMonitor) { localCacheMonitor = i }
