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

type ExpertAdoptionRecordApi struct{}

var expertAdoptionRecordService = service.ServiceGroupApp.ExpertDatabaseServiceGroup.ExpertAdoptionRecordService

// CreateExpertAdoptionRecord 创建专家决策影响记录
// @Tags ExpertAdoptionRecord
// @Summary 创建专家决策影响记录
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body ExpertDatabase.ExpertAdoptionRecord true "创建专家决策影响记录"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"创建成功"}"
// @Router /expertAdoptionRecord/createExpertAdoptionRecord [post]
func (expertAdoptionRecordApi *ExpertAdoptionRecordApi) CreateExpertAdoptionRecord(c *gin.Context) {
	var record ExpertDatabase.ExpertAdoptionRecord
	err := c.ShouldBindJSON(&record)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	record.CreatedBy = utils.GetUserID(c)
	if err := expertAdoptionRecordService.CreateExpertAdoptionRecord(&record); err != nil {
		global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败", c)
	} else {
		response.OkWithMessage("创建成功", c)
	}
}

// DeleteExpertAdoptionRecord 删除专家决策影响记录
// @Tags ExpertAdoptionRecord
// @Summary 删除专家决策影响记录
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /expertAdoptionRecord/deleteExpertAdoptionRecord [delete]
func (expertAdoptionRecordApi *ExpertAdoptionRecordApi) DeleteExpertAdoptionRecord(c *gin.Context) {
	ID := c.Query("ID")
	userID := utils.GetUserID(c)
	if err := expertAdoptionRecordService.DeleteExpertAdoptionRecord(ID, userID); err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败", c)
	} else {
		response.OkWithMessage("删除成功", c)
	}
}

// DeleteExpertAdoptionRecordByIds 批量删除专家决策影响记录
// @Tags ExpertAdoptionRecord
// @Summary 批量删除专家决策影响记录
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"批量删除成功"}"
// @Router /expertAdoptionRecord/deleteExpertAdoptionRecordByIds [delete]
func (expertAdoptionRecordApi *ExpertAdoptionRecordApi) DeleteExpertAdoptionRecordByIds(c *gin.Context) {
	IDs := c.QueryArray("IDs[]")
	userID := utils.GetUserID(c)
	if err := expertAdoptionRecordService.DeleteExpertAdoptionRecordByIds(IDs, userID); err != nil {
		global.GVA_LOG.Error("批量删除失败!", zap.Error(err))
		response.FailWithMessage("批量删除失败", c)
	} else {
		response.OkWithMessage("批量删除成功", c)
	}
}

// UpdateExpertAdoptionRecord 更新专家决策影响记录
// @Tags ExpertAdoptionRecord
// @Summary 更新专家决策影响记录
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body ExpertDatabase.ExpertAdoptionRecord true "更新专家决策影响记录"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /expertAdoptionRecord/updateExpertAdoptionRecord [put]
func (expertAdoptionRecordApi *ExpertAdoptionRecordApi) UpdateExpertAdoptionRecord(c *gin.Context) {
	var record ExpertDatabase.ExpertAdoptionRecord
	err := c.ShouldBindJSON(&record)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	record.UpdatedBy = utils.GetUserID(c)
	if err := expertAdoptionRecordService.UpdateExpertAdoptionRecord(record); err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败", c)
	} else {
		response.OkWithMessage("更新成功", c)
	}
}

// FindExpertAdoptionRecord 用id查询专家决策影响记录
// @Tags ExpertAdoptionRecord
// @Summary 用id查询专家决策影响记录
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /expertAdoptionRecord/findExpertAdoptionRecord [get]
func (expertAdoptionRecordApi *ExpertAdoptionRecordApi) FindExpertAdoptionRecord(c *gin.Context) {
	ID := c.Query("ID")
	if record, err := expertAdoptionRecordService.GetExpertAdoptionRecord(ID); err != nil {
		global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败", c)
	} else {
		response.OkWithData(gin.H{"reExpertAdoptionRecord": record}, c)
	}
}

// GetExpertAdoptionRecordList 分页获取专家决策影响记录列表
// @Tags ExpertAdoptionRecord
// @Summary 分页获取专家决策影响记录列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query ExpertDatabaseReq.ExpertAdoptionRecordSearch true "分页获取专家决策影响记录列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /expertAdoptionRecord/getExpertAdoptionRecordList [get]
func (expertAdoptionRecordApi *ExpertAdoptionRecordApi) GetExpertAdoptionRecordList(c *gin.Context) {
	var pageInfo ExpertDatabaseReq.ExpertAdoptionRecordSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if list, total, err := expertAdoptionRecordService.GetExpertAdoptionRecordInfoList(pageInfo); err != nil {
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
