package model

import "fmt"

type Condition struct {
	Fields []string
	Value  string
}

type EsQueryCondition struct {
	Conditions         []Condition
	MinimumShouldMatch int
	From               int
	Size               int
}

// EsModel 定义了一个ES模型
type EsModel interface {
	// SetScore 设置对象的评分
	SetScore(score float64)

	// GetScore 获取对象的评分
	GetScore() float64

	// GetID 获取对象的唯一标识符
	GetID() string

	// 获取对象的索引名称
	GetIndexName() string

	// 获取内容
	GetContent() string

	// 获取标题
	GetTitle() string
}

// BookInfoES 表示 Elasticsearch 中存储的图书信息
type BookInfoES struct {
	ID          int64   `json:"id"`
	Name        string  `json:"name"`
	Author      string  `json:"author"`
	CoverImage  string  `json:"cover_image"`
	Publisher   string  `json:"publisher"`
	PublishYear int     `json:"publish_year"`
	ISBN        string  `json:"isbn"`
	Score       float64 `json:"score"`
	Intro       string  `json:"intro"`
	Category    string  `json:"category"`
}

// SetScore 设置评分
func (b BookInfoES) SetScore(score float64) {
	b.Score = score
}

// GetScore 获取评分
func (b BookInfoES) GetScore() float64 {
	return b.Score
}

// GetID 获取对象的唯一标识符
func (b BookInfoES) GetID() string {
	return fmt.Sprintf("%d", b.ID)
}

// GetIndexName 获取对象的索引名称
func (b BookInfoES) GetIndexName() string {
	return "book_index"
}

// GetContent 获取内容
func (b BookInfoES) GetContent() string {
	return b.Intro //图书简介作为内容
}

// GetTitle 获取标题
func (b BookInfoES) GetTitle() string {
	return b.Name //书名作为标题
}
