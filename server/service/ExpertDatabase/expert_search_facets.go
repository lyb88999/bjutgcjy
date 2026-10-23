package ExpertDatabase

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/ExpertDatabase"
	ExpertDatabaseReq "github.com/flipped-aurora/gin-vue-admin/server/model/ExpertDatabase/request"
	ExpertDatabaseRes "github.com/flipped-aurora/gin-vue-admin/server/model/ExpertDatabase/response"
	"gorm.io/gorm"
)

// filterClause 一条筛选子句（SQL 片段 + 占位符参数），structuralClauses 按维度分组返回，
// 分面计数时可以按维度整组地跳过——不是逐个 if 复制黏贴四五遍筛选条件
type filterClause struct {
	sql  string
	args []interface{}
}

// structuralClauses 把 ExpertSearchReq 里跟关键词无关的筛选条件按维度分组：name 自成一组（不参与
// 排除逻辑，任何分面计数都要遵守姓名筛选）、unit（UnitName+OrgId）、discipline、region、
// title（TechTitle+TitleTier）各自一组。prefix 是列名前缀——applyStructuralFilters/浏览模式分页
// 走 db.Model() 不用前缀；单位分面要 LEFT JOIN sys_organizations，expert_profile 和
// sys_organizations 都有同名的 name/status 列，必须显式加 "p." 前缀，不然 MySQL 会报列不明确
func structuralClauses(req ExpertDatabaseReq.ExpertSearchReq, prefix string) map[string][]filterClause {
	col := func(name string) string { return prefix + name }
	clauses := map[string][]filterClause{}

	if req.Name != "" {
		clauses["name"] = []filterClause{{col("name") + " LIKE ?", []interface{}{"%" + req.Name + "%"}}}
	}
	var unit []filterClause
	if req.UnitName != "" {
		unit = append(unit, filterClause{col("unit_name") + " LIKE ?", []interface{}{"%" + req.UnitName + "%"}})
	}
	if req.OrgId != nil {
		unit = append(unit, filterClause{col("org_id") + " = ?", []interface{}{*req.OrgId}})
	}
	if len(unit) > 0 {
		clauses["unit"] = unit
	}
	if req.DisciplineL1 != "" {
		clauses["discipline"] = []filterClause{{col("discipline_l1") + " = ?", []interface{}{req.DisciplineL1}}}
	}
	if req.RegionExpertise != "" {
		clauses["region"] = []filterClause{{col("region_expertise") + " LIKE ?", []interface{}{"%" + req.RegionExpertise + "%"}}}
	}
	var title []filterClause
	if req.TechTitle != "" {
		title = append(title, filterClause{col("tech_title") + " = ?", []interface{}{req.TechTitle}})
	}
	if req.TitleTier != "" {
		title = append(title, filterClause{col("title_tier") + " = ?", []interface{}{req.TitleTier}})
	}
	if len(title) > 0 {
		clauses["title"] = title
	}
	return clauses
}

// applyClausesExcept 把 structuralClauses 里除了 except 那一组之外的全部子句套到 db 上。分面
// 计数算某个维度自己的数字时，要叠上"其它维度当前选中的筛选"，但不能筛自己这一维——不然选中
// 某个单位之后，单位分面里除了它自己全部会归零
func applyClausesExcept(db *gorm.DB, req ExpertDatabaseReq.ExpertSearchReq, prefix string, except string) *gorm.DB {
	for group, clauses := range structuralClauses(req, prefix) {
		if group == except {
			continue
		}
		for _, c := range clauses {
			db = db.Where(c.sql, c.args...)
		}
	}
	return db
}

// applyStructuralFilters 给已发布专家套上全部结构化筛选条件（姓名/单位/学科/职称档位等），
// rankedCandidates（关键词模式+导出）、browseModeQuery（浏览模式分页+导出）、facet 计数的 Total
// 共用同一份，不在三处各写一遍——新增筛选字段只用改这一个函数
func applyStructuralFilters(req ExpertDatabaseReq.ExpertSearchReq) *gorm.DB {
	db := global.GVA_DB.Model(&ExpertDatabase.ExpertProfile{}).Where("status = ?", StatusPublished)
	for _, clauses := range structuralClauses(req, "") {
		for _, c := range clauses {
			db = db.Where(c.sql, c.args...)
		}
	}
	return db
}

// GetFacetCounts 分面检索侧栏各维度计数：所在单位（按 org_id/机构名分组，不按 unit_name 原文——
// 同一个机构名下常年混着好几种写法）、一级学科、职称档位、区域国别研究专长。每个维度的计数只按
// "其它维度当前选中的筛选"收窄（applyClausesExcept 排除自身），不按关键词收窄——关键词相关性
// 要对每个候选人跑一次语义 embedding 匹配，四个维度各跑一次的开销扛不住，所以侧栏数字反映的是
// "这些筛选条件在全库范围的分布"，不是"在关键词命中范围内的分布"，这一点是有意的产品取舍，
// 不是遗漏
func (s *ExpertSearchService) GetFacetCounts(req ExpertDatabaseReq.ExpertSearchReq) (result ExpertDatabaseRes.ExpertFacetCounts, err error) {
	// 单位分面：用 Table() 手写 JOIN 时 GORM 不会像 Model() 那样自动带 deleted_at 过滤，要手动加；
	// COALESCE(p.org_id, 0) 是为了让没关联上单位的少数记录也能被稳定分到同一个 value，避免 NULL
	// 扫描进 Go 的 string 字段报错
	orgDB := applyClausesExcept(
		global.GVA_DB.Table("expert_profile p").
			Select("COALESCE(p.org_id, 0) as value, COALESCE(o.name, '未关联单位') as label, COUNT(*) as count").
			Joins("LEFT JOIN sys_organizations o ON o.id = p.org_id").
			Where("p.status = ? AND p.deleted_at IS NULL", StatusPublished),
		req, "p.", "unit",
	)
	if err = orgDB.Group("p.org_id, label").Order("count DESC").Scan(&result.Org).Error; err != nil {
		return
	}

	// 一级学科分面：空值单独归一桶，跟 dashboard 的分布统计用一样的写法
	disciplineDB := applyClausesExcept(
		global.GVA_DB.Model(&ExpertDatabase.ExpertProfile{}).Where("status = ?", StatusPublished).
			Select("CASE WHEN discipline_l1 = '' THEN '未标注' ELSE discipline_l1 END as label, discipline_l1 as value, count(*) as count"),
		req, "", "discipline",
	)
	if err = disciplineDB.Group("label, value").Order("count DESC").Limit(30).Scan(&result.DisciplineL1).Error; err != nil {
		return
	}

	// 职称档位分面：title_tier 是 RecomputeExpertScore 写入的干净枚举值（院士/教授/副教授/讲师/
	// 未分类），不用像学科那样处理空值分桶
	titleTierDB := applyClausesExcept(
		global.GVA_DB.Model(&ExpertDatabase.ExpertProfile{}).Where("status = ?", StatusPublished).
			Select("title_tier as value, title_tier as label, count(*) as count"),
		req, "", "title",
	)
	if err = titleTierDB.Group("title_tier").Order("count DESC").Scan(&result.TitleTier).Error; err != nil {
		return
	}

	// 区域/国别研究专长：目前全库已发布记录没有一条有值（2026-08 数据现状），真的跑 GROUP BY
	// 只会返回空列表——直接返回空切片，前端按"这个分面维度是空的就不渲染"处理（跟
	// GetScoreColumnAvailability 隐藏全零得分列是同一套逻辑）。等这个字段以后真有数据了，把这里
	// 换成跟学科分面一样的 GROUP BY 查询即可，前端不用改
	result.RegionExpertise = []ExpertDatabaseRes.ExpertFacetCount{}

	if err = applyStructuralFilters(req).Count(&result.Total).Error; err != nil {
		return
	}
	return result, nil
}
