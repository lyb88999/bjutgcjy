package ExpertDatabase

import (
	"testing"

	"github.com/flipped-aurora/gin-vue-admin/server/model/ExpertDatabase"
	ExpertDatabaseReq "github.com/flipped-aurora/gin-vue-admin/server/model/ExpertDatabase/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
)

// TestGetFacetCounts_OrgFacetCountsScopedByOtherActiveFilters 分面计数的核心产品要求：选中一个
// 维度的筛选后，其它维度的计数只统计"在这个筛选条件下"的分布，不是全库固定数字——这里筛了
// 一级学科=教育学，单位分面里每个单位的数字应该只算这个单位下学科是教育学的人数，不含法学的
func TestGetFacetCounts_OrgFacetCountsScopedByOtherActiveFilters(t *testing.T) {
	db := setupTestDB(t)

	orgA := system.SysOrganization{Name: "A大学"}
	orgB := system.SysOrganization{Name: "B大学"}
	if err := db.Create(&orgA).Error; err != nil {
		t.Fatalf("建单位失败: %v", err)
	}
	if err := db.Create(&orgB).Error; err != nil {
		t.Fatalf("建单位失败: %v", err)
	}

	create := func(orgID uint, discipline string) {
		p := ExpertDatabase.ExpertProfile{Name: "专家", Status: StatusPublished, OrgId: &orgID, DisciplineL1: discipline}
		if err := db.Create(&p).Error; err != nil {
			t.Fatalf("建档案失败: %v", err)
		}
	}
	create(orgA.ID, "教育学")
	create(orgA.ID, "法学")
	create(orgB.ID, "教育学")
	create(orgB.ID, "教育学")

	svc := &ExpertSearchService{}
	result, err := svc.GetFacetCounts(ExpertDatabaseReq.ExpertSearchReq{DisciplineL1: "教育学"})
	if err != nil {
		t.Fatalf("查询分面计数不应该报错: %v", err)
	}

	counts := map[string]int64{}
	for _, o := range result.Org {
		counts[o.Label] = o.Count
	}
	if counts["A大学"] != 1 {
		t.Fatalf("筛了教育学之后，A大学在单位分面里应该是 1（原本 2 条，1 条是法学被筛掉了），实际 %d", counts["A大学"])
	}
	if counts["B大学"] != 2 {
		t.Fatalf("筛了教育学之后，B大学在单位分面里应该是 2，实际 %d", counts["B大学"])
	}
	if result.Total != 3 {
		t.Fatalf("筛了教育学之后总数应该是 3（A大学1条+B大学2条），实际 %d", result.Total)
	}
}

// TestGetFacetCounts_DoesNotApplyKeyword 锁定一个有意的产品口径：分面计数不按关键词收窄，带不带
// 关键词算出来的计数应该一样——这不是遗漏，是因为按关键词收窄要对每个候选人跑一次语义匹配，
// 四个维度各跑一次的开销扛不住，见 GetFacetCounts 顶部注释
func TestGetFacetCounts_DoesNotApplyKeyword(t *testing.T) {
	db := setupTestDB(t)
	if err := db.Create(&ExpertDatabase.ExpertProfile{Name: "专家一", Status: StatusPublished, ResearchKeywords: "人工智能"}).Error; err != nil {
		t.Fatalf("建档案失败: %v", err)
	}
	if err := db.Create(&ExpertDatabase.ExpertProfile{Name: "专家二", Status: StatusPublished, ResearchKeywords: "无关主题"}).Error; err != nil {
		t.Fatalf("建档案失败: %v", err)
	}

	svc := &ExpertSearchService{}
	withoutKeyword, err := svc.GetFacetCounts(ExpertDatabaseReq.ExpertSearchReq{})
	if err != nil {
		t.Fatalf("查询分面计数不应该报错: %v", err)
	}
	withKeyword, err := svc.GetFacetCounts(ExpertDatabaseReq.ExpertSearchReq{Keyword: "人工智能"})
	if err != nil {
		t.Fatalf("查询分面计数不应该报错: %v", err)
	}

	if withoutKeyword.Total != withKeyword.Total {
		t.Fatalf("分面计数不应该受关键词影响，不带关键词总数 %d，带关键词总数 %d", withoutKeyword.Total, withKeyword.Total)
	}
	if withoutKeyword.Total != 2 {
		t.Fatalf("想要总数 2，实际 %d", withoutKeyword.Total)
	}
}
