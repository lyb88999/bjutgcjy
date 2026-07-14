package ExpertDatabase

import (
	"strconv"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/ExpertDatabase"
	ExpertDatabaseReq "github.com/flipped-aurora/gin-vue-admin/server/model/ExpertDatabase/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/service"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ExpertTagApi struct{}

var expertTagService = service.ServiceGroupApp.ExpertDatabaseServiceGroup.ExpertTagService

// CreateExpertTag 创建标签
// @Tags ExpertTag
// @Summary 创建标签
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body ExpertDatabase.ExpertTag true "创建标签"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"创建成功"}"
// @Router /expertTag/createExpertTag [post]
func (expertTagApi *ExpertTagApi) CreateExpertTag(c *gin.Context) {
	var tag ExpertDatabase.ExpertTag
	err := c.ShouldBindJSON(&tag)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	tag.CreatedBy = utils.GetUserID(c)
	if err := expertTagService.CreateExpertTag(&tag); err != nil {
		global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败", c)
	} else {
		response.OkWithMessage("创建成功", c)
	}
}

// DeleteExpertTag 删除标签
// @Tags ExpertTag
// @Summary 删除标签
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /expertTag/deleteExpertTag [delete]
func (expertTagApi *ExpertTagApi) DeleteExpertTag(c *gin.Context) {
	ID := c.Query("ID")
	userID := utils.GetUserID(c)
	if err := expertTagService.DeleteExpertTag(ID, userID); err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败", c)
	} else {
		response.OkWithMessage("删除成功", c)
	}
}

// UpdateExpertTag 更新标签
// @Tags ExpertTag
// @Summary 更新标签
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body ExpertDatabase.ExpertTag true "更新标签"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /expertTag/updateExpertTag [put]
func (expertTagApi *ExpertTagApi) UpdateExpertTag(c *gin.Context) {
	var tag ExpertDatabase.ExpertTag
	err := c.ShouldBindJSON(&tag)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	tag.UpdatedBy = utils.GetUserID(c)
	if err := expertTagService.UpdateExpertTag(tag); err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败", c)
	} else {
		response.OkWithMessage("更新成功", c)
	}
}

// FindExpertTag 用id查询标签
// @Tags ExpertTag
// @Summary 用id查询标签
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /expertTag/findExpertTag [get]
func (expertTagApi *ExpertTagApi) FindExpertTag(c *gin.Context) {
	ID := c.Query("ID")
	if tag, err := expertTagService.GetExpertTag(ID); err != nil {
		global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败", c)
	} else {
		response.OkWithData(gin.H{"reExpertTag": tag}, c)
	}
}

// GetExpertTagList 分页获取标签列表
// @Tags ExpertTag
// @Summary 分页获取标签列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query ExpertDatabaseReq.ExpertTagSearch true "分页获取标签列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /expertTag/getExpertTagList [get]
func (expertTagApi *ExpertTagApi) GetExpertTagList(c *gin.Context) {
	var pageInfo ExpertDatabaseReq.ExpertTagSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if list, total, err := expertTagService.GetExpertTagInfoList(pageInfo); err != nil {
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

// SetExpertTagRelations 覆盖式设置某专家的标签关联
// @Tags ExpertTag
// @Summary 设置专家标签关联
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body ExpertDatabaseReq.ExpertTagRelationSetReq true "设置专家标签关联"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"设置成功"}"
// @Router /expertTag/setExpertTagRelations [post]
func (expertTagApi *ExpertTagApi) SetExpertTagRelations(c *gin.Context) {
	var req ExpertDatabaseReq.ExpertTagRelationSetReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := expertTagService.SetExpertTags(req.ExpertId, req.TagIds); err != nil {
		global.GVA_LOG.Error("设置失败!", zap.Error(err))
		response.FailWithMessage("设置失败", c)
	} else {
		response.OkWithMessage("设置成功", c)
	}
}

// GetExpertTagsByExpertId 获取某专家关联的全部标签
// @Tags ExpertTag
// @Summary 获取专家标签
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /expertTag/getExpertTagsByExpertId [get]
func (expertTagApi *ExpertTagApi) GetExpertTagsByExpertId(c *gin.Context) {
	expertId, err := strconv.Atoi(c.Query("expertId"))
	if err != nil {
		response.FailWithMessage("expertId 参数错误", c)
		return
	}
	tags, err := expertTagService.GetExpertTagsByExpertId(uint(expertId))
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithDetailed(gin.H{"tags": tags}, "获取成功", c)
}
