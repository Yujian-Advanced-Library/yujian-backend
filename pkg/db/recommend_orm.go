package db

import (
	"gorm.io/gorm"
	"yujian-backend/pkg/model"
)

type RecommendRepository struct {
	DB *gorm.DB
}

var recommendRepository RecommendRepository

func GetRecommendRepository() *RecommendRepository {
	return &recommendRepository
}

func (r *RecommendRepository) CreateRecRecord(dto *model.UserRecommendRecordDTO) error {
	return r.DB.Create(dto.Convert2DO()).Error
}

func (r *RecommendRepository) QueryByUserId(userId int64) (*model.UserRecommendRecordDTO, error) {
	var do model.UserRecommendRecordDO
	if err := recommendRepository.DB.Where("user_id = ?", userId).First(&do).Error; err != nil {
		return nil, err
	}
	return do.Convert2DTO(), nil
}

func (r *RecommendRepository) UpdateRecommend(dto *model.UserRecommendRecordDTO) error {
	return r.DB.Save(dto.Convert2DO()).Error
}
