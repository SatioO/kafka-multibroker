package controllers

import (
	"encoding/json"
	"net"
	"net/http"
	"strconv"

	"github.com/kafka/stream/dto"
	"github.com/segmentio/kafka-go"
)

// ListTopic ...
func ListTopic(w http.ResponseWriter, r *http.Request) {
	conn, err := kafka.Dial("tcp", "localhost:9092")

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	defer conn.Close()

	partitions, err := conn.ReadPartitions()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	m := map[string]struct{}{}

	for _, p := range partitions {
		m[p.Topic] = struct{}{}
	}

	response, err := json.Marshal(m)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

// CreateTopic ...
func CreateTopic(w http.ResponseWriter, r *http.Request) {
	var topicRequest dto.CreateTopicRequest
	json.NewDecoder(r.Body).Decode(&topicRequest)

	conn, err := kafka.Dial("tcp", "localhost:9092")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer conn.Close()

	controller, err := conn.Controller()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var ctrlConn *kafka.Conn
	ctrlConn, err = kafka.Dial("tcp", net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port)))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer ctrlConn.Close()

	topicConfig := []kafka.TopicConfig{
		kafka.TopicConfig{
			Topic:             topicRequest.Topic,
			NumPartitions:     topicRequest.Partitions,
			ReplicationFactor: topicRequest.Replications,
			ConfigEntries: []kafka.ConfigEntry{
				kafka.ConfigEntry{
					ConfigName:  "cleanup.policy",
					ConfigValue: "compact",
				},
				kafka.ConfigEntry{
					ConfigName:  "delete.retention.ms",
					ConfigValue: "100",
				},
				kafka.ConfigEntry{
					ConfigName:  "segment.ms",
					ConfigValue: "100",
				},
				kafka.ConfigEntry{
					ConfigName:  "min.cleanable.dirty.ratio",
					ConfigValue: "0.01",
				},
			},
		},
	}

	err = ctrlConn.CreateTopics(topicConfig...)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Topic created successfully."))
}
