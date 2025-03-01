package recommand

import (
	"github.com/gin-gonic/gin"
	"yujian-backend/pkg/db"
	"yujian-backend/pkg/log"
	"yujian-backend/pkg/model"
)

func recordUserAction(c *gin.Context, bookId int64, keywords ...string) {
	bookRepository := db.GetBookRepository()

	book, err := bookRepository.GetBookById(bookId)
	if err != nil {
		log.GetLogger().Errorf("failed to find book: %s", err.Error())
		return
	}

	obj, exists := c.Get("user")
	if !exists {
		return
	}

	user, _ := obj.(*model.UserDTO)

	repository := db.GetRecommendRepository()

	rec, err := repository.QueryByUserId(user.Id)
	if err != nil {
		return
	}

	if rec == nil {
		rec = &model.UserRecommendRecordDTO{
			UserId:   user.Id,
			Category: []string{book.Category},
			KeyWords: []string{book.Name},
		}
		err = repository.CreateRecRecord(rec)
	} else {
		categoryExist := false
		for _, category := range rec.Category {
			if category == book.Category {
				categoryExist = true
				break
			}
		}
		if !categoryExist {
			rec.Category = append(rec.Category, book.Category)
		}

		for _, keyword := range keywords {
			keyWordExist := false
			for _, key := range rec.KeyWords {
				if key == keyword {
					keyWordExist = true
					break
				}
			}
			if !keyWordExist {
				rec.KeyWords = append(rec.KeyWords, keyword)
			}
		}
		err = repository.UpdateRecommend(rec)
	}
}
