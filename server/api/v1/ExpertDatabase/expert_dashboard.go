package ExpertDatabase

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/service"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ExpertDashboardApi struct{}

var expertDashboardService = service.ServiceGroupApp.ExpertDatabaseServiceGroup.ExpertDashboardService

// GetDashboardStats 专家库统计概览
// @Tags ExpertDashboard
// @Summary 专家库统计概览
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /expertDatabase/dashboardStats [get]
func (expertDashboardApi *ExpertDashboardApi) GetDashboardStats(c *gin.Context) {
	stats, err := expertDashboardService.GetDashboardStats(utils.GetUserID(c))
	if err != nil {
		global.GVA_LOG.Error("获取统计概览失败!", zap.Error(err))
		response.FailWithMessage("获取统计概览失败", c)
		return
	}
	response.OkWithData(stats, c)
}
