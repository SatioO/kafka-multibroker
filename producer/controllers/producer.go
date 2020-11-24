package controllers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/segmentio/kafka-go"

	"github.com/kafka/stream/dto"

	snappy "github.com/segmentio/kafka-go/snappy"
)

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
		Balancer:         &kafka.RoundRobin{},
		BatchTimeout:     10 * time.Millisecond,
		CompressionCodec: snappy.NewCompressionCodec(),
	})

	value, err := json.Marshal(request.Message)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = writer.WriteMessages(context.Background(), kafka.Message{
		Key:   []byte(request.Key),
		Value: []byte(value),
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := writer.Close(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Message produced successfully."))
}
