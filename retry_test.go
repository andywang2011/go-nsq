package nsq

import (
	"context"
	"errors"
	"testing"
	"time"
)

func Test_normalizeRetry(t *testing.T) {
	r := RetryConfig{}
	normalizeRetry(&r)
	if r.MaxAttempts != DefaultRetryConfig().MaxAttempts {
		t.Fatal(r)
	}
	r2 := RetryConfig{MaxAttempts: 2, InitialDelay: 0, MaxDelay: 0, Multiplier: 0}
	normalizeRetry(&r2)
	if r2.InitialDelay <= 0 || r2.MaxDelay < r2.InitialDelay {
		t.Fatal(r2)
	}
	r3 := RetryConfig{MaxAttempts: 2, InitialDelay: time.Second, MaxDelay: time.Millisecond, Multiplier: 2}
	normalizeRetry(&r3)
	if r3.MaxDelay < r3.InitialDelay {
		t.Fatal("max delay should be >= initial")
	}
}

func Test_retryManager_execute_successFirst(t *testing.T) {
	rm := newRetryManager(DefaultRetryConfig())
	if err := rm.execute(context.Background(), func(context.Context) error { return nil }); err != nil {
		t.Fatal(err)
	}
}

func Test_retryManager_jitter(t *testing.T) {
	rm := newRetryManager(DefaultRetryConfig())
	d := rm.jitter(time.Millisecond)
	if d < 0 {
		t.Fatal(d)
	}
	d0 := rm.jitter(0)
	if d0 < 0 {
		t.Fatal(d0)
	}
}

func Test_retryManager_sleepError(t *testing.T) {
	rm := newRetryManager(RetryConfig{MaxAttempts: 2, InitialDelay: time.Millisecond, MaxDelay: time.Millisecond, Multiplier: 1})
	rm.sleep = func(context.Context, time.Duration) error { return errors.New("sleep") }
	err := rm.execute(context.Background(), func(context.Context) error { return errors.New("op") })
	if err == nil {
		t.Fatal("expected error")
	}
}

func Test_ProducerConfig_Validate_normalizesRetry(t *testing.T) {
	cfg := ProducerConfig{
		Address: "127.0.0.1:4150",
		Retry:   RetryConfig{MaxAttempts: 1, InitialDelay: time.Millisecond, MaxDelay: time.Second, Multiplier: 2},
	}
	if err := cfg.Validate(); err != nil {
		t.Fatal(err)
	}
}
