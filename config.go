package nsq

import (
	"fmt"
	"time"

	nsqio "github.com/nsqio/go-nsq"
)

// ProducerConfig holds nsqd connection settings for a producer.
type ProducerConfig struct {
	// Address is the nsqd host:port (required).
	Address string
	// DialTimeout, ReadTimeout, WriteTimeout map to nsqio.Config when NSQConfig is nil.
	DialTimeout time.Duration
	// ReadTimeout bounds read operations on the producer connection.
	ReadTimeout time.Duration
	// WriteTimeout bounds write operations on the producer connection.
	WriteTimeout time.Duration
	// Snappy enables Snappy compression when NSQConfig is nil.
	Snappy bool
	// Deflate enables Deflate compression when NSQConfig is nil.
	Deflate bool
	// Retry optional startup connect retry (connect + ping only).
	Retry RetryConfig
	// NSQConfig optional full override of driver config.
	NSQConfig *nsqio.Config
}

// ConsumerConfig holds nsqlookupd discovery settings.
type ConsumerConfig struct {
	// LookupAddress single lookupd (used when LookupAddresses is empty).
	LookupAddress string
	// LookupAddresses multiple lookupd nodes for HA.
	LookupAddresses []string
	// LookupdPollInterval is how often the consumer polls lookupd.
	LookupdPollInterval time.Duration
	// NSQConfig optional full override before topology fields are applied.
	NSQConfig *nsqio.Config
}

// Topology defines topic/channel subscription and consumer concurrency.
type Topology struct {
	// Topic is the NSQ topic name (required).
	Topic string
	// Channel is the NSQ channel name (required).
	Channel string
	// Concurrency is the number of concurrent message handlers (default 10).
	Concurrency int
	// MaxAttempts is the max delivery attempts before giving up (default 5).
	MaxAttempts uint16
	// RequeueDelay is the delay before requeueing a failed message (default 90s).
	RequeueDelay time.Duration
}

// RetryConfig controls exponential-backoff retry during producer startup.
type RetryConfig struct {
	// MaxAttempts is the total number of attempts (including the first).
	MaxAttempts int
	// InitialDelay is the first backoff pause.
	InitialDelay time.Duration
	// MaxDelay caps exponential backoff.
	MaxDelay time.Duration
	// Multiplier scales delay after each failed attempt.
	Multiplier float64
}

const (
	defaultConcurrency  = 10
	defaultMaxAttempts  = 5
	defaultRequeueDelay = 90 * time.Second
	defaultDialTimeout  = time.Second
	defaultReadTimeout  = 60 * time.Second
	defaultWriteTimeout = time.Second
)

// ConfigDefault returns producer-oriented defaults (Address must still be set).
func ConfigDefault() ProducerConfig {
	return ProducerConfig{
		DialTimeout:  defaultDialTimeout,
		ReadTimeout:  defaultReadTimeout,
		WriteTimeout: defaultWriteTimeout,
		Retry:        DefaultRetryConfig(),
	}
}

// DefaultTopology returns consumer topology defaults (Topic/Channel must still be set).
func DefaultTopology() Topology {
	return Topology{
		Concurrency:  defaultConcurrency,
		MaxAttempts:  defaultMaxAttempts,
		RequeueDelay: defaultRequeueDelay,
	}
}

// DefaultRetryConfig returns production-safe retry defaults for producer startup.
func DefaultRetryConfig() RetryConfig {
	return RetryConfig{
		MaxAttempts:  3,
		InitialDelay: 100 * time.Millisecond,
		MaxDelay:     5 * time.Second,
		Multiplier:   2.0,
	}
}

// Validate checks producer config.
func (c *ProducerConfig) Validate() error {
	if c.Address == "" {
		return fmt.Errorf("%w: producer address is required", ErrInvalidConfig)
	}
	if c.Retry.MaxAttempts > 0 {
		normalizeRetry(&c.Retry)
	}
	return nil
}

// Validate checks consumer lookup config.
func (c *ConsumerConfig) Validate() error {
	if c.LookupAddress == "" && len(c.LookupAddresses) == 0 {
		return fmt.Errorf("%w: lookup address or lookup_addresses is required", ErrInvalidConfig)
	}
	return nil
}

// Validate checks topology and applies defaults.
func (t *Topology) Validate() error {
	if t.Topic == "" || t.Channel == "" {
		return fmt.Errorf("%w: topic and channel are required", ErrInvalidConfig)
	}
	if t.Concurrency < 1 {
		t.Concurrency = defaultConcurrency
	}
	if t.MaxAttempts == 0 {
		t.MaxAttempts = defaultMaxAttempts
	}
	if t.RequeueDelay == 0 {
		t.RequeueDelay = defaultRequeueDelay
	}
	return nil
}

func normalizeRetry(r *RetryConfig) {
	if r.MaxAttempts <= 0 {
		*r = DefaultRetryConfig()
	}
	if r.InitialDelay <= 0 {
		r.InitialDelay = DefaultRetryConfig().InitialDelay
	}
	if r.MaxDelay <= 0 {
		r.MaxDelay = DefaultRetryConfig().MaxDelay
	}
	if r.Multiplier <= 0 {
		r.Multiplier = DefaultRetryConfig().Multiplier
	}
	if r.MaxDelay < r.InitialDelay {
		r.MaxDelay = r.InitialDelay
	}
}
