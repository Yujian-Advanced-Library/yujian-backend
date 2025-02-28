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
func (r *UserRepository) CreateUser(userDTO *model.UserDTO) (int64, error) {
	userDO := userDTO.Transfer()
	if err := r.DB.Create(userDO).Error; err != nil {
		return 0, err
	} else {
		return userDO.Id, nil
	}
}

// GetUserById 根据ID获取用户
func (r *UserRepository) GetUserById(id int64) (*model.UserDTO, error) {
	var user model.UserDO
	if err := r.DB.First(&user, id).Error; err != nil {
		return nil, err
	} else {
		return user.Transfer(), nil
	}
}

func (r *UserRepository) GetUserByName(name string) (*model.UserDTO, error) {
	var userDO model.UserDO
	if err := r.DB.Where("name IS NOT NULL").
		Where("name = ?", name).First(&userDO).Error; err != nil {
		return nil, err
	} else {
		return userDO.Transfer(), nil
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

// PasswordChange 修改密码
// 接收id和新密码，根据id在数据库中查找，没找到返回err，找到则把其密码更改为新密码
func (r *UserRepository) PasswordChange(id int64, newPassword string) error {
	var user model.UserDO
	if err := r.DB.First(&user, id).Error; err != nil {
		return err
	}
	user.Password = newPassword
	return r.DB.Save(&user).Error
}
