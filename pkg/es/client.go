package es

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/elastic/go-elasticsearch/v8"

	"yujian-backend/pkg/log"
	"yujian-backend/pkg/model"
)

var es *elasticsearch.Client

// ensureIndex 确保索引存在
func ensureIndex(indexName string) error {
	// 检查索引是否存在
	exists, err := es.Indices.Exists([]string{indexName})
	if err != nil {
		return fmt.Errorf("检查索引是否存在时出错: %v", err)
	}

	if exists.StatusCode == 404 {
		// 创建索引
		_, err := es.Indices.Create(indexName)
		if err != nil {
			return fmt.Errorf("创建索引失败: %v", err)
		}
	} else {
		// 确保索引是打开的
		_, err := es.Indices.Open([]string{indexName})
		if err != nil {
			return fmt.Errorf("打开索引失败: %v", err)
		}
	}

	return nil
}

// Create 创建内容
func Create(ctx context.Context, item model.EsModel) error {
	// 先确保索引存在且打开
	index := item.GetIndexName()
	if err := ensureIndex(index); err != nil {
		return fmt.Errorf("确保索引存在时出错: %v", err)
	}

	body, err := json.Marshal(item)
	if err != nil {
		return err
	}

	res, err := es.Index(
		index,
		bytes.NewReader(body),
		es.Index.WithContext(ctx),
		es.Index.WithDocumentID(item.GetID()),
	)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			return err
		}
		return fmt.Errorf("elasticsearch error: %v", e)
	}

	return nil
}

// PutSimilarityMapping 设置相似度映射
func PutSimilarityMapping(ctx context.Context, indexName string, fieldName string, similarity string) error {
	// 设置相似度映射
	body := map[string]interface{}{
		"properties": map[string]interface{}{
			fieldName: map[string]interface{}{
				"type": "text",
			},
		},
	}

	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return err
	}

	res, err := es.Indices.PutMapping(
		[]string{indexName},
		bytes.NewReader(bodyBytes),
		es.Indices.PutMapping.WithContext(ctx),
	)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("elasticsearch error: %v", res.String())
	}

	return nil
}

// Search 搜索内容
func Search[T model.EsModel](ctx context.Context, indexName string, condition model.EsQueryCondition) ([]T, error) {
	var shouldFields []interface{}
	for _, condition := range condition.Conditions {
		shouldCondition := map[string]interface{}{
			"multi_match": map[string]interface{}{
				"query":  condition.Value,
				"fields": condition.Fields,
			},
		}
		shouldFields = append(shouldFields, shouldCondition)
	}

	searchQuery := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"should":               shouldFields,
				"minimum_should_match": condition.MinimumShouldMatch, // 至少匹配一个 should 子句
			},
		},
		"from": condition.From,
		"size": condition.Size,
	}

	body, err := json.Marshal(searchQuery)
	if err != nil {
		log.GetLogger().Warnf("Error marshaling es-search query: %v", err)
		return nil, err
	}

	res, err := es.Search(
		es.Search.WithContext(ctx),
		es.Search.WithIndex(indexName),
		es.Search.WithBody(bytes.NewReader(body)),
	)
	if err != nil {
		log.GetLogger().Warnf("Error searching documents: %v", err)
		return nil, err
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, errors.New("error searching documents")
	}

	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, err
	}

	var items []T
	hits := result["hits"].(map[string]interface{})["hits"].([]interface{})

	for _, hit := range hits {
		source := hit.(map[string]interface{})["_source"]
		articleBytes, err := json.Marshal(source)
		if err != nil {
			log.GetLogger().Error("Error marshaling article: %v", err)
			continue
		}

		var item T
		if err := json.Unmarshal(articleBytes, &item); err != nil {
			log.GetLogger().Error("Error unmarshaling article: %v", err)
			continue
		}
		items = append(items, item)
	}

	return items, nil
}

// UpdateArticle 更新文章索引
func UpdateArticle(ctx context.Context, item model.EsModel) error {
	log.GetLogger().Info("Updating article with ID: %s", item.GetID())

	articleBytes, err := json.Marshal(item)
	if err != nil {
		log.GetLogger().Error("Error marshaling article: %v", err)
		return err
	}

	res, err := es.Index(
		item.GetIndexName(),
		bytes.NewReader(articleBytes),
		es.Index.WithContext(ctx),
		es.Index.WithDocumentID(item.GetID()), // 指定文档ID
		es.Index.WithRefresh("true"),
	)
	if err != nil {
		log.GetLogger().Error("Error updating article: %v", err)
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		log.GetLogger().Error("Error response from Elasticsearch: %s", res.String())
		return errors.New("error updating article")
	}

	log.GetLogger().Info("Article updated successfully: %s", item.GetID())
	return nil
}

// DeleteArticle 删除文章索引
func DeleteArticle(ctx context.Context, item model.EsModel) error {
	log.GetLogger().Info("Deleting article with ID: %s", item.GetID())

	res, err := es.Delete(
		item.GetIndexName(),
		item.GetID(),
		es.Delete.WithContext(ctx),
		es.Delete.WithRefresh("true"),
	)
	if err != nil {
		log.GetLogger().Error("Error deleting article: %v", err)
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		log.GetLogger().Error("Error response from Elasticsearch: %s", res.String())
		return errors.New("error deleting article")
	}

	log.GetLogger().Info("Article deleted successfully: %s", item.GetID())
	return nil
}

// FindSimilarArticles 根据文章ID查找相似文章
func FindSimilarArticles[T model.EsModel](ctx context.Context, indexName string, item T) ([]T, error) {
	id := item.GetID()
	log.GetLogger().Info("Finding similar articles for ID: %s", id)

	// 首先获取目标文章
	res, err := es.Get(item.GetIndexName(), id)
	if err != nil {
		log.GetLogger().Error("Error getting article: %v", err)
		return nil, err
	}
	defer res.Body.Close()

	if res.IsError() {
		log.GetLogger().Error("Article not found: %s", id)
		return nil, errors.New("article not found")
	}

	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		log.GetLogger().Error("Error decoding response: %v", err)
		return nil, err
	}

	source := result["_source"]
	articleBytes, err := json.Marshal(source)
	if err != nil {
		log.GetLogger().Error("Error marshaling source: %v", err)
		return nil, err
	}

	var article T
	if err := json.Unmarshal(articleBytes, &article); err != nil {
		log.GetLogger().Error("Error unmarshaling article: %v", err)
		return nil, err
	}

	log.GetLogger().Info("Found original article: %+v", article)

	// 构建相似度查询
	searchQuery := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": []map[string]interface{}{
					{
						"multi_match": map[string]interface{}{
							"query":  item.GetTitle() + " " + item.GetContent(),
							"fields": []string{"title^3", "content^2"},
							"type":   "best_fields",
						},
					},
				},
				"must_not": []map[string]interface{}{
					{
						"term": map[string]interface{}{
							"_id": id,
						},
					},
				},
			},
		},
	}

	body, err := json.Marshal(searchQuery)
	if err != nil {
		log.GetLogger().Error("Error marshaling search query: %v", err)
		return nil, err
	}

	log.GetLogger().Info("Search query: %s", string(body))

	searchRes, err := es.Search(
		es.Search.WithContext(ctx),
		es.Search.WithIndex(indexName),
		es.Search.WithBody(bytes.NewReader(body)),
		es.Search.WithTrackScores(true),
	)
	if err != nil {
		log.GetLogger().Error("Error searching: %v", err)
		return nil, err
	}
	defer searchRes.Body.Close()

	if searchRes.IsError() {
		log.GetLogger().Error("Search error: %s", searchRes.String())
		return nil, errors.New("error searching similar documents")
	}

	var searchResult map[string]interface{}
	if err := json.NewDecoder(searchRes.Body).Decode(&searchResult); err != nil {
		log.GetLogger().Error("Error decoding search result: %v", err)
		return nil, err
	}

	log.GetLogger().Info("Search result: %+v", searchResult)

	var items []T
	hits := searchResult["hits"].(map[string]interface{})["hits"].([]interface{})

	for _, hit := range hits {
		source := hit.(map[string]interface{})["_source"]
		score := hit.(map[string]interface{})["_score"]
		articleBytes, err := json.Marshal(source)
		if err != nil {
			log.GetLogger().Error("Error marshaling article: %v", err)
			continue
		}

		var similarItem T
		if err := json.Unmarshal(articleBytes, &similarItem); err != nil {
			log.GetLogger().Error("Error unmarshaling article: %v", err)
			continue
		}
		similarItem.SetScore(score.(float64))
		items = append(items, similarItem)
	}

	log.GetLogger().Info("Found %d similar articles", len(items))
	return items, nil
}

// SearchArticlesWithScores 搜索文章并返回相似度分数
func SearchArticlesWithScores[T model.EsModel](ctx context.Context, indexName string, condition model.EsQueryCondition) ([]T, error) {

	var matchFields []interface{}
	for _, condition := range condition.Conditions {
		shouldCondition := map[string]interface{}{
			"multi_match": map[string]interface{}{
				"query":  condition.Value,
				"fields": condition.Fields,
			},
		}
		matchFields = append(matchFields, shouldCondition)
	}

	searchQuery := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"should":               matchFields,
				"minimum_should_match": condition.MinimumShouldMatch, // 至少匹配一个 should 子句
			},
		},
	}

	body, err := json.Marshal(searchQuery)
	if err != nil {
		return nil, err
	}

	res, err := es.Search(
		es.Search.WithContext(ctx),
		es.Search.WithIndex(indexName),
		es.Search.WithBody(bytes.NewReader(body)),
		es.Search.WithTrackScores(true),
	)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("elasticsearch error: %v", e)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, err
	}

	var items []T
	hits := result["hits"].(map[string]interface{})["hits"].([]interface{})

	for _, hit := range hits {
		hitMap := hit.(map[string]interface{})
		source := hitMap["_source"]
		score := hitMap["_score"].(float64)
		log.GetLogger().Info("Id: %v, Score: %v\n", hitMap["_id"], score)

		articleBytes, err := json.Marshal(source)
		if err != nil {
			log.GetLogger().Error("Error marshaling article: %v", err)
			continue
		}

		var item T
		if err := json.Unmarshal(articleBytes, &item); err != nil {
			log.GetLogger().Error("Error unmarshalling article: %v", err)
			continue
		}

		item.SetScore(score)
		items = append(items, item)
	}

	return items, nil
}

func GetContentById(ctx context.Context, indexName, id string) (string, error) {
	// 根据documentId和indexName查询文章
	res, err := esClient.Get(indexName, id)
	if err != nil {
		log.GetLogger().Error("查询文章失败: %v", err)
		return "", err
	}
	defer res.Body.Close()

	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			log.GetLogger().Error("解析错误响应失败: %v", err)
			return "", err
		}
		log.GetLogger().Error("Elasticsearch错误: %v", e)
		return "", err
	}

	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		log.GetLogger().Error("解析响应失败: %v", err)
		return "", err
	}

	source := result["_source"]
	content, ok := source.(map[string]interface{})["content"].(string)
	if !ok {
		log.GetLogger().Error("获取content字段失败")
		return "", errors.New("获取content字段失败")
	}

	return content, nil
}
