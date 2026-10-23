package ExpertDatabase

import (
	"encoding/json"
	"net/http/httptest"
	"strings"
	"testing"

	systemReq "github.com/flipped-aurora/gin-vue-admin/server/model/system/request"
	"github.com/gin-gonic/gin"
)

func TestAgentChatRejectsMissingIdentityAndOversizedBody(t *testing.T) {
	gin.SetMode(gin.TestMode)
	for _, test := range []struct {
		name string
		user uint
		body string
		want string
	}{
		{"missing identity", 0, `{"message":"找养老专家"}`, "请先登录"},
		{"oversized request", 1, `{"message":"` + strings.Repeat("x", 17<<10) + `"}`, "16KB"},
		{"invalid request", 1, `{`, "有效的遴选需求"},
	} {
		t.Run(test.name, func(t *testing.T) {
			writer := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(writer)
			claims := &systemReq.CustomClaims{}
			claims.BaseClaims.ID = test.user
			ctx.Set("claims", claims)
			ctx.Request = httptest.NewRequest("POST", "/expertDatabase/agent/chat", strings.NewReader(test.body))
			ctx.Request.Header.Set("Content-Type", "application/json")
			(&ExpertAgentApi{}).ChatWithAgent(ctx)
			var result struct {
				Code int    `json:"code"`
				Msg  string `json:"msg"`
			}
			if err := json.Unmarshal(writer.Body.Bytes(), &result); err != nil || result.Code == 0 || !strings.Contains(result.Msg, test.want) {
				t.Fatalf("request must be rejected before model access: %s", writer.Body.String())
			}
		})
	}
}
