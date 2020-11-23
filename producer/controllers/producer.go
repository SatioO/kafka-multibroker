package controllers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/segmentio/kafka-go"

	"github.com/kafka/stream/dto"

	snappy "github.com/segmentio/kafka-go/snappy"
)

func generator(writer *kafka.Writer) <-chan int {
	c := make(chan int)

	go func() {
		defer close(c)
		for i := 0; i < 100; i++ {
			writer.WriteMessages(context.Background(), kafka.Message{
				Key:   []byte("Key" + strconv.Itoa(i)),
				Value: []byte("$" + strconv.Itoa(i) + "000"),
			})
		}
	}()

	return c
}

// Producer ...
func Producer(w http.ResponseWriter, r *http.Request) {
	var request dto.ProducerRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:          []string{"localhost:9091", "localhost:9092", "localhost:9093"},
		Topic:            "balance_details",
		Balancer:         &kafka.LeastBytes{},
		CompressionCodec: snappy.NewCompressionCodec(),
	})

	c := generator(writer)

	for i := 0; i < 100; i++ {
		<-c
	}

	if err := writer.Close(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Message produced successfully."))
}
