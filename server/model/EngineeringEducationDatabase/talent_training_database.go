// 自动生成模板TalentTrainingDatabase
package EngineeringEducationDatabase

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	
	
)

// 人才培养库 结构体  TalentTrainingDatabase
type TalentTrainingDatabase struct {
 global.GVA_MODEL 
      UniversityName  string `json:"universityName" form:"universityName" gorm:"primarykey;column:university_name;comment:学校名称;" binding:"required"`  //所在学校 
      UndergraduateCount  *int `json:"undergraduateCount" form:"undergraduateCount" gorm:"column:undergraduate_count;comment:其中工科人数;" binding:"required"`  //本科生数 
      MasterStudentCount  *int `json:"masterStudentCount" form:"masterStudentCount" gorm:"column:master_student_count;comment:其中工科人数;" binding:"required"`  //硕士生数 
      DoctoralStudentCount  *int `json:"doctoralStudentCount" form:"doctoralStudentCount" gorm:"column:doctoral_student_count;comment:其中工科人数;" binding:"required"`  //博士生数 
      NationalFirstClassCourseCount  *int `json:"nationalFirstClassCourseCount" form:"nationalFirstClassCourseCount" gorm:"column:national_first_class_course_count;comment:国家级一流课程数;" binding:"required"`  //国家级一流课程数 
      ProvincialFirstClassCourseCount  *int `json:"provincialFirstClassCourseCount" form:"provincialFirstClassCourseCount" gorm:"column:provincial_first_class_course_count;comment:省部级一流课程数;" binding:"required"`  //省部级一流课程数 
      AtionalTeachingDemonstrationCenterCount  *int `json:"ationalTeachingDemonstrationCenterCount" form:"ationalTeachingDemonstrationCenterCount" gorm:"column:ational_teaching_demonstration_center_count;comment:国家级教学示范中心数量;" binding:"required"`  //国家级教学示范中心 
      ProvincialTeachingDemonstrationCenterCount  *int `json:"provincialTeachingDemonstrationCenterCount" form:"provincialTeachingDemonstrationCenterCount" gorm:"column:provincial_teaching_demonstration_center_count;comment:省部级教学示范中心数量;" binding:"required"`  //省部级教学示范中心 
      NationalPlanningTextbookCount  *int `json:"nationalPlanningTextbookCount" form:"nationalPlanningTextbookCount" gorm:"column:national_planning_textbook_count;comment:其中工科教材数;" binding:"required"`  //国家级规划教材数 
      NationalExcellentTextbookCount  *int `json:"nationalExcellentTextbookCount" form:"nationalExcellentTextbookCount" gorm:"column:national_excellent_textbook_count;comment:其中工科教材数;" binding:"required"`  //全国优秀教材数 
      ProvincialExcellentTextbookCount  *int `json:"provincialExcellentTextbookCount" form:"provincialExcellentTextbookCount" gorm:"column:provincial_excellent_textbook_count;comment:其中工科教材数;" binding:"required"`  //省级优秀教材数 
      InternationalAwardsCount  *int `json:"internationalAwardsCount" form:"internationalAwardsCount" gorm:"column:international_awards_count;comment:其中工科学生获奖数;" binding:"required"`  //学生获国际奖项数 
      NationalAwardsCount  *int `json:"nationalAwardsCount" form:"nationalAwardsCount" gorm:"column:national_awards_count;comment:其中工科学生获奖数;" binding:"required"`  //学生获国家级奖项数 
      ProvincialAwardsCount  *int `json:"provincialAwardsCount" form:"provincialAwardsCount" gorm:"column:provincial_awards_count;comment:其中工科学生获奖数;" binding:"required"`  //学生获省部级奖项数 
      UndergraduateEmploymentRate  *float64 `json:"undergraduateEmploymentRate" form:"undergraduateEmploymentRate" gorm:"column:undergraduate_employment_rate;comment:其中工科学生平均落实率;" binding:"required"`  //本科生平均毕业去向落实率 
      MasterEmploymentRate  *float64 `json:"masterEmploymentRate" form:"masterEmploymentRate" gorm:"column:master_employment_rate;comment:其中工科学生平均落实率;" binding:"required"`  //硕士生平均毕业去向落实率 
      DoctoralEmploymentRate  *float64 `json:"doctoralEmploymentRate" form:"doctoralEmploymentRate" gorm:"column:doctoral_employment_rate;comment:其中工科学生平均落实率;" binding:"required"`  //博士生平均毕业去向落实率 
      ProjectCount  *int `json:"projectCount" form:"projectCount" gorm:"column:project_count;comment:其中工科类项目数;" binding:"required"`  //项目数 
      ParticipatingStudentCount  *int `json:"participatingStudentCount" form:"participatingStudentCount" gorm:"column:participating_student_count;comment:其中工科学生参与数;" binding:"required"`  //参与学生数 
      TransferredIntoEngineeringCount  *int `json:"transferredIntoEngineeringCount" form:"transferredIntoEngineeringCount" gorm:"column:transferred_into_engineering_count;comment:转入工科学生数;" binding:"required"`  //转入工科学生数 
      TransferredOutOfEngineeringCount  *int `json:"transferredOutOfEngineeringCount" form:"transferredOutOfEngineeringCount" gorm:"column:transferred_out_of_engineering_count;comment:转出工科学生数;" binding:"required"`  //转出工科学生数 
      CreatedBy  uint   `gorm:"column:created_by;comment:创建者"`
      UpdatedBy  uint   `gorm:"column:updated_by;comment:更新者"`
      DeletedBy  uint   `gorm:"column:deleted_by;comment:删除者"`
}


// TableName 人才培养库 TalentTrainingDatabase自定义表名 talent_training_database
func (TalentTrainingDatabase) TableName() string {
  return "talent_training_database"
}

