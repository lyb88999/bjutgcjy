package ExpertDatabase

import (
	v1 "github.com/flipped-aurora/gin-vue-admin/server/api/v1"
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type ExpertTagRouter struct{}

// InitExpertTagRouter 初始化 专家标签 路由信息
func (s *ExpertTagRouter) InitExpertTagRouter(Router *gin.RouterGroup) {
	expertTagRouter := Router.Group("expertTag").Use(middleware.OperationRecord())
	expertTagRouterWithoutRecord := Router.Group("expertTag")
	expertTagApi := v1.ApiGroupApp.ExpertDatabaseApiGroup.ExpertTagApi
	{
		expertTagRouter.POST("createExpertTag", expertTagApi.CreateExpertTag)             // 新建标签
		expertTagRouter.DELETE("deleteExpertTag", expertTagApi.DeleteExpertTag)           // 删除标签
		expertTagRouter.PUT("updateExpertTag", expertTagApi.UpdateExpertTag)              // 更新标签
		expertTagRouter.POST("setExpertTagRelations", expertTagApi.SetExpertTagRelations) // 设置专家标签关联
	}
	{
		expertTagRouterWithoutRecord.GET("findExpertTag", expertTagApi.FindExpertTag)                     // 根据ID获取标签
		expertTagRouterWithoutRecord.GET("getExpertTagList", expertTagApi.GetExpertTagList)               // 分页获取标签列表
		expertTagRouterWithoutRecord.GET("getExpertTagsByExpertId", expertTagApi.GetExpertTagsByExpertId) // 获取某专家关联的标签
	}
}
