package application

import (
	"context"

	notifDomain "taskflow-api/internal/notification/domain"
)

type NotificationService struct {
	notifs notifDomain.NotificationRepository
	prefs  notifDomain.PreferencesRepository
}

func NewNotificationService(n notifDomain.NotificationRepository, p notifDomain.PreferencesRepository) *NotificationService {
	return &NotificationService{notifs: n, prefs: p}
}

func (s *NotificationService) ListByUser(ctx context.Context, userID string) ([]*NotificationDTO, error) {
	items, err := s.notifs.ListByUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	out := make([]*NotificationDTO, len(items))
	for i, n := range items {
		out[i] = toDTO(n)
	}
	return out, nil
}

func (s *NotificationService) MarkAsRead(ctx context.Context, notifID, userID string) (*NotificationDTO, error) {
	n, err := s.notifs.FindByID(ctx, notifID)
	if err != nil {
		return nil, err
	}
	if n.UserID != userID {
		return nil, ErrForbidden
	}
	n.MarkAsRead()
	if err := s.notifs.Update(ctx, n); err != nil {
		return nil, err
	}
	return toDTO(n), nil
}

func (s *NotificationService) GetPreferences(ctx context.Context, userID string) (*PreferencesDTO, error) {
	p, err := s.prefs.Get(ctx, userID)
	if err != nil {
		return nil, err
	}
	return &PreferencesDTO{UserID: p.UserID, Enabled: p.Enabled}, nil
}

func (s *NotificationService) UpdatePreferences(ctx context.Context, dto UpdatePreferencesDTO) (*PreferencesDTO, error) {
	p, err := s.prefs.Get(ctx, dto.UserID)
	if err != nil {
		p = notifDomain.DefaultPreferences(dto.UserID)
	}
	for k, v := range dto.Enabled {
		p.Set(k, v)
	}
	if err := s.prefs.Save(ctx, p); err != nil {
		return nil, err
	}
	return &PreferencesDTO{UserID: p.UserID, Enabled: p.Enabled}, nil
}

func toDTO(n *notifDomain.Notification) *NotificationDTO {
	return &NotificationDTO{
		ID:        n.ID,
		UserID:    n.UserID,
		Type:      n.Type,
		Title:     n.Title,
		Body:      n.Body,
		ReadAt:    n.ReadAt,
		CreatedAt: n.CreatedAt,
	}
}
