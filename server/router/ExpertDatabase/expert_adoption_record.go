package ExpertDatabase

import (
	v1 "github.com/flipped-aurora/gin-vue-admin/server/api/v1"
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type ExpertAdoptionRecordRouter struct{}

// InitExpertAdoptionRecordRouter 初始化 专家决策影响记录 路由信息
func (s *ExpertAdoptionRecordRouter) InitExpertAdoptionRecordRouter(Router *gin.RouterGroup) {
	expertAdoptionRecordRouter := Router.Group("expertAdoptionRecord").Use(middleware.OperationRecord())
	expertAdoptionRecordRouterWithoutRecord := Router.Group("expertAdoptionRecord")
	expertAdoptionRecordApi := v1.ApiGroupApp.ExpertDatabaseApiGroup.ExpertAdoptionRecordApi
	{
		expertAdoptionRecordRouter.POST("createExpertAdoptionRecord", expertAdoptionRecordApi.CreateExpertAdoptionRecord)             // 新建专家决策影响记录
		expertAdoptionRecordRouter.DELETE("deleteExpertAdoptionRecord", expertAdoptionRecordApi.DeleteExpertAdoptionRecord)           // 删除专家决策影响记录
		expertAdoptionRecordRouter.DELETE("deleteExpertAdoptionRecordByIds", expertAdoptionRecordApi.DeleteExpertAdoptionRecordByIds) // 批量删除专家决策影响记录
		expertAdoptionRecordRouter.PUT("updateExpertAdoptionRecord", expertAdoptionRecordApi.UpdateExpertAdoptionRecord)              // 更新专家决策影响记录
	}
	{
		expertAdoptionRecordRouterWithoutRecord.GET("findExpertAdoptionRecord", expertAdoptionRecordApi.FindExpertAdoptionRecord)       // 根据ID获取专家决策影响记录
		expertAdoptionRecordRouterWithoutRecord.GET("getExpertAdoptionRecordList", expertAdoptionRecordApi.GetExpertAdoptionRecordList) // 分页获取专家决策影响记录列表
	}
}
