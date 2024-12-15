package auth

import (
	"errors"
	"net/http"
	"yujian-backend/pkg/db"
	"yujian-backend/pkg/model"

	"github.com/gin-gonic/gin"
)

// UserLogin 返回一个处理用户登录的中间件函数
// 该函数验证用户身份信息，并在成功验证后返回一个令牌
func UserLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		var err error
		// 初始化用户仓库
		userRepository := db.GetUserRepository()

		// 获取请求体json
		var authInfo model.LoginRequestDTO
		if err = c.ShouldBindJSON(&authInfo); err != nil {
			// 当请求体无法被正确解析时，返回错误响应
			badBody := model.LoginResponseDTO{
				BaseResp: model.BaseResp{
					Code:  model.InvalidRequestBody,
					Error: errors.New("invalid request body"),
				},
			}
			c.JSON(http.StatusBadRequest, badBody)
			return
		}

		// 查数据库
		var userDTO *model.UserDTO
		if userDTO, err = userRepository.GetUserByName(authInfo.UserName); err != nil {
			// 当数据库中找不到指定用户名的用户时，返回错误响应
			userNotFound := model.LoginResponseDTO{
				BaseResp: model.BaseResp{
					Code:  model.UserNotExists,
					Error: errors.New("user not found"),
				},
			}
			c.JSON(http.StatusOK, userNotFound)
			return
		} else {
			// todo[xinhui] 用JWT来解决
			// 验证用户密码
			if userDTO.Password == authInfo.Password {
				// 当密码匹配时，返回包含令牌和用户信息的成功响应
				okResp := model.LoginResponseDTO{
					Token: "123",
					User:  *userDTO,
					BaseResp: model.BaseResp{
						Code:  model.Success,
						Error: nil,
					},
				}
				c.JSON(http.StatusOK, okResp)
				return
			} else {
				// 当密码不匹配时，返回错误响应
				invalidPassWord := model.LoginResponseDTO{
					BaseResp: model.BaseResp{
						Code:  model.PasswordError,
						Error: errors.New("invalid password"),
					},
				}
				c.JSON(http.StatusOK, invalidPassWord)
				return
			}
		}
	}
}

// UserRegister 返回一个处理用户注册的中间件函数
// 该函数接收用户注册信息，并在成功注册后返回一个令牌
func UserRegister() gin.HandlerFunc {
	return func(c *gin.Context) {
		var err error
		// 初始化用户仓库
		userRepository := db.GetUserRepository()

		// 获取请求体json
		var registerInfo model.RegisterRequestDTO
		if err = c.ShouldBindJSON(&registerInfo); err != nil {
			// 当请求体无法被正确解析时，返回错误响应
			badBody := model.RegisterResponseDTO{
				BaseResp: model.BaseResp{
					Code:  model.InvalidRequestBody,
					Error: errors.New("invalid request body"),
				},
			}
			c.JSON(http.StatusBadRequest, badBody)
			return
		}

		// 检查用户名是否已存在
		var existingUser *model.UserDTO
		if existingUser, err = userRepository.GetUserByName(registerInfo.UserName); err != nil {
			internalErr := model.RegisterResponseDTO{
				BaseResp: model.BaseResp{
					Code:  model.InternalError,
					Error: errors.New("internal server error"),
				},
			}
			c.JSON(http.StatusInternalServerError, internalErr)
			return
		} else if existingUser != nil {
			// 当用户名已存在时，返回错误响应
			userExists := model.RegisterResponseDTO{
				BaseResp: model.BaseResp{
					Code:  model.UserExists,
					Error: errors.New("user already exists"),
				},
			}
			c.JSON(http.StatusOK, userExists)
			return
		}

		// 创建新用户
		newUser := &model.UserDTO{
			Name:     registerInfo.UserName,
			Password: registerInfo.Password, // 注意：在实际应用中，应存储加密后的密码
		}
		if id, err := userRepository.CreateUser(newUser); err != nil {
			// 当用户创建失败时，返回错误响应
			createFailed := model.RegisterResponseDTO{
				BaseResp: model.BaseResp{
					Code:  model.InternalError,
					Error: errors.New("failed to create user"),
				},
			}
			c.JSON(http.StatusInternalServerError, createFailed)
			return
		} else {
			newUser.Id = id
		}

		// 注册成功，返回包含令牌和用户信息的成功响应
		okResp := model.RegisterResponseDTO{
			Token: "123", // todo[xinhui] 用JWT来解决
			User:  *newUser,
			BaseResp: model.BaseResp{
				Code:  model.Success,
				Error: nil,
			},
		}
		c.JSON(http.StatusOK, okResp)
	}
}
