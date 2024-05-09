package EngineeringEducationDatabase

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
    "github.com/flipped-aurora/gin-vue-admin/server/model/EngineeringEducationDatabase"
    EngineeringEducationDatabaseReq "github.com/flipped-aurora/gin-vue-admin/server/model/EngineeringEducationDatabase/request"
    "github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
    "github.com/flipped-aurora/gin-vue-admin/server/service"
    "github.com/gin-gonic/gin"
    "go.uber.org/zap"
    "github.com/flipped-aurora/gin-vue-admin/server/utils"
)

type ConditionalGuaranteeDatabaseApi struct {
}

var CGDService = service.ServiceGroupApp.EngineeringEducationDatabaseServiceGroup.ConditionalGuaranteeDatabaseService


// CreateConditionalGuaranteeDatabase 创建条件保障库
// @Tags ConditionalGuaranteeDatabase
// @Summary 创建条件保障库
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body EngineeringEducationDatabase.ConditionalGuaranteeDatabase true "创建条件保障库"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"创建成功"}"
// @Router /CGD/createConditionalGuaranteeDatabase [post]
func (CGDApi *ConditionalGuaranteeDatabaseApi) CreateConditionalGuaranteeDatabase(c *gin.Context) {
	var CGD EngineeringEducationDatabase.ConditionalGuaranteeDatabase
	err := c.ShouldBindJSON(&CGD)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
    CGD.CreatedBy = utils.GetUserID(c)

	if err := CGDService.CreateConditionalGuaranteeDatabase(&CGD); err != nil {
        global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败", c)
	} else {
		response.OkWithMessage("创建成功", c)
	}
}

// DeleteConditionalGuaranteeDatabase 删除条件保障库
// @Tags ConditionalGuaranteeDatabase
// @Summary 删除条件保障库
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body EngineeringEducationDatabase.ConditionalGuaranteeDatabase true "删除条件保障库"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /CGD/deleteConditionalGuaranteeDatabase [delete]
func (CGDApi *ConditionalGuaranteeDatabaseApi) DeleteConditionalGuaranteeDatabase(c *gin.Context) {
	ID := c.Query("ID")
    	userID := utils.GetUserID(c)
	if err := CGDService.DeleteConditionalGuaranteeDatabase(ID,userID); err != nil {
        global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败", c)
	} else {
		response.OkWithMessage("删除成功", c)
	}
}

// DeleteConditionalGuaranteeDatabaseByIds 批量删除条件保障库
// @Tags ConditionalGuaranteeDatabase
// @Summary 批量删除条件保障库
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"批量删除成功"}"
// @Router /CGD/deleteConditionalGuaranteeDatabaseByIds [delete]
func (CGDApi *ConditionalGuaranteeDatabaseApi) DeleteConditionalGuaranteeDatabaseByIds(c *gin.Context) {
	IDs := c.QueryArray("IDs[]")
    userID := utils.GetUserID(c)
	if err := CGDService.DeleteConditionalGuaranteeDatabaseByIds(IDs,userID); err != nil {
        global.GVA_LOG.Error("批量删除失败!", zap.Error(err))
		response.FailWithMessage("批量删除失败", c)
	} else {
		response.OkWithMessage("批量删除成功", c)
	}
}

// UpdateConditionalGuaranteeDatabase 更新条件保障库
// @Tags ConditionalGuaranteeDatabase
// @Summary 更新条件保障库
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body EngineeringEducationDatabase.ConditionalGuaranteeDatabase true "更新条件保障库"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /CGD/updateConditionalGuaranteeDatabase [put]
func (CGDApi *ConditionalGuaranteeDatabaseApi) UpdateConditionalGuaranteeDatabase(c *gin.Context) {
	var CGD EngineeringEducationDatabase.ConditionalGuaranteeDatabase
	err := c.ShouldBindJSON(&CGD)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
    CGD.UpdatedBy = utils.GetUserID(c)

	if err := CGDService.UpdateConditionalGuaranteeDatabase(CGD); err != nil {
        global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败", c)
	} else {
		response.OkWithMessage("更新成功", c)
	}
}

// FindConditionalGuaranteeDatabase 用id查询条件保障库
// @Tags ConditionalGuaranteeDatabase
// @Summary 用id查询条件保障库
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query EngineeringEducationDatabase.ConditionalGuaranteeDatabase true "用id查询条件保障库"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /CGD/findConditionalGuaranteeDatabase [get]
func (CGDApi *ConditionalGuaranteeDatabaseApi) FindConditionalGuaranteeDatabase(c *gin.Context) {
	ID := c.Query("ID")
	if reCGD, err := CGDService.GetConditionalGuaranteeDatabase(ID); err != nil {
        global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败", c)
	} else {
		response.OkWithData(gin.H{"reCGD": reCGD}, c)
	}
}

// GetConditionalGuaranteeDatabaseList 分页获取条件保障库列表
// @Tags ConditionalGuaranteeDatabase
// @Summary 分页获取条件保障库列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query EngineeringEducationDatabaseReq.ConditionalGuaranteeDatabaseSearch true "分页获取条件保障库列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /CGD/getConditionalGuaranteeDatabaseList [get]
func (CGDApi *ConditionalGuaranteeDatabaseApi) GetConditionalGuaranteeDatabaseList(c *gin.Context) {
	var pageInfo EngineeringEducationDatabaseReq.ConditionalGuaranteeDatabaseSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if list, total, err := CGDService.GetConditionalGuaranteeDatabaseInfoList(pageInfo); err != nil {
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
