package nsq

import (
	"context"
	"fmt"

	nsqio "github.com/nsqio/go-nsq"
)

var (
	newNSQProducer = nsqio.NewProducer
	producerPing   = func(p *nsqio.Producer) error { return p.Ping() }
)

var connectProducer = defaultConnectProducer

func defaultConnectProducer(cfg ProducerConfig) (Producer, error) {
	nsqCfg := applyProducerConfig(cfg)
	producer, err := newNSQProducer(cfg.Address, nsqCfg)
	if err != nil {
		return nil, fmt.Errorf("nsq new producer: %w", err)
	}
	if err := producerPing(producer); err != nil {
		producer.Stop()
		return nil, fmt.Errorf("nsq ping: %w", err)
	}
	return producer, nil
}

// NewProducer creates a producer and pings nsqd (returns mock-friendly Producer).
// ctx bounds the connect and ping; on cancellation a late producer is stopped.
func NewProducer(ctx context.Context, cfg ProducerConfig) (Producer, error) {
	ctx = ctxOrBackground(ctx)
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	return buildWithCtx(ctx, func() (Producer, error) { return connectProducer(cfg) }, func(p Producer) { p.Stop() })
}

// NewProducerWithRetry creates a producer with exponential-backoff startup retry.
// Retries (and their backoff) stop as soon as ctx is cancelled.
func NewProducerWithRetry(ctx context.Context, cfg ProducerConfig) (Producer, error) {
	ctx = ctxOrBackground(ctx)
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	retry := cfg.Retry
	if retry.MaxAttempts <= 0 {
		retry = DefaultRetryConfig()
	}
	normalizeRetry(&retry)
	rm := newRetryManager(retry)
	var producer Producer
	var lastErr error
	err := rm.execute(ctx, func(ctx context.Context) error {
		p, err := NewProducer(ctx, cfg)
		if err != nil {
			lastErr = err
			return err
		}
		producer = p
		return nil
	})
	if err != nil {
		if lastErr == nil {
			lastErr = err
		}
		return nil, fmt.Errorf("nsq producer retry failed after %d attempts: %w", retry.MaxAttempts, lastErr)
	}
	return producer, nil
}
