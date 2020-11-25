package modes

import (
	"github.com/Shopify/sarama"
)

// MultiAsyncConsumerConfig ...
type MultiAsyncConsumerConfig struct {
	BufChan chan *ConsumerSessionMessage
}

type multiAsyncConsumerGroupHandler struct {
	cfg *MultiAsyncConsumerConfig

	ready chan bool
}

// NewMultiAsyncConsumerGroupHandler ...
func NewMultiAsyncConsumerGroupHandler(cfg *MultiAsyncConsumerConfig) ConsumerGroupHandler {
	handler := multiAsyncConsumerGroupHandler{
		ready: make(chan bool, 0),
	}

	handler.cfg = cfg

	return &handler
}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (h *multiAsyncConsumerGroupHandler) Setup(sarama.ConsumerGroupSession) error {
	// Mark the consumer as ready
	close(h.ready)
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (h *multiAsyncConsumerGroupHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (h *multiAsyncConsumerGroupHandler) WaitReady() {
	<-h.ready
	return
}

func (h *multiAsyncConsumerGroupHandler) Reset() {
	h.ready = make(chan bool, 0)
	return
}

func (h *multiAsyncConsumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	claimMsgChan := claim.Messages()

	for message := range claimMsgChan {
		h.cfg.BufChan <- &ConsumerSessionMessage{
			Session: session,
			Message: message,
		}
	}

	return nil
}
