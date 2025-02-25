package db

import (
	"gorm.io/gorm"
	"yujian-backend/pkg/model"
)

type BookRepository struct {
	DB *gorm.DB
}

var bookRepository BookRepository

func GetBookRepository() BookRepository {
	return bookRepository
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
func (r *BookRepository) CreateBookComment(commentDTO *model.BookCommentDTO) (int64, error) {
	commentDO := commentDTO.Transfer()
	if err := r.DB.Create(commentDO).Error; err != nil {
		return 0, err
	}
	return commentDO.Id, nil
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
//改用了es查询符合的book_id然后再用id查，但是我不确定es的部分
func (r *BookRepository) SearchBooks(keyword, category string, page, pageSize int) ([]model.BookInfoDTO, error) {
	// 调用es查询符合条件的book_id
	//es
	ctx := context.Background()
	bookIDs, err := es.SearchBooks(ctx, keyword, category, page, pageSize)
	if err != nil {
		return nil, fmt.Errorf("failed to search books in ES: %v", err)
	}
	//如查询结果为空,直接返回空列表
	if len(bookIDs) == 0 {
		return []model.BookInfoDTO{}, nil
	}

	//根据book_id从数据库中查询图书信息
	var bookDTOs []model.BookInfoDTO
	for _, id := range bookIDs {
		book, err := r.GetBookById(id)
		if err != nil {
			// 如果某本书查询失败，记录日志并跳过
			fmt.Printf("failed to get book by id %d: %v\n", id, err)
			continue
		}
		bookDTOs = append(bookDTOs, *book)
	}
	return bookDTOs, nil
}
