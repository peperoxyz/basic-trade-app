package models

import (
	"time"

	"github.com/asaskevich/govalidator"
	"gorm.io/gorm"
)

type Variant struct {
	ID          uint      `gorm:"primaryKey"`
	UUID        string    `gorm:"type:uuid;default:uuid_generate_v4()"`
	VariantName string    `gorm:"type:varchar(100);not null"`
	Quantity    int       `gorm:"not null"`
	ProductID   uint      `gorm:"not null"`
	CreatedAt   time.Time 
	UpdatedAt   time.Time 
}

// hook for validation
func (v *Variant) BeforeCreate(tx *gorm.DB) (err error) {
	_, errCreate := govalidator.ValidateStruct(v)

	if errCreate != nil {
		err = errCreate
		return
	}
	return
}