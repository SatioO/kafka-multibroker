package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/kafka/consumer/modes"
)

func main() {
	topic := []string{"balance_details"}
	broker := []string{"localhost:9091", "localhost:9092", "localhost:9093"}

	consumer, _ := modes.StartMultiBatchConsumer(broker, topic)
	defer consumer.Close()
	log.Println("Sarama consumer up and running!...")

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	fmt.Println("received signal", <-c)
}
