package model

// GetCaptchaResponse 获取验证码返回结构体
type GetCaptchaResponse struct {
	BaseResp
	Captcha string `json:"captcha"`
}
