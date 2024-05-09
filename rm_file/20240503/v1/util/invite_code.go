package util

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
    "github.com/flipped-aurora/gin-vue-admin/server/model/util"
    utilReq "github.com/flipped-aurora/gin-vue-admin/server/model/util/request"
    "github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
    "github.com/flipped-aurora/gin-vue-admin/server/service"
    "github.com/gin-gonic/gin"
    "go.uber.org/zap"
    "github.com/flipped-aurora/gin-vue-admin/server/utils"
)

type InviteCodeApi struct {
}

var inviteCodeService = service.ServiceGroupApp.UtilServiceGroup.InviteCodeService


// CreateInviteCode 创建邀请码
// @Tags InviteCode
// @Summary 创建邀请码
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body util.InviteCode true "创建邀请码"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"创建成功"}"
// @Router /inviteCode/createInviteCode [post]
func (inviteCodeApi *InviteCodeApi) CreateInviteCode(c *gin.Context) {
	var inviteCode util.InviteCode
	err := c.ShouldBindJSON(&inviteCode)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
    inviteCode.CreatedBy = utils.GetUserID(c)

	if err := inviteCodeService.CreateInviteCode(&inviteCode); err != nil {
        global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败", c)
	} else {
		response.OkWithMessage("创建成功", c)
	}
}

// DeleteInviteCode 删除邀请码
// @Tags InviteCode
// @Summary 删除邀请码
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body util.InviteCode true "删除邀请码"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /inviteCode/deleteInviteCode [delete]
func (inviteCodeApi *InviteCodeApi) DeleteInviteCode(c *gin.Context) {
	ID := c.Query("ID")
    	userID := utils.GetUserID(c)
	if err := inviteCodeService.DeleteInviteCode(ID,userID); err != nil {
        global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败", c)
	} else {
		response.OkWithMessage("删除成功", c)
	}
}

// DeleteInviteCodeByIds 批量删除邀请码
// @Tags InviteCode
// @Summary 批量删除邀请码
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"批量删除成功"}"
// @Router /inviteCode/deleteInviteCodeByIds [delete]
func (inviteCodeApi *InviteCodeApi) DeleteInviteCodeByIds(c *gin.Context) {
	IDs := c.QueryArray("IDs[]")
    userID := utils.GetUserID(c)
	if err := inviteCodeService.DeleteInviteCodeByIds(IDs,userID); err != nil {
        global.GVA_LOG.Error("批量删除失败!", zap.Error(err))
		response.FailWithMessage("批量删除失败", c)
	} else {
		response.OkWithMessage("批量删除成功", c)
	}
}

// UpdateInviteCode 更新邀请码
// @Tags InviteCode
// @Summary 更新邀请码
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body util.InviteCode true "更新邀请码"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /inviteCode/updateInviteCode [put]
func (inviteCodeApi *InviteCodeApi) UpdateInviteCode(c *gin.Context) {
	var inviteCode util.InviteCode
	err := c.ShouldBindJSON(&inviteCode)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
    inviteCode.UpdatedBy = utils.GetUserID(c)

	if err := inviteCodeService.UpdateInviteCode(inviteCode); err != nil {
        global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败", c)
	} else {
		response.OkWithMessage("更新成功", c)
	}
}

// FindInviteCode 用id查询邀请码
// @Tags InviteCode
// @Summary 用id查询邀请码
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query util.InviteCode true "用id查询邀请码"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /inviteCode/findInviteCode [get]
func (inviteCodeApi *InviteCodeApi) FindInviteCode(c *gin.Context) {
	ID := c.Query("ID")
	if reinviteCode, err := inviteCodeService.GetInviteCode(ID); err != nil {
        global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败", c)
	} else {
		response.OkWithData(gin.H{"reinviteCode": reinviteCode}, c)
	}
}

// GetInviteCodeList 分页获取邀请码列表
// @Tags InviteCode
// @Summary 分页获取邀请码列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query utilReq.InviteCodeSearch true "分页获取邀请码列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /inviteCode/getInviteCodeList [get]
func (inviteCodeApi *InviteCodeApi) GetInviteCodeList(c *gin.Context) {
	var pageInfo utilReq.InviteCodeSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if list, total, err := inviteCodeService.GetInviteCodeInfoList(pageInfo); err != nil {
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
