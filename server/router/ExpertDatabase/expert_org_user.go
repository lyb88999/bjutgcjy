package ExpertDatabase

import (
	v1 "github.com/flipped-aurora/gin-vue-admin/server/api/v1"
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type ExpertOrgUserRouter struct{}

// InitExpertOrgUserRouter 初始化 本单位账号管理 路由信息（单位审核员为本单位新建/管理个人申报人账号）
func (s *ExpertOrgUserRouter) InitExpertOrgUserRouter(Router *gin.RouterGroup) {
	expertOrgUserRouter := Router.Group("expertOrgUser").Use(middleware.OperationRecord())
	expertOrgUserRouterWithoutRecord := Router.Group("expertOrgUser")
	expertOrgUserApi := v1.ApiGroupApp.ExpertDatabaseApiGroup.ExpertOrgUserApi
	{
		expertOrgUserRouter.POST("createOrgApplicant", expertOrgUserApi.CreateOrgApplicant)   // 新建本单位个人申报人账号
		expertOrgUserRouter.POST("toggleOrgUserEnable", expertOrgUserApi.ToggleOrgUserEnable) // 启用/冻结本单位账号
	}
	{
		expertOrgUserRouterWithoutRecord.GET("getOrgUserList", expertOrgUserApi.GetOrgUserList) // 分页获取本单位账号列表
	}
}
