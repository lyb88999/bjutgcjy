package ExpertDatabase

import (
	"reflect"
	"strings"
	"testing"

	"github.com/flipped-aurora/gin-vue-admin/server/model/ExpertDatabase"
	ExpertDatabaseRes "github.com/flipped-aurora/gin-vue-admin/server/model/ExpertDatabase/response"
)

// 这些用例只使用内存记录和固定向量，不连接数据库、模型服务或网络。
func evidenceTestAchievement(id, expertID uint, title string) ExpertDatabase.ExpertAchievement {
	achievement := ExpertDatabase.ExpertAchievement{ExpertId: expertID, Title: title}
	achievement.ID = id
	return achievement
}

func evidenceTestProfile() ExpertDatabase.ExpertProfile {
	profile := ExpertDatabase.ExpertProfile{Name: "样例专家", ResearchDirections: "公共管理", ProfileLink: "https://example.org/expert"}
	profile.ID = 7
	return profile
}

func achievementEvidenceIDs(evidence []ExpertDatabaseRes.ExpertAgentEvidence) []uint {
	ids := []uint{}
	for _, entry := range evidence {
		if entry.Kind == "achievement" {
			ids = append(ids, entry.SourceID)
		}
	}
	return ids
}

func TestExpertAgentSemanticEvidence_IdenticalTruncatedTitlesResolveOnlyMatchedID(t *testing.T) {
	profile := evidenceTestProfile()
	prefix := strings.Repeat("相同标题", 8)
	first := evidenceTestAchievement(101, profile.ID, prefix+"：老龄社会研究")
	second := evidenceTestAchievement(102, profile.ID, prefix+"：财政制度研究")
	label := "成果：《" + truncateLabel(first.Title, 30) + "》"
	if label != "成果：《"+truncateLabel(second.Title, 30)+"》" {
		t.Fatal("fixture must reproduce a collision in the old truncated-title lookup")
	}
	term := keywordTerm{term: "晚年照护", expanded: []string{"晚年照护"}, vectors: [][]float32{{1, 0}}}
	items := []parsedEmbeddingItem{
		{ItemType: embeddingItemTypeAchievement, ItemRefID: second.ID, ItemLabel: label, Vector: []float32{0, 1}},
		{ItemType: embeddingItemTypeAchievement, ItemRefID: first.ID, ItemLabel: label, Vector: []float32{1, 0}},
	}
	relevance, matched, reason, source := (&ExpertSearchService{}).termMatchWithSource("不含查询原词的语料", items, term, true)
	if !matched || source == nil || source.ItemRefID != first.ID || reason != label {
		t.Fatalf("semantic match did not retain exact source: matched=%v source=%+v reason=%q", matched, source, reason)
	}
	evidence := buildExpertAgentEvidence(profile, []ExpertDatabase.ExpertAchievement{second, first}, nil, []keywordTerm{term}, []expertAgentSemanticMatch{{Source: *source, Relevance: relevance}})
	if len(evidence) == 0 || evidence[0].SourceID != first.ID || evidence[0].Title != first.Title {
		t.Fatalf("best evidence should be the actual matched achievement: %+v", evidence)
	}
	if got := achievementEvidenceIDs(evidence); !reflect.DeepEqual(got, []uint{first.ID}) {
		t.Fatalf("same prefix must not make the other achievement evidence: %v", got)
	}
}

func TestExpertAgentSemanticEvidence_BestSourceSurvivesFourItemLimit(t *testing.T) {
	profile := evidenceTestProfile()
	best := evidenceTestAchievement(1, profile.ID, "老年群体生活支持比较")
	achievements := []ExpertDatabase.ExpertAchievement{
		evidenceTestAchievement(104, profile.ID, "教育政策研究甲"),
		evidenceTestAchievement(103, profile.ID, "教育政策研究乙"),
		evidenceTestAchievement(102, profile.ID, "教育政策研究丙"),
		evidenceTestAchievement(101, profile.ID, "教育政策研究丁"),
		best,
	}
	terms := []keywordTerm{{term: "晚年照护", expanded: []string{"晚年照护"}}, {term: "教育", expanded: []string{"教育"}}}
	evidence := buildExpertAgentEvidence(profile, achievements, nil, terms, []expertAgentSemanticMatch{{
		Source: parsedEmbeddingItem{ItemType: embeddingItemTypeAchievement, ItemRefID: best.ID, ItemLabel: "成果：《老年群体生活支持比较》"}, Relevance: 0.95,
	}})
	if len(evidence) != 4 || evidence[0].Kind != "achievement" || evidence[0].SourceID != best.ID {
		t.Fatalf("newer literal records or profile context displaced the best source: %+v", evidence)
	}
}

func TestExpertAgentSemanticEvidence_MultipleTermsRetainTheirSourcesInRelevanceOrder(t *testing.T) {
	profile := evidenceTestProfile()
	first := evidenceTestAchievement(201, profile.ID, "学龄群体资源配置研究")
	second := evidenceTestAchievement(202, profile.ID, "高龄居民照料机制研究")
	items := []parsedEmbeddingItem{
		{ItemType: embeddingItemTypeAchievement, ItemRefID: first.ID, ItemLabel: "成果：《学龄群体资源配置研究》", Vector: []float32{0.8, 0.6}},
		{ItemType: embeddingItemTypeAchievement, ItemRefID: second.ID, ItemLabel: "成果：《高龄居民照料机制研究》", Vector: []float32{0, 1}},
	}
	terms := []keywordTerm{
		{term: "教育公平", expanded: []string{"教育公平"}, vectors: [][]float32{{1, 0}}},
		{term: "晚年照护", expanded: []string{"晚年照护"}, vectors: [][]float32{{0, 1}}},
		{term: "养老服务", expanded: []string{"养老服务"}, vectors: [][]float32{{0, 1}}},
	}
	matches := []expertAgentSemanticMatch{}
	for i, term := range terms {
		relevance, matched, _, source := (&ExpertSearchService{}).termMatchWithSource("抽象语料", items, term, true)
		if !matched || source == nil {
			t.Fatalf("term %d lost its semantic source", i)
		}
		wantID := first.ID
		if i > 0 {
			wantID = second.ID
		}
		if source.ItemRefID != wantID {
			t.Fatalf("term %d was associated with source %d instead of %d", i, source.ItemRefID, wantID)
		}
		matches = append(matches, expertAgentSemanticMatch{Source: *source, Relevance: relevance})
	}
	evidence := buildExpertAgentEvidence(profile, []ExpertDatabase.ExpertAchievement{first, second}, nil, terms, matches)
	if got := achievementEvidenceIDs(evidence); !reflect.DeepEqual(got, []uint{second.ID, first.ID}) {
		t.Fatalf("expected highest source first and duplicate terms deduplicated, got %v", got)
	}
	if matches[0].Source.ItemRefID != first.ID {
		t.Fatal("evidence ordering mutated the input term/source associations")
	}
}

func TestExpertAgentSemanticEvidence_RejectsMissingAndOtherExpertSourceIDs(t *testing.T) {
	profile := evidenceTestProfile()
	other := evidenceTestAchievement(301, profile.ID+1, "同名成果")
	local := evidenceTestAchievement(302, profile.ID, "同名成果")
	matches := []expertAgentSemanticMatch{
		{Source: parsedEmbeddingItem{ItemType: embeddingItemTypeAchievement, ItemRefID: other.ID, ItemLabel: "成果：《同名成果》"}, Relevance: 1},
		{Source: parsedEmbeddingItem{ItemType: embeddingItemTypeAchievement, ItemRefID: 999, ItemLabel: "成果：《同名成果》"}, Relevance: 0.9},
	}
	evidence := buildExpertAgentEvidence(profile, []ExpertDatabase.ExpertAchievement{other, local}, nil, nil, matches)
	if ids := achievementEvidenceIDs(evidence); len(ids) != 0 {
		t.Fatalf("missing or another expert's source must not resolve by title: %v", ids)
	}
}

func TestExpertAgentSemanticEvidence_ProfileAndTagUseCurrentSourceData(t *testing.T) {
	profile := evidenceTestProfile()
	profile.ResearchKeywords = "基层治理"
	evidence := buildExpertAgentEvidence(profile, nil, []string{"社区服务", "社会政策"}, nil, []expertAgentSemanticMatch{
		{Source: parsedEmbeddingItem{ItemType: embeddingItemTypeProfile, ItemLabel: "研究方向：旧标签"}, Relevance: 0.9},
		{Source: parsedEmbeddingItem{ItemType: embeddingItemTypeTag, ItemLabel: "标签：旧标签"}, Relevance: 0.8},
	})
	if len(evidence) < 2 || evidence[0].Kind != "profile" || evidence[0].SourceID != profile.ID || !strings.Contains(evidence[0].Title, "基层治理") {
		t.Fatalf("profile evidence must retain the profile identity and actual current fields: %+v", evidence)
	}
	if evidence[1].Kind != "tag" || evidence[1].SourceID != 0 || evidence[1].Title != "社区服务、社会政策" {
		t.Fatalf("aggregate tag evidence must not invent an individual tag ID: %+v", evidence[1])
	}
}

func TestTermMatchWithSource_LiteralSynonymFallbackAndThreshold(t *testing.T) {
	service := &ExpertSearchService{}
	term := keywordTerm{term: "城市更新", expanded: []string{"城市更新", "旧城改造"}, vectors: [][]float32{{1, 0}}}
	for _, test := range []struct {
		name, corpus, reason string
		semantic             bool
		vector               []float32
		matched              bool
		relevance            float64
	}{
		{"literal-original", "城市更新与旧城改造", "包含关键词「城市更新」", true, []float32{1, 0}, true, 1},
		{"synonym-only", "旧城改造研究", "命中关键词「城市更新」的同义词「旧城改造」", true, []float32{0, 1}, true, literalHitRelevanceFloor},
		{"synonym-fallback", "旧城改造研究", "命中关键词「城市更新」的同义词「旧城改造」", false, []float32{1, 0}, true, 1},
		{"no-semantic-fallback", "其它研究主题", "", false, []float32{1, 0}, false, 0},
		{"below-threshold", "其它研究主题", "成果：《弱相关条目》", true, []float32{0.5, 0.8660254}, false, 0},
	} {
		t.Run(test.name, func(t *testing.T) {
			items := []parsedEmbeddingItem{{ItemType: embeddingItemTypeAchievement, ItemRefID: 7, ItemLabel: "成果：《弱相关条目》", Vector: test.vector}}
			relevance, matched, reason, source := service.termMatchWithSource(test.corpus, items, term, test.semantic)
			if matched != test.matched || relevance != test.relevance || reason != test.reason || source != nil {
				t.Fatalf("unexpected literal/fallback result: relevance=%v matched=%v reason=%q source=%+v", relevance, matched, reason, source)
			}
			legacyRelevance, legacyMatched, legacyReason := service.termMatch(test.corpus, items, term, test.semantic)
			if legacyRelevance != relevance || legacyMatched != matched || legacyReason != reason {
				t.Fatal("legacy termMatch wrapper changed search behavior")
			}
		})
	}
}

func TestTermMatchWithSource_ReturnsIndependentMetadataCopy(t *testing.T) {
	items := []parsedEmbeddingItem{{ItemType: embeddingItemTypeAchievement, ItemRefID: 1, ItemLabel: "最佳成果", Vector: []float32{1, 0}}}
	term := keywordTerm{term: "未出现的原词", expanded: []string{"未出现的原词"}, vectors: [][]float32{{1, 0}}}
	_, matched, _, source := (&ExpertSearchService{}).termMatchWithSource("其它研究主题", items, term, true)
	if !matched || source == nil || source.Vector != nil {
		t.Fatalf("expected independent metadata without vector alias: %+v", source)
	}
	source.ItemRefID = 99
	if items[0].ItemRefID != 1 {
		t.Fatal("changing returned metadata mutated the shared cache")
	}
	items[0].ItemLabel = "后来更新的缓存标签"
	if source.ItemLabel != "最佳成果" {
		t.Fatal("returned metadata aliases mutable cached metadata")
	}
}
