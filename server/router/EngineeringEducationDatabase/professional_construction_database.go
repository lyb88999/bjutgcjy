package EngineeringEducationDatabase

import (
	"github.com/flipped-aurora/gin-vue-admin/server/api/v1"
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type ProfessionalConstructionDatabaseRouter struct {
}

// InitProfessionalConstructionDatabaseRouter 初始化 专业建设库 路由信息
func (s *ProfessionalConstructionDatabaseRouter) InitProfessionalConstructionDatabaseRouter(Router *gin.RouterGroup) {
	PCDRouter := Router.Group("PCD").Use(middleware.OperationRecord())
	PCDRouterWithoutRecord := Router.Group("PCD")
	var PCDApi = v1.ApiGroupApp.EngineeringEducationDatabaseApiGroup.ProfessionalConstructionDatabaseApi
	{
		PCDRouter.POST("createProfessionalConstructionDatabase", PCDApi.CreateProfessionalConstructionDatabase)   // 新建专业建设库
		PCDRouter.DELETE("deleteProfessionalConstructionDatabase", PCDApi.DeleteProfessionalConstructionDatabase) // 删除专业建设库
		PCDRouter.DELETE("deleteProfessionalConstructionDatabaseByIds", PCDApi.DeleteProfessionalConstructionDatabaseByIds) // 批量删除专业建设库
		PCDRouter.PUT("updateProfessionalConstructionDatabase", PCDApi.UpdateProfessionalConstructionDatabase)    // 更新专业建设库
	}
	{
		PCDRouterWithoutRecord.GET("findProfessionalConstructionDatabase", PCDApi.FindProfessionalConstructionDatabase)        // 根据ID获取专业建设库
		PCDRouterWithoutRecord.GET("getProfessionalConstructionDatabaseList", PCDApi.GetProfessionalConstructionDatabaseList)  // 获取专业建设库列表
	}
}
