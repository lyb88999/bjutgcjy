package EngineeringEducationDatabase

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/EngineeringEducationDatabase"
    EngineeringEducationDatabaseReq "github.com/flipped-aurora/gin-vue-admin/server/model/EngineeringEducationDatabase/request"
)

type InternationalExchangeDatabaseService struct {
}

// CreateInternationalExchangeDatabase 创建国际交流库记录
// Author [piexlmax](https://github.com/piexlmax)
func (IEDService *InternationalExchangeDatabaseService) CreateInternationalExchangeDatabase(IED *EngineeringEducationDatabase.InternationalExchangeDatabase) (err error) {
	err = global.GVA_DB.Create(IED).Error
	return err
}

// DeleteInternationalExchangeDatabase 删除国际交流库记录
// Author [piexlmax](https://github.com/piexlmax)
func (IEDService *InternationalExchangeDatabaseService)DeleteInternationalExchangeDatabase(ID string) (err error) {
	err = global.GVA_DB.Delete(&EngineeringEducationDatabase.InternationalExchangeDatabase{},"id = ?",ID).Error
	return err
}

// DeleteInternationalExchangeDatabaseByIds 批量删除国际交流库记录
// Author [piexlmax](https://github.com/piexlmax)
func (IEDService *InternationalExchangeDatabaseService)DeleteInternationalExchangeDatabaseByIds(IDs []string) (err error) {
	err = global.GVA_DB.Delete(&[]EngineeringEducationDatabase.InternationalExchangeDatabase{},"id in ?",IDs).Error
	return err
}

// UpdateInternationalExchangeDatabase 更新国际交流库记录
// Author [piexlmax](https://github.com/piexlmax)
func (IEDService *InternationalExchangeDatabaseService)UpdateInternationalExchangeDatabase(IED EngineeringEducationDatabase.InternationalExchangeDatabase) (err error) {
	err = global.GVA_DB.Save(&IED).Error
	return err
}

// GetInternationalExchangeDatabase 根据ID获取国际交流库记录
// Author [piexlmax](https://github.com/piexlmax)
func (IEDService *InternationalExchangeDatabaseService)GetInternationalExchangeDatabase(ID string) (IED EngineeringEducationDatabase.InternationalExchangeDatabase, err error) {
	err = global.GVA_DB.Where("id = ?", ID).First(&IED).Error
	return
}

// GetInternationalExchangeDatabaseInfoList 分页获取国际交流库记录
// Author [piexlmax](https://github.com/piexlmax)
func (IEDService *InternationalExchangeDatabaseService)GetInternationalExchangeDatabaseInfoList(info EngineeringEducationDatabaseReq.InternationalExchangeDatabaseSearch) (list []EngineeringEducationDatabase.InternationalExchangeDatabase, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
    // 创建db
	db := global.GVA_DB.Model(&EngineeringEducationDatabase.InternationalExchangeDatabase{})
    var IEDs []EngineeringEducationDatabase.InternationalExchangeDatabase
    // 如果有条件搜索 下方会自动创建搜索语句
    if info.StartCreatedAt !=nil && info.EndCreatedAt !=nil {
     db = db.Where("created_at BETWEEN ? AND ?", info.StartCreatedAt, info.EndCreatedAt)
    }
    if info.UniversityName != "" {
        db = db.Where("university_name LIKE ?","%"+ info.UniversityName+"%")
    }
    if info.InternationalConferencesHosted != nil {
        db = db.Where("international_conferences_hosted <> ?",info.InternationalConferencesHosted)
    }
    if info.ForeignFacultyCount != nil {
        db = db.Where("foreign_faculty_count <> ?",info.ForeignFacultyCount)
    }
    if info.InternationalStudentsCount != nil {
        db = db.Where("international_students_count <> ?",info.InternationalStudentsCount)
    }
    if info.StudentsAbroadExchangeCount != nil {
        db = db.Where("students_abroad_exchange_count <> ?",info.StudentsAbroadExchangeCount)
    }
	err = db.Count(&total).Error
	if err!=nil {
    	return
    }
        var OrderStr string
        orderMap := make(map[string]bool)
         	orderMap["international_conferences_hosted"] = true
         	orderMap["foreign_faculty_count"] = true
         	orderMap["international_students_count"] = true
         	orderMap["students_abroad_exchange_count"] = true
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
	
	err = db.Find(&IEDs).Error
	return  IEDs, total, err
}
