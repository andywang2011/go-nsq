package nsq

import (
	"context"
	"errors"
	"testing"

	nsqio "github.com/nsqio/go-nsq"
)

func TestNewConsumer_nilTopology(t *testing.T) {
	_, err := NewConsumer(context.Background(), nil, ConsumerConfig{LookupAddress: "127.0.0.1:4161"}, nsqio.HandlerFunc(func(*nsqio.Message) error { return nil }))
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestNewConsumer_invalidTopology(t *testing.T) {
	topo := &Topology{Topic: "", Channel: "c"}
	_, err := NewConsumer(context.Background(), topo, ConsumerConfig{LookupAddress: "127.0.0.1:4161"}, nil)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestNewConsumer_invalidConsumerConfig(t *testing.T) {
	topo := &Topology{Topic: "t", Channel: "c"}
	_, err := NewConsumer(context.Background(), topo, ConsumerConfig{}, nsqio.HandlerFunc(func(*nsqio.Message) error { return nil }))
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestNewConsumer_success(t *testing.T) {
	oldConnect := connectLookupd
	defer func() { connectLookupd = oldConnect }()
	connectLookupd = func(*nsqio.Consumer, ConsumerConfig) error { return nil }

	topo := &Topology{Topic: "t", Channel: "c"}
	c, err := NewConsumer(context.Background(), topo, ConsumerConfig{LookupAddress: "127.0.0.1:4161"}, nsqio.HandlerFunc(func(*nsqio.Message) error { return nil }))
	if err != nil {
		t.Fatal(err)
	}
	defer c.Stop()
	var _ Consumer = c
	if c == nil {
		t.Fatal("expected consumer")
	}
}

func TestNewConsumer_connectLookupdError(t *testing.T) {
	oldConnect := connectLookupd
	defer func() { connectLookupd = oldConnect }()
	connectLookupd = func(*nsqio.Consumer, ConsumerConfig) error {
		return errors.New("lookupd failed")
	}
	topo := &Topology{Topic: "t", Channel: "c"}
	_, err := NewConsumer(context.Background(), topo, ConsumerConfig{LookupAddresses: []string{"a", "b"}}, nsqio.HandlerFunc(func(*nsqio.Message) error { return nil }))
	if err == nil {
		t.Fatal("expected error")
	}
}

func Test_defaultConnectLookupd(t *testing.T) {
	c, err := nsqio.NewConsumer("t", "ch", nsqio.NewConfig())
	if err != nil {
		t.Fatal(err)
	}
	defer c.Stop()

	if err := defaultConnectLookupd(c, ConsumerConfig{}); err == nil {
		t.Fatal("expected no addresses error")
	}

	oldOne := connectOneLookupd
	oldMany := connectManyLookupd
	defer func() {
		connectOneLookupd = oldOne
		connectManyLookupd = oldMany
	}()
	var single, multi bool
	connectOneLookupd = func(*nsqio.Consumer, string) error { single = true; return nil }
	connectManyLookupd = func(*nsqio.Consumer, []string) error { multi = true; return nil }

	if err := defaultConnectLookupd(c, ConsumerConfig{LookupAddress: "a"}); err != nil {
		t.Fatal(err)
	}
	if !single || multi {
		t.Fatalf("single=%v multi=%v", single, multi)
	}
	single, multi = false, false
	if err := defaultConnectLookupd(c, ConsumerConfig{LookupAddresses: []string{"a", "b"}}); err != nil {
		t.Fatal(err)
	}
	if single || !multi {
		t.Fatalf("single=%v multi=%v", single, multi)
	}
}

func TestNewConsumer_newConsumerError(t *testing.T) {
	oldNew := newNSQConsumer
	defer func() { newNSQConsumer = oldNew }()
	newNSQConsumer = func(string, string, *nsqio.Config) (*nsqio.Consumer, error) {
		return nil, errors.New("new consumer failed")
	}
	topo := &Topology{Topic: "t", Channel: "c"}
	_, err := NewConsumer(context.Background(), topo, ConsumerConfig{LookupAddress: "127.0.0.1:4161"}, nil)
	if err == nil {
		t.Fatal("expected error")
	}
}
