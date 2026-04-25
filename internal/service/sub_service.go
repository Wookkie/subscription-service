package service

import (
	"context"

	"github.com/Wookkie/subscription-service/internal/domain"
	"github.com/google/uuid"
)

type SubService struct {
}

func NewSubService() *SubService {
	return &SubService{}
}

func (s *SubService) CreateSubscription(ctx context.Context, sub domain.Subscription) (*domain.Subscription, error) {
	sub.ID = uuid.New().String()
	return &sub, nil
}

func (s *SubService) GetAllSubscriptions() ([]domain.Subscription, error) {
	return []domain.Subscription{}, nil
}

func (s *SubService) GetSubscriptionByID(id string) (*domain.Subscription, error) {
	return nil, nil
}

func (s *SubService) UpdateSubscription(id string, updated domain.Subscription) (*domain.Subscription, error) {
	updated.ID = id
	return &updated, nil
}

func (s *SubService) DeleteSubscription(id string) error {
	return nil
}
