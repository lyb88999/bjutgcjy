package request

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"time"
	
)

type ProfessionalConstructionDatabaseSearch struct{
    
        StartCreatedAt *time.Time `json:"startCreatedAt" form:"startCreatedAt"`
        EndCreatedAt   *time.Time `json:"endCreatedAt" form:"endCreatedAt"`
    
                      MajorName  string `json:"majorName" form:"majorName" `
                      MajorCode  string `json:"majorCode" form:"majorCode" `
                      DegreeType  string `json:"degreeType" form:"degreeType"`
                      IsNewMajor  *bool `json:"isNewMajor" form:"isNewMajor" `
                      IsBroadAdmissionCategory  *bool `json:"isBroadAdmissionCategory" form:"isBroadAdmissionCategory" `
                      IsFirstClassMajor  *bool `json:"isFirstClassMajor" form:"isFirstClassMajor" `
                      IsMajorAccredited  *bool `json:"isMajorAccredited" form:"isMajorAccredited" `
                      IsSelectedForExcellenceProgram  *bool `json:"isSelectedForExcellenceProgram" form:"isSelectedForExcellenceProgram" `
                      IsDoubleDegree  *bool `json:"isDoubleDegree" form:"isDoubleDegree" `
                      IsNewEngineeringDiscipline  *bool `json:"isNewEngineeringDiscipline" form:"isNewEngineeringDiscipline" `
                      MajorSatisfaction  string `json:"majorSatisfaction" form:"majorSatisfaction"`
                      UniversityName  string `json:"universityName" form:"universityName" `
                      AdmissionMajor  string `json:"admissionMajor" form:"admissionMajor" `
    request.PageInfo
    Sort  string `json:"sort" form:"sort"`
    Order string `json:"order" form:"order"`
}
