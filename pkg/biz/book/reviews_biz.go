package book

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"yujian-backend/pkg/db"
	"yujian-backend/pkg/model"
)

// CreatReview 书评发布
func CreatReview(c *gin.Context) {
	//实例
	reviewsRepository := db.GetBookRepository()
	//解析请求体
	var ReviewRequest model.CreatReviewRequest
	if err := c.ShouldBindJSON(ReviewRequest); err != nil {
		//绑定失败
		c.JSON(http.StatusBadRequest, model.CreatReviewResponse{
			BaseResp: model.BaseResp{
				Error:  err,
				Code:   http.StatusBadRequest,
				ErrMsg: "invalid request parameters",
			},
		})
		return
	}
	var review model.BookCommentDTO
	review.Content = ReviewRequest.Content
	review.BookId = ReviewRequest.BookId
	review.Score = ReviewRequest.Score
	// CreateBookComment 创建书评
	if err := reviewsRepository.CreateBookComment(&review); err != nil {
		c.JSON(http.StatusInternalServerError, model.CreatReviewResponse{
			BaseResp: model.BaseResp{
				Error:  err,
				Code:   http.StatusInternalServerError,
				ErrMsg: "create review comment failed",
			},
		})
		return
	}
	c.JSON(http.StatusOK, model.CreatReviewResponse{
		BaseResp: model.BaseResp{
			Error:  nil,
			Code:   http.StatusOK,
			ErrMsg: "",
		},
	})
	return
}

// GetReviews 根据书的id获取书评
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
	return
}

// ClickLike 点赞处理函数
func ClickLike(c *gin.Context) {
	//验证登录凭证

	reviewRepository := db.GetBookRepository()
	reviewId, err := strconv.ParseInt(c.Param("reviewId"), 10, 64)
	if err != nil { //绑定失败
		c.JSON(http.StatusBadRequest, model.ClickLikeResponse{
			BaseResp: model.BaseResp{
				Error:  err,
				Code:   http.StatusBadRequest,
				ErrMsg: "invalid review id",
			},
		})
		return
	}
	//用reviewid查到具体的Reviews
	ReviewDTO, err := reviewRepository.GetBookCommentById(reviewId)
	if err != nil { // 没查到
		c.JSON(http.StatusNotFound, model.ClickLikeResponse{
			BaseResp: model.BaseResp{
				Error:  err,
				Code:   http.StatusNotFound,
				ErrMsg: "failed to find review",
			},
		})
		return
	}
	ReviewDO := ReviewDTO.Transfer()
	ReviewDO.Like++
	//111这里其实有点没搞懂怎么区分是点赞还是点踩，从地址到请求参数都没看到区分
	if err := reviewRepository.UpdateBookComment(ReviewDO); err != nil {
		c.JSON(http.StatusInternalServerError, model.ClickLikeResponse{ //修改失败
			BaseResp: model.BaseResp{
				Error:  err,
				Code:   http.StatusInternalServerError,
				ErrMsg: "failed to update like",
			},
		})
		return
	}
	//成功
	c.JSON(http.StatusOK, model.ClickLikeResponse{
		BaseResp: model.BaseResp{
			Error:  nil,
			Code:   http.StatusOK,
			ErrMsg: "",
		},
		Like:    ReviewDO.Like,
		Dislike: ReviewDO.Dislike,
	})
	return
}
