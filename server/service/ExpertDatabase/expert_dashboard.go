package ExpertDatabase

import (
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/ExpertDatabase"
	ExpertDatabaseResp "github.com/flipped-aurora/gin-vue-admin/server/model/ExpertDatabase/response"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
	"gorm.io/gorm"
)

type ExpertDashboardService struct{}

// statusOrder/statusLabel 固定展示顺序，即使某个状态当前一条记录都没有，也要在图表里占一个位置（0），
// 不然图表每次刷新柱子的个数和顺序都会跳动
var statusOrder = []string{StatusDraft, StatusPendingOrgReview, StatusOrgRejected, StatusPendingCityReview, StatusCityRejected, StatusPublished}
var statusLabel = map[string]string{
	StatusDraft: "草稿", StatusPendingOrgReview: "待单位审核", StatusOrgRejected: "单位已退回",
	StatusPendingCityReview: "待市级审核", StatusCityRejected: "市级已退回", StatusPublished: "已发布",
}

// GetDashboardStats 专家库统计概览：管理员和市级审核员看全库范围，单位审核员只看本单位范围内的数据
// （呼应 GetExpertProfileInfoList 的角色数据域收敛，不能让单位审核员通过统计图表侧面看到别的单位的数据）
func (s *ExpertDashboardService) GetDashboardStats(operatorID uint) (ExpertDatabaseResp.ExpertDashboardStats, error) {
	var stats ExpertDatabaseResp.ExpertDashboardStats

	var operator system.SysUser
	if err := global.GVA_DB.Where("id = ?", operatorID).First(&operator).Error; err != nil {
		return stats, err
	}
	orgScope := operator.AuthorityId == authorityOrgReviewer && operator.OrgId != nil

	baseQuery := func() *gorm.DB {
		db := global.GVA_DB.Model(&ExpertDatabase.ExpertProfile{})
		if orgScope {
			db = db.Where("org_id = ?", *operator.OrgId)
		}
		return db
	}

	// 状态分布
	type statusCountRow struct {
		Status string
		Count  int64
	}
	var statusRows []statusCountRow
	if err := baseQuery().Select("status, count(*) as count").Group("status").Scan(&statusRows).Error; err != nil {
		return stats, err
	}
	statusCounts := make(map[string]int64, len(statusRows))
	for _, r := range statusRows {
		statusCounts[r.Status] = r.Count
	}
	for _, st := range statusOrder {
		stats.StatusBreakdown = append(stats.StatusBreakdown, ExpertDatabaseResp.ExpertDashboardCount{
			Label: statusLabel[st], Count: statusCounts[st],
		})
	}
	stats.TotalDraft = statusCounts[StatusDraft]
	stats.PendingOrgReview = statusCounts[StatusPendingOrgReview]
	stats.PendingCityReview = statusCounts[StatusPendingCityReview]
	stats.TotalPublished = statusCounts[StatusPublished]

	// 免审核发布数量（审核开关关闭期间绕过审核直接发布、且尚未被真正的市级审核覆盖过的记录）
	if err := baseQuery().Where("review_bypassed = ?", true).Count(&stats.ReviewBypassed).Error; err != nil {
		return stats, err
	}

	// 按单位分布（只统计已发布的，草稿/审核中的不适合公开到统计图表里）；单位审核员本来就已经被
	// baseQuery 收敛到自己单位，这里自然只会算出一行，不会看到其他单位
	// 用 Table() 手写 join 时 GORM 不会像 Model() 那样自动带上软删除过滤，要手动加，
	// 不然已经被删除（deleted_at 不为空）的历史记录会被计入统计
	orgDB := global.GVA_DB.Table("expert_profile p").
		Select("COALESCE(o.name, '未关联单位') as label, COUNT(*) as count").
		Joins("LEFT JOIN sys_organizations o ON o.id = p.org_id").
		Where("p.status = ? AND p.deleted_at IS NULL", StatusPublished)
	if orgScope {
		orgDB = orgDB.Where("p.org_id = ?", *operator.OrgId)
	}
	if err := orgDB.Group("label").Order("count DESC").Limit(10).Scan(&stats.OrgBreakdown).Error; err != nil {
		return stats, err
	}

	// 按一级学科分布（同样只看已发布的）
	disciplineDB := baseQuery().Where("status = ?", StatusPublished).
		Select("CASE WHEN discipline_l1 = '' THEN '未标注' ELSE discipline_l1 END as label, count(*) as count").
		Group("label").Order("count DESC").Limit(10)
	if err := disciplineDB.Scan(&stats.DisciplineBreakdown).Error; err != nil {
		return stats, err
	}

	// 最近 30 天新增趋势（按创建日期计数，没有新增的当天补 0，避免图表因为缺天数而变形）
	since := time.Now().AddDate(0, 0, -29)
	type trendRow struct {
		Date  string
		Count int64
	}
	var trendRows []trendRow
	if err := baseQuery().Where("created_at >= ?", since).
		Select("DATE(created_at) as date, count(*) as count").
		Group("date").Order("date").Scan(&trendRows).Error; err != nil {
		return stats, err
	}
	trendMap := make(map[string]int64, len(trendRows))
	for _, r := range trendRows {
		trendMap[r.Date] = r.Count
	}
	for i := 0; i < 30; i++ {
		day := since.AddDate(0, 0, i).Format("2006-01-02")
		stats.RecentTrend = append(stats.RecentTrend, ExpertDatabaseResp.ExpertDashboardTrendPoint{
			Date: day, Count: trendMap[day],
		})
	}

	// 覆盖单位数：有至少一位已发布专家的单位数量（org_id 为空的不算数，那些还没关联到具体单位）
	if err := baseQuery().Where("status = ? AND org_id IS NOT NULL", StatusPublished).
		Distinct("org_id").Count(&stats.CoveredOrgCount).Error; err != nil {
		return stats, err
	}

	// 平均综合排序得分、高级职称（教授/研究员）占比，都是衡量"已发布数据含金量"的指标，
	// 跟"库里有多少人"是两件事——哪怕只有几十条，也应该看得出质量高不高
	type scoreAggRow struct {
		AvgScore    float64
		SeniorCount int64
	}
	var scoreAgg scoreAggRow
	if err := baseQuery().Where("status = ?", StatusPublished).
		Select("COALESCE(AVG(composite_score), 0) as avg_score, SUM(CASE WHEN tech_title IN ('教授','研究员') THEN 1 ELSE 0 END) as senior_count").
		Scan(&scoreAgg).Error; err != nil {
		return stats, err
	}
	stats.AvgCompositeScore = scoreAgg.AvgScore
	if stats.TotalPublished > 0 {
		stats.SeniorTitleRatio = float64(scoreAgg.SeniorCount) / float64(stats.TotalPublished) * 100
	}

	// 标签覆盖率：打过标签才谈得上被检索关键词准确命中，这个比例低说明检索大概率会失灵
	var taggedCount int64
	tagDB := global.GVA_DB.Table("expert_profile p").
		Joins("JOIN expert_tag_relation r ON r.expert_id = p.id AND r.deleted_at IS NULL").
		Where("p.status = ? AND p.deleted_at IS NULL", StatusPublished)
	if orgScope {
		tagDB = tagDB.Where("p.org_id = ?", *operator.OrgId)
	}
	if err := tagDB.Distinct("p.id").Count(&taggedCount).Error; err != nil {
		return stats, err
	}
	if stats.TotalPublished > 0 {
		stats.TaggedRatio = float64(taggedCount) / float64(stats.TotalPublished) * 100
	}

	// 已发布专家名下研究成果按级别分布，看看"国家级"这种含金量高的成果占比多不多
	achievementDB := global.GVA_DB.Table("expert_achievement a").
		Select("CASE WHEN a.level IS NULL OR a.level = '' THEN '未标注' ELSE a.level END as label, COUNT(*) as count").
		Joins("JOIN expert_profile p ON p.id = a.expert_id").
		Where("p.status = ? AND p.deleted_at IS NULL AND a.deleted_at IS NULL", StatusPublished)
	if orgScope {
		achievementDB = achievementDB.Where("p.org_id = ?", *operator.OrgId)
	}
	if err := achievementDB.Group("label").Order("count DESC").Scan(&stats.AchievementLevelBreakdown).Error; err != nil {
		return stats, err
	}

	// 已发布专家名下决策影响记录，按采纳/批示单位级别分布，跟成果级别分布是同一个思路
	adoptionDB := global.GVA_DB.Table("expert_adoption_record r").
		Select("CASE WHEN r.adopting_unit_level IS NULL OR r.adopting_unit_level = '' THEN '未标注' ELSE r.adopting_unit_level END as label, COUNT(*) as count").
		Joins("JOIN expert_profile p ON p.id = r.expert_id").
		Where("p.status = ? AND p.deleted_at IS NULL AND r.deleted_at IS NULL", StatusPublished)
	if orgScope {
		adoptionDB = adoptionDB.Where("p.org_id = ?", *operator.OrgId)
	}
	if err := adoptionDB.Group("label").Order("count DESC").Scan(&stats.AdoptionLevelBreakdown).Error; err != nil {
		return stats, err
	}

	return stats, nil
}
