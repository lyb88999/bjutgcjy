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

type SocialServiceDatabaseApi struct {
}

var SSDService = service.ServiceGroupApp.EngineeringEducationDatabaseServiceGroup.SocialServiceDatabaseService


// CreateSocialServiceDatabase 创建社会服务库
// @Tags SocialServiceDatabase
// @Summary 创建社会服务库
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body EngineeringEducationDatabase.SocialServiceDatabase true "创建社会服务库"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"创建成功"}"
// @Router /SSD/createSocialServiceDatabase [post]
func (SSDApi *SocialServiceDatabaseApi) CreateSocialServiceDatabase(c *gin.Context) {
	var SSD EngineeringEducationDatabase.SocialServiceDatabase
	err := c.ShouldBindJSON(&SSD)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
    SSD.CreatedBy = utils.GetUserID(c)

	if err := SSDService.CreateSocialServiceDatabase(&SSD); err != nil {
        global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败", c)
	} else {
		response.OkWithMessage("创建成功", c)
	}
}

// DeleteSocialServiceDatabase 删除社会服务库
// @Tags SocialServiceDatabase
// @Summary 删除社会服务库
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body EngineeringEducationDatabase.SocialServiceDatabase true "删除社会服务库"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /SSD/deleteSocialServiceDatabase [delete]
func (SSDApi *SocialServiceDatabaseApi) DeleteSocialServiceDatabase(c *gin.Context) {
	ID := c.Query("ID")
    	userID := utils.GetUserID(c)
	if err := SSDService.DeleteSocialServiceDatabase(ID,userID); err != nil {
        global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败", c)
	} else {
		response.OkWithMessage("删除成功", c)
	}
}

// DeleteSocialServiceDatabaseByIds 批量删除社会服务库
// @Tags SocialServiceDatabase
// @Summary 批量删除社会服务库
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"批量删除成功"}"
// @Router /SSD/deleteSocialServiceDatabaseByIds [delete]
func (SSDApi *SocialServiceDatabaseApi) DeleteSocialServiceDatabaseByIds(c *gin.Context) {
	IDs := c.QueryArray("IDs[]")
    userID := utils.GetUserID(c)
	if err := SSDService.DeleteSocialServiceDatabaseByIds(IDs,userID); err != nil {
        global.GVA_LOG.Error("批量删除失败!", zap.Error(err))
		response.FailWithMessage("批量删除失败", c)
	} else {
		response.OkWithMessage("批量删除成功", c)
	}
}

// UpdateSocialServiceDatabase 更新社会服务库
// @Tags SocialServiceDatabase
// @Summary 更新社会服务库
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body EngineeringEducationDatabase.SocialServiceDatabase true "更新社会服务库"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /SSD/updateSocialServiceDatabase [put]
func (SSDApi *SocialServiceDatabaseApi) UpdateSocialServiceDatabase(c *gin.Context) {
	var SSD EngineeringEducationDatabase.SocialServiceDatabase
	err := c.ShouldBindJSON(&SSD)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
    SSD.UpdatedBy = utils.GetUserID(c)

	if err := SSDService.UpdateSocialServiceDatabase(SSD); err != nil {
        global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败", c)
	} else {
		response.OkWithMessage("更新成功", c)
	}
}

// FindSocialServiceDatabase 用id查询社会服务库
// @Tags SocialServiceDatabase
// @Summary 用id查询社会服务库
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query EngineeringEducationDatabase.SocialServiceDatabase true "用id查询社会服务库"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /SSD/findSocialServiceDatabase [get]
func (SSDApi *SocialServiceDatabaseApi) FindSocialServiceDatabase(c *gin.Context) {
	ID := c.Query("ID")
	if reSSD, err := SSDService.GetSocialServiceDatabase(ID); err != nil {
        global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败", c)
	} else {
		response.OkWithData(gin.H{"reSSD": reSSD}, c)
	}
}

// GetSocialServiceDatabaseList 分页获取社会服务库列表
// @Tags SocialServiceDatabase
// @Summary 分页获取社会服务库列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query EngineeringEducationDatabaseReq.SocialServiceDatabaseSearch true "分页获取社会服务库列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /SSD/getSocialServiceDatabaseList [get]
func (SSDApi *SocialServiceDatabaseApi) GetSocialServiceDatabaseList(c *gin.Context) {
	var pageInfo EngineeringEducationDatabaseReq.SocialServiceDatabaseSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if list, total, err := SSDService.GetSocialServiceDatabaseInfoList(pageInfo); err != nil {
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
