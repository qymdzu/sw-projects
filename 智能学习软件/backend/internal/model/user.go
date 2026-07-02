// Package model 定义 GORM 数据模型，与 docs/design/数据模型设计.md 1:1 映射。
package model

import (
	"time"

	"github.com/google/uuid"
)

// User 用户表（数据模型设计 2.1）。
type User struct {
	ID           uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Name         string     `gorm:"type:varchar(100);not null" json:"name"`
	Phone        *string    `gorm:"type:varchar(20);uniqueIndex:idx_users_phone" json:"phone,omitempty"`
	Email        *string    `gorm:"type:varchar(255);uniqueIndex:idx_users_email" json:"email,omitempty"`
	PasswordHash string     `gorm:"type:varchar(255);not null" json:"-"`
	Role         string     `gorm:"type:varchar(20);not null;default:student" json:"role"`
	AvatarURL    *string    `gorm:"type:varchar(500)" json:"avatar_url,omitempty"`
	ParentID     *uuid.UUID `gorm:"type:uuid;index:idx_users_parent_id" json:"parent_id,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

// TableName 指定 GORM 表名。
func (User) TableName() string { return "users" }