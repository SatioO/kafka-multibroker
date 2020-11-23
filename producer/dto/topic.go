package dto

// CreateTopicRequest ...
type CreateTopicRequest struct {
	Topic        string `json:"topic"`
	Partitions   int    `json:"partitions"`
	Replications int    `json:"replications"`
}
