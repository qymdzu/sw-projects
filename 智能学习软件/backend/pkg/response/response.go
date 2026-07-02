// Package response 提供统一的 HTTP 响应封装。
//
// 与 API 设计 1.2 统一响应格式保持一致。
package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// 业务错误码常量（与 API设计 1.3 章节保持一致）。
const (
	CodeSuccess        = 0
	CodeBadRequest     = 10001
	CodeResourceExists = 10002
	CodeResourceState  = 10003
	CodeUnauthorized   = 20001
	CodeTokenExpired   = 20002
	CodeTokenInvalid   = 20003
	CodeForbidden      = 20004
	CodeNotFound       = 30001
	CodeConflict       = 30002
	CodeRateLimited    = 40001
	CodeServerError    = 50001
	CodeUpstreamError  = 50002
)

// Response 是统一响应结构。
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Detail  interface{} `json:"detail,omitempty"`
}

// PaginationData 是分页数据格式。
type PaginationData struct {
	Items    interface{} `json:"items"`
	Total    int64       `json:"total"`
	Page     int         `json:"page"`
	PageSize int         `json:"page_size"`
}

// OK 业务成功响应（HTTP 200）。
func OK(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    CodeSuccess,
		Message: "success",
		Data:    data,
	})
}

// Created 资源创建成功（HTTP 201）。
func Created(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, Response{
		Code:    CodeSuccess,
		Message: "success",
		Data:    data,
	})
}

// Page 返回分页响应。
func Page(c *gin.Context, items interface{}, total int64, page, pageSize int) {
	c.JSON(http.StatusOK, Response{
		Code:    CodeSuccess,
		Message: "success",
		Data: PaginationData{
			Items:    items,
			Total:    total,
			Page:     page,
			PageSize: pageSize,
		},
	})
}

// Fail 自定义错误响应。
func Fail(c *gin.Context, httpStatus int, code int, msg string, detail interface{}) {
	c.AbortWithStatusJSON(httpStatus, Response{
		Code:    code,
		Message: msg,
		Detail:  detail,
	})
}

// BadRequest 参数错误。
func BadRequest(c *gin.Context, msg string, detail interface{}) {
	Fail(c, http.StatusBadRequest, CodeBadRequest, msg, detail)
}

// Unauthorized 未认证。
func Unauthorized(c *gin.Context, msg string) {
	Fail(c, http.StatusUnauthorized, CodeUnauthorized, msg, nil)
}

// TokenExpired Token 过期。
func TokenExpired(c *gin.Context) {
	Fail(c, http.StatusUnauthorized, CodeTokenExpired, "Token 已过期", nil)
}

// TokenInvalid Token 无效。
func TokenInvalid(c *gin.Context) {
	Fail(c, http.StatusUnauthorized, CodeTokenInvalid, "Token 无效", nil)
}

// Forbidden 权限不足。
func Forbidden(c *gin.Context, msg string) {
	Fail(c, http.StatusForbidden, CodeForbidden, msg, nil)
}

// NotFound 资源不存在。
func NotFound(c *gin.Context, msg string) {
	Fail(c, http.StatusNotFound, CodeNotFound, msg, nil)
}

// Conflict 资源冲突。
func Conflict(c *gin.Context, msg string) {
	Fail(c, http.StatusConflict, CodeConflict, msg, nil)
}

// ServerError 服务器内部错误。
func ServerError(c *gin.Context, msg string) {
	Fail(c, http.StatusInternalServerError, CodeServerError, msg, nil)
}