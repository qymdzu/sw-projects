package model

import "time"

// Subject 科目表（数据模型设计 2.2）。
type Subject struct {
	ID          uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string    `gorm:"type:varchar(100);not null;uniqueIndex" json:"name"`
	Description *string   `gorm:"type:text" json:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

// TableName 指定 GORM 表名。
func (Subject) TableName() string { return "subjects" }