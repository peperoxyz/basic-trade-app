package models

import (
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Admin struct {
	ID        uint   `gorm:"primaryKey"`
	UUID      string `gorm:"type:varchar(36);uniqueIndex"`
	Name      string `gorm:"type:varchar(100);not null"`
	Email     string `gorm:"type:varchar(50);unique;not null"`
	Password  string `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Products  []Product `gorm:"foreignKey:AdminID"`// for association
}

// hook for validation
func (a *Admin) BeforeCreate(tx *gorm.DB) (err error) {
	a.UUID = uuid.New().String()
	_, errCreate := govalidator.ValidateStruct(a)

	if errCreate != nil {
		err = errCreate
		return
	}
	return
}