package request

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"time"
	
)

type PolicyDatabaseSearch struct{
    
        StartCreatedAt *time.Time `json:"startCreatedAt" form:"startCreatedAt"`
        EndCreatedAt   *time.Time `json:"endCreatedAt" form:"endCreatedAt"`
    
                StartPolicyReleaseDate  *time.Time  `json:"startPolicyReleaseDate" form:"startPolicyReleaseDate"`
                EndPolicyReleaseDate  *time.Time  `json:"endPolicyReleaseDate" form:"endPolicyReleaseDate"`
                      PolicyCountryOrRegion  string `json:"policyCountryOrRegion" form:"policyCountryOrRegion" `
    request.PageInfo
    Sort  string `json:"sort" form:"sort"`
    Order string `json:"order" form:"order"`
}
