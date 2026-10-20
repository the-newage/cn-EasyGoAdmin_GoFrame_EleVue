## 📚 项目介绍

一款基于 GoFrame v2、Vue3、TypeScript、ElementPlus 打造的模块化、高性能前后端分离敏捷开发框架，可快速搭建企业级后台管理系统。框架本着简化开发、提升效率的初衷，内置代码生成器，根据已建表结构一键生成整个模块的增删改查代码，并同步创建菜单、权限节点与字典数据，真正实现低代码开发，缩短研发周期。

框架提供体系化的后端能力：多租户数据隔离、RBAC 权限与四档数据权限、JWT 双令牌认证、接口限流、操作日志审计、在线用户/服务器/缓存三组监控、可视化定时任务等开箱即用；并配有全局钩子机制，审计字段、数据软删、租户归属等横切逻辑统一在数据层自动完成，业务开发零负担。

数据库全面支持 MySQL、PostgreSQL、SqlServer、Oracle、Sqlite 五种主流数据库并可灵活切换，配套整库迁移工具，跨数据库迁移一条命令完成；前端提供动态菜单、25+ 公共组件与统一请求层，开箱即用；同时提供 Docker 全栈编排，从开发到部署一条龙。

## 🍻 项目特点

+ 低代码开发：代码生成器读表结构一键产出后端八件套 + 前端整套 CRUD 代码，并同步创建菜单、权限节点、字典数据与路由注册（幂等），支持 HTTP 与 CLI 两种方式
+ 多数据库支持：内置 MySQL、PostgreSQL、SqlServer、Oracle、Sqlite 五方言驱动适配可灵活切换；配套整库迁移工具，源库内省 → 方言 DDL 翻译 → 分批搬运 → 序列修复，跨库迁移一条命令完成
+ 企业级多租户：数据层全局租户隔离（全局钩子改写 + dao.Ctx 唯一查询入口），租户 CRUD、专属域名与平台管理开箱即用
+ 细粒度权限体系：RBAC（用户-角色-菜单）+ 四档数据权限（全部/本部门及以下/本部门/仅本人），多角色自动取最宽松，超管/平台绕过
+ 安全基座：JWT 双令牌（访问+刷新轮换、黑名单登出/强踢）、图形验证码、固定窗口接口限流（认证入口更严阈值）、操作日志审计、演示环境写保护
+ 全局钩子机制：审计字段盖章、软删除改写、租户隔离统一在数据层生效，业务代码零侵入
+ 统一响应与规范：全链路驼峰、统一错误码、分页/筛选三件套定稿，枚举显示名由字典驱动（唯一真相，文案改字典即生效）
+ 现代化前端：Vue3 + TypeScript + Element Plus + Vite，动态菜单路由、封装请求层、25+ 公共组件（表格/表单/上传/富文本/Excel/行政区划等）
+ 完善监控运维：在线用户监控（可强制下线）、服务器监控、缓存监控；Redis 故障自动降级内存，服务不中断
+ 定时任务框架：数据库可视化配置任务，支持 cron/固定频率/仅一次三种模式与 1 并行/单例/放弃三种执行策略，多实例经分布式锁不双跑
+ Docker 容器化：内置多阶段 Dockerfile（静态编译、非 root、健康检查）与 docker-compose 全栈编排（后端+MySQL+Redis），敏感配置环境变量注入
+ 文档与质量保障：140+ 篇开发者文档（架构/实战/部署/规范）、OpenAPI 原生输出（/swagger）、黄金样本接口回归基线

## 🍪 内置模块
+ 用户管理：维护管理系统的用户，常规信息维护与账号设置。
+ 角色管理：角色菜单管理与权限分配、设置角色所拥有的菜单权限。
+ 菜单管理：配置系统菜单、操作权限、按钮权限标识等。
+ 部门管理：管理系统组织架构，对组织架构进行统一管理维护。
+ 职级管理：管理用户担任的职级。
+ 岗位管理：管理用户担任的岗位。
+ 租户管理：多租户开通、配置与停启用，租户级数据全局隔离。
+ 字典管理：对系统中常用的较为固定的数据进行统一维护，枚举文案改字典即生效。
+ 配置管理：对系统的常规配置信息进行维护。
+ 参数管理：系统运行参数的统一维护与读取。
+ 城市管理：统一对全国行政区划进行维护，对其他模块提供行政区划数据支撑。
+ 分类管理：树形分类数据的统一维护，支撑内容模块。
+ 文章管理：文章的发布、维护与富文本内容处理。
+ 通知公告：系统通知公告信息发布维护。
+ 友链管理：对系统友情链接、合作伙伴等相关外链进行集成维护管理。
+ 文件模板管理：富文本/上传等文件模板的统一管理。
+ 定时任务：可视化配置调度任务（cron/固定频率/仅一次），含执行日志查询。
+ 操作日志：系统正常操作日志记录和查询；系统异常信息日志记录和查询。
+ 登录日志：系统登录日志记录查询，包含登录异常。
+ 在线用户监控：当前在线用户查询与强制下线。
+ 服务器监控：运行时、磁盘与数据库/Redis 连通性监测。
+ 缓存监控：Redis 运行状态与内存使用监测（降级内存模式可见）。
+ 个人中心：当前登录用户的个人信息便捷修改。
+ 代码生成：一键生成模块 CRUD 功能，包括后端和前端等相关代码，并同步创建菜单与权限。
+ 案例演示：代码生成器一键生成后的演示案例。

## 👷 软件信息
* 软件名称：EasyGoAdmin敏捷开发框架GoFrame+ElementPlus版本
* 软件作者：@小蚂蚁云 团队
* 软件协议：Apache-2.0
* 官网网址：[https://www.easygoadmin.vip](https://www.easygoadmin.vip)
* 文档网址：[http://docs.goframe.elevue.easygoadmin.vip](http://docs.goframe.elevue.easygoadmin.vip)
* 演示地址：[http://manage.goframe.elevue.easygoadmin.vip](http://manage.goframe.elevue.easygoadmin.vip)

## 📌 版本说明

版本名称 | 版本说明 | 版本地址
---|-------------------------------------|---
GoFrame+Layui混编版 | 采用GoFrame v2、Layui等框架研发 | https://gitee.com/easygoadmin/EasyGoAdmin_GoFrame_Layui
GoFrame+ElementPlus前后端分离版 | 采用GoFrame v2、Vue3、TypeScript、ElementPlus等框架研发前后端分离版本 | https://gitee.com/easygoadmin/EasyGoAdmin_GoFrame_EleVue
GoFrame+AntDesign前后端分离版 | 采用GoFrame v2、Vue3、TypeScript、AntDesign等框架研发前后端分离版本 | https://gitee.com/easygoadmin/EasyGoAdmin_GoFrame_AntdVue
GoFrame+NaiveUI前后端分离版 | 采用GoFrame v2、Vue3、TypeScript、NaiveUI等框架研发前后端分离版本 | https://gitee.com/easygoadmin/EasyGoAdmin_GoFrame_NaiveVue
Gin+Layui混编版 | 采用Gin、Layui等框架研发 | https://gitee.com/easygoadmin/EasyGoAdmin_Gin_Layui
Gin+ElementPlus前后端分离版 | 采用Gin、Vue3、TypeScript、ElementPlus等框架研发前后端分离版本 | https://gitee.com/easygoadmin/EasyGoAdmin_Gin_EleVue
Gin+AntDesign前后端分离版 | 采用Gin、Vue3、TypeScript、AntDesign等框架研发前后端分离版本 | https://gitee.com/easygoadmin/EasyGoAdmin_Gin_AntdVue
Gin+NaiveUI前后端分离版 | 采用Gin、Vue3、TypeScript、NaiveUI等框架研发前后端分离版本 | https://gitee.com/easygoadmin/EasyGoAdmin_Gin_NaiveVue
Beego+Layui混编版 | 采用Beego、Layui等框架研发 | https://gitee.com/easygoadmin/EasyGoAdmin_Beego_Layui
Beego+ElementPlus前后端分离版 | 采用Beego、Vue3、TypeScript、ElementPlus等框架研发前后端分离版本 | https://gitee.com/easygoadmin/EasyGoAdmin_Beego_EleVue
Beego+AntDesign前后端分离版 | 采用Beego、Vue3、TypeScript、AntDesign等框架研发前后端分离版本 | https://gitee.com/easygoadmin/EasyGoAdmin_Beego_AntdVue
Beego+NaiveUI前后端分离版 | 采用Beego、Vue3、TypeScript、NaiveUI等框架研发前后端分离版本 | https://gitee.com/easygoadmin/EasyGoAdmin_Beego_NaiveVue
Iris+Layui混编版 | 采用Iris、Layui等框架研发 | https://gitee.com/easygoadmin/EasyGoAdmin_Iris_Layui
Iris+ElementPlus前后端分离版 | 采用Iris、Vue3、TypeScript、ElementPlus等框架研发前后端分离版本 | https://gitee.com/easygoadmin/EasyGoAdmin_Iris_EleVue
Iris+AntDesign前后端分离版 | 采用Iris、Vue3、TypeScript、AntDesign等框架研发前后端分离版本 | https://gitee.com/easygoadmin/EasyGoAdmin_Iris_AntdVue
Iris+NaiveUI前后端分离版 | 采用Iris、Vue3、TypeScript、NaiveUI等框架研发前后端分离版本 | https://gitee.com/easygoadmin/EasyGoAdmin_Iris_NaiveVue

## 🍪 项目结构

```text
EasyGoAdmin_GoFrame_EleVue/
├── main.go                     # 服务入口：注册五库驱动 + Redis 适配器 → cmd.Run
├── Makefile                    # run/build/test/vet/fmt/gen(gen-dao)/clean
├── api/                        # 全部请求/响应结构体（单包 v1，每模块一文件）
├── cmd/                        # CLI 入口
│   ├── generator/              #   代码生成器（文件+菜单+字典+路由，幂等）
│   └── migrator/               #   整库迁移（源库内省 → 五方言 DDL → 分批搬运）
├── internal/                   # 后端核心代码
│   ├── cmd/                    #   启动编排：cache.Init → 依赖接线 → 路由 → 调度器
│   ├── controller/             #   控制器（单包薄层，仅转发 service）
│   ├── logic/                  #   业务逻辑（单包，每模块一文件，同包 init 注册）
│   ├── service/                #   服务接口层（手维护，勿运行 gf gen service）
│   ├── dao/                    #   数据访问层（gf gen dao 产出，挂钩子 + Ctx 覆写）
│   ├── model/                  #   entity / do / hooks（审计盖章、软删、租户隔离）
│   ├── middleware/             #   中间件链：CORS→Response→RateLimit→Demo→JWT→Tenant→OperationLog
│   ├── router/                 #   路由注册（router.go）+ gmeta 反射绑定（bind.go）
│   ├── task/                   #   定时任务处理器（init 自注册，依赖单向 task→logic）
│   ├── consts/                 #   业务错误码 / 缓存键 / 日志类型枚举
│   └── migrate/                #   整库迁移引擎（内省/翻译/DDL 渲染/搬运）
├── utility/                    # 基础能力包（无业务依赖）
│   ├── jwt/ cache/ upload/     #   双令牌 / 缓存降级 / 两段式上传
│   ├── captcha/ crypto/ excel  #   图形验证码 / 加解密 / Excel 导入导出
│   └── ip2region/ response/ ua #   IP 归属地 / 统一响应 / 设备识别
├── templates/                  # 代码生成器模板（后端 {{ }} 定界；ui* 前端 [[ ]] 定界）
├── manifest/
│   ├── config/config.yaml      # 运行时配置（唯一配置文件，GF_ 环境变量可覆盖）
│   └── data/ip2region.xdb      # IP 归属地库
├── resource/
│   ├── static/                 # 验证码字体等静态资源
│   ├── swagger/                # 本地托管 API 文档页面
│   └── upload/                 # 上传文件静态服务根（/upload/**）
├── hack/                       # gf gen dao 配置、自定义模板、黄金样本回归基线
├── docker/                     # Docker 化文件（构建上下文为项目根）
│   ├── Dockerfile              #   多阶段构建：静态编译 → alpine 运行（非 root + 健康检查）
│   ├── docker-compose.yml      #   三服务编排：后端 + MySQL8 + Redis7
│   └── .env.example            #   compose 变量模板（复制为 .env，密钥不入库）
├── .dockerignore               # Docker 构建上下文排除清单
├── ui/                         # 前端工程（Vue3 + TypeScript + Element Plus + Vite）
├── document/                   # 项目文档（架构定稿、五方言 SQL、wiki）
├── upload/                     # 上传物理目录（file.uploadDir，挂载持久卷）
├── log/                        # 运行日志
└── bin/                        # 编译产物
```

## 🚀 快速开始

**环境要求**：Go ≥ 1.27、Node.js 22.18.0（pnpm 推荐）、MySQL 8.0、Redis

```sh
# 1. 初始化数据库：导入定稿脚本
#    document/mysql/easygoadmin.goframe.elevue.sql → 库 easygoadmin.goframe.elevue

# 2. 启动后端（默认 :8000，配置见 manifest/config/config.yaml）
go run main.go

# 3. 启动前端（默认 :8001，代理转发至 :8000）
cd ui && pnpm install && pnpm dev

# 或 Docker 一条命令拉起全栈（后端 + MySQL + Redis）
cd docker && copy .env.example .env && docker compose up -d --build
```

登录账号 `admin / 123456`。更多启动细节与常见问题见[在线文档·快速入门](http://docs.goframe.elevue.easygoadmin.vip/zh/quickstart/environment.html)。

## 🔧 模块展示


|                                  |                                  |
|----------------------------------|----------------------------------|
| ![模块展示1](./upload/demo/1.png)   | ![模块展示2](./upload/demo/2.png)   |
| ![模块展示3](./upload/demo/3.png)   | ![模块展示4](./upload/demo/4.png)   |
| ![模块展示5](./upload/demo/5.png)   | ![模块展示6](./upload/demo/6.png)   |
| ![模块展示7](./upload/demo/7.png)   | ![模块展示8](./upload/demo/8.png)   |
| ![模块展示9](./upload/demo/9.png)   | ![模块展示10](./upload/demo/10.png) |
| ![模块展示11](./upload/demo/11.png) | ![模块展示12](./upload/demo/12.png) |
| ![模块展示13](./upload/demo/13.png) | ![模块展示14](./upload/demo/14.png) |
| ![模块展示15](./upload/demo/15.png) | ![模块展示16](./upload/demo/16.png) |
| ![模块展示17](./upload/demo/17.png) | ![模块展示18](./upload/demo/18.png) |
| ![模块展示19](./upload/demo/19.png) | ![模块展示20](./upload/demo/20.png) |
| ![模块展示21](./upload/demo/21.png) | ![模块展示22](./upload/demo/22.png) |
| ![模块展示23](./upload/demo/23.png) | ![模块展示24](./upload/demo/24.png) |
| ![模块展示25](./upload/demo/25.png) | ![模块展示26](./upload/demo/26.png) |
| ![模块展示27](./upload/demo/27.png) | ![模块展示28](./upload/demo/28.png) |
| ![模块展示29](./upload/demo/29.png) | ![模块展示30](./upload/demo/30.png) |
| ![模块展示31](./upload/demo/31.png) | ![模块展示32](./upload/demo/32.png) |
| ![模块展示33](./upload/demo/33.png) | ![模块展示34](./upload/demo/34.png) |

## ✨  特别鸣谢
感谢[GoFrame](https://goframe.org/#all-updates)、[Vue.js](https://cn.vuejs.org/)、[Element Plus](https://element-plus.org/zh-CN/)、[TypeScript](https://www.typescriptlang.org/)、[Vite](https://cn.vitejs.dev/)、[Pinia](https://pinia.vuejs.org/)等优秀开源项目。