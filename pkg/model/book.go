package model

import (
	"time"
)

// BookInfoDTO 书信息DTO
type BookInfoDTO struct {
	Id          int64   `json:"id"`
	Name        string  `json:"name"`
	Author      string  `json:"author"`
	CoverImage  string  `json:"cover_image"`  //封面
	Publisher   string  `json:"publisher"`    //出版社
	PublishYear int     `json:"publish_year"` //出版年份
	ISBN        string  `json:"ISBN"`
	Score       float64 `json:"score"`
	Intro       string  `json:"intro"`
	Category    string  `json:"Category"` //分类
}

// BookInfoDO 书信息数据库对象
type BookInfoDO struct {
	Id          int64   `gorm:"column:id;primaryKey" json:"id"`
	Name        string  `gorm:"column:name" json:"name"`
	Author      string  `gorm:"column:author" json:"author"`
	CoverImage  string  `json:"cover_image"`  //封面
	Publisher   string  `json:"publisher"`    //出版社
	PublishYear int     `json:"publish_year"` //出版年份
	ISBN        string  `gorm:"column:isbn" json:"ISBN"`
	Score       float64 `gorm:"column:score" json:"score"`
	Intro       string  `gorm:"column:intro" json:"intro"`
	Category    string  `json:"Category"` //分类
}

// TransformToDTO 将BookInfoDO转换为BookInfoDTO
func (bookInfoDO *BookInfoDO) Transfer() *BookInfoDTO {
	return &BookInfoDTO{
		Id:          bookInfoDO.Id,
		Name:        bookInfoDO.Name,
		Author:      bookInfoDO.Author,
		CoverImage:  bookInfoDO.CoverImage,
		Publisher:   bookInfoDO.Publisher,
		PublishYear: bookInfoDO.PublishYear,
		ISBN:        bookInfoDO.ISBN,
		Score:       bookInfoDO.Score,
		Intro:       bookInfoDO.Intro,
		Category:    bookInfoDO.Category,
	}
}

// TransformToDO 将BookInfoDTO转换为BookInfoDO
func (bookInfoDTO *BookInfoDTO) TransformToDO() *BookInfoDO {
	return &BookInfoDO{
		Id:          bookInfoDTO.Id,
		Name:        bookInfoDTO.Name,
		Author:      bookInfoDTO.Author,
		CoverImage:  bookInfoDTO.CoverImage,
		Publisher:   bookInfoDTO.Publisher,
		PublishYear: bookInfoDTO.PublishYear,
		ISBN:        bookInfoDTO.ISBN,
		Score:       bookInfoDTO.Score,
		Intro:       bookInfoDTO.Intro,
		Category:    bookInfoDTO.Category,
	}
}

// 搜索请求结构体
type BookSearchRequest struct {
	Keyword  string `json:"Keyword"`  //关键词
	Category string `json:"Category"` //分类
	Page     int    `json:"Page"`     //页码
	PageSize int    `json:"PageSize"` //页码数量
}

// 搜索返回请求结构体
type SearchResponse struct {
	BaseResp
	Books []*BookInfoDTO `json:"books"`
}

// BookDetailResponse 图书详情返回
type BookDetailResponse struct {
	BaseResp
	Data BookInfoDTO `json:"data"` // 图书详情数据
}

// BookCommentDTO 书评DTO
type BookCommentDTO struct {
	Id          int64     `json:"id"`           //书评id
	BookId      int64     `json:"book_id"`      //书的id
	PublisherId int64     `json:"publisher_id"` //发布者id
	Content     string    `json:"content"`      //书评内容
	Score       float64   `json:"score"`        //评分
	PostTime    time.Time `json:"post_time"`    //发布时间
	Like        int64     `json:"like"`         //赞数
	Dislike     int64     `json:"dislike"`      //踩数
}

// BookCommentDO 书评数据库对象
type BookCommentDO struct {
	Id          int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	BookId      int64     `gorm:"column:book_id" json:"book_id"`
	PublisherId int64     `gorm:"column:publisher_id" json:"publisher_id"`
	Content     string    `gorm:"column:content" json:"content"`
	Score       float64   `gorm:"column:score" json:"score"`
	PostTime    time.Time `gorm:"column:post_time" json:"post_time"`
	Like        int64     `gorm:"column:like" json:"like"`
	Dislike     int64     `gorm:"column:dislike" json:"dislike"`
}

// TransformToDO 将BookCommentDTO转换为BookCommentDO
func (bookCommentDTO *BookCommentDTO) Transfer() *BookCommentDO {
	return &BookCommentDO{
		Id:          bookCommentDTO.Id,
		BookId:      bookCommentDTO.BookId,
		PublisherId: bookCommentDTO.PublisherId,
		Content:     bookCommentDTO.Content,
		Score:       bookCommentDTO.Score,
		PostTime:    bookCommentDTO.PostTime,
		Like:        bookCommentDTO.Like,
		Dislike:     bookCommentDTO.Dislike,
	}
}

// TransformToDTO 将BookCommentDO转换为BookCommentDTO
func (bookCommentDO *BookCommentDO) TransformToDTO() *BookCommentDTO {
	return &BookCommentDTO{
		Id:          bookCommentDO.Id,
		BookId:      bookCommentDO.BookId,
		PublisherId: bookCommentDO.PublisherId,
		Content:     bookCommentDO.Content,
		Score:       bookCommentDO.Score,
		PostTime:    bookCommentDO.PostTime,
		Like:        bookCommentDO.Like,
		Dislike:     bookCommentDO.Dislike,
	}
}

// CreatReviewRequest 书评发布请求结构体
type CreatReviewRequest struct {
	BookId      int64   `json:"book_id"`      //图书id
	Content     string  `json:"content"`      //书评内容
	Score       float64 `json:"score"`        //评分
	PublisherId int64   `json:"publisher_id"` //111发布者id，这个接口文档里没有，但是不是需要啊，我不确定因为我看书评结构里存这个了
}

// CreatReviewResponse 书评发布返回结构体
type CreatReviewResponse struct {
	BaseResp
}

// ReviewsResponse 获取书评的返回结构体
type ReviewsResponse struct {
	BaseResp
	Reviews []BookCommentDTO `json:"book_reviews"`
}

// ClickLikeResponse 点赞/踩返回结构体
type ClickLikeResponse struct {
	BaseResp
	Like    int64 `json:"like"`
	Dislike int64 `json:"dislike"`
}
