package repository

import "github.com/Wookkie/subscription-service/internal/domain"

type SubscriptionRepository interface {
	GetAllSubscriptions() ([]domain.Subscription, error)
	GetSubscriptionByID(id string) (*domain.Subscription, error)
	CreateSubscription(sub domain.Subscription) (*domain.Subscription, error)
	UpdateSubscription(id string, sub domain.Subscription) (*domain.Subscription, error)
	DeleteSubscription(id string) error
}
