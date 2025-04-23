package models

import(
	"gorm.io/gorm"
)

type Inventory struct{
	gorm.Model
	ActivityID uint `gorm:"uniqueIndex;not null"`
	SalesCount int  `gorm:"not null;default:0"`
	Stock      int  `gorm:"not null"`
	Version    int  `gorm:"not null;default:0"` // 乐观锁版本号，0表示未锁定
}