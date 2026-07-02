package model

import (
	"time"

	"github.com/google/uuid"
)

// MistakeBook 错题本（数据模型设计 2.8）。
type MistakeBook struct {
	ID               uint64     `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID           uuid.UUID  `gorm:"type:uuid;not null;uniqueIndex:idx_mb_user_question,priority:1" json:"user_id"`
	QuestionID       uint64     `gorm:"not null;uniqueIndex:idx_mb_user_question,priority:2" json:"question_id"`
	KnowledgePointID uint64     `gorm:"not null;index:idx_mb_user_kp,priority:2" json:"knowledge_point_id"`
	WrongAnswer      string     `gorm:"type:text;not null" json:"wrong_answer"`
	MistakeCount     int        `gorm:"not null;default:1" json:"mistake_count"`
	Mastered         bool       `gorm:"not null;default:false;index:idx_mb_user_mastered,priority:2" json:"mastered"`
	MasteredAt       *time.Time `json:"mastered_at,omitempty"`
	LastReviewedAt   time.Time  `gorm:"not null;index:idx_mb_last_reviewed" json:"last_reviewed_at"`
	CreatedAt        time.Time  `json:"created_at"`
}

// TableName 指定 GORM 表名。
func (MistakeBook) TableName() string { return "mistake_books" }