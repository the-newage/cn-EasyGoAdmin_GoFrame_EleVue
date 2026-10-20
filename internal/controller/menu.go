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

// ============================== 菜单控制器 ==============================

// CMenu 菜单控制器
// 薄层控制器，仅负责参数透传，业务逻辑由 service.Menu() 处理
type CMenu struct{}

// Page 菜单分页
// 按查询条件分页返回平铺的菜单列表
func (c *CMenu) Page(ctx context.Context, req *v1.MenuPageReq) (res *v1.MenuPageRes, err error) {
	return service.Menu().Page(ctx, req)
}

// List 菜单平铺列表
// 返回菜单平铺列表，由前端自行构建树形结构
func (c *CMenu) List(ctx context.Context, req *v1.MenuListReq) (res v1.MenuListRes, err error) {
	return service.Menu().List(ctx, req)
}

// Detail 菜单详情
// 按菜单ID查询详情
func (c *CMenu) Detail(ctx context.Context, req *v1.MenuDetailReq) (res *v1.MenuDetailRes, err error) {
	return service.Menu().Detail(ctx, req)
}

// Add 新增菜单
// 新增成功后返回新记录ID
func (c *CMenu) Add(ctx context.Context, req *v1.MenuAddReq) (res *v1.IdRes, err error) {
	return service.Menu().Add(ctx, req)
}

// Update 修改菜单
// 按菜单ID更新记录
func (c *CMenu) Update(ctx context.Context, req *v1.MenuUpdateReq) (res *v1.MenuUpdateRes, err error) {
	return service.Menu().Update(ctx, req)
}

// Delete 删除菜单
// 按菜单ID删除单条记录
func (c *CMenu) Delete(ctx context.Context, req *v1.MenuDeleteReq) (res *v1.MenuDeleteRes, err error) {
	return service.Menu().Delete(ctx, req)
}

// BatchDelete 批量删除菜单
// 按菜单ID列表批量删除记录
func (c *CMenu) BatchDelete(ctx context.Context, req *v1.MenuBatchDeleteReq) (res *v1.MenuBatchDeleteRes, err error) {
	return service.Menu().BatchDelete(ctx, req)
}
