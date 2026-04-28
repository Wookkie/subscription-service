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

	sub := domain.Subscription{
		ServiceName: req.ServiceName,
		Price:       req.Price,
		UserID:      req.UserID,
		StartDate:   t,
	}

	subscription, err := s.service.CreateSubscription(sub)
	if err != nil {
		log.Error().Err(err).Msg("failed to create subscription")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, subscription)

}

func (s *SubHandler) GetAllSubscriptions(ctx *gin.Context) {
	subscription, err := s.service.GetAllSubscriptions()
	if err != nil {
		log.Error().Err(err).Msg("failed to get subscriptions")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, subscription)
}

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

func (s *SubHandler) UpdateSubscription(ctx *gin.Context) {
	id := ctx.Param("id")

	var sub domain.SubscriptionUpdateRequest

	if err := ctx.ShouldBindJSON(&sub); err != nil {
		log.Error().Err(err).Msg("invalid update request")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	req := domain.Subscription{
		Price:   sub.Price,
		EndDate: time.Time{},
	}

	updated, err := s.service.UpdateSubscription(id, req)
	if err != nil {
		log.Error().Err(err).Str("id", id).Msg("update failed")
		ctx.JSON(http.StatusNotFound, gin.H{"error": "not found subscription"})
		return
	}
	ctx.JSON(http.StatusOK, updated)

}

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
