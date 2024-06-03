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

type ProfessionalConstructionDatabaseApi struct {
}

var PCDService = service.ServiceGroupApp.EngineeringEducationDatabaseServiceGroup.ProfessionalConstructionDatabaseService

// CreateProfessionalConstructionDatabase 创建专业建设库
// @Tags ProfessionalConstructionDatabase
// @Summary 创建专业建设库
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body EngineeringEducationDatabase.ProfessionalConstructionDatabase true "创建专业建设库"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"创建成功"}"
// @Router /PCD/createProfessionalConstructionDatabase [post]
func (PCDApi *ProfessionalConstructionDatabaseApi) CreateProfessionalConstructionDatabase(c *gin.Context) {
	var PCD EngineeringEducationDatabase.ProfessionalConstructionDatabase
	err := c.ShouldBindJSON(&PCD)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	PCD.CreatedBy = utils.GetUserID(c)

	if err := PCDService.CreateProfessionalConstructionDatabase(&PCD); err != nil {
		global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败", c)
	} else {
		response.OkWithMessage("创建成功", c)
	}
}

// DeleteProfessionalConstructionDatabase 删除专业建设库
// @Tags ProfessionalConstructionDatabase
// @Summary 删除专业建设库
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body EngineeringEducationDatabase.ProfessionalConstructionDatabase true "删除专业建设库"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /PCD/deleteProfessionalConstructionDatabase [delete]
func (PCDApi *ProfessionalConstructionDatabaseApi) DeleteProfessionalConstructionDatabase(c *gin.Context) {
	ID := c.Query("ID")
	userID := utils.GetUserID(c)
	if err := PCDService.DeleteProfessionalConstructionDatabase(ID, userID); err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败", c)
	} else {
		response.OkWithMessage("删除成功", c)
	}
}

// DeleteProfessionalConstructionDatabaseByIds 批量删除专业建设库
// @Tags ProfessionalConstructionDatabase
// @Summary 批量删除专业建设库
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"批量删除成功"}"
// @Router /PCD/deleteProfessionalConstructionDatabaseByIds [delete]
func (PCDApi *ProfessionalConstructionDatabaseApi) DeleteProfessionalConstructionDatabaseByIds(c *gin.Context) {
	IDs := c.QueryArray("IDs[]")
	userID := utils.GetUserID(c)
	if err := PCDService.DeleteProfessionalConstructionDatabaseByIds(IDs, userID); err != nil {
		global.GVA_LOG.Error("批量删除失败!", zap.Error(err))
		response.FailWithMessage("批量删除失败", c)
	} else {
		response.OkWithMessage("批量删除成功", c)
	}
}

// UpdateProfessionalConstructionDatabase 更新专业建设库
// @Tags ProfessionalConstructionDatabase
// @Summary 更新专业建设库
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body EngineeringEducationDatabase.ProfessionalConstructionDatabase true "更新专业建设库"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /PCD/updateProfessionalConstructionDatabase [put]
func (PCDApi *ProfessionalConstructionDatabaseApi) UpdateProfessionalConstructionDatabase(c *gin.Context) {
	var PCD EngineeringEducationDatabase.ProfessionalConstructionDatabase
	err := c.ShouldBindJSON(&PCD)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	PCD.UpdatedBy = utils.GetUserID(c)

	if err := PCDService.UpdateProfessionalConstructionDatabase(PCD); err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败", c)
	} else {
		response.OkWithMessage("更新成功", c)
	}
}

// FindProfessionalConstructionDatabase 用id查询专业建设库
// @Tags ProfessionalConstructionDatabase
// @Summary 用id查询专业建设库
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query EngineeringEducationDatabase.ProfessionalConstructionDatabase true "用id查询专业建设库"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /PCD/findProfessionalConstructionDatabase [get]
func (PCDApi *ProfessionalConstructionDatabaseApi) FindProfessionalConstructionDatabase(c *gin.Context) {
	ID := c.Query("ID")
	if rePCD, err := PCDService.GetProfessionalConstructionDatabase(ID); err != nil {
		global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败", c)
	} else {
		response.OkWithData(gin.H{"rePCD": rePCD}, c)
	}
}

// GetProfessionalConstructionDatabaseList 分页获取专业建设库列表
// @Tags ProfessionalConstructionDatabase
// @Summary 分页获取专业建设库列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query EngineeringEducationDatabaseReq.ProfessionalConstructionDatabaseSearch true "分页获取专业建设库列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /PCD/getProfessionalConstructionDatabaseList [get]
func (PCDApi *ProfessionalConstructionDatabaseApi) GetProfessionalConstructionDatabaseList(c *gin.Context) {
	var pageInfo EngineeringEducationDatabaseReq.ProfessionalConstructionDatabaseSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if list, total, err := PCDService.GetProfessionalConstructionDatabaseInfoList(pageInfo); err != nil {
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
