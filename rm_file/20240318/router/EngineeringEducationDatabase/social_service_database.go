package EngineeringEducationDatabase

import (
	"github.com/flipped-aurora/gin-vue-admin/server/api/v1"
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type SocialServiceDatabaseRouter struct {
}

// InitSocialServiceDatabaseRouter 初始化 社会服务库 路由信息
func (s *SocialServiceDatabaseRouter) InitSocialServiceDatabaseRouter(Router *gin.RouterGroup) {
	SSDRouter := Router.Group("SSD").Use(middleware.OperationRecord())
	SSDRouterWithoutRecord := Router.Group("SSD")
	var SSDApi = v1.ApiGroupApp.EngineeringEducationDatabaseApiGroup.SocialServiceDatabaseApi
	{
		SSDRouter.POST("createSocialServiceDatabase", SSDApi.CreateSocialServiceDatabase)   // 新建社会服务库
		SSDRouter.DELETE("deleteSocialServiceDatabase", SSDApi.DeleteSocialServiceDatabase) // 删除社会服务库
		SSDRouter.DELETE("deleteSocialServiceDatabaseByIds", SSDApi.DeleteSocialServiceDatabaseByIds) // 批量删除社会服务库
		SSDRouter.PUT("updateSocialServiceDatabase", SSDApi.UpdateSocialServiceDatabase)    // 更新社会服务库
	}
	{
		SSDRouterWithoutRecord.GET("findSocialServiceDatabase", SSDApi.FindSocialServiceDatabase)        // 根据ID获取社会服务库
		SSDRouterWithoutRecord.GET("getSocialServiceDatabaseList", SSDApi.GetSocialServiceDatabaseList)  // 获取社会服务库列表
	}
}
