package book

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"yujian-backend/pkg/db"
	"yujian-backend/pkg/model"
)


// SearchBooks 搜索书
func SearchBooks(c *gin.Context) {
	var req model.BookSearchRequest
	//绑定
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.SearchResponse{
			BaseResp: model.BaseResp{
				Error:  err,
				Code:   http.StatusBadRequest,
				ErrMsg: "invalid request parameters",
			},
			Books: nil,
		})
		return
	}
	//默认值
	if req.Page == 0 {
		req.Page = 1
	}
	if req.PageSize == 0 {
		req.PageSize = 10
	}
	bookRepository := db.GetBookRepository()
	books, err := bookRepository.SearchBooks(req.Keyword, req.Category, req.Page, req.PageSize)
	if err != nil {
		//没查到
		c.JSON(http.StatusBadRequest, model.SearchResponse{
			BaseResp: model.BaseResp{
				Error:  err,
				Code:   http.StatusNotFound,
				ErrMsg: "failed to search books",
			},
			Books: nil,
		})
		return
	}
	c.JSON(http.StatusOK, model.SearchResponse{
		BaseResp: model.BaseResp{
			Error:  nil,
			Code:   http.StatusOK,
			ErrMsg: "",
		},
		Books: books,
	})
}

// GetBookDetail 图书详情获取接口
func GetBookDetail(c *gin.Context) {
	//获取id
	bookId, err := strconv.ParseInt(c.Param("bookId"), 10, 64)
	//解析失败
	if err != nil {
		c.JSON(http.StatusBadRequest, model.BookDetailResponse{
			BaseResp: model.BaseResp{
				Error:  err,
				Code:   http.StatusBadRequest,
				ErrMsg: "invalid book ID",
			},
			Data: model.BookInfoDTO{}, //空
		})
		return
	}

	bookRepository := db.GetBookRepository()
	// 查询详情
	bookDTO, err := bookRepository.GetBookById(bookId)
	if err != nil {
		c.JSON(http.StatusNotFound, model.BookDetailResponse{
			BaseResp: model.BaseResp{
				Error:  err,
				Code:   http.StatusNotFound,
				ErrMsg: "failed to find book",
			},
			Data: model.BookInfoDTO{},
		})
		return
	}
	// 找到
	c.JSON(http.StatusOK, model.BookDetailResponse{
		BaseResp: model.BaseResp{
			Error:  nil,
			Code:   http.StatusOK,
			ErrMsg: "",
		},
		Data: *bookDTO,
	})
}
