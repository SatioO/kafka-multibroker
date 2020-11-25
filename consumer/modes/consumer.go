package modes

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"runtime"
	"sync/atomic"
	"time"

	"github.com/Shopify/sarama"
	"github.com/kafka/consumer/dto"
)

// ConsumerGroupHandler ...
type ConsumerGroupHandler interface {
	sarama.ConsumerGroupHandler
	WaitReady()
	Reset()
}

// ConsumerGroup ...
type ConsumerGroup struct {
	cg sarama.ConsumerGroup
}

// NewConsumerGroup ...
func NewConsumerGroup(broker []string, topics []string, group string, handler ConsumerGroupHandler) (*ConsumerGroup, error) {
	ctx := context.Background()
	cfg := sarama.NewConfig()
	cfg.Version = sarama.V0_10_2_0
	cfg.Consumer.Offsets.Initial = sarama.OffsetOldest
	client, err := sarama.NewConsumerGroup(broker, group, cfg)
	if err != nil {
		panic(err)
	}

	go func() {
		for {
			err := client.Consume(ctx, topics, handler)
			if err != nil {
				if err == sarama.ErrClosedConsumerGroup {
					break
				} else {
					panic(err)
				}
			}
			if ctx.Err() != nil {
				return
			}
			handler.Reset()
		}
	}()

	handler.WaitReady() // Await till the consumer has been set up

	return &ConsumerGroup{
		cg: client,
	}, nil
}

// Close ...
func (c *ConsumerGroup) Close() error {
	return c.cg.Close()
}

// ConsumerSessionMessage ...
type ConsumerSessionMessage struct {
	Session sarama.ConsumerGroupSession
	Message *sarama.ConsumerMessage
}

func decodeMessage(data []byte) error {
	var msg dto.Message
	err := json.Unmarshal(data, &msg)
	if err != nil {
		return err
	}
	return nil
}

// StartSyncConsumer ...
func StartSyncConsumer(broker, topic []string) (*ConsumerGroup, error) {
	var count int64
	var start = time.Now()
	handler := NewSyncConsumerGroupHandler(func(data []byte) error {
		if err := decodeMessage(data); err != nil {
			return err
		}
		count++
		fmt.Printf("sync consumer consumed %d messages at speed %.2f/s\n", count, float64(count)/time.Since(start).Seconds())

		return nil
	})

	consumer, err := NewConsumerGroup(broker, topic, "sync-consumer-"+fmt.Sprintf("%d", time.Now().Unix()), handler)
	if err != nil {
		return nil, err
	}
	return consumer, nil
}

// StartMultiAsyncConsumer ...
func StartMultiAsyncConsumer(broker, topic []string) (*ConsumerGroup, error) {
	var count int64
	var start = time.Now()
	var bufChan = make(chan *ConsumerSessionMessage, 1000)

	handler := NewMultiAsyncConsumerGroupHandler(&MultiAsyncConsumerConfig{
		BufChan: bufChan,
	})

	for i := 0; i < 8; i++ {
		go func() {
			for message := range bufChan {
				if err := decodeMessage(message.Message.Value); err == nil {
					message.Session.MarkMessage(message.Message, "")
				}
				cur := atomic.AddInt64(&count, 1)
				if cur%100 == 0 {
					fmt.Printf("multi async consumer consumed %d messages at speed %.2f/s\n", cur, float64(cur)/time.Since(start).Seconds())
				}
			}
		}()
	}
	consumer, err := NewConsumerGroup(broker, topic, "multi-async-consumer-"+fmt.Sprintf("%d", time.Now().Unix()), handler)
	if err != nil {
		return nil, err
	}
	return consumer, nil
}

// StartMultiBatchConsumer ...
func StartMultiBatchConsumer(broker, topic []string) (*ConsumerGroup, error) {
	var count int64
	var start = time.Now()
	var bufChan = make(chan batchMessages, 1000)
	log.Println(runtime.NumCPU())
	for i := 0; i < runtime.NumCPU(); i++ {
		go func() {
			for messages := range bufChan {
				for j := range messages {
					err := decodeMessage(messages[j].Message.Value)
					if err == nil {
						messages[j].Session.MarkMessage(messages[j].Message, "")
					}
				}
				cur := atomic.AddInt64(&count, int64(len(messages)))
				fmt.Printf("multi batch consumer consumed %d messages at speed %.2f/s\n", cur, float64(cur)/time.Since(start).Seconds())
			}
		}()
	}
	handler := NewMultiBatchConsumerGroupHandler(&MultiBatchConsumerConfig{
		MaxBufSize: 1000,
		BufChan:    bufChan,
	})
	consumer, err := NewConsumerGroup(broker, topic, "multi-batch-consumer", handler)
	if err != nil {
		return nil, err
	}
	log.Println("goroutines", runtime.NumGoroutine())
	return consumer, nil
}
