package request

import (
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
)

type ExpertProfileSearch struct {
	StartCreatedAt *time.Time `json:"startCreatedAt" form:"startCreatedAt"`
	EndCreatedAt   *time.Time `json:"endCreatedAt" form:"endCreatedAt"`

	Name             string `json:"name" form:"name"`
	TechTitle        string `json:"techTitle" form:"techTitle"`
	DisciplineL1     string `json:"disciplineL1" form:"disciplineL1"`
	DisciplineL2     string `json:"disciplineL2" form:"disciplineL2"`
	ResearchKeywords string `json:"researchKeywords" form:"researchKeywords"`
	Status           string `json:"status" form:"status"`
	OrgId            *uint  `json:"orgId" form:"orgId"`

	request.PageInfo
	Sort  string `json:"sort" form:"sort"`
	Order string `json:"order" form:"order"`
}
