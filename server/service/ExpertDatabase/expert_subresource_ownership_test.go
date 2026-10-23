package ExpertDatabase

import (
	"testing"

	"github.com/flipped-aurora/gin-vue-admin/server/model/ExpertDatabase"
)

// 研究成果/学术兼职/决策影响/标签这四个子资源之前完全没有数据域校验——个人申报人的账号能给
// 任意专家（不管是不是自己创建的）增删改这些子资源，跟专家档案本身"只能改自己创建的"这条规则
// 完全脱节。这里覆盖修复后的行为：子资源能不能改，看它挂的那条专家档案是不是操作人能改的。

func TestCreateExpertAchievement_OwnershipDenied(t *testing.T) {
	db := setupTestDB(t)
	svc := &ExpertAchievementService{}

	owner := createTestUser(t, db, authorityIndividualApplicant, 0)
	other := createTestUser(t, db, authorityIndividualApplicant, 0)
	profile := createTestProfile(t, db, StatusDraft, 0)
	profile.CreatedBy = owner.ID
	if err := db.Save(&profile).Error; err != nil {
		t.Fatalf("准备数据失败: %v", err)
	}

	achievement := &ExpertDatabase.ExpertAchievement{ExpertId: profile.ID, AchievementType: "学术论文", Title: "别人的成果"}
	if err := svc.CreateExpertAchievement(achievement, other.ID); err != ErrProfileAccessDenied {
		t.Fatalf("非档案创建者应被拒绝新增成果，实际 err=%v", err)
	}

	mine := &ExpertDatabase.ExpertAchievement{ExpertId: profile.ID, AchievementType: "学术论文", Title: "自己的成果"}
	if err := svc.CreateExpertAchievement(mine, owner.ID); err != nil {
		t.Fatalf("档案创建者本人新增成果不应报错: %v", err)
	}
}

func TestUpdateExpertAchievement_OwnershipAndExpertIdPinned(t *testing.T) {
	db := setupTestDB(t)
	svc := &ExpertAchievementService{}

	owner := createTestUser(t, db, authorityIndividualApplicant, 0)
	other := createTestUser(t, db, authorityIndividualApplicant, 0)
	profile := createTestProfile(t, db, StatusDraft, 0)
	profile.CreatedBy = owner.ID
	if err := db.Save(&profile).Error; err != nil {
		t.Fatalf("准备数据失败: %v", err)
	}
	otherProfile := createTestProfile(t, db, StatusDraft, 0)
	otherProfile.CreatedBy = other.ID
	if err := db.Save(&otherProfile).Error; err != nil {
		t.Fatalf("准备数据失败: %v", err)
	}

	achievement := ExpertDatabase.ExpertAchievement{ExpertId: profile.ID, AchievementType: "学术论文", Title: "原标题"}
	if err := db.Create(&achievement).Error; err != nil {
		t.Fatalf("准备数据失败: %v", err)
	}

	// 非档案创建者：拒绝
	tampered := achievement
	tampered.Title = "被篡改"
	if err := svc.UpdateExpertAchievement(tampered, other.ID); err != ErrProfileAccessDenied {
		t.Fatalf("他人档案下的成果应拒绝更新，实际 err=%v", err)
	}

	// 档案创建者本人：允许，且不能借更新把成果挪到别的专家名下
	moved := achievement
	moved.Title = "改过的标题"
	moved.ExpertId = otherProfile.ID
	if err := svc.UpdateExpertAchievement(moved, owner.ID); err != nil {
		t.Fatalf("本人档案下的成果更新不应报错: %v", err)
	}
	var reloaded ExpertDatabase.ExpertAchievement
	db.First(&reloaded, achievement.ID)
	if reloaded.Title != "改过的标题" {
		t.Fatalf("标题应正常更新，实际 title=%s", reloaded.Title)
	}
	if reloaded.ExpertId != profile.ID {
		t.Fatalf("expertId 不应被更新请求改动，实际 expertId=%d", reloaded.ExpertId)
	}
}

func TestDeleteExpertAchievement_OwnershipDenied(t *testing.T) {
	db := setupTestDB(t)
	svc := &ExpertAchievementService{}

	owner := createTestUser(t, db, authorityIndividualApplicant, 0)
	other := createTestUser(t, db, authorityIndividualApplicant, 0)
	profile := createTestProfile(t, db, StatusDraft, 0)
	profile.CreatedBy = owner.ID
	if err := db.Save(&profile).Error; err != nil {
		t.Fatalf("准备数据失败: %v", err)
	}
	achievement := ExpertDatabase.ExpertAchievement{ExpertId: profile.ID, AchievementType: "学术论文", Title: "成果"}
	if err := db.Create(&achievement).Error; err != nil {
		t.Fatalf("准备数据失败: %v", err)
	}

	if err := svc.DeleteExpertAchievement(idStr(achievement.ID), other.ID); err != ErrProfileAccessDenied {
		t.Fatalf("他人档案下的成果应拒绝删除，实际 err=%v", err)
	}
	if err := svc.DeleteExpertAchievement(idStr(achievement.ID), owner.ID); err != nil {
		t.Fatalf("本人档案下的成果删除不应报错: %v", err)
	}
}

func TestCreateExpertAcademicPosition_OwnershipDenied(t *testing.T) {
	db := setupTestDB(t)
	svc := &ExpertAcademicPositionService{}

	owner := createTestUser(t, db, authorityIndividualApplicant, 0)
	other := createTestUser(t, db, authorityIndividualApplicant, 0)
	profile := createTestProfile(t, db, StatusDraft, 0)
	profile.CreatedBy = owner.ID
	if err := db.Save(&profile).Error; err != nil {
		t.Fatalf("准备数据失败: %v", err)
	}

	position := &ExpertDatabase.ExpertAcademicPosition{ExpertId: profile.ID, OrganizationName: "某学会"}
	if err := svc.CreateExpertAcademicPosition(position, other.ID); err != ErrProfileAccessDenied {
		t.Fatalf("非档案创建者应被拒绝新增学术兼职，实际 err=%v", err)
	}
	if err := svc.CreateExpertAcademicPosition(position, owner.ID); err != nil {
		t.Fatalf("档案创建者本人新增学术兼职不应报错: %v", err)
	}
}

func TestCreateExpertAdoptionRecord_OwnershipDenied(t *testing.T) {
	db := setupTestDB(t)
	svc := &ExpertAdoptionRecordService{}

	owner := createTestUser(t, db, authorityIndividualApplicant, 0)
	other := createTestUser(t, db, authorityIndividualApplicant, 0)
	profile := createTestProfile(t, db, StatusDraft, 0)
	profile.CreatedBy = owner.ID
	if err := db.Save(&profile).Error; err != nil {
		t.Fatalf("准备数据失败: %v", err)
	}

	record := &ExpertDatabase.ExpertAdoptionRecord{ExpertId: profile.ID, AdoptionType: "领导批示"}
	if err := svc.CreateExpertAdoptionRecord(record, other.ID); err != ErrProfileAccessDenied {
		t.Fatalf("非档案创建者应被拒绝新增决策影响记录，实际 err=%v", err)
	}
	if err := svc.CreateExpertAdoptionRecord(record, owner.ID); err != nil {
		t.Fatalf("档案创建者本人新增决策影响记录不应报错: %v", err)
	}
}

func TestSetExpertTags_OwnershipDenied(t *testing.T) {
	db := setupTestDB(t)
	svc := &ExpertTagService{}

	owner := createTestUser(t, db, authorityIndividualApplicant, 0)
	other := createTestUser(t, db, authorityIndividualApplicant, 0)
	profile := createTestProfile(t, db, StatusDraft, 0)
	profile.CreatedBy = owner.ID
	if err := db.Save(&profile).Error; err != nil {
		t.Fatalf("准备数据失败: %v", err)
	}
	tag := ExpertDatabase.ExpertTag{TagType: "keyword", TagValue: "测试标签"}
	if err := db.Create(&tag).Error; err != nil {
		t.Fatalf("准备数据失败: %v", err)
	}

	if err := svc.SetExpertTags(profile.ID, []uint{tag.ID}, other.ID); err != ErrProfileAccessDenied {
		t.Fatalf("非档案创建者应被拒绝设置标签，实际 err=%v", err)
	}
	if err := svc.SetExpertTags(profile.ID, []uint{tag.ID}, owner.ID); err != nil {
		t.Fatalf("档案创建者本人设置标签不应报错: %v", err)
	}
}

// 管理员等不受 canModifyProfile 限制的角色，子资源照样不受限
func TestCreateExpertAchievement_AdminUnrestricted(t *testing.T) {
	db := setupTestDB(t)
	svc := &ExpertAchievementService{}

	admin := createTestUser(t, db, 1, 0)
	applicant := createTestUser(t, db, authorityIndividualApplicant, 0)
	profile := createTestProfile(t, db, StatusDraft, 0)
	profile.CreatedBy = applicant.ID
	if err := db.Save(&profile).Error; err != nil {
		t.Fatalf("准备数据失败: %v", err)
	}

	achievement := &ExpertDatabase.ExpertAchievement{ExpertId: profile.ID, AchievementType: "学术论文", Title: "管理员代录入的成果"}
	if err := svc.CreateExpertAchievement(achievement, admin.ID); err != nil {
		t.Fatalf("管理员不应受数据域限制，实际 err=%v", err)
	}
}
