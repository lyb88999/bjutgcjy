-- 专家库权限种子（幂等）：把专家库的菜单 + 接口权限授予"总管理员"角色（authority_id = 1）
-- 使用前提：docs/expert-database-menu-seed.sql 已经执行过，sys_base_menus 里已经有专家库的 8 个菜单行。
-- 如果你的环境里"超级管理员"不是 authority_id = 1，把下面所有的 1 换成实际的角色ID即可。

-- 菜单权限：把专家库父菜单 + 子菜单 + 隐藏的专家详情整合页都挂给该角色
INSERT IGNORE INTO sys_authority_menus (sys_base_menu_id, sys_authority_authority_id)
SELECT id, 1 FROM sys_base_menus
WHERE path IN ('expertDatabase','expertApproval','expertProfile','expertAchievement','expertAdoptionRecord','expertAcademicPosition','expertTag','expertSearch','expertProfileDetail/:id','sysOrganization','expertDashboard','expertHelp')
  AND deleted_at IS NULL;

-- 接口权限：Casbin 策略（v0=角色ID, v1=接口路径, v2=HTTP方法）
INSERT IGNORE INTO casbin_rule (ptype, v0, v1, v2, v3, v4, v5) VALUES
('p','1','/sysOrganization/createSysOrganization','POST','','',''),
('p','1','/sysOrganization/deleteSysOrganization','DELETE','','',''),
('p','1','/sysOrganization/updateSysOrganization','PUT','','',''),
('p','1','/sysOrganization/findSysOrganization','GET','','',''),
('p','1','/sysOrganization/getSysOrganizationList','GET','','',''),
('p','1','/sysOrganization/getSysOrganizationTree','GET','','',''),

('p','1','/expertProfile/createExpertProfile','POST','','',''),
('p','1','/expertProfile/deleteExpertProfile','DELETE','','',''),
('p','1','/expertProfile/deleteExpertProfileByIds','DELETE','','',''),
('p','1','/expertProfile/updateExpertProfile','PUT','','',''),
('p','1','/expertProfile/findExpertProfile','GET','','',''),
('p','1','/expertProfile/getExpertProfileList','GET','','',''),
('p','1','/expertProfile/exportExpertProfiles','GET','','',''),

('p','1','/expertAchievement/createExpertAchievement','POST','','',''),
('p','1','/expertAchievement/deleteExpertAchievement','DELETE','','',''),
('p','1','/expertAchievement/deleteExpertAchievementByIds','DELETE','','',''),
('p','1','/expertAchievement/updateExpertAchievement','PUT','','',''),
('p','1','/expertAchievement/findExpertAchievement','GET','','',''),
('p','1','/expertAchievement/getExpertAchievementList','GET','','',''),

('p','1','/expertAdoptionRecord/createExpertAdoptionRecord','POST','','',''),
('p','1','/expertAdoptionRecord/deleteExpertAdoptionRecord','DELETE','','',''),
('p','1','/expertAdoptionRecord/deleteExpertAdoptionRecordByIds','DELETE','','',''),
('p','1','/expertAdoptionRecord/updateExpertAdoptionRecord','PUT','','',''),
('p','1','/expertAdoptionRecord/findExpertAdoptionRecord','GET','','',''),
('p','1','/expertAdoptionRecord/getExpertAdoptionRecordList','GET','','',''),

('p','1','/expertAcademicPosition/createExpertAcademicPosition','POST','','',''),
('p','1','/expertAcademicPosition/deleteExpertAcademicPosition','DELETE','','',''),
('p','1','/expertAcademicPosition/deleteExpertAcademicPositionByIds','DELETE','','',''),
('p','1','/expertAcademicPosition/updateExpertAcademicPosition','PUT','','',''),
('p','1','/expertAcademicPosition/findExpertAcademicPosition','GET','','',''),
('p','1','/expertAcademicPosition/getExpertAcademicPositionList','GET','','',''),

('p','1','/expertTag/createExpertTag','POST','','',''),
('p','1','/expertTag/deleteExpertTag','DELETE','','',''),
('p','1','/expertTag/updateExpertTag','PUT','','',''),
('p','1','/expertTag/setExpertTagRelations','POST','','',''),
('p','1','/expertTag/findExpertTag','GET','','',''),
('p','1','/expertTag/getExpertTagList','GET','','',''),
('p','1','/expertTag/getExpertTagsByExpertId','GET','','',''),

('p','1','/expertApproval/submit','POST','','',''),
('p','1','/expertApproval/orgApprove','POST','','',''),
('p','1','/expertApproval/batchOrgApprove','POST','','',''),
('p','1','/expertApproval/orgReject','POST','','',''),
('p','1','/expertApproval/cityApprove','POST','','',''),
('p','1','/expertApproval/batchCityApprove','POST','','',''),
('p','1','/expertApproval/cityReject','POST','','',''),
('p','1','/expertApproval/adminSetStatus','POST','','',''),
('p','1','/expertApproval/myDrafts','GET','','',''),
('p','1','/expertApproval/pendingOrgReview','GET','','',''),
('p','1','/expertApproval/pendingCityReview','GET','','',''),
('p','1','/expertApproval/getApprovalLogList','GET','','',''),

('p','1','/expertDatabase/recomputeScore','POST','','',''),
('p','1','/expertDatabase/search','GET','','',''),
('p','1','/expertDatabase/exportSearchResults','GET','','',''),

('p','1','/expertProfile/downloadImportTemplate','GET','','',''),
('p','1','/expertProfile/importBatch','POST','','',''),

('p','1','/expertDatabase/dashboardStats','GET','','','');
