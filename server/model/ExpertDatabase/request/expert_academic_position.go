package request

import (
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
)

type ExpertAcademicPositionSearch struct {
	StartCreatedAt *time.Time `json:"startCreatedAt" form:"startCreatedAt"`
	EndCreatedAt   *time.Time `json:"endCreatedAt" form:"endCreatedAt"`

	ExpertId         *uint  `json:"expertId" form:"expertId"`
	PositionType     string `json:"positionType" form:"positionType"`
	OrganizationName string `json:"organizationName" form:"organizationName"`

	request.PageInfo
	Sort  string `json:"sort" form:"sort"`
	Order string `json:"order" form:"order"`
}
