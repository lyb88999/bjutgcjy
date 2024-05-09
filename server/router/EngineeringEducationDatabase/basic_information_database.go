package EngineeringEducationDatabase

import (
	"github.com/flipped-aurora/gin-vue-admin/server/api/v1"
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type BasicInfomationDatabaseRouter struct {
}

// InitBasicInfomationDatabaseRouter 初始化 基础信息库 路由信息
func (s *BasicInfomationDatabaseRouter) InitBasicInfomationDatabaseRouter(Router *gin.RouterGroup) {
	BIDRouter := Router.Group("BID").Use(middleware.OperationRecord())
	BIDRouterWithoutRecord := Router.Group("BID")
	var BIDApi = v1.ApiGroupApp.EngineeringEducationDatabaseApiGroup.BasicInfomationDatabaseApi
	{
		BIDRouter.POST("createBasicInfomationDatabase", BIDApi.CreateBasicInfomationDatabase)   // 新建基础信息库
		BIDRouter.DELETE("deleteBasicInfomationDatabase", BIDApi.DeleteBasicInfomationDatabase) // 删除基础信息库
		BIDRouter.DELETE("deleteBasicInfomationDatabaseByIds", BIDApi.DeleteBasicInfomationDatabaseByIds) // 批量删除基础信息库
		BIDRouter.PUT("updateBasicInfomationDatabase", BIDApi.UpdateBasicInfomationDatabase)    // 更新基础信息库
	}
	{
		BIDRouterWithoutRecord.GET("findBasicInfomationDatabase", BIDApi.FindBasicInfomationDatabase)        // 根据ID获取基础信息库
		BIDRouterWithoutRecord.GET("getBasicInfomationDatabaseList", BIDApi.GetBasicInfomationDatabaseList)  // 获取基础信息库列表
	}
}
