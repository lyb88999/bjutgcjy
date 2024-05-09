// 自动生成模板ProfessionalConstructionDatabase
package EngineeringEducationDatabase

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// 专业建设库 结构体  ProfessionalConstructionDatabase
type ProfessionalConstructionDatabase struct {
	global.GVA_MODEL
	MajorName                      string `json:"majorName" form:"majorName" gorm:"column:major_name;comment:专业的名称;" binding:"required"`                                                                                //专业名称
	MajorCode                      string `json:"majorCode" form:"majorCode" gorm:"primarykey;column:major_code;comment:专业代码;" binding:"required"`                                                                      //专业代码
	DegreeType                     string `json:"degreeType" form:"degreeType" gorm:"column:degree_type;type:enum('学士','硕士','博士');comment:授予的学位类型，例如学士、硕士或博士;" binding:"required"`                                      //授予学位类别
	IsNewMajor                     *bool  `json:"isNewMajor" form:"isNewMajor" gorm:"column:is_new_major;comment:表示专业是否为新设立的专业;" binding:"required"`                                                                    //是否新专业
	IsBroadAdmissionCategory       *bool  `json:"isBroadAdmissionCategory" form:"isBroadAdmissionCategory" gorm:"column:is_broad_admission_category;comment:表示是否采用宽泛类别进行招生;" binding:"required"`                        //是否大类招生
	IsFirstClassMajor              *bool  `json:"isFirstClassMajor" form:"isFirstClassMajor" gorm:"column:is_first_class_major;comment:表示该专业是否被划分为一流专业;" binding:"required"`                                            //是否一流专业
	IsMajorAccredited              *bool  `json:"isMajorAccredited" form:"isMajorAccredited" gorm:"column:is_major_accredited;comment:表示该专业是否通过了某种形式的认证;" binding:"required"`                                           //是否专业认证
	IsSelectedForExcellenceProgram *bool  `json:"isSelectedForExcellenceProgram" form:"isSelectedForExcellenceProgram" gorm:"column:is_selected_for_excellence_program;comment:表示该专业是否被选入某个卓越教育计划;" binding:"required"` //是否入选卓越计划
	IsDoubleDegree                 *bool  `json:"isDoubleDegree" form:"isDoubleDegree" gorm:"column:is_double_degree;comment:表示该专业是否提供双学位选项;" binding:"required"`                                                       //是否双学位
	IsNewEngineeringDiscipline     *bool  `json:"isNewEngineeringDiscipline" form:"isNewEngineeringDiscipline" gorm:"column:is_new_engineering_discipline;comment:表示该专业是否属于新设立的工科类别;" binding:"required"`               //是否新工科
	MajorSatisfaction              string `json:"majorSatisfaction" form:"majorSatisfaction" gorm:"column:major_satisfaction;type:enum('非常不满意','不满意','中等','满意','非常满意');comment:反映学生或者校友对专业满意度的指标;" binding:"required"`  //专业满意度
	UniversityName                 string `json:"universityName" form:"universityName" gorm:"column:university_name;comment:所在学校名称;" binding:"required"`                                                                //所在学校
	AdmissionMajor                 string `json:"admissionMajor" form:"admissionMajor" gorm:"column:admission_major;comment:招生专业;" binding:"required"`                                                                  //招生专业
	PlannedEnrollment              *int   `json:"plannedEnrollment" form:"plannedEnrollment" gorm:"column:planned_enrollment;comment:针对该专业的计划招生数量;" binding:"required"`                                                 //计划招生数
	CreatedBy                      uint   `gorm:"column:created_by;comment:创建者"`
	UpdatedBy                      uint   `gorm:"column:updated_by;comment:更新者"`
	DeletedBy                      uint   `gorm:"column:deleted_by;comment:删除者"`
}

// TableName 专业建设库 ProfessionalConstructionDatabase自定义表名 professional_construction_database
func (ProfessionalConstructionDatabase) TableName() string {
	return "professional_construction_database"
}
