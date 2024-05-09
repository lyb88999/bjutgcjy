package EngineeringEducationDatabase

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/EngineeringEducationDatabase"
    EngineeringEducationDatabaseReq "github.com/flipped-aurora/gin-vue-admin/server/model/EngineeringEducationDatabase/request"
)

type SubjectConstructionDatabaseService struct {
}

// CreateSubjectConstructionDatabase 创建学科建设库记录
// Author [piexlmax](https://github.com/piexlmax)
func (SCDService *SubjectConstructionDatabaseService) CreateSubjectConstructionDatabase(SCD *EngineeringEducationDatabase.SubjectConstructionDatabase) (err error) {
	err = global.GVA_DB.Create(SCD).Error
	return err
}

// DeleteSubjectConstructionDatabase 删除学科建设库记录
// Author [piexlmax](https://github.com/piexlmax)
func (SCDService *SubjectConstructionDatabaseService)DeleteSubjectConstructionDatabase(ID string) (err error) {
	err = global.GVA_DB.Delete(&EngineeringEducationDatabase.SubjectConstructionDatabase{},"id = ?",ID).Error
	return err
}

// DeleteSubjectConstructionDatabaseByIds 批量删除学科建设库记录
// Author [piexlmax](https://github.com/piexlmax)
func (SCDService *SubjectConstructionDatabaseService)DeleteSubjectConstructionDatabaseByIds(IDs []string) (err error) {
	err = global.GVA_DB.Delete(&[]EngineeringEducationDatabase.SubjectConstructionDatabase{},"id in ?",IDs).Error
	return err
}

// UpdateSubjectConstructionDatabase 更新学科建设库记录
// Author [piexlmax](https://github.com/piexlmax)
func (SCDService *SubjectConstructionDatabaseService)UpdateSubjectConstructionDatabase(SCD EngineeringEducationDatabase.SubjectConstructionDatabase) (err error) {
	err = global.GVA_DB.Save(&SCD).Error
	return err
}

// GetSubjectConstructionDatabase 根据ID获取学科建设库记录
// Author [piexlmax](https://github.com/piexlmax)
func (SCDService *SubjectConstructionDatabaseService)GetSubjectConstructionDatabase(ID string) (SCD EngineeringEducationDatabase.SubjectConstructionDatabase, err error) {
	err = global.GVA_DB.Where("id = ?", ID).First(&SCD).Error
	return
}

// GetSubjectConstructionDatabaseInfoList 分页获取学科建设库记录
// Author [piexlmax](https://github.com/piexlmax)
func (SCDService *SubjectConstructionDatabaseService)GetSubjectConstructionDatabaseInfoList(info EngineeringEducationDatabaseReq.SubjectConstructionDatabaseSearch) (list []EngineeringEducationDatabase.SubjectConstructionDatabase, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
    // 创建db
	db := global.GVA_DB.Model(&EngineeringEducationDatabase.SubjectConstructionDatabase{})
    var SCDs []EngineeringEducationDatabase.SubjectConstructionDatabase
    // 如果有条件搜索 下方会自动创建搜索语句
    if info.StartCreatedAt !=nil && info.EndCreatedAt !=nil {
     db = db.Where("created_at BETWEEN ? AND ?", info.StartCreatedAt, info.EndCreatedAt)
    }
    if info.PostdocStationName != "" {
        db = db.Where("postdoc_station_name LIKE ?","%"+ info.PostdocStationName+"%")
    }
    if info.PostdocSubjectCategory != "" {
        db = db.Where("postdoc_subject_category LIKE ?","%"+ info.PostdocSubjectCategory+"%")
    }
    if info.MasterSubjectCode != "" {
        db = db.Where("master_subject_code = ?",info.MasterSubjectCode)
    }
    if info.MasterSubjectCategory != "" {
        db = db.Where("master_subject_category LIKE ?","%"+ info.MasterSubjectCategory+"%")
    }
    if info.MasterType != "" {
        db = db.Where("master_type = ?",info.MasterType)
    }
    if info.MasterIsFirstClass != nil {
        db = db.Where("master_is_first_class = ?",info.MasterIsFirstClass)
    }
    if info.MasterUniversityName != "" {
        db = db.Where("master_university_name LIKE ?","%"+ info.MasterUniversityName+"%")
    }
    if info.PhdSubjectCode != "" {
        db = db.Where("phd_subject_code = ?",info.PhdSubjectCode)
    }
    if info.PhdSubjectCategory != "" {
        db = db.Where("phd_subject_category LIKE ?","%"+ info.PhdSubjectCategory+"%")
    }
    if info.PhdType != "" {
        db = db.Where("phd_type = ?",info.PhdType)
    }
    if info.PhdIsFirstClass != nil {
        db = db.Where("phd_is_first_class = ?",info.PhdIsFirstClass)
    }
    if info.PhdUniversityName != "" {
        db = db.Where("phd_university_name LIKE ?","%"+ info.PhdUniversityName+"%")
    }
	err = db.Count(&total).Error
	if err!=nil {
    	return
    }
        var OrderStr string
        orderMap := make(map[string]bool)
         	orderMap["postdoc_station_count"] = true
         	orderMap["master_subject_code"] = true
         	orderMap["master_soft_science_ranking"] = true
         	orderMap["master_admission_number"] = true
         	orderMap["phd_subject_code"] = true
         	orderMap["phd_soft_science_ranking"] = true
         	orderMap["phd_admission_number"] = true
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
	
	err = db.Find(&SCDs).Error
	return  SCDs, total, err
}
