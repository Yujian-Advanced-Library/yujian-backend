package es

import (
	"context"
	"fmt"
	"yujian-backend/pkg/model"
)

const (
	book_index = "book_index"
)

// SearchBooks 搜索图书
func SearchBooks(ctx context.Context, keyword, category string, page, pageSize int) ([]int64, error) {
	//调用Search函数
	esResult, err := Search[model.BookInfoES](ctx, book_index, model.EsQueryCondition{
		From:               page,
		Size:               pageSize,
		MinimumShouldMatch: 1,
		Conditions: []model.Condition{
			{
				Fields: []string{"name", "author", "isbn"},
				Value:  keyword,
			},
			{
				Fields: []string{"category"},
				Value:  category,
			},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to search books in ES: %v", err)
	}

	//提取ES查询结果中的book id
	var bookIds []int64
	for _, book := range esResult {
		bookIds = append(bookIds, book.ID)
	}
	return bookIds, nil
}
