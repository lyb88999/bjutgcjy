package initialize

import (
	"os"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/example"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"

	"github.com/flipped-aurora/gin-vue-admin/server/model/EngineeringEducationDatabase"
	"github.com/flipped-aurora/gin-vue-admin/server/model/ExpertDatabase"
	"github.com/flipped-aurora/gin-vue-admin/server/model/util"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func Gorm() *gorm.DB {
	switch global.GVA_CONFIG.System.DbType {
	case "mysql":
		return GormMysql()
	case "pgsql":
		return GormPgSql()
	case "oracle":
		return GormOracle()
	case "mssql":
		return GormMssql()
	case "sqlite":
		return GormSqlite()
	default:
		return GormMysql()
	}
}

func RegisterTables() {
	db := global.GVA_DB
	err := db.AutoMigrate(

		system.SysApi{},
		system.SysUser{},
		system.SysBaseMenu{},
		system.JwtBlacklist{},
		system.SysAuthority{},
		system.SysDictionary{},
		system.SysOperationRecord{},
		system.SysAutoCodeHistory{},
		system.SysDictionaryDetail{},
		system.SysBaseMenuParameter{},
		system.SysBaseMenuBtn{},
		system.SysAuthorityBtn{},
		system.SysAutoCode{},
		system.SysExportTemplate{},
		system.Condition{},
		system.SysOrganization{},

		example.ExaFile{},
		example.ExaCustomer{},
		example.ExaFileChunk{},
		example.ExaFileUploadAndDownload{}, EngineeringEducationDatabase.SocialServiceDatabase{}, EngineeringEducationDatabase.ConditionalGuaranteeDatabase{}, EngineeringEducationDatabase.PolicyDatabase{}, EngineeringEducationDatabase.ProfessionalConstructionDatabase{}, EngineeringEducationDatabase.SubjectConstructionDatabase{}, EngineeringEducationDatabase.TalentTrainingDatabase{}, EngineeringEducationDatabase.InternationalExchangeDatabase{}, EngineeringEducationDatabase.BasicInfomationDatabase{}, util.InviteCode{},

		ExpertDatabase.ExpertProfile{}, ExpertDatabase.ExpertAchievement{}, ExpertDatabase.ExpertAdoptionRecord{}, ExpertDatabase.ExpertTag{}, ExpertDatabase.ExpertTagRelation{}, ExpertDatabase.ExpertAcademicPosition{}, ExpertDatabase.ExpertApprovalLog{}, ExpertDatabase.ExpertSearchEmbedding{},
	)
	if err != nil {
		global.GVA_LOG.Error("register table failed", zap.Error(err))
		os.Exit(0)
	}
	global.GVA_LOG.Info("register table success")
}
