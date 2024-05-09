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

type PolicyDatabaseApi struct {
}

var PDService = service.ServiceGroupApp.EngineeringEducationDatabaseServiceGroup.PolicyDatabaseService


// CreatePolicyDatabase 创建政策数据库
// @Tags PolicyDatabase
// @Summary 创建政策数据库
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body EngineeringEducationDatabase.PolicyDatabase true "创建政策数据库"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"创建成功"}"
// @Router /PD/createPolicyDatabase [post]
func (PDApi *PolicyDatabaseApi) CreatePolicyDatabase(c *gin.Context) {
	var PD EngineeringEducationDatabase.PolicyDatabase
	err := c.ShouldBindJSON(&PD)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
    PD.CreatedBy = utils.GetUserID(c)

	if err := PDService.CreatePolicyDatabase(&PD); err != nil {
        global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败", c)
	} else {
		response.OkWithMessage("创建成功", c)
	}
}

// DeletePolicyDatabase 删除政策数据库
// @Tags PolicyDatabase
// @Summary 删除政策数据库
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body EngineeringEducationDatabase.PolicyDatabase true "删除政策数据库"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /PD/deletePolicyDatabase [delete]
func (PDApi *PolicyDatabaseApi) DeletePolicyDatabase(c *gin.Context) {
	ID := c.Query("ID")
    	userID := utils.GetUserID(c)
	if err := PDService.DeletePolicyDatabase(ID,userID); err != nil {
        global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败", c)
	} else {
		response.OkWithMessage("删除成功", c)
	}
}

// DeletePolicyDatabaseByIds 批量删除政策数据库
// @Tags PolicyDatabase
// @Summary 批量删除政策数据库
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"批量删除成功"}"
// @Router /PD/deletePolicyDatabaseByIds [delete]
func (PDApi *PolicyDatabaseApi) DeletePolicyDatabaseByIds(c *gin.Context) {
	IDs := c.QueryArray("IDs[]")
    userID := utils.GetUserID(c)
	if err := PDService.DeletePolicyDatabaseByIds(IDs,userID); err != nil {
        global.GVA_LOG.Error("批量删除失败!", zap.Error(err))
		response.FailWithMessage("批量删除失败", c)
	} else {
		response.OkWithMessage("批量删除成功", c)
	}
}

// UpdatePolicyDatabase 更新政策数据库
// @Tags PolicyDatabase
// @Summary 更新政策数据库
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body EngineeringEducationDatabase.PolicyDatabase true "更新政策数据库"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /PD/updatePolicyDatabase [put]
func (PDApi *PolicyDatabaseApi) UpdatePolicyDatabase(c *gin.Context) {
	var PD EngineeringEducationDatabase.PolicyDatabase
	err := c.ShouldBindJSON(&PD)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
    PD.UpdatedBy = utils.GetUserID(c)

	if err := PDService.UpdatePolicyDatabase(PD); err != nil {
        global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败", c)
	} else {
		response.OkWithMessage("更新成功", c)
	}
}

// FindPolicyDatabase 用id查询政策数据库
// @Tags PolicyDatabase
// @Summary 用id查询政策数据库
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query EngineeringEducationDatabase.PolicyDatabase true "用id查询政策数据库"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /PD/findPolicyDatabase [get]
func (PDApi *PolicyDatabaseApi) FindPolicyDatabase(c *gin.Context) {
	ID := c.Query("ID")
	if rePD, err := PDService.GetPolicyDatabase(ID); err != nil {
        global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败", c)
	} else {
		response.OkWithData(gin.H{"rePD": rePD}, c)
	}
}

// GetPolicyDatabaseList 分页获取政策数据库列表
// @Tags PolicyDatabase
// @Summary 分页获取政策数据库列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query EngineeringEducationDatabaseReq.PolicyDatabaseSearch true "分页获取政策数据库列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /PD/getPolicyDatabaseList [get]
func (PDApi *PolicyDatabaseApi) GetPolicyDatabaseList(c *gin.Context) {
	var pageInfo EngineeringEducationDatabaseReq.PolicyDatabaseSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if list, total, err := PDService.GetPolicyDatabaseInfoList(pageInfo); err != nil {
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
