package models

import (
	"time"

	"gorm.io/gorm"
)

const (
	// 订单状态常量
	StatusPending   = 0 // 订单创建中/订单未支付
	StatusCreated   = 1 // 订单创建成功
	StatusPaid      = 2 // 订单已支付
	StatusFailed    = 3 // 订单创建失败
	StatusCancelled = 4 // 订单已取消
)

type Order struct{
	gorm.Model
	UserID 		uint 	`gorm:"not null;index"`
	ProductID 	uint 	`gorm:"not null;index"`
	ActivityID 	uint 	`gorm:"not null;index"`
	Product 	Product `gorm:"foreignKey:ProductID"`
	Activity 	Activity`gorm:"foreignKey:ActivityID"`
	OrderSn 	string 	`gorm:"not null;unique"`
	Amount 		float64 `gorm:"type:decimal(10,2); not null"`
	CreateTime 	time.Time `gorm:"not null"`
	Price     	float64 `gorm:"type:decimal(10,2); not null"`
	Quantity 	int 	`gorm:"not null"`
	Status 		int 	`gorm:"not null;default:0"`
}