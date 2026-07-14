package system

import (
	"strconv"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type OrganizationApi struct{}

// CreateSysOrganization
// @Tags      SysOrganization
// @Summary   创建单位
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      system.SysOrganization         true  "单位模型"
// @Success   200   {object}  response.Response{msg=string}  "创建单位"
// @Router    /sysOrganization/createSysOrganization [post]
func (o *OrganizationApi) CreateSysOrganization(c *gin.Context) {
	var org system.SysOrganization
	err := c.ShouldBindJSON(&org)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = organizationService.CreateSysOrganization(org)
	if err != nil {
		global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("创建成功", c)
}

// DeleteSysOrganization
// @Tags      SysOrganization
// @Summary   删除单位
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     id  query     int                             true  "单位ID"
// @Success   200 {object}  response.Response{msg=string}  "删除单位"
// @Router    /sysOrganization/deleteSysOrganization [delete]
func (o *OrganizationApi) DeleteSysOrganization(c *gin.Context) {
	id, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		response.FailWithMessage("id 参数错误", c)
		return
	}
	err = organizationService.DeleteSysOrganization(uint(id))
	if err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// UpdateSysOrganization
// @Tags      SysOrganization
// @Summary   更新单位
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      system.SysOrganization          true  "单位模型"
// @Success   200   {object}  response.Response{msg=string}  "更新单位"
// @Router    /sysOrganization/updateSysOrganization [put]
func (o *OrganizationApi) UpdateSysOrganization(c *gin.Context) {
	var org system.SysOrganization
	err := c.ShouldBindJSON(&org)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = organizationService.UpdateSysOrganization(&org)
	if err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("更新成功", c)
}

// FindSysOrganization
// @Tags      SysOrganization
// @Summary   用id查询单位
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     id  query     int                                                        true  "单位ID"
// @Success   200 {object}  response.Response{data=map[string]interface{},msg=string}  "用id查询单位"
// @Router    /sysOrganization/findSysOrganization [get]
func (o *OrganizationApi) FindSysOrganization(c *gin.Context) {
	id, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		response.FailWithMessage("id 参数错误", c)
		return
	}
	org, err := organizationService.GetSysOrganization(uint(id))
	if err != nil {
		global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败", c)
		return
	}
	response.OkWithDetailed(gin.H{"sysOrganization": org}, "查询成功", c)
}

// GetSysOrganizationList
// @Tags      SysOrganization
// @Summary   分页获取单位列表
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  query     request.PageInfo                                        true  "分页获取单位列表"
// @Success   200   {object}  response.Response{data=response.PageResult,msg=string}  "分页获取单位列表"
// @Router    /sysOrganization/getSysOrganizationList [get]
func (o *OrganizationApi) GetSysOrganizationList(c *gin.Context) {
	var pageInfo request.PageInfo
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := organizationService.GetSysOrganizationInfoList(pageInfo)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, "获取成功", c)
}

// GetSysOrganizationTree
// @Tags      SysOrganization
// @Summary   获取单位树，供审核路由/表单级联选择使用
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Success   200 {object}  response.Response{data=map[string]interface{},msg=string}  "获取单位树"
// @Router    /sysOrganization/getSysOrganizationTree [get]
func (o *OrganizationApi) GetSysOrganizationTree(c *gin.Context) {
	tree, err := organizationService.GetSysOrganizationTree()
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithDetailed(gin.H{"tree": tree}, "获取成功", c)
}
