package nsq

import (
	"context"
	"errors"
	"testing"
	"time"

	nsqio "github.com/nsqio/go-nsq"
)

func TestNewProducer_invalidConfig(t *testing.T) {
	_, err := NewProducer(context.Background(), ProducerConfig{})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestNewProducer_success(t *testing.T) {
	old := connectProducer
	defer func() { connectProducer = old }()
	connectProducer = func(ProducerConfig) (Producer, error) {
		return &fakeProducer{}, nil
	}
	p, err := NewProducer(context.Background(), ProducerConfig{Address: "127.0.0.1:4150"})
	if err != nil {
		t.Fatal(err)
	}
	p.Stop()
}

func TestNewProducer_connectError(t *testing.T) {
	old := connectProducer
	defer func() { connectProducer = old }()
	connectProducer = func(ProducerConfig) (Producer, error) {
		return nil, errors.New("connect failed")
	}
	_, err := NewProducer(context.Background(), ProducerConfig{Address: "127.0.0.1:4150"})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestNewProducer_newNSQProducerError(t *testing.T) {
	oldNew := newNSQProducer
	defer func() { newNSQProducer = oldNew }()
	newNSQProducer = func(string, *nsqio.Config) (*nsqio.Producer, error) {
		return nil, errors.New("new producer failed")
	}
	_, err := NewProducer(context.Background(), ProducerConfig{Address: "127.0.0.1:4150"})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestNewProducer_pingError(t *testing.T) {
	oldNew := newNSQProducer
	oldPing := producerPing
	defer func() {
		newNSQProducer = oldNew
		producerPing = oldPing
	}()
	newNSQProducer = func(addr string, cfg *nsqio.Config) (*nsqio.Producer, error) {
		return nsqio.NewProducer(addr, cfg)
	}
	producerPing = func(*nsqio.Producer) error {
		return errors.New("ping failed")
	}
	_, err := NewProducer(context.Background(), ProducerConfig{Address: "127.0.0.1:4150"})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestNewProducerWithRetry_eventuallySuccess(t *testing.T) {
	old := connectProducer
	defer func() { connectProducer = old }()
	attempts := 0
	connectProducer = func(ProducerConfig) (Producer, error) {
		attempts++
		if attempts < 2 {
			return nil, errors.New("temporary")
		}
		return &fakeProducer{}, nil
	}
	cfg := ProducerConfig{
		Address: "127.0.0.1:4150",
		Retry: RetryConfig{
			MaxAttempts:  3,
			InitialDelay: time.Millisecond,
			MaxDelay:     2 * time.Millisecond,
			Multiplier:   1.1,
		},
	}
	p, err := NewProducerWithRetry(context.Background(), cfg)
	if err != nil {
		t.Fatal(err)
	}
	if attempts != 2 {
		t.Fatalf("attempts = %d", attempts)
	}
	p.Stop()
}

func TestNewProducerWithRetry_exhausted(t *testing.T) {
	old := connectProducer
	defer func() { connectProducer = old }()
	connectProducer = func(ProducerConfig) (Producer, error) {
		return nil, errors.New("fail")
	}
	cfg := ProducerConfig{
		Address: "127.0.0.1:4150",
		Retry: RetryConfig{
			MaxAttempts:  2,
			InitialDelay: time.Millisecond,
			MaxDelay:     time.Millisecond,
			Multiplier:   1,
		},
	}
	_, err := NewProducerWithRetry(context.Background(), cfg)
	if err == nil {
		t.Fatal("expected error")
	}
}

func Test_defaultConnectProducer_integrationHook(t *testing.T) {
	oldNew := newNSQProducer
	oldPing := producerPing
	defer func() {
		newNSQProducer = oldNew
		producerPing = oldPing
	}()
	newNSQProducer = func(string, *nsqio.Config) (*nsqio.Producer, error) {
		return nsqio.NewProducer("127.0.0.1:4150", nsqio.NewConfig())
	}
	producerPing = func(p *nsqio.Producer) error { return nil }
	p, err := defaultConnectProducer(ProducerConfig{Address: "127.0.0.1:4150"})
	if err != nil {
		t.Fatal(err)
	}
	p.Stop()
}

func TestNewProducerWithRetry_usesDefaultRetry(t *testing.T) {
	old := connectProducer
	defer func() { connectProducer = old }()
	connectProducer = func(ProducerConfig) (Producer, error) {
		return &fakeProducer{}, nil
	}
	p, err := NewProducerWithRetry(context.Background(), ProducerConfig{Address: "127.0.0.1:4150"})
	if err != nil {
		t.Fatal(err)
	}
	p.Stop()
}

func Test_retryManager_contextCancel(t *testing.T) {
	rm := newRetryManager(DefaultRetryConfig())
	rm.sleep = func(context.Context, time.Duration) error { return nil }
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	err := rm.execute(ctx, func(context.Context) error { return errors.New("x") })
	if err == nil {
		t.Fatal("expected error")
	}
}
