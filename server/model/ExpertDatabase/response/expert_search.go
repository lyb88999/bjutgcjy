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
