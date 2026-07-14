package ExpertDatabase

import (
	v1 "github.com/flipped-aurora/gin-vue-admin/server/api/v1"
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type ExpertProfileRouter struct{}

// InitExpertProfileRouter 初始化 专家主档 路由信息
func (s *ExpertProfileRouter) InitExpertProfileRouter(Router *gin.RouterGroup) {
	expertProfileRouter := Router.Group("expertProfile").Use(middleware.OperationRecord())
	expertProfileRouterWithoutRecord := Router.Group("expertProfile")
	expertProfileApi := v1.ApiGroupApp.ExpertDatabaseApiGroup.ExpertProfileApi
	{
		expertProfileRouter.POST("createExpertProfile", expertProfileApi.CreateExpertProfile)             // 新建专家主档
		expertProfileRouter.DELETE("deleteExpertProfile", expertProfileApi.DeleteExpertProfile)           // 删除专家主档
		expertProfileRouter.DELETE("deleteExpertProfileByIds", expertProfileApi.DeleteExpertProfileByIds) // 批量删除专家主档
		expertProfileRouter.PUT("updateExpertProfile", expertProfileApi.UpdateExpertProfile)              // 更新专家主档
	}
	{
		expertProfileRouterWithoutRecord.GET("findExpertProfile", expertProfileApi.FindExpertProfile)       // 根据ID获取专家主档
		expertProfileRouterWithoutRecord.GET("getExpertProfileList", expertProfileApi.GetExpertProfileList) // 分页获取专家主档列表
	}
}
