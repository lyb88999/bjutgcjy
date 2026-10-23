package ExpertDatabase

import (
	"context"
	"errors"
	"net/http"
	"time"

	ExpertDatabaseReq "github.com/flipped-aurora/gin-vue-admin/server/model/ExpertDatabase/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/service"
	ExpertDatabaseService "github.com/flipped-aurora/gin-vue-admin/server/service/ExpertDatabase"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/gin-gonic/gin"
)

type ExpertAgentApi struct{}

var expertAgentService = service.ServiceGroupApp.ExpertDatabaseServiceGroup.ExpertAgentService

func (*ExpertAgentApi) GetAgentStatus(c *gin.Context) {
	response.OkWithData(expertAgentService.Status(), c)
}

func (*ExpertAgentApi) ChatWithAgent(c *gin.Context) {
	operatorID := utils.GetUserID(c)
	if operatorID == 0 {
		response.FailWithMessage("请先登录", c)
		return
	}
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 16<<10)
	var req ExpertDatabaseReq.ExpertAgentChatReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage("请提交有效的遴选需求，请求内容不能超过 16KB", c)
		return
	}
	ctx, cancel := context.WithTimeout(c.Request.Context(), 80*time.Second)
	defer cancel()
	result, err := expertAgentService.Chat(ctx, operatorID, req)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			return
		}
		if errors.Is(err, context.DeadlineExceeded) {
			response.FailWithMessage("本次遴选超时，请稍后重试或缩小研究主题", c)
			return
		}
		for _, known := range []error{
			ExpertDatabaseService.ErrExpertAgentInvalidRequest,
			ExpertDatabaseService.ErrExpertAgentNotConfigured,
			ExpertDatabaseService.ErrExpertAgentRateLimited,
			ExpertDatabaseService.ErrExpertAgentUpstream,
			ExpertDatabaseService.ErrExpertAgentSessionNotFound,
			ExpertDatabaseService.ErrExpertAgentSessionBusy,
			ExpertDatabaseService.ErrExpertAgentSessionLimit,
		} {
			if errors.Is(err, known) {
				response.FailWithMessage(known.Error(), c)
				return
			}
		}
		// SQL and provider error details can contain private data; expose only known messages.
		response.FailWithMessage("专家遴选暂时未能完成，请稍后重试", c)
		return
	}
	response.OkWithData(result, c)
}
