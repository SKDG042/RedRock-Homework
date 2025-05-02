package data

import (
	"context"
	"time"

	"Redrock/seckill/internal/pkg/database"
	"Redrock/seckill/internal/pkg/models"

	"gorm.io/gorm"
)

type ActivityData struct {
	db *gorm.DB
}

// 创建活动实例
func NewActivityData() *ActivityData{
	return &ActivityData{
		db: database.GetDB(),
	}
}

// Create 创建活动
func (d *ActivityData) Create(ctx context.Context, activity *models.Activity) error{
	return d.db.WithContext(ctx).Create(activity).Error
}

// GetByID 通过activityID获取活动
func (d *ActivityData) GetByID(ctx context.Context, id uint) (*models.Activity, error){
	var activity models.Activity
	err := d.db.WithContext(ctx).First(&activity, id).Error
	if err != nil{
		return nil, err
	}
	return &activity, nil
}

// List 获取活动
func (d *ActivityData) List(ctx context.Context, status int) ([]*models.Activity, int64, error){
	var activities []*models.Activity
	var count int64

	query := d.db.WithContext(ctx).Model(&models.Activity{})

	// 因为我们定义当status = -1时获取所有活动
	// 因此定义一个过滤条件
	if status >= 0{
		query = query.Where("status = ?", status)
	}

	// 获取记录的数量
	if err := query.Count(&count).Error; err != nil{
		return nil, 0, err
	}

	// 查询所需记录
	err := query.Preload("Product").Order("created_at DESC").Find(&activities).Error
	if err != nil{
		return nil, 0, err
	}

	return activities, count, nil
}

// UpdateStock 更新库存
func (d *ActivityData) UpdateStock(ctx context.Context, id uint, stock int64) error{
	err := d.db.WithContext(ctx).Model(&models.Activity{}).Where("id = ?", id).Update("available_stock",stock).Error
	
	return err
}

// 接下来来处理活动的状态

// UpdataStatus 更新活动状态 
func (d *ActivityData) UpdateActivityStatus(ctx context.Context, id uint, status int) error{
	err := d.db.WithContext(ctx).Model(&models.Activity{}).Where("id = ?", id).Update("status",status).Error
	
	return err
}

// AutoUpdateActivityStatus 根据活动的时间自动更新活动状态
func (d *ActivityData) AutoUpdateActivityStatus(ctx context.Context) error{
	now := time.Now()

	// 修改截止时间已过但未结束的活动
	err := d.db.WithContext(ctx).Model(models.Activity{}).
		Where("end_time < ? AND status != ?", now, 2).
		Update("status",2).Error
	if err != nil{
		return err
	}

	// 修改开始时间未到但未开始的活动
	err = d.db.WithContext(ctx).Model(models.Activity{}).
		Where("start_time <= ? AND end_time > ? AND status = ?",now,now,0).
		Update("status", 1).Error

	return err
}
