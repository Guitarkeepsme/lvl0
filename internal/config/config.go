package config

import (
	"os"
	"strconv"
)

type Config struct {
	Port          string
	DBUrl         string
	NatsUrl       string
	NatsClusterID string
	OrdersBufSize int
}

func Load() Config {
	ordersBufSize, err := strconv.Atoi(os.Getenv("ORDERS_BUF_SIZE"))
	if err != nil {
		ordersBufSize = 10 // default value
	}

	return Config{
		Port:          os.Getenv("PORT"),
		DBUrl:         os.Getenv("DB_URL"),
		NatsUrl:       os.Getenv("NATS_URL"),
		NatsClusterID: os.Getenv("NATS_CLUSTER_ID"),
		OrdersBufSize: ordersBufSize,
	}
}
