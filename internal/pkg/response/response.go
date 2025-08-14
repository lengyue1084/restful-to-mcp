package response

import (
	"fmt"
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

type OptFunc func(*Response)

// Success 创建成功响应
func Success(message string, opts ...OptFunc) *Response {
	if message == "" {
		message = GetMessage(SuccessCode)
	}

	resp := &Response{
		Code:    SuccessCode,
		Message: message,
		Data:    nil,
	}

	// 应用选项
	for _, opt := range opts {
		opt(resp)
	}

	return resp
}

// WithData 设置响应数据
func WithData(data interface{}) OptFunc {
	return func(r *Response) {
		r.Data = data
	}
}

// Fail 创建失败响应
func Fail(message string, data interface{}) *Response {
	if message == "" {
		message = GetMessage(FailCode)
	}
	return &Response{
		Code:    FailCode,
		Message: message,
		Data:    data,
	}
}

// NewResponse 创建自定义响应
func NewResponse(code int, message string, data interface{}) *Response {
	if message == "" {
		message = GetMessage(code)
	}
	return &Response{
		Code:    code,
		Message: message,
		Data:    data,
	}
}

// Error 实现 error 接口
func (r *Response) Error() string {
	return fmt.Sprintf("code=%d message=%s data=%+v", r.Code, r.Message, r.Data)
}

// GetMessage 获取错误代码对应的消息
func GetMessage(code int) string {
	if msg, exists := errorMessages[code]; exists {
		return msg
	}
	return "unknown error"
}

// Gin 集成函数

// JSON 返回标准JSON响应
func JSON(c *gin.Context, resp *Response) {
	statusCode := getHTTPStatusCode(resp.Code)
	c.JSON(statusCode, resp)
}

// SuccessJSON 返回成功响应
func SuccessJSON(c *gin.Context, message string, opts ...OptFunc) {
	resp := Success(message, opts...)
	c.JSON(http.StatusOK, resp)
}

// ErrorJSON 返回错误响应
func ErrorJSON(c *gin.Context, resp *Response) {
	statusCode := getHTTPStatusCode(resp.Code)
	c.JSON(statusCode, resp)
}

// AbortWithJSON 中断请求并返回JSON响应
func AbortWithJSON(c *gin.Context, resp *Response) {
	statusCode := getHTTPStatusCode(resp.Code)
	c.AbortWithStatusJSON(statusCode, resp)
}

// 根据错误代码获取HTTP状态码
func getHTTPStatusCode(code int) int {
	switch {
	case code == SuccessCode:
		return http.StatusOK
	case code >= 400 && code < 500:
		return http.StatusBadRequest
	case code >= 500:
		return http.StatusInternalServerError
	case code >= 10000:
		return http.StatusBadRequest
	default:
		return http.StatusOK
	}
}

// 便捷函数创建特定错误并直接返回给客户端
func ClientNotSupportJSON(c *gin.Context, message string) {
	resp := ClientNotSupport(message)
	ErrorJSON(c, resp)
}

func ServerNotSupportJSON(c *gin.Context, message string) {
	resp := ServerNotSupport(message)
	ErrorJSON(c, resp)
}

func RequestInvalidJSON(c *gin.Context, message string) {
	resp := RequestInvalid(message)
	ErrorJSON(c, resp)
}

func MethodNotSupportJSON(c *gin.Context, message string) {
	resp := MethodNotSupport(message)
	ErrorJSON(c, resp)
}

func JSONUnmarshalErrorJSON(c *gin.Context, message string) {
	resp := JSONUnmarshalError(message)
	ErrorJSON(c, resp)
}

func SessionNotInitializedJSON(c *gin.Context, message string) {
	resp := SessionNotInitialized(message)
	ErrorJSON(c, resp)
}

func SessionClosedJSON(c *gin.Context, message string) {
	resp := SessionClosed(message)
	ErrorJSON(c, resp)
}

// ClientNotSupport 便捷函数创建特定错误
func ClientNotSupport(message string) *Response {
	if message == "" {
		message = GetMessage(ClientNotSupportCode)
	}
	return NewResponse(ClientNotSupportCode, message, nil)
}

func ServerNotSupport(message string) *Response {
	if message == "" {
		message = GetMessage(ServerNotSupportCode)
	}
	return NewResponse(ServerNotSupportCode, message, nil)
}

func RequestInvalid(message string) *Response {
	if message == "" {
		message = GetMessage(RequestInvalidCode)
	}
	return NewResponse(RequestInvalidCode, message, nil)
}

func MethodNotSupport(message string) *Response {
	if message == "" {
		message = GetMessage(MethodNotSupportCode)
	}
	return NewResponse(MethodNotSupportCode, message, nil)
}

func JSONUnmarshalError(message string) *Response {
	if message == "" {
		message = GetMessage(JSONUnmarshalErrorCode)
	}
	return NewResponse(JSONUnmarshalErrorCode, message, nil)
}

func SessionNotInitialized(message string) *Response {
	if message == "" {
		message = GetMessage(SessionNotInitializedCode)
	}
	return NewResponse(SessionNotInitializedCode, message, nil)
}

func SessionClosed(message string) *Response {
	if message == "" {
		message = GetMessage(SessionClosedCode)
	}
	return NewResponse(SessionClosedCode, message, nil)
}

// IsSuccess 判断是否为成功响应
func (r *Response) IsSuccess() bool {
	return r.Code == SuccessCode
}

// IsError 判断是否为错误响应
func (r *Response) IsError() bool {
	return r.Code != SuccessCode
}
