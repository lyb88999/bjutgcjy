package ExpertDatabase

import (
	"testing"

	"github.com/flipped-aurora/gin-vue-admin/server/model/ExpertDatabase"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
	"gorm.io/gorm"
)

func TestTransition_AllowedMoveUpdatesStatusAndWritesLog(t *testing.T) {
	db := setupTestDB(t)
	profile := createTestProfile(t, db, StatusDraft, 0)
	svc := &ExpertApprovalService{}

	err := db.Transaction(func(tx *gorm.DB) error {
		return svc.transition(tx, &profile, []string{StatusDraft}, StatusPendingOrgReview, 1, "")
	})
	if err != nil {
		t.Fatalf("合法流转不应该报错: %v", err)
	}

	var reloaded ExpertDatabase.ExpertProfile
	db.First(&reloaded, profile.ID)
	if reloaded.Status != StatusPendingOrgReview {
		t.Fatalf("状态应该变成 %s，实际是 %s", StatusPendingOrgReview, reloaded.Status)
	}

	var logs []ExpertDatabase.ExpertApprovalLog
	db.Where("expert_id = ?", profile.ID).Find(&logs)
	if len(logs) != 1 {
		t.Fatalf("应该写入 1 条审核日志，实际 %d 条", len(logs))
	}
	if logs[0].FromStatus != StatusDraft || logs[0].ToStatus != StatusPendingOrgReview {
		t.Fatalf("日志的流转记录不对：from=%s to=%s", logs[0].FromStatus, logs[0].ToStatus)
	}
}

func TestTransition_DisallowedMoveRejectedAndNoLogWritten(t *testing.T) {
	db := setupTestDB(t)
	// 已经发布的档案不应该允许直接"提交审核"（只有草稿/被退回状态才能提交）
	profile := createTestProfile(t, db, StatusPublished, 0)
	svc := &ExpertApprovalService{}

	err := db.Transaction(func(tx *gorm.DB) error {
		return svc.transition(tx, &profile, []string{StatusDraft, StatusOrgRejected, StatusCityRejected}, StatusPendingOrgReview, 1, "")
	})
	if err == nil {
		t.Fatal("已发布状态不应该允许流转到待单位审核，应该报错")
	}

	var reloaded ExpertDatabase.ExpertProfile
	db.First(&reloaded, profile.ID)
	if reloaded.Status != StatusPublished {
		t.Fatalf("非法流转不应该改动状态，实际变成了 %s", reloaded.Status)
	}

	var count int64
	db.Model(&ExpertDatabase.ExpertApprovalLog{}).Where("expert_id = ?", profile.ID).Count(&count)
	if count != 0 {
		t.Fatalf("非法流转不应该写入审核日志，实际写了 %d 条", count)
	}
}

func TestSubmit_WithOrgRoutesToOrgReview(t *testing.T) {
	db := setupTestDB(t)
	profile := createTestProfile(t, db, StatusDraft, 2)
	svc := &ExpertApprovalService{}

	if err := svc.Submit(profile.ID, 1); err != nil {
		t.Fatalf("提交不应该报错: %v", err)
	}

	var reloaded ExpertDatabase.ExpertProfile
	db.First(&reloaded, profile.ID)
	if reloaded.Status != StatusPendingOrgReview {
		t.Fatalf("有所属单位的档案提交后应该是待单位审核，实际是 %s", reloaded.Status)
	}
}

func TestSubmit_WithoutOrgFallsBackToCityReview(t *testing.T) {
	db := setupTestDB(t)
	// org_id 为空——没有单位可以归口，Submit 应该跳过单位审核，直接进市级审核，不能卡住流程
	profile := createTestProfile(t, db, StatusDraft, 0)
	svc := &ExpertApprovalService{}

	if err := svc.Submit(profile.ID, 1); err != nil {
		t.Fatalf("提交不应该报错: %v", err)
	}

	var reloaded ExpertDatabase.ExpertProfile
	db.First(&reloaded, profile.ID)
	if reloaded.Status != StatusPendingCityReview {
		t.Fatalf("没有所属单位的档案提交后应该直接进待市级审核，实际是 %s", reloaded.Status)
	}
}

func TestAssertSameOrg_DifferentOrgRejected(t *testing.T) {
	db := setupTestDB(t)
	reviewer := createTestUser(t, db, authorityOrgReviewer, 2)     // 挂靠单位2
	profile := createTestProfile(t, db, StatusPendingOrgReview, 3) // 档案属于单位3
	svc := &ExpertApprovalService{}

	if err := svc.OrgApprove(profile.ID, reviewer.ID); err == nil {
		t.Fatal("单位审核员不应该能审核别的单位提交的档案")
	}

	var reloaded ExpertDatabase.ExpertProfile
	db.First(&reloaded, profile.ID)
	if reloaded.Status != StatusPendingOrgReview {
		t.Fatalf("越权审核不应该改动状态，实际变成了 %s", reloaded.Status)
	}
}

func TestAssertSameOrg_SameOrgApproved(t *testing.T) {
	db := setupTestDB(t)
	reviewer := createTestUser(t, db, authorityOrgReviewer, 2)
	profile := createTestProfile(t, db, StatusPendingOrgReview, 2)
	svc := &ExpertApprovalService{}

	if err := svc.OrgApprove(profile.ID, reviewer.ID); err != nil {
		t.Fatalf("同单位审核不应该报错: %v", err)
	}

	var reloaded ExpertDatabase.ExpertProfile
	db.First(&reloaded, profile.ID)
	if reloaded.Status != StatusPendingCityReview {
		t.Fatalf("单位审核通过后应该进入待市级审核，实际是 %s", reloaded.Status)
	}
}

func TestAssertSameOrg_SuperAdminBypassesOrgCheck(t *testing.T) {
	db := setupTestDB(t)
	admin := createTestUser(t, db, superAdminAuthorityId, 0)       // 管理员不挂靠任何单位
	profile := createTestProfile(t, db, StatusPendingOrgReview, 5) // 档案属于单位5，管理员不属于这个单位
	svc := &ExpertApprovalService{}

	if err := svc.OrgApprove(profile.ID, admin.ID); err != nil {
		t.Fatalf("管理员应该能审核任意单位提交的档案，不应该报错: %v", err)
	}
}

func TestCityApprove_PublishesRecomputesScoreAndClearsBypassFlag(t *testing.T) {
	db := setupTestDB(t)
	profile := ExpertDatabase.ExpertProfile{Name: "测试专家", Status: StatusPendingCityReview, ReviewBypassed: true}
	if err := db.Create(&profile).Error; err != nil {
		t.Fatalf("建档案失败: %v", err)
	}
	svc := &ExpertApprovalService{}

	if err := svc.CityApprove(profile.ID, 1); err != nil {
		t.Fatalf("市级审核通过不应该报错: %v", err)
	}

	var reloaded ExpertDatabase.ExpertProfile
	db.First(&reloaded, profile.ID)
	if reloaded.Status != StatusPublished {
		t.Fatalf("市级审核通过后应该是已发布，实际是 %s", reloaded.Status)
	}
	if reloaded.ReviewBypassed {
		t.Fatal("走完真实审核之后，免审核发布标记应该被清掉")
	}
	if reloaded.ScoreUpdatedAt == nil {
		t.Fatal("发布后应该立即算过一次得分，score_updated_at 不应该是空的")
	}
}

func TestBatchOrgApprove_PartialFailureDoesNotBlockOthers(t *testing.T) {
	db := setupTestDB(t)
	reviewer := createTestUser(t, db, authorityOrgReviewer, 2)
	ok1 := createTestProfile(t, db, StatusPendingOrgReview, 2)
	ok2 := createTestProfile(t, db, StatusPendingOrgReview, 2)
	wrongOrg := createTestProfile(t, db, StatusPendingOrgReview, 3) // 不是审核员自己的单位，这条应该失败
	svc := &ExpertApprovalService{}

	result := svc.BatchOrgApprove([]uint{ok1.ID, ok2.ID, wrongOrg.ID}, reviewer.ID)

	if result.SuccessCount != 2 {
		t.Fatalf("应该成功 2 条，实际 %d 条", result.SuccessCount)
	}
	if result.FailCount != 1 {
		t.Fatalf("应该失败 1 条，实际 %d 条", result.FailCount)
	}
	if len(result.Failures) != 1 || result.Failures[0].ExpertId != wrongOrg.ID {
		t.Fatalf("失败列表应该只包含 wrongOrg 这一条，实际: %+v", result.Failures)
	}

	var reloadedOK1, reloadedWrong ExpertDatabase.ExpertProfile
	db.First(&reloadedOK1, ok1.ID)
	db.First(&reloadedWrong, wrongOrg.ID)
	if reloadedOK1.Status != StatusPendingCityReview {
		t.Fatalf("成功的那条应该流转到待市级审核，实际是 %s", reloadedOK1.Status)
	}
	if reloadedWrong.Status != StatusPendingOrgReview {
		t.Fatalf("失败的那条状态不应该被改动，实际是 %s", reloadedWrong.Status)
	}
}

func TestReviewRequired_DefaultsToTrueWhenDictionaryMissing(t *testing.T) {
	setupTestDB(t)
	// 字典没配置时必须默认"需要审核"——这是个安全默认值，不能因为漏配字典就意外把审核整个绕过去
	if !reviewRequired() {
		t.Fatal("字典缺失时应该默认需要审核")
	}
}

func TestReviewRequired_ReadsFromDictionary(t *testing.T) {
	db := setupTestDB(t)
	dict := system.SysDictionary{Name: "专家库-审核开关", Type: DictTypeReviewSwitch}
	if err := db.Create(&dict).Error; err != nil {
		t.Fatalf("建字典失败: %v", err)
	}
	detail := system.SysDictionaryDetail{Label: "启用审核", Value: 0, SysDictionaryID: int(dict.ID)}
	if err := db.Create(&detail).Error; err != nil {
		t.Fatalf("建字典明细失败: %v", err)
	}

	if reviewRequired() {
		t.Fatal("字典明细 value=0 时应该读到关闭审核")
	}

	db.Model(&detail).Update("value", 1)
	if !reviewRequired() {
		t.Fatal("字典明细 value=1 时应该读到需要审核")
	}
}
