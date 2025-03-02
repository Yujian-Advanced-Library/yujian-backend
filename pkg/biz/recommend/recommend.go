package recommend

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"yujian-backend/pkg/db"
	"yujian-backend/pkg/log"
	"yujian-backend/pkg/model"
	"yujian-backend/pkg/utils"
)

func RecordUserAction(user *model.UserDTO, bookId int64, keywords ...string) {
	if user == nil {
		return
	}
	bookRepository := db.GetBookRepository()
	book, err := bookRepository.GetBookById(bookId)
	if err != nil {
		log.GetLogger().Errorf("failed to find book: %s", err.Error())
		return
	}
	utils.AddUv(bookId)

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
		if err = repository.CreateRecRecord(rec); err != nil {
			log.GetLogger().Errorf("failed to create recommend: %s", err.Error())
		}
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

		keywords = append(keywords, book.Name)
		for _, keyword := range keywords {
			if keyword == "" {
				continue
			}
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

		if err = repository.UpdateRecommend(rec); err != nil {
			log.GetLogger().Errorf("failed to update recommend: %s", err.Error())
		}
	}
}

func Personal() gin.HandlerFunc {
	return func(c *gin.Context) {
		obj, _ := c.Get("user")
		user, _ := obj.(*model.UserDTO)

		var req model.RecommendPersonalRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, model.RecommendResponse{
				BaseResp: model.BaseResp{Code: http.StatusBadRequest, ErrMsg: "failed to bind json", Error: err},
			})
			return
		}

		repository := db.GetRecommendRepository()
		bookRepository := db.GetBookRepository()
		rec, err := repository.QueryByUserId(user.Id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, model.RecommendResponse{
				BaseResp: model.BaseResp{Code: http.StatusInternalServerError, ErrMsg: "failed to query recommend", Error: err},
			})
			return
		}
		if rec == nil {
			// fallback to normal
			if query, err := bookRepository.RandomQuery(req.Page, req.PageSize); err != nil {
				c.JSON(http.StatusInternalServerError, model.RecommendResponse{
					BaseResp: model.BaseResp{Code: http.StatusInternalServerError, Error: err, ErrMsg: "failed to query recommend"},
				})
			} else {
				c.JSON(http.StatusOK, model.RecommendResponse{
					BaseResp: model.BaseResp{Code: http.StatusOK, Error: nil},
					Books:    query,
				})
			}
		}

		var retBooks []*model.BookInfoDTO
		for _, keyword := range rec.KeyWords {
			if books, err := bookRepository.SearchBooksWithScore(keyword, req.Page, req.PageSize, "title", "content"); err == nil && len(books) > 0 {
				retBooks = append(retBooks, books...)
			}
		}
		for _, category := range rec.Category {
			if books, err := bookRepository.SearchBooksWithScore(category, req.Page, req.PageSize, "category"); err == nil && len(books) > 0 {
				retBooks = append(retBooks, books...)
			}
		}
		c.JSON(http.StatusOK, model.RecommendResponse{
			BaseResp: model.BaseResp{Code: http.StatusOK, Error: nil},
			Books:    retBooks,
		})
	}
}

func Topic() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req model.RecommendHotRequest
		err := c.ShouldBindJSON(&req)
		if err != nil {
			c.JSON(http.StatusBadRequest, model.RecommendResponse{
				BaseResp: model.BaseResp{Code: http.StatusBadRequest, ErrMsg: "failed to bind json", Error: err},
			})
			return
		}

		repository := db.GetBookRepository()
		if score, err := repository.SearchBooksWithScore(req.Topic, req.Page, req.PageSize, "category"); err != nil {
			c.JSON(http.StatusInternalServerError, model.RecommendResponse{
				BaseResp: model.BaseResp{Code: http.StatusInternalServerError, ErrMsg: "failed to query recommend", Error: err},
			})
		} else {
			c.JSON(http.StatusOK, model.RecommendResponse{
				BaseResp: model.BaseResp{Code: http.StatusOK, Error: nil},
				Books:    score,
			})
		}
	}
}

func Hot() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req model.RecommendHotRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, model.RecommendResponse{
				BaseResp: model.BaseResp{Code: http.StatusBadRequest, ErrMsg: "Failed to bind request body", Error: err},
			})
		}
		topK := utils.GetTopK(req.Page, req.PageSize)
		repository := db.GetBookRepository()
		var res []*model.BookInfoDTO
		for _, uv := range topK {
			if book, err := repository.GetBookById(uv.BookId); err == nil {
				book.Score = float64(uv.Count)
				res = append(res, book)
			}
		}
		c.JSON(http.StatusOK, model.RecommendResponse{
			BaseResp: model.BaseResp{Code: http.StatusOK, Error: nil},
			Books:    res,
		})
		return
	}
}
