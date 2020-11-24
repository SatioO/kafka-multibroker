package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/kafka/producer/dto"
)

// ListTopic ...
func ListTopic(w http.ResponseWriter, r *http.Request) {
	// conn, err := kafka.Dial("tcp", "kafka-2:9092")

	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusBadRequest)
	// 	return
	// }

	// defer conn.Close()

	// partitions, err := conn.ReadPartitions()
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusBadRequest)
	// 	return
	// }

	// m := map[string]struct{}{}

	// for _, p := range partitions {
	// 	m[p.Topic] = struct{}{}
	// }

	// response, err := json.Marshal(m)

	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusBadRequest)
	// 	return
	// }

	// w.Header().Set("content-type", "application/json")
	// w.WriteHeader(http.StatusOK)
	// w.Write(response)
}

// CreateTopic ...
func CreateTopic(w http.ResponseWriter, r *http.Request) {
	var topicRequest dto.CreateTopicRequest
	json.NewDecoder(r.Body).Decode(&topicRequest)

	// kafka.NewAdminClient(&kafka.ConfigMap{})

	// if err != nil {
	// 	fmt.Printf("Failed to create Admin client: %s\n", err)
	// 	os.Exit(1)
	// }

	// results, err := client.GetMetadata(&topicRequest.Topic, false, 1000)

	// if err != nil {
	// 	fmt.Printf("Failed to create topic: %v\n", err)
	// 	os.Exit(1)
	// }

	// log.Println(results)

	// client.Close()

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Topic created successfully."))
}
