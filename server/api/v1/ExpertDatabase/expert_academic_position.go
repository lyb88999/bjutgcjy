package ExpertDatabase

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/ExpertDatabase"
	ExpertDatabaseReq "github.com/flipped-aurora/gin-vue-admin/server/model/ExpertDatabase/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/service"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ExpertAcademicPositionApi struct{}

var expertAcademicPositionService = service.ServiceGroupApp.ExpertDatabaseServiceGroup.ExpertAcademicPositionService

// CreateExpertAcademicPosition 创建专家学术兼职
// @Tags ExpertAcademicPosition
// @Summary 创建专家学术兼职
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body ExpertDatabase.ExpertAcademicPosition true "创建专家学术兼职"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"创建成功"}"
// @Router /expertAcademicPosition/createExpertAcademicPosition [post]
func (expertAcademicPositionApi *ExpertAcademicPositionApi) CreateExpertAcademicPosition(c *gin.Context) {
	var position ExpertDatabase.ExpertAcademicPosition
	err := c.ShouldBindJSON(&position)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	position.CreatedBy = utils.GetUserID(c)
	if err := expertAcademicPositionService.CreateExpertAcademicPosition(&position); err != nil {
		global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败", c)
	} else {
		response.OkWithMessage("创建成功", c)
	}
}

// DeleteExpertAcademicPosition 删除专家学术兼职
// @Tags ExpertAcademicPosition
// @Summary 删除专家学术兼职
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /expertAcademicPosition/deleteExpertAcademicPosition [delete]
func (expertAcademicPositionApi *ExpertAcademicPositionApi) DeleteExpertAcademicPosition(c *gin.Context) {
	ID := c.Query("ID")
	userID := utils.GetUserID(c)
	if err := expertAcademicPositionService.DeleteExpertAcademicPosition(ID, userID); err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败", c)
	} else {
		response.OkWithMessage("删除成功", c)
	}
}

// DeleteExpertAcademicPositionByIds 批量删除专家学术兼职
// @Tags ExpertAcademicPosition
// @Summary 批量删除专家学术兼职
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"批量删除成功"}"
// @Router /expertAcademicPosition/deleteExpertAcademicPositionByIds [delete]
func (expertAcademicPositionApi *ExpertAcademicPositionApi) DeleteExpertAcademicPositionByIds(c *gin.Context) {
	IDs := c.QueryArray("IDs[]")
	userID := utils.GetUserID(c)
	if err := expertAcademicPositionService.DeleteExpertAcademicPositionByIds(IDs, userID); err != nil {
		global.GVA_LOG.Error("批量删除失败!", zap.Error(err))
		response.FailWithMessage("批量删除失败", c)
	} else {
		response.OkWithMessage("批量删除成功", c)
	}
}

// UpdateExpertAcademicPosition 更新专家学术兼职
// @Tags ExpertAcademicPosition
// @Summary 更新专家学术兼职
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body ExpertDatabase.ExpertAcademicPosition true "更新专家学术兼职"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /expertAcademicPosition/updateExpertAcademicPosition [put]
func (expertAcademicPositionApi *ExpertAcademicPositionApi) UpdateExpertAcademicPosition(c *gin.Context) {
	var position ExpertDatabase.ExpertAcademicPosition
	err := c.ShouldBindJSON(&position)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	position.UpdatedBy = utils.GetUserID(c)
	if err := expertAcademicPositionService.UpdateExpertAcademicPosition(position); err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败", c)
	} else {
		response.OkWithMessage("更新成功", c)
	}
}

// FindExpertAcademicPosition 用id查询专家学术兼职
// @Tags ExpertAcademicPosition
// @Summary 用id查询专家学术兼职
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /expertAcademicPosition/findExpertAcademicPosition [get]
func (expertAcademicPositionApi *ExpertAcademicPositionApi) FindExpertAcademicPosition(c *gin.Context) {
	ID := c.Query("ID")
	if position, err := expertAcademicPositionService.GetExpertAcademicPosition(ID); err != nil {
		global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败", c)
	} else {
		response.OkWithData(gin.H{"reExpertAcademicPosition": position}, c)
	}
}

// GetExpertAcademicPositionList 分页获取专家学术兼职列表
// @Tags ExpertAcademicPosition
// @Summary 分页获取专家学术兼职列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query ExpertDatabaseReq.ExpertAcademicPositionSearch true "分页获取专家学术兼职列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /expertAcademicPosition/getExpertAcademicPositionList [get]
func (expertAcademicPositionApi *ExpertAcademicPositionApi) GetExpertAcademicPositionList(c *gin.Context) {
	var pageInfo ExpertDatabaseReq.ExpertAcademicPositionSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if list, total, err := expertAcademicPositionService.GetExpertAcademicPositionInfoList(pageInfo); err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
	} else {
		response.OkWithDetailed(response.PageResult{
			List:     list,
			Total:    total,
			Page:     pageInfo.Page,
			PageSize: pageInfo.PageSize,
		}, "获取成功", c)
	}
}
