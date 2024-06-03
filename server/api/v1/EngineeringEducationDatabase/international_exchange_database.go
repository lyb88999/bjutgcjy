package EngineeringEducationDatabase

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/EngineeringEducationDatabase"
	EngineeringEducationDatabaseReq "github.com/flipped-aurora/gin-vue-admin/server/model/EngineeringEducationDatabase/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/service"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type InternationalExchangeDatabaseApi struct {
}

var IEDService = service.ServiceGroupApp.EngineeringEducationDatabaseServiceGroup.InternationalExchangeDatabaseService

// CreateInternationalExchangeDatabase 创建国际交流库
// @Tags InternationalExchangeDatabase
// @Summary 创建国际交流库
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body EngineeringEducationDatabase.InternationalExchangeDatabase true "创建国际交流库"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"创建成功"}"
// @Router /IED/createInternationalExchangeDatabase [post]
func (IEDApi *InternationalExchangeDatabaseApi) CreateInternationalExchangeDatabase(c *gin.Context) {
	var IED EngineeringEducationDatabase.InternationalExchangeDatabase
	err := c.ShouldBindJSON(&IED)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	IED.CreatedBy = utils.GetUserID(c)

	if err := IEDService.CreateInternationalExchangeDatabase(&IED); err != nil {
		global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败", c)
	} else {
		response.OkWithMessage("创建成功", c)
	}
}

// DeleteInternationalExchangeDatabase 删除国际交流库
// @Tags InternationalExchangeDatabase
// @Summary 删除国际交流库
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body EngineeringEducationDatabase.InternationalExchangeDatabase true "删除国际交流库"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /IED/deleteInternationalExchangeDatabase [delete]
func (IEDApi *InternationalExchangeDatabaseApi) DeleteInternationalExchangeDatabase(c *gin.Context) {
	ID := c.Query("ID")
	userID := utils.GetUserID(c)
	if err := IEDService.DeleteInternationalExchangeDatabase(ID, userID); err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败", c)
	} else {
		response.OkWithMessage("删除成功", c)
	}
}

// DeleteInternationalExchangeDatabaseByIds 批量删除国际交流库
// @Tags InternationalExchangeDatabase
// @Summary 批量删除国际交流库
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"批量删除成功"}"
// @Router /IED/deleteInternationalExchangeDatabaseByIds [delete]
func (IEDApi *InternationalExchangeDatabaseApi) DeleteInternationalExchangeDatabaseByIds(c *gin.Context) {
	IDs := c.QueryArray("IDs[]")
	userID := utils.GetUserID(c)
	if err := IEDService.DeleteInternationalExchangeDatabaseByIds(IDs, userID); err != nil {
		global.GVA_LOG.Error("批量删除失败!", zap.Error(err))
		response.FailWithMessage("批量删除失败", c)
	} else {
		response.OkWithMessage("批量删除成功", c)
	}
}

// UpdateInternationalExchangeDatabase 更新国际交流库
// @Tags InternationalExchangeDatabase
// @Summary 更新国际交流库
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body EngineeringEducationDatabase.InternationalExchangeDatabase true "更新国际交流库"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /IED/updateInternationalExchangeDatabase [put]
func (IEDApi *InternationalExchangeDatabaseApi) UpdateInternationalExchangeDatabase(c *gin.Context) {
	var IED EngineeringEducationDatabase.InternationalExchangeDatabase
	err := c.ShouldBindJSON(&IED)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	IED.UpdatedBy = utils.GetUserID(c)

	if err := IEDService.UpdateInternationalExchangeDatabase(IED); err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败", c)
	} else {
		response.OkWithMessage("更新成功", c)
	}
}

// FindInternationalExchangeDatabase 用id查询国际交流库
// @Tags InternationalExchangeDatabase
// @Summary 用id查询国际交流库
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query EngineeringEducationDatabase.InternationalExchangeDatabase true "用id查询国际交流库"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /IED/findInternationalExchangeDatabase [get]
func (IEDApi *InternationalExchangeDatabaseApi) FindInternationalExchangeDatabase(c *gin.Context) {
	ID := c.Query("ID")
	if reIED, err := IEDService.GetInternationalExchangeDatabase(ID); err != nil {
		global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败", c)
	} else {
		response.OkWithData(gin.H{"reIED": reIED}, c)
	}
}

// GetInternationalExchangeDatabaseList 分页获取国际交流库列表
// @Tags InternationalExchangeDatabase
// @Summary 分页获取国际交流库列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query EngineeringEducationDatabaseReq.InternationalExchangeDatabaseSearch true "分页获取国际交流库列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /IED/getInternationalExchangeDatabaseList [get]
func (IEDApi *InternationalExchangeDatabaseApi) GetInternationalExchangeDatabaseList(c *gin.Context) {
	var pageInfo EngineeringEducationDatabaseReq.InternationalExchangeDatabaseSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if list, total, err := IEDService.GetInternationalExchangeDatabaseInfoList(pageInfo); err != nil {
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
