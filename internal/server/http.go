package server

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"service/config"
	"service/internal/domain"
)

type CacheReader interface {
	GetOrder(id string) *domain.Order
}

type server struct {
	cache  CacheReader
	server *http.Server
}

func NewServer(config config.Config, cache CacheReader) *server {
	s := &server{
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

func (s *server) Stop() {
	s.server.Shutdown(context.TODO())
}

func (s *server) GetOrder(c *gin.Context) {
	order := s.cache.GetOrder(c.Param("id"))
	if order == nil {
		c.Status(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, order)
}
