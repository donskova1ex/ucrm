package client_application

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ClientApplication struct {
	gorm.Model
	UUID        *uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();not null" json:"uuid"`
	Email       string     `gorm:"not null; size:30" json:"email"`
	Company     string     `gorm:"not null; size:30" json:"company"`
	PhoneNumber string     `gorm:"not null; size:30" json:"phone_number"`
	AppText     string     `gorm:"not null; size:255" json:"app_text"`
}
