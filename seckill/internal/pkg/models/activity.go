package models

import(
	"gorm.io/gorm"
	"time"
)

type Activity struct{
	gorm.Model
	ProductID 		uint 		`gorm:"not null;index"`
	StartTime 		time.Time 	`gorm:"not null"`
	EndTime   		time.Time 	`gorm:"not null"`
	Available   	bool 		`gorm:"not null,default:true"`
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
	return a.Available && time.Now().After(a.StartTime) && time.Now().Before(a.EndTime)
}