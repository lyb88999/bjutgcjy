package request

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"time"
	
)

type BasicInfomationDatabaseSearch struct{
    
        StartCreatedAt *time.Time `json:"startCreatedAt" form:"startCreatedAt"`
        EndCreatedAt   *time.Time `json:"endCreatedAt" form:"endCreatedAt"`
    
                      RegionName  string `json:"regionName" form:"regionName" `
                      SchoolName  string `json:"schoolName" form:"schoolName" `
                      SchoolCode  string `json:"schoolCode" form:"schoolCode" `
                      IsDoubleFirstClass  *bool `json:"isDoubleFirstClass" form:"isDoubleFirstClass" `
    request.PageInfo
    Sort  string `json:"sort" form:"sort"`
    Order string `json:"order" form:"order"`
}
