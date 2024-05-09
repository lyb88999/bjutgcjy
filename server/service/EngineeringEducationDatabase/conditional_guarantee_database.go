package EngineeringEducationDatabase

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/EngineeringEducationDatabase"
	EngineeringEducationDatabaseReq "github.com/flipped-aurora/gin-vue-admin/server/model/EngineeringEducationDatabase/request"
	"gorm.io/gorm"
)

type ConditionalGuaranteeDatabaseService struct {
}

// CreateConditionalGuaranteeDatabase 创建条件保障库记录
// Author [piexlmax](https://github.com/piexlmax)
func (CGDService *ConditionalGuaranteeDatabaseService) CreateConditionalGuaranteeDatabase(CGD *EngineeringEducationDatabase.ConditionalGuaranteeDatabase) (err error) {
	err = global.GVA_DB.Create(CGD).Error
	return err
}

// DeleteConditionalGuaranteeDatabase 删除条件保障库记录
// Author [piexlmax](https://github.com/piexlmax)
func (CGDService *ConditionalGuaranteeDatabaseService) DeleteConditionalGuaranteeDatabase(ID string, userID uint) (err error) {
	err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&EngineeringEducationDatabase.ConditionalGuaranteeDatabase{}).Where("id = ?", ID).Update("deleted_by", userID).Error; err != nil {
			return err
		}
		if err = tx.Delete(&EngineeringEducationDatabase.ConditionalGuaranteeDatabase{}, "id = ?", ID).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

// DeleteConditionalGuaranteeDatabaseByIds 批量删除条件保障库记录
// Author [piexlmax](https://github.com/piexlmax)
func (CGDService *ConditionalGuaranteeDatabaseService) DeleteConditionalGuaranteeDatabaseByIds(IDs []string, deleted_by uint) (err error) {
	err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&EngineeringEducationDatabase.ConditionalGuaranteeDatabase{}).Where("id in ?", IDs).Update("deleted_by", deleted_by).Error; err != nil {
			return err
		}
		if err := tx.Where("id in ?", IDs).Delete(&EngineeringEducationDatabase.ConditionalGuaranteeDatabase{}).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

// UpdateConditionalGuaranteeDatabase 更新条件保障库记录
// Author [piexlmax](https://github.com/piexlmax)
func (CGDService *ConditionalGuaranteeDatabaseService) UpdateConditionalGuaranteeDatabase(CGD EngineeringEducationDatabase.ConditionalGuaranteeDatabase) (err error) {
	err = global.GVA_DB.Save(&CGD).Error
	return err
}

// GetConditionalGuaranteeDatabase 根据ID获取条件保障库记录
// Author [piexlmax](https://github.com/piexlmax)
func (CGDService *ConditionalGuaranteeDatabaseService) GetConditionalGuaranteeDatabase(ID string) (CGD EngineeringEducationDatabase.ConditionalGuaranteeDatabase, err error) {
	err = global.GVA_DB.Where("id = ?", ID).First(&CGD).Error
	return
}

// GetConditionalGuaranteeDatabaseInfoList 分页获取条件保障库记录
// Author [piexlmax](https://github.com/piexlmax)
func (CGDService *ConditionalGuaranteeDatabaseService) GetConditionalGuaranteeDatabaseInfoList(info EngineeringEducationDatabaseReq.ConditionalGuaranteeDatabaseSearch) (list []EngineeringEducationDatabase.ConditionalGuaranteeDatabase, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	// 创建db
	db := global.GVA_DB.Model(&EngineeringEducationDatabase.ConditionalGuaranteeDatabase{})
	var CGDs []EngineeringEducationDatabase.ConditionalGuaranteeDatabase
	// 如果有条件搜索 下方会自动创建搜索语句
	if info.StartCreatedAt != nil && info.EndCreatedAt != nil {
		db = db.Where("created_at BETWEEN ? AND ?", info.StartCreatedAt, info.EndCreatedAt)
	}
	if info.UniversityName != "" {
		// db = db.Where("university_name = ?",info.UniversityName)
		db = db.Where("university_name LIKE ?", "%"+info.UniversityName+"%")
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	var OrderStr string
	orderMap := make(map[string]bool)
	orderMap["full_time_faculty_count"] = true
	orderMap["foreign_high_level_talent_count"] = true
	orderMap["academician_count"] = true
	orderMap["national_talent_project_winner_count"] = true
	orderMap["provincial_talent_project_winner_count"] = true
	orderMap["off_campus_internship_base_count"] = true
	orderMap["large_equipment_sharing_platform_count"] = true
	orderMap["national_experimental_teaching_demo_center_count"] = true
	if orderMap[info.Sort] {
		OrderStr = info.Sort
		if info.Order == "descending" {
			OrderStr = OrderStr + " desc"
		}
		db = db.Order(OrderStr)
	}

	if limit != 0 {
		db = db.Limit(limit).Offset(offset)
	}

	err = db.Find(&CGDs).Error
	return CGDs, total, err
}
