package response

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/ExpertDatabase"
)

// ExpertSearchItem 检索结果项：在专家主档基础上附带本次检索的相关性系数与实时综合得分，
// 让使用者能看懂"为什么这个专家排在前面"
type ExpertSearchItem struct {
	ExpertDatabase.ExpertProfile
	Relevance     float64 `json:"relevance"`
	RealtimeScore float64 `json:"realtimeScore"`
	// MatchReason 相关性最主要是由哪一条语料贡献的（比如具体是哪条研究成果、还是研究方向本身），
	// 关键词为空的默认浏览场景下留空。给用户一个"为什么排在前面"的直观依据，而不是只有一个百分比
	MatchReason string `json:"matchReason"`
}

// ExpertScoreColumnAvailability 各分项得分在全库已发布专家里是否至少有一个人非零——批量导入的
// 花名册类数据往往还没有决策影响/学术兼职记录，这种情况下相应分项对谁都是 0，界面上继续摆一列
// 清一色的 0 没有信息量，前端拿这个标记决定要不要显示对应的列/卡片
type ExpertScoreColumnAvailability struct {
	HasAchievement bool `json:"hasAchievement"`
	HasInfluence   bool `json:"hasInfluence"`
	HasSocial      bool `json:"hasSocial"`
}

// ExpertFacetCount 分面检索侧栏里某个维度下的一个可选值
type ExpertFacetCount struct {
	Value string `json:"value"` // 单位分面是 org_id 的字符串形式；其余维度是原始文本值
	Label string `json:"label"`
	Count int64  `json:"count"`
}

// ExpertFacetCounts 分面检索侧栏一次性拿到的全部维度计数。每个维度的计数口径：只按"结构化
// 筛选"（其余分面维度当前选中的值）收窄，不按关键词收窄——关键词相关性要对每个候选人跑一次
// embedding 匹配，四个维度各跑一次的开销扛不住，所以侧栏数字反映的是"这些筛选条件在全库范围
// 的分布"，不是"在关键词命中范围内的分布"
type ExpertFacetCounts struct {
	Org             []ExpertFacetCount `json:"org"`
	DisciplineL1    []ExpertFacetCount `json:"disciplineL1"`
	TitleTier       []ExpertFacetCount `json:"titleTier"`
	RegionExpertise []ExpertFacetCount `json:"regionExpertise"`
	// Total 应用当前全部结构化筛选后的候选总数（不含关键词）
	Total int64 `json:"total"`
}
