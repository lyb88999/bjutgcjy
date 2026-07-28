package request

import (
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
)

type ExpertProfileSearch struct {
	StartCreatedAt *time.Time `json:"startCreatedAt" form:"startCreatedAt"`
	EndCreatedAt   *time.Time `json:"endCreatedAt" form:"endCreatedAt"`

	Name             string `json:"name" form:"name"`
	UnitName         string `json:"unitName" form:"unitName"`
	TechTitle        string `json:"techTitle" form:"techTitle"`
	DisciplineL1     string `json:"disciplineL1" form:"disciplineL1"`
	DisciplineL2     string `json:"disciplineL2" form:"disciplineL2"`
	ResearchKeywords string `json:"researchKeywords" form:"researchKeywords"`
	Status           string `json:"status" form:"status"`
	OrgId            *uint  `json:"orgId" form:"orgId"`
	OrgUnmatched     bool   `json:"orgUnmatched" form:"orgUnmatched"` // 只看 org_id 为空（单位没能匹配到"单位管理"里任何记录）的档案，供管理员清理历史数据用

	request.PageInfo
	Sort  string `json:"sort" form:"sort"`
	Order string `json:"order" form:"order"`
}
