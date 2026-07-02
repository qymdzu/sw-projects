// Package validator 提供轻量的入参校验辅助函数。
package validator

import (
	"errors"
	"regexp"
	"strings"
)

// 业务错误码定义（与 API设计 1.3 一致）。
var (
	ErrEmptyField   = errors.New("字段不能为空")
	ErrInvalidEmail = errors.New("邮箱格式不正确")
	ErrInvalidPhone = errors.New("手机号格式不正确")
	ErrPasswordWeak = errors.New("密码至少 8 位，且包含字母与数字")
	ErrInvalidRole  = errors.New("角色取值必须为 student/teacher/admin/parent")
)

// 中国大陆手机号校验（11 位，1 开头）。
var phoneRe = regexp.MustCompile(`^1[3-9]\d{9}$`)

// 邮箱校验（RFC 5322 简化版）。
var emailRe = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

// RequireString 校验非空。
func RequireString(field, value string) error {
	if strings.TrimSpace(value) == "" {
		return errors.New(field + " 不能为空")
	}
	return nil
}

// Phone 校验手机号。
func Phone(v string) error {
	if v == "" {
		return ErrEmptyField
	}
	if !phoneRe.MatchString(v) {
		return ErrInvalidPhone
	}
	return nil
}

// Email 校验邮箱（允许空字符串）。
func Email(v string) error {
	if v == "" {
		return nil
	}
	if !emailRe.MatchString(v) {
		return ErrInvalidEmail
	}
	return nil
}

// Password 校验密码强度（≥8 位，包含字母与数字）。
func Password(v string) error {
	if len(v) < 8 {
		return ErrPasswordWeak
	}
	hasLetter := regexp.MustCompile(`[A-Za-z]`).MatchString
	hasDigit := regexp.MustCompile(`\d`).MatchString
	if !hasLetter(v) || !hasDigit(v) {
		return ErrPasswordWeak
	}
	return nil
}

// Role 校验角色取值。
func Role(v string) error {
	switch v {
	case "student", "teacher", "admin", "parent":
		return nil
	default:
		return ErrInvalidRole
	}
}