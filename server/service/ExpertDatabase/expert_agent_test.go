package ExpertDatabase

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/ExpertDatabase"
	ExpertDatabaseReq "github.com/flipped-aurora/gin-vue-admin/server/model/ExpertDatabase/request"
	ExpertDatabaseRes "github.com/flipped-aurora/gin-vue-admin/server/model/ExpertDatabase/response"
)

// Every provider call stays on this local mock. A fresh store per test prevents
// ownership, concurrency and rate-limit state from leaking between test cases.
func newTestExpertAgent(t *testing.T, handler http.HandlerFunc) *ExpertAgentService {
	t.Helper()
	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)
	t.Setenv("EXPERT_AGENT_BASE_URL", server.URL+"/v1")
	t.Setenv("EXPERT_AGENT_API_KEY", "test-agent-key-not-a-real-secret")
	t.Setenv("EXPERT_AGENT_MODEL", "test-agent-model")
	return &ExpertAgentService{client: server.Client(), sessions: newExpertAgentSessionStore()}
}

func writeTestExpertAgentPlan(t *testing.T, w http.ResponseWriter, arguments string) {
	t.Helper()
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(map[string]any{
		"choices": []any{map[string]any{"message": map[string]any{
			"role": "assistant",
			// A provider's prose is untrusted; only database candidates may become results.
			"content": "推荐不存在的专家：ID 999999，姓名为模型编造姓名",
			"tool_calls": []any{map[string]any{
				"id": "call_search_1", "type": "function",
				"function": map[string]any{"name": "search_experts", "arguments": arguments},
			}},
		}}},
	})
	if err != nil {
		t.Errorf("write mock provider response: %v", err)
	}
}

func testExpertAgentCandidateIDs(candidates []ExpertDatabaseRes.ExpertAgentCandidate) []uint {
	ids := make([]uint, len(candidates))
	for i := range candidates {
		ids[i] = candidates[i].ID
	}
	return ids
}

func TestExpertAgentChat_ExtractsToolPlanAndEnforcesSelectionConstraints(t *testing.T) {
	var calls atomic.Int32
	service := newTestExpertAgent(t, func(w http.ResponseWriter, r *http.Request) {
		calls.Add(1)
		if r.Method != http.MethodPost || r.URL.Path != "/v1/chat/completions" {
			t.Errorf("unexpected provider endpoint: %s %s", r.Method, r.URL.Path)
		}
		if r.Header.Get("Authorization") != "Bearer test-agent-key-not-a-real-secret" {
			t.Error("provider request must authenticate with the server-side test key")
		}
		var payload struct {
			Model    string `json:"model"`
			Messages []struct {
				Role    string `json:"role"`
				Content string `json:"content"`
			} `json:"messages"`
			Tools []struct {
				Function struct {
					Name string `json:"name"`
				} `json:"function"`
			} `json:"tools"`
		}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Errorf("decode provider request: %v", err)
		}
		if payload.Model != "test-agent-model" || len(payload.Messages) == 0 {
			t.Errorf("provider request missing selected model or conversation: %+v", payload)
		}
		if len(payload.Tools) != 1 || payload.Tools[0].Function.Name != "search_experts" {
			t.Errorf("provider must receive the bounded search tool: %+v", payload.Tools)
		}
		writeTestExpertAgentPlan(t, w, `{"keywords":["基层治理"],"disciplineL1":"政治学","titleTier":"教授","count":5,"maxPerUnit":3,"excludedExpertIds":[4],"diversify":true}`)
	})
	var gotPlan ExpertDatabaseRes.ExpertAgentPlan
	service.search = func(ctx context.Context, plan ExpertDatabaseRes.ExpertAgentPlan) ([]ExpertDatabaseRes.ExpertAgentCandidate, int64, error) {
		gotPlan = plan
		return []ExpertDatabaseRes.ExpertAgentCandidate{
			{ID: 1, Name: "甲", UnitName: "甲大学", OrgID: 7, Score: 10},
			{ID: 1, Name: "重复甲", UnitName: "不同拼写单位", Score: 100},
			{ID: 0, Name: "无效编号", UnitName: "无效单位", Score: 100},
			// The same organization may have different imported display names;
			// its stable ID must still enforce the one-per-unit constraint.
			{ID: 2, Name: "乙", UnitName: "甲大学旧名称", OrgID: 7, Score: 9},
			{ID: 3, Name: "丙", UnitName: "乙大学", OrgID: 8, Score: 8},
			{ID: 4, Name: "丁", UnitName: "丙大学", Score: 7},
		}, 4, nil
	}
	result, err := service.Chat(context.Background(), 11, ExpertDatabaseReq.ExpertAgentChatReq{Message: "找两位基层治理方向教授，每个单位最多一位", Count: 2, MaxPerUnit: 1})
	if err != nil {
		t.Fatalf("chat failed: %v", err)
	}
	if !reflect.DeepEqual(gotPlan.Keywords, []string{"基层治理"}) || gotPlan.DisciplineL1 != "政治学" || gotPlan.TitleTier != "教授" {
		t.Fatalf("structured filters were not passed to database search: %+v", gotPlan)
	}
	if gotPlan.Count != 2 || gotPlan.MaxPerUnit != 1 {
		t.Fatalf("explicit user controls must override model selection limits: %+v", gotPlan)
	}
	if len(gotPlan.ExcludedExpertIDs) != 0 {
		t.Fatalf("model must not exclude IDs it has not seen in this user's session: %+v", gotPlan.ExcludedExpertIDs)
	}
	if got := testExpertAgentCandidateIDs(result.Recommendations); !reflect.DeepEqual(got, []uint{1, 3}) {
		t.Fatalf("selection must enforce count, unit limit and excluded IDs, got %v", got)
	}
	if calls.Load() != 1 || result.SessionID == "" {
		t.Fatalf("expected one provider turn and a reusable session, calls=%d session=%q", calls.Load(), result.SessionID)
	}
}

func TestExpertAgentChat_OnlyPublishedDatabaseCandidatesCanBeRecommended(t *testing.T) {
	db := setupTestDB(t)
	operator := createTestUser(t, db, superAdminAuthorityId, 0)
	oldEmbeddingURL := global.GVA_CONFIG.ExpertEmbedding.Url
	global.GVA_CONFIG.ExpertEmbedding.Url = ""
	t.Cleanup(func() { global.GVA_CONFIG.ExpertEmbedding.Url = oldEmbeddingURL })
	profiles := []ExpertDatabase.ExpertProfile{
		{Name: "已发布甲", Status: StatusPublished, UnitName: "甲大学", ResearchDirections: "基层治理", CompositeScore: 10, Phone: "private-office-phone", Mobile: "private-mobile", Email: "private-email@example.invalid"},
		{Name: "已发布乙", Status: StatusPublished, UnitName: "乙大学", ResearchDirections: "基层治理", CompositeScore: 8},
		{Name: "未公开草稿", Status: StatusDraft, UnitName: "丙大学", ResearchDirections: "基层治理", CompositeScore: 1000},
		{Name: "尚未审核", Status: StatusPendingCityReview, UnitName: "丁大学", ResearchDirections: "基层治理", CompositeScore: 2000},
	}
	if err := db.Create(&profiles).Error; err != nil {
		t.Fatalf("seed database profiles: %v", err)
	}
	service := newTestExpertAgent(t, func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Errorf("read model input: %v", err)
		}
		if strings.Contains(string(body), "private-office-phone") || strings.Contains(string(body), "private-mobile") || strings.Contains(string(body), "private-email@example.invalid") {
			t.Error("contact details must not be transmitted to the model")
		}
		writeTestExpertAgentPlan(t, w, `{"count":10}`)
	})
	result, err := service.Chat(context.Background(), operator.ID, ExpertDatabaseReq.ExpertAgentChatReq{Message: "浏览全部已发布专家"})
	if err != nil {
		t.Fatalf("chat failed: %v", err)
	}
	if len(result.Candidates) != 2 {
		t.Fatalf("only the two published database rows may be returned: %+v", result.Candidates)
	}
	allowed := map[uint]string{profiles[0].ID: profiles[0].Name, profiles[1].ID: profiles[1].Name}
	for _, candidate := range result.Candidates {
		if expectedName, ok := allowed[candidate.ID]; !ok || expectedName != candidate.Name {
			t.Fatalf("candidate was not a published database record: %+v", candidate)
		}
	}
	encoded, err := json.Marshal(result)
	if err != nil {
		t.Fatalf("marshal result: %v", err)
	}
	if strings.Contains(string(encoded), "999999") || strings.Contains(string(encoded), "模型编造姓名") || strings.Contains(string(encoded), "未公开草稿") || strings.Contains(string(encoded), "private-") {
		t.Fatalf("model prose or private records leaked into the final result: %s", encoded)
	}
}

func TestExpertAgentChat_RejectsOtherUsersSessionBeforeCallingModel(t *testing.T) {
	var calls atomic.Int32
	service := newTestExpertAgent(t, func(w http.ResponseWriter, r *http.Request) {
		calls.Add(1)
		writeTestExpertAgentPlan(t, w, `{"keywords":["基层治理"],"count":1}`)
	})
	service.search = func(context.Context, ExpertDatabaseRes.ExpertAgentPlan) ([]ExpertDatabaseRes.ExpertAgentCandidate, int64, error) {
		return []ExpertDatabaseRes.ExpertAgentCandidate{{ID: 1, Name: "专家", UnitName: "甲大学"}}, 1, nil
	}
	first, err := service.Chat(context.Background(), 11, ExpertDatabaseReq.ExpertAgentChatReq{Message: "找一位基层治理专家"})
	if err != nil {
		t.Fatalf("create owner session: %v", err)
	}
	_, err = service.Chat(context.Background(), 22, ExpertDatabaseReq.ExpertAgentChatReq{SessionID: first.SessionID, Message: "继续"})
	if !errors.Is(err, ErrExpertAgentSessionNotFound) {
		t.Fatalf("another user must not read or reuse this session, got %v", err)
	}
	if calls.Load() != 1 {
		t.Fatalf("cross-user session request reached provider, calls=%d", calls.Load())
	}
}

func TestExpertAgentChat_Provider429DoesNotExposeProviderBodyOrCredentials(t *testing.T) {
	service := newTestExpertAgent(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTooManyRequests)
		_, _ = w.Write([]byte(`{"error":"provider-private-debug test-agent-key-not-a-real-secret"}`))
	})
	_, err := service.Chat(context.Background(), 11, ExpertDatabaseReq.ExpertAgentChatReq{Message: "找一位专家"})
	if !errors.Is(err, ErrExpertAgentRateLimited) {
		t.Fatalf("provider 429 should become the safe rate-limit error, got %v", err)
	}
	if strings.Contains(err.Error(), "provider-private-debug") || strings.Contains(err.Error(), "test-agent-key") {
		t.Fatalf("provider details or credentials exposed: %v", err)
	}
}

func TestExpertAgentChat_FollowUpPreservesFiltersAndExcludesPreviousResults(t *testing.T) {
	var calls atomic.Int32
	service := newTestExpertAgent(t, func(w http.ResponseWriter, r *http.Request) {
		if calls.Add(1) == 1 {
			writeTestExpertAgentPlan(t, w, `{"keywords":["基层治理"],"disciplineL1":"政治学","titleTier":"教授","count":1,"maxPerUnit":1}`)
			return
		}
		writeTestExpertAgentPlan(t, w, `{"excludePrevious":true}`)
	})
	var plans []ExpertDatabaseRes.ExpertAgentPlan
	service.search = func(ctx context.Context, plan ExpertDatabaseRes.ExpertAgentPlan) ([]ExpertDatabaseRes.ExpertAgentCandidate, int64, error) {
		plans = append(plans, plan)
		return []ExpertDatabaseRes.ExpertAgentCandidate{
			{ID: 1, Name: "甲", UnitName: "甲大学", Score: 10},
			{ID: 2, Name: "乙", UnitName: "乙大学", Score: 9},
		}, 2, nil
	}
	first, err := service.Chat(context.Background(), 11, ExpertDatabaseReq.ExpertAgentChatReq{Message: "找一位研究基层治理的政治学教授"})
	if err != nil {
		t.Fatalf("first chat: %v", err)
	}
	second, err := service.Chat(context.Background(), 11, ExpertDatabaseReq.ExpertAgentChatReq{SessionID: first.SessionID, Message: "按刚才的条件换一位"})
	if err != nil {
		t.Fatalf("follow-up chat: %v", err)
	}
	if len(plans) != 2 {
		t.Fatalf("expected two searches, got %d", len(plans))
	}
	plan := plans[1]
	if !reflect.DeepEqual(plan.Keywords, []string{"基层治理"}) || plan.DisciplineL1 != "政治学" || plan.TitleTier != "教授" || plan.Count != 1 || plan.MaxPerUnit != 1 {
		t.Fatalf("follow-up lost previous filters: %+v", plan)
	}
	if !reflect.DeepEqual(plan.ExcludedExpertIDs, []uint{1}) {
		t.Fatalf("previous selection must be excluded, got %+v", plan.ExcludedExpertIDs)
	}
	if got := testExpertAgentCandidateIDs(second.Recommendations); !reflect.DeepEqual(got, []uint{2}) {
		t.Fatalf("follow-up should return a fresh real candidate, got %v", got)
	}
}

func TestExpertAgentChat_CanceledTurnDoesNotPolluteSessionOrLeaveItBusy(t *testing.T) {
	var calls atomic.Int32
	started := make(chan struct{})
	release := make(chan struct{})
	service := newTestExpertAgent(t, func(w http.ResponseWriter, r *http.Request) {
		switch calls.Add(1) {
		case 1:
			writeTestExpertAgentPlan(t, w, `{"keywords":["基层治理"],"count":1}`)
		case 2:
			close(started)
			select {
			case <-r.Context().Done():
			case <-release:
			}
		default:
			body, err := io.ReadAll(r.Body)
			if err != nil {
				t.Errorf("read follow-up model input: %v", err)
			}
			if strings.Contains(string(body), "取消前临时改成完全不同的主题") || strings.Contains(string(body), "并发继续") {
				t.Error("canceled or rejected turn polluted the committed model conversation")
			}
			writeTestExpertAgentPlan(t, w, `{"excludePrevious":true}`)
		}
	})
	// This cleanup executes before httptest.Close so a failing assertion cannot
	// leave the deliberately blocked mock handler hanging during cleanup.
	t.Cleanup(func() { close(release) })
	var plans []ExpertDatabaseRes.ExpertAgentPlan
	service.search = func(ctx context.Context, plan ExpertDatabaseRes.ExpertAgentPlan) ([]ExpertDatabaseRes.ExpertAgentCandidate, int64, error) {
		plans = append(plans, plan)
		return []ExpertDatabaseRes.ExpertAgentCandidate{
			{ID: 1, Name: "甲", UnitName: "甲大学"}, {ID: 2, Name: "乙", UnitName: "乙大学"},
		}, 2, nil
	}
	first, err := service.Chat(context.Background(), 11, ExpertDatabaseReq.ExpertAgentChatReq{Message: "找一位基层治理专家"})
	if err != nil {
		t.Fatalf("first chat: %v", err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	done := make(chan error, 1)
	go func() {
		_, chatErr := service.Chat(ctx, 11, ExpertDatabaseReq.ExpertAgentChatReq{SessionID: first.SessionID, Message: "取消前临时改成完全不同的主题"})
		done <- chatErr
	}()
	select {
	case <-started:
	case <-time.After(3 * time.Second):
		t.Fatal("provider did not receive the turn that will be canceled")
	}
	_, busyErr := service.Chat(context.Background(), 11, ExpertDatabaseReq.ExpertAgentChatReq{SessionID: first.SessionID, Message: "并发继续"})
	if !errors.Is(busyErr, ErrExpertAgentSessionBusy) {
		t.Fatalf("concurrent session turns should be rejected, got %v", busyErr)
	}
	cancel()
	select {
	case err := <-done:
		if !errors.Is(err, context.Canceled) {
			t.Fatalf("expected caller cancellation, got %v", err)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("canceling the context did not stop the provider call")
	}
	last, err := service.Chat(context.Background(), 11, ExpertDatabaseReq.ExpertAgentChatReq{SessionID: first.SessionID, Message: "仍按成功那次的条件换一位"})
	if err != nil {
		t.Fatalf("session must remain usable after cancellation: %v", err)
	}
	if calls.Load() != 3 || len(plans) != 2 || !reflect.DeepEqual(plans[1].Keywords, []string{"基层治理"}) || !reflect.DeepEqual(plans[1].ExcludedExpertIDs, []uint{1}) {
		t.Fatalf("canceled/concurrent turn changed session state: calls=%d plans=%+v", calls.Load(), plans)
	}
	if got := testExpertAgentCandidateIDs(last.Recommendations); !reflect.DeepEqual(got, []uint{2}) {
		t.Fatalf("last successful selection was not preserved after cancellation: %v", got)
	}
}

func TestExpertAgentChat_RejectsUntrustedModelOutputWithoutExecutingSearch(t *testing.T) {
	cases := []struct {
		name string
		body string
	}{
		{"prose instead of tool", `{"choices":[{"message":{"content":"{\"keywords\":[\"基层治理\"]}"}}]}`},
		{"unauthorized tool", `{"choices":[{"message":{"tool_calls":[{"type":"function","function":{"name":"delete_experts","arguments":"{}"}}]}}]}`},
		{"malformed arguments", `{"choices":[{"message":{"tool_calls":[{"type":"function","function":{"name":"search_experts","arguments":"private-provider-debug"}}]}}]}`},
		{"fabricated selected IDs", `{"choices":[{"message":{"tool_calls":[{"type":"function","function":{"name":"search_experts","arguments":"{\"selectedIds\":[999999]}"}}]}}]}`},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			service := newTestExpertAgent(t, func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				_, _ = io.WriteString(w, tc.body)
			})
			var searchCalls int
			service.search = func(context.Context, ExpertDatabaseRes.ExpertAgentPlan) ([]ExpertDatabaseRes.ExpertAgentCandidate, int64, error) {
				searchCalls++
				return nil, 0, nil
			}
			_, err := service.Chat(context.Background(), 11, ExpertDatabaseReq.ExpertAgentChatReq{Message: "找一位基层治理专家"})
			if !errors.Is(err, ErrExpertAgentUpstream) {
				t.Fatalf("untrusted model output should return a safe upstream error, got %v", err)
			}
			if searchCalls != 0 || strings.Contains(err.Error(), "private-provider-debug") || strings.Contains(err.Error(), "999999") {
				t.Fatalf("untrusted output reached search or leaked in the error: searchCalls=%d err=%v", searchCalls, err)
			}
		})
	}
}

func TestExpertAgentChat_ValidatesRequestBeforeCallingProvider(t *testing.T) {
	var calls atomic.Int32
	service := newTestExpertAgent(t, func(w http.ResponseWriter, r *http.Request) {
		calls.Add(1)
		writeTestExpertAgentPlan(t, w, `{}`)
	})
	cases := []struct {
		name       string
		operatorID uint
		request    ExpertDatabaseReq.ExpertAgentChatReq
	}{
		{"unauthenticated", 0, ExpertDatabaseReq.ExpertAgentChatReq{Message: "找专家"}},
		{"blank message", 11, ExpertDatabaseReq.ExpertAgentChatReq{Message: "  \n "}},
		{"oversized message", 11, ExpertDatabaseReq.ExpertAgentChatReq{Message: strings.Repeat("专", 2001)}},
		{"negative count", 11, ExpertDatabaseReq.ExpertAgentChatReq{Message: "找专家", Count: -1}},
		{"excessive count", 11, ExpertDatabaseReq.ExpertAgentChatReq{Message: "找专家", Count: 31}},
		{"negative unit limit", 11, ExpertDatabaseReq.ExpertAgentChatReq{Message: "找专家", MaxPerUnit: -1}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := service.Chat(context.Background(), tc.operatorID, tc.request)
			if !errors.Is(err, ErrExpertAgentInvalidRequest) {
				t.Fatalf("invalid request must be rejected, got %v", err)
			}
		})
	}
	if calls.Load() != 0 {
		t.Fatalf("invalid request reached provider, calls=%d", calls.Load())
	}
}

func TestExpertAgentChat_MissingKeyReportsNotConfiguredWithoutCallingProvider(t *testing.T) {
	var calls atomic.Int32
	service := newTestExpertAgent(t, func(w http.ResponseWriter, r *http.Request) {
		calls.Add(1)
		writeTestExpertAgentPlan(t, w, `{}`)
	})
	t.Setenv("EXPERT_AGENT_API_KEY", "")
	if service.Status().Configured {
		t.Fatal("missing API key must be reflected in assistant status")
	}
	_, err := service.Chat(context.Background(), 11, ExpertDatabaseReq.ExpertAgentChatReq{Message: "找一位专家"})
	if !errors.Is(err, ErrExpertAgentNotConfigured) || calls.Load() != 0 {
		t.Fatalf("unconfigured assistant must not call any provider: calls=%d err=%v", calls.Load(), err)
	}
}

func TestExpertAgentChat_GreetingAsksForCriteriaWithoutSearchingAllExperts(t *testing.T) {
	service := newTestExpertAgent(t, func(w http.ResponseWriter, r *http.Request) {
		writeTestExpertAgentPlan(t, w, `{}`)
	})
	var searches int
	service.search = func(context.Context, ExpertDatabaseRes.ExpertAgentPlan) ([]ExpertDatabaseRes.ExpertAgentCandidate, int64, error) {
		searches++
		return []ExpertDatabaseRes.ExpertAgentCandidate{{ID: 1, Name: "不应出现的专家"}}, 1, nil
	}
	result, err := service.Chat(context.Background(), 11, ExpertDatabaseReq.ExpertAgentChatReq{Message: "你好"})
	if err != nil {
		t.Fatalf("greeting should produce clarification, got %v", err)
	}
	if searches != 0 || len(result.Candidates) != 0 || len(result.Recommendations) != 0 {
		t.Fatalf("a greeting must not silently trigger unrestricted browsing: searches=%d result=%+v", searches, result)
	}
	if result.Reply == "" || (!strings.Contains(result.Reply, "主题") && !strings.Contains(result.Reply, "条件")) {
		t.Fatalf("clarification should ask for useful selection criteria, got %q", result.Reply)
	}
}

func TestExpertAgentChat_CancellationAfterSearchDoesNotCommitNewCriteria(t *testing.T) {
	var calls atomic.Int32
	service := newTestExpertAgent(t, func(w http.ResponseWriter, r *http.Request) {
		switch calls.Add(1) {
		case 1:
			writeTestExpertAgentPlan(t, w, `{"keywords":["基层治理"],"count":1}`)
		case 2:
			writeTestExpertAgentPlan(t, w, `{"keywords":["不应提交的新主题"]}`)
		default:
			writeTestExpertAgentPlan(t, w, `{"excludePrevious":true}`)
		}
	})
	var plans []ExpertDatabaseRes.ExpertAgentPlan
	var cancelSearch context.CancelFunc
	service.search = func(ctx context.Context, plan ExpertDatabaseRes.ExpertAgentPlan) ([]ExpertDatabaseRes.ExpertAgentCandidate, int64, error) {
		plans = append(plans, plan)
		if len(plans) == 2 {
			// Cancellation arriving after the database returns must still prevent
			// the new plan and selection from replacing the committed session.
			cancelSearch()
		}
		return []ExpertDatabaseRes.ExpertAgentCandidate{{ID: 1, UnitName: "甲大学"}, {ID: 2, UnitName: "乙大学"}}, 2, nil
	}
	first, err := service.Chat(context.Background(), 11, ExpertDatabaseReq.ExpertAgentChatReq{Message: "找一位基层治理专家"})
	if err != nil {
		t.Fatalf("first chat: %v", err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancelSearch = cancel
	defer cancel()
	_, err = service.Chat(ctx, 11, ExpertDatabaseReq.ExpertAgentChatReq{SessionID: first.SessionID, Message: "换一个研究主题"})
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("late cancellation must be honored, got %v", err)
	}
	last, err := service.Chat(context.Background(), 11, ExpertDatabaseReq.ExpertAgentChatReq{SessionID: first.SessionID, Message: "仍按成功条件换一位"})
	if err != nil {
		t.Fatalf("resume after late cancellation: %v", err)
	}
	if len(plans) != 3 || !reflect.DeepEqual(plans[1].Keywords, []string{"不应提交的新主题"}) || !reflect.DeepEqual(plans[2].Keywords, []string{"基层治理"}) {
		t.Fatalf("late cancellation overwrote the prior committed plan: %+v", plans)
	}
	if got := testExpertAgentCandidateIDs(last.Recommendations); !reflect.DeepEqual(got, []uint{2}) {
		t.Fatalf("late cancellation overwrote the prior selection: %v", got)
	}
}

func TestExpertAgentChat_ExplicitBrowseAllContinuesOnFollowUp(t *testing.T) {
	var calls atomic.Int32
	service := newTestExpertAgent(t, func(w http.ResponseWriter, r *http.Request) {
		if calls.Add(1) == 1 {
			writeTestExpertAgentPlan(t, w, `{"count":1}`)
			return
		}
		writeTestExpertAgentPlan(t, w, `{"excludePrevious":true}`)
	})
	var searches int
	service.search = func(context.Context, ExpertDatabaseRes.ExpertAgentPlan) ([]ExpertDatabaseRes.ExpertAgentCandidate, int64, error) {
		searches++
		return []ExpertDatabaseRes.ExpertAgentCandidate{{ID: 1, UnitName: "甲大学"}, {ID: 2, UnitName: "乙大学"}}, 2, nil
	}
	first, err := service.Chat(context.Background(), 11, ExpertDatabaseReq.ExpertAgentChatReq{Message: "浏览全部专家，推荐一位"})
	if err != nil || len(first.Recommendations) != 1 {
		t.Fatalf("explicit full-database browse should select one candidate: result=%+v err=%v", first, err)
	}
	second, err := service.Chat(context.Background(), 11, ExpertDatabaseReq.ExpertAgentChatReq{SessionID: first.SessionID, Message: "换一批"})
	if err != nil {
		t.Fatalf("follow-up of an explicit browse should continue: %v", err)
	}
	if searches != 2 || !reflect.DeepEqual(testExpertAgentCandidateIDs(second.Recommendations), []uint{2}) {
		t.Fatalf("follow-up forgot the committed full-database browse: searches=%d result=%+v", searches, second)
	}
}

func TestExpertAgentSessionStore_ExpiredSessionCannotBeReused(t *testing.T) {
	store := newExpertAgentSessionStore()
	id, session, _, err := store.acquire(11, "")
	if err != nil {
		t.Fatalf("create session: %v", err)
	}
	store.release(id, false)
	session.updated = time.Now().Add(-agentSessionTTL - time.Second)
	_, _, _, err = store.acquire(11, id)
	if !errors.Is(err, ErrExpertAgentSessionNotFound) {
		t.Fatalf("expired session must not be reused, got %v", err)
	}
	if _, exists := store.entries[id]; exists {
		t.Fatal("expired session should be removed from memory")
	}
}

func TestExpertAgentSessionStore_EnforcesQuotasAndReclaimsExpiredCapacity(t *testing.T) {
	cases := []struct {
		name  string
		count int
		owner func(int) uint
	}{
		{"per-user quota", agentMaxSessionsPerUser, func(int) uint { return 11 }},
		{"global quota", agentMaxSessions, func(i int) uint { return uint(i + 100) }},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			store := newExpertAgentSessionStore()
			for i := 0; i < tc.count; i++ {
				store.entries[idStr(uint(i+1))] = &expertAgentSession{owner: tc.owner(i), updated: time.Now()}
			}
			_, _, _, err := store.acquire(11, "")
			if !errors.Is(err, ErrExpertAgentSessionLimit) {
				t.Fatalf("session quota must reject additional allocations, got %v", err)
			}
			store.entries["1"].updated = time.Now().Add(-agentSessionTTL - time.Second)
			id, _, created, err := store.acquire(11, "")
			if err != nil || !created || id == "" || len(store.entries) != tc.count {
				t.Fatalf("expired entry must free capacity for one fresh session: id=%q created=%v count=%d err=%v", id, created, len(store.entries), err)
			}
			if _, exists := store.entries["1"]; exists {
				t.Fatal("expired entry was retained after capacity reclamation")
			}
		})
	}
}

func TestExpertAgentChat_DeniedBrowseDoesNotReadCandidates(t *testing.T) {
	for _, message := range []string{"不要全库浏览，先问清我需要什么", "别从全部专家里随机挑"} {
		t.Run(message, func(t *testing.T) {
			service := newTestExpertAgent(t, func(w http.ResponseWriter, r *http.Request) {
				writeTestExpertAgentPlan(t, w, `{}`)
			})
			searches := 0
			service.search = func(context.Context, ExpertDatabaseRes.ExpertAgentPlan) ([]ExpertDatabaseRes.ExpertAgentCandidate, int64, error) {
				searches++
				return []ExpertDatabaseRes.ExpertAgentCandidate{{ID: 1}}, 1, nil
			}
			result, err := service.Chat(context.Background(), 11, ExpertDatabaseReq.ExpertAgentChatReq{Message: message})
			if err != nil || searches != 0 || len(result.Candidates) != 0 || !strings.Contains(result.Reply, "请补充") {
				t.Fatalf("refusal should ask for criteria without querying candidates: searches=%d result=%+v err=%v", searches, result, err)
			}
		})
	}
}

func TestExpertAgentChat_DeniedBrowseClearsPreviousPermission(t *testing.T) {
	service := newTestExpertAgent(t, func(w http.ResponseWriter, r *http.Request) {
		writeTestExpertAgentPlan(t, w, `{}`)
	})
	searches := 0
	service.search = func(context.Context, ExpertDatabaseRes.ExpertAgentPlan) ([]ExpertDatabaseRes.ExpertAgentCandidate, int64, error) {
		searches++
		return []ExpertDatabaseRes.ExpertAgentCandidate{{ID: 1}}, 1, nil
	}
	first, err := service.Chat(context.Background(), 11, ExpertDatabaseReq.ExpertAgentChatReq{Message: "浏览全部专家"})
	if err != nil {
		t.Fatal(err)
	}
	second, err := service.Chat(context.Background(), 11, ExpertDatabaseReq.ExpertAgentChatReq{SessionID: first.SessionID, Message: "不要全库浏览"})
	if err != nil || searches != 1 || len(second.Recommendations) != 0 {
		t.Fatalf("refusal must revoke the prior broad-browse intent: searches=%d result=%+v err=%v", searches, second, err)
	}
}
