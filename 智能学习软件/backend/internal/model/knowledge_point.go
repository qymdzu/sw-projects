package model

import "time"

// KnowledgePoint 知识点表（数据模型设计 2.3）。
type KnowledgePoint struct {
	ID        uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	SubjectID uint64    `gorm:"not null;index:idx_kp_subject_level,priority:1" json:"subject_id"`
	ParentID  *uint64   `gorm:"index:idx_kp_parent_id" json:"parent_id,omitempty"`
	Name      string    `gorm:"type:varchar(200);not null" json:"name"`
	Level     int       `gorm:"not null" json:"level"`
	Path      string    `gorm:"type:varchar(500);not null;index:idx_kp_path" json:"path"`
	CreatedAt time.Time `json:"created_at"`
}

// TableName 指定 GORM 表名。
func (KnowledgePoint) TableName() string { return "knowledge_points" }