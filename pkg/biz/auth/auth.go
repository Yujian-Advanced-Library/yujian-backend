package auth

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"yujian-backend/pkg/db"
	"yujian-backend/pkg/model"
)


// jwt密钥
var jwtKey = []byte("your_secret_key")

// Claims 定义jwt的claims结构
type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

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
			c.JSON(http.StatusBadRequest, model.LoginResponseDTO{
				BaseResp: model.BaseResp{
					Error:  err,
					Code:   http.StatusBadRequest,
					ErrMsg: "invalid request body",
				},
			})
			return
		}
		// 查数据库
		var userDTO *model.UserDTO
		if userDTO, err = userRepository.GetUserByName(authInfo.UserName); err != nil {
			// 当数据库中找不到指定用户名的用户时，返回错误响应
			c.JSON(http.StatusUnauthorized, model.LoginResponseDTO{
				BaseResp: model.BaseResp{
					Error:  err,
					Code:   http.StatusUnauthorized,
					ErrMsg: "user not found",
				},
			})
			return
		} else {
			// todo 用JWT来解决
			// 验证用户密码
			if userDTO.Password == authInfo.Password {
				// 当密码匹配时，返回包含令牌和用户信息的成功响应
				expirationTime := time.Now().Add(24 * time.Hour) //Token有效期24h
				claims := &Claims{
					Username: authInfo.UserName,
					RegisteredClaims: jwt.RegisteredClaims{
						ExpiresAt: jwt.NewNumericDate(expirationTime),
					},
				}
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
				tokenString, err := token.SignedString(jwtKey)
				if err != nil {
					// 如果生成 Token 失败，返回错误响应
					c.JSON(http.StatusInternalServerError, model.LoginResponseDTO{
						BaseResp: model.BaseResp{
							Error:  err,
							Code:   http.StatusInternalServerError,
							ErrMsg: "failed to generate token",
						},
					})
					return
				} //成功
				okResp := model.LoginResponseDTO{
					Token: tokenString,
					User:  *userDTO,
					BaseResp: model.BaseResp{
						Error:  nil,
						Code:   http.StatusOK,
						ErrMsg: "",
					},
				}
				c.JSON(http.StatusOK, okResp)
				return

			} else {
				// 当密码不匹配时，返回错误响应
				invalidPassWord := model.LoginResponseDTO{
					BaseResp: model.BaseResp{
						Error:  nil,
						Code:   http.StatusUnauthorized,
						ErrMsg: "invalid password",
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
					Error:  err,
					Code:   http.StatusBadRequest,
					ErrMsg: "Invalid request body",
				},
			}
			c.JSON(http.StatusBadRequest, badBody)
			return
		}

		// 检验密码与确认密码是否相同
		if registerInfo.Password != registerInfo.ConfirmPassword {
			passwordNotMatch := model.RegisterResponseDTO{
				BaseResp: model.BaseResp{
					Error:  err,
					Code:   http.StatusBadRequest,
					ErrMsg: "Password and confirm password do not match",
				},
			}
			c.JSON(http.StatusBadRequest, passwordNotMatch)
			return
		}

		// 检查用户名是否已存在
		var existingUser *model.UserDTO
		if existingUser, err = userRepository.GetUserByName(registerInfo.UserName); err != nil {
			internalErr := model.RegisterResponseDTO{
				BaseResp: model.BaseResp{
					Error:  err,
					Code:   http.StatusInternalServerError,
					ErrMsg: "Internal server error",
				},
			}
			c.JSON(http.StatusInternalServerError, internalErr)
			return
		} else if err == nil && existingUser != nil {
			// 当用户名已存在时，返回错误响应
			userExists := model.RegisterResponseDTO{
				BaseResp: model.BaseResp{
					Error:  err,
					Code:   http.StatusConflict,
					ErrMsg: "User already exists",
				},
			}
			c.JSON(http.StatusConflict, userExists)
			return
		}

		// 创建新用户
		newUser := &model.UserDTO{
			Name:     registerInfo.UserName,
			Password: registerInfo.Password,
		}
		if id, err := userRepository.CreateUser(newUser); err != nil {
			// 当用户创建失败时，返回错误响应
			createFailed := model.RegisterResponseDTO{
				BaseResp: model.BaseResp{
					Error:  err,
					Code:   http.StatusInternalServerError,
					ErrMsg: "Failed to create user",
				},
			}
			c.JSON(http.StatusInternalServerError, createFailed)
			return
		} else {
			newUser.Id = id
		}

		//生成jwt令牌
		expirationTime := time.Now().Add(24 * time.Hour) // 令牌有效期为 24 小时
		claims := &Claims{
			Username: registerInfo.UserName,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(expirationTime),
			},
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString(jwtKey)
		if err != nil {
			//生成令牌失败
			tokenFailed := model.RegisterResponseDTO{
				BaseResp: model.BaseResp{
					Error:  err,
					Code:   http.StatusInternalServerError,
					ErrMsg: "Failed to generate token",
				},
			}
			c.JSON(http.StatusInternalServerError, tokenFailed)
			return
		}

		// 注册成功，返回包含令牌和用户信息的成功响应
		okResp := model.RegisterResponseDTO{
			BaseResp: model.BaseResp{
				Error:  nil,
				Code:   http.StatusOK,
				ErrMsg: "",
			},
			Token: tokenString,
			User:  *newUser,
		}
		c.JSON(http.StatusOK, okResp)
	}
}


