package EngineeringEducationDatabase

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/EngineeringEducationDatabase"
    EngineeringEducationDatabaseReq "github.com/flipped-aurora/gin-vue-admin/server/model/EngineeringEducationDatabase/request"
    "gorm.io/gorm"
)

type TalentTrainingDatabaseService struct {
}

// CreateTalentTrainingDatabase 创建人才培养库记录
// Author [piexlmax](https://github.com/piexlmax)
func (TTDService *TalentTrainingDatabaseService) CreateTalentTrainingDatabase(TTD *EngineeringEducationDatabase.TalentTrainingDatabase) (err error) {
	err = global.GVA_DB.Create(TTD).Error
	return err
}

// DeleteTalentTrainingDatabase 删除人才培养库记录
// Author [piexlmax](https://github.com/piexlmax)
func (TTDService *TalentTrainingDatabaseService)DeleteTalentTrainingDatabase(ID string,userID uint) (err error) {
	err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
	    if err := tx.Model(&EngineeringEducationDatabase.TalentTrainingDatabase{}).Where("id = ?", ID).Update("deleted_by", userID).Error; err != nil {
              return err
        }
        if err = tx.Delete(&EngineeringEducationDatabase.TalentTrainingDatabase{},"id = ?",ID).Error; err != nil {
              return err
        }
        return nil
	})
	return err
}

// DeleteTalentTrainingDatabaseByIds 批量删除人才培养库记录
// Author [piexlmax](https://github.com/piexlmax)
func (TTDService *TalentTrainingDatabaseService)DeleteTalentTrainingDatabaseByIds(IDs []string,deleted_by uint) (err error) {
	err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
	    if err := tx.Model(&EngineeringEducationDatabase.TalentTrainingDatabase{}).Where("id in ?", IDs).Update("deleted_by", deleted_by).Error; err != nil {
            return err
        }
        if err := tx.Where("id in ?", IDs).Delete(&EngineeringEducationDatabase.TalentTrainingDatabase{}).Error; err != nil {
            return err
        }
        return nil
    })
	return err
}

// UpdateTalentTrainingDatabase 更新人才培养库记录
// Author [piexlmax](https://github.com/piexlmax)
func (TTDService *TalentTrainingDatabaseService)UpdateTalentTrainingDatabase(TTD EngineeringEducationDatabase.TalentTrainingDatabase) (err error) {
	err = global.GVA_DB.Save(&TTD).Error
	return err
}

// GetTalentTrainingDatabase 根据ID获取人才培养库记录
// Author [piexlmax](https://github.com/piexlmax)
func (TTDService *TalentTrainingDatabaseService)GetTalentTrainingDatabase(ID string) (TTD EngineeringEducationDatabase.TalentTrainingDatabase, err error) {
	err = global.GVA_DB.Where("id = ?", ID).First(&TTD).Error
	return
}

// GetTalentTrainingDatabaseInfoList 分页获取人才培养库记录
// Author [piexlmax](https://github.com/piexlmax)
func (TTDService *TalentTrainingDatabaseService)GetTalentTrainingDatabaseInfoList(info EngineeringEducationDatabaseReq.TalentTrainingDatabaseSearch) (list []EngineeringEducationDatabase.TalentTrainingDatabase, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
    // 创建db
	db := global.GVA_DB.Model(&EngineeringEducationDatabase.TalentTrainingDatabase{})
    var TTDs []EngineeringEducationDatabase.TalentTrainingDatabase
    // 如果有条件搜索 下方会自动创建搜索语句
    if info.StartCreatedAt !=nil && info.EndCreatedAt !=nil {
     db = db.Where("created_at BETWEEN ? AND ?", info.StartCreatedAt, info.EndCreatedAt)
    }
    if info.UniversityName != "" {
        db = db.Where("university_name LIKE ?","%"+ info.UniversityName+"%")
    }
	err = db.Count(&total).Error
	if err!=nil {
    	return
    }
        var OrderStr string
        orderMap := make(map[string]bool)
         	orderMap["undergraduate_count"] = true
         	orderMap["master_student_count"] = true
         	orderMap["doctoral_student_count"] = true
         	orderMap["national_first_class_course_count"] = true
         	orderMap["provincial_first_class_course_count"] = true
         	orderMap["ational_teaching_demonstration_center_count"] = true
         	orderMap["provincial_teaching_demonstration_center_count"] = true
         	orderMap["national_planning_textbook_count"] = true
         	orderMap["national_excellent_textbook_count"] = true
         	orderMap["provincial_excellent_textbook_count"] = true
         	orderMap["international_awards_count"] = true
         	orderMap["national_awards_count"] = true
         	orderMap["provincial_awards_count"] = true
         	orderMap["undergraduate_employment_rate"] = true
         	orderMap["master_employment_rate"] = true
         	orderMap["doctoral_employment_rate"] = true
         	orderMap["project_count"] = true
         	orderMap["participating_student_count"] = true
         	orderMap["transferred_into_engineering_count"] = true
         	orderMap["transferred_out_of_engineering_count"] = true
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
	
	err = db.Find(&TTDs).Error
	return  TTDs, total, err
}
