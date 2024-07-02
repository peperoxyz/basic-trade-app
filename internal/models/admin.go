package models

import (
	"basic-trade-app/internal/helpers"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Admin struct {
	ID        uint   `gorm:"primaryKey" `
	UUID      string `gorm:"type:varchar(36);uniqueIndex"`
	Name      string `gorm:"type:varchar(100);not null" json:"name"`
	Email     string `gorm:"type:varchar(50);unique;not null" json:"email"`
	Password  string `gorm:"not null" json:"password"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Products  []Product `gorm:"foreignKey:AdminID"`// for association
}

// hook for validation
func (a *Admin) BeforeCreate(tx *gorm.DB) (err error) {
	_, errCreate := govalidator.ValidateStruct(a)

	if errCreate != nil {
		err = errCreate
		return
	}

	// set uuid
	a.UUID = uuid.New().String()

	// hash password
	a.Password = helpers.HashPass(a.Password)

	err = nil
	return
}