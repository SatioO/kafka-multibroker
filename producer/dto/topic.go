package dto

// CreateTopicRequest ...
type CreateTopicRequest struct {
	Topic        string `json:"topic"`
	Partitions   int32  `json:"partitions"`
	Replications int16  `json:"replications"`
}
