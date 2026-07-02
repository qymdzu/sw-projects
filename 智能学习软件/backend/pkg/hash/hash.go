// Package hash 提供基于 bcrypt 的密码哈希与校验能力。
//
// 密码哈希 cost = 10（满足 NF-06 密码安全）。
package hash

import (
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// cost 是 bcrypt 计算成本，与架构设计 1.2 一致。
const cost = 10

// ErrPasswordMismatch 密码不匹配错误。
var ErrPasswordMismatch = errors.New("密码不匹配")

// Hash 使用 bcrypt 对明文密码进行哈希。
func Hash(password string) (string, error) {
	if password == "" {
		return "", errors.New("密码不能为空")
	}
	h, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		return "", fmt.Errorf("密码哈希失败: %w", err)
	}
	return string(h), nil
}

// Verify 校验明文密码与哈希是否匹配。
// 匹配返回 nil；不匹配返回 ErrPasswordMismatch；其他错误为系统错误。
func Verify(hashed, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password))
	if err == nil {
		return nil
	}
	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		return ErrPasswordMismatch
	}
	return fmt.Errorf("密码校验失败: %w", err)
}