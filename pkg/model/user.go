package model

// UserDTO `用户`DTO结构体
type UserDTO struct {
	Id       int64  `json:"id"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

// UserDO `用户`存储数据结构体
type UserDO struct {
	Id       int64  `json:"id"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

func (userDO UserDO) TableName() string {
	return "user"
}

func (userDTO *UserDTO) Transfer() *UserDO {
	return &UserDO{
		Id:       userDTO.Id,
		Email:    userDTO.Email,
		Name:     userDTO.Name,
		Password: userDTO.Password,
	}
}

func (userDO *UserDO) Transfer() *UserDTO {
	return &UserDTO{
		Id:       userDO.Id,
		Email:    userDO.Email,
		Name:     userDO.Name,
		Password: userDO.Password,
	}
}

// GetUserByIdResponse 根据ID获取用户的返回体
type GetUserByIdResponse struct {
	BaseResp
	Userinfo UserDTO `json:"userinfo"`
}

// UpdateUserRequest 更新用户信息请求体
type UpdateUserRequest struct {
	Id       int64  `json:"id"`
	Name     string `json:"name"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

// UpdateUserResponse 更新用户信息返回体
type UpdateUserResponse struct {
	BaseResp
}

// UserDeleteResponse 删除用户返回体
type UserDeleteResponse struct {
	BaseResp
	DeleteId int64 `json:"deleted_user_id"` //被删除用户的id,我感觉是不是还是应该返回下这个比较好
}

// ChangePasswordRequest 修改密码的请求结构体
type ChangePasswordRequest struct {
	OldPassword     string `json:"OldPassword"`     // 旧密码
	NewPassword     string `json:"NewPassword"`     // 新密码
	ConfirmPassword string `json:"ConfirmPassword"` // 确认新密码
}

// ChangePasswordResponse 修改密码的返回结构体
type ChangePasswordResponse struct {
	BaseResp
}
