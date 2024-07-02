package models

import (
	"time"

	"github.com/asaskevich/govalidator"
	"gorm.io/gorm"
)

type Product struct {
	ID        uint   `gorm:"primaryKey"`
	UUID      string `gorm:"type:uuid;default:uuid_generate_v4()"`
	Name      string `gorm:"type:varchar(100);not null"`
	ImageUrl  string `gorm:"type:varchar(255)"`
	AdminID   uint   `gorm:"not null"` // otomatis foreignKey karena User(Struct di model lain)+ID(PK)
	CreatedAt time.Time
	UpdatedAt time.Time
}

// hook for validation
func (p *Product) BeforeCreate(tx *gorm.DB) (err error) {
	_, errCreate := govalidator.ValidateStruct(p)

	if errCreate != nil {
		err = errCreate
		return
	}
	return
}