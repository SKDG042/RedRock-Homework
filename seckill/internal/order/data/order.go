package data

import(
	"context"
	"errors"

	"gorm.io/gorm"

	"Redrock/seckill/internal/pkg/database"
	"Redrock/seckill/internal/pkg/models"
)

type OrderData struct{
	db *gorm.DB
}

func NewOrderData() *OrderData{
	return &OrderData{
		db: database.GetDB(),
	}
}

// Create 创建订单
func (d *OrderData) Create(ctx context.Context, order *models.Order) error{
	err := d.db.WithContext(ctx).Create(order).Error
	
	return err
}

// GetByOrderSn 根据订单号获取订单 
func (d *OrderData) GetByOrderSn(ctx context.Context, orderSn string) (*models.Order, error){
	var order models.Order

	err := d.db.WithContext(ctx).Where("order_sn = ?", orderSn).First(&order).Error
	if err != nil{
		if errors.Is(err, gorm.ErrRecordNotFound){
			return nil, errors.New("订单不存在")
		}
		return nil, err
	}

	return &order, nil
}

// GetByUserIDAndOrderSn 根据用户ID和订单号获取订单详情
func (d *OrderData) GetByUserIDAndOrderSn(ctx context.Context, userID uint, orderSn string) (*models.Order, error){
	var order models.Order

	err := d.db.WithContext(ctx).Where("user_id = ? AND order_sn = ?", userID, orderSn).
		Preload("Product").
		Preload("Activity").
		First(&order).Error
	if err != nil{
		if errors.Is(err, gorm.ErrRecordNotFound){
			return nil, errors.New("订单不存在或无权限查看")
		}
		return nil, err
	}

	return &order, nil
}

// ListByUserID 根据用户ID获取订单列表
func (d *OrderData) ListByUserID(ctx context.Context, userID uint, status int) ([]*models.Order, int64, error){
	var orders []*models.Order
	var count int64

	query := d.db.WithContext(ctx).Where("user_id = ?", userID)

	if status != 0{
		query = query.Where("status = ?", status)
	}

	// 获取count
	err := query.Model(&models.Order{}).Count(&count).Error
	if err != nil{
		return nil, 0, err
	}

	err = query.Preload("Product").Preload("Activity").Order("created_at DESC").Find(&orders).Error
	if err != nil{
		return nil, 0, err
	}

	return orders, count, nil
}

// UpdateStatus 更新订单状态
func (d *OrderData) UpdateStatus(ctx context.Context, orderSn string, status int) error{
	err := d.db.WithContext(ctx).Model(&models.Order{}).Where("order_sn = ?", orderSn).Update("status", status).Error

	return err
}


