# 北京市哲学社会科学专家库 技术方案

版本：v0.1 草案　日期：2026-07-14

---

## 0. 现状梳理

当前仓库是基于 **gin-vue-admin**（Go + Gin + GORM + Vue3 + Element-Plus）的后台管理系统，已实现的业务是"工程教育发展监测数据采集"：

- 后端 `server/model|service|router|api/v1/EngineeringEducationDatabase`：8 张业务表（基础信息库、学科建设库、人才培养库…），均由 gin-vue-admin 的 **autocode** 代码生成器生成（模型文件头部注释"自动生成模板"），CRUD 风格高度一致。
- 权限体系：沿用 gin-vue-admin 原生 RBAC —— `sys_user` / `sys_authority`（角色，树形 ParentId）/ `sys_authority_menus`（角色-菜单）/ Casbin（`sys_casbin`，接口级鉴权）/ `DataAuthorityId`（数据权限，角色间可见数据范围）。
- 登录页 `web/src/view/login/index.vue` 已经改过一次皮，写死了"北京工业大学采集点"字样；品牌名取自 `web/src/core/config.js` 的 `appName`（当前值"工程教育发展监测数据采集"），Logo 取自 `appLogo`。`index_new.vue` 是未改造的原始模板，未被引用。
- 前端已存在一个**空目录** `web/src/view/beijingExpertDatabase/`，无文件、无路由、无对应后端代码——说明专家库此前只建了壳，尚未真正开发。
- 用户提供的《专家库项目指标建设 0713》文档定义了专家画像的 **五个维度**：① 背景信息画像 ② 学科领域画像 ③ 研究主题画像 ④ 研究成果画像 ⑤ 决策影响画像，每个维度下有若干具体指标（详见第 1 节）。

**架构决策（已与用户确认）**：专家库不新起部署，作为**同一系统内的新业务模块**接入，复用现有 `sys_user`/Casbin/菜单体系；登录页整体品牌切换为"北京市哲学社会科学专家库"，工程教育数据库继续作为菜单里的一个业务模块保留。

本方案覆盖四块内容：
1. 专家库指标体系 → 数据模型
2. 成果分级赋值 + 综合推荐排序算法
3. 专家数据录入的三级审核工作流（个人申报 → 单位审核 → 市级审核，管理员全权限）
4. 登录页品牌重塑

---

## 1. 专家库指标体系与数据模型

### 1.1 指标体系还原（来自 0713 文档）

| 画像维度 | 具体指标 |
|---|---|
| 一、背景信息画像 | 基本信息（姓名/性别/出生年月/民族/政治面貌）、单位信息（单位/院系部门/行政职务/专技职称）、联系方式、学历背景、荣誉标识（人才计划/荣誉称号）、学术兼职 |
| 二、学科领域画像 | 一级学科、二级学科、交叉学科领域、所属学科平台/研究基地 |
| 三、研究主题画像 | 近五年研究方向、研究关键词、所属政策领域、主要研究对象、当前关注议题、区域/国别研究专长、研究方法专长 |
| 四、研究成果画像 | 代表性著作、代表性学术论文、代表性研究报告、代表性决策咨询成果、代表性获奖成果及荣誉、主持重要课题项目（级别+来源） |
| 五、决策影响画像 | 成果被内参/专报采用情况、成果被党政部门采纳/领导批示情况、研究观点进入政策文件/规划情况、参与政策起草/规划编制/咨询论证情况 |

维度四、五的指标本质是**多值、可分级、可计数**的列表型数据（一个专家可以有 N 篇论文、N 次采纳），不适合像工程教育库那样"一行打平所有字段"，需要拆成主表 + 明细表结构。这也正是第 2 节赋值算法所需要的数据基础。

### 1.2 表结构设计

沿用现有代码风格（`global.GVA_MODEL` 基类含 ID/CreatedAt/UpdatedAt/DeletedAt，`CreatedBy/UpdatedBy/DeletedBy` 审计字段），新建 Go module `server/model/ExpertDatabase`，建议表：

**主表**
- `expert_profile`（专家主档，对应维度一、二、三 —— 这些是"一对一"属性，直接打平存字段即可）
  - 背景信息：`name, gender, birth_date, ethnicity, political_status, unit_name, department, admin_title, tech_title, phone, mobile, email, address, highest_degree, highest_education, graduate_school, major, talent_program, honor_title` …
  - 学术兼职：`academic_positions`（JSON 数组或独立明细表 `expert_academic_position`，兼职条数一般不多，JSON 即可，如需按机构类型筛选再拆表）
  - 学科领域：`discipline_l1, discipline_l2, cross_discipline, discipline_platform`
  - 研究主题：`research_directions, research_keywords, policy_fields, research_objects, focus_topics, region_expertise, method_expertise`（这些字段本身是"多值打标签"，用 `,` 分隔字符串或独立标签表 `expert_tag`，**建议用独立标签表**，因为第 3 节的"综合排序/检索匹配"要按标签做关联查询和词频统计，打平成字符串会导致后续 LIKE 查询效率低、无法做词频统计）
  - 状态字段（用于第 3 节审核流）：`status, org_id, submitted_by, current_step` 等，见 3.2
  - 综合评分缓存字段（用于第 2 节排序，避免每次检索都实时算分）：`achievement_score, influence_score, composite_score, score_updated_at`

- `expert_tag`（标签表，多对多）：`id, tag_type(discipline_l1/keyword/policy_field/region/method...), tag_value`，专家-标签关联表 `expert_tag_relation(expert_id, tag_id)`

**成果明细表（维度四）**，统一用一张"成果表"而不是五张表，用 `achievement_type` 区分，便于第 2 节统一算分：

```
expert_achievement
  id
  expert_id            -- 关联专家
  achievement_type      -- enum: 著作/论文/研究报告/决策咨询成果/获奖成果/课题项目
  title                 -- 成果名称
  level                 -- 成果级别（国家级/省部级/厅局级/一般，字典项，见2.1）
  level_source          -- 级别认定来源（如"课题下达单位"、"获奖颁发单位"）
  publish_org           -- 发表/出版/立项单位
  publish_date
  keywords              -- 该成果的关键词/摘要分词结果，用于第2节相关性计算
  ref_count             -- 检索相关次数（第2节"篇数"的计数依据）
  remark
  created_by / updated_by / deleted_by
```

**决策影响明细表（维度五）**，这是"成果采纳单位级别"的核心承载表：

```
expert_adoption_record
  id
  expert_id
  achievement_id        -- 可关联到具体成果，也可独立记录（比如笼统的"某次被采纳"）
  adoption_type         -- enum: 内参/专报采用、部门采纳、领导批示、进入政策文件、参与政策起草/咨询论证
  adopting_unit_level   -- 采纳/批示单位级别（国家级/省部级/厅局级/区县级，字典项，见2.1，这是用户要求新增的关键字段）
  adopting_unit_name    -- 具体单位名称
  adoption_date
  evidence_desc         -- 佐证材料描述/文号
  attachment_url        -- 佐证材料附件
  created_by / updated_by / deleted_by
```

> 这张表直接对应用户需求里"评价指标体系后两部分是否要加上成果采纳单位级别的判断"——`adopting_unit_level` 就是这个判断字段，级别字典可扩展、赋值权重可配置（见 2.1），不写死在代码里。

前端目录直接复用已存在的空壳 `web/src/view/beijingExpertDatabase/`，内部按 `expertList`（专家主档管理）、`expertAchievement`（成果管理）、`expertApproval`（审核台）、`expertSearch`（检索推荐）拆子目录。

后端沿用 autocode 生成 CRUD 基础代码，审核流转、评分计算部分再手写 service 补充（autocode 生成的代码不会覆盖手写扩展，这是 gin-vue-admin 的常规用法）。

---

## 2. 成果分级赋值 + 综合推荐排序算法

### 2.1 分级赋值模型

需求原文："重要指标按级别赋值，普通指标按次数赋值。比如省部级赋5，厅局级赋3；篇数按检索相关次数赋值；相关性可以按成果库词频检索赋值。计算方法为：级别权重 × 相应级别篇数的求和 × 相关性系数"。

设计为**可配置评分卡**，而非硬编码，落地方式：复用 gin-vue-admin 已有的 `sys_dictionary`/`sys_dictionary_detail`（数据字典）机制，新增字典：

| 字典 | 示例项 | 用途 |
|---|---|---|
| `expert_achievement_level`（成果/课题级别） | 国家级=10，省部级=5，厅局级=3，一般级=1 | `expert_achievement.level` 的权重来源 |
| `expert_adoption_level`（采纳单位级别） | 中央/国家级=10，省部级=5，厅局级=3，区县级=1 | `expert_adoption_record.adopting_unit_level` 的权重来源 |

字典详情表 `sys_dictionary_detail` 本身有 `label/value/extend` 字段，`extend` 可以直接存权重分值 JSON，管理员可在系统工具-数据字典页面里调整级别与权重，不需要发版。

**单个专家的成果得分**（对应"级别权重 × 相应级别篇数的求和 × 相关性系数"）：

```
achievement_score(expert, query) =
    Σ over level L [ weight(L) × count(expert, L) ] × relevance(expert, query)
```

- `weight(L)`：从字典表取该级别权重
- `count(expert, L)`：该专家在该级别下的成果篇数（`expert_achievement` 按 level 分组 count）
- `relevance(expert, query)`：相关性系数，见 2.2

**决策影响得分**（对应"成果采纳单位级别"新增判断），结构完全一致，只是数据源换成 `expert_adoption_record`：

```
influence_score(expert, query) =
    Σ over level L [ weight(L) × count(expert, L) ] × relevance(expert, query)
```

两者都是"重要指标按级别赋值、按次数累加"的同一套公式，只是作用在不同的表上，这样代码可以写成一个通用函数 `calcLeveledScore(records []Record, levelWeights map[string]float64, relevance float64) float64`，两处复用。

### 2.2 相关性系数（检索词频）

"相关性可以按成果库词频检索赋值"——即给定一个检索主题（如"防汛"），衡量专家的成果集合与该主题的匹配程度：

- **V1（快速上线，够用）**：对 `expert_achievement.title/keywords` 与专家标签 `expert_tag` 做分词（Go 侧用 `gojieba` 或 `sego` 中文分词），构建关键词计数；检索时对 query 分词，计算 query 词在专家成果关键词集合里的**命中词频 / 总词数**，作为 `relevance ∈ [0,1]`。可以先用简单的 `LIKE` + 计数实现 MVP，避免引入新中间件。
- **V2（检索量上来后再演进）**：接入 Elasticsearch / Meilisearch，用 BM25 或 TF-IDF 算相关性，检索性能和相关性精度都更好，同时能支持同义词（比如"防汛"↔"防洪""汛期"）。
- 相关性系数**不建议**做成 0/1 布尔命中，而是连续值，这样综合排序时同一主题下的专家还能按匹配程度精细区分。

### 2.3 综合推荐排序（检索场景）

需求原文的场景是"我需要提防汛方案，需要按发表成果等级、数量、职称、社会贡献等生成综合排序"。这是一个**多因子加权排序**，不是纯搜索匹配：

```
composite_score(expert, query) =
    w1 × achievement_score(expert, query)      -- 成果等级×数量×相关性
  + w2 × influence_score(expert, query)        -- 决策影响/采纳单位级别×数量×相关性
  + w3 × title_score(expert)                    -- 职称权重（教授/研究员 > 副教授 > 讲师，字典配置）
  + w4 × social_contribution_score(expert)      -- 社会贡献（学术兼职级别、荣誉称号，字典配置）
```

- `w1..w4` 同样做成可配置项（后台一个"排序权重配置"页面，或直接放 `sys_dictionary`），而不是写死，方便运营根据不同课题类型调整权重（比如应急类课题可能更看重决策影响，理论研究类课题更看重论文等级）。
- 检索接口 `POST /expertDatabase/search`：入参 `keyword`（如"防汛方案"）、筛选条件（学科/地区/职称等，命中 `expert_tag`/字段过滤）、分页；服务端先按筛选条件圈定候选集，再对候选集算 `composite_score` 排序返回。
- **性能考虑**：`composite_score` 中职称分、社会贡献分与 query 无关，可以离线预计算存 `expert_profile.achievement_score/influence_score`（专家成果变更、审核通过时触发重算，用异步任务 `server/task` 或写入时同步算），检索时只需实时算 `relevance(expert, query)` 这一个依赖 query 的因子，其余用缓存值，大幅降低检索时的计算量。

---

## 3. 专家数据录入三级审核工作流

### 3.1 组织与角色

新增 `sys_organization`（单位）表（当前系统没有"单位"概念，只有 `sys_user`，需要补上）：

```
sys_organization
  id, name, level(省部级/市级/区级...), parent_id, region_code
```

`sys_user` 增加 `org_id`，标识该用户归属的申报单位。

角色沿用 `sys_authority`（树形角色，天然支持"父角色"），新增 4 个角色：

| 角色 | AuthorityId 示例 | 权限范围 |
|---|---|---|
| 专家/个人申报人 | 用户自己 | 只能填报/编辑/查看自己申报的专家档案，未通过前可撤回重填 |
| 单位审核员 | 挂在具体 `org_id` 下 | 审核本单位提交的专家档案（Casbin 接口权限 + `DataAuthorityId` 数据权限限定 `org_id`） |
| 市级审核员 | 顶层单位 | 审核所有已通过单位审核、待市级审核的档案 |
| 系统管理员 | `authority_id = 888`（超管，系统已有） | 全部权限：可直接编辑/审核/退回任意环节，无需走流程 |

沿用 gin-vue-admin 原生的"角色-菜单-接口(Casbin)-数据权限(DataAuthorityId)"四层模型即可实现"谁能看什么、谁能审什么"，不需要额外引入工作流引擎，量级上三级审核用状态机足够，没必要上 Camunda/Flowable 这类重型 BPM。

### 3.2 状态机设计

`expert_profile.status` 枚举：

```
draft(草稿/个人填报中)
  → pending_org_review(待单位审核)
    → org_rejected(单位退回，回到 draft，附退回意见)
    → pending_city_review(待市级审核，单位审核通过)
      → city_rejected(市级退回，回到 draft 或 pending_org_review，附退回意见)
      → published(市级审核通过，正式收录，可被检索排序)
```

- `expert_profile` 增加 `submitted_by`（申报人 user_id）、`org_id`（申报单位，审核路由依据）、`current_step`。
- 新增 `expert_approval_log` 审核留痕表：`id, expert_id, from_status, to_status, operator_id, opinion, created_at`，每次提交/审核/退回都写一条，满足留痕审计要求，也是后续做"审核时效统计"的数据基础。
- 管理员对 `expert_approval_log`/`status` 拥有直接改写权限（跳过流程），满足"管理员有所有权限"的要求，但仍应记录日志，避免绕过审计。
- 只有 `status = published` 的专家档案才参与第 2 节的检索排序，草稿/审核中的数据不对外可见，保证检索结果的权威性。

### 3.3 接口设计（示意）

```
POST   /expertDatabase/expert                 创建草稿（个人）
PUT    /expertDatabase/expert                  编辑草稿
POST   /expertDatabase/expert/submit           提交单位审核（draft → pending_org_review）
POST   /expertDatabase/expert/orgApprove       单位审核通过（→ pending_city_review）
POST   /expertDatabase/expert/orgReject        单位审核退回（→ org_rejected/draft，需填 opinion）
POST   /expertDatabase/expert/cityApprove      市级审核通过（→ published）
POST   /expertDatabase/expert/cityReject       市级审核退回（需填 opinion）
GET    /expertDatabase/expert/approvalLog      查看某专家的审核历史
```

前端 `expertApproval` 子模块做成一个统一的"审核台"页面，根据当前登录角色（单位审核员/市级审核员/管理员）在同一个页面组件里过滤展示待审列表，区分"我发起的""待我审核的""历史记录"三个 Tab。

---

## 4. 登录页品牌重塑

- 品牌名与 Logo 来自 `web/src/core/config.js`：`appName: '工程教育发展监测数据采集'` → 改为 `'北京市哲学社会科学专家库'`；`appLogo` 换成专家库专用 Logo 资源。
- `web/src/view/login/index.vue` 中硬编码的副标题"北京工业大学采集点"字样（对应工程教育库的品牌语境）替换为专家库的定位语，例如"北京市哲学社会科学专家库 · 智库专家资源管理与推荐平台"。
- 视觉上可以顺带调整背景色 `bg-[#194bfb]`/斜切色块为专家库的主题色，与工程教育库的蓝色系做区分，避免用户混淆"我登录的是哪个系统"。
- `index_new.vue`（未使用的原始模板）建议保留不动或直接删除，避免误用。
- 工程教育数据库模块继续挂在登录后的菜单里（`server/router/EngineeringEducationDatabase` 不动），只是不再作为"登录页第一印象"出现。

---

## 5. 实施路线图（建议）

| 阶段 | 内容 | 产出 |
|---|---|---|
| S1 基础建模 | `sys_organization` 建表；`ExpertDatabase` module 的 `expert_profile/expert_achievement/expert_adoption_record/expert_tag` 建表 + autocode 生成基础 CRUD | 可录入、可查看的专家档案（无审核、无排序） |
| S2 审核流 | 状态机字段、`expert_approval_log`、4 类角色的 Casbin+菜单+数据权限配置、审核台前端 | 完整的个人申报→单位审核→市级审核闭环 |
| S3 评分与检索 | 字典驱动的级别权重、`calcLeveledScore` 通用算分、关键词分词 relevance V1、检索排序接口与前端检索页 | 支持"防汛方案"这类主题检索并给出综合排序 |
| S4 品牌与打磨 | 登录页/Logo/主题色替换、菜单入口整理、排序权重可配置后台页面 | 对外可用的专家库门户 |
| S5（可选）检索升级 | 引入 ES/Meilisearch 替换 V1 分词计数方案 | 检索量大时的性能与相关性提升 |

---

## 6. 待确认事项

1. **成果/采纳级别字典的具体档位与分值**：本方案先用"省部级=5、厅局级=3"作为示例，实际字典项（如是否要区分国家级、区县级、行业级）、以及 2.3 中 `w1..w4` 权重初始值，需要业务方给出第一版参考值，后续可在后台动态调整。
2. **相关性计算 V1 的分词方案**：是否已有现成的"成果库"全文检索能力可以复用，还是需要新建；如果专家成果数据量较大（万级以上），建议直接规划 S5 引入独立检索引擎，跳过 V1。
3. **单位数据来源**：`sys_organization` 是否需要对接现有的组织机构代码/编制数据，还是本项目内手工维护一份简化的单位列表。
