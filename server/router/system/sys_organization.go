package system

import (
	v1 "github.com/flipped-aurora/gin-vue-admin/server/api/v1"
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type OrganizationRouter struct{}

func (s *OrganizationRouter) InitSysOrganizationRouter(Router *gin.RouterGroup) {
	sysOrganizationRouter := Router.Group("sysOrganization").Use(middleware.OperationRecord())
	sysOrganizationRouterWithoutRecord := Router.Group("sysOrganization")
	sysOrganizationApi := v1.ApiGroupApp.SystemApiGroup.OrganizationApi
	{
		sysOrganizationRouter.POST("createSysOrganization", sysOrganizationApi.CreateSysOrganization)   // 新建单位
		sysOrganizationRouter.DELETE("deleteSysOrganization", sysOrganizationApi.DeleteSysOrganization) // 删除单位
		sysOrganizationRouter.PUT("updateSysOrganization", sysOrganizationApi.UpdateSysOrganization)    // 更新单位
	}
	{
		sysOrganizationRouterWithoutRecord.GET("findSysOrganization", sysOrganizationApi.FindSysOrganization)       // 根据ID获取单位
		sysOrganizationRouterWithoutRecord.GET("getSysOrganizationList", sysOrganizationApi.GetSysOrganizationList) // 分页获取单位列表
		sysOrganizationRouterWithoutRecord.GET("getSysOrganizationTree", sysOrganizationApi.GetSysOrganizationTree) // 获取单位树
	}
}
