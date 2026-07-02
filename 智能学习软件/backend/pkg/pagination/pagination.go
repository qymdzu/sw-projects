// Package pagination 提供分页参数解析与默认值管理。
package pagination

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

// 默认与边界值。
const (
	DefaultPage     = 1
	DefaultPageSize = 20
	MaxPageSize     = 100
)

// Params 是分页参数。
type Params struct {
	Page     int
	PageSize int
}

// Parse 从 gin.Context 中解析分页参数；缺省值与上限见常量定义。
func Parse(c *gin.Context) Params {
	page := parseInt(c.Query("page"), DefaultPage)
	if page < 1 {
		page = DefaultPage
	}
	pageSize := parseInt(c.Query("page_size"), DefaultPageSize)
	if pageSize < 1 {
		pageSize = DefaultPageSize
	}
	if pageSize > MaxPageSize {
		pageSize = MaxPageSize
	}
	return Params{Page: page, PageSize: pageSize}
}

// Offset 计算偏移量。
func (p Params) Offset() int { return (p.Page - 1) * p.PageSize }

// Limit 计算 limit。
func (p Params) Limit() int { return p.PageSize }

func parseInt(s string, def int) int {
	if s == "" {
		return def
	}
	v, err := strconv.Atoi(s)
	if err != nil {
		return def
	}
	return v
}