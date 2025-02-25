package book

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"yujian-backend/pkg/db"
	"yujian-backend/pkg/model"
)

//根据书的id获取书评
//响应数据：指定图书的所有书评列表，包含书评内容、发布者信息、发布时间、点赞数、踩数等

// GetReviews 图书详情获取接口
func GetReviews(c *gin.Context) {
	//获取id
	bookId, err := strconv.ParseInt(c.Param("bookId"), 10, 64)
	//解析失败
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ReviewsResponse{
			BaseResp: model.BaseResp{
				Error:  err,
				Code:   http.StatusBadRequest,
				ErrMsg: "invalid book ID",
			},
			Reviews: nil, //空
		})
		return
	}
	// GetBookCommentsByBookId 根据书ID获取书评
	reviewsRepository := db.GetBookRepository()
	// 查询详情
	ReviewsDTO, err := reviewsRepository.GetBookCommentsByBookId(bookId)
	if err != nil { //没查到
		c.JSON(http.StatusNotFound, model.ReviewsResponse{
			BaseResp: model.BaseResp{
				Error:  err,
				Code:   http.StatusNotFound,
				ErrMsg: "failed to find reviews",
			},
			Reviews: nil,
		})
		return
	}
	// 找到
	reviews := make([]model.BookCommentDTO, len(ReviewsDTO))
	for i, review := range ReviewsDTO {
		reviews[i] = *review //解引用指针
	}
	c.JSON(http.StatusOK, model.ReviewsResponse{
		BaseResp: model.BaseResp{
			Error:  nil,
			Code:   http.StatusOK,
			ErrMsg: "",
		},
		Reviews: reviews,
	})
}
