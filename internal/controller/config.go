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

package controller

import (
	"context"

	"easygoadmin/api/v1"
	"easygoadmin/internal/service"
)

// ============================== 配置控制器 ==============================

// CConfig 配置控制器
// 薄层控制器，仅负责参数透传，业务逻辑由 service.Config() 处理
type CConfig struct{}

// Page 配置分页
// 按查询条件分页返回配置列表
func (c *CConfig) Page(ctx context.Context, req *v1.ConfigPageReq) (res *v1.ConfigPageRes, err error) {
	return service.Config().Page(ctx, req)
}

// List 配置列表
// 返回配置列表，每项携带其下的配置项集合
func (c *CConfig) List(ctx context.Context, req *v1.ConfigListReq) (res v1.ConfigListRes, err error) {
	return service.Config().List(ctx, req)
}

// Detail 配置详情
// 按配置ID查询详情及其下的配置项集合
func (c *CConfig) Detail(ctx context.Context, req *v1.ConfigDetailReq) (res *v1.ConfigDetailRes, err error) {
	return service.Config().Detail(ctx, req)
}

// Add 新增配置
// 新增成功后返回新记录ID
func (c *CConfig) Add(ctx context.Context, req *v1.ConfigAddReq) (res *v1.IdRes, err error) {
	return service.Config().Add(ctx, req)
}

// Update 修改配置
// 按配置ID更新记录
func (c *CConfig) Update(ctx context.Context, req *v1.ConfigUpdateReq) (res *v1.ConfigUpdateRes, err error) {
	return service.Config().Update(ctx, req)
}

// Delete 删除配置
// 按配置ID删除单条记录
func (c *CConfig) Delete(ctx context.Context, req *v1.ConfigDeleteReq) (res *v1.ConfigDeleteRes, err error) {
	return service.Config().Delete(ctx, req)
}

// BatchDelete 批量删除配置
// 按配置ID列表批量删除记录
func (c *CConfig) BatchDelete(ctx context.Context, req *v1.ConfigBatchDeleteReq) (res *v1.ConfigBatchDeleteRes, err error) {
	return service.Config().BatchDelete(ctx, req)
}

// WebIndex 获取网站配置
// 返回配置分组及其下的配置项，供网站配置页渲染
func (c *CConfig) WebIndex(ctx context.Context, req *v1.ConfigWebIndexReq) (res *v1.ConfigWebIndexRes, err error) {
	return service.Config().WebIndex(ctx, req)
}

// WebSave 保存网站配置
// 接收配置项编码到值的平铺映射并批量保存
func (c *CConfig) WebSave(ctx context.Context, req *v1.ConfigWebSaveReq) (res *v1.ConfigWebSaveRes, err error) {
	return service.Config().WebSave(ctx, req)
}
