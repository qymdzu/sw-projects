// Package service 提供业务逻辑层。
//
// 所有 Service 通过接口暴露，依赖 Repository 接口；
// 便于单元测试中使用 mock，符合 docs/design/目录结构.md 第 2 节。
package service

import "errors"

// 业务层通用错误。
var (
	ErrInvalidParam    = errors.New("参数错误")
	ErrUnauthorized    = errors.New("未认证")
	ErrPasswordInvalid = errors.New("密码错误")
	ErrAccountConflict = errors.New("账号已存在")
	ErrResourceMissing = errors.New("资源不存在")
	ErrForbidden       = errors.New("权限不足")
	ErrAccountNotFound = errors.New("账号不存在")
)