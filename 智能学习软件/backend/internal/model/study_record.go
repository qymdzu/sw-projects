package model

import (
	"time"

	"github.com/google/uuid"
)

// StudyRecord 学习记录（数据模型设计 2.6）。
type StudyRecord struct {
	ID          uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID      uuid.UUID `gorm:"type:uuid;not null;index:idx_sr_user_date,priority:1" json:"user_id"`
	PlanID      uint64    `gorm:"not null" json:"plan_id"`
	Date        time.Time `gorm:"type:date;not null;index:idx_sr_user_date,priority:2" json:"date"`
	DurationMin int       `gorm:"not null" json:"duration_min"`
	Status      string    `gorm:"type:varchar(20);not null;default:done" json:"status"`
	Memo        *string   `gorm:"type:text" json:"memo,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

// TableName 指定 GORM 表名。
func (StudyRecord) TableName() string { return "study_records" }