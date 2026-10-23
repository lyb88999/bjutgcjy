package response

type ExpertAgentStatus struct {
	Configured bool   `json:"configured"`
	Model      string `json:"model"`
}

type ExpertAgentPlan struct {
	Keywords          []string `json:"keywords"`
	Name              string   `json:"name"`
	UnitName          string   `json:"unitName"`
	DisciplineL1      string   `json:"disciplineL1"`
	TitleTier         string   `json:"titleTier"`
	Count             int      `json:"count"`
	MaxPerUnit        int      `json:"maxPerUnit"`
	ExcludedExpertIDs []uint   `json:"excludedExpertIds"`
	Diversify         bool     `json:"diversify"`
}

type ExpertAgentEvidence struct {
	Kind     string `json:"kind"`
	Title    string `json:"title"`
	SourceID uint   `json:"sourceId,omitempty"`
	URL      string `json:"url,omitempty"`
}

// 候选用独立DTO，不能嵌入含联系方式和其他个人信息的ExpertProfile。
type ExpertAgentCandidate struct {
	OrgID              uint                  `json:"-"`
	TopicRelevance     float64               `json:"-"` // 内部主题相关性，只用于遴选，不作为概率展示
	ID                 uint                  `json:"id"`
	Name               string                `json:"name"`
	UnitName           string                `json:"unitName"`
	Department         string                `json:"department"`
	TechTitle          string                `json:"techTitle"`
	DisciplineL1       string                `json:"disciplineL1"`
	ResearchDirections string                `json:"researchDirections"`
	ProfileLink        string                `json:"profileLink"`
	MatchReason        string                `json:"matchReason"`
	Score              float64               `json:"score"`
	Evidence           []ExpertAgentEvidence `json:"evidence"`
	Reason             string                `json:"reason"`
}

type ExpertAgentDistribution struct {
	Label string `json:"label"`
	Count int    `json:"count"`
}

type ExpertAgentStatistics struct {
	CandidateCount  int                       `json:"candidateCount"`
	SelectedCount   int                       `json:"selectedCount"`
	UnitCount       int                       `json:"unitCount"`
	DisciplineCount int                       `json:"disciplineCount"`
	Units           []ExpertAgentDistribution `json:"units"`
	Disciplines     []ExpertAgentDistribution `json:"disciplines"`
}

type ExpertAgentStep struct {
	Label  string `json:"label"`
	Detail string `json:"detail"`
}

type ExpertAgentChatResult struct {
	SessionID       string                 `json:"sessionId"`
	Reply           string                 `json:"reply"`
	Plan            ExpertAgentPlan        `json:"plan"`
	Candidates      []ExpertAgentCandidate `json:"candidates"`
	Recommendations []ExpertAgentCandidate `json:"recommendations"`
	Statistics      ExpertAgentStatistics  `json:"statistics"`
	Steps           []ExpertAgentStep      `json:"steps"`
	Warnings        []string               `json:"warnings"`
}
