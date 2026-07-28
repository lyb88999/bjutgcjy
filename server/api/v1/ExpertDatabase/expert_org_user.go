package ExpertDatabase

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	ExpertDatabaseReq "github.com/flipped-aurora/gin-vue-admin/server/model/ExpertDatabase/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/service"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ExpertOrgUserApi struct{}

var expertOrgUserService = service.ServiceGroupApp.ExpertDatabaseServiceGroup.ExpertOrgUserService

// CreateOrgApplicant 单位审核员为本单位新建个人申报人账号
// @Tags ExpertOrgUser
// @Summary 单位审核员为本单位新建个人申报人账号
// @Security ApiKeyAuth
// @Router /expertOrgUser/createOrgApplicant [post]
func (a *ExpertOrgUserApi) CreateOrgApplicant(c *gin.Context) {
	var req ExpertDatabaseReq.CreateOrgApplicantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if len(req.Password) < 6 {
		response.FailWithMessage("密码至少6位", c)
		return
	}
	if err := expertOrgUserService.CreateOrgApplicant(utils.GetUserID(c), req); err != nil {
		global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("创建成功", c)
}

// GetOrgUserList 分页获取本单位个人申报人账号列表
// @Tags ExpertOrgUser
// @Summary 分页获取本单位个人申报人账号列表
// @Security ApiKeyAuth
// @Router /expertOrgUser/getOrgUserList [get]
func (a *ExpertOrgUserApi) GetOrgUserList(c *gin.Context) {
	var pageInfo ExpertDatabaseReq.OrgUserSearch
	if err := c.ShouldBindQuery(&pageInfo); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := expertOrgUserService.GetOrgUserList(utils.GetUserID(c), pageInfo)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List: list, Total: total, Page: pageInfo.Page, PageSize: pageInfo.PageSize,
	}, "获取成功", c)
}

// ToggleOrgUserEnable 启用/冻结本单位个人申报人账号
// @Tags ExpertOrgUser
// @Summary 启用/冻结本单位个人申报人账号
// @Security ApiKeyAuth
// @Router /expertOrgUser/toggleOrgUserEnable [post]
func (a *ExpertOrgUserApi) ToggleOrgUserEnable(c *gin.Context) {
	var req ExpertDatabaseReq.ToggleOrgUserEnableRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := expertOrgUserService.ToggleOrgUserEnable(utils.GetUserID(c), req.ID, req.Enable); err != nil {
		global.GVA_LOG.Error("操作失败!", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("操作成功", c)
}
