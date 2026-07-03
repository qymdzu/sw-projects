// Package repository 提供数据访问层接口与 GORM 实现。

package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"smart-learning/internal/model"
)

// SettingRepository 是 model_settings 表的数据访问接口。
//
// 与 UserRepository 风格一致：暴露接口给 service 层，私有 GORM 实现。
type SettingRepository interface {
	GetByUserAndProvider(ctx context.Context, userID uuid.UUID, provider string) (*model.ModelSetting, error)
	GetActiveByUser(ctx context.Context, userID uuid.UUID) (*model.ModelSetting, error)
	ListByUser(ctx context.Context, userID uuid.UUID) ([]*model.ModelSetting, error)
	Upsert(ctx context.Context, m *model.ModelSetting) error
	SetActive(ctx context.Context, userID uuid.UUID, provider string) error
	DeleteByUserAndProvider(ctx context.Context, userID uuid.UUID, provider string) error
	CountActiveByUser(ctx context.Context, userID uuid.UUID) (int64, error)
}

type settingRepo struct {
	db *gorm.DB
}

// NewSettingRepository 构造 SettingRepository 实现。
func NewSettingRepository(db *gorm.DB) SettingRepository {
	return &settingRepo{db: db}
}

// GetByUserAndProvider 取指定 user+provider 的一条配置。
func (r *settingRepo) GetByUserAndProvider(ctx context.Context, userID uuid.UUID, provider string) (*model.ModelSetting, error) {
	var m model.ModelSetting
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND provider = ?", userID, provider).
		First(&m).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("query model setting: %w", err)
	}
	return &m, nil
}

// GetActiveByUser 取用户的默认（is_default=true）配置。
func (r *settingRepo) GetActiveByUser(ctx context.Context, userID uuid.UUID) (*model.ModelSetting, error) {
	var m model.ModelSetting
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND is_default = ?", userID, true).
		First(&m).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("query active model setting: %w", err)
	}
	return &m, nil
}

// ListByUser 列出用户所有配置（按 default 在前 + 最新更新时间排序）。
func (r *settingRepo) ListByUser(ctx context.Context, userID uuid.UUID) ([]*model.ModelSetting, error) {
	var ms []*model.ModelSetting
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("is_default DESC, updated_at DESC").
		Find(&ms).Error
	if err != nil {
		return nil, fmt.Errorf("list model settings: %w", err)
	}
	return ms, nil
}

// Upsert 按 (user_id, provider) 唯一键冲突时更新记录。
func (r *settingRepo) Upsert(ctx context.Context, m *model.ModelSetting) error {
	return r.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "user_id"}, {Name: "provider"}},
			DoUpdates: clause.AssignmentColumns([]string{
				"api_endpoint", "api_key_ciphertext", "api_key_nonce",
				"model", "config", "is_default", "updated_at",
			}),
		}).
		Create(m).Error
}

// SetActive 事务化切换用户的 default 配置：
//   1) 把该 user 的所有配置 is_default 置为 false
//   2) 把目标 provider 的配置 is_default 置为 true
//
// 若目标 provider 下用户没有配置（rows_affected=0），返回 ErrNotFound。
func (r *settingRepo) SetActive(ctx context.Context, userID uuid.UUID, provider string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&model.ModelSetting{}).
			Where("user_id = ?", userID).
			Update("is_default", false).Error; err != nil {
			return err
		}
		res := tx.Model(&model.ModelSetting{}).
			Where("user_id = ? AND provider = ?", userID, provider).
			Update("is_default", true)
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected == 0 {
			return ErrNotFound
		}
		return nil
	})
}

// DeleteByUserAndProvider 删除一条配置；记录不存在时返回 ErrNotFound。
func (r *settingRepo) DeleteByUserAndProvider(ctx context.Context, userID uuid.UUID, provider string) error {
	res := r.db.WithContext(ctx).
		Where("user_id = ? AND provider = ?", userID, provider).
		Delete(&model.ModelSetting{})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

// CountActiveByUser 统计用户当前 active 配置数量（用于 service 校验"不允许删完"）。
func (r *settingRepo) CountActiveByUser(ctx context.Context, userID uuid.UUID) (int64, error) {
	var n int64
	err := r.db.WithContext(ctx).
		Model(&model.ModelSetting{}).
		Where("user_id = ? AND is_default = ?", userID, true).
		Count(&n).Error
	return n, err
}