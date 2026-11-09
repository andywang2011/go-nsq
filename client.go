package nsq

import (
	"context"
	"fmt"
	"time"

	nsqio "github.com/nsqio/go-nsq"
)

// Client holds a process-scoped NSQ producer and optional consumer lookup config.
// Use it from Wire as the single injectible NSQ dependency; create per-topic
// consumers with NewConsumer while publishing through the shared producer.
type Client struct {
	producer Producer
	lookup   ConsumerConfig
}

// NewClient returns a Client. producer may be nil when the process only consumes;
// lookup may be empty when the process only publishes.
func NewClient(producer Producer, lookup ConsumerConfig) *Client {
	return &Client{producer: producer, lookup: lookup}
}

// Publish sends a message to topic. Returns an error when the producer is unset.
func (c *Client) Publish(topic string, body []byte) error {
	if c == nil || c.producer == nil {
		return fmt.Errorf("nsq: producer is nil")
	}
	return c.producer.Publish(topic, body)
}

// MultiPublish sends multiple messages to topic.
func (c *Client) MultiPublish(topic string, bodies [][]byte) error {
	if c == nil || c.producer == nil {
		return fmt.Errorf("nsq: producer is nil")
	}
	return c.producer.MultiPublish(topic, bodies)
}

// DeferredPublish sends a delayed message to topic.
func (c *Client) DeferredPublish(topic string, delay time.Duration, body []byte) error {
	if c == nil || c.producer == nil {
		return fmt.Errorf("nsq: producer is nil")
	}
	return c.producer.DeferredPublish(topic, delay, body)
}

// Close stops the producer. It is safe on a nil Client or nil producer.
func (c *Client) Close() error {
	if c == nil || c.producer == nil {
		return nil
	}
	c.producer.Stop()
	return nil
}

// NewConsumer creates a consumer for the given topology using the Client lookup config.
func (c *Client) NewConsumer(ctx context.Context, topo *Topology, handler nsqio.Handler) (Consumer, error) {
	if c == nil {
		return nil, fmt.Errorf("nsq: client is nil")
	}
	return NewConsumer(ctx, topo, c.lookup, handler)
}
