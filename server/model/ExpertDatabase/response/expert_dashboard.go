package response

// ExpertDashboardCount 通用的"标签 -> 数量"统计项，用于状态/单位/学科分布这几类柱状图/饼图
type ExpertDashboardCount struct {
	Label string `json:"label"`
	Count int64  `json:"count"`
}

// ExpertDashboardTrendPoint 按日期统计的新增数量，用于趋势折线图
type ExpertDashboardTrendPoint struct {
	Date  string `json:"date"`
	Count int64  `json:"count"`
}

// ExpertDashboardStats 专家库统计概览；单位审核员看到的是本单位范围内的统计，
// 管理员/市级审核员看到的是全库统计（见 service 层 GetDashboardStats 的角色分支）
type ExpertDashboardStats struct {
	TotalPublished      int64                       `json:"totalPublished"`
	TotalDraft          int64                       `json:"totalDraft"`
	PendingOrgReview    int64                       `json:"pendingOrgReview"`
	PendingCityReview   int64                       `json:"pendingCityReview"`
	ReviewBypassed      int64                       `json:"reviewBypassed"`
	StatusBreakdown     []ExpertDashboardCount      `json:"statusBreakdown"`
	OrgBreakdown        []ExpertDashboardCount      `json:"orgBreakdown"`
	DisciplineBreakdown []ExpertDashboardCount      `json:"disciplineBreakdown"`
	RecentTrend         []ExpertDashboardTrendPoint `json:"recentTrend"`
}
