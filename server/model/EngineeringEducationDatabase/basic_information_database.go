// 自动生成模板BasicInfomationDatabase
package EngineeringEducationDatabase

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	
	
)

// 基础信息库 结构体  BasicInfomationDatabase
type BasicInfomationDatabase struct {
 global.GVA_MODEL 
      RegionName  string `json:"regionName" form:"regionName" gorm:"column:region_name;comment:省区市名称;" binding:"required"`  //省区市名称 
      DivisionCode  string `json:"divisionCode" form:"divisionCode" gorm:"column:division_code;comment:区划代码;" binding:"required"`  //区划代码 
      GeographicArea  string `json:"geographicArea" form:"geographicArea" gorm:"column:geographic_area;comment:地理区域;" binding:"required"`  //地理区域 
      RegionGDP  *float64 `json:"regionGDP" form:"regionGDP" gorm:"column:region_g_d_p;comment:地区生产总值;" binding:"required"`  //地区生产总值 
      SchoolName  string `json:"schoolName" form:"schoolName" gorm:"column:school_name;comment:学校名称;" binding:"required"`  //学校名称 
      SchoolCode  string `json:"schoolCode" form:"schoolCode" gorm:"primarykey;column:school_code;comment:学校代码;" binding:"required"`  //学校代码 
      SchoolType  string `json:"schoolType" form:"schoolType" gorm:"column:school_type;comment:学校类型;" binding:"required"`  //学校类型 
      SchoolCategory  string `json:"schoolCategory" form:"schoolCategory" gorm:"column:school_category;comment:学校类别;" binding:"required"`  //学校类别 
      EducationLevel  string `json:"educationLevel" form:"educationLevel" gorm:"column:education_level;comment:学校层次;" binding:"required"`  //学校层次 
      IsDoubleFirstClass  *bool `json:"isDoubleFirstClass" form:"isDoubleFirstClass" gorm:"column:is_double_first_class;comment:是否双一流;" binding:"required"`  //是否双一流 
      LocatedRegion  string `json:"locatedRegion" form:"locatedRegion" gorm:"column:located_region;comment:所在地区;" binding:"required"`  //所在地区 
      UndergraduateNumber  *int `json:"undergraduateNumber" form:"undergraduateNumber" gorm:"column:undergraduate_number;comment:在校本科生数;" binding:"required"`  //在校本科生数 
      GraduateStudentNumber  *int `json:"graduateStudentNumber" form:"graduateStudentNumber" gorm:"column:graduate_student_number;comment:在校硕士生数;" binding:"required"`  //在校硕士生数 
      DoctoralStudentNumber  *int `json:"doctoralStudentNumber" form:"doctoralStudentNumber" gorm:"column:doctoral_student_number;comment:在校博士生数;" binding:"required"`  //在校博士生数 
      InternationalStudentNumber  *int `json:"internationalStudentNumber" form:"internationalStudentNumber" gorm:"column:international_student_number;comment:在校留学生数;" binding:"required"`  //在校留学生数 
      FacultyStaffNumber  *int `json:"facultyStaffNumber" form:"facultyStaffNumber" gorm:"column:faculty_staff_number;comment:在校教职员工数;" binding:"required"`  //在校教职员工数 
      FullTimeTeacherNumber  *int `json:"fullTimeTeacherNumber" form:"fullTimeTeacherNumber" gorm:"column:full_time_teacher_number;comment:在校专任教师数;" binding:"required"`  //在校专任教师数 
      UndergraduateProgramNumber  *int `json:"undergraduateProgramNumber" form:"undergraduateProgramNumber" gorm:"column:undergraduate_program_number;comment:本科专业数;" binding:"required"`  //本科专业数 
      MastersAutorizationPointerNumber  *int `json:"mastersAutorizationPointerNumber" form:"mastersAutorizationPointerNumber" gorm:"column:masters_autorization_pointer_number;comment:硕士授权点数;" binding:"required"`  //硕士授权点数 
      DoctoralAuthorizationPointerNumber  *int `json:"doctoralAuthorizationPointerNumber" form:"doctoralAuthorizationPointerNumber" gorm:"column:doctoral_authorization_pointer_number;comment:博士授权点数;" binding:"required"`  //博士授权点数 
      MastersDegreeCategoryNumber  *int `json:"mastersDegreeCategoryNumber" form:"mastersDegreeCategoryNumber" gorm:"column:masters_degree_category_number;comment:硕士专业学位类别数;" binding:"required"`  //硕士专业学位类别数 
      DoctoralDegreeCategoryNumber  *int `json:"doctoralDegreeCategoryNumber" form:"doctoralDegreeCategoryNumber" gorm:"column:doctoral_degree_category_number;comment:博士专业学位类别数;" binding:"required"`  //博士专业学位类别数 
      PostDoctoralMobileStationNumber  *int `json:"postDoctoralMobileStationNumber" form:"postDoctoralMobileStationNumber" gorm:"column:post_doctoral_mobile_station_number;comment:博士后流动站数;" binding:"required"`  //博士后流动站数 
      SoftScienceRanking  *int `json:"softScienceRanking" form:"softScienceRanking" gorm:"column:soft_science_ranking;comment:软科排名;" binding:"required"`  //软科排名 
      CreatedBy  uint   `gorm:"column:created_by;comment:创建者"`
      UpdatedBy  uint   `gorm:"column:updated_by;comment:更新者"`
      DeletedBy  uint   `gorm:"column:deleted_by;comment:删除者"`
}


// TableName 基础信息库 BasicInfomationDatabase自定义表名 basic_infomation_database
func (BasicInfomationDatabase) TableName() string {
  return "basic_infomation_database"
}

