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

func TestRelevance_EmptyKeywordAlwaysReturnsOne(t *testing.T) {
	svc := &ExpertSearchService{}
	// 没输入关键词就是"浏览全部"，不应该按相关性筛掉任何人
	if got := svc.relevance("随便什么语料", ""); got != 1 {
		t.Fatalf("空关键词应该返回 1，实际 %v", got)
	}
}

func TestRelevance_PartialMatchIsFractional(t *testing.T) {
	svc := &ExpertSearchService{}
	got := svc.relevance("研究防汛工程", "防汛,教育")
	if got != 0.5 {
		t.Fatalf("两个词命中一个，相关性应该是 0.5，实际 %v", got)
	}
}

func TestRelevance_NoMatchReturnsZero(t *testing.T) {
	svc := &ExpertSearchService{}
	got := svc.relevance("跟关键词完全不沾边的语料", "医生")
	if got != 0 {
		t.Fatalf("完全不命中应该是 0，实际 %v", got)
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
