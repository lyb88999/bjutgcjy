package ExpertDatabase

import (
	"reflect"
	"testing"

	"github.com/flipped-aurora/gin-vue-admin/server/model/ExpertDatabase"
	ExpertDatabaseReq "github.com/flipped-aurora/gin-vue-admin/server/model/ExpertDatabase/request"
)

func TestSplitKeyword_SplitsOnCommonDelimiters(t *testing.T) {
	got := splitKeyword("防汛,基层治理、数字政府；乡村振兴")
	want := []string{"防汛", "基层治理", "数字政府", "乡村振兴"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("分词结果不对：想要 %v，实际 %v", want, got)
	}
}

// TestRankedCandidates_FiltersOutZeroRelevanceWhenKeywordGiven 复现并回归验证之前那个真实 bug：
// 搜"医生"的时候返回了全部 419 个已发布专家，因为职称/社会贡献分不受相关性影响，
// 完全不相关的人也能凭职称权重排到前面。修复后，填了关键词但一个词都不命中的候选人应该被整体排除。
func TestRankedCandidates_FiltersOutZeroRelevanceWhenKeywordGiven(t *testing.T) {
	db := setupTestDB(t)

	relevant := ExpertDatabase.ExpertProfile{Name: "相关专家", Status: StatusPublished, TechTitle: "讲师"}
	irrelevantButHighTitle := ExpertDatabase.ExpertProfile{Name: "不相关但职称高", Status: StatusPublished, TechTitle: "教授"}
	if err := db.Create(&relevant).Error; err != nil {
		t.Fatalf("建档案失败: %v", err)
	}
	if err := db.Create(&irrelevantButHighTitle).Error; err != nil {
		t.Fatalf("建档案失败: %v", err)
	}
	achievement := ExpertDatabase.ExpertAchievement{ExpertId: relevant.ID, Title: "关于医生职业发展的研究", AchievementType: "学术论文"}
	if err := db.Create(&achievement).Error; err != nil {
		t.Fatalf("建成果失败: %v", err)
	}

	svc := &ExpertSearchService{}
	items, err := svc.rankedCandidates(ExpertDatabaseReq.ExpertSearchReq{Keyword: "医生"})
	if err != nil {
		t.Fatalf("检索不应该报错: %v", err)
	}

	if len(items) != 1 {
		t.Fatalf("搜'医生'应该只返回 1 条真正相关的记录，实际返回 %d 条", len(items))
	}
	if items[0].Name != "相关专家" {
		t.Fatalf("返回的应该是语料里真正命中关键词的那位，实际是 %s", items[0].Name)
	}
}

func TestRankedCandidates_EmptyKeywordReturnsEveryone(t *testing.T) {
	db := setupTestDB(t)
	for i := 0; i < 3; i++ {
		p := ExpertDatabase.ExpertProfile{Name: "专家", Status: StatusPublished}
		if err := db.Create(&p).Error; err != nil {
			t.Fatalf("建档案失败: %v", err)
		}
	}

	svc := &ExpertSearchService{}
	items, err := svc.rankedCandidates(ExpertDatabaseReq.ExpertSearchReq{})
	if err != nil {
		t.Fatalf("检索不应该报错: %v", err)
	}
	if len(items) != 3 {
		t.Fatalf("不带关键词应该是默认浏览，返回全部已发布记录，想要 3 条，实际 %d 条", len(items))
	}
}

// 标签关联是软删除的：解除关联之后再搜这个标签词，不应该还能命中这位专家
func TestRankedCandidates_RemovedTagRelationExcludedFromCorpus(t *testing.T) {
	db := setupTestDB(t)

	expert := ExpertDatabase.ExpertProfile{Name: "曾有标签的专家", Status: StatusPublished}
	if err := db.Create(&expert).Error; err != nil {
		t.Fatalf("建档案失败: %v", err)
	}
	tag := ExpertDatabase.ExpertTag{TagType: "keyword", TagValue: "碳中和"}
	if err := db.Create(&tag).Error; err != nil {
		t.Fatalf("建标签失败: %v", err)
	}
	relation := ExpertDatabase.ExpertTagRelation{ExpertId: expert.ID, TagId: tag.ID}
	if err := db.Create(&relation).Error; err != nil {
		t.Fatalf("建关联失败: %v", err)
	}

	svc := &ExpertSearchService{}
	items, err := svc.rankedCandidates(ExpertDatabaseReq.ExpertSearchReq{Keyword: "碳中和"})
	if err != nil {
		t.Fatalf("检索不应该报错: %v", err)
	}
	if len(items) != 1 {
		t.Fatalf("标签关联存在时应命中 1 条，实际 %d 条", len(items))
	}

	// 解除关联（软删除），语料里不应再有这个标签
	if err := db.Delete(&relation).Error; err != nil {
		t.Fatalf("解除关联失败: %v", err)
	}
	items, err = svc.rankedCandidates(ExpertDatabaseReq.ExpertSearchReq{Keyword: "碳中和"})
	if err != nil {
		t.Fatalf("检索不应该报错: %v", err)
	}
	if len(items) != 0 {
		t.Fatalf("关联已解除，不应再命中，实际 %d 条", len(items))
	}
}

// TestRankedCandidates_MultiKeywordUnionMatchesAnyTermRanksMoreHitsHigher 复现并回归验证之前
// "2-3 个关键词无法联合搜索"的 bug：旧实现把整个输入串（比如"人工智能 京津冀"）当一个词去跟
// 语料整段比对，几乎不可能原样命中，等于多关键词检索直接搜不出东西。修复后应该拆成独立词分别
// 判断命中（命中任意一个词就收录），且命中词更多的候选人排名更靠前。
func TestRankedCandidates_MultiKeywordUnionMatchesAnyTermRanksMoreHitsHigher(t *testing.T) {
	db := setupTestDB(t)

	onlyOneTerm := ExpertDatabase.ExpertProfile{Name: "只命中一个词", Status: StatusPublished, ResearchKeywords: "人工智能"}
	bothTerms := ExpertDatabase.ExpertProfile{Name: "两个词都命中", Status: StatusPublished, ResearchKeywords: "人工智能 京津冀协同发展"}
	if err := db.Create(&onlyOneTerm).Error; err != nil {
		t.Fatalf("建档案失败: %v", err)
	}
	if err := db.Create(&bothTerms).Error; err != nil {
		t.Fatalf("建档案失败: %v", err)
	}

	svc := &ExpertSearchService{}
	items, err := svc.rankedCandidates(ExpertDatabaseReq.ExpertSearchReq{Keyword: "人工智能 京津冀"})
	if err != nil {
		t.Fatalf("检索不应该报错: %v", err)
	}
	if len(items) != 2 {
		t.Fatalf("命中任意一个词都该收录，想要 2 条，实际 %d 条", len(items))
	}
	if items[0].Name != "两个词都命中" {
		t.Fatalf("命中词更多的应该排第一，实际排第一的是 %s", items[0].Name)
	}
	if items[0].Relevance <= items[1].Relevance {
		t.Fatalf("命中两个词的相关性应该高于只命中一个词的，实际 %v vs %v", items[0].Relevance, items[1].Relevance)
	}
}

// TestRankedCandidates_KeywordMatchAlwaysMovesCompositeScoreEvenWithZeroBaseScores 复现并回归
// 验证之前那个真实 bug："综合得分"公式是 w1*成果分*相关性 + w2*决策影响分*相关性 + w3*职称权重 +
// w4*社会贡献分——成果分/决策影响分是 0（还没攒够成果）的专家，不管关键词匹配得多好，前两项乘出来
// 还是 0，综合得分对这批人形同一个跟检索词毫无关系的固定值。修复后要单独加一项只看相关性本身的
// 得分，让关键词匹配程度总能体现到综合得分里。
func TestRankedCandidates_KeywordMatchAlwaysMovesCompositeScoreEvenWithZeroBaseScores(t *testing.T) {
	db := setupTestDB(t)

	// 两人成果分、决策影响分、职称、社会贡献分完全一样（都是零基础），唯一的差别是关键词命中程度
	partialMatch := ExpertDatabase.ExpertProfile{Name: "只命中一个词", Status: StatusPublished, ResearchKeywords: "人工智能"}
	fullMatch := ExpertDatabase.ExpertProfile{Name: "两个词都命中", Status: StatusPublished, ResearchKeywords: "人工智能 京津冀协同发展"}
	if err := db.Create(&partialMatch).Error; err != nil {
		t.Fatalf("建档案失败: %v", err)
	}
	if err := db.Create(&fullMatch).Error; err != nil {
		t.Fatalf("建档案失败: %v", err)
	}

	svc := &ExpertSearchService{}
	items, err := svc.rankedCandidates(ExpertDatabaseReq.ExpertSearchReq{Keyword: "人工智能 京津冀"})
	if err != nil {
		t.Fatalf("检索不应该报错: %v", err)
	}
	if len(items) != 2 {
		t.Fatalf("想要 2 条，实际 %d 条", len(items))
	}
	if items[0].RealtimeScore <= items[1].RealtimeScore {
		t.Fatalf("成果分/决策影响分/职称/社会贡献分都相同时，命中更多关键词的综合得分应该更高，实际 %v vs %v",
			items[0].RealtimeScore, items[1].RealtimeScore)
	}
}
