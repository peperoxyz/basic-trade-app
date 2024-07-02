package models

import (
	"time"

	"github.com/asaskevich/govalidator"
	"gorm.io/gorm"
)

type Admin struct {
	ID        uint   `gorm:"primaryKey"`
	UUID      string `gorm:"type:uuid;default:uuid_generate_v4()"`
	Name      string `gorm:"type:varchar(100);not null"`
	Email     string `gorm:"type:varchar(50);unique;not null"`
	Password  string `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Products  []Product `gorm:"foreignKey:AdminID"`// for association
}

// hook for validation
func (u *Admin) BeforeCreate(tx *gorm.DB) (err error) {
	_, errCreate := govalidator.ValidateStruct(u)

	if errCreate != nil {
		err = errCreate
		return
	}
	return
}