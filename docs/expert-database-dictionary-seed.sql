-- 专家库推荐排序权重 + 审核开关 字典种子（幂等）
--
-- 背景：server/service/ExpertDatabase/expert_score.go 的算分逻辑会优先从前 5 个字典类型读取
-- 权重，字典不存在时才退回代码里写死的默认值。这份种子把默认值实体化成"系统工具-字典管理"
-- 里可见、可调整的真实字典，这样检索页面"权重可在字典管理里调整"这句提示才是真的——
-- 在此之前，管理员打开字典管理会发现根本没有这些字典类型。
--
-- 调整方法：以后要改权重，直接去"系统工具-字典管理"里编辑对应字典的明细项即可，
-- 不需要改代码、不需要重启服务（下次算分请求会实时读取最新值）。
--
-- 每个字典明细的 value 和 extend 都存了同一个权重数字：value 是字典本身的常规字段（界面上
-- 展示/排序用），extend 是 server/service/ExpertDatabase/expert_score.go 实际解析读取的字段
-- （dictWeights() 优先 strconv.ParseFloat(extend)，取不到才退回 value）。
--
-- 最后一个字典类型 expert_review_switch 是审核开关：只有一条明细"启用审核"，value=1 表示
-- 批量导入/手动新增/提交审核 要走完整的单位->市级三级审核流程（默认），改成 value=0 就变成
-- 直接落库为已发布、跳过整个审核流程（见 server/service/ExpertDatabase/expert_approval.go 的
-- reviewRequired()）。字典/明细缺失时按 value=1（启用审核）处理，不会因为没配置就意外放开。

INSERT IGNORE INTO sys_dictionaries (created_at, updated_at, name, type, status, `desc`) VALUES
  (NOW(3), NOW(3), '专家库-成果级别权重', 'expert_achievement_level', 1, '成果/课题级别 -> 权重，用于成果得分 = 级别权重 x 该级别篇数求和'),
  (NOW(3), NOW(3), '专家库-采纳单位级别权重', 'expert_adoption_level', 1, '决策影响记录的采纳/批示单位级别 -> 权重'),
  (NOW(3), NOW(3), '专家库-职称权重', 'expert_title_level', 1, '专业技术职称 -> 权重，用于综合排序里的职称分项'),
  (NOW(3), NOW(3), '专家库-社会贡献权重', 'expert_social_weight', 1, '学术兼职每条 / 荣誉称号是否非空 的加分权重'),
  (NOW(3), NOW(3), '专家库-综合排序权重', 'expert_ranking_weight', 1, '综合排序四个分项 achievement/influence/title/social 各自的权重 w1..w4'),
  (NOW(3), NOW(3), '专家库-审核开关', 'expert_review_switch', 1, '批量导入/手动新增/提交审核 是否要走单位->市级三级审核流程；明细 value=1 表示启用审核（默认），value=0 表示关闭审核、直接落库为已发布');

INSERT IGNORE INTO sys_dictionary_details (created_at, updated_at, label, value, extend, status, sort, sys_dictionary_id)
SELECT NOW(3), NOW(3), d.label, d.value, d.value, 1, d.sort, sd.id
FROM sys_dictionaries sd
JOIN (
  SELECT 'expert_achievement_level' AS type, '国家级' AS label, 10 AS value, 1 AS sort UNION ALL
  SELECT 'expert_achievement_level', '省部级', 5, 2 UNION ALL
  SELECT 'expert_achievement_level', '厅局级', 3, 3 UNION ALL
  SELECT 'expert_achievement_level', '一般级', 1, 4 UNION ALL

  SELECT 'expert_adoption_level', '国家级', 10, 1 UNION ALL
  SELECT 'expert_adoption_level', '中央', 10, 2 UNION ALL
  SELECT 'expert_adoption_level', '省部级', 5, 3 UNION ALL
  SELECT 'expert_adoption_level', '厅局级', 3, 4 UNION ALL
  SELECT 'expert_adoption_level', '区县级', 1, 5 UNION ALL

  SELECT 'expert_title_level', '教授', 5, 1 UNION ALL
  SELECT 'expert_title_level', '研究员', 5, 2 UNION ALL
  SELECT 'expert_title_level', '副教授', 3, 3 UNION ALL
  SELECT 'expert_title_level', '副研究员', 3, 4 UNION ALL
  SELECT 'expert_title_level', '讲师', 1, 5 UNION ALL
  SELECT 'expert_title_level', '助理研究员', 1, 6 UNION ALL

  SELECT 'expert_social_weight', 'academic_position', 1, 1 UNION ALL
  SELECT 'expert_social_weight', 'honor_title', 3, 2 UNION ALL

  SELECT 'expert_ranking_weight', 'achievement', 1, 1 UNION ALL
  SELECT 'expert_ranking_weight', 'influence', 1, 2 UNION ALL
  SELECT 'expert_ranking_weight', 'title', 1, 3 UNION ALL
  SELECT 'expert_ranking_weight', 'social', 1, 4 UNION ALL

  SELECT 'expert_review_switch', '启用审核', 1, 1
) d ON d.type = sd.type
WHERE sd.type LIKE 'expert_%'
  AND sd.deleted_at IS NULL
  AND NOT EXISTS (
    SELECT 1 FROM sys_dictionary_details sdd
    WHERE sdd.sys_dictionary_id = sd.id AND sdd.label = d.label AND sdd.deleted_at IS NULL
  );
