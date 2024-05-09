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

type BasicInfomationDatabaseApi struct {
}

var BIDService = service.ServiceGroupApp.EngineeringEducationDatabaseServiceGroup.BasicInfomationDatabaseService


// CreateBasicInfomationDatabase 创建基础信息库
// @Tags BasicInfomationDatabase
// @Summary 创建基础信息库
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body EngineeringEducationDatabase.BasicInfomationDatabase true "创建基础信息库"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"创建成功"}"
// @Router /BID/createBasicInfomationDatabase [post]
func (BIDApi *BasicInfomationDatabaseApi) CreateBasicInfomationDatabase(c *gin.Context) {
	var BID EngineeringEducationDatabase.BasicInfomationDatabase
	err := c.ShouldBindJSON(&BID)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
    BID.CreatedBy = utils.GetUserID(c)

	if err := BIDService.CreateBasicInfomationDatabase(&BID); err != nil {
        global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败", c)
	} else {
		response.OkWithMessage("创建成功", c)
	}
}

// DeleteBasicInfomationDatabase 删除基础信息库
// @Tags BasicInfomationDatabase
// @Summary 删除基础信息库
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body EngineeringEducationDatabase.BasicInfomationDatabase true "删除基础信息库"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /BID/deleteBasicInfomationDatabase [delete]
func (BIDApi *BasicInfomationDatabaseApi) DeleteBasicInfomationDatabase(c *gin.Context) {
	ID := c.Query("ID")
    	userID := utils.GetUserID(c)
	if err := BIDService.DeleteBasicInfomationDatabase(ID,userID); err != nil {
        global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败", c)
	} else {
		response.OkWithMessage("删除成功", c)
	}
}

// DeleteBasicInfomationDatabaseByIds 批量删除基础信息库
// @Tags BasicInfomationDatabase
// @Summary 批量删除基础信息库
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"批量删除成功"}"
// @Router /BID/deleteBasicInfomationDatabaseByIds [delete]
func (BIDApi *BasicInfomationDatabaseApi) DeleteBasicInfomationDatabaseByIds(c *gin.Context) {
	IDs := c.QueryArray("IDs[]")
    userID := utils.GetUserID(c)
	if err := BIDService.DeleteBasicInfomationDatabaseByIds(IDs,userID); err != nil {
        global.GVA_LOG.Error("批量删除失败!", zap.Error(err))
		response.FailWithMessage("批量删除失败", c)
	} else {
		response.OkWithMessage("批量删除成功", c)
	}
}

// UpdateBasicInfomationDatabase 更新基础信息库
// @Tags BasicInfomationDatabase
// @Summary 更新基础信息库
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body EngineeringEducationDatabase.BasicInfomationDatabase true "更新基础信息库"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /BID/updateBasicInfomationDatabase [put]
func (BIDApi *BasicInfomationDatabaseApi) UpdateBasicInfomationDatabase(c *gin.Context) {
	var BID EngineeringEducationDatabase.BasicInfomationDatabase
	err := c.ShouldBindJSON(&BID)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
    BID.UpdatedBy = utils.GetUserID(c)

	if err := BIDService.UpdateBasicInfomationDatabase(BID); err != nil {
        global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败", c)
	} else {
		response.OkWithMessage("更新成功", c)
	}
}

// FindBasicInfomationDatabase 用id查询基础信息库
// @Tags BasicInfomationDatabase
// @Summary 用id查询基础信息库
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query EngineeringEducationDatabase.BasicInfomationDatabase true "用id查询基础信息库"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /BID/findBasicInfomationDatabase [get]
func (BIDApi *BasicInfomationDatabaseApi) FindBasicInfomationDatabase(c *gin.Context) {
	ID := c.Query("ID")
	if reBID, err := BIDService.GetBasicInfomationDatabase(ID); err != nil {
        global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败", c)
	} else {
		response.OkWithData(gin.H{"reBID": reBID}, c)
	}
}

// GetBasicInfomationDatabaseList 分页获取基础信息库列表
// @Tags BasicInfomationDatabase
// @Summary 分页获取基础信息库列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query EngineeringEducationDatabaseReq.BasicInfomationDatabaseSearch true "分页获取基础信息库列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /BID/getBasicInfomationDatabaseList [get]
func (BIDApi *BasicInfomationDatabaseApi) GetBasicInfomationDatabaseList(c *gin.Context) {
	var pageInfo EngineeringEducationDatabaseReq.BasicInfomationDatabaseSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if list, total, err := BIDService.GetBasicInfomationDatabaseInfoList(pageInfo); err != nil {
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
