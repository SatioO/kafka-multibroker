package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/Shopify/sarama"
	"github.com/gorilla/mux"
	"github.com/kafka/producer/dto"
)

// Producer ...
type Producer struct {
	p sarama.AsyncProducer
}

// NewProducer ...
func NewProducer(broker []string) (*Producer, error) {
	config := sarama.NewConfig()
	// config.Producer.Idempotent = true
	config.Producer.Compression = sarama.CompressionSnappy

	producer, err := sarama.NewAsyncProducer(broker, config)
	if err != nil {
		return nil, err
	}
	return &Producer{p: producer}, nil
}

var id int

// StartProduce ...
func (p *Producer) StartProduce(w http.ResponseWriter, r *http.Request) {
	var request dto.ProducerRequest
	params := mux.Vars(r)

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	value, _ := json.Marshal(request.Message)
	id++

	p.p.Input() <- &sarama.ProducerMessage{
		Topic: params["topic_name"],
		Key:   sarama.ByteEncoder([]byte(request.Key)),
		Value: sarama.ByteEncoder([]byte(value)),
	}

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Message produced successfully."))
}

// Close ...
func (p *Producer) Close() error {
	if p != nil {
		return p.p.Close()
	}
	return nil
}
