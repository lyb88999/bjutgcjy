package request

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"time"
	
)

type BasicInfomationDatabaseSearch struct{
    
        StartCreatedAt *time.Time `json:"startCreatedAt" form:"startCreatedAt"`
        EndCreatedAt   *time.Time `json:"endCreatedAt" form:"endCreatedAt"`
    
                      RegionName  string `json:"regionName" form:"regionName" `
                      DivisionCode  string `json:"divisionCode" form:"divisionCode" `
                      GeographicArea  string `json:"geographicArea" form:"geographicArea" `
                      SchoolName  string `json:"schoolName" form:"schoolName" `
                      SchoolCode  string `json:"schoolCode" form:"schoolCode" `
                      SchoolType  string `json:"schoolType" form:"schoolType" `
                      SchoolCategory  string `json:"schoolCategory" form:"schoolCategory" `
                      EducationLevel  string `json:"educationLevel" form:"educationLevel" `
                      IsDoubleFirstClass  *bool `json:"isDoubleFirstClass" form:"isDoubleFirstClass" `
                      LocatedRegion  string `json:"locatedRegion" form:"locatedRegion" `
    request.PageInfo
    Sort  string `json:"sort" form:"sort"`
    Order string `json:"order" form:"order"`
}
