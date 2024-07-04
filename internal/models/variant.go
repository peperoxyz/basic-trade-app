package models

import (
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Variant struct {
	ID          uint      `gorm:"primaryKey"`
	UUID        string    `gorm:"type:varchar(36);uniqueIndex"`
	VariantName string    `gorm:"type:varchar(100);not null" form:"variant_name"`
	Quantity    int       `gorm:"not null" form:"quantity"`
	ProductID   uint      `gorm:"not null"` // remove the form tag, not needed
	CreatedAt   time.Time 
	UpdatedAt   time.Time 
}

// hook for validation
func (v *Variant) BeforeCreate(tx *gorm.DB) (err error) {
	v.UUID = uuid.New().String()
	_, errCreate := govalidator.ValidateStruct(v)

	if errCreate != nil {
		err = errCreate
		return
	}
	return
}