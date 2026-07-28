# 系统架构

本文档描述这套系统（基于 gin-vue-admin 二次开发）的整体架构，重点是专家库模块（`ExpertDatabase`）如何嵌入到既有框架里，方便后续按同样的套路加新功能。

## 技术栈

- 后端：Go + Gin + GORM + MySQL，鉴权用 JWT + Casbin（RBAC）
- 前端：Vue3（`<script setup>`）+ Element-Plus + Pinia + Vite，路由由后端菜单动态生成
- 权限模型：JWT 认证 + Casbin 接口级鉴权 + 自建的"单位"数据域隔离（不是 gin-vue-admin 自带的 `DataAuthorityId`）

## 后端目录结构

```
server/
  api/v1/<Module>/       控制器层：解析请求、调 service、包装响应
  service/<Module>/      业务逻辑层：真正的读写数据库逻辑
  model/<Module>/        GORM 数据模型 + request/、response/ 子包（查询参数、返回结构体）
  router/<Module>/       路由注册：把 api 层的方法绑到具体路径+方法上
  middleware/            JWTAuth、CasbinHandler 等全局中间件
  initialize/            启动时的装配逻辑：gorm.go（建表）、router.go（挂载所有模块路由）
  config/                config.yaml 对应的 Go 结构体定义
  global/                全局单例：GVA_DB、GVA_CONFIG、GVA_LOG
  source/                种子数据
  main.go, config.yaml
```

每个业务模块（`ExpertDatabase`、`EngineeringEducationDatabase`、系统自带的 `system` 等）都是四层结构的一个"横切面"，比如专家主档一张表对应：

```
model/ExpertDatabase/expert_profile.go              专家主档的 GORM 结构体
model/ExpertDatabase/request/expert_profile.go       列表查询参数结构体
service/ExpertDatabase/expert_profile.go             CRUD 与业务逻辑
api/v1/ExpertDatabase/expert_profile.go               HTTP handler
router/ExpertDatabase/expert_profile.go                路由注册
```

每层还有一个 `enter.go`，把这个模块下所有表的 service/api/router 聚合成一个结构体（如 `ExpertDatabaseServiceGroup`），最终在 `server/initialize/router.go` 里被显式调用挂载——**新模块不是自动发现的，必须手动在 `router.go` 里加一行**。

路由分两组：

- `PublicGroup`：健康检查、登录、初始化，不过中间件
- `PrivateGroup`：`.Use(middleware.JWTAuth()).Use(middleware.CasbinHandler())`，所有业务接口都在这组下面

## 鉴权链路

请求进来后依次经过：

1. **`JWTAuth()`**（`middleware/jwt.go`）：读 `x-token` 请求头，校验签名和黑名单（`sys_jwt_blacklist` 表，用于登出失效），把用户信息（`ID`、`AuthorityId` 等）塞进 `gin.Context`
2. **`CasbinHandler()`**（`middleware/casbin_rbac.go`）：从 context 里取 `AuthorityId` 作为 `sub`，请求路径作为 `obj`，HTTP 方法作为 `act`，调用 `casbinService.Casbin().Enforce(sub, obj, act)` 判断放行——底层是 `*casbin.SyncedCachedEnforcer`，策略存在 `casbin_rule` 表，**内存里有缓存，直接改数据库不会立即生效，需要重启后端**（这是本项目开发过程中踩过好几次的坑）
3. 角色/菜单关系：`sys_authorities`（角色）、`sys_base_menus`（菜单/路由）、`sys_authority_menus`（角色-菜单关联），都在 `server/model/system/` 下

专家库额外做了一层 Casbin 之外的数据域隔离：`sys_organizations`（单位表）+ `sys_users.org_id`（用户所属单位）+ `expert_profile.org_id`（专家档案的申报单位），审核流程里"单位审核员只能看/审自己单位提交的记录"就是靠这个字段手动过滤实现的，不是靠 gin-vue-admin 自带的角色数据权限（那个是角色维度，不是单位维度，本项目明确没用它）。

## 数据库

没有迁移文件，`server/initialize/gorm.go` 在启动时对一份显式的结构体列表跑 `AutoMigrate`，新加表/改字段直接改 GORM 结构体标签，重启后端即可生效。

## 前端目录结构

```
web/src/
  api/            每个资源一个文件，封装 axios 请求（如 api/expertProfile.js）
  view/           页面组件，一个业务模块一个目录（如 view/beijingExpertDatabase/）
  router/         静态路由（登录页、404 等）
  pinia/modules/  user.js（登录态/token）、router.js（动态路由存储）、dictionary.js（数据字典缓存）
  permission.js   全局路由守卫
  utils/          request.js（axios 实例，自动挂 x-token 头）、format.js 等
```

路由是**登录时动态生成**的：登录成功后前端调 `/menu/getMenu` 拿当前角色能看到的菜单树，`permission.js` 里据此调用 `router.addRoute()` 把这些路由注册进 vue-router，存进 `pinia/modules/router.js` 的 `asyncRouters`。这意味着：

- 新加的路由必须先在 `sys_base_menus` 表里有一条记录，否则前端直接 404（哪怕组件文件本身没问题）
- 改了菜单/角色权限后，已登录的浏览器标签页不会自动感知，**必须重新登录**才能刷新路由表
- 不想出现在侧边栏、但又要能被路由到的页面（比如带参数的详情页），把菜单记录的 `hidden` 设成 1 即可，参照 `expertProfileDetail/:id` 这条

## 专家库模块（ExpertDatabase）业务设计

### 数据模型

- `expert_profile`：专家主档，背景信息 + 学科画像 + 审核状态机字段（`status`/`org_id`/`submitted_by`）+ 打分缓存字段（`achievement_score`/`influence_score`/`social_score`/`composite_score`）打平在一张表里
- `expert_achievement`：研究成果，`expert_id` 外键 + `level`（打分用）
- `expert_adoption_record`：决策影响记录，`expert_id` 外键 + `adopting_unit_level`（打分用）
- `expert_academic_position`：学术兼职
- `expert_tag` / `expert_tag_relation`：标签库，多对多
- `expert_approval_log`：审核流转日志，每次状态变更追加一条，不可变

### 审核状态机

```
draft ──提交──▶ pending_org_review ──单位通过──▶ pending_city_review ──市级通过──▶ published
  ▲                    │                              │
  └────────单位退回─────┘              └──────市级退回──────┘（回到 draft）
```

例外：档案没有 `org_id`（所属单位没配置/匹配不上）时，提交会跳过 `pending_org_review`，直接进 `pending_city_review`——否则会卡死在一个没有任何单位审核员能看到的状态里，详见 `expert_approval.go` 里 `Submit()` 的实现。

审核相关的数据可见性规则（不是 Casbin 能表达的，都是手写在 service 层的 SQL 过滤）：

- 单位审核员（角色 9002）：`专家主档` 列表只能看到"已发布"+"本单位已提交的非草稿记录"；`审核台`只显示"待本单位审核"区块
- 市级审核员（角色 9003）：`专家主档` 只能看到"已发布"+"待市级审核"；`审核台`只显示"待市级审核"区块
- 超级管理员（角色 1）：不受上述限制，两边都能看、都能审——`assertSameOrg`（单位匹配校验）和 `GetPendingOrgReview` 都对 `authority_id == 1` 做了放行

### 打分算法

`server/service/ExpertDatabase/expert_score.go`：成果分/决策影响分是"级别权重 × 该级别篇数"求和，职称分是查表映射，社会贡献分是学术兼职篇数 + 荣誉称号加成，综合排序分是四项按可配置权重线性加权。权重全部走数据字典（`sys_dictionaries`，类型如 `expert_achievement_level`），字典里查不到就退回代码里写死的默认值，管理员可以在"系统工具 → 字典管理"里改权重不用发版。只有 `status = published` 的档案会被重算（`recomputeIfPublished`），审核通过瞬间触发一次，另外有个兜底的批量重算方法供定时任务调用。

### 前端页面

- `expertProfile.vue`：专家主档列表，CRUD + 直接提交审核（不用跳去审核台）
- `expertProfileDetail.vue`：专家详情整合页（隐藏路由 `expertProfileDetail/:id`），背景信息 + 研究成果/决策影响记录/学术兼职三个 Tab 一次看全，是从"信息散落 4 个独立页面"这个设计问题上迭代出来的
- `expertApproval.vue`：审核台，"我发起的"/"待我审核的"两个 tab，按当前角色只渲染有权限的区块
- `expertAchievement.vue`/`expertAdoptionRecord.vue`/`expertAcademicPosition.vue`：各自独立的全量管理页（管理员跨专家批量管理用）
- `expertTag.vue`：标签库
- `expertSearch.vue`：检索推荐，展示相关性/成果分/决策影响分/社会贡献分/综合得分排序结果

只读角色（9002/9003）在所有这些页面上，新增/编辑/删除按钮直接不渲染（`isReadOnlyReviewer` 计算属性判断 `authorityId`），而不是渲染出来点了再被后端拒绝——但后端权限本身才是真正的边界，前端隐藏只是体验优化。

## 登录页

`web/src/view/login/index.vue` 现在是一个页面里塞了两套皮肤（专家库品牌 / 教育数据库品牌），右上角按钮切换，选择存 localStorage。两套皮肤共用同一套登录表单逻辑和账号体系，纯粹是视觉差异——这套系统本身没有"多租户"概念，教育数据库和专家库是同一个后端、同一批用户表，只是历史上先做了教育数据库、后加了专家库，两边为了品牌区分度各自保留了一套登录页外观。
