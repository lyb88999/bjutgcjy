// 自动生成模板InternationalExchangeDatabase
package EngineeringEducationDatabase

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	
	
)

// 国际交流库 结构体  InternationalExchangeDatabase
type InternationalExchangeDatabase struct {
 global.GVA_MODEL 
      UniversityName  string `json:"universityName" form:"universityName" gorm:"primarykey;column:university_name;comment:学校名称;" binding:"required"`  //学校名称 
      InternationalConferencesHosted  *int `json:"internationalConferencesHosted" form:"internationalConferencesHosted" gorm:"column:international_conferences_hosted;comment:其中工科国际会议数;" binding:"required"`  //举办的国际会议数 
      ForeignFacultyCount  *int `json:"foreignFacultyCount" form:"foreignFacultyCount" gorm:"column:foreign_faculty_count;comment:其中工科外籍教师数;" binding:"required"`  //外籍教师数 
      InternationalStudentsCount  *int `json:"internationalStudentsCount" form:"internationalStudentsCount" gorm:"column:international_students_count;comment:其中工科留学生数;" binding:"required"`  //来华留学生数 
      StudentsAbroadExchangeCount  *int `json:"studentsAbroadExchangeCount" form:"studentsAbroadExchangeCount" gorm:"column:students_abroad_exchange_count;comment:其中工科学生出国交流数;" binding:"required"`  //出国交流学生数 
      CreatedBy  uint   `gorm:"column:created_by;comment:创建者"`
      UpdatedBy  uint   `gorm:"column:updated_by;comment:更新者"`
      DeletedBy  uint   `gorm:"column:deleted_by;comment:删除者"`
}


// TableName 国际交流库 InternationalExchangeDatabase自定义表名 internationalExchangeDatabase
func (InternationalExchangeDatabase) TableName() string {
  return "internationalExchangeDatabase"
}

