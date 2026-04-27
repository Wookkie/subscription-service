package domain

type SubscriptionTotalRequest struct {
	UserID      string `form:"user_id"`
	ServiceName string `form:"service_name"`
	From        string `form:"from"`
	To          string `form:"to"`
}

type SubscriptionTotalResponse struct {
	Total int `json:"total"`
}
