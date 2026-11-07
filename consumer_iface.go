package nsq

import nsqio "github.com/nsqio/go-nsq"

// ConsumerHandle abstracts consumer shutdown for tests and graceful stop helpers.
type ConsumerHandle interface {
	// Stop signals the underlying NSQ consumer to stop.
	Stop()
}

// Consumer is a mock-friendly NSQ consumer (subset of github.com/nsqio/go-nsq.Consumer).
type Consumer interface {
	ConsumerHandle
	// StopChan is closed (or receives) once the consumer has fully stopped after Stop.
	StopChan() <-chan int
	// ChangeMaxInFlight updates the max in-flight message count at runtime.
	ChangeMaxInFlight(maxInFlight int)
	// IsStarved reports whether the consumer is starved for messages.
	IsStarved() bool
	// Stats returns a snapshot of consumer counters.
	Stats() *nsqio.ConsumerStats
}

type consumerWrapper struct {
	c *nsqio.Consumer
}

var _ Consumer = (*consumerWrapper)(nil)

// WrapConsumer wraps an upstream consumer in the Consumer interface.
// It returns nil when c is nil.
func WrapConsumer(c *nsqio.Consumer) Consumer {
	if c == nil {
		return nil
	}
	return &consumerWrapper{c: c}
}

func (w *consumerWrapper) Stop() {
	if w != nil && w.c != nil {
		w.c.Stop()
	}
}

func (w *consumerWrapper) StopChan() <-chan int {
	if w == nil || w.c == nil {
		return nil
	}
	return w.c.StopChan
}

func (w *consumerWrapper) ChangeMaxInFlight(maxInFlight int) {
	if w != nil && w.c != nil {
		w.c.ChangeMaxInFlight(maxInFlight)
	}
}

func (w *consumerWrapper) IsStarved() bool {
	if w == nil || w.c == nil {
		return false
	}
	return w.c.IsStarved()
}

func (w *consumerWrapper) Stats() *nsqio.ConsumerStats {
	if w == nil || w.c == nil {
		return nil
	}
	return w.c.Stats()
}
