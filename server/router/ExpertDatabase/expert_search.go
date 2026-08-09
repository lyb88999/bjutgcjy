package ExpertDatabase

import (
	v1 "github.com/flipped-aurora/gin-vue-admin/server/api/v1"
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type ExpertSearchRouter struct{}

// InitExpertSearchRouter 初始化 专家综合推荐排序检索 路由信息
func (s *ExpertSearchRouter) InitExpertSearchRouter(Router *gin.RouterGroup) {
	expertDatabaseRouter := Router.Group("expertDatabase").Use(middleware.OperationRecord())
	expertDatabaseRouterWithoutRecord := Router.Group("expertDatabase")
	expertSearchApi := v1.ApiGroupApp.ExpertDatabaseApiGroup.ExpertSearchApi
	{
		expertDatabaseRouter.POST("recomputeScore", expertSearchApi.RecomputeExpertScore) // 手动重算专家得分
	}
	{
		expertDatabaseRouterWithoutRecord.GET("search", expertSearchApi.SearchExpert)                                       // 综合推荐排序检索
		expertDatabaseRouterWithoutRecord.GET("exportSearchResults", expertSearchApi.ExportSearchResults)                   // 导出检索结果
		expertDatabaseRouterWithoutRecord.GET("scoreColumnAvailability", expertSearchApi.GetScoreColumnAvailability)        // 分项得分是否有数据
	}
}
