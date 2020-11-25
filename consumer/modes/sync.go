package modes

import "github.com/Shopify/sarama"

// Consumer ...
type Consumer struct {
	ready chan bool
	cb    func(data []byte) error
}

// NewSyncConsumerGroupHandler ...
func NewSyncConsumerGroupHandler(cb func(data []byte) error) ConsumerGroupHandler {
	handler := Consumer{
		ready: make(chan bool, 0),
		cb:    cb,
	}
	return &handler
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (consumer *Consumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (consumer *Consumer) Setup(sarama.ConsumerGroupSession) error {
	// Mark the consumer as ready
	close(consumer.ready)
	return nil
}

// WaitReady ...
func (consumer *Consumer) WaitReady() {
	// Mark the consumer as ready
	<-consumer.ready
	return
}

// Reset ...
func (consumer *Consumer) Reset() {
	consumer.ready = make(chan bool, 0)
	return
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
func (consumer *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {

	claimMsgChan := claim.Messages()

	for message := range claimMsgChan {
		if consumer.cb(message.Value) == nil {
			session.MarkMessage(message, "")
		}
	}

	return nil
}
