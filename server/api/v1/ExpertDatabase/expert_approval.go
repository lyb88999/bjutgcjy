package ExpertDatabase

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	ExpertDatabaseReq "github.com/flipped-aurora/gin-vue-admin/server/model/ExpertDatabase/request"
	commonRequest "github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/service"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ExpertApprovalApi struct{}

var expertApprovalService = service.ServiceGroupApp.ExpertDatabaseServiceGroup.ExpertApprovalService

// SubmitExpertProfile 提交审核（草稿/单位退回/市级退回 -> 待单位审核）
// @Tags ExpertApproval
// @Summary 提交专家档案审核
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body ExpertDatabaseReq.ExpertApprovalAction true "提交审核"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"提交成功"}"
// @Router /expertApproval/submit [post]
func (expertApprovalApi *ExpertApprovalApi) SubmitExpertProfile(c *gin.Context) {
	var req ExpertDatabaseReq.ExpertApprovalAction
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := expertApprovalService.Submit(req.ExpertId, utils.GetUserID(c)); err != nil {
		global.GVA_LOG.Error("提交失败!", zap.Error(err))
		response.FailWithMessage("提交失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("提交成功", c)
}

// OrgApproveExpertProfile 单位审核通过
// @Tags ExpertApproval
// @Summary 单位审核通过
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body ExpertDatabaseReq.ExpertApprovalAction true "单位审核通过"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"审核通过"}"
// @Router /expertApproval/orgApprove [post]
func (expertApprovalApi *ExpertApprovalApi) OrgApproveExpertProfile(c *gin.Context) {
	var req ExpertDatabaseReq.ExpertApprovalAction
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := expertApprovalService.OrgApprove(req.ExpertId, utils.GetUserID(c)); err != nil {
		global.GVA_LOG.Error("审核失败!", zap.Error(err))
		response.FailWithMessage("审核失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("审核通过", c)
}

// BatchOrgApproveExpertProfile 批量单位审核通过
// @Tags ExpertApproval
// @Summary 批量单位审核通过
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body ExpertDatabaseReq.ExpertBatchApprovalAction true "批量单位审核通过"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"批量审核完成"}"
// @Router /expertApproval/batchOrgApprove [post]
func (expertApprovalApi *ExpertApprovalApi) BatchOrgApproveExpertProfile(c *gin.Context) {
	var req ExpertDatabaseReq.ExpertBatchApprovalAction
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	result := expertApprovalService.BatchOrgApprove(req.ExpertIds, utils.GetUserID(c))
	response.OkWithDetailed(result, "批量审核完成", c)
}

// OrgRejectExpertProfile 单位审核退回
// @Tags ExpertApproval
// @Summary 单位审核退回
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body ExpertDatabaseReq.ExpertApprovalAction true "单位审核退回，需填写 opinion"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"已退回"}"
// @Router /expertApproval/orgReject [post]
func (expertApprovalApi *ExpertApprovalApi) OrgRejectExpertProfile(c *gin.Context) {
	var req ExpertDatabaseReq.ExpertApprovalAction
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := expertApprovalService.OrgReject(req.ExpertId, utils.GetUserID(c), req.Opinion); err != nil {
		global.GVA_LOG.Error("退回失败!", zap.Error(err))
		response.FailWithMessage("退回失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("已退回", c)
}

// CityApproveExpertProfile 市级审核通过
// @Tags ExpertApproval
// @Summary 市级审核通过
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body ExpertDatabaseReq.ExpertApprovalAction true "市级审核通过"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"审核通过"}"
// @Router /expertApproval/cityApprove [post]
func (expertApprovalApi *ExpertApprovalApi) CityApproveExpertProfile(c *gin.Context) {
	var req ExpertDatabaseReq.ExpertApprovalAction
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := expertApprovalService.CityApprove(req.ExpertId, utils.GetUserID(c)); err != nil {
		global.GVA_LOG.Error("审核失败!", zap.Error(err))
		response.FailWithMessage("审核失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("审核通过", c)
}

// BatchCityApproveExpertProfile 批量市级审核通过
// @Tags ExpertApproval
// @Summary 批量市级审核通过
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body ExpertDatabaseReq.ExpertBatchApprovalAction true "批量市级审核通过"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"批量审核完成"}"
// @Router /expertApproval/batchCityApprove [post]
func (expertApprovalApi *ExpertApprovalApi) BatchCityApproveExpertProfile(c *gin.Context) {
	var req ExpertDatabaseReq.ExpertBatchApprovalAction
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	result := expertApprovalService.BatchCityApprove(req.ExpertIds, utils.GetUserID(c))
	response.OkWithDetailed(result, "批量审核完成", c)
}

// CityRejectExpertProfile 市级审核退回
// @Tags ExpertApproval
// @Summary 市级审核退回
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body ExpertDatabaseReq.ExpertApprovalAction true "市级审核退回，需填写 opinion"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"已退回"}"
// @Router /expertApproval/cityReject [post]
func (expertApprovalApi *ExpertApprovalApi) CityRejectExpertProfile(c *gin.Context) {
	var req ExpertDatabaseReq.ExpertApprovalAction
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := expertApprovalService.CityReject(req.ExpertId, utils.GetUserID(c), req.Opinion); err != nil {
		global.GVA_LOG.Error("退回失败!", zap.Error(err))
		response.FailWithMessage("退回失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("已退回", c)
}

// AdminSetExpertProfileStatus 管理员直接改写审核状态（跳过流程，仍记录日志）
// @Tags ExpertApproval
// @Summary 管理员直接改写审核状态
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body ExpertDatabaseReq.ExpertAdminSetStatus true "管理员直接改写状态"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"设置成功"}"
// @Router /expertApproval/adminSetStatus [post]
func (expertApprovalApi *ExpertApprovalApi) AdminSetExpertProfileStatus(c *gin.Context) {
	var req ExpertDatabaseReq.ExpertAdminSetStatus
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := expertApprovalService.AdminSetStatus(req.ExpertId, req.Status, utils.GetUserID(c), req.Opinion); err != nil {
		global.GVA_LOG.Error("设置失败!", zap.Error(err))
		response.FailWithMessage("设置失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("设置成功", c)
}

// GetMyDrafts 我发起的专家档案列表
// @Tags ExpertApproval
// @Summary 我发起的专家档案列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query commonRequest.PageInfo true "分页获取我发起的专家档案列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /expertApproval/myDrafts [get]
func (expertApprovalApi *ExpertApprovalApi) GetMyDrafts(c *gin.Context) {
	var pageInfo commonRequest.PageInfo
	if err := c.ShouldBindQuery(&pageInfo); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := expertApprovalService.GetMyDrafts(utils.GetUserID(c), pageInfo)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithDetailed(response.PageResult{List: list, Total: total, Page: pageInfo.Page, PageSize: pageInfo.PageSize}, "获取成功", c)
}

// GetPendingOrgReview 待本单位审核的专家档案列表
// @Tags ExpertApproval
// @Summary 待本单位审核的专家档案列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query commonRequest.PageInfo true "分页获取待本单位审核的专家档案列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /expertApproval/pendingOrgReview [get]
func (expertApprovalApi *ExpertApprovalApi) GetPendingOrgReview(c *gin.Context) {
	var pageInfo commonRequest.PageInfo
	if err := c.ShouldBindQuery(&pageInfo); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := expertApprovalService.GetPendingOrgReview(utils.GetUserID(c), pageInfo)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败:"+err.Error(), c)
		return
	}
	response.OkWithDetailed(response.PageResult{List: list, Total: total, Page: pageInfo.Page, PageSize: pageInfo.PageSize}, "获取成功", c)
}

// GetPendingCityReview 待市级审核的专家档案列表
// @Tags ExpertApproval
// @Summary 待市级审核的专家档案列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query commonRequest.PageInfo true "分页获取待市级审核的专家档案列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /expertApproval/pendingCityReview [get]
func (expertApprovalApi *ExpertApprovalApi) GetPendingCityReview(c *gin.Context) {
	var pageInfo commonRequest.PageInfo
	if err := c.ShouldBindQuery(&pageInfo); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := expertApprovalService.GetPendingCityReview(pageInfo)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithDetailed(response.PageResult{List: list, Total: total, Page: pageInfo.Page, PageSize: pageInfo.PageSize}, "获取成功", c)
}

// GetExpertApprovalLogList 专家档案审核历史
// @Tags ExpertApproval
// @Summary 专家档案审核历史
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query ExpertDatabaseReq.ExpertApprovalLogSearch true "分页获取审核历史"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /expertApproval/getApprovalLogList [get]
func (expertApprovalApi *ExpertApprovalApi) GetExpertApprovalLogList(c *gin.Context) {
	var pageInfo ExpertDatabaseReq.ExpertApprovalLogSearch
	if err := c.ShouldBindQuery(&pageInfo); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := expertApprovalService.GetExpertApprovalLogList(pageInfo)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithDetailed(response.PageResult{List: list, Total: total, Page: pageInfo.Page, PageSize: pageInfo.PageSize}, "获取成功", c)
}
