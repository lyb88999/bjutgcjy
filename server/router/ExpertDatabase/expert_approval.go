package ExpertDatabase

import (
	v1 "github.com/flipped-aurora/gin-vue-admin/server/api/v1"
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type ExpertApprovalRouter struct{}

// InitExpertApprovalRouter 初始化 专家档案审核流转 路由信息
func (s *ExpertApprovalRouter) InitExpertApprovalRouter(Router *gin.RouterGroup) {
	expertApprovalRouter := Router.Group("expertApproval").Use(middleware.OperationRecord())
	expertApprovalRouterWithoutRecord := Router.Group("expertApproval")
	expertApprovalApi := v1.ApiGroupApp.ExpertDatabaseApiGroup.ExpertApprovalApi
	{
		expertApprovalRouter.POST("submit", expertApprovalApi.SubmitExpertProfile)                     // 提交审核
		expertApprovalRouter.POST("orgApprove", expertApprovalApi.OrgApproveExpertProfile)             // 单位审核通过
		expertApprovalRouter.POST("batchOrgApprove", expertApprovalApi.BatchOrgApproveExpertProfile)   // 批量单位审核通过
		expertApprovalRouter.POST("orgReject", expertApprovalApi.OrgRejectExpertProfile)               // 单位审核退回
		expertApprovalRouter.POST("cityApprove", expertApprovalApi.CityApproveExpertProfile)           // 市级审核通过
		expertApprovalRouter.POST("batchCityApprove", expertApprovalApi.BatchCityApproveExpertProfile) // 批量市级审核通过
		expertApprovalRouter.POST("cityReject", expertApprovalApi.CityRejectExpertProfile)             // 市级审核退回
		expertApprovalRouter.POST("adminSetStatus", expertApprovalApi.AdminSetExpertProfileStatus)     // 管理员直接改写状态
	}
	{
		expertApprovalRouterWithoutRecord.GET("myDrafts", expertApprovalApi.GetMyDrafts)                        // 我发起的
		expertApprovalRouterWithoutRecord.GET("pendingOrgReview", expertApprovalApi.GetPendingOrgReview)        // 待本单位审核
		expertApprovalRouterWithoutRecord.GET("pendingCityReview", expertApprovalApi.GetPendingCityReview)      // 待市级审核
		expertApprovalRouterWithoutRecord.GET("getApprovalLogList", expertApprovalApi.GetExpertApprovalLogList) // 审核历史
	}
}
