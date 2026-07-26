package ExpertDatabase

import (
	"testing"

	"github.com/flipped-aurora/gin-vue-admin/server/model/ExpertDatabase"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
)

func TestCreateExpertProfile_ReviewRequiredKeepsDraftStatus(t *testing.T) {
	setupTestDB(t)
	svc := &ExpertProfileService{}
	profile := ExpertDatabase.ExpertProfile{Name: "测试专家"}

	if err := svc.CreateExpertProfile(&profile); err != nil {
		t.Fatalf("创建不应该报错: %v", err)
	}
	if profile.ReviewBypassed {
		t.Fatal("审核开关开启时不应该带免审核发布标记")
	}
	// 字典缺失时 reviewRequired() 默认 true，正常走草稿流程；这里只验证没有被强改成 published
	if profile.Status == StatusPublished {
		t.Fatal("审核开关开启时新建档案不应该被直接发布")
	}
}

func TestCreateExpertProfile_ReviewSwitchOffPublishesDirectlyAndRecomputesScore(t *testing.T) {
	db := setupTestDB(t)
	dict := system.SysDictionary{Name: "专家库-审核开关", Type: DictTypeReviewSwitch}
	if err := db.Create(&dict).Error; err != nil {
		t.Fatalf("建字典失败: %v", err)
	}
	if err := db.Create(&system.SysDictionaryDetail{Label: "启用审核", Value: 0, SysDictionaryID: int(dict.ID)}).Error; err != nil {
		t.Fatalf("建字典明细失败: %v", err)
	}

	svc := &ExpertProfileService{}
	profile := ExpertDatabase.ExpertProfile{Name: "测试专家", TechTitle: "教授"}
	if err := svc.CreateExpertProfile(&profile); err != nil {
		t.Fatalf("创建不应该报错: %v", err)
	}

	if profile.Status != StatusPublished {
		t.Fatalf("审核开关关闭时应该直接发布，实际状态是 %s", profile.Status)
	}
	if !profile.ReviewBypassed {
		t.Fatal("审核开关关闭时应该带上免审核发布标记")
	}

	var reloaded ExpertDatabase.ExpertProfile
	db.First(&reloaded, profile.ID)
	if reloaded.ScoreUpdatedAt == nil {
		t.Fatal("免审核直接发布后应该立即重算一次得分")
	}
}
