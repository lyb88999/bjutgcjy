package EngineeringEducationDatabase

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
    "github.com/flipped-aurora/gin-vue-admin/server/model/EngineeringEducationDatabase"
    EngineeringEducationDatabaseReq "github.com/flipped-aurora/gin-vue-admin/server/model/EngineeringEducationDatabase/request"
    "github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
    "github.com/flipped-aurora/gin-vue-admin/server/service"
    "github.com/gin-gonic/gin"
    "go.uber.org/zap"
)

type SubjectConstructionDatabaseApi struct {
}

var SCDService = service.ServiceGroupApp.EngineeringEducationDatabaseServiceGroup.SubjectConstructionDatabaseService


// CreateSubjectConstructionDatabase 创建学科建设库
// @Tags SubjectConstructionDatabase
// @Summary 创建学科建设库
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body EngineeringEducationDatabase.SubjectConstructionDatabase true "创建学科建设库"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"创建成功"}"
// @Router /SCD/createSubjectConstructionDatabase [post]
func (SCDApi *SubjectConstructionDatabaseApi) CreateSubjectConstructionDatabase(c *gin.Context) {
	var SCD EngineeringEducationDatabase.SubjectConstructionDatabase
	err := c.ShouldBindJSON(&SCD)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if err := SCDService.CreateSubjectConstructionDatabase(&SCD); err != nil {
        global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败", c)
	} else {
		response.OkWithMessage("创建成功", c)
	}
}

// DeleteSubjectConstructionDatabase 删除学科建设库
// @Tags SubjectConstructionDatabase
// @Summary 删除学科建设库
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body EngineeringEducationDatabase.SubjectConstructionDatabase true "删除学科建设库"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /SCD/deleteSubjectConstructionDatabase [delete]
func (SCDApi *SubjectConstructionDatabaseApi) DeleteSubjectConstructionDatabase(c *gin.Context) {
	ID := c.Query("ID")
	if err := SCDService.DeleteSubjectConstructionDatabase(ID); err != nil {
        global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败", c)
	} else {
		response.OkWithMessage("删除成功", c)
	}
}

// DeleteSubjectConstructionDatabaseByIds 批量删除学科建设库
// @Tags SubjectConstructionDatabase
// @Summary 批量删除学科建设库
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"批量删除成功"}"
// @Router /SCD/deleteSubjectConstructionDatabaseByIds [delete]
func (SCDApi *SubjectConstructionDatabaseApi) DeleteSubjectConstructionDatabaseByIds(c *gin.Context) {
	IDs := c.QueryArray("IDs[]")
	if err := SCDService.DeleteSubjectConstructionDatabaseByIds(IDs); err != nil {
        global.GVA_LOG.Error("批量删除失败!", zap.Error(err))
		response.FailWithMessage("批量删除失败", c)
	} else {
		response.OkWithMessage("批量删除成功", c)
	}
}

// UpdateSubjectConstructionDatabase 更新学科建设库
// @Tags SubjectConstructionDatabase
// @Summary 更新学科建设库
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body EngineeringEducationDatabase.SubjectConstructionDatabase true "更新学科建设库"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /SCD/updateSubjectConstructionDatabase [put]
func (SCDApi *SubjectConstructionDatabaseApi) UpdateSubjectConstructionDatabase(c *gin.Context) {
	var SCD EngineeringEducationDatabase.SubjectConstructionDatabase
	err := c.ShouldBindJSON(&SCD)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if err := SCDService.UpdateSubjectConstructionDatabase(SCD); err != nil {
        global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败", c)
	} else {
		response.OkWithMessage("更新成功", c)
	}
}

// FindSubjectConstructionDatabase 用id查询学科建设库
// @Tags SubjectConstructionDatabase
// @Summary 用id查询学科建设库
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query EngineeringEducationDatabase.SubjectConstructionDatabase true "用id查询学科建设库"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /SCD/findSubjectConstructionDatabase [get]
func (SCDApi *SubjectConstructionDatabaseApi) FindSubjectConstructionDatabase(c *gin.Context) {
	ID := c.Query("ID")
	if reSCD, err := SCDService.GetSubjectConstructionDatabase(ID); err != nil {
        global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败", c)
	} else {
		response.OkWithData(gin.H{"reSCD": reSCD}, c)
	}
}

// GetSubjectConstructionDatabaseList 分页获取学科建设库列表
// @Tags SubjectConstructionDatabase
// @Summary 分页获取学科建设库列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query EngineeringEducationDatabaseReq.SubjectConstructionDatabaseSearch true "分页获取学科建设库列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /SCD/getSubjectConstructionDatabaseList [get]
func (SCDApi *SubjectConstructionDatabaseApi) GetSubjectConstructionDatabaseList(c *gin.Context) {
	var pageInfo EngineeringEducationDatabaseReq.SubjectConstructionDatabaseSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if list, total, err := SCDService.GetSubjectConstructionDatabaseInfoList(pageInfo); err != nil {
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
