package system

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/gofrs/uuid/v5"
)

type InviteCode struct {
	global.GVA_MODEL
	UUID   uuid.UUID `json:"uuid" gorm:"index;comment:用户UUID"`
	Code   string    `json:"code" gorm:"not null"`
	Status uint      `json:"status" gorm:"not null"`
}

func (InviteCode) TableName() string { return "invite_code" }
