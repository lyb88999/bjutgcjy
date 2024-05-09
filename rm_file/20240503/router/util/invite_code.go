package util

import (
	"github.com/flipped-aurora/gin-vue-admin/server/api/v1"
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type InviteCodeRouter struct {
}

// InitInviteCodeRouter 初始化 邀请码 路由信息
func (s *InviteCodeRouter) InitInviteCodeRouter(Router *gin.RouterGroup) {
	inviteCodeRouter := Router.Group("inviteCode").Use(middleware.OperationRecord())
	inviteCodeRouterWithoutRecord := Router.Group("inviteCode")
	var inviteCodeApi = v1.ApiGroupApp.UtilApiGroup.InviteCodeApi
	{
		inviteCodeRouter.POST("createInviteCode", inviteCodeApi.CreateInviteCode)   // 新建邀请码
		inviteCodeRouter.DELETE("deleteInviteCode", inviteCodeApi.DeleteInviteCode) // 删除邀请码
		inviteCodeRouter.DELETE("deleteInviteCodeByIds", inviteCodeApi.DeleteInviteCodeByIds) // 批量删除邀请码
		inviteCodeRouter.PUT("updateInviteCode", inviteCodeApi.UpdateInviteCode)    // 更新邀请码
	}
	{
		inviteCodeRouterWithoutRecord.GET("findInviteCode", inviteCodeApi.FindInviteCode)        // 根据ID获取邀请码
		inviteCodeRouterWithoutRecord.GET("getInviteCodeList", inviteCodeApi.GetInviteCodeList)  // 获取邀请码列表
	}
}
