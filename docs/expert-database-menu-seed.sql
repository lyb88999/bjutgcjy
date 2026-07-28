-- 专家库菜单种子（幂等，可反复执行，已存在的菜单不会重复插入）
--
-- 背景：本仓库现有的业务菜单（如"工程教育数据库"）都不是通过 Go 代码种子创建的，
-- 而是通过 autocode 生成器或"系统工具 -> 菜单管理"在运行时写入 sys_base_menus 表。
-- 专家库照同样的方式创建菜单后，把最终效果导出成本文件提交到仓库，方便其他环境复现，
-- 而不必每个环境都重新手工点一遍菜单管理。
--
-- 使用前请确认：
--   1. 目标环境 sys_base_menus 表结构与本仓库一致（字段：menu_level/parent_id/path/name/
--      hidden/component/sort/active_name/keep_alive/default_menu/title/icon/close_tab）。
--   2. 执行完本脚本后，仍需要在"角色管理"里把这些菜单勾选给对应角色（至少是管理员角色），
--      并在"API管理 + 角色管理"里给角色勾选专家库相关接口（Casbin 策略不会因为建了菜单就自动生效）。
--   3. 建议先在开发库跑一遍确认无误，再应用到其他环境。

INSERT INTO sys_base_menus
  (created_at, updated_at, menu_level, parent_id, path, name, hidden, component, sort, active_name, keep_alive, default_menu, title, icon, close_tab)
SELECT NOW(3), NOW(3), 0, '0', 'expertDatabase', 'expertDatabase', 0, 'view/routerHolder.vue', 0, '', 0, 0, '专家库', 'user', 0
WHERE NOT EXISTS (
  SELECT 1 FROM sys_base_menus WHERE path = 'expertDatabase' AND deleted_at IS NULL
);

SET @expert_parent_id := (
  SELECT id FROM sys_base_menus WHERE path = 'expertDatabase' AND deleted_at IS NULL LIMIT 1
);

INSERT INTO sys_base_menus
  (created_at, updated_at, menu_level, parent_id, path, name, hidden, component, sort, active_name, keep_alive, default_menu, title, icon, close_tab)
SELECT NOW(3), NOW(3), 0, @expert_parent_id, 'expertApproval', 'expertApproval', 0, 'view/beijingExpertDatabase/expertApproval/expertApproval.vue', 0, '', 0, 0, '审核台', 'finished', 0
WHERE NOT EXISTS (
  SELECT 1 FROM sys_base_menus WHERE path = 'expertApproval' AND deleted_at IS NULL
);

INSERT INTO sys_base_menus
  (created_at, updated_at, menu_level, parent_id, path, name, hidden, component, sort, active_name, keep_alive, default_menu, title, icon, close_tab)
SELECT NOW(3), NOW(3), 0, @expert_parent_id, 'expertProfile', 'expertProfile', 0, 'view/beijingExpertDatabase/expertProfile/expertProfile.vue', 1, '', 0, 0, '专家主档', 'user', 0
WHERE NOT EXISTS (
  SELECT 1 FROM sys_base_menus WHERE path = 'expertProfile' AND deleted_at IS NULL
);

INSERT INTO sys_base_menus
  (created_at, updated_at, menu_level, parent_id, path, name, hidden, component, sort, active_name, keep_alive, default_menu, title, icon, close_tab)
SELECT NOW(3), NOW(3), 0, @expert_parent_id, 'expertAchievement', 'expertAchievement', 0, 'view/beijingExpertDatabase/expertAchievement/expertAchievement.vue', 2, '', 0, 0, '研究成果', 'collection', 0
WHERE NOT EXISTS (
  SELECT 1 FROM sys_base_menus WHERE path = 'expertAchievement' AND deleted_at IS NULL
);

INSERT INTO sys_base_menus
  (created_at, updated_at, menu_level, parent_id, path, name, hidden, component, sort, active_name, keep_alive, default_menu, title, icon, close_tab)
SELECT NOW(3), NOW(3), 0, @expert_parent_id, 'expertAdoptionRecord', 'expertAdoptionRecord', 0, 'view/beijingExpertDatabase/expertAdoptionRecord/expertAdoptionRecord.vue', 3, '', 0, 0, '决策影响记录', 'data-board', 0
WHERE NOT EXISTS (
  SELECT 1 FROM sys_base_menus WHERE path = 'expertAdoptionRecord' AND deleted_at IS NULL
);

INSERT INTO sys_base_menus
  (created_at, updated_at, menu_level, parent_id, path, name, hidden, component, sort, active_name, keep_alive, default_menu, title, icon, close_tab)
SELECT NOW(3), NOW(3), 0, @expert_parent_id, 'expertAcademicPosition', 'expertAcademicPosition', 0, 'view/beijingExpertDatabase/expertAcademicPosition/expertAcademicPosition.vue', 4, '', 0, 0, '学术兼职', 'avatar', 0
WHERE NOT EXISTS (
  SELECT 1 FROM sys_base_menus WHERE path = 'expertAcademicPosition' AND deleted_at IS NULL
);

INSERT INTO sys_base_menus
  (created_at, updated_at, menu_level, parent_id, path, name, hidden, component, sort, active_name, keep_alive, default_menu, title, icon, close_tab)
SELECT NOW(3), NOW(3), 0, @expert_parent_id, 'expertTag', 'expertTag', 0, 'view/beijingExpertDatabase/expertTag/expertTag.vue', 5, '', 0, 0, '标签库', 'notebook', 0
WHERE NOT EXISTS (
  SELECT 1 FROM sys_base_menus WHERE path = 'expertTag' AND deleted_at IS NULL
);

INSERT INTO sys_base_menus
  (created_at, updated_at, menu_level, parent_id, path, name, hidden, component, sort, active_name, keep_alive, default_menu, title, icon, close_tab)
SELECT NOW(3), NOW(3), 0, @expert_parent_id, 'expertSearch', 'expertSearch', 0, 'view/beijingExpertDatabase/expertSearch/expertSearch.vue', 6, '', 0, 0, '专家检索推荐', 'search', 0
WHERE NOT EXISTS (
  SELECT 1 FROM sys_base_menus WHERE path = 'expertSearch' AND deleted_at IS NULL
);

-- 专家详情整合页（专家主档"详情"按钮 + 审核台"查看详情"跳转的目标页，一次性展示背景信息
-- 加上研究成果/决策影响记录/学术兼职三个 Tab，不用再手动记专家 ID 去另外三个独立页面筛选）。
-- hidden=1：不出现在侧边栏，仅作为 sys_base_menu 里的一条路由记录供 vue-router 识别，
-- 参照本仓库已有的 autoCodeEdit/:id 这条隐藏参数路由的写法。
INSERT INTO sys_base_menus
  (created_at, updated_at, menu_level, parent_id, path, name, hidden, component, sort, active_name, keep_alive, default_menu, title, icon, close_tab)
SELECT NOW(3), NOW(3), 0, @expert_parent_id, 'expertProfileDetail/:id', 'expertProfileDetail', 1, 'view/beijingExpertDatabase/expertProfile/expertProfileDetail.vue', 7, '', 0, 0, '专家详情', 'user', 0
WHERE NOT EXISTS (
  SELECT 1 FROM sys_base_menus WHERE path = 'expertProfileDetail/:id' AND deleted_at IS NULL
);

-- 单位管理：sys_organization 的增删改查页面，只给管理员用（见 expert-database-permission-seed.sql），
-- 用来维护"单位审核员只能审自己单位提交的档案"这套数据域隔离依赖的单位树，
-- 以及给批量导入/真实数据回填 org_id 提供可持续维护的入口，不用再靠手工跑 SQL
INSERT INTO sys_base_menus
  (created_at, updated_at, menu_level, parent_id, path, name, hidden, component, sort, active_name, keep_alive, default_menu, title, icon, close_tab)
SELECT NOW(3), NOW(3), 0, @expert_parent_id, 'sysOrganization', 'sysOrganization', 0, 'view/beijingExpertDatabase/sysOrganization/sysOrganization.vue', 8, '', 0, 0, '单位管理', 'office-building', 0
WHERE NOT EXISTS (
  SELECT 1 FROM sys_base_menus WHERE path = 'sysOrganization' AND deleted_at IS NULL
);

-- 本单位账号管理：只给单位审核员用（见 expert-database-roles-seed.sql），
-- 让审核员自己给本单位新建个人申报人账号，不用事事找超级管理员建账号
INSERT INTO sys_base_menus
  (created_at, updated_at, menu_level, parent_id, path, name, hidden, component, sort, active_name, keep_alive, default_menu, title, icon, close_tab)
SELECT NOW(3), NOW(3), 0, @expert_parent_id, 'expertOrgUser', 'expertOrgUser', 0, 'view/beijingExpertDatabase/expertOrgUser/expertOrgUser.vue', 9, '', 0, 0, '本单位账号管理', 'user-filled', 0
WHERE NOT EXISTS (
  SELECT 1 FROM sys_base_menus WHERE path = 'expertOrgUser' AND deleted_at IS NULL
);

-- 统计概览：只给管理员/单位审核员/市级审核员用（见 expert-database-roles-seed.sql），
-- 个人申报人只关心自己提交的那几条，不需要看全库/全单位的统计图表
INSERT INTO sys_base_menus
  (created_at, updated_at, menu_level, parent_id, path, name, hidden, component, sort, active_name, keep_alive, default_menu, title, icon, close_tab)
SELECT NOW(3), NOW(3), 0, @expert_parent_id, 'expertDashboard', 'expertDashboard', 0, 'view/beijingExpertDatabase/expertDashboard/expertDashboard.vue', 10, '', 0, 0, '统计概览', 'data-analysis', 0
WHERE NOT EXISTS (
  SELECT 1 FROM sys_base_menus WHERE path = 'expertDashboard' AND deleted_at IS NULL
);

-- 使用帮助：纯前端静态内容页，不调用任何后端接口，四个角色都给（内容按角色自动折叠/展开，
-- 见 expertHelp.vue），不需要额外的 Casbin 接口权限
INSERT INTO sys_base_menus
  (created_at, updated_at, menu_level, parent_id, path, name, hidden, component, sort, active_name, keep_alive, default_menu, title, icon, close_tab)
SELECT NOW(3), NOW(3), 0, @expert_parent_id, 'expertHelp', 'expertHelp', 0, 'view/beijingExpertDatabase/expertHelp/expertHelp.vue', 11, '', 0, 0, '使用帮助', 'question-filled', 0
WHERE NOT EXISTS (
  SELECT 1 FROM sys_base_menus WHERE path = 'expertHelp' AND deleted_at IS NULL
);
