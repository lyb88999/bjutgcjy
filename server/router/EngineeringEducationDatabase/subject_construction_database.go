package EngineeringEducationDatabase

import (
	"github.com/flipped-aurora/gin-vue-admin/server/api/v1"
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type SubjectConstructionDatabaseRouter struct {
}

// InitSubjectConstructionDatabaseRouter 初始化 学科建设库 路由信息
func (s *SubjectConstructionDatabaseRouter) InitSubjectConstructionDatabaseRouter(Router *gin.RouterGroup) {
	SCDRouter := Router.Group("SCD").Use(middleware.OperationRecord())
	SCDRouterWithoutRecord := Router.Group("SCD")
	var SCDApi = v1.ApiGroupApp.EngineeringEducationDatabaseApiGroup.SubjectConstructionDatabaseApi
	{
		SCDRouter.POST("createSubjectConstructionDatabase", SCDApi.CreateSubjectConstructionDatabase)   // 新建学科建设库
		SCDRouter.DELETE("deleteSubjectConstructionDatabase", SCDApi.DeleteSubjectConstructionDatabase) // 删除学科建设库
		SCDRouter.DELETE("deleteSubjectConstructionDatabaseByIds", SCDApi.DeleteSubjectConstructionDatabaseByIds) // 批量删除学科建设库
		SCDRouter.PUT("updateSubjectConstructionDatabase", SCDApi.UpdateSubjectConstructionDatabase)    // 更新学科建设库
	}
	{
		SCDRouterWithoutRecord.GET("findSubjectConstructionDatabase", SCDApi.FindSubjectConstructionDatabase)        // 根据ID获取学科建设库
		SCDRouterWithoutRecord.GET("getSubjectConstructionDatabaseList", SCDApi.GetSubjectConstructionDatabaseList)  // 获取学科建设库列表
	}
}
