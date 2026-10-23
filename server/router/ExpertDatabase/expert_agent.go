package ExpertDatabase

import (
	v1 "github.com/flipped-aurora/gin-vue-admin/server/api/v1"
	"github.com/gin-gonic/gin"
)

type ExpertAgentRouter struct{}

func (*ExpertAgentRouter) InitExpertAgentRouter(router *gin.RouterGroup) {
	group := router.Group("expertDatabase/agent")
	api := v1.ApiGroupApp.ExpertDatabaseApiGroup.ExpertAgentApi
	// These routes inherit JWT and Casbin. Conversation bodies are not copied to operation logs.
	group.GET("status", api.GetAgentStatus)
	group.POST("chat", api.ChatWithAgent)
}
