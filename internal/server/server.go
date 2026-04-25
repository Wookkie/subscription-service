package server

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Server struct {
	cfg       any
	httpServe *http.Server
}

func New(host string, port int) *Server {
	httpServe := http.Server{
		Addr: host + ":" + strconv.Itoa(port),
	}

	myServer := Server{
		httpServe: &httpServe,
	}

	myServer.configRoutes()

	return &myServer
}

func (s *Server) configRoutes() {
	router := gin.Default()
	router.GET("/")
	subscriptions := router.Group("/subscriptions")

	subscriptions.GET("/")
	subscriptions.GET("/:id")
	subscriptions.POST("/")
	subscriptions.PUT("/:id")
	subscriptions.DELETE("/:id")

	s.httpServe.Handler = router

}

func (s *Server) Run() error {
	return s.httpServe.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	return s.httpServe.Shutdown(ctx)
}
