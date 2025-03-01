package db

import (
	"sort"
	"time"
	"yujian-backend/pkg/model"

	"gorm.io/gorm"
)

var postRepository PostRepository

type PostRepository struct {
	DB *gorm.DB
}

func GetPostRepository() *PostRepository {
	return &postRepository
}

// CreatePost 创建帖子
func (r *PostRepository) CreatePost(postDTO *model.PostDTO) (int64, error) {
	postDO := postDTO.TransformToDO()
	if err := r.DB.Create(postDO).Error; err != nil {
		return 0, err
	}
	return postDO.Id, nil
}

// GetPostById 根据ID获取帖子
func (r *PostRepository) GetPostById(ids []int64) ([]*model.PostDTO, error) {
	var post []model.PostDO
	if err := r.DB.Where("id IN (?)", ids).Find(&post).Error; err != nil {
		return nil, err
	}

	postDTOs := make([]*model.PostDTO, len(post))
	for i, post := range post {
		userDTO, err := userRepository.GetUserById(post.AuthorId)
		if err != nil {
			return nil, err
		}

		comments, err := r.GetPostCommentsByPostId(post.Id)
		if err != nil {
			return nil, err
		}

		postDTOs[i] = post.TransformToDTO(userDTO, comments)
	}

	return postDTOs, nil
}

// UpdatePost 更新帖子
func (r *PostRepository) UpdatePost(postDTO *model.PostDTO) error {
	postDO := postDTO.TransformToDO()
	return r.DB.Save(postDO).Error
}

// DeletePost 删除帖子
func (r *PostRepository) DeletePost(id int64) error {
	return r.DB.Delete(&model.PostDO{}, id).Error
}

// ListPosts 获取帖子列表
func (r *PostRepository) ListPosts(offset, limit int) ([]*model.PostDTO, error) {
	var posts []model.PostDO
	if err := r.DB.Offset(offset).Limit(limit).Find(&posts).Error; err != nil {
		return nil, err
	}

	postDTOs := make([]*model.PostDTO, len(posts))
	for i, post := range posts {
		userDTO, err := userRepository.GetUserById(post.AuthorId)
		if err != nil {
			return nil, err
		}

		comments, err := r.GetPostCommentsByPostId(post.Id)
		if err != nil {
			return nil, err
		}

		postDTOs[i] = post.TransformToDTO(userDTO, comments)
	}
	return postDTOs, nil
}

// GetPostCommentsByPostId 根据帖子id获取帖子评论
func (r *PostRepository) GetPostCommentsByPostId(postId int64) ([]*model.PostCommentDTO, error) {
	var comments []model.PostCommentDO
	if err := r.DB.Where("post_id = ?", postId).Find(&comments).Error; err != nil {
		return nil, err
	}

	// 将PostCommentDO转换为PostCommentDTO
	postCommentDTOs := make([]*model.PostCommentDTO, len(comments))
	for i, comment := range comments {
		postCommentDTOs[i] = comment.TransformToDTO()
	}
	return postCommentDTOs, nil
}

// BatchGetPostCommentById 批量获取帖子评论
func (r *PostRepository) BatchGetPostCommentById(ids []int64) ([]*model.PostCommentDTO, error) {
	var comments []model.PostCommentDO
	if err := r.DB.Where("id IN (?)", ids).Find(&comments).Error; err != nil {
		return nil, err
	}

	// 将PostCommentDO转换为PostCommentDTO
	postCommentDTOs := make([]*model.PostCommentDTO, len(comments))
	for i, comment := range comments {
		postCommentDTOs[i] = comment.TransformToDTO()
	}
	return postCommentDTOs, nil
}

// GetPostByTimeLine 根据时间范围获取帖子
func (r *PostRepository) GetPostByTimeLine(startTime time.Time, endTime time.Time, category string, page int, pageSize int) ([]*model.PostDTO, int64, error) {
	var posts []*model.PostDO
	var total int64

	offset := (page - 1) * pageSize
	if err := r.DB.Model(&model.PostDO{}).Where("edit_time BETWEEN ? AND ?", startTime, endTime).
		Count(&total).Order("edit_time DESC").Offset(offset).Limit(pageSize).Find(&posts).Error; err != nil {
		return nil, 0, err
	}

	// 转换为DTO
	postDTOs := make([]*model.PostDTO, len(posts))
	for i, post := range posts {
		userDTO, err := userRepository.GetUserById(post.AuthorId)
		if err != nil {
			return nil, 0, err
		}

		comments, err := r.GetPostCommentsByPostId(post.Id)
		if err != nil {
			return nil, 0, err
		}
		postDTOs[i] = post.TransformToDTO(userDTO, comments)
	}

	// go代码里再排序一次
	sort.Slice(postDTOs, func(i, j int) bool {
		return postDTOs[i].EditTime.After(postDTOs[j].EditTime)
	})

	return postDTOs, total, nil
}

// GetPostByUserId 根据用户ID获取帖子
func (r *PostRepository) GetPostByUserId(userId int64, page int, pageSize int) ([]*model.PostDTO, int64, error) {
	var posts []*model.PostDO
	var total int64

	offset := (page - 1) * pageSize
	if err := r.DB.Model(&model.PostDO{}).Where("author_id = ?", userId).
		Count(&total).Order("edit_time DESC").Offset(offset).Limit(pageSize).Find(&posts).Error; err != nil {
		return nil, 0, err
	}

	postDTOs := make([]*model.PostDTO, len(posts))
	for i, post := range posts {
		userDTO, err := userRepository.GetUserById(post.AuthorId)
		if err != nil {
			return nil, 0, err
		}

		comments, err := r.GetPostCommentsByPostId(post.Id)
		if err != nil {
			return nil, 0, err
		}
		postDTOs[i] = post.TransformToDTO(userDTO, comments)
	}

	return postDTOs, total, nil
}

// CreatePostComment 创建帖子评论
func (r *PostRepository) CreatePostComment(user *model.UserDTO, postId int64, content string) error {
	comment := model.PostCommentDO{
		PostId:         postId,
		AuthorId:       user.Id,
		AuthorName:     user.Name,
		EditTime:       time.Now(),
		Content:        content,
		LikeUserIds:    "",
		DislikeUserIds: "",
	}
	return r.DB.Create(&comment).Error
}

func (r *PostRepository) UpdateComment(comment *model.PostCommentDTO) error {
	do := comment.TransformToDO()
	return r.DB.Save(do).Error
}
