package EngineeringEducationDatabase

import (
	"github.com/flipped-aurora/gin-vue-admin/server/api/v1"
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type ConditionalGuaranteeDatabaseRouter struct {
}

// InitConditionalGuaranteeDatabaseRouter 初始化 条件保障库 路由信息
func (s *ConditionalGuaranteeDatabaseRouter) InitConditionalGuaranteeDatabaseRouter(Router *gin.RouterGroup) {
	CGDRouter := Router.Group("CGD").Use(middleware.OperationRecord())
	CGDRouterWithoutRecord := Router.Group("CGD")
	var CGDApi = v1.ApiGroupApp.EngineeringEducationDatabaseApiGroup.ConditionalGuaranteeDatabaseApi
	{
		CGDRouter.POST("createConditionalGuaranteeDatabase", CGDApi.CreateConditionalGuaranteeDatabase)   // 新建条件保障库
		CGDRouter.DELETE("deleteConditionalGuaranteeDatabase", CGDApi.DeleteConditionalGuaranteeDatabase) // 删除条件保障库
		CGDRouter.DELETE("deleteConditionalGuaranteeDatabaseByIds", CGDApi.DeleteConditionalGuaranteeDatabaseByIds) // 批量删除条件保障库
		CGDRouter.PUT("updateConditionalGuaranteeDatabase", CGDApi.UpdateConditionalGuaranteeDatabase)    // 更新条件保障库
	}
	{
		CGDRouterWithoutRecord.GET("findConditionalGuaranteeDatabase", CGDApi.FindConditionalGuaranteeDatabase)        // 根据ID获取条件保障库
		CGDRouterWithoutRecord.GET("getConditionalGuaranteeDatabaseList", CGDApi.GetConditionalGuaranteeDatabaseList)  // 获取条件保障库列表
	}
}
