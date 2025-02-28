package post

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"yujian-backend/pkg/db"
	"yujian-backend/pkg/es"
	"yujian-backend/pkg/log"
	"yujian-backend/pkg/model"
)

var postBizInstance *PostBiz

// PostBiz 帖子业务逻辑
type PostBiz struct {
	postRepo *db.PostRepository
}

// CreatePost 创建帖子
func CreatePost() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求中获取参数
		var req model.CreatePostRequestDTO
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// todo 这里外层的逻辑需要修改
		resp := postBizInstance.CreatePost(&req)
		if resp.Code != model.Success {
			c.JSON(http.StatusInternalServerError, gin.H{"error": resp.Error.Error()})
			return
		}
		c.JSON(http.StatusOK, resp)

		return
	}
}

func (b *PostBiz) CreatePost(req *model.CreatePostRequestDTO) *model.CreatePostResponseDTO {
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
		Title: req.Title,
		Author: &model.UserDTO{
			Id: req.UserId,
		},
		EditTime: time.Now(),
		Comments: []*model.PostCommentDTO{},
	}

	// 保存帖子
	id, err := b.postRepo.CreatePost(postDTO)
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
		if errRollback := b.postRepo.DeletePost(id); errRollback != nil {
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

		resp := postBizInstance.GetPostByTimeLine(&req)
		if resp.Code != model.Success {
			c.JSON(http.StatusInternalServerError, gin.H{"error": resp.Error.Error()})
			return
		}
		c.JSON(http.StatusOK, resp)
		return
	}
}

func (b *PostBiz) GetPostByTimeLine(req *model.GetPostByTimeLineRequestDTO) *model.GetPostByTimeLineResponseDTO {
	resp := &model.GetPostByTimeLineResponseDTO{
		BaseResp: model.BaseResp{
			Code: model.Success,
		},
	}

	// 参数校验
	if req.StartTime.IsZero() || req.EndTime.IsZero() || req.Page <= 0 || req.PageSize <= 0 {
		resp.Code = model.InvalidRequestBody
		resp.Error = errors.New("参数错误")
		return resp
	}

	// 获取帖子
	posts, total, err := b.postRepo.GetPostByTimeLine(req.StartTime, req.EndTime, req.Page, req.PageSize)
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

		resp := postBizInstance.GetPostByUserId(&req)
		if resp.Code != model.Success {
			c.JSON(http.StatusInternalServerError, gin.H{"error": resp.Error.Error()})
			return
		}
		c.JSON(http.StatusOK, resp)
		return
	}
}

func (b *PostBiz) GetPostByUserId(req *model.GetPostByUserIdRequestDTO) *model.GetPostByUserIdResponseDTO {
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
	posts, total, err := b.postRepo.GetPostByUserId(req.UserId, req.Page, req.PageSize)
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

		resp := postBizInstance.GetPostById(&req)
		if resp.Code != model.Success {
			c.JSON(http.StatusInternalServerError, gin.H{"error": resp.Error.Error()})
			return
		}
		c.JSON(http.StatusOK, resp)
		return
	}
}

func (b *PostBiz) GetPostById(req *model.GetPostByIdRequestDTO) *model.GetPostByIdResponseDTO {
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
	posts, err := b.postRepo.GetPostById(req.PostId)
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

		resp := postBizInstance.GetPostContentByPostId(&req)
		if resp.Code != model.Success {
			c.JSON(http.StatusInternalServerError, gin.H{"error": resp.Error.Error()})
			return
		}
		c.JSON(http.StatusOK, resp)
		return
	}
}

func (b *PostBiz) GetPostContentByPostId(req *model.GetPostContentByPostIdRequestDTO) *model.GetPostContentByPostIdResponseDTO {
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
