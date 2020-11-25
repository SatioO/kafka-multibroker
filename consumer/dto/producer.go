package dto

// Message ...
type Message struct {
	Name    string `json:"name"`
	Balance string `json:"balance"`
}

// ProducerRequest ...
type ProducerRequest struct {
	Key     string  `json:"key"`
	Message Message `json:"message"`
}
