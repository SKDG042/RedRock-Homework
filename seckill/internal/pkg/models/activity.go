package models

import(
	"gorm.io/gorm"
	"time"
)

type Activity struct{
	gorm.Model
	Name			string		`gorm:"not null;index"`
	Product			Product
	ProductID 		uint 		`gorm:"not null;index"`
	StartTime 		time.Time 	`gorm:"not null"`
	EndTime   		time.Time 	`gorm:"not null"`
	TotalStock 		int64 		`gorm:"not null"`
	AvailableStock 	int64 		`gorm:"not null"`
	Status 			int 		`gorm:"not null"` // 0: 未开始, 1: 进行中, 2: 已结束
	SeckillPrice 	float64 	`gorm:"type:decimal(10,2); not null"`
}

// 活动是否开始
func (a *Activity) IsStarted() bool{
	return time.Now().After(a.StartTime)
}

// 活动是否结束
func (a *Activity) IsEnded() bool{
	return time.Now().After(a.EndTime)
}

// 活动是否可用
func (a *Activity) IsAvailable() bool{
	return a.Status == 1
}