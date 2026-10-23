package ExpertDatabase

import (
	"testing"

	"github.com/flipped-aurora/gin-vue-admin/server/model/ExpertDatabase"
	ExpertDatabaseReq "github.com/flipped-aurora/gin-vue-admin/server/model/ExpertDatabase/request"
)

// ---------- combineTermResults 纯单元测试，不需要数据库 ----------

func TestCombineTermResults_AndTakesWeakestLink(t *testing.T) {
	got := combineTermResults(termResult{matched: true, relevance: 0.9}, "and", termResult{matched: true, relevance: 0.6})
	if !got.matched || got.relevance != 0.6 {
		t.Fatalf("且应该取两边较小的相关性（较保守），想要 matched=true relevance=0.6，实际 %+v", got)
	}
}

func TestCombineTermResults_AndFailsWhenEitherSideUnmatched(t *testing.T) {
	got := combineTermResults(termResult{matched: true, relevance: 0.9}, "and", termResult{matched: false, relevance: 0})
	if got.matched || got.relevance != 0 {
		t.Fatalf("且：只要有一边没命中就该整体不命中，实际 %+v", got)
	}
}

func TestCombineTermResults_OrTakesStrongestMatch(t *testing.T) {
	got := combineTermResults(termResult{matched: true, relevance: 0.6}, "or", termResult{matched: true, relevance: 0.9})
	if !got.matched || got.relevance != 0.9 {
		t.Fatalf("或应该取两边较大的相关性，想要 matched=true relevance=0.9，实际 %+v", got)
	}
}

func TestCombineTermResults_OrMatchesWhenOnlyOneSideMatches(t *testing.T) {
	got := combineTermResults(termResult{matched: false, relevance: 0}, "or", termResult{matched: true, relevance: 0.7})
	if !got.matched || got.relevance != 0.7 {
		t.Fatalf("或：只要有一边命中就该整体命中，实际 %+v", got)
	}
}

func TestCombineTermResults_NotPassesThroughAccRelevanceWhenExclusionHolds(t *testing.T) {
	got := combineTermResults(termResult{matched: true, relevance: 0.8}, "not", termResult{matched: false, relevance: 0})
	if !got.matched || got.relevance != 0.8 {
		t.Fatalf("不含：排除词没命中时应该维持原来的命中状态和相关性不变，实际 %+v", got)
	}
}

func TestCombineTermResults_NotFailsWhenExcludedTermMatches(t *testing.T) {
	got := combineTermResults(termResult{matched: true, relevance: 0.8}, "not", termResult{matched: true, relevance: 0.5})
	if got.matched {
		t.Fatalf("不含：排除词命中了就该让整体判定失败，实际 %+v", got)
	}
}

// ---------- parseAdvancedRows ----------

func TestParseAdvancedRows_FiltersBlankTermRows(t *testing.T) {
	rows, err := parseAdvancedRows(`[{"term":"人工智能"},{"connector":"and","term":"  "},{"connector":"or","term":"京津冀"}]`)
	if err != nil {
		t.Fatalf("解析不应该报错: %v", err)
	}
	if len(rows) != 2 {
		t.Fatalf("空白词的行应该被过滤掉，想要 2 行，实际 %d 行: %+v", len(rows), rows)
	}
}

func TestParseAdvancedRows_InvalidJSONReturnsError(t *testing.T) {
	if _, err := parseAdvancedRows("not-json"); err == nil {
		t.Fatal("格式不对的 JSON 应该返回 error")
	}
}

// ---------- rankedCandidates 集成测试 ----------

func advancedRowsJSON(t *testing.T, rows []ExpertDatabaseReq.AdvancedSearchRow) string {
	t.Helper()
	var sb string
	sb = "["
	for i, r := range rows {
		if i > 0 {
			sb += ","
		}
		sb += `{"connector":"` + r.Connector + `","term":"` + r.Term + `"}`
	}
	sb += "]"
	return sb
}

func TestRankedCandidates_AdvancedAndExcludesCandidateMatchingOnlyOneTerm(t *testing.T) {
	db := setupTestDB(t)
	both := ExpertDatabase.ExpertProfile{Name: "两个词都命中", Status: StatusPublished, ResearchKeywords: "人工智能 京津冀协同发展"}
	onlyOne := ExpertDatabase.ExpertProfile{Name: "只命中一个词", Status: StatusPublished, ResearchKeywords: "人工智能"}
	if err := db.Create(&both).Error; err != nil {
		t.Fatalf("建档案失败: %v", err)
	}
	if err := db.Create(&onlyOne).Error; err != nil {
		t.Fatalf("建档案失败: %v", err)
	}

	rows := []ExpertDatabaseReq.AdvancedSearchRow{{Term: "人工智能"}, {Connector: "and", Term: "京津冀"}}
	svc := &ExpertSearchService{}
	items, err := svc.rankedCandidates(ExpertDatabaseReq.ExpertSearchReq{AdvancedRows: advancedRowsJSON(t, rows)})
	if err != nil {
		t.Fatalf("检索不应该报错: %v", err)
	}
	if len(items) != 1 || items[0].Name != "两个词都命中" {
		t.Fatalf("且：只有两个词都命中的应该留下，实际 %v", items)
	}
}

func TestRankedCandidates_AdvancedOrIncludesEitherMatch(t *testing.T) {
	db := setupTestDB(t)
	matchA := ExpertDatabase.ExpertProfile{Name: "命中A", Status: StatusPublished, ResearchKeywords: "人工智能"}
	matchB := ExpertDatabase.ExpertProfile{Name: "命中B", Status: StatusPublished, ResearchKeywords: "京津冀协同发展"}
	if err := db.Create(&matchA).Error; err != nil {
		t.Fatalf("建档案失败: %v", err)
	}
	if err := db.Create(&matchB).Error; err != nil {
		t.Fatalf("建档案失败: %v", err)
	}

	rows := []ExpertDatabaseReq.AdvancedSearchRow{{Term: "人工智能"}, {Connector: "or", Term: "京津冀"}}
	svc := &ExpertSearchService{}
	items, err := svc.rankedCandidates(ExpertDatabaseReq.ExpertSearchReq{AdvancedRows: advancedRowsJSON(t, rows)})
	if err != nil {
		t.Fatalf("检索不应该报错: %v", err)
	}
	if len(items) != 2 {
		t.Fatalf("或：两边任意命中都该收录，想要 2 条，实际 %d 条: %v", len(items), items)
	}
}

func TestRankedCandidates_AdvancedNotExcludesCandidateMatchingExcludedTerm(t *testing.T) {
	db := setupTestDB(t)
	clean := ExpertDatabase.ExpertProfile{Name: "干净命中", Status: StatusPublished, ResearchKeywords: "人工智能"}
	tainted := ExpertDatabase.ExpertProfile{Name: "命中但含排除词", Status: StatusPublished, ResearchKeywords: "人工智能 养老服务"}
	if err := db.Create(&clean).Error; err != nil {
		t.Fatalf("建档案失败: %v", err)
	}
	if err := db.Create(&tainted).Error; err != nil {
		t.Fatalf("建档案失败: %v", err)
	}

	rows := []ExpertDatabaseReq.AdvancedSearchRow{{Term: "人工智能"}, {Connector: "not", Term: "养老"}}
	svc := &ExpertSearchService{}
	items, err := svc.rankedCandidates(ExpertDatabaseReq.ExpertSearchReq{AdvancedRows: advancedRowsJSON(t, rows)})
	if err != nil {
		t.Fatalf("检索不应该报错: %v", err)
	}
	if len(items) != 1 || items[0].Name != "干净命中" {
		t.Fatalf("不含：提到排除词的候选人应该被排除掉，实际 %v", items)
	}
}

// TestRankedCandidates_AdvancedMixedChainEvaluatesStrictlyLeftToRight 专门验证"没有运算符优先级"
// 这条规则：A 或 B 且 C 严格按从左到右求值成 (A 或 B) 且 C，不是"且比或优先级高"那种算法（A 或
// (B 且 C)）。构造一个两种算法会给出不同答案的候选人——A 命中、B/C 都不命中：
//   - 从左到右（本实现）：(A或B)=命中，(A或B)且C = 命中 且 不命中 = 不命中 → 应该被排除
//   - 若误把"且"实现成优先级更高：B且C=不命中，A或(B且C) = 命中 或 不命中 = 命中 → 会被误收录
//
// 如果这个测试失败，八成是combineTermResults被改成了非严格从左到右的求值顺序
func TestRankedCandidates_AdvancedMixedChainEvaluatesStrictlyLeftToRight(t *testing.T) {
	db := setupTestDB(t)
	onlyA := ExpertDatabase.ExpertProfile{Name: "只命中A", Status: StatusPublished, ResearchKeywords: "人工智能"}
	if err := db.Create(&onlyA).Error; err != nil {
		t.Fatalf("建档案失败: %v", err)
	}

	rows := []ExpertDatabaseReq.AdvancedSearchRow{
		{Term: "人工智能"},
		{Connector: "or", Term: "京津冀协同"},
		{Connector: "and", Term: "养老服务"},
	}
	svc := &ExpertSearchService{}
	items, err := svc.rankedCandidates(ExpertDatabaseReq.ExpertSearchReq{AdvancedRows: advancedRowsJSON(t, rows)})
	if err != nil {
		t.Fatalf("检索不应该报错: %v", err)
	}
	if len(items) != 0 {
		t.Fatalf("严格从左到右求值 (A或B)且C 应该排除掉只命中 A 的候选人，实际返回了 %v", items)
	}
}

func TestRankedCandidates_AdvancedRowsTakesPrecedenceOverPlainKeyword(t *testing.T) {
	db := setupTestDB(t)
	profile := ExpertDatabase.ExpertProfile{Name: "专家", Status: StatusPublished, ResearchKeywords: "医生"}
	if err := db.Create(&profile).Error; err != nil {
		t.Fatalf("建档案失败: %v", err)
	}

	rows := []ExpertDatabaseReq.AdvancedSearchRow{{Term: "不存在的关键词"}}
	svc := &ExpertSearchService{}
	items, err := svc.rankedCandidates(ExpertDatabaseReq.ExpertSearchReq{
		Keyword:      "医生", // 这个词其实能命中，但 AdvancedRows 非空时应该被完全忽略
		AdvancedRows: advancedRowsJSON(t, rows),
	})
	if err != nil {
		t.Fatalf("检索不应该报错: %v", err)
	}
	if len(items) != 0 {
		t.Fatalf("AdvancedRows 非空时应该完全替代 Keyword，Keyword 不该生效，实际返回了 %v", items)
	}
}

func TestRankedCandidates_AdvancedRowsInvalidJSONReturnsError(t *testing.T) {
	setupTestDB(t)
	svc := &ExpertSearchService{}
	if _, err := svc.rankedCandidates(ExpertDatabaseReq.ExpertSearchReq{AdvancedRows: "not-json"}); err == nil {
		t.Fatal("AdvancedRows 格式不对时应该报错")
	}
}

// TestSearch_AdvancedRowsWithEmptyKeywordStillGoesThroughRankedCandidates 这是最容易改错的一行
// 路由条件：Keyword 为空但 AdvancedRows 非空时，不能被误判成"没填关键词的浏览模式"而跳过关键词
// 匹配——用两个候选人（一个命中高级检索条件、一个不命中）区分：如果误路由到 browseModeQuery，
// 两个都会被返回（浏览模式不做关键词过滤）；正确路由到 rankedCandidates 应该只返回命中的那个
func TestSearch_AdvancedRowsWithEmptyKeywordStillGoesThroughRankedCandidates(t *testing.T) {
	db := setupTestDB(t)
	match := ExpertDatabase.ExpertProfile{Name: "命中", Status: StatusPublished, ResearchKeywords: "人工智能"}
	noMatch := ExpertDatabase.ExpertProfile{Name: "不命中", Status: StatusPublished, ResearchKeywords: "无关内容"}
	if err := db.Create(&match).Error; err != nil {
		t.Fatalf("建档案失败: %v", err)
	}
	if err := db.Create(&noMatch).Error; err != nil {
		t.Fatalf("建档案失败: %v", err)
	}

	rows := []ExpertDatabaseReq.AdvancedSearchRow{{Term: "人工智能"}}
	svc := &ExpertSearchService{}
	list, total, err := svc.Search(ExpertDatabaseReq.ExpertSearchReq{
		AdvancedRows: advancedRowsJSON(t, rows),
		Page:         1, PageSize: 10,
	})
	if err != nil {
		t.Fatalf("检索不应该报错: %v", err)
	}
	if total != 1 || len(list) != 1 || list[0].Name != "命中" {
		t.Fatalf("Keyword 为空但 AdvancedRows 非空时应该走 rankedCandidates 做关键词过滤，"+
			"想要只有 1 条命中结果，实际 total=%d list=%v（如果这里返回了 2 条，说明被误路由到了 browseModeQuery）", total, list)
	}
	if list[0].MatchReason == "" {
		t.Fatal("走 rankedCandidates 的结果应该带 MatchReason；browseModeQuery 的结果 MatchReason 永远是空的，这里是空的说明路由错了")
	}
}
