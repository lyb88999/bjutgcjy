// 自动生成模板SubjectConstructionDatabase
package EngineeringEducationDatabase

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	
	
)

// 学科建设库 结构体  SubjectConstructionDatabase
type SubjectConstructionDatabase struct {
 global.GVA_MODEL 
      PostdocStationName  string `json:"postdocStationName" form:"postdocStationName" gorm:"column:postdoc_station_name;comment:博士后流动站名称;"`  //流动站名称 
      PostdocSubjectCategory  string `json:"postdocSubjectCategory" form:"postdocSubjectCategory" gorm:"column:postdoc_subject_category;comment:博士后学科门类;"`  //学科门类 
      PostdocStationCount  *int `json:"postdocStationCount" form:"postdocStationCount" gorm:"column:postdoc_station_count;comment:博士后流动站数;"`  //博士后流动站数 
      MasterSubjectCode  string `json:"masterSubjectCode" form:"masterSubjectCode" gorm:"column:master_subject_code;comment:硕士学位点学科代码;"`  //学科代码 
      MasterSubjectCategory  string `json:"masterSubjectCategory" form:"masterSubjectCategory" gorm:"column:master_subject_category;comment:硕士学位点学科门类;"`  //学科门类 
      MasterSoftScienceRanking  *int `json:"masterSoftScienceRanking" form:"masterSoftScienceRanking" gorm:"column:master_soft_science_ranking;comment:硕士学位点软科排名;"`  //软科排名 
      MasterAdmissionNumber  *int `json:"masterAdmissionNumber" form:"masterAdmissionNumber" gorm:"column:master_admission_number;comment:硕士学位点招生人数;"`  //招生人数 
      MasterType  string `json:"masterType" form:"masterType" gorm:"column:master_type;type:enum('学术学位','专业学位');comment:硕士学位点类型;"`  //类型 
      MasterIsFirstClass  *bool `json:"masterIsFirstClass" form:"masterIsFirstClass" gorm:"column:master_is_first_class;comment:硕士学位点是否一流学科;"`  //是否一流学科 
      MasterUniversityName  string `json:"masterUniversityName" form:"masterUniversityName" gorm:"column:master_university_name;comment:硕士学位点所在学校;"`  //所在学校 
      PhdSubjectCode  string `json:"phdSubjectCode" form:"phdSubjectCode" gorm:"column:phd_subject_code;comment:博士学科代码;"`  //学科代码 
      PhdSubjectCategory  string `json:"phdSubjectCategory" form:"phdSubjectCategory" gorm:"column:phd_subject_category;comment:博士学科门类;"`  //学科门类 
      PhdSoftScienceRanking  *int `json:"phdSoftScienceRanking" form:"phdSoftScienceRanking" gorm:"column:phd_soft_science_ranking;comment:软科排名;"`  //软科排名 
      PhdAdmissionNumber  *int `json:"phdAdmissionNumber" form:"phdAdmissionNumber" gorm:"column:phd_admission_number;comment:博士学位点招生人数;"`  //招生人数 
      PhdType  string `json:"phdType" form:"phdType" gorm:"column:phd_type;type:enum('学术学位','专业学位');comment:博士学位点类型;"`  //类型 
      PhdIsFirstClass  *bool `json:"phdIsFirstClass" form:"phdIsFirstClass" gorm:"column:phd_is_first_class;comment:博士学位点是否一流学科;"`  //是否一流学科 
      PhdUniversityName  string `json:"phdUniversityName" form:"phdUniversityName" gorm:"column:phd_university_name;comment:博士学位点所在学校;"`  //所在学校 
}


// TableName 学科建设库 SubjectConstructionDatabase自定义表名 subject_construction_database
func (SubjectConstructionDatabase) TableName() string {
  return "subject_construction_database"
}

