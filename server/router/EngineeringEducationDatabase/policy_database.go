package EngineeringEducationDatabase

import (
	"github.com/flipped-aurora/gin-vue-admin/server/api/v1"
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type PolicyDatabaseRouter struct {
}

// InitPolicyDatabaseRouter 初始化 政策数据库 路由信息
func (s *PolicyDatabaseRouter) InitPolicyDatabaseRouter(Router *gin.RouterGroup) {
	PDRouter := Router.Group("PD").Use(middleware.OperationRecord())
	PDRouterWithoutRecord := Router.Group("PD")
	var PDApi = v1.ApiGroupApp.EngineeringEducationDatabaseApiGroup.PolicyDatabaseApi
	{
		PDRouter.POST("createPolicyDatabase", PDApi.CreatePolicyDatabase)   // 新建政策数据库
		PDRouter.DELETE("deletePolicyDatabase", PDApi.DeletePolicyDatabase) // 删除政策数据库
		PDRouter.DELETE("deletePolicyDatabaseByIds", PDApi.DeletePolicyDatabaseByIds) // 批量删除政策数据库
		PDRouter.PUT("updatePolicyDatabase", PDApi.UpdatePolicyDatabase)    // 更新政策数据库
	}
	{
		PDRouterWithoutRecord.GET("findPolicyDatabase", PDApi.FindPolicyDatabase)        // 根据ID获取政策数据库
		PDRouterWithoutRecord.GET("getPolicyDatabaseList", PDApi.GetPolicyDatabaseList)  // 获取政策数据库列表
	}
}
