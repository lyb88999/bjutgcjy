package task

import (
	"github.com/flipped-aurora/gin-vue-admin/server/service"
)

// RecomputeExpertScores 兜底批量重算所有已发布专家的推荐排序得分缓存
// 正常情况下成果/决策影响记录变更、审核发布时都会即时触发单个专家的重算，
// 这里作为定时兜底（例如手工改过字典权重后，需要让所有专家的缓存分重新生效）
func RecomputeExpertScores() error {
	return service.ServiceGroupApp.ExpertDatabaseServiceGroup.ExpertScoreService.RecomputeAllPublishedExpertScores()
}
