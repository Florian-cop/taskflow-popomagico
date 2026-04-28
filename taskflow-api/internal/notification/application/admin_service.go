package application

import (
	"context"

	notifDomain "taskflow-api/internal/notification/domain"
)

// FailureToggle est satisfait par tout canal capable de simuler une panne (ex: FaultInjectingChannel).
// Permet au service admin de manipuler le canal sans dépendre de l'implémentation concrète.
type FailureToggle interface {
	Name() string
	SetFailing(failing bool)
	IsFailing() bool
}

type AdminService struct {
	dispatcher *Dispatcher
	failedRepo notifDomain.FailedNotificationRepository
	toggles    map[string]FailureToggle
}

func NewAdminService(d *Dispatcher, failedRepo notifDomain.FailedNotificationRepository, toggles ...FailureToggle) *AdminService {
	m := make(map[string]FailureToggle, len(toggles))
	for _, t := range toggles {
		m[t.Name()] = t
	}
	return &AdminService{dispatcher: d, failedRepo: failedRepo, toggles: m}
}

type FailedNotificationDTO struct {
	ID             string
	NotificationID string
	UserID         string
	Channel        string
	Type           string
	Title          string
	Body           string
	Error          string
	RetryCount     int
	Status         string
	OccurredAt     string
	LastRetriedAt  *string
}

type ChannelStatusDTO struct {
	Name    string
	Failing bool
}

func (s *AdminService) ListFailed(ctx context.Context, limit int) ([]*FailedNotificationDTO, error) {
	items, err := s.failedRepo.ListPending(ctx, limit)
	if err != nil {
		return nil, err
	}
	out := make([]*FailedNotificationDTO, len(items))
	for i, f := range items {
		dto := &FailedNotificationDTO{
			ID: f.ID, NotificationID: f.NotificationID, UserID: f.UserID,
			Channel: f.Channel, Type: f.Type, Title: f.Title, Body: f.Body,
			Error: f.Error, RetryCount: f.RetryCount, Status: string(f.Status),
			OccurredAt: f.OccurredAt.Format("2006-01-02T15:04:05Z07:00"),
		}
		if f.LastRetriedAt != nil {
			s := f.LastRetriedAt.Format("2006-01-02T15:04:05Z07:00")
			dto.LastRetriedAt = &s
		}
		out[i] = dto
	}
	return out, nil
}

func (s *AdminService) RetryFailed(ctx context.Context, id string) error {
	return s.dispatcher.RetryFailed(ctx, id)
}

func (s *AdminService) ListChannels() []*ChannelStatusDTO {
	out := make([]*ChannelStatusDTO, 0, len(s.dispatcher.Channels()))
	for _, c := range s.dispatcher.Channels() {
		failing := false
		if t, ok := s.toggles[c.Name()]; ok {
			failing = t.IsFailing()
		}
		out = append(out, &ChannelStatusDTO{Name: c.Name(), Failing: failing})
	}
	return out
}

func (s *AdminService) SetChannelFailing(channel string, failing bool) error {
	t, ok := s.toggles[channel]
	if !ok {
		return ErrChannelUnknown
	}
	t.SetFailing(failing)
	return nil
}
