package domain

import "time"

type Subscription struct {
	ID          string    `json:"id"`
	ServiceName string    `json:"service_name"`
	Price       int       `json:"price"`
	UserID      string    `json:"user_id"`
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
}

type SubscriptionRequest struct {
	ServiceName string `json:"service_name"`
	Price       int    `json:"price"`
	UserID      string `json:"user_id"`
	StartDate   string `json:"start_date"`
}

type SubscriptionUpdateRequest struct {
	Price   int    `json:"price"`
	EndDate string `json:"end_date"`
}
