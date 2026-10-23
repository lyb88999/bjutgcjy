-- Apply to the development database after the standard expert database permission seeds.
-- Grant the assistant only to roles already allowed to search published experts.
START TRANSACTION;

INSERT INTO sys_apis (created_at, updated_at, path, description, api_group, method)
SELECT NOW(3), NOW(3), '/expertDatabase/agent/status', '专家遴选助手状态', 'ExpertAgent', 'GET'
WHERE NOT EXISTS (SELECT 1 FROM sys_apis WHERE path = '/expertDatabase/agent/status' AND method = 'GET' AND deleted_at IS NULL);

INSERT INTO sys_apis (created_at, updated_at, path, description, api_group, method)
SELECT NOW(3), NOW(3), '/expertDatabase/agent/chat', '专家智能遴选', 'ExpertAgent', 'POST'
WHERE NOT EXISTS (SELECT 1 FROM sys_apis WHERE path = '/expertDatabase/agent/chat' AND method = 'POST' AND deleted_at IS NULL);

INSERT IGNORE INTO casbin_rule (ptype, v0, v1, v2, v3, v4, v5)
SELECT DISTINCT 'p', v0, '/expertDatabase/agent/status', 'GET', '', '', ''
FROM casbin_rule WHERE ptype = 'p' AND v1 = '/expertDatabase/search' AND v2 = 'GET';

INSERT IGNORE INTO casbin_rule (ptype, v0, v1, v2, v3, v4, v5)
SELECT DISTINCT 'p', v0, '/expertDatabase/agent/chat', 'POST', '', '', ''
FROM casbin_rule WHERE ptype = 'p' AND v1 = '/expertDatabase/search' AND v2 = 'GET';

COMMIT;
