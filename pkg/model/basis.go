package model

// BaseResp 基础响应结构
type BaseResp struct {
	Error  error     `json:"error"`     // 错误
	ErrMsg string    `json:"error_msg"` // 前端提示的错误信息
	Code   ErrorCode `json:"code"`      // 错误码
}
