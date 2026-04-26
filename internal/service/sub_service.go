package service

import (
	"github.com/Wookkie/subscription-service/internal/domain"
	"github.com/Wookkie/subscription-service/internal/repository"
	"github.com/google/uuid"
)

type SubService struct {
	repo repository.SubscriptionRepository
}

func NewSubService(repo repository.SubscriptionRepository) *SubService {
	return &SubService{repo: repo}
}

func (s *SubService) CreateSubscription(sub domain.Subscription) (*domain.Subscription, error) {
	sub.ID = uuid.New().String()
	return s.repo.CreateSubscription(sub)
}

func (s *SubService) GetAllSubscriptions() ([]domain.Subscription, error) {
	return s.repo.GetAllSubscriptions()
}

func (s *SubService) GetSubscriptionByID(id string) (*domain.Subscription, error) {
	return s.repo.GetSubscriptionByID(id)
}

func (s *SubService) UpdateSubscription(id string, updated domain.Subscription) (*domain.Subscription, error) {
	return s.repo.UpdateSubscription(id, updated)
}

func (s *SubService) DeleteSubscription(id string) error {
	return s.repo.DeleteSubscription(id)
}
