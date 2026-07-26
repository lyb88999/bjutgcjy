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

// 详情接口的数据域收敛必须跟列表页一致：列表上看不到的记录，拿着 ID 直接查详情也必须被拒
func TestGetExpertProfile_RoleScope(t *testing.T) {
	db := setupTestDB(t)
	svc := &ExpertProfileService{}

	orgReviewer := createTestUser(t, db, authorityOrgReviewer, 1)
	cityReviewer := createTestUser(t, db, authorityCityReviewer, 0)
	admin := createTestUser(t, db, 1, 0)

	otherOrgDraft := createTestProfile(t, db, StatusDraft, 2)
	sameOrgPending := createTestProfile(t, db, StatusPendingOrgReview, 1)
	published := createTestProfile(t, db, StatusPublished, 2)

	cases := []struct {
		name       string
		operator   uint
		profileID  uint
		wantDenied bool
	}{
		{"单位审核员查其他单位草稿应被拒", orgReviewer.ID, otherOrgDraft.ID, true},
		{"单位审核员查本单位待审记录应放行", orgReviewer.ID, sameOrgPending.ID, false},
		{"单位审核员查已发布记录应放行", orgReviewer.ID, published.ID, false},
		{"市级审核员查待单位审核记录应被拒", cityReviewer.ID, sameOrgPending.ID, true},
		{"市级审核员查已发布记录应放行", cityReviewer.ID, published.ID, false},
		{"管理员查任何记录都放行", admin.ID, otherOrgDraft.ID, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := svc.GetExpertProfile(idStr(tc.profileID), tc.operator)
			if tc.wantDenied && err != ErrProfileAccessDenied {
				t.Fatalf("期望被数据域校验拒绝，实际 err=%v", err)
			}
			if !tc.wantDenied && err != nil {
				t.Fatalf("期望放行，实际 err=%v", err)
			}
		})
	}
}

// 个人申报人只能改自己创建的档案；不能借更新接口夹带 status/得分字段绕过审核流程
func TestUpdateExpertProfile_OwnershipAndFieldProtection(t *testing.T) {
	db := setupTestDB(t)
	svc := &ExpertProfileService{}

	applicant := createTestUser(t, db, authorityIndividualApplicant, 0)
	other := createTestUser(t, db, authorityIndividualApplicant, 0)

	profile := createTestProfile(t, db, StatusDraft, 0)
	profile.CreatedBy = applicant.ID
	if err := db.Save(&profile).Error; err != nil {
		t.Fatalf("准备数据失败: %v", err)
	}

	// 别人的档案：拒绝
	profile.Name = "被篡改"
	if err := svc.UpdateExpertProfile(profile, other.ID); err != ErrProfileAccessDenied {
		t.Fatalf("他人档案应被拒绝更新，实际 err=%v", err)
	}

	// 自己的档案：允许，但请求体里夹带的状态/得分字段必须被库里现值覆盖
	tampered := profile
	tampered.Name = "正常改名"
	tampered.Status = StatusPublished
	tampered.CompositeScore = 99
	tampered.ReviewBypassed = true
	if err := svc.UpdateExpertProfile(tampered, applicant.ID); err != nil {
		t.Fatalf("本人档案更新不应报错: %v", err)
	}
	var reloaded ExpertDatabase.ExpertProfile
	db.First(&reloaded, profile.ID)
	if reloaded.Name != "正常改名" {
		t.Fatalf("画像字段应正常更新，实际 name=%s", reloaded.Name)
	}
	if reloaded.Status != StatusDraft || reloaded.CompositeScore != 0 || reloaded.ReviewBypassed {
		t.Fatalf("状态机/得分字段不应被更新接口改动，实际 status=%s score=%v bypassed=%v",
			reloaded.Status, reloaded.CompositeScore, reloaded.ReviewBypassed)
	}
}

// 批量删除只要有一条越出数据域就整体拒绝
func TestDeleteExpertProfileByIds_ScopeIsAllOrNothing(t *testing.T) {
	db := setupTestDB(t)
	svc := &ExpertProfileService{}

	orgReviewer := createTestUser(t, db, authorityOrgReviewer, 1)
	mine := createTestProfile(t, db, StatusPendingOrgReview, 1)
	others := createTestProfile(t, db, StatusPendingOrgReview, 2)

	err := svc.DeleteExpertProfileByIds([]string{idStr(mine.ID), idStr(others.ID)}, orgReviewer.ID)
	if err != ErrProfileAccessDenied {
		t.Fatalf("混入他单位记录的批量删除应整体拒绝，实际 err=%v", err)
	}
	var count int64
	db.Model(&ExpertDatabase.ExpertProfile{}).Count(&count)
	if count != 2 {
		t.Fatalf("整体拒绝时不应删除任何记录，剩余 %d 条", count)
	}

	if err := svc.DeleteExpertProfileByIds([]string{idStr(mine.ID)}, orgReviewer.ID); err != nil {
		t.Fatalf("本单位记录删除不应报错: %v", err)
	}
}
