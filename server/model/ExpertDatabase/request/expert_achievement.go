package request

import (
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
)

type ExpertAchievementSearch struct {
	StartCreatedAt *time.Time `json:"startCreatedAt" form:"startCreatedAt"`
	EndCreatedAt   *time.Time `json:"endCreatedAt" form:"endCreatedAt"`

	ExpertId        *uint  `json:"expertId" form:"expertId"`
	AchievementType string `json:"achievementType" form:"achievementType"`
	Level           string `json:"level" form:"level"`
	Title           string `json:"title" form:"title"`

	request.PageInfo
	Sort  string `json:"sort" form:"sort"`
	Order string `json:"order" form:"order"`
}
