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

type ExpertAchievementApi struct{}

var expertAchievementService = service.ServiceGroupApp.ExpertDatabaseServiceGroup.ExpertAchievementService

// CreateExpertAchievement 创建专家成果
// @Tags ExpertAchievement
// @Summary 创建专家成果
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body ExpertDatabase.ExpertAchievement true "创建专家成果"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"创建成功"}"
// @Router /expertAchievement/createExpertAchievement [post]
func (expertAchievementApi *ExpertAchievementApi) CreateExpertAchievement(c *gin.Context) {
	var achievement ExpertDatabase.ExpertAchievement
	err := c.ShouldBindJSON(&achievement)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	achievement.CreatedBy = utils.GetUserID(c)
	if err := expertAchievementService.CreateExpertAchievement(&achievement); err != nil {
		global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败", c)
	} else {
		response.OkWithMessage("创建成功", c)
	}
}

// DeleteExpertAchievement 删除专家成果
// @Tags ExpertAchievement
// @Summary 删除专家成果
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /expertAchievement/deleteExpertAchievement [delete]
func (expertAchievementApi *ExpertAchievementApi) DeleteExpertAchievement(c *gin.Context) {
	ID := c.Query("ID")
	userID := utils.GetUserID(c)
	if err := expertAchievementService.DeleteExpertAchievement(ID, userID); err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败", c)
	} else {
		response.OkWithMessage("删除成功", c)
	}
}

// DeleteExpertAchievementByIds 批量删除专家成果
// @Tags ExpertAchievement
// @Summary 批量删除专家成果
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"批量删除成功"}"
// @Router /expertAchievement/deleteExpertAchievementByIds [delete]
func (expertAchievementApi *ExpertAchievementApi) DeleteExpertAchievementByIds(c *gin.Context) {
	IDs := c.QueryArray("IDs[]")
	userID := utils.GetUserID(c)
	if err := expertAchievementService.DeleteExpertAchievementByIds(IDs, userID); err != nil {
		global.GVA_LOG.Error("批量删除失败!", zap.Error(err))
		response.FailWithMessage("批量删除失败", c)
	} else {
		response.OkWithMessage("批量删除成功", c)
	}
}

// UpdateExpertAchievement 更新专家成果
// @Tags ExpertAchievement
// @Summary 更新专家成果
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body ExpertDatabase.ExpertAchievement true "更新专家成果"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /expertAchievement/updateExpertAchievement [put]
func (expertAchievementApi *ExpertAchievementApi) UpdateExpertAchievement(c *gin.Context) {
	var achievement ExpertDatabase.ExpertAchievement
	err := c.ShouldBindJSON(&achievement)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	achievement.UpdatedBy = utils.GetUserID(c)
	if err := expertAchievementService.UpdateExpertAchievement(achievement); err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败", c)
	} else {
		response.OkWithMessage("更新成功", c)
	}
}

// FindExpertAchievement 用id查询专家成果
// @Tags ExpertAchievement
// @Summary 用id查询专家成果
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /expertAchievement/findExpertAchievement [get]
func (expertAchievementApi *ExpertAchievementApi) FindExpertAchievement(c *gin.Context) {
	ID := c.Query("ID")
	if achievement, err := expertAchievementService.GetExpertAchievement(ID); err != nil {
		global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败", c)
	} else {
		response.OkWithData(gin.H{"reExpertAchievement": achievement}, c)
	}
}

// GetExpertAchievementList 分页获取专家成果列表
// @Tags ExpertAchievement
// @Summary 分页获取专家成果列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query ExpertDatabaseReq.ExpertAchievementSearch true "分页获取专家成果列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /expertAchievement/getExpertAchievementList [get]
func (expertAchievementApi *ExpertAchievementApi) GetExpertAchievementList(c *gin.Context) {
	var pageInfo ExpertDatabaseReq.ExpertAchievementSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if list, total, err := expertAchievementService.GetExpertAchievementInfoList(pageInfo); err != nil {
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
