package es

import (
	"context"
	"encoding/json"
	"fmt"
	"yujian-backend/pkg/model"
)

// SearchBooks 搜索图书
func SearchBooks(ctx context.Context, keyword, category string, page, pageSize int) ([]int64, error) {
	//构建es查询条件
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": []map[string]interface{}{},
			},
		},
		"from": (page - 1) * pageSize, // 分页起始位置
		"size": pageSize,              // 每页大小
	}
	//添加keyword查询条件
	if keyword != "" {
		keywordCondition := map[string]interface{}{
			"multi_match": map[string]interface{}{
				"query":  keyword,
				"fields": []string{"name", "author", "isbn"},
			},
		}
		mustSlice := query["query"].(map[string]interface{})["bool"].(map[string]interface{})["must"].([]map[string]interface{})
		mustSlice = append(mustSlice, keywordCondition)
		query["query"].(map[string]interface{})["bool"].(map[string]interface{})["must"] = mustSlice
	}

	//添加category查询条件
	if category != "" {
		categoryCondition := map[string]interface{}{
			"term": map[string]interface{}{
				"category": category,
			},
		}
		mustSlice := query["query"].(map[string]interface{})["bool"].(map[string]interface{})["must"].([]map[string]interface{})
		mustSlice = append(mustSlice, categoryCondition)
		query["query"].(map[string]interface{})["bool"].(map[string]interface{})["must"] = mustSlice
	}

	//将query转换为json
	queryBytes, err := json.Marshal(query)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal query: %v", err)
	}

	//调用Search函数
	esResult, err := Search[model.BookInfoES](ctx, "book_index", string(queryBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to search books in ES: %v", err)
	}

	//提取ES查询结果中的book id
	var bookIDs []int64
	for _, item := range esResult {
		bookIDs = append(bookIDs, item.ID)
	}
	return bookIDs, nil
}
