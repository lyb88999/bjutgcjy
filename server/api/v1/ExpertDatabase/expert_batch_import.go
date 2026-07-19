package ExpertDatabase

import (
	"net/http"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// DownloadImportTemplate 下载专家批量导入模板
// @Tags ExpertProfile
// @Summary 下载专家批量导入模板
// @Security ApiKeyAuth
// @Produce application/octet-stream
// @Router /expertProfile/downloadImportTemplate [get]
func (expertProfileApi *ExpertProfileApi) DownloadImportTemplate(c *gin.Context) {
	f, err := expertProfileService.BuildImportTemplate()
	if err != nil {
		global.GVA_LOG.Error("生成模板失败!", zap.Error(err))
		response.FailWithMessage("生成模板失败", c)
		return
	}
	buf, err := f.WriteToBuffer()
	if err != nil {
		global.GVA_LOG.Error("生成模板失败!", zap.Error(err))
		response.FailWithMessage("生成模板失败", c)
		return
	}
	c.Header("Content-Disposition", "attachment; filename=expert_batch_import_template.xlsx")
	c.Header("success", "true")
	c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", buf.Bytes())
}

// ImportBatch 批量导入专家背景信息与研究成果
// @Tags ExpertProfile
// @Summary 批量导入专家背景信息与研究成果
// @Security ApiKeyAuth
// @accept multipart/form-data
// @Produce application/json
// @Router /expertProfile/importBatch [post]
func (expertProfileApi *ExpertProfileApi) ImportBatch(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		response.FailWithMessage("文件获取失败", c)
		return
	}
	src, err := file.Open()
	if err != nil {
		response.FailWithMessage("文件打开失败", c)
		return
	}
	defer src.Close()

	result, err := expertProfileService.ImportBatch(src, utils.GetUserID(c))
	if err != nil {
		global.GVA_LOG.Error("批量导入失败!", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	if !result.Success {
		response.FailWithDetailed(result, "校验未通过，未导入任何数据，请按提示修改后重新上传", c)
		return
	}
	response.OkWithDetailed(result, "导入成功", c)
}
