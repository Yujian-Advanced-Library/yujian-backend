package user

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"yujian-backend/pkg/db"
	"yujian-backend/pkg/model"
)


// CreateUser 创建用户的处理函数
func CreateUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		userRepository := db.GetUserRepository()
		var userDTO model.UserDTO
		if err := c.ShouldBindJSON(&userDTO); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if id, err := userRepository.CreateUser(&userDTO); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		} else {
			c.JSON(http.StatusCreated, gin.H{"message": "User created successfully", "id": id})
		}
	}
}

// GetUserById 根据ID获取用户的处理函数
func GetUserById() gin.HandlerFunc {
	return func(c *gin.Context) {
		userRepository := db.GetUserRepository()
		id := c.Param("id")
		userId, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}

		userDO, err := userRepository.GetUserById(userId)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"user": userDO})
	}
}

// UpdateUser 更新用户的处理函数
func UpdateUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		userRepository := db.GetUserRepository()

		id := c.Param("id")
		userId, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}

		var userDTO model.UserDTO
		if err := c.ShouldBindJSON(&userDTO); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		userDO := userDTO.Transfer()
		userDO.Id = userId
		if err := userRepository.UpdateUser(userDO); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
	}
}

// DeleteUser 删除用户的处理函数
func DeleteUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		userRepository := db.GetUserRepository()

		id := c.Param("id")
		userId, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}

		if err := userRepository.DeleteUser(userId); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
	}
}



// PasswordChange 更新密码的处理函数
func PasswordChange() gin.HandlerFunc {
	return func(c *gin.Context) {
		userRepository := db.GetUserRepository()
		id := c.Param("id")
		userId, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}

		var requestBody model.ChangePasswordRequest
		if err := c.ShouldBindJSON(&requestBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		//使用UserRepository提供的GetUserById查询
		userDTO, err := userRepository.GetUserById(userId)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Failed to find user"})
			return
		} //id不存在

		//转换
		user := userDTO.Transfer()
		// 校验旧密码是否正确
		if user.Password != requestBody.OldPassword {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Incorrect password"})
			return
		}
		//密码和确认新密码是否一样
		if requestBody.NewPassword != requestBody.ConfirmPassword {
			c.JSON(http.StatusBadRequest, gin.H{"error": "New password and confirm password do not match"})
			return
		}

		// 更新密码
		if err := userRepository.PasswordChange(userId, requestBody.NewPassword); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update password"})
			return
		}
		// 返回成功响应，修改
		c.JSON(http.StatusOK, model.BaseResp{
			Error:  nil,
			Code:   0, // 这个不确定，搜了一下感觉是
			ErrMsg: "Success to update password",
		})
	}
}

