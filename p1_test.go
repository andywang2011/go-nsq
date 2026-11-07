package nsq

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	nsqio "github.com/nsqio/go-nsq"
)

type stopProducer struct {
	fakeProducer
	stopped atomic.Int32
}

func (s *stopProducer) Stop() { s.stopped.Add(1) }

func TestBuildWithCtx(t *testing.T) {
	v, err := buildWithCtx(context.Background(), func() (int, error) { return 1, nil }, func(int) {})
	if err != nil || v != 1 {
		t.Fatal(v, err)
	}
	if _, err := buildWithCtx(context.Background(), func() (int, error) { return 0, errors.New("x") }, func(int) {}); err == nil {
		t.Fatal("expected error")
	}
	cctx, ccancel := context.WithCancel(context.Background())
	ccancel()
	if _, err := buildWithCtx(cctx, func() (int, error) { return 1, nil }, func(int) {}); !errors.Is(err, context.Canceled) {
		t.Fatal(err)
	}
	if ctxOrBackground(nil) == nil { //nolint:staticcheck
		t.Fatal("nil ctx must become Background")
	}

	ctx, cancel := context.WithCancel(context.Background())
	release := make(chan struct{})
	cleaned := make(chan int, 1)
	go func() { time.Sleep(10 * time.Millisecond); cancel() }()
	_, err = buildWithCtx(ctx, func() (int, error) { <-release; return 7, nil }, func(v int) { cleaned <- v })
	if !errors.Is(err, context.Canceled) {
		t.Fatal(err)
	}
	close(release)
	select {
	case v := <-cleaned:
		if v != 7 {
			t.Fatal(v)
		}
	case <-time.After(time.Second):
		t.Fatal("late value must be cleaned up")
	}
}

func TestNewProducer_ctxCancelStopsLateProducer(t *testing.T) {
	old := connectProducer
	defer func() { connectProducer = old }()
	release := make(chan struct{})
	late := &stopProducer{}
	connectProducer = func(ProducerConfig) (Producer, error) { <-release; return late, nil }
	ctx, cancel := context.WithCancel(context.Background())
	go func() { time.Sleep(10 * time.Millisecond); cancel() }()
	if _, err := NewProducer(ctx, ProducerConfig{Address: "a:1"}); !errors.Is(err, context.Canceled) {
		t.Fatal(err)
	}
	close(release)
	deadline := time.Now().Add(time.Second)
	for late.stopped.Load() == 0 && time.Now().Before(deadline) {
		time.Sleep(5 * time.Millisecond)
	}
	if late.stopped.Load() != 1 {
		t.Fatal("late producer not stopped")
	}
}

func TestNewProducerWithRetry_ctxCancelStopsRetry(t *testing.T) {
	old := connectProducer
	defer func() { connectProducer = old }()
	var calls atomic.Int32
	connectProducer = func(ProducerConfig) (Producer, error) { calls.Add(1); return nil, errors.New("down") }
	ctx, cancel := context.WithCancel(context.Background())
	go func() { time.Sleep(30 * time.Millisecond); cancel() }()
	start := time.Now()
	_, err := NewProducerWithRetry(ctx, ProducerConfig{
		Address: "a:1",
		Retry:   RetryConfig{MaxAttempts: 50, InitialDelay: time.Hour, MaxDelay: time.Hour, Multiplier: 1},
	})
	if err == nil || !errors.Is(err, context.Canceled) && calls.Load() < 1 || time.Since(start) > time.Second {
		t.Fatal(err, calls.Load())
	}
	// Already cancelled: no connect attempt and still a wrapped error.
	calls.Store(0)
	if _, err := NewProducerWithRetry(ctx, ProducerConfig{Address: "a:1"}); err == nil || calls.Load() != 0 {
		t.Fatal(err, calls.Load())
	}
	if _, err := NewProducerWithRetry(nil, ProducerConfig{}); err == nil { //nolint:staticcheck
		t.Fatal("expected validation error")
	}
}

func TestNewConsumer_ctxCancelled(t *testing.T) {
	old := connectLookupd
	defer func() { connectLookupd = old }()
	release := make(chan struct{})
	exited := make(chan struct{})
	connectLookupd = func(*nsqio.Consumer, ConsumerConfig) error {
		defer close(exited)
		<-release
		return nil
	}
	ctx, cancel := context.WithCancel(context.Background())
	go func() { time.Sleep(10 * time.Millisecond); cancel() }()
	topo := &Topology{Topic: "t", Channel: "c"}
	_, err := NewConsumer(ctx, topo, ConsumerConfig{LookupAddress: "a:1"}, nsqio.HandlerFunc(func(*nsqio.Message) error { return nil }))
	if !errors.Is(err, context.Canceled) {
		t.Fatal(err)
	}
	close(release)
	<-exited // the background goroutine must finish before the seam is restored

	if _, err := NewConsumer(nil, &Topology{Topic: "t", Channel: "c"}, ConsumerConfig{}, nil); err == nil { //nolint:staticcheck
		t.Fatal("expected validation error")
	}
}

func TestRetryManager_ctxPaths(t *testing.T) {
	rm := newRetryManager(RetryConfig{MaxAttempts: 5, InitialDelay: time.Hour, MaxDelay: time.Hour, Multiplier: 1})
	ctx, cancel := context.WithCancel(context.Background())
	go func() { time.Sleep(20 * time.Millisecond); cancel() }()
	if err := rm.execute(ctx, func(context.Context) error { return errors.New("x") }); !errors.Is(err, context.Canceled) {
		t.Fatal(err)
	}
	ctx2, cancel2 := context.WithCancel(context.Background())
	if err := rm.execute(ctx2, func(context.Context) error { cancel2(); return errors.New("x") }); !errors.Is(err, context.Canceled) {
		t.Fatal(err)
	}
	if err := rm.execute(nil, func(context.Context) error { return nil }); err != nil { //nolint:staticcheck
		t.Fatal(err)
	}
	if err := rm.sleep(context.Background(), 0); err != nil {
		t.Fatal(err)
	}
}
