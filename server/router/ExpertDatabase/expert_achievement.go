package ExpertDatabase

import (
	v1 "github.com/flipped-aurora/gin-vue-admin/server/api/v1"
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type ExpertAchievementRouter struct{}

// InitExpertAchievementRouter 初始化 专家成果 路由信息
func (s *ExpertAchievementRouter) InitExpertAchievementRouter(Router *gin.RouterGroup) {
	expertAchievementRouter := Router.Group("expertAchievement").Use(middleware.OperationRecord())
	expertAchievementRouterWithoutRecord := Router.Group("expertAchievement")
	expertAchievementApi := v1.ApiGroupApp.ExpertDatabaseApiGroup.ExpertAchievementApi
	{
		expertAchievementRouter.POST("createExpertAchievement", expertAchievementApi.CreateExpertAchievement)             // 新建专家成果
		expertAchievementRouter.DELETE("deleteExpertAchievement", expertAchievementApi.DeleteExpertAchievement)           // 删除专家成果
		expertAchievementRouter.DELETE("deleteExpertAchievementByIds", expertAchievementApi.DeleteExpertAchievementByIds) // 批量删除专家成果
		expertAchievementRouter.PUT("updateExpertAchievement", expertAchievementApi.UpdateExpertAchievement)              // 更新专家成果
	}
	{
		expertAchievementRouterWithoutRecord.GET("findExpertAchievement", expertAchievementApi.FindExpertAchievement)       // 根据ID获取专家成果
		expertAchievementRouterWithoutRecord.GET("getExpertAchievementList", expertAchievementApi.GetExpertAchievementList) // 分页获取专家成果列表
	}
}
