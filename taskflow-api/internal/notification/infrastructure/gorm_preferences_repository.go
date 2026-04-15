package infrastructure

import (
	"context"
	"encoding/json"
	"errors"

	"gorm.io/gorm"

	"taskflow-api/internal/notification/domain"
)

type GormPreferencesRepository struct {
	db *gorm.DB
}

func NewGormPreferencesRepository(db *gorm.DB) *GormPreferencesRepository {
	return &GormPreferencesRepository{db: db}
}

func (r *GormPreferencesRepository) Get(ctx context.Context, userID string) (*domain.Preferences, error) {
	var m PreferencesModel
	err := r.db.WithContext(ctx).First(&m, "user_id = ?", userID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return domain.DefaultPreferences(userID), nil
	}
	if err != nil {
		return nil, err
	}
	enabled := map[string]bool{}
	if m.Enabled != "" {
		_ = json.Unmarshal([]byte(m.Enabled), &enabled)
	}
	return &domain.Preferences{UserID: m.UserID, Enabled: enabled}, nil
}

func (r *GormPreferencesRepository) Save(ctx context.Context, p *domain.Preferences) error {
	raw, err := json.Marshal(p.Enabled)
	if err != nil {
		return err
	}
	return r.db.WithContext(ctx).Save(&PreferencesModel{
		UserID:  p.UserID,
		Enabled: string(raw),
	}).Error
}
