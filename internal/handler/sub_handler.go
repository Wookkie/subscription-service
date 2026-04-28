package handler

import (
	"net/http"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/Wookkie/subscription-service/internal/domain"
	"github.com/Wookkie/subscription-service/internal/service"
	"github.com/gin-gonic/gin"
)

type SubHandler struct {
	service *service.SubService
}

func NewSubHandler(service *service.SubService) *SubHandler {
	return &SubHandler{service: service}
}

// CreateSubscription godoc
// @Summary Create subscription
// @Description Create new subscription for user
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param input body domain.SubscriptionRequest true "subscription request"
// @Success 200 {object} domain.Subscription
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /subscriptions [post]
func (s *SubHandler) CreateSubscription(ctx *gin.Context) {
	var req domain.SubscriptionRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Error().Err(err).Msg("invalid create subscription request")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	t, err := time.Parse("01-2006", req.StartDate)
	if err != nil {
		log.Error().Err(err).Msg("invalid start_date format")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid start_date"})
		return
	}

	endDate := t.AddDate(0, 1, 0)

	sub := domain.Subscription{
		ServiceName: req.ServiceName,
		Price:       req.Price,
		UserID:      req.UserID,
		StartDate:   t,
		EndDate:     endDate,
	}

	subscription, err := s.service.CreateSubscription(sub)
	if err != nil {
		log.Error().Err(err).Msg("failed to create subscription")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, subscription)

}

// GetAllSubscriptions godoc
// @Summary Get all subscriptions
// @Description Get list of all subscriptions
// @Tags subscriptions
// @Produce json
// @Success 200 {array} domain.Subscription
// @Failure 500 {object} map[string]string
// @Router /subscriptions [get]
func (s *SubHandler) GetAllSubscriptions(ctx *gin.Context) {
	subscription, err := s.service.GetAllSubscriptions()
	if err != nil {
		log.Error().Err(err).Msg("failed to get subscriptions")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, subscription)
}

// GetSubscriptionByID godoc
// @Summary Get subscription by ID
// @Description Get subscription by UUID
// @Tags subscriptions
// @Produce json
// @Param id path string true "subscription id"
// @Success 200 {object} domain.Subscription
// @Failure 404 {object} map[string]string
// @Router /subscriptions/{id} [get]
func (s *SubHandler) GetSubscriptionByID(ctx *gin.Context) {
	id := ctx.Param("id")

	subscription, err := s.service.GetSubscriptionByID(id)
	if err != nil {
		log.Error().Err(err).Str("id", id).Msg("subscription not found")
		ctx.JSON(http.StatusNotFound, gin.H{"error": "not found subscription"})
		return
	}
	ctx.JSON(http.StatusOK, subscription)
}

// UpdateSubscription godoc
// @Summary Update subscription
// @Description Update price and end date
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param id path string true "subscription id"
// @Param input body domain.SubscriptionUpdateRequest true "update request"
// @Success 200 {object} domain.Subscription
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /subscriptions/{id} [put]
func (s *SubHandler) UpdateSubscription(ctx *gin.Context) {
	id := ctx.Param("id")

	var sub domain.SubscriptionUpdateRequest

	if err := ctx.ShouldBindJSON(&sub); err != nil {
		log.Error().Err(err).Msg("invalid update request")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	var endDate time.Time

	if sub.EndDate != "" {
		parsed, err := time.Parse("01-2006", sub.EndDate)
		if err != nil {
			log.Error().Err(err).Msg("failed to parse end_date")
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid end_date"})
			return
		}
		endDate = parsed
	}

	req := domain.Subscription{
		Price:   sub.Price,
		EndDate: endDate,
	}

	updated, err := s.service.UpdateSubscription(id, req)
	if err != nil {
		log.Error().Err(err).Str("id", id).Msg("update failed")
		ctx.JSON(http.StatusNotFound, gin.H{"error": "not found subscription"})
		return
	}
	ctx.JSON(http.StatusOK, updated)

}

// DeleteSubscription godoc
// @Summary Delete subscription
// @Description Delete subscription by ID
// @Tags subscriptions
// @Param id path string true "subscription id"
// @Success 204
// @Failure 404 {object} map[string]string
// @Router /subscriptions/{id} [delete]
func (s *SubHandler) DeleteSubscription(ctx *gin.Context) {
	id := ctx.Param("id")

	err := s.service.DeleteSubscription(id)
	if err != nil {
		log.Error().Err(err).Str("id", id).Msg("delete failed")
		ctx.JSON(http.StatusNotFound, gin.H{"error": "not found subscription"})
		return
	}
	ctx.Status(http.StatusNoContent)
}

// GetTotalCost godoc
// @Summary Get total cost
// @Description Calculate total subscription cost for period
// @Tags subscriptions
// @Produce json
// @Param user_id query string false "user id"
// @Param service_name query string false "service name"
// @Param from query string true "start month (MM-YYYY)"
// @Param to query string true "end month (MM-YYYY)"
// @Success 200 {object} domain.SubscriptionTotalResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /subscriptions/total [get]
func (h *SubHandler) GetTotalCost(c *gin.Context) {
	var req domain.SubscriptionTotalRequest

	if err := c.ShouldBindQuery(&req); err != nil {
		log.Error().Err(err).Msg("invalid query params")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid query"})
		return
	}

	if req.From == "" || req.To == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "from and to are required"})
		return
	}

	total, err := h.service.CalculateTotalCost(req)
	if err != nil {
		log.Error().Err(err).Msg("failed to calculate total")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Info().Int("total", total).Msg("total calculated")
	c.JSON(http.StatusOK, domain.SubscriptionTotalResponse{
		Total: total,
	})
}
