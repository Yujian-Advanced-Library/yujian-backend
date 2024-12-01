package db

import (
	"gorm.io/gorm"
	"yujian-backend/pkg/model"
)

var userRepository UserRepository

type UserRepository struct {
	DB *gorm.DB
}

func GetUserRepository() *UserRepository {
	return &userRepository
}

// CreateUser 创建用户
func (r *UserRepository) CreateUser(user *model.UserDO) error {
	return r.DB.Create(user).Error
}

// GetUserById 根据ID获取用户
func (r *UserRepository) GetUserById(id int64) (*model.UserDO, error) {
	var user model.UserDO
	if err := r.DB.First(&user, id).Error; err != nil {
		return nil, err
	} else {
		return &user, nil
	}
}

// UpdateUser 更新用户信息
func (r *UserRepository) UpdateUser(user *model.UserDO) error {
	return r.DB.Save(user).Error
}

// DeleteUser 删除用户
func (r *UserRepository) DeleteUser(id int64) error {
	return r.DB.Delete(&model.UserDO{}, id).Error
}