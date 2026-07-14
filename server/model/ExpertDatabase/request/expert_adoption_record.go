package request

import (
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
)

type ExpertAdoptionRecordSearch struct {
	StartCreatedAt *time.Time `json:"startCreatedAt" form:"startCreatedAt"`
	EndCreatedAt   *time.Time `json:"endCreatedAt" form:"endCreatedAt"`

	ExpertId          *uint  `json:"expertId" form:"expertId"`
	AdoptionType      string `json:"adoptionType" form:"adoptionType"`
	AdoptingUnitLevel string `json:"adoptingUnitLevel" form:"adoptingUnitLevel"`

	request.PageInfo
	Sort  string `json:"sort" form:"sort"`
	Order string `json:"order" form:"order"`
}
