package model

import (
	"time"

	"github.com/google/uuid"
)

// ExerciseRecord 练习记录（数据模型设计 2.7）。
type ExerciseRecord struct {
	ID          uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID      uuid.UUID `gorm:"type:uuid;not null;index:idx_er_user_question,priority:1" json:"user_id"`
	QuestionID  uint64    `gorm:"not null;index:idx_er_user_question,priority:2" json:"question_id"`
	Answer      string    `gorm:"type:text;not null" json:"answer"`
	IsCorrect   bool      `gorm:"not null;index:idx_er_user_correct_time,priority:2" json:"is_correct"`
	Score       *int      `json:"score,omitempty"`
	DurationSec *int      `json:"duration_sec,omitempty"`
	CreatedAt   time.Time `gorm:"index:idx_er_created_at" json:"created_at"`
}

// TableName 指定 GORM 表名。
func (ExerciseRecord) TableName() string { return "exercise_records" }