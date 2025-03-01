package model

import (
	"encoding/json"
	"time"
)

// PostDTO 帖子DTO
type PostDTO struct {
	Id            int64             `json:"id"`
	Author        *UserDTO          `json:"author"`
	Title         string            `json:"title"`
	EditTime      time.Time         `json:"edit_time"`
	Category      string            `json:"category"`
	Comments      []*PostCommentDTO `json:"comments"`
	LikeUserIds   []int64           `json:"like_user_ids"`
	UnlikeUserIds []int64           `json:"unlike_user_ids"`
}

// TransformToDO 将PostDTO转换为PostDO
func (p *PostDTO) TransformToDO() *PostDO {
	likeIds, _ := json.Marshal(p.LikeUserIds)
	unlikeIds, _ := json.Marshal(p.UnlikeUserIds)
	return &PostDO{
		Id:            p.Id,
		AuthorId:      p.Author.Id,
		AuthorName:    p.Author.Name,
		Title:         p.Title,
		EditTime:      p.EditTime,
		Category:      p.Category,
		LikeUserIds:   string(likeIds),
		UnlikeUserIds: string(unlikeIds),
	}
}

// PostDO 帖子DO
type PostDO struct {
	Id            int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	AuthorId      int64     `gorm:"column:author_id" json:"author_id"`
	AuthorName    string    `gorm:"column:author_name" json:"author_name"`
	Title         string    `gorm:"column:title" json:"title"`
	Category      string    `gorm:"column:category" json:"category"`
	EditTime      time.Time `gorm:"column:edit_time" json:"edit_time"`
	LikeUserIds   string    `gorm:"like_user_ids" json:"likes"`
	UnlikeUserIds string    `gorm:"unlike_user_ids" json:"unlike_user_ids"`
}

func (p PostDO) TableName() string {
	return "post"
}

// TransformToDTO 将PostDO转换为PostDTO
func (p *PostDO) TransformToDTO(userDTO *UserDTO, comments []*PostCommentDTO) *PostDTO {
	dto := &PostDTO{
		Id:       p.Id,
		Author:   userDTO,
		Title:    p.Title,
		EditTime: p.EditTime,
		Comments: comments,
		Category: p.Category,
	}
	_ = json.Unmarshal([]byte(p.UnlikeUserIds), &dto.UnlikeUserIds)
	_ = json.Unmarshal([]byte(p.LikeUserIds), &dto.LikeUserIds)
	return dto
}

// PostEsModel 帖子ES模型
type PostEsModel struct {
	Id      string  `json:"id"`
	Title   string  `json:"title"`
	Content string  `json:"content"`
	Score   float64 `json:"score"`
}

func (p *PostEsModel) GetID() string {
	return p.Id
}

func (p *PostEsModel) SetScore(score float64) {
	p.Score = score
}

func (p *PostEsModel) GetScore() float64 {
	return p.Score
}

func (p *PostEsModel) GetIndexName() string {
	return "post"
}

func (p *PostEsModel) GetContent() string {
	return p.Content
}

func (p *PostEsModel) GetTitle() string {
	return p.Title
}

// PostCommentDTO 帖子评论DTO
type PostCommentDTO struct {
	Id             int64     `json:"id"`
	PostId         int64     `json:"post_id"`
	Author         UserDTO   `json:"author"`
	EditTime       time.Time `json:"edit_time"`
	Content        string    `json:"content"`          // 评论的内容不会很长,直接存mysql
	LikeUserIds    []int64   `json:"like_user_ids"`    // 点赞的用户ID列表
	DislikeUserIds []int64   `json:"dislike_user_ids"` // 点踩的用户ID列表
}

// PostCommentDO 帖子评论DO
type PostCommentDO struct {
	Id             int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	PostId         int64     `gorm:"column:post_id" json:"post_id"`
	AuthorId       int64     `gorm:"column:author_id" json:"author_id"`
	AuthorName     string    `gorm:"column:author_name" json:"author_name"`
	EditTime       time.Time `gorm:"column:edit_time" json:"edit_time"`
	Content        string    `gorm:"column:content" json:"content"` // 评论的内容不会很长,直接存mysql
	LikeUserIds    string    `gorm:"column:like_user_ids" json:"like_user_ids"`
	DislikeUserIds string    `gorm:"column:dislike_user_ids" json:"dislike_user_ids"`
}

func (p PostCommentDO) TableName() string {
	return "post_comment"
}

// TransformToDTO 将PostCommentDO转换为PostCommentDTO
func (p *PostCommentDO) TransformToDTO() *PostCommentDTO {
	return &PostCommentDTO{
		Id:       p.Id,
		PostId:   p.PostId,
		Author:   UserDTO{Id: p.AuthorId, Name: p.AuthorName},
		EditTime: p.EditTime,
		Content:  p.Content,
	}
}

func (p *PostCommentDTO) TransformToDO() *PostCommentDO {
	return &PostCommentDO{
		Id:         p.Id,
		PostId:     p.PostId,
		AuthorId:   p.Author.Id,
		AuthorName: p.Author.Name,
		EditTime:   p.EditTime,
		Content:    p.Content,
	}
}

// CreatePostRequestDTO 创建帖子请求DTO
type CreatePostRequestDTO struct {
	Title    string `json:"title"`
	Content  string `json:"content"`
	Category string `json:"category"`
}

// CreatePostResponseDTO 创建帖子响应DTO
type CreatePostResponseDTO struct {
	BaseResp
	PostId int64 `json:"post_id"`
}

// GetPostByTimeLineRequestDTO 获取帖子时间线请求DTO
type GetPostByTimeLineRequestDTO struct {
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	Page      int       `json:"page"`
	PageSize  int       `json:"page_size"`
	Category  string    `json:"category"`
}

// GetPostByTimeLineResponseDTO 获取帖子时间线响应DTO
type GetPostByTimeLineResponseDTO struct {
	BaseResp
	Posts []*PostDTO `json:"posts"`
	Total int64      `json:"total"`
}

// GetPostByUserIdRequestDTO 获取用户帖子请求DTO
type GetPostByUserIdRequestDTO struct {
	UserId   int64 `json:"user_id"`
	Page     int   `json:"page"`
	PageSize int   `json:"page_size"`
}

// GetPostByUserIdResponseDTO 获取用户帖子响应DTO
type GetPostByUserIdResponseDTO struct {
	BaseResp
	Posts []*PostDTO `json:"posts"`
	Total int64      `json:"total"`
}

// GetPostByIdRequestDTO 获取帖子请求DTO
type GetPostByIdRequestDTO struct {
	PostId []int64 `json:"post_id"`
}

// GetPostByIdResponseDTO 获取帖子响应DTO
type GetPostByIdResponseDTO struct {
	BaseResp
	Posts []*PostDTO `json:"posts"`
}

// GetPostContentByPostIdRequestDTO 获取帖子内容请求DTO
type GetPostContentByPostIdRequestDTO struct {
	PostId int64 `json:"post_id"`
}

// GetPostContentByPostIdResponseDTO 获取帖子内容响应DTO
type GetPostContentByPostIdResponseDTO struct {
	BaseResp
	Content string `json:"content"`
}

// CreatePostCommentReq 创建帖子评论请求
type CreatePostCommentReq struct {
	Content string
	PostId  int64
}
