package service

import (
	"time"

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

func (s *SubService) CalculateTotalCost(req domain.SubscriptionTotalRequest) (int, error) {
	from, err := time.Parse("01-2006", req.From)
	if err != nil {
		return 0, err
	}

	to, err := time.Parse("01-2006", req.To)
	if err != nil {
		return 0, err
	}

	from = time.Date(from.Year(), from.Month(), 1, 0, 0, 0, 0, time.UTC)
	to = time.Date(to.Year(), to.Month()+1, 0, 23, 59, 59, 0, time.UTC)

	subs, err := s.repo.GetSubscriptionsForPeriod(from, to, req.UserID, req.ServiceName)
	if err != nil {
		return 0, err
	}

	total := 0

	for _, sub := range subs {
		months := countMonths(sub.StartDate, sub.EndDate, from, to)
		total += months * sub.Price
	}

	return total, nil
}

func countMonths(start, end, from, to time.Time) int {
	if end.IsZero() {
		end = to
	}

	if start.Before(from) {
		start = from
	}
	if end.After(to) {
		end = to
	}

	if start.After(end) {
		return 0
	}

	months := (end.Year()-start.Year())*12 + int(end.Month()-start.Month()) + 1
	return months
}
