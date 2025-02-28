package model

type ErrorCode int

const (
	Success ErrorCode = 0

	UserExists    ErrorCode = 301
	UserNotExists ErrorCode = 302
	PasswordError ErrorCode = 303

	InternalError      ErrorCode = 500
	InvalidRequestBody ErrorCode = 501
)
