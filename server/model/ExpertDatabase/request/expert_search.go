package request

// ExpertSearchReq 专家综合推荐排序检索入参：先按筛选条件圈定候选集，再按关键词相关性 + 缓存得分排序
type ExpertSearchReq struct {
	Keyword         string `json:"keyword" form:"keyword"`
	DisciplineL1    string `json:"disciplineL1" form:"disciplineL1"`
	RegionExpertise string `json:"regionExpertise" form:"regionExpertise"`
	TechTitle       string `json:"techTitle" form:"techTitle"`
	Page            int    `json:"page" form:"page"`
	PageSize        int    `json:"pageSize" form:"pageSize"`
}
