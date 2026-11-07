package nsq

import (
	"testing"
	"time"
)

func TestProducerConfig_Validate(t *testing.T) {
	if err := (&ProducerConfig{Address: "127.0.0.1:4150"}).Validate(); err != nil {
		t.Fatal(err)
	}
	if err := (&ProducerConfig{}).Validate(); err == nil {
		t.Fatal("expected error")
	}
}

func TestConsumerConfig_Validate(t *testing.T) {
	if err := (&ConsumerConfig{LookupAddress: "127.0.0.1:4161"}).Validate(); err != nil {
		t.Fatal(err)
	}
	if err := (&ConsumerConfig{LookupAddresses: []string{"h1:4161", "h2:4161"}}).Validate(); err != nil {
		t.Fatal(err)
	}
	if err := (&ConsumerConfig{}).Validate(); err == nil {
		t.Fatal("expected error")
	}
}

func TestTopology_Validate_defaults(t *testing.T) {
	topo := &Topology{Topic: "t", Channel: "c"}
	if err := topo.Validate(); err != nil {
		t.Fatal(err)
	}
	if topo.Concurrency != defaultConcurrency {
		t.Fatalf("Concurrency = %d", topo.Concurrency)
	}
	if topo.MaxAttempts != defaultMaxAttempts {
		t.Fatalf("MaxAttempts = %d", topo.MaxAttempts)
	}
	if topo.RequeueDelay != defaultRequeueDelay {
		t.Fatalf("RequeueDelay = %s", topo.RequeueDelay)
	}
}

func TestTopology_Validate_errors(t *testing.T) {
	if err := (&Topology{Topic: "", Channel: "c", Concurrency: 1}).Validate(); err == nil {
		t.Fatal("expected error")
	}
}

func TestConfigDefault(t *testing.T) {
	c := ConfigDefault()
	if c.DialTimeout != defaultDialTimeout {
		t.Fatal(c.DialTimeout)
	}
}

func TestDefaultTopology(t *testing.T) {
	topo := DefaultTopology()
	if topo.Concurrency != defaultConcurrency {
		t.Fatal(topo.Concurrency)
	}
}

func TestTopology_Validate_custom(t *testing.T) {
	topo := &Topology{
		Topic: "t", Channel: "c", Concurrency: 2, MaxAttempts: 3,
		RequeueDelay: time.Second,
	}
	if err := topo.Validate(); err != nil {
		t.Fatal(err)
	}
	if topo.Concurrency != 2 {
		t.Fatal("custom concurrency overwritten")
	}
}
