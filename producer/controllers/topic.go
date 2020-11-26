package controllers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Shopify/sarama"

	"github.com/kafka/producer/dto"
)

// ListTopic ...
func ListTopic(w http.ResponseWriter, r *http.Request) {
	admin, err := sarama.NewClusterAdmin(
		[]string{"localhost:9091", "localhost:9092", "localhost:9093"},
		sarama.NewConfig(),
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	defer admin.Close()

	topics, err := admin.ListTopics()

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
		[]string{"localhost:9091", "localhost:9092", "localhost:9093"},
		config,
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	defer admin.Close()

	cleanupPolicy := "compact"

	err = admin.CreateTopic(topicRequest.Topic, &sarama.TopicDetail{
		NumPartitions:     topicRequest.Partitions,
		ReplicationFactor: topicRequest.Replications,
		ConfigEntries: map[string]*string{
			"cleanup.policy": &cleanupPolicy,
		},
	}, false)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Topic created successfully."))
}
