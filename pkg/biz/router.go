package biz

import (
	"github.com/gin-gonic/gin"
	"yujian-backend/pkg/biz/auth"
	"yujian-backend/pkg/biz/book"
	"yujian-backend/pkg/biz/file"
	"yujian-backend/pkg/biz/post"
	"yujian-backend/pkg/biz/recommend"
	"yujian-backend/pkg/biz/user"
)

// SetupRouter 设置路由
func SetupRouter(r *gin.Engine) {

	r.Use(auth.MiddleWareAuth())
	r.POST("/login", auth.UserLogin())       //登录
	r.POST("/register", auth.UserRegister()) //注册

	// 用户相关的路由
	userGroup := r.Group("/api/user")
	{
		userGroup.GET("/info", user.GetUserById())                   //信息获取
		userGroup.PUT("/update", user.UpdateUser())                  //更新
		userGroup.PUT("/password/change/:id", user.PasswordChange()) //修改密码
		userGroup.DELETE("/delete/:id", user.DeleteUser())           //删除用户
	}

	bookGroup := r.Group("/api/books")
	{
		bookGroup.GET("/search", book.SearchBooks())    // 图书搜索
		bookGroup.GET("/:bookId", book.GetBookDetail()) // 图书详情获取
	}

	//书评相关路由
	reviewsGroup := r.Group("/api/reviews")
	{
		reviewsGroup.POST("/post", book.CreatReview())  //书评发布接口
		reviewsGroup.GET("/:bookId", book.GetReviews()) //书评获取接口

		reviewsGroup.POST("/:reviewId/like", book.ClickLike())      //书评点赞接口
		reviewsGroup.POST("/:reviewId/dislike", book.ClickUnlike()) //书评点踩接口
	}

	posts := r.Group("/api/forum")
	{
		posts.POST("/posts/publish", post.CreatePost())
		posts.POST("/posts", post.GetPostByTimeLine())
		posts.GET("/posts/:postId/content", post.GetPostContentByPostId())
		posts.GET("/posts/:postId", post.GetPostById())
		posts.POST("/posts/:postId/comments/post", post.CreateComment())

		posts.POST("/posts/:postId/like", post.Like())
		posts.POST("/posts/:postId/like", post.DisLike())
		posts.POST("/posts/comments/:commentId/like", post.LikeComment())
		posts.POST("/posts/comments/:commentId/dislike", post.DisLikeComment())
	}

	recom := r.Group("/api/recommendation")
	{
		recom.GET("/personal", recommend.Personal())
		recom.GET("/topic", recommend.Topic())
		recom.GET("/hot", recommend.Hot())
	}

	image := r.Group("/image")
	{
		image.GET("/:imageId", file.FetchFile())
	}

}
