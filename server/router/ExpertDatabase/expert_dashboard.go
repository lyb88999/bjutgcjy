package ExpertDatabase

import (
	v1 "github.com/flipped-aurora/gin-vue-admin/server/api/v1"
	"github.com/gin-gonic/gin"
)

type ExpertDashboardRouter struct{}

// InitExpertDashboardRouter 初始化 专家库统计概览 路由信息
func (s *ExpertDashboardRouter) InitExpertDashboardRouter(Router *gin.RouterGroup) {
	expertDatabaseRouterWithoutRecord := Router.Group("expertDatabase")
	expertDashboardApi := v1.ApiGroupApp.ExpertDatabaseApiGroup.ExpertDashboardApi
	{
		expertDatabaseRouterWithoutRecord.GET("dashboardStats", expertDashboardApi.GetDashboardStats) // 统计概览，只读，不用记操作日志
	}
}
