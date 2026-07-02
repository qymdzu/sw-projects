package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

// StudyReport 学习报告（数据模型设计 2.9）。
type StudyReport struct {
	ID          uint64         `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID      uuid.UUID      `gorm:"type:uuid;not null;uniqueIndex:idx_sr_user_period,priority:1" json:"user_id"`
	PeriodType  string         `gorm:"type:varchar(10);not null;uniqueIndex:idx_sr_user_period,priority:2" json:"period_type"`
	PeriodStart time.Time      `gorm:"type:date;not null;uniqueIndex:idx_sr_user_period,priority:3" json:"period_start"`
	PeriodEnd   time.Time      `gorm:"type:date;not null" json:"period_end"`
	Stats       datatypes.JSON `gorm:"type:jsonb;not null;default:'{}'" json:"stats"`
	CreatedAt   time.Time      `json:"created_at"`
}

// TableName 指定 GORM 表名。
func (StudyReport) TableName() string { return "study_reports" }