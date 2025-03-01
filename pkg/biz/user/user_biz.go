package user

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"yujian-backend/pkg/db"
	"yujian-backend/pkg/model"
)

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

		var updateReq model.UpdateUserRequest
		if err := c.ShouldBindJSON(&updateReq); err != nil {
			c.JSON(http.StatusBadRequest, model.BaseResp{Error: err, Code: http.StatusBadRequest, ErrMsg: "Invalid request body"})
			return
		}

		userDTO := model.UserDTO{
			Id:       updateReq.Id,
			Email:    updateReq.Email,
			Name:     updateReq.Name,
			Password: updateReq.Password,
		}
		userDO := userDTO.Transfer()
		if err := userRepository.UpdateUser(userDO); err != nil {
			c.JSON(http.StatusInternalServerError, model.BaseResp{Error: err, Code: http.StatusInternalServerError, ErrMsg: "Failed to update user"})
			return
		}

		c.JSON(http.StatusOK, model.BaseResp{Code: http.StatusOK})
	}
}

// DeleteUser 删除用户的处理函数
func DeleteUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		userRepository := db.GetUserRepository()

		id := c.Param("id")
		userId, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, model.BaseResp{Code: http.StatusBadRequest, Error: err, ErrMsg: "Invalid user ID"})
			return
		}

		if err := userRepository.DeleteUser(userId); err != nil {
			c.JSON(http.StatusInternalServerError, model.BaseResp{
				Error:  err,
				ErrMsg: "Delete user failed",
				Code:   http.StatusInternalServerError,
			})
			return
		}

		c.JSON(http.StatusOK, model.BaseResp{
			Error:  nil,
			ErrMsg: "",
			Code:   http.StatusOK,
		})
	}
}

// PasswordChange 更新密码的处理函数
func PasswordChange() gin.HandlerFunc {
	return func(c *gin.Context) {
		userRepository := db.GetUserRepository()
		id := c.Param("id")
		userId, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, model.BaseResp{Error: err, Code: http.StatusBadRequest, ErrMsg: "Invalid user ID"})
			return
		}

		var requestBody model.ChangePasswordRequest
		if err := c.ShouldBindJSON(&requestBody); err != nil {
			c.JSON(http.StatusBadRequest, model.BaseResp{Error: err, Code: http.StatusBadRequest, ErrMsg: "Invalid request body"})
			return
		}

		//使用UserRepository提供的GetUserById查询
		userDTO, err := userRepository.GetUserById(userId)
		if err != nil {
			c.JSON(http.StatusNotFound, model.BaseResp{Error: err, Code: http.StatusNotFound, ErrMsg: "User not found"})
			return
		} //id不存在

		//转换
		user := userDTO.Transfer()
		// 校验旧密码是否正确
		if user.Password != requestBody.OldPassword {
			c.JSON(http.StatusUnauthorized, model.BaseResp{Error: errors.New("old password is incorrect"), Code: http.StatusUnauthorized, ErrMsg: "Old password is incorrect"})
			return
		}
		//密码和确认新密码是否一样
		if requestBody.NewPassword != requestBody.ConfirmPassword {
			c.JSON(http.StatusUnauthorized, model.BaseResp{Error: errors.New("new password and confirm password do not match"), Code: http.StatusBadRequest, ErrMsg: "New password and confirm password do not match"})
			return
		}

		// 更新密码
		if err := userRepository.PasswordChange(userId, requestBody.NewPassword); err != nil {
			c.JSON(http.StatusInternalServerError, model.BaseResp{Error: err, Code: http.StatusInternalServerError, ErrMsg: "Failed to update password"})
			return
		}

		// 返回成功响应，修改
		c.JSON(http.StatusOK, model.BaseResp{
			Error:  nil,
			Code:   model.Success,
			ErrMsg: "Success to update password",
		})
	}
}
