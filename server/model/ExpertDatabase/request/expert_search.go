package request

// ExpertSearchReq 专家综合推荐排序检索入参：先按筛选条件圈定候选集，再按关键词相关性 + 缓存得分排序
type ExpertSearchReq struct {
	Keyword         string `json:"keyword" form:"keyword"`
	Name            string `json:"name" form:"name"`
	UnitName        string `json:"unitName" form:"unitName"`
	DisciplineL1    string `json:"disciplineL1" form:"disciplineL1"`
	RegionExpertise string `json:"regionExpertise" form:"regionExpertise"`
	TechTitle       string `json:"techTitle" form:"techTitle"`
	// SortBy 只在没填 Keyword 的默认浏览场景生效："updatedAt"（最近更新在前）、"name"（姓名排序），
	// 留空则按综合实力排序。带关键词检索时必须按相关性排序，不接受这个参数覆盖——不然搜索结果的
	// 排序跟检索意图脱节
	SortBy   string `json:"sortBy" form:"sortBy"`
	Page     int    `json:"page" form:"page"`
	PageSize int    `json:"pageSize" form:"pageSize"`
}
