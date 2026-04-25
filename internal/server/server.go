package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Wookkie/subscription-service/internal/config"
	"github.com/Wookkie/subscription-service/internal/handler"
	"github.com/Wookkie/subscription-service/internal/service"
	"github.com/gin-gonic/gin"
)

type Server struct {
	cfg       *config.Config
	httpServe *http.Server
}

func New(cfg *config.Config) *Server {
	httpServe := http.Server{
		Addr: fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
	}

	myServer := Server{
		httpServe: &httpServe,
		cfg:       cfg,
	}

	myServer.configRoutes()

	return &myServer
}

func (s *Server) configRoutes() {
	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
	service := service.NewSubService()
	handler := handler.NewSubHandler(service)
	subscriptions := router.Group("/subscriptions")

	subscriptions.GET("/", handler.GetAllSubscriptions)
	subscriptions.GET("/:id", handler.GetSubscriptionByID)
	subscriptions.POST("/", handler.CreateSubscription)
	subscriptions.PUT("/:id", handler.UpdateSubscription)
	subscriptions.DELETE("/:id", handler.DeleteSubscription)

	s.httpServe.Handler = router

}

func (s *Server) Run() error {
	return s.httpServe.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	return s.httpServe.Shutdown(ctx)
}
