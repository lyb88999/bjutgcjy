package ExpertDatabase

import (
	"testing"

	"github.com/flipped-aurora/gin-vue-admin/server/model/ExpertDatabase"
	ExpertDatabaseReq "github.com/flipped-aurora/gin-vue-admin/server/model/ExpertDatabase/request"
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
	tampered.TechTitle = "副教授"
	tampered.Status = StatusPublished
	tampered.TitleScore = 99
	tampered.TitleTier = "院士"
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
	if reloaded.TechTitle != "副教授" || reloaded.TitleScore != 0 || reloaded.TitleTier != "" || reloaded.ScoreUpdatedAt != nil {
		t.Fatalf("草稿可改职称原文，但不能写入或重算派生得分，实际 title=%s titleScore=%v tier=%q scoreUpdatedAt=%v",
			reloaded.TechTitle, reloaded.TitleScore, reloaded.TitleTier, reloaded.ScoreUpdatedAt)
	}
}

func TestUpdateExpertProfile_PublishedRecomputesTitleAndScores(t *testing.T) {
	db := setupTestDB(t)
	applicant := createTestUser(t, db, authorityIndividualApplicant, 0)
	profile := ExpertDatabase.ExpertProfile{
		Name: "已发布专家", Status: StatusPublished, TechTitle: "教授", CreatedBy: applicant.ID,
	}
	if err := db.Create(&profile).Error; err != nil {
		t.Fatalf("创建档案失败: %v", err)
	}
	achievement := ExpertDatabase.ExpertAchievement{ExpertId: profile.ID, AchievementType: "课题项目", Title: "已有成果", Level: "国家级"}
	if err := db.Create(&achievement).Error; err != nil {
		t.Fatalf("创建成果失败: %v", err)
	}
	if err := expertScoreSvc.RecomputeExpertScore(profile.ID); err != nil {
		t.Fatalf("初始化得分失败: %v", err)
	}
	if err := db.First(&profile, profile.ID).Error; err != nil {
		t.Fatalf("读取档案失败: %v", err)
	}

	// 更新普通画像字段时夹带伪造派生值，最终结果必须来自更新后的职称、荣誉及已有成果。
	profile.TechTitle = "副教授"
	profile.HonorTitle = "测试荣誉"
	profile.TitleTier = "院士"
	profile.TitleScore = 99
	profile.AchievementScore = 99
	profile.SocialScore = 99
	profile.CompositeScore = 99
	profile.Status = StatusDraft
	profile.ScoreUpdatedAt = nil
	svc := &ExpertProfileService{}
	if err := svc.UpdateExpertProfile(profile, applicant.ID); err != nil {
		t.Fatalf("已发布档案更新不应报错: %v", err)
	}

	var reloaded ExpertDatabase.ExpertProfile
	if err := db.First(&reloaded, profile.ID).Error; err != nil {
		t.Fatalf("读取更新后的档案失败: %v", err)
	}
	if reloaded.Status != StatusPublished || reloaded.TitleTier != "副教授" || reloaded.TitleScore != 3 {
		t.Fatalf("更新后应保持已发布并按实际职称派生档位，实际 status=%s tier=%q titleScore=%v",
			reloaded.Status, reloaded.TitleTier, reloaded.TitleScore)
	}
	if reloaded.AchievementScore != 10 || reloaded.SocialScore != 3 || reloaded.CompositeScore != 16 || reloaded.ScoreUpdatedAt == nil {
		t.Fatalf("成果、荣誉与职称应立即重算为综合分16，实际 achievement=%v social=%v composite=%v scoreUpdatedAt=%v",
			reloaded.AchievementScore, reloaded.SocialScore, reloaded.CompositeScore, reloaded.ScoreUpdatedAt)
	}
	search := &ExpertSearchService{}
	list, total, err := search.Search(ExpertDatabaseReq.ExpertSearchReq{TitleTier: "副教授", Page: 1, PageSize: 10})
	if err != nil || total != 1 || len(list) != 1 || list[0].ID != profile.ID || list[0].RealtimeScore != 16 {
		t.Fatalf("新职称分面应立即命中，并使用更新后的缓存得分，实际 total=%d list=%v err=%v", total, list, err)
	}
	_, oldTotal, err := search.Search(ExpertDatabaseReq.ExpertSearchReq{TitleTier: "教授", Page: 1, PageSize: 10})
	if err != nil || oldTotal != 0 {
		t.Fatalf("旧职称分面不应继续命中，实际 total=%d err=%v", oldTotal, err)
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
