package post

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"
	"yujian-backend/pkg/db"

	"github.com/gin-gonic/gin"

	"yujian-backend/pkg/es"
	"yujian-backend/pkg/log"
	"yujian-backend/pkg/model"
)

// CreatePost 创建帖子
func CreatePost() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求中获取参数
		var req model.CreatePostRequestDTO
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, model.BaseResp{
				Code:   http.StatusUnauthorized,
				Error:  errors.New("invalid request body"),
				ErrMsg: "Invalid request body",
			})
			return
		}

		obj, exists := c.Get("userName")
		if !exists {
			c.JSON(http.StatusUnauthorized, model.BaseResp{
				Code:   http.StatusUnauthorized,
				Error:  errors.New("用户未登录"),
				ErrMsg: "用户未登录",
			})
			return
		}
		userDTO := obj.(*model.UserDTO)

		resp := createPost(&req, userDTO)
		if resp.Code != model.Success {
			c.JSON(http.StatusInternalServerError, model.BaseResp{
				Code:   http.StatusInternalServerError,
				Error:  errors.New("failed to create post"),
				ErrMsg: "failed to create post",
			})
			return
		}
		c.JSON(http.StatusOK, resp)

		return
	}
}

func createPost(req *model.CreatePostRequestDTO, user *model.UserDTO) *model.CreatePostResponseDTO {
	resp := &model.CreatePostResponseDTO{
		BaseResp: model.BaseResp{
			Code: model.Success,
		},
	}

	// 参数校验
	if req.Title == "" || req.Content == "" {
		resp.Code = model.UserNotExists
		resp.Error = errors.New("标题或内容不能为空")
		return resp
	}

	// 构建帖子DO
	postDTO := &model.PostDTO{
		Title:    req.Title,
		Author:   user,
		EditTime: time.Now(),
		Category: req.Category,
		Comments: []*model.PostCommentDTO{},
	}

	// 保存帖子
	repository := db.GetPostRepository()
	id, err := repository.CreatePost(postDTO)
	if err != nil {
		log.GetLogger().Error("创建帖子失败: %v", err)
		resp.Code = model.InternalError
		resp.Error = errors.New("创建帖子失败")
		return resp
	} else {
		resp.PostId = id
	}

	// 保存到ES
	postEsModel := &model.PostEsModel{
		Id:      strconv.FormatInt(id, 10),
		Title:   req.Title,
		Content: req.Content,
	}
	err = es.Create(context.Background(), postEsModel)
	if err != nil {
		resp.Code = model.InternalError
		resp.Error = fmt.Errorf("帖子创建失败保存到ES失败: %v", err)
		// 保存到ES失败，删除帖子
		if errRollback := repository.DeletePost(id); errRollback != nil {
			resp.Code = model.InternalError
			resp.Error = fmt.Errorf("帖子创建失败保存到ES失败: %v, 回滚删除帖子失败: %v", err, errRollback)
		}
		return resp
	}

	return resp
}

// GetPostByTimeLine 根据时间线获取帖子(不会获取内容,只获取帖子信息)
func GetPostByTimeLine() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req model.GetPostByTimeLineRequestDTO
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		resp := getPostByTimeLine(&req)
		if resp.Code != model.Success {
			c.JSON(http.StatusInternalServerError, gin.H{"error": resp.Error.Error()})
			return
		}
		c.JSON(http.StatusOK, resp)
		return
	}
}

func getPostByTimeLine(req *model.GetPostByTimeLineRequestDTO) *model.GetPostByTimeLineResponseDTO {
	resp := &model.GetPostByTimeLineResponseDTO{
		BaseResp: model.BaseResp{
			Code: model.Success,
		},
	}

	if req.StartTime.IsZero() {
		req.StartTime = time.Now().Add(-24 * time.Hour)
	}
	if req.EndTime.IsZero() {
		req.EndTime = time.Now()
	}
	if req.Page < 0 {
		req.Page = 1
	}
	if req.PageSize < 0 {
		req.PageSize = 10
	}

	// 获取帖子
	repository := db.GetPostRepository()
	posts, total, err := repository.GetPostByTimeLine(req.StartTime, req.EndTime, req.Category, req.Page, req.PageSize)
	if err != nil {
		resp.Code = model.InternalError
		resp.Error = errors.New("获取帖子失败")
		return resp
	}

	resp.Posts = posts
	resp.Total = total
	return resp
}

// GetPostByUserId 根据用户ID获取帖子
func GetPostByUserId() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req model.GetPostByUserIdRequestDTO
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		resp := getPostByUserId(&req)
		if resp.Code != model.Success {
			c.JSON(http.StatusInternalServerError, gin.H{"error": resp.Error.Error()})
			return
		}
		c.JSON(http.StatusOK, resp)
		return
	}
}

func getPostByUserId(req *model.GetPostByUserIdRequestDTO) *model.GetPostByUserIdResponseDTO {
	resp := &model.GetPostByUserIdResponseDTO{
		BaseResp: model.BaseResp{
			Code: model.Success,
		},
	}

	// 参数校验
	if req.UserId <= 0 || req.Page <= 0 || req.PageSize <= 0 {
		resp.Code = model.InvalidRequestBody
		resp.Error = errors.New("参数错误")
		return resp
	}

	// 获取帖子
	repository := db.GetPostRepository()
	posts, total, err := repository.GetPostByUserId(req.UserId, req.Page, req.PageSize)
	if err != nil {
		resp.Code = model.InternalError
		resp.Error = errors.New("获取帖子失败")
		return resp
	}

	resp.Posts = posts
	resp.Total = total
	return resp
}

// GetPostById 根据帖子ID获取帖子
func GetPostById() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req model.GetPostByIdRequestDTO
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		resp := getPostById(&req)
		if resp.Code != model.Success {
			c.JSON(http.StatusInternalServerError, gin.H{"error": resp.Error.Error()})
			return
		}
		c.JSON(http.StatusOK, resp)
		return
	}
}

func getPostById(req *model.GetPostByIdRequestDTO) *model.GetPostByIdResponseDTO {
	resp := &model.GetPostByIdResponseDTO{
		BaseResp: model.BaseResp{
			Code: model.Success,
		},
	}

	// 参数校验
	if len(req.PostId) == 0 {
		resp.Code = model.InvalidRequestBody
		resp.Error = errors.New("参数错误")
		return resp
	}

	// 获取帖子
	repository := db.GetPostRepository()
	posts, err := repository.GetPostById(req.PostId)
	if err != nil {
		resp.Code = model.InternalError
		resp.Error = errors.New("获取帖子失败")
		return resp
	}

	resp.Posts = posts

	return resp
}

// GetPostContentByPostId	 获取帖子内容
func GetPostContentByPostId() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req model.GetPostContentByPostIdRequestDTO
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		resp := getPostContentByPostId(&req)
		if resp.Code != model.Success {
			c.JSON(http.StatusInternalServerError, gin.H{"error": resp.Error.Error()})
			return
		}
		c.JSON(http.StatusOK, resp)
		return
	}
}

func getPostContentByPostId(req *model.GetPostContentByPostIdRequestDTO) *model.GetPostContentByPostIdResponseDTO {
	resp := &model.GetPostContentByPostIdResponseDTO{
		BaseResp: model.BaseResp{
			Code: model.Success,
		},
	}

	// 参数校验
	if req.PostId <= 0 {
		resp.Code = model.InvalidRequestBody
		resp.Error = errors.New("参数错误")
		return resp
	}

	// 通过es获取帖子内容
	content, err := es.GetContentById(context.Background(), "post", strconv.FormatInt(req.PostId, 10))
	if err != nil {
		resp.Code = model.InternalError
		resp.Error = errors.New("获取帖子内容失败")
		return resp
	}

	resp.Content = content
	return resp
}

func Like() gin.HandlerFunc {
	return func(c *gin.Context) {
		param := c.Param("postId")
		postId, err := strconv.ParseInt(param, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, model.BaseResp{Code: http.StatusBadRequest, ErrMsg: "invalid post ID", Error: err})
			return
		}

		obj, _ := c.Get("user")
		user := obj.(*model.UserDTO)

		if err := updateLikeNum(postId, true, user.Id); err != nil {
			c.JSON(http.StatusInternalServerError, model.BaseResp{Code: http.StatusInternalServerError, ErrMsg: "update like num failed", Error: err})
			return
		}
		c.JSON(http.StatusOK, model.BaseResp{Code: http.StatusOK, ErrMsg: "success"})
	}
}

func DisLike() gin.HandlerFunc {
	return func(c *gin.Context) {
		param := c.Param("postId")
		postId, err := strconv.ParseInt(param, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, model.BaseResp{Code: http.StatusBadRequest, ErrMsg: "invalid post ID", Error: err})
			return
		}

		obj, _ := c.Get("user")
		user := obj.(*model.UserDTO)

		if err := updateLikeNum(postId, false, user.Id); err != nil {
			c.JSON(http.StatusInternalServerError, model.BaseResp{Code: http.StatusInternalServerError, ErrMsg: "update like num failed", Error: err})
			return
		}
		c.JSON(http.StatusOK, model.BaseResp{Code: http.StatusOK, ErrMsg: "success"})
	}
}

func updateLikeNum(postId int64, like bool, userId int64) error {
	repository := db.GetPostRepository()
	posts, err := repository.GetPostById([]int64{postId})
	if err != nil {
		return err
	}
	if len(posts) != 1 {
		return err
	}

	post := posts[0]
	if like {
		post.LikeUserIds = append(post.LikeUserIds, userId)
	} else {
		post.UnlikeUserIds = append(post.UnlikeUserIds, userId)
	}
	if err = repository.UpdatePost(post); err != nil {
		return err
	}
	return nil
}

func CreateComment() gin.HandlerFunc {
	return func(c *gin.Context) {
		postIdStr := c.Param("postId")
		postId, err := strconv.ParseInt(postIdStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, model.BaseResp{Code: http.StatusBadRequest, ErrMsg: "invalid post ID", Error: err})
			return
		}

		obj, exists := c.Get("user")
		if !exists {
			c.JSON(http.StatusUnauthorized, model.BaseResp{Code: http.StatusUnauthorized, ErrMsg: "unauthorized", Error: errors.New("unauthorized")})
			return
		}
		user, _ := obj.(*model.UserDTO)

		var req model.CreatePostCommentReq

		if err = c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, model.BaseResp{Code: http.StatusBadRequest, ErrMsg: "invalid request body", Error: err})
			return
		}

		repository := db.GetPostRepository()
		if err = repository.CreatePostComment(user, postId, req.Content); err != nil {
			c.JSON(http.StatusInternalServerError, model.BaseResp{Code: http.StatusInternalServerError, ErrMsg: "failed to create comment", Error: err})
			return
		} else {
			c.JSON(http.StatusOK, model.BaseResp{Code: http.StatusOK, ErrMsg: "success"})
		}
	}
}

func LikeComment() gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("commentId")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, model.BaseResp{Code: http.StatusBadRequest, ErrMsg: "invalid comment ID", Error: err})
			return
		}

		obj, _ := c.Get("user")
		user := obj.(*model.UserDTO)

		if err = updateCommentLikeNum(id, true, user.Id); err != nil {
			c.JSON(http.StatusInternalServerError, model.BaseResp{Code: http.StatusInternalServerError, ErrMsg: err.Error(), Error: err})
			return
		} else {
			c.JSON(http.StatusOK, model.BaseResp{Code: http.StatusOK, ErrMsg: "success"})
		}
	}
}

func DisLikeComment() gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("commentId")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, model.BaseResp{Code: http.StatusBadRequest, ErrMsg: "invalid comment ID", Error: err})
			return
		}

		obj, _ := c.Get("user")
		user := obj.(*model.UserDTO)

		if err = updateCommentLikeNum(id, false, user.Id); err != nil {
			c.JSON(http.StatusInternalServerError, model.BaseResp{Code: http.StatusInternalServerError, ErrMsg: err.Error(), Error: err})
			return
		} else {
			c.JSON(http.StatusOK, model.BaseResp{Code: http.StatusOK, ErrMsg: "success"})
		}
	}
}

func updateCommentLikeNum(commentId int64, like bool, userId int64) error {
	repository := db.GetPostRepository()
	comments, err := repository.BatchGetPostCommentById([]int64{commentId})
	if err != nil {
		return err
	}
	if len(comments) != 1 {
		return errors.New("comment not found")
	}

	comment := comments[0]
	if like {
		comment.LikeUserIds = append(comment.LikeUserIds, userId)
	} else {
		comment.DislikeUserIds = append(comment.DislikeUserIds, userId)
	}
	if err = repository.UpdateComment(comment); err != nil {
		return err
	}
	return nil
}
