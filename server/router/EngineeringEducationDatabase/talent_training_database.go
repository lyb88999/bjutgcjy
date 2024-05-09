package EngineeringEducationDatabase

import (
	"github.com/flipped-aurora/gin-vue-admin/server/api/v1"
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type TalentTrainingDatabaseRouter struct {
}

// InitTalentTrainingDatabaseRouter 初始化 人才培养库 路由信息
func (s *TalentTrainingDatabaseRouter) InitTalentTrainingDatabaseRouter(Router *gin.RouterGroup) {
	TTDRouter := Router.Group("TTD").Use(middleware.OperationRecord())
	TTDRouterWithoutRecord := Router.Group("TTD")
	var TTDApi = v1.ApiGroupApp.EngineeringEducationDatabaseApiGroup.TalentTrainingDatabaseApi
	{
		TTDRouter.POST("createTalentTrainingDatabase", TTDApi.CreateTalentTrainingDatabase)   // 新建人才培养库
		TTDRouter.DELETE("deleteTalentTrainingDatabase", TTDApi.DeleteTalentTrainingDatabase) // 删除人才培养库
		TTDRouter.DELETE("deleteTalentTrainingDatabaseByIds", TTDApi.DeleteTalentTrainingDatabaseByIds) // 批量删除人才培养库
		TTDRouter.PUT("updateTalentTrainingDatabase", TTDApi.UpdateTalentTrainingDatabase)    // 更新人才培养库
	}
	{
		TTDRouterWithoutRecord.GET("findTalentTrainingDatabase", TTDApi.FindTalentTrainingDatabase)        // 根据ID获取人才培养库
		TTDRouterWithoutRecord.GET("getTalentTrainingDatabaseList", TTDApi.GetTalentTrainingDatabaseList)  // 获取人才培养库列表
	}
}
