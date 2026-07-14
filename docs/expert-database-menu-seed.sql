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
