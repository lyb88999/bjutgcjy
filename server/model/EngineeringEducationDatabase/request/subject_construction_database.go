package request

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"time"
	
)

type SubjectConstructionDatabaseSearch struct{
    
        StartCreatedAt *time.Time `json:"startCreatedAt" form:"startCreatedAt"`
        EndCreatedAt   *time.Time `json:"endCreatedAt" form:"endCreatedAt"`
    
                      PostdocStationName  string `json:"postdocStationName" form:"postdocStationName" `
                      PostdocSubjectCategory  string `json:"postdocSubjectCategory" form:"postdocSubjectCategory" `
                      MasterSubjectCode  string `json:"masterSubjectCode" form:"masterSubjectCode" `
                      MasterSubjectCategory  string `json:"masterSubjectCategory" form:"masterSubjectCategory" `
                      MasterType  string `json:"masterType" form:"masterType"`
                      MasterIsFirstClass  *bool `json:"masterIsFirstClass" form:"masterIsFirstClass" `
                      MasterUniversityName  string `json:"masterUniversityName" form:"masterUniversityName" `
                      PhdSubjectCode  string `json:"phdSubjectCode" form:"phdSubjectCode" `
                      PhdSubjectCategory  string `json:"phdSubjectCategory" form:"phdSubjectCategory" `
                      PhdType  string `json:"phdType" form:"phdType"`
                      PhdIsFirstClass  *bool `json:"phdIsFirstClass" form:"phdIsFirstClass" `
                      PhdUniversityName  string `json:"phdUniversityName" form:"phdUniversityName" `
    request.PageInfo
    Sort  string `json:"sort" form:"sort"`
    Order string `json:"order" form:"order"`
}
