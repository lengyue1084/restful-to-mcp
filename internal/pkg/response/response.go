package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// 定义错误代码枚举
const (
	SuccessCode int = iota
	FailCode
)

const (
	ClientNotSupportCode      = 10001
	ServerNotSupportCode      = 10002
	RequestInvalidCode        = 10003
	MethodNotSupportCode      = 10004
	JSONUnmarshalErrorCode    = 10005
	SessionNotInitializedCode = 10006
	SessionClosedCode         = 10007
	SendEOFCode               = 10008
)

// 预定义错误消息
var errorMessages = map[int]string{
	SuccessCode:               "success",
	FailCode:                  "operation failed",
	ClientNotSupportCode:      "this feature client not support",
	ServerNotSupportCode:      "this feature server not support",
	RequestInvalidCode:        "request invalid",
	MethodNotSupportCode:      "method not support",
	JSONUnmarshalErrorCode:    "json unmarshal error",
	SessionNotInitializedCode: "the session has not been initialized",
	SessionClosedCode:         "session closed",
	SendEOFCode:               "send EOF",
}

// Response 错误响应结构体
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// Success 创建并返回成功响应
func Success(ctx *gin.Context, message string, data interface{}) {
	if message == "" {
		message = GetMessage(SuccessCode)
	}
	resp := &Response{
		Code:    SuccessCode,
		Message: message,
		Data:    data,
	}
	ctx.JSON(http.StatusOK, resp)
}

// Error 创建并返回错误响应
func Error(ctx *gin.Context, message string, data interface{}) {
	if message == "" {
		message = GetMessage(FailCode)
	}
	resp := &Response{
		Code:    FailCode,
		Message: message,
		Data:    data,
	}
	//statusCode := getHTTPStatusCode(FailCode)
	ctx.JSON(http.StatusOK, resp)
}

// CustomError 创建并返回自定义错误响应
func CustomError(ctx *gin.Context, code int, message string, data interface{}) {
	if message == "" {
		message = GetMessage(code)
	}

	resp := &Response{
		Code:    code,
		Message: message,
		Data:    data,
	}

	statusCode := getHTTPStatusCode(code)
	ctx.JSON(statusCode, resp)
}

// AbortWithErrorResponse 中断请求并返回错误响应
func AbortWithErrorResponse(ctx *gin.Context, code int, message string, data interface{}) {
	if message == "" {
		message = GetMessage(code)
	}

	resp := &Response{
		Code:    code,
		Message: message,
		Data:    data,
	}

	statusCode := getHTTPStatusCode(code)
	ctx.AbortWithStatusJSON(statusCode, resp)
}

// GetMessage 获取错误代码对应的消息
func GetMessage(code int) string {
	if msg, exists := errorMessages[code]; exists {
		return msg
	}
	return "unknown error"
}

// 根据错误代码获取HTTP状态码
func getHTTPStatusCode(code int) int {
	switch {
	case code == SuccessCode:
		return http.StatusOK
	case code >= 400 && code < 500:
		return http.StatusBadRequest
	case code >= 10000:
		return http.StatusBadRequest
	case code >= 500:
		return http.StatusInternalServerError
	default:
		return http.StatusOK
	}
}

// 便捷函数 - 直接返回特定错误

func ClientNotSupport(ctx *gin.Context, message string) {
	if message == "" {
		message = GetMessage(ClientNotSupportCode)
	}
	CustomError(ctx, ClientNotSupportCode, message, nil)
}

func ServerNotSupport(ctx *gin.Context, message string) {
	if message == "" {
		message = GetMessage(ServerNotSupportCode)
	}
	CustomError(ctx, ServerNotSupportCode, message, nil)
}

func RequestInvalid(ctx *gin.Context, message string) {
	if message == "" {
		message = GetMessage(RequestInvalidCode)
	}
	CustomError(ctx, RequestInvalidCode, message, nil)
}

func MethodNotSupport(ctx *gin.Context, message string) {
	if message == "" {
		message = GetMessage(MethodNotSupportCode)
	}
	CustomError(ctx, MethodNotSupportCode, message, nil)
}

func JSONUnmarshalError(ctx *gin.Context, message string) {
	if message == "" {
		message = GetMessage(JSONUnmarshalErrorCode)
	}
	CustomError(ctx, JSONUnmarshalErrorCode, message, nil)
}

func SessionNotInitialized(ctx *gin.Context, message string) {
	if message == "" {
		message = GetMessage(SessionNotInitializedCode)
	}
	CustomError(ctx, SessionNotInitializedCode, message, nil)
}

func SessionClosed(ctx *gin.Context, message string) {
	if message == "" {
		message = GetMessage(SessionClosedCode)
	}
	CustomError(ctx, SessionClosedCode, message, nil)
}

// IsSuccess 判断是否为成功响应
func (r *Response) IsSuccess() bool {
	return r.Code == SuccessCode
}

// IsError 判断是否为错误响应
func (r *Response) IsError() bool {
	return r.Code != SuccessCode
}
