package models

import (
	"gorm.io/gorm"
)

type Product struct{
	gorm.Model
	Name 		string `gorm:"type:varchar(255); not null"`
	Description string `gorm:"type:text"`
	Price 		float64 `gorm:"type:decimal(10,2); not null"`
}