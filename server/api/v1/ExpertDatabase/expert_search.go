package ExpertDatabase

import (
	"net/http"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	ExpertDatabaseReq "github.com/flipped-aurora/gin-vue-admin/server/model/ExpertDatabase/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ExpertSearchApi struct{}

var expertSearchService = service.ServiceGroupApp.ExpertDatabaseServiceGroup.ExpertSearchService
var expertScoreService = service.ServiceGroupApp.ExpertDatabaseServiceGroup.ExpertScoreService

// SearchExpert 专家综合推荐排序检索
// @Tags ExpertSearch
// @Summary 专家综合推荐排序检索
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query ExpertDatabaseReq.ExpertSearchReq true "检索关键词与筛选条件"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"检索成功"}"
// @Router /expertDatabase/search [get]
func (expertSearchApi *ExpertSearchApi) SearchExpert(c *gin.Context) {
	var req ExpertDatabaseReq.ExpertSearchReq
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := expertSearchService.Search(req)
	if err != nil {
		global.GVA_LOG.Error("检索失败!", zap.Error(err))
		response.FailWithMessage("检索失败", c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, "检索成功", c)
}

// ExportSearchResults 导出当前检索条件下命中的全部结果
// @Tags ExpertSearch
// @Summary 导出检索结果
// @Security ApiKeyAuth
// @Produce application/octet-stream
// @Param data query ExpertDatabaseReq.ExpertSearchReq true "检索关键词与筛选条件"
// @Router /expertDatabase/exportSearchResults [get]
func (expertSearchApi *ExpertSearchApi) ExportSearchResults(c *gin.Context) {
	var req ExpertDatabaseReq.ExpertSearchReq
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	f, err := expertSearchService.ExportSearchResults(req)
	if err != nil {
		global.GVA_LOG.Error("导出失败!", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	buf, err := f.WriteToBuffer()
	if err != nil {
		global.GVA_LOG.Error("导出失败!", zap.Error(err))
		response.FailWithMessage("导出失败", c)
		return
	}
	c.Header("Content-Disposition", "attachment; filename=expert_search_export.xlsx")
	c.Header("success", "true")
	c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", buf.Bytes())
}

// RecomputeExpertScore 手动触发单个专家的得分重算（管理员/数据变更排查用）
// @Tags ExpertSearch
// @Summary 手动重算专家得分
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param ID query int true "专家ID"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"重算成功"}"
// @Router /expertDatabase/recomputeScore [post]
func (expertSearchApi *ExpertSearchApi) RecomputeExpertScore(c *gin.Context) {
	var req struct {
		ExpertId uint `json:"expertId" form:"expertId" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := expertScoreService.RecomputeExpertScore(req.ExpertId); err != nil {
		global.GVA_LOG.Error("重算失败!", zap.Error(err))
		response.FailWithMessage("重算失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("重算成功", c)
}
