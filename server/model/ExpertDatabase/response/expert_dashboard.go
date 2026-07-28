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

	// 内容质量类指标：跟"库里有多少人"无关，反映已发布数据本身的完整度和含金量
	CoveredOrgCount           int64                  `json:"coveredOrgCount"`           // 有至少一位已发布专家的单位数量
	TaggedRatio               float64                `json:"taggedRatio"`               // 已发布专家里，打过至少一个标签的比例（0~100）
	AvgCompositeScore         float64                `json:"avgCompositeScore"`         // 已发布专家的平均综合排序得分
	SeniorTitleRatio          float64                `json:"seniorTitleRatio"`          // 已发布专家里，教授/研究员级别职称的比例（0~100）
	AchievementLevelBreakdown []ExpertDashboardCount `json:"achievementLevelBreakdown"` // 已发布专家名下研究成果，按级别分布
	AdoptionLevelBreakdown    []ExpertDashboardCount `json:"adoptionLevelBreakdown"`    // 已发布专家名下决策影响记录，按采纳/批示单位级别分布
}
