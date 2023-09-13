package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"service/internal/config"
	"service/internal/reader"
	"service/internal/repository"
	"service/internal/server"
)

func main() {
	config := config.Load()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	db, err := repository.NewDB(config)
	if err != nil {
		log.Fatalf("repository error: %v", err)
	}
	defer db.Stop()

	cache, err := repository.NewCache(db)
	if err != nil {
		log.Fatalf("cache error: %v", err)
	}
	defer cache.Stop()

	nats := reader.NewReader(config, db, cache)
	defer nats.Stop()

	server := server.NewServer(config, cache)
	defer server.Stop()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		s := (<-sig).String()
		log.Printf("stopping with %s signal\n", s)
		cancel()
	}()

	<-ctx.Done()
}
