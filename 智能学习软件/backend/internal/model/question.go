package model

import (
	"time"

	"gorm.io/datatypes"
)

// Question 题目表（数据模型设计 2.4）。
type Question struct {
	ID               uint64           `gorm:"primaryKey;autoIncrement" json:"id"`
	SubjectID        uint64           `gorm:"not null;index:idx_q_subject_kp_diff,priority:1" json:"subject_id"`
	KnowledgePointID uint64           `gorm:"not null;index:idx_q_subject_kp_diff,priority:2" json:"knowledge_point_id"`
	Type             string           `gorm:"type:varchar(20);not null;index:idx_q_type" json:"type"`
	Difficulty       int              `gorm:"not null;default:3" json:"difficulty"`
	Content          datatypes.JSON   `gorm:"type:jsonb;not null" json:"content"`
	Options          datatypes.JSON   `gorm:"type:jsonb" json:"options,omitempty"`
	Answer           string           `gorm:"type:text;not null" json:"answer"`
	Analysis         *string          `gorm:"type:text" json:"analysis,omitempty"`
	CreatedAt        time.Time        `json:"created_at"`
	UpdatedAt        time.Time        `json:"updated_at"`
	KnowledgePoint   *KnowledgePoint  `gorm:"foreignKey:KnowledgePointID;references:ID" json:"knowledge_point,omitempty"`
}

// TableName 指定 GORM 表名。
func (Question) TableName() string { return "questions" }