package models

import (
	"gorm.io/gorm"
)

type Order struct{
	gorm.Model
	UserID 		uint 	`gorm:"not null;index"`
	ProductID 	uint 	`gorm:"not null;index"`
	ActivityID 	uint 	`gorm:"not null;index"`
	Product 	Product `gorm:"foreignKey:ProductID"`
	Activity 	Activity`gorm:"foreignKey:ActivityID"`
	Price     	float64 `gorm:"type:decimal(10,2); not null"`
	Quantity 	int 	`gorm:"not null"`
	Status 		int 	`gorm:"not null;default:1"` // 1:待支付，2:已支付，3:已取消
}