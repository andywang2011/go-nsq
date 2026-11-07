package nsq

import (
	"testing"
	"time"

	nsqio "github.com/nsqio/go-nsq"
)

func Test_applyProducerConfig_defaults(t *testing.T) {
	cfg := ProducerConfig{Snappy: true}
	c := applyProducerConfig(cfg)
	if c.DialTimeout != defaultDialTimeout {
		t.Fatalf("DialTimeout %s", c.DialTimeout)
	}
	if !c.Snappy {
		t.Fatal("expected snappy")
	}
}

func Test_applyProducerConfig_override(t *testing.T) {
	custom := nsqio.NewConfig()
	cfg := ProducerConfig{NSQConfig: custom}
	if applyProducerConfig(cfg) != custom {
		t.Fatal("expected override")
	}
}

func Test_applyConsumerConfig(t *testing.T) {
	topo := &Topology{Topic: "t", Channel: "c", Concurrency: 4, MaxAttempts: 2, RequeueDelay: time.Second}
	_ = topo.Validate()
	cfg := ConsumerConfig{LookupAddress: "127.0.0.1:4161", LookupdPollInterval: 30 * time.Second}
	c := applyConsumerConfig(topo, cfg)
	if c.MaxInFlight != 4 {
		t.Fatalf("MaxInFlight %d", c.MaxInFlight)
	}
	if c.LookupdPollInterval != 30*time.Second {
		t.Fatal(c.LookupdPollInterval)
	}
}

func Test_applyProducerConfig_customTimeouts(t *testing.T) {
	cfg := ProducerConfig{
		DialTimeout:  2 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 4 * time.Second,
		Deflate:      true,
	}
	c := applyProducerConfig(cfg)
	if c.DialTimeout != 2*time.Second || c.ReadTimeout != 3*time.Second || c.WriteTimeout != 4*time.Second {
		t.Fatal(c)
	}
	if !c.Deflate {
		t.Fatal("expected deflate")
	}
}

func Test_applyConsumerConfig_override(t *testing.T) {
	custom := nsqio.NewConfig()
	topo := &Topology{Topic: "t", Channel: "c"}
	_ = topo.Validate()
	c := applyConsumerConfig(topo, ConsumerConfig{NSQConfig: custom})
	if c != custom {
		t.Fatal("expected override")
	}
}

func Test_lookupAddresses_empty(t *testing.T) {
	if lookupAddresses(ConsumerConfig{}) != nil {
		t.Fatal("expected nil")
	}
}

func Test_lookupAddresses(t *testing.T) {
	multi := lookupAddresses(ConsumerConfig{LookupAddresses: []string{"a", "b"}})
	if len(multi) != 2 {
		t.Fatal(multi)
	}
	single := lookupAddresses(ConsumerConfig{LookupAddress: "a"})
	if len(single) != 1 || single[0] != "a" {
		t.Fatal(single)
	}
}
