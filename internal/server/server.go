package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Wookkie/subscription-service/internal/config"
	"github.com/Wookkie/subscription-service/internal/database"
	"github.com/Wookkie/subscription-service/internal/handler"
	"github.com/Wookkie/subscription-service/internal/repository"
	"github.com/Wookkie/subscription-service/internal/service"
	"github.com/gin-gonic/gin"
)

type Server struct {
	cfg       *config.Config
	httpServe *http.Server
	db        *database.DBStorage
}

func New(cfg *config.Config) *Server {
	db, err := database.New(context.Background(), cfg.DBConn)
	if err != nil {
		panic(err)
	}

	if err := database.ApplyMigrations(cfg.DBConn); err != nil {
		panic(err)
	}

	httpServe := http.Server{
		Addr: fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
	}

	myServer := Server{
		httpServe: &httpServe,
		cfg:       cfg,
		db:        db,
	}

	myServer.configRoutes()

	return &myServer
}

func (s *Server) configRoutes() {
	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	repo := repository.NewSubscriptionRepository(s.db.DB)
	service := service.NewSubService(repo)
	handler := handler.NewSubHandler(service)

	subscriptions := router.Group("/subscriptions")

	subscriptions.GET("/", handler.GetAllSubscriptions)
	subscriptions.GET("/:id", handler.GetSubscriptionByID)
	subscriptions.POST("/", handler.CreateSubscription)
	subscriptions.PUT("/:id", handler.UpdateSubscription)
	subscriptions.DELETE("/:id", handler.DeleteSubscription)
	subscriptions.GET("/total", handler.GetTotalCost)

	s.httpServe.Handler = router

}

func (s *Server) Run() error {
	return s.httpServe.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	if s.db != nil {
		s.db.Close()
	}

	return s.httpServe.Shutdown(ctx)
}
