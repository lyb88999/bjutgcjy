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

type ExpertProfileApi struct{}

var expertProfileService = service.ServiceGroupApp.ExpertDatabaseServiceGroup.ExpertProfileService

// CreateExpertProfile 创建专家主档
// @Tags ExpertProfile
// @Summary 创建专家主档
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body ExpertDatabase.ExpertProfile true "创建专家主档"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"创建成功"}"
// @Router /expertProfile/createExpertProfile [post]
func (expertProfileApi *ExpertProfileApi) CreateExpertProfile(c *gin.Context) {
	var profile ExpertDatabase.ExpertProfile
	err := c.ShouldBindJSON(&profile)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	profile.CreatedBy = utils.GetUserID(c)
	profile.SubmittedBy = &profile.CreatedBy
	if err := expertProfileService.CreateExpertProfile(&profile); err != nil {
		global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败", c)
	} else {
		response.OkWithMessage("创建成功", c)
	}
}

// DeleteExpertProfile 删除专家主档
// @Tags ExpertProfile
// @Summary 删除专家主档
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /expertProfile/deleteExpertProfile [delete]
func (expertProfileApi *ExpertProfileApi) DeleteExpertProfile(c *gin.Context) {
	ID := c.Query("ID")
	userID := utils.GetUserID(c)
	if err := expertProfileService.DeleteExpertProfile(ID, userID); err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败", c)
	} else {
		response.OkWithMessage("删除成功", c)
	}
}

// DeleteExpertProfileByIds 批量删除专家主档
// @Tags ExpertProfile
// @Summary 批量删除专家主档
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"批量删除成功"}"
// @Router /expertProfile/deleteExpertProfileByIds [delete]
func (expertProfileApi *ExpertProfileApi) DeleteExpertProfileByIds(c *gin.Context) {
	IDs := c.QueryArray("IDs[]")
	userID := utils.GetUserID(c)
	if err := expertProfileService.DeleteExpertProfileByIds(IDs, userID); err != nil {
		global.GVA_LOG.Error("批量删除失败!", zap.Error(err))
		response.FailWithMessage("批量删除失败", c)
	} else {
		response.OkWithMessage("批量删除成功", c)
	}
}

// UpdateExpertProfile 更新专家主档
// @Tags ExpertProfile
// @Summary 更新专家主档
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body ExpertDatabase.ExpertProfile true "更新专家主档"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /expertProfile/updateExpertProfile [put]
func (expertProfileApi *ExpertProfileApi) UpdateExpertProfile(c *gin.Context) {
	var profile ExpertDatabase.ExpertProfile
	err := c.ShouldBindJSON(&profile)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	profile.UpdatedBy = utils.GetUserID(c)
	if err := expertProfileService.UpdateExpertProfile(profile); err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败", c)
	} else {
		response.OkWithMessage("更新成功", c)
	}
}

// FindExpertProfile 用id查询专家主档
// @Tags ExpertProfile
// @Summary 用id查询专家主档
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /expertProfile/findExpertProfile [get]
func (expertProfileApi *ExpertProfileApi) FindExpertProfile(c *gin.Context) {
	ID := c.Query("ID")
	if profile, err := expertProfileService.GetExpertProfile(ID); err != nil {
		global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败", c)
	} else {
		response.OkWithData(gin.H{"reExpertProfile": profile}, c)
	}
}

// GetExpertProfileList 分页获取专家主档列表
// @Tags ExpertProfile
// @Summary 分页获取专家主档列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query ExpertDatabaseReq.ExpertProfileSearch true "分页获取专家主档列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /expertProfile/getExpertProfileList [get]
func (expertProfileApi *ExpertProfileApi) GetExpertProfileList(c *gin.Context) {
	var pageInfo ExpertDatabaseReq.ExpertProfileSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if list, total, err := expertProfileService.GetExpertProfileInfoList(pageInfo); err != nil {
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
