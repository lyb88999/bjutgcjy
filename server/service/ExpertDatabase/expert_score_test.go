package ExpertDatabase

import (
	"testing"

	"github.com/flipped-aurora/gin-vue-admin/server/model/ExpertDatabase"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
)

func TestCalcLeveledScore_SumsWeightTimesCount(t *testing.T) {
	weights := map[string]float64{"国家级": 10, "省部级": 5}
	levels := []string{"国家级", "国家级", "省部级"}

	got := calcLeveledScore(levels, weights)
	want := 10*2 + 5*1.0
	if got != want {
		t.Fatalf("成果得分算错了：想要 %v，实际 %v", want, got)
	}
}

func TestCalcLeveledScore_UnknownLevelFallsBackToWeightOne(t *testing.T) {
	weights := map[string]float64{"国家级": 10}
	// 字典里没配置的级别，不能让这条成果完全不计分，应该给保底权重 1
	levels := []string{"国家级", "没配置过的级别"}

	got := calcLeveledScore(levels, weights)
	want := 10.0 + 1.0
	if got != want {
		t.Fatalf("未配置级别应该按权重 1 保底计分：想要 %v，实际 %v", want, got)
	}
}

func TestDictWeights_FallsBackWhenDictionaryMissing(t *testing.T) {
	setupTestDB(t)
	fallback := map[string]float64{"国家级": 10, "一般级": 1}
	svc := &ExpertScoreService{}

	got := svc.dictWeights(DictTypeAchievementLevel, fallback)
	if got["国家级"] != 10 || got["一般级"] != 1 {
		t.Fatalf("字典不存在时应该原样返回兜底权重，实际: %+v", got)
	}
}

func TestDictWeights_PrefersDictionaryOverFallback(t *testing.T) {
	db := setupTestDB(t)
	dict := system.SysDictionary{Name: "专家库-成果级别权重", Type: DictTypeAchievementLevel}
	if err := db.Create(&dict).Error; err != nil {
		t.Fatalf("建字典失败: %v", err)
	}
	// extend 存的是实际权重数字，dictWeights 优先解析 extend，取不到才退回 value
	detail := system.SysDictionaryDetail{Label: "国家级", Value: 999, Extend: "20", SysDictionaryID: int(dict.ID)}
	if err := db.Create(&detail).Error; err != nil {
		t.Fatalf("建字典明细失败: %v", err)
	}

	svc := &ExpertScoreService{}
	got := svc.dictWeights(DictTypeAchievementLevel, map[string]float64{"国家级": 10})
	if got["国家级"] != 20 {
		t.Fatalf("字典存在时应该用字典里配置的权重（20），实际用了 %v", got["国家级"])
	}
}

func TestRecomputeExpertScore_WritesCompositeScoreAndTimestamp(t *testing.T) {
	db := setupTestDB(t)
	profile := createTestProfile(t, db, StatusPublished, 0)
	profile.TechTitle = "教授"
	if err := db.Save(&profile).Error; err != nil {
		t.Fatalf("更新档案失败: %v", err)
	}

	svc := &ExpertScoreService{}
	if err := svc.RecomputeExpertScore(profile.ID); err != nil {
		t.Fatalf("重算得分不应该报错: %v", err)
	}

	var updated ExpertDatabase.ExpertProfile
	db.First(&updated, profile.ID)
	if updated.ScoreUpdatedAt == nil {
		t.Fatal("重算之后 score_updated_at 不应该是空的")
	}
	// 默认字典缺失时走兜底权重：教授职称权重是 5（见 defaultTitleLevelWeights）
	if updated.CompositeScore != 5 {
		t.Fatalf("没有成果/决策影响/学术兼职时，综合得分应该只有职称分 5，实际是 %v", updated.CompositeScore)
	}
}
