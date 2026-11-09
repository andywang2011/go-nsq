package nsq

import (
	"context"
	"errors"
	"testing"
	"time"

	nsqio "github.com/nsqio/go-nsq"
)

func TestNewClient_Publish(t *testing.T) {
	p := &fakeProducer{}
	c := NewClient(p, ConsumerConfig{})
	if err := c.Publish("t", []byte("x")); err != nil {
		t.Fatal(err)
	}
	if err := c.MultiPublish("t", [][]byte{[]byte("a"), []byte("b")}); err != nil {
		t.Fatal(err)
	}
	if err := c.DeferredPublish("t", time.Second, []byte("d")); err != nil {
		t.Fatal(err)
	}
	if len(p.published) != 4 {
		t.Fatalf("published=%v", p.published)
	}
}

func TestClient_Publish_NilProducer(t *testing.T) {
	c := NewClient(nil, ConsumerConfig{})
	if err := c.Publish("t", []byte("x")); err == nil {
		t.Fatal("expected error")
	}
	if err := c.MultiPublish("t", nil); err == nil {
		t.Fatal("expected error")
	}
	if err := c.DeferredPublish("t", 0, nil); err == nil {
		t.Fatal("expected error")
	}
	var nilClient *Client
	if err := nilClient.Publish("t", nil); err == nil {
		t.Fatal("expected error")
	}
}

func TestClient_Close(t *testing.T) {
	var nilClient *Client
	if err := nilClient.Close(); err != nil {
		t.Fatal(err)
	}
	c := NewClient(nil, ConsumerConfig{})
	if err := c.Close(); err != nil {
		t.Fatal(err)
	}
	stopped := false
	p := &stopTrackingProducer{onStop: func() { stopped = true }}
	c = NewClient(p, ConsumerConfig{})
	if err := c.Close(); err != nil {
		t.Fatal(err)
	}
	if !stopped {
		t.Fatal("expected Stop")
	}
}

type stopTrackingProducer struct {
	fakeProducer
	onStop func()
}

func (s *stopTrackingProducer) Stop() {
	if s.onStop != nil {
		s.onStop()
	}
}

func TestClient_NewConsumer_NilClient(t *testing.T) {
	var c *Client
	_, err := c.NewConsumer(context.Background(), &Topology{Topic: "t", Channel: "c"}, nsqio.HandlerFunc(func(*nsqio.Message) error { return nil }))
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestClient_NewConsumer_InvalidLookup(t *testing.T) {
	c := NewClient(nil, ConsumerConfig{})
	_, err := c.NewConsumer(context.Background(), &Topology{Topic: "t", Channel: "ch"}, nsqio.HandlerFunc(func(*nsqio.Message) error { return nil }))
	if err == nil {
		t.Fatal("expected error")
	}
	if !errors.Is(err, ErrInvalidConfig) && err.Error() == "" {
		t.Fatalf("unexpected err: %v", err)
	}
}
