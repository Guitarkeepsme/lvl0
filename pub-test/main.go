package main

import (
	"io"
	"os"

	"github.com/nats-io/stan.go"
)

func main() {
	sc, err := stan.Connect("my-cluster", "publisher", stan.NatsURL("nats://localhost:4222"))
	if err != nil {
		panic(err)
	}
	defer sc.Close()

	file, err := os.Open("model.json")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	json, _ := io.ReadAll(file)

	sc.Publish("order", json)

	file2, err := os.Open("model2.json")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	json2, _ := io.ReadAll(file2)

	sc.Publish("order", json2)

	file3, err := os.Open("model3.json")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	json3, _ := io.ReadAll(file3)

	sc.Publish("order", json3)
}
