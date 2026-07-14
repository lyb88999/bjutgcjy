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
}
