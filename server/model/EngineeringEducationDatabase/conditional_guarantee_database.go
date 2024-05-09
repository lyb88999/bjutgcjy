// 自动生成模板ConditionalGuaranteeDatabase
package EngineeringEducationDatabase

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	
	
)

// 条件保障库 结构体  ConditionalGuaranteeDatabase
type ConditionalGuaranteeDatabase struct {
 global.GVA_MODEL 
      UniversityName  string `json:"universityName" form:"universityName" gorm:"primarykey;column:university_name;comment:学校名称;" binding:"required"`  //学校名称 
      FullTimeFacultyCount  *int `json:"fullTimeFacultyCount" form:"fullTimeFacultyCount" gorm:"column:full_time_faculty_count;comment:其中工科专任教师数;" binding:"required"`  //专任教师数 
      ForeignHighLevelTalentCount  *int `json:"foreignHighLevelTalentCount" form:"foreignHighLevelTalentCount" gorm:"column:foreign_high_level_talent_count;comment:其中工科外籍高层次人才数;" binding:"required"`  //外籍高层次人才数 
      AcademicianCount  *int `json:"academicianCount" form:"academicianCount" gorm:"column:academician_count;comment:院士人数;" binding:"required"`  //院士人数 
      NationalTalentProjectWinnerCount  *int `json:"nationalTalentProjectWinnerCount" form:"nationalTalentProjectWinnerCount" gorm:"column:national_talent_project_winner_count;comment:其中工科国家级人才项目人次;" binding:"required"`  //国家级人才项目获得者人次 
      ProvincialTalentProjectWinnerCount  *int `json:"provincialTalentProjectWinnerCount" form:"provincialTalentProjectWinnerCount" gorm:"column:provincial_talent_project_winner_count;comment:其中工科省部级人才项目人次;" binding:"required"`  //省部级人才项目获得者人次 
      OffCampusInternshipBaseCount  *int `json:"offCampusInternshipBaseCount" form:"offCampusInternshipBaseCount" gorm:"column:off_campus_internship_base_count;comment:校外实习、实践、实训基地数;" binding:"required"`  //校外实习、实践、实训基地数 
      LargeEquipmentSharingPlatformCount  *int `json:"largeEquipmentSharingPlatformCount" form:"largeEquipmentSharingPlatformCount" gorm:"column:large_equipment_sharing_platform_count;comment:大型仪器设备共享系统平台数;" binding:"required"`  //大型仪器设备共享系统平台数 
      NationalExperimentalTeachingDemoCenterCount  *int `json:"nationalExperimentalTeachingDemoCenterCount" form:"nationalExperimentalTeachingDemoCenterCount" gorm:"column:national_experimental_teaching_demo_center_count;comment:国家级实验教学示范中心数;" binding:"required"`  //国家级实验教学示范中心数 
      CreatedBy  uint   `gorm:"column:created_by;comment:创建者"`
      UpdatedBy  uint   `gorm:"column:updated_by;comment:更新者"`
      DeletedBy  uint   `gorm:"column:deleted_by;comment:删除者"`
}


// TableName 条件保障库 ConditionalGuaranteeDatabase自定义表名 conditional_guarantee_database
func (ConditionalGuaranteeDatabase) TableName() string {
  return "conditional_guarantee_database"
}

