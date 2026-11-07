package nsq

import (
	"time"

	nsqio "github.com/nsqio/go-nsq"
)

// Producer is a mock-friendly NSQ producer (method set matches github.com/nsqio/go-nsq.Producer).
type Producer interface {
	Publish(topic string, body []byte) error
	DeferredPublish(topic string, delay time.Duration, body []byte) error
	MultiPublish(topic string, body [][]byte) error
	Stop()
}

var _ Producer = (*nsqio.Producer)(nil)
