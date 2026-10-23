package request

// AdvancedSearchRow 高级检索一行；Connector 取值 "and"/"or"/"not"，首行的 Connector 会被忽略
// （首行只是起点，没有更早的累计结果可供连接）
type AdvancedSearchRow struct {
	Connector string `json:"connector"`
	Term      string `json:"term"`
}

// ExpertSearchReq 专家综合推荐排序检索入参：先按筛选条件圈定候选集，再按关键词相关性 + 缓存得分排序
type ExpertSearchReq struct {
	Keyword string `json:"keyword" form:"keyword"`
	// AdvancedRows 高级检索（多行且/或/不含）整体 JSON 数组的字符串形式——GET 请求走
	// ShouldBindQuery，gin 的 form 绑定不支持 slice-of-struct，一个 JSON 字符串参数最简单、
	// 两端都好写。非空且解析出≥1个有效行时完全替代 Keyword，见 rankedCandidates
	AdvancedRows string `json:"advancedRows" form:"advancedRows"`
	Name         string `json:"name" form:"name"`
	// UnitName 旧版自由文本单位筛选，分面检索改版后前端不再发这个字段，改发 OrgId；字段留着
	// 只为了后端向后兼容（万一有别的调用方还在用），两个字段可以同时生效（AND 关系）
	UnitName string `json:"unitName" form:"unitName"`
	// OrgId 分面检索"所在单位"筛选，对应 sys_organizations.id，比 UnitName 干净——一个机构名下
	// 常年混着好几种 unit_name 写法（比如"中国社会科学院大学"和"中国社会科学院"）
	OrgId           *uint  `json:"orgId" form:"orgId"`
	DisciplineL1    string `json:"disciplineL1" form:"disciplineL1"`
	RegionExpertise string `json:"regionExpertise" form:"regionExpertise"`
	// TechTitle 旧版自由文本职称筛选，同 UnitName 一样留作后端兼容，新前端改发 TitleTier
	TechTitle string `json:"techTitle" form:"techTitle"`
	// TitleTier 分面检索"职称档位"筛选，取值是 matchTitleTier 判出的档位标签
	// （院士/教授/副教授/讲师/未分类），不是 TechTitle 原文——TechTitle 原文有 313 种自由写法，
	// 没法直接当筛选值用
	TitleTier string `json:"titleTier" form:"titleTier"`
	// SortBy 只在没填 Keyword 的默认浏览场景生效："updatedAt"（最近更新在前）、"name"（姓名排序），
	// 留空则按综合实力排序。带关键词检索时必须按相关性排序，不接受这个参数覆盖——不然搜索结果的
	// 排序跟检索意图脱节
	SortBy   string `json:"sortBy" form:"sortBy"`
	Page     int    `json:"page" form:"page"`
	PageSize int    `json:"pageSize" form:"pageSize"`
}
