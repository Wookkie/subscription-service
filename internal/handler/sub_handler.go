package handler

import (
	"net/http"

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
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"}) 
		return
	}

	sub := domain.Subscription{
		ServiceName: req.ServiceName,
		Price:       req.Price,
		UserID:      req.UserID,
	}

	subscription, err := s.service.CreateSubscription(ctx.Request.Context(), sub)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, subscription)

}

func (s *SubHandler) GetAllSubscriptions(ctx *gin.Context) {
	subscription, err := s.service.GetAllSubscriptions()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, subscription)
}

func (s *SubHandler) GetSubscriptionByID(ctx *gin.Context) {
	id := ctx.Param("id")

	subscription, err := s.service.GetSubscriptionByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "not found subscription"})
		return
	}
	ctx.JSON(http.StatusOK, subscription)
}

func (s *SubHandler) UpdateSubscription(ctx *gin.Context) {
	id := ctx.Param("id")

	var sub domain.Subscription

	updated, err := s.service.UpdateSubscription(id, sub)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "not found subscription"})
		return
	}
	ctx.JSON(http.StatusOK, updated)
}

func (s *SubHandler) DeleteSubscription(ctx *gin.Context) {
	id := ctx.Param("id")

	err := s.service.DeleteSubscription(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "not found subscription"})
		return
	}
	ctx.Status(http.StatusNoContent)
}
