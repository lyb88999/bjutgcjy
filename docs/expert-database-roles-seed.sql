-- 专家库三级审核角色种子：个人申报人 / 单位审核员 / 市级审核员
-- 管理员角色沿用系统自带的超级管理员（无需新建，见 expert-database-permission-seed.sql）。
--
-- 每个角色的菜单都给了完整的专家库菜单树（是否能操作由下面的 Casbin 策略决定，
-- 而不是靠隐藏菜单——单位/市级审核员就算能看到"专家主档"列表，也只有只读权限）。
--
-- 使用前提：docs/expert-database-menu-seed.sql 已经执行过。
-- 使用后续步骤：在"用户管理"里把具体用户的角色改成对应角色、并设置"所属单位"
-- （单位审核员必须设置所属单位，否则 assertSameOrg 会拒绝所有审核操作）。

-- parent_id 必须是 0（顶级角色），不能是 NULL——角色树构建逻辑只认 0，NULL 会导致这些角色
-- 从"用户管理"的角色选择器里彻底消失（即使 sys_authorities 表里已经有记录）。
INSERT IGNORE INTO sys_authorities (created_at, updated_at, authority_id, authority_name, parent_id, default_router) VALUES
  (NOW(3), NOW(3), 9001, '专家库-个人申报人', 0, 'dashboard'),
  (NOW(3), NOW(3), 9002, '专家库-单位审核员', 0, 'dashboard'),
  (NOW(3), NOW(3), 9003, '专家库-市级审核员', 0, 'dashboard');

-- 菜单：三个角色都给仪表盘 + 专家库全树，具体操作权限由 Casbin 控制
INSERT IGNORE INTO sys_authority_menus (sys_base_menu_id, sys_authority_authority_id)
SELECT m.id, a.authority_id
FROM sys_base_menus m
JOIN (SELECT 9001 AS authority_id UNION SELECT 9002 UNION SELECT 9003) a ON 1=1
WHERE m.id = 1 OR m.path IN ('expertDatabase','expertApproval','expertProfile','expertAchievement','expertAdoptionRecord','expertAcademicPosition','expertTag','expertSearch','expertProfileDetail/:id')
  AND m.deleted_at IS NULL;

-- 本单位账号管理这个菜单只给单位审核员（9002），个人申报人/市级审核员没有"自己的单位"这个概念
INSERT IGNORE INTO sys_authority_menus (sys_base_menu_id, sys_authority_authority_id)
SELECT id, 9002 FROM sys_base_menus WHERE path = 'expertOrgUser' AND deleted_at IS NULL;

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

-- 审核台是三个角色共用的同一个页面：一进页面就会请求"我发起的"（myDrafts），
-- 点开"待我审核的"tab 会同时请求单位待审（pendingOrgReview）和市级待审（pendingCityReview）——
-- 不管当前角色是不是真的用得上那部分数据。三个角色都必须能读这三个只读列表接口，
-- 否则某个角色一进审核台或者切 tab 就会先弹一次"权限不足"，即使他真正需要看的那部分数据是正常的。
-- （expert_approval.go 的 GetPendingOrgReview 对没有关联单位的账号会返回空列表而不是报错，
-- 所以这里放开权限是安全的，不会导致越权查看其他单位的数据。）
INSERT IGNORE INTO casbin_rule (ptype, v0, v1, v2, v3, v4, v5)
SELECT 'p', a.authority_id, ep.path, 'GET', '', '', ''
FROM (SELECT 9001 AS authority_id UNION SELECT 9002 UNION SELECT 9003) a
JOIN (
  SELECT '/expertApproval/myDrafts' AS path UNION ALL
  SELECT '/expertApproval/pendingOrgReview' UNION ALL
  SELECT '/expertApproval/pendingCityReview'
) ep ON 1=1;

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
('p','9001','/expertApproval/getApprovalLogList','GET','','',''),

('p','9001','/expertProfile/downloadImportTemplate','GET','','',''),
('p','9001','/expertProfile/importBatch','POST','','','');

-- 9002 专家库-单位审核员：只读专家档案（含成果/决策影响/学术兼职/标签的列表+详情）+ 本单位审核操作
-- （跨单位由后端 assertSameOrg 拦截）。故意不给 create/update/delete —— 审核员不应该能改动
-- 正在审核的申报人自己的成果数据，否则"审核通过"就失去意义了。
INSERT IGNORE INTO casbin_rule (ptype, v0, v1, v2, v3, v4, v5) VALUES
('p','9002','/expertProfile/findExpertProfile','GET','','',''),
('p','9002','/expertProfile/getExpertProfileList','GET','','',''),
('p','9002','/expertAchievement/findExpertAchievement','GET','','',''),
('p','9002','/expertAchievement/getExpertAchievementList','GET','','',''),
('p','9002','/expertAdoptionRecord/findExpertAdoptionRecord','GET','','',''),
('p','9002','/expertAdoptionRecord/getExpertAdoptionRecordList','GET','','',''),
('p','9002','/expertAcademicPosition/findExpertAcademicPosition','GET','','',''),
('p','9002','/expertAcademicPosition/getExpertAcademicPositionList','GET','','',''),
('p','9002','/expertTag/findExpertTag','GET','','',''),
('p','9002','/expertTag/getExpertTagList','GET','','',''),
('p','9002','/expertTag/getExpertTagsByExpertId','GET','','',''),
('p','9002','/expertApproval/pendingOrgReview','GET','','',''),
('p','9002','/expertApproval/orgApprove','POST','','',''),
('p','9002','/expertApproval/orgReject','POST','','',''),
('p','9002','/expertApproval/getApprovalLogList','GET','','',''),

-- 本单位账号管理：单位审核员给本单位新建个人申报人账号，不用事事找超级管理员。
-- org_id/角色都由后端固定推导（见 expert_org_user.go），这三个接口本身不接受越权参数，
-- 但仍然只放给 9002——申报人和市级审核员没有"自己的单位"这个概念，不适用这套能力
('p','9002','/expertOrgUser/createOrgApplicant','POST','','',''),
('p','9002','/expertOrgUser/getOrgUserList','GET','','',''),
('p','9002','/expertOrgUser/toggleOrgUserEnable','POST','','','');

-- 9003 专家库-市级审核员：只读专家档案（含成果/决策影响/学术兼职/标签的列表+详情）+ 市级审核操作
-- （不限单位）。同样故意不给 create/update/delete，理由同 9002。
INSERT IGNORE INTO casbin_rule (ptype, v0, v1, v2, v3, v4, v5) VALUES
('p','9003','/expertProfile/findExpertProfile','GET','','',''),
('p','9003','/expertProfile/getExpertProfileList','GET','','',''),
('p','9003','/expertAchievement/findExpertAchievement','GET','','',''),
('p','9003','/expertAchievement/getExpertAchievementList','GET','','',''),
('p','9003','/expertAdoptionRecord/findExpertAdoptionRecord','GET','','',''),
('p','9003','/expertAdoptionRecord/getExpertAdoptionRecordList','GET','','',''),
('p','9003','/expertAcademicPosition/findExpertAcademicPosition','GET','','',''),
('p','9003','/expertAcademicPosition/getExpertAcademicPositionList','GET','','',''),
('p','9003','/expertTag/findExpertTag','GET','','',''),
('p','9003','/expertTag/getExpertTagList','GET','','',''),
('p','9003','/expertTag/getExpertTagsByExpertId','GET','','',''),
('p','9003','/expertApproval/pendingCityReview','GET','','',''),
('p','9003','/expertApproval/cityApprove','POST','','',''),
('p','9003','/expertApproval/cityReject','POST','','',''),
('p','9003','/expertApproval/getApprovalLogList','GET','','','');
