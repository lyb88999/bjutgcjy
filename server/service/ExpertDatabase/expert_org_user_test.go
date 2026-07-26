package ExpertDatabase

import (
	"testing"

	ExpertDatabaseReq "github.com/flipped-aurora/gin-vue-admin/server/model/ExpertDatabase/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
)

func TestCreateOrgApplicant_ForcesOperatorOrgAndApplicantRole(t *testing.T) {
	db := setupTestDB(t)
	svc := &ExpertOrgUserService{}
	reviewer := createTestUser(t, db, authorityOrgReviewer, 3)

	err := svc.CreateOrgApplicant(reviewer.ID, ExpertDatabaseReq.CreateOrgApplicantRequest{
		Username: "applicant01", Password: "Test@123456",
	})
	if err != nil {
		t.Fatalf("建账号不应报错: %v", err)
	}

	var created system.SysUser
	if err := db.Where("username = ?", "applicant01").First(&created).Error; err != nil {
		t.Fatalf("查新账号失败: %v", err)
	}
	if created.OrgId == nil || *created.OrgId != 3 {
		t.Fatalf("新账号单位应强制取操作人的单位 3，实际 %v", created.OrgId)
	}
	if created.AuthorityId != authorityIndividualApplicant {
		t.Fatalf("新账号角色应强制为个人申报人，实际 %d", created.AuthorityId)
	}
	if created.Password == "Test@123456" {
		t.Fatal("密码应加密存储，不能落明文")
	}

	// 用户名占用应被拒
	err = svc.CreateOrgApplicant(reviewer.ID, ExpertDatabaseReq.CreateOrgApplicantRequest{
		Username: "applicant01", Password: "Another@123",
	})
	if err == nil {
		t.Fatal("重复用户名应被拒绝")
	}
}

func TestCreateOrgApplicant_OperatorWithoutOrgRejected(t *testing.T) {
	db := setupTestDB(t)
	svc := &ExpertOrgUserService{}
	noOrgReviewer := createTestUser(t, db, authorityOrgReviewer, 0)

	err := svc.CreateOrgApplicant(noOrgReviewer.ID, ExpertDatabaseReq.CreateOrgApplicantRequest{
		Username: "whoever", Password: "Test@123456",
	})
	if err == nil {
		t.Fatal("未关联单位的操作人不应能建账号")
	}
}

func TestGetOrgUserList_OnlyOwnOrgApplicants(t *testing.T) {
	db := setupTestDB(t)
	svc := &ExpertOrgUserService{}
	reviewer := createTestUser(t, db, authorityOrgReviewer, 3)

	mkUser := func(username string, authorityID uint, orgID uint) {
		u := system.SysUser{Username: username, AuthorityId: authorityID}
		if orgID != 0 {
			u.OrgId = &orgID
		}
		if err := db.Create(&u).Error; err != nil {
			t.Fatalf("建用户失败: %v", err)
		}
	}
	mkUser("mine", authorityIndividualApplicant, 3)     // 本单位申报人：应出现
	mkUser("otherOrg", authorityIndividualApplicant, 4) // 别的单位申报人：不应出现
	mkUser("sameOrgReviewer", authorityOrgReviewer, 3)  // 本单位审核员：不应出现

	list, total, err := svc.GetOrgUserList(reviewer.ID, ExpertDatabaseReq.OrgUserSearch{})
	if err != nil {
		t.Fatalf("查列表不应报错: %v", err)
	}
	if total != 1 || len(list) != 1 || list[0].Username != "mine" {
		t.Fatalf("只应看到本单位的申报人账号，实际 total=%d list=%+v", total, list)
	}
}

func TestToggleOrgUserEnable_ScopeChecks(t *testing.T) {
	db := setupTestDB(t)
	svc := &ExpertOrgUserService{}
	reviewer := createTestUser(t, db, authorityOrgReviewer, 3)

	org4 := uint(4)
	otherOrgUser := system.SysUser{Username: "otherOrg", AuthorityId: authorityIndividualApplicant, OrgId: &org4}
	if err := db.Create(&otherOrgUser).Error; err != nil {
		t.Fatalf("建用户失败: %v", err)
	}
	org3 := uint(3)
	sameOrgReviewer := system.SysUser{Username: "peerReviewer", AuthorityId: authorityOrgReviewer, OrgId: &org3}
	if err := db.Create(&sameOrgReviewer).Error; err != nil {
		t.Fatalf("建用户失败: %v", err)
	}
	mine := system.SysUser{Username: "mine", AuthorityId: authorityIndividualApplicant, OrgId: &org3, Enable: 1}
	if err := db.Create(&mine).Error; err != nil {
		t.Fatalf("建用户失败: %v", err)
	}

	if err := svc.ToggleOrgUserEnable(reviewer.ID, otherOrgUser.ID, 2); err == nil {
		t.Fatal("冻结别的单位的账号应被拒绝")
	}
	if err := svc.ToggleOrgUserEnable(reviewer.ID, sameOrgReviewer.ID, 2); err == nil {
		t.Fatal("冻结本单位其他审核员应被拒绝（只能管申报人账号）")
	}
	if err := svc.ToggleOrgUserEnable(reviewer.ID, mine.ID, 2); err != nil {
		t.Fatalf("冻结本单位申报人不应报错: %v", err)
	}
	var reloaded system.SysUser
	db.First(&reloaded, mine.ID)
	if reloaded.Enable != 2 {
		t.Fatalf("账号应已被冻结，实际 enable=%d", reloaded.Enable)
	}
}
