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

type TalentTrainingDatabaseApi struct {
}

var TTDService = service.ServiceGroupApp.EngineeringEducationDatabaseServiceGroup.TalentTrainingDatabaseService


// CreateTalentTrainingDatabase 创建人才培养库
// @Tags TalentTrainingDatabase
// @Summary 创建人才培养库
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body EngineeringEducationDatabase.TalentTrainingDatabase true "创建人才培养库"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"创建成功"}"
// @Router /TTD/createTalentTrainingDatabase [post]
func (TTDApi *TalentTrainingDatabaseApi) CreateTalentTrainingDatabase(c *gin.Context) {
	var TTD EngineeringEducationDatabase.TalentTrainingDatabase
	err := c.ShouldBindJSON(&TTD)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
    TTD.CreatedBy = utils.GetUserID(c)

	if err := TTDService.CreateTalentTrainingDatabase(&TTD); err != nil {
        global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败", c)
	} else {
		response.OkWithMessage("创建成功", c)
	}
}

// DeleteTalentTrainingDatabase 删除人才培养库
// @Tags TalentTrainingDatabase
// @Summary 删除人才培养库
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body EngineeringEducationDatabase.TalentTrainingDatabase true "删除人才培养库"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /TTD/deleteTalentTrainingDatabase [delete]
func (TTDApi *TalentTrainingDatabaseApi) DeleteTalentTrainingDatabase(c *gin.Context) {
	ID := c.Query("ID")
    	userID := utils.GetUserID(c)
	if err := TTDService.DeleteTalentTrainingDatabase(ID,userID); err != nil {
        global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败", c)
	} else {
		response.OkWithMessage("删除成功", c)
	}
}

// DeleteTalentTrainingDatabaseByIds 批量删除人才培养库
// @Tags TalentTrainingDatabase
// @Summary 批量删除人才培养库
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"批量删除成功"}"
// @Router /TTD/deleteTalentTrainingDatabaseByIds [delete]
func (TTDApi *TalentTrainingDatabaseApi) DeleteTalentTrainingDatabaseByIds(c *gin.Context) {
	IDs := c.QueryArray("IDs[]")
    userID := utils.GetUserID(c)
	if err := TTDService.DeleteTalentTrainingDatabaseByIds(IDs,userID); err != nil {
        global.GVA_LOG.Error("批量删除失败!", zap.Error(err))
		response.FailWithMessage("批量删除失败", c)
	} else {
		response.OkWithMessage("批量删除成功", c)
	}
}

// UpdateTalentTrainingDatabase 更新人才培养库
// @Tags TalentTrainingDatabase
// @Summary 更新人才培养库
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body EngineeringEducationDatabase.TalentTrainingDatabase true "更新人才培养库"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /TTD/updateTalentTrainingDatabase [put]
func (TTDApi *TalentTrainingDatabaseApi) UpdateTalentTrainingDatabase(c *gin.Context) {
	var TTD EngineeringEducationDatabase.TalentTrainingDatabase
	err := c.ShouldBindJSON(&TTD)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
    TTD.UpdatedBy = utils.GetUserID(c)

	if err := TTDService.UpdateTalentTrainingDatabase(TTD); err != nil {
        global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败", c)
	} else {
		response.OkWithMessage("更新成功", c)
	}
}

// FindTalentTrainingDatabase 用id查询人才培养库
// @Tags TalentTrainingDatabase
// @Summary 用id查询人才培养库
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query EngineeringEducationDatabase.TalentTrainingDatabase true "用id查询人才培养库"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /TTD/findTalentTrainingDatabase [get]
func (TTDApi *TalentTrainingDatabaseApi) FindTalentTrainingDatabase(c *gin.Context) {
	ID := c.Query("ID")
	if reTTD, err := TTDService.GetTalentTrainingDatabase(ID); err != nil {
        global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败", c)
	} else {
		response.OkWithData(gin.H{"reTTD": reTTD}, c)
	}
}

// GetTalentTrainingDatabaseList 分页获取人才培养库列表
// @Tags TalentTrainingDatabase
// @Summary 分页获取人才培养库列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query EngineeringEducationDatabaseReq.TalentTrainingDatabaseSearch true "分页获取人才培养库列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /TTD/getTalentTrainingDatabaseList [get]
func (TTDApi *TalentTrainingDatabaseApi) GetTalentTrainingDatabaseList(c *gin.Context) {
	var pageInfo EngineeringEducationDatabaseReq.TalentTrainingDatabaseSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if list, total, err := TTDService.GetTalentTrainingDatabaseInfoList(pageInfo); err != nil {
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
