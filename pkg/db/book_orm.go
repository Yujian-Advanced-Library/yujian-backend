package db

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"sort"
	"yujian-backend/pkg/es"
	"yujian-backend/pkg/log"
	"yujian-backend/pkg/model"
)

type BookRepository struct {
	DB *gorm.DB
}

var bookRepository BookRepository

func GetBookRepository() *BookRepository {
	return &bookRepository
}

// 书

// CreateBook 创建书
func (r *BookRepository) CreateBook(bookDTO *model.BookInfoDTO) (int64, error) {
	bookDO := bookDTO.TransformToDO()
	if err := r.DB.Create(bookDO).Error; err != nil {
		return 0, err
	}
	return bookDO.Id, nil
}

// GetBookById 根据ID获取书
func (r *BookRepository) GetBookById(id int64) (*model.BookInfoDTO, error) {
	var book model.BookInfoDO
	if err := r.DB.First(&book, id).Error; err != nil {
		return nil, err
	}
	return book.Transfer(), nil
}

// UpdateBook 更新书
func (r *BookRepository) UpdateBook(bookDTO *model.BookInfoDTO) error {
	bookDO := bookDTO.TransformToDO()
	return r.DB.Save(bookDO).Error
}

// DeleteBook 删除书
func (r *BookRepository) DeleteBook(id int64) error {
	return r.DB.Delete(&model.BookInfoDO{}, id).Error
}

// 书评

// CreateBookComment 创建书评
func (r *BookRepository) CreateBookComment(commentDTO *model.BookCommentDTO) error {
	commentDO := commentDTO.Transfer()
	if err := r.DB.Create(commentDO).Error; err != nil {
		return err
	}
	return nil
}

// GetBookCommentById 根据书评ID获取书评
func (r *BookRepository) GetBookCommentById(id int64) (*model.BookCommentDTO, error) {
	var comment model.BookCommentDO
	if err := r.DB.First(&comment, id).Error; err != nil {
		return nil, err
	}
	return comment.TransformToDTO(), nil
}

// GetBookCommentsByBookId 根据书ID获取书评
func (r *BookRepository) GetBookCommentsByBookId(bookId int64) ([]*model.BookCommentDTO, error) {
	var commentDOs []*model.BookCommentDO
	if err := r.DB.Where("book_id = ?", bookId).Find(&commentDOs).Error; err != nil {
		return nil, err
	}
	commentDTOs := make([]*model.BookCommentDTO, len(commentDOs))
	for i, commentDO := range commentDOs {
		commentDTOs[i] = commentDO.TransformToDTO()
	}
	return commentDTOs, nil
}

// UpdateBookComment 更新书评
func (r *BookRepository) UpdateBookComment(comment *model.BookCommentDO) error {
	return r.DB.Save(comment).Error
}

// DeleteBookComment 删除书评
func (r *BookRepository) DeleteBookComment(id int64) error {
	return r.DB.Delete(&model.BookCommentDO{}, id).Error
}

// SearchBooks 搜索书
func (r *BookRepository) SearchBooks(keyword, category string, page, pageSize int) ([]*model.BookInfoDTO, error) {
	// 调用es查询符合条件的book_id
	ctx := context.Background()
	bookIDs, err := es.SearchBooks(ctx, keyword, category, page, pageSize)
	if err != nil {
		return nil, fmt.Errorf("failed to search books in ES: %v", err)
	}
	//如查询结果为空,直接返回空列表
	if len(bookIDs) == 0 {
		return []*model.BookInfoDTO{}, nil
	}

	//根据book_id从数据库中查询图书信息
	var bookDOs []*model.BookInfoDO
	if err = r.DB.Where("id IN ?", bookIDs).Find(&bookDOs).Error; err != nil {
		log.GetLogger().Warnf("failed to search books in DB: %v", err)
		return nil, fmt.Errorf("failed to search books in DB: %v", err)
	}

	bookDTOs := make([]*model.BookInfoDTO, len(bookDOs))
	for i, bookDO := range bookDOs {
		bookDTOs[i] = bookDO.Transfer()
	}
	return bookDTOs, nil
}

func (r *BookRepository) RandomQuery(page, size int) ([]*model.BookInfoDTO, error) {
	var bookDOs []*model.BookInfoDO
	offset := (page - 1) * size
	if err := r.DB.Limit(size).Offset(offset).Find(&bookDOs).Error; err != nil {
		log.GetLogger().Warnf("failed to search books in DB: %v", err)
		return nil, err
	}
	bookDTOs := make([]*model.BookInfoDTO, len(bookDOs))
	for i, bookDO := range bookDOs {
		bookDTOs[i] = bookDO.Transfer()
	}
	return bookDTOs, nil
}

func (r *BookRepository) SearchBooksWithScore(keyword string, page, size int, fields ...string) ([]*model.BookInfoDTO, error) {
	var esConditions []model.Condition
	for _, field := range fields {
		esConditions = append(esConditions, model.Condition{
			Fields: []string{field},
			Value:  keyword,
		})
	}
	condition := model.EsQueryCondition{
		Conditions:         esConditions,
		MinimumShouldMatch: 1,
		From:               (page - 1) * size,
		Size:               size,
	}

	if books, err := es.SearchArticlesWithScores[model.BookInfoES](context.Background(), "book", condition); err != nil {
		return nil, fmt.Errorf("failed to search books in ES: %v", err)
	} else {
		sort.Sort(model.BookInfoESArr(books))
		var bookDTOs []*model.BookInfoDTO
		for _, book := range books {
			var bookDO model.BookInfoDO
			if err = r.DB.Find(&bookDO, book.ID).Error; err != nil {
				continue
			} else {
				bookDTOs = append(bookDTOs, bookDO.Transfer())
			}
		}
		return bookDTOs, nil
	}
}
