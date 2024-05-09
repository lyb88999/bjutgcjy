package EngineeringEducationDatabase

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/EngineeringEducationDatabase"
    EngineeringEducationDatabaseReq "github.com/flipped-aurora/gin-vue-admin/server/model/EngineeringEducationDatabase/request"
    "gorm.io/gorm"
)

type ProfessionalConstructionDatabaseService struct {
}

// CreateProfessionalConstructionDatabase 创建专业建设库记录
// Author [piexlmax](https://github.com/piexlmax)
func (PCDService *ProfessionalConstructionDatabaseService) CreateProfessionalConstructionDatabase(PCD *EngineeringEducationDatabase.ProfessionalConstructionDatabase) (err error) {
	err = global.GVA_DB.Create(PCD).Error
	return err
}

// DeleteProfessionalConstructionDatabase 删除专业建设库记录
// Author [piexlmax](https://github.com/piexlmax)
func (PCDService *ProfessionalConstructionDatabaseService)DeleteProfessionalConstructionDatabase(ID string,userID uint) (err error) {
	err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
	    if err := tx.Model(&EngineeringEducationDatabase.ProfessionalConstructionDatabase{}).Where("id = ?", ID).Update("deleted_by", userID).Error; err != nil {
              return err
        }
        if err = tx.Delete(&EngineeringEducationDatabase.ProfessionalConstructionDatabase{},"id = ?",ID).Error; err != nil {
              return err
        }
        return nil
	})
	return err
}

// DeleteProfessionalConstructionDatabaseByIds 批量删除专业建设库记录
// Author [piexlmax](https://github.com/piexlmax)
func (PCDService *ProfessionalConstructionDatabaseService)DeleteProfessionalConstructionDatabaseByIds(IDs []string,deleted_by uint) (err error) {
	err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
	    if err := tx.Model(&EngineeringEducationDatabase.ProfessionalConstructionDatabase{}).Where("id in ?", IDs).Update("deleted_by", deleted_by).Error; err != nil {
            return err
        }
        if err := tx.Where("id in ?", IDs).Delete(&EngineeringEducationDatabase.ProfessionalConstructionDatabase{}).Error; err != nil {
            return err
        }
        return nil
    })
	return err
}

// UpdateProfessionalConstructionDatabase 更新专业建设库记录
// Author [piexlmax](https://github.com/piexlmax)
func (PCDService *ProfessionalConstructionDatabaseService)UpdateProfessionalConstructionDatabase(PCD EngineeringEducationDatabase.ProfessionalConstructionDatabase) (err error) {
	err = global.GVA_DB.Save(&PCD).Error
	return err
}

// GetProfessionalConstructionDatabase 根据ID获取专业建设库记录
// Author [piexlmax](https://github.com/piexlmax)
func (PCDService *ProfessionalConstructionDatabaseService)GetProfessionalConstructionDatabase(ID string) (PCD EngineeringEducationDatabase.ProfessionalConstructionDatabase, err error) {
	err = global.GVA_DB.Where("id = ?", ID).First(&PCD).Error
	return
}

// GetProfessionalConstructionDatabaseInfoList 分页获取专业建设库记录
// Author [piexlmax](https://github.com/piexlmax)
func (PCDService *ProfessionalConstructionDatabaseService)GetProfessionalConstructionDatabaseInfoList(info EngineeringEducationDatabaseReq.ProfessionalConstructionDatabaseSearch) (list []EngineeringEducationDatabase.ProfessionalConstructionDatabase, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
    // 创建db
	db := global.GVA_DB.Model(&EngineeringEducationDatabase.ProfessionalConstructionDatabase{})
    var PCDs []EngineeringEducationDatabase.ProfessionalConstructionDatabase
    // 如果有条件搜索 下方会自动创建搜索语句
    if info.StartCreatedAt !=nil && info.EndCreatedAt !=nil {
     db = db.Where("created_at BETWEEN ? AND ?", info.StartCreatedAt, info.EndCreatedAt)
    }
    if info.MajorName != "" {
        db = db.Where("major_name LIKE ?","%"+ info.MajorName+"%")
    }
    if info.MajorCode != "" {
        db = db.Where("major_code = ?",info.MajorCode)
    }
    if info.DegreeType != "" {
        db = db.Where("degree_type = ?",info.DegreeType)
    }
    if info.IsNewMajor != nil {
        db = db.Where("is_new_major = ?",info.IsNewMajor)
    }
    if info.IsBroadAdmissionCategory != nil {
        db = db.Where("is_broad_admission_category = ?",info.IsBroadAdmissionCategory)
    }
    if info.IsFirstClassMajor != nil {
        db = db.Where("is_first_class_major = ?",info.IsFirstClassMajor)
    }
    if info.IsMajorAccredited != nil {
        db = db.Where("is_major_accredited = ?",info.IsMajorAccredited)
    }
    if info.IsSelectedForExcellenceProgram != nil {
        db = db.Where("is_selected_for_excellence_program = ?",info.IsSelectedForExcellenceProgram)
    }
    if info.IsDoubleDegree != nil {
        db = db.Where("is_double_degree = ?",info.IsDoubleDegree)
    }
    if info.IsNewEngineeringDiscipline != nil {
        db = db.Where("is_new_engineering_discipline = ?",info.IsNewEngineeringDiscipline)
    }
    if info.MajorSatisfaction != "" {
        db = db.Where("major_satisfaction = ?",info.MajorSatisfaction)
    }
    if info.UniversityName != "" {
        db = db.Where("university_name LIKE ?","%"+ info.UniversityName+"%")
    }
    if info.AdmissionMajor != "" {
        db = db.Where("admission_major LIKE ?","%"+ info.AdmissionMajor+"%")
    }
	err = db.Count(&total).Error
	if err!=nil {
    	return
    }
        var OrderStr string
        orderMap := make(map[string]bool)
         	orderMap["major_code"] = true
         	orderMap["planned_enrollment"] = true
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
	
	err = db.Find(&PCDs).Error
	return  PCDs, total, err
}
