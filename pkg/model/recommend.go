package model

import "encoding/json"

type UserRecommendRecordDO struct {
	Id       int64  `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UserId   int64  `gorm:"column:user_id" json:"user_id"`
	Category string `gorm:"column:category" json:"category"`
	KeyWords string `gorm:"column:key_words" json:"key_words"`
}

type UserRecommendRecordDTO struct {
	Id       int64    `json:"id"`
	UserId   int64    `json:"user_id"`
	Category []string `json:"category"`
	KeyWords []string `json:"key_words"`
}

func (u UserRecommendRecordDO) TableName() string {
	return "user_recommend"
}

func (u *UserRecommendRecordDTO) Convert2DO() *UserRecommendRecordDO {
	categoryStr, _ := json.Marshal(u.Category)
	keyWordsStr, _ := json.Marshal(u.KeyWords)
	return &UserRecommendRecordDO{
		Id:       u.Id,
		UserId:   u.UserId,
		Category: string(categoryStr),
		KeyWords: string(keyWordsStr),
	}
}

func (u *UserRecommendRecordDO) Convert2DTO() *UserRecommendRecordDTO {
	dto := &UserRecommendRecordDTO{
		Id:     u.Id,
		UserId: u.UserId,
	}
	_ = json.Unmarshal([]byte(u.Category), &dto.Category)
	_ = json.Unmarshal([]byte(u.KeyWords), &dto.KeyWords)
	return dto
}

type BaseRecommendRequest struct {
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
}

type RecommendResponse struct {
	BaseResp
	Books []*BookInfoDTO `json:"books"`
}

type RecommendPersonalRequest struct {
	BaseRecommendRequest
}

type RecommendHotRequest struct {
	Topic string `json:"topic"`
	BaseRecommendRequest
}
