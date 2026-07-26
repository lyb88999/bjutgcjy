package ExpertDatabase

import (
	"testing"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/ExpertDatabase"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// setupTestDB 给每个测试用例起一个全新的内存 sqlite 库，AutoMigrate 出业务逻辑用得到的表，
// 换掉全局的 global.GVA_DB——测完即弃，测试之间互不影响，不需要真的连 MySQL
func setupTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	// sqlite 的普通 :memory: DSN 是"每个连接一份独立的库"——业务代码里有些地方（比如
	// assertSameOrg）会在事务内部另外用 global.GVA_DB（而不是事务的 tx）发起查询，这种写法在
	// MySQL 下没问题，但如果测试库只有一个连接，就会跟当前占着这条连接的事务自己互相等待、死锁。
	// 用带名字的 shared-cache 内存库，同一个测试用例里开出的所有连接都指向同一份数据；
	// 名字用 t.Name() 保证测试用例之间各用各的库，不会串数据。
	dsn := "file:" + t.Name() + "?mode=memory&cache=shared"
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		t.Fatalf("打开内存数据库失败: %v", err)
	}
	if err := db.AutoMigrate(
		&system.SysAuthority{},
		&system.SysUser{},
		&system.SysOrganization{},
		&system.SysDictionary{},
		&system.SysDictionaryDetail{},
		&ExpertDatabase.ExpertProfile{},
		&ExpertDatabase.ExpertApprovalLog{},
		&ExpertDatabase.ExpertAchievement{},
		&ExpertDatabase.ExpertAdoptionRecord{},
		&ExpertDatabase.ExpertAcademicPosition{},
		&ExpertDatabase.ExpertTag{},
		&ExpertDatabase.ExpertTagRelation{},
	); err != nil {
		t.Fatalf("建表失败: %v", err)
	}
	global.GVA_DB = db
	return db
}

// createTestUser 建一个测试用户，orgID 为 0 表示不挂靠任何单位（对应 profile.OrgId == nil 的场景）
func createTestUser(t *testing.T, db *gorm.DB, authorityID uint, orgID uint) system.SysUser {
	t.Helper()
	user := system.SysUser{AuthorityId: authorityID, Username: "test-user"}
	if orgID != 0 {
		user.OrgId = &orgID
	}
	if err := db.Create(&user).Error; err != nil {
		t.Fatalf("建用户失败: %v", err)
	}
	return user
}

// createTestProfile 建一条测试专家档案，orgID 为 0 表示 org_id 为空
func createTestProfile(t *testing.T, db *gorm.DB, status string, orgID uint) ExpertDatabase.ExpertProfile {
	t.Helper()
	profile := ExpertDatabase.ExpertProfile{Name: "测试专家", Status: status}
	if orgID != 0 {
		profile.OrgId = &orgID
	}
	if err := db.Create(&profile).Error; err != nil {
		t.Fatalf("建专家档案失败: %v", err)
	}
	return profile
}
