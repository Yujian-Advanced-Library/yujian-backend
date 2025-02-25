package biz

import (
	"github.com/gin-gonic/gin"
	"yujian-backend/pkg/biz/auth"

	"yujian-backend/pkg/biz/user"
)

// SetupRouter 设置路由
func SetupRouter(r *gin.Engine) {
	// 用户相关的路由
	userGroup := r.Group("/api/users")
	{
		userGroup.GET("/info", user.GetUserById())                   //信息获取
		userGroup.PUT("/update", user.UpdateUser())                  //更新
		userGroup.PUT("/password/change/:id", user.PasswordChange()) //修改密码
		userGroup.DELETE("/delete/:id", user.DeleteUser())           //删除用户
	}

	// 登录相关的路由
	r.POST("/api/user/login", auth.UserLogin())       //登录
	r.POST("/api/user/register", auth.UserRegister()) //注册

	bookGroup := r.Group("/api/books")
	{
		bookGroup.GET("/search", book.SearchBooks)    // 图书搜索
		bookGroup.GET("/:bookId", book.GetBookDetail) // 图书详情获取
	}

	//书评相关路由
	reviewsGroup := r.Group("/api/reviews")
	{
		reviewsGroup.POST("/post", book.PostReview)              //书评发布接口
		reviewsGroup.GET("/:bookId", book.GetReviews)            //书评获取接口
		reviewsGroup.POST("/:reviewId/like", book.LikeonReviews) //书评点赞/踩接口
	}

}
