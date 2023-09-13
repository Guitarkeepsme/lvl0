package server

import (
	"context"
	"log"
	"net/http"

	"service/internal/config"
	"service/internal/domain"

	"github.com/gin-gonic/gin"
)

type CacheReader interface {
	GetOrder(id string) *domain.Order
}

type Server struct {
	cache  CacheReader
	server *http.Server
}

func NewServer(config config.Config, cache CacheReader) *Server {
	s := &Server{
		cache: cache,
	}

	g := gin.Default()

	g.GET("/order/:id", s.GetOrder)

	s.server = &http.Server{
		Addr:    ":" + config.Port,
		Handler: g,
	}

	go func() {
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen and serve error: %v\n", err)
		}
	}()

	return s
}

func (s *Server) Stop() {
	s.server.Shutdown(context.TODO())
}

func (s *Server) GetOrder(c *gin.Context) {
	order := s.cache.GetOrder(c.Param("id"))
	if order == nil {
		c.Status(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, order)
}
