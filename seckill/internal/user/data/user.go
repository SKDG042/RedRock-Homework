package data

import (
	"context"
	"crypto/md5"
	"errors"
	"fmt"

	"gorm.io/gorm"

	"Redrock/seckill/internal/pkg/database"
	"Redrock/seckill/internal/pkg/models"
)

// UserData 用户数据访问层
type UserData struct {
	db *gorm.DB
}

// NewUserData 创建用户数据访问对象
func NewUserData() *UserData {
	return &UserData{
		db: database.GetDB(),
	}
}

// Create 创建用户
func (d *UserData) Create(ctx context.Context, user *models.User) error {
	// 先检查用户名是否已存在
	var count int64
	err := d.db.WithContext(ctx).Model(&models.User{}).Where("username = ?", user.Username).Count(&count).Error
	if err != nil {
		return err
	}
	if count > 0 {
		return fmt.Errorf("用户名已存在")
	}

	// 加密密码
	user.Password = fmt.Sprintf("%x", md5.Sum([]byte(user.Password)))

	// 创建用户
	return d.db.WithContext(ctx).Create(user).Error
}

// CheckPassword 检查密码是否正确
func (d *UserData) CheckPassword(ctx context.Context, username, password string) (*models.User, error) {
	var user models.User
	err := d.db.WithContext(ctx).Where("username = ?", username).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("用户不存在")
		}
		return nil, err
	}

	// 验证密码
	encryptedPassword := fmt.Sprintf("%x", md5.Sum([]byte(password)))
	if user.Password != encryptedPassword {
		return nil, fmt.Errorf("密码错误")
	}

	return &user, nil
}
