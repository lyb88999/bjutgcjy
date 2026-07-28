package ExpertDatabase

import (
	v1 "github.com/flipped-aurora/gin-vue-admin/server/api/v1"
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type ExpertAcademicPositionRouter struct{}

// InitExpertAcademicPositionRouter 初始化 专家学术兼职 路由信息
func (s *ExpertAcademicPositionRouter) InitExpertAcademicPositionRouter(Router *gin.RouterGroup) {
	expertAcademicPositionRouter := Router.Group("expertAcademicPosition").Use(middleware.OperationRecord())
	expertAcademicPositionRouterWithoutRecord := Router.Group("expertAcademicPosition")
	expertAcademicPositionApi := v1.ApiGroupApp.ExpertDatabaseApiGroup.ExpertAcademicPositionApi
	{
		expertAcademicPositionRouter.POST("createExpertAcademicPosition", expertAcademicPositionApi.CreateExpertAcademicPosition)             // 新建专家学术兼职
		expertAcademicPositionRouter.DELETE("deleteExpertAcademicPosition", expertAcademicPositionApi.DeleteExpertAcademicPosition)           // 删除专家学术兼职
		expertAcademicPositionRouter.DELETE("deleteExpertAcademicPositionByIds", expertAcademicPositionApi.DeleteExpertAcademicPositionByIds) // 批量删除专家学术兼职
		expertAcademicPositionRouter.PUT("updateExpertAcademicPosition", expertAcademicPositionApi.UpdateExpertAcademicPosition)              // 更新专家学术兼职
	}
	{
		expertAcademicPositionRouterWithoutRecord.GET("findExpertAcademicPosition", expertAcademicPositionApi.FindExpertAcademicPosition)       // 根据ID获取专家学术兼职
		expertAcademicPositionRouterWithoutRecord.GET("getExpertAcademicPositionList", expertAcademicPositionApi.GetExpertAcademicPositionList) // 分页获取专家学术兼职列表
	}
}
