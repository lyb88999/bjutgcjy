-- 专家库三级审核角色种子：个人申报人 / 单位审核员 / 市级审核员
-- 管理员角色沿用系统自带的超级管理员（无需新建，见 expert-database-permission-seed.sql）。
--
-- 每个角色的菜单都给了完整的专家库菜单树（是否能操作由下面的 Casbin 策略决定，
-- 而不是靠隐藏菜单——单位/市级审核员就算能看到"专家主档"列表，也只有只读权限）。
--
-- 使用前提：docs/expert-database-menu-seed.sql 已经执行过。
-- 使用后续步骤：在"用户管理"里把具体用户的角色改成对应角色、并设置"所属单位"
-- （单位审核员必须设置所属单位，否则 assertSameOrg 会拒绝所有审核操作）。

INSERT IGNORE INTO sys_authorities (created_at, updated_at, authority_id, authority_name, parent_id, default_router) VALUES
  (NOW(3), NOW(3), 9001, '专家库-个人申报人', NULL, 'dashboard'),
  (NOW(3), NOW(3), 9002, '专家库-单位审核员', NULL, 'dashboard'),
  (NOW(3), NOW(3), 9003, '专家库-市级审核员', NULL, 'dashboard');

-- 菜单：三个角色都给仪表盘 + 专家库全树，具体操作权限由 Casbin 控制
INSERT IGNORE INTO sys_authority_menus (sys_base_menu_id, sys_authority_authority_id)
SELECT m.id, a.authority_id
FROM sys_base_menus m
JOIN (SELECT 9001 AS authority_id UNION SELECT 9002 UNION SELECT 9003) a ON 1=1
WHERE m.id = 1 OR m.path IN ('expertDatabase','expertApproval','expertProfile','expertAchievement','expertAdoptionRecord','expertAcademicPosition','expertTag','expertSearch')
  AND m.deleted_at IS NULL;

-- 每个角色都需要的登录/基础接口（对应 systemReq.DefaultCasbin()，新建角色的标准最小权限集）
INSERT IGNORE INTO casbin_rule (ptype, v0, v1, v2, v3, v4, v5)
SELECT 'p', a.authority_id, base.path, base.method, '', '', ''
FROM (SELECT 9001 AS authority_id UNION SELECT 9002 UNION SELECT 9003) a
JOIN (
  SELECT '/menu/getMenu' AS path, 'POST' AS method UNION ALL
  SELECT '/jwt/jsonInBlacklist', 'POST' UNION ALL
  SELECT '/user/admin_register', 'POST' UNION ALL
  SELECT '/user/changePassword', 'POST' UNION ALL
  SELECT '/user/setUserAuthority', 'POST' UNION ALL
  SELECT '/user/setUserInfo', 'PUT' UNION ALL
  SELECT '/user/getUserInfo', 'GET'
) base ON 1=1;

-- 专家检索：三个角色都能查看已发布专家的综合排序结果
INSERT IGNORE INTO casbin_rule (ptype, v0, v1, v2, v3, v4, v5)
SELECT 'p', a.authority_id, '/expertDatabase/search', 'GET', '', '', ''
FROM (SELECT 9001 AS authority_id UNION SELECT 9002 UNION SELECT 9003) a;

-- 9001 专家库-个人申报人：维护自己的专家档案与成果/影响/兼职记录，提交审核
INSERT IGNORE INTO casbin_rule (ptype, v0, v1, v2, v3, v4, v5) VALUES
('p','9001','/expertProfile/createExpertProfile','POST','','',''),
('p','9001','/expertProfile/updateExpertProfile','PUT','','',''),
('p','9001','/expertProfile/findExpertProfile','GET','','',''),
('p','9001','/expertProfile/getExpertProfileList','GET','','',''),

('p','9001','/expertAchievement/createExpertAchievement','POST','','',''),
('p','9001','/expertAchievement/deleteExpertAchievement','DELETE','','',''),
('p','9001','/expertAchievement/deleteExpertAchievementByIds','DELETE','','',''),
('p','9001','/expertAchievement/updateExpertAchievement','PUT','','',''),
('p','9001','/expertAchievement/findExpertAchievement','GET','','',''),
('p','9001','/expertAchievement/getExpertAchievementList','GET','','',''),

('p','9001','/expertAdoptionRecord/createExpertAdoptionRecord','POST','','',''),
('p','9001','/expertAdoptionRecord/deleteExpertAdoptionRecord','DELETE','','',''),
('p','9001','/expertAdoptionRecord/deleteExpertAdoptionRecordByIds','DELETE','','',''),
('p','9001','/expertAdoptionRecord/updateExpertAdoptionRecord','PUT','','',''),
('p','9001','/expertAdoptionRecord/findExpertAdoptionRecord','GET','','',''),
('p','9001','/expertAdoptionRecord/getExpertAdoptionRecordList','GET','','',''),

('p','9001','/expertAcademicPosition/createExpertAcademicPosition','POST','','',''),
('p','9001','/expertAcademicPosition/deleteExpertAcademicPosition','DELETE','','',''),
('p','9001','/expertAcademicPosition/deleteExpertAcademicPositionByIds','DELETE','','',''),
('p','9001','/expertAcademicPosition/updateExpertAcademicPosition','PUT','','',''),
('p','9001','/expertAcademicPosition/findExpertAcademicPosition','GET','','',''),
('p','9001','/expertAcademicPosition/getExpertAcademicPositionList','GET','','',''),

('p','9001','/expertTag/getExpertTagList','GET','','',''),
('p','9001','/expertTag/getExpertTagsByExpertId','GET','','',''),
('p','9001','/expertTag/setExpertTagRelations','POST','','',''),

('p','9001','/expertApproval/submit','POST','','',''),
('p','9001','/expertApproval/myDrafts','GET','','',''),
('p','9001','/expertApproval/getApprovalLogList','GET','','','');

-- 9002 专家库-单位审核员：只读专家档案 + 本单位审核操作（跨单位由后端 assertSameOrg 拦截）
INSERT IGNORE INTO casbin_rule (ptype, v0, v1, v2, v3, v4, v5) VALUES
('p','9002','/expertProfile/findExpertProfile','GET','','',''),
('p','9002','/expertProfile/getExpertProfileList','GET','','',''),
('p','9002','/expertAchievement/getExpertAchievementList','GET','','',''),
('p','9002','/expertAdoptionRecord/getExpertAdoptionRecordList','GET','','',''),
('p','9002','/expertAcademicPosition/getExpertAcademicPositionList','GET','','',''),
('p','9002','/expertApproval/pendingOrgReview','GET','','',''),
('p','9002','/expertApproval/orgApprove','POST','','',''),
('p','9002','/expertApproval/orgReject','POST','','',''),
('p','9002','/expertApproval/getApprovalLogList','GET','','','');

-- 9003 专家库-市级审核员：只读专家档案 + 市级审核操作（不限单位）
INSERT IGNORE INTO casbin_rule (ptype, v0, v1, v2, v3, v4, v5) VALUES
('p','9003','/expertProfile/findExpertProfile','GET','','',''),
('p','9003','/expertProfile/getExpertProfileList','GET','','',''),
('p','9003','/expertAchievement/getExpertAchievementList','GET','','',''),
('p','9003','/expertAdoptionRecord/getExpertAdoptionRecordList','GET','','',''),
('p','9003','/expertAcademicPosition/getExpertAcademicPositionList','GET','','',''),
('p','9003','/expertApproval/pendingCityReview','GET','','',''),
('p','9003','/expertApproval/cityApprove','POST','','',''),
('p','9003','/expertApproval/cityReject','POST','','',''),
('p','9003','/expertApproval/getApprovalLogList','GET','','','');
