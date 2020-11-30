package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/Shopify/sarama"

	"github.com/kafka/producer/dto"
)

// ListTopic ...
func ListTopic(w http.ResponseWriter, r *http.Request) {
	log.Println("my-cluster-kafka-bootstrap:9092")
	admin, err := sarama.NewClusterAdmin(
		[]string{"my-cluster-kafka-bootstrap:9092"},
		sarama.NewConfig(),
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	defer admin.Close()

	topics, err := admin.ListTopics()
	log.Println(topics)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response, err := json.Marshal(topics)

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

// CreateTopic ...
func CreateTopic(w http.ResponseWriter, r *http.Request) {
	var topicRequest dto.CreateTopicRequest
	json.NewDecoder(r.Body).Decode(&topicRequest)
	
	config := sarama.NewConfig()
	config.Admin.Timeout = 10 * time.Second

	admin, err := sarama.NewClusterAdmin(
		[]string{"my-cluster-kafka-bootstrap:9092"},
		config,
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	defer admin.Close()

	err = admin.CreateTopic(topicRequest.Topic, &sarama.TopicDetail{
		NumPartitions:     topicRequest.Partitions,
		ReplicationFactor: topicRequest.Replications,
	}, false)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Topic created successfully."))
}
