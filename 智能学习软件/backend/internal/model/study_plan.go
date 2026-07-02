package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

// StudyPlan 学习计划表（数据模型设计 2.5）。
type StudyPlan struct {
	ID          uint64         `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID      uuid.UUID      `gorm:"type:uuid;not null;index:idx_sp_user_id" json:"user_id"`
	Goal        string         `gorm:"type:text;not null" json:"goal"`
	StartDate   time.Time      `gorm:"type:date;not null" json:"start_date"`
	EndDate     time.Time      `gorm:"type:date;not null" json:"end_date"`
	Items       datatypes.JSON `gorm:"type:jsonb;not null;default:'[]'" json:"items"`
	Status      string         `gorm:"type:varchar(20);not null;default:active;index:idx_sp_status" json:"status"`
	AIGenerated bool           `gorm:"not null;default:false" json:"ai_generated"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

// TableName 指定 GORM 表名。
func (StudyPlan) TableName() string { return "study_plans" }