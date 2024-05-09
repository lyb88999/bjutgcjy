package EngineeringEducationDatabase

import (
	"github.com/flipped-aurora/gin-vue-admin/server/api/v1"
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type InternationalExchangeDatabaseRouter struct {
}

// InitInternationalExchangeDatabaseRouter 初始化 国际交流库 路由信息
func (s *InternationalExchangeDatabaseRouter) InitInternationalExchangeDatabaseRouter(Router *gin.RouterGroup) {
	IEDRouter := Router.Group("IED").Use(middleware.OperationRecord())
	IEDRouterWithoutRecord := Router.Group("IED")
	var IEDApi = v1.ApiGroupApp.EngineeringEducationDatabaseApiGroup.InternationalExchangeDatabaseApi
	{
		IEDRouter.POST("createInternationalExchangeDatabase", IEDApi.CreateInternationalExchangeDatabase)   // 新建国际交流库
		IEDRouter.DELETE("deleteInternationalExchangeDatabase", IEDApi.DeleteInternationalExchangeDatabase) // 删除国际交流库
		IEDRouter.DELETE("deleteInternationalExchangeDatabaseByIds", IEDApi.DeleteInternationalExchangeDatabaseByIds) // 批量删除国际交流库
		IEDRouter.PUT("updateInternationalExchangeDatabase", IEDApi.UpdateInternationalExchangeDatabase)    // 更新国际交流库
	}
	{
		IEDRouterWithoutRecord.GET("findInternationalExchangeDatabase", IEDApi.FindInternationalExchangeDatabase)        // 根据ID获取国际交流库
		IEDRouterWithoutRecord.GET("getInternationalExchangeDatabaseList", IEDApi.GetInternationalExchangeDatabaseList)  // 获取国际交流库列表
	}
}
