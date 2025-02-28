package model

// BaseResp 基础响应结构
type BaseResp struct {
	Error error     `json:"error"` // 错误
	Code  ErrorCode `json:"code"`  // 错误码
}
