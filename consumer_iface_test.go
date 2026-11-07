package nsq

import (
	"testing"

	nsqio "github.com/nsqio/go-nsq"
)

func TestWrapConsumer_nil(t *testing.T) {
	if WrapConsumer(nil) != nil {
		t.Fatal("expected nil")
	}
}

func TestWrapConsumer_stop(t *testing.T) {
	c, err := nsqio.NewConsumer("t", "ch", nsqio.NewConfig())
	if err != nil {
		t.Fatal(err)
	}
	h := WrapConsumer(c)
	h.Stop()
}

func TestWrapConsumer_methods(t *testing.T) {
	c, err := nsqio.NewConsumer("t", "ch", nsqio.NewConfig())
	if err != nil {
		t.Fatal(err)
	}
	w := WrapConsumer(c)
	if w.StopChan() == nil {
		t.Fatal("expected stop chan")
	}
	w.ChangeMaxInFlight(5)
	_ = w.IsStarved()
	if w.Stats() == nil {
		t.Fatal("expected stats")
	}
	w.Stop()
}

func TestConsumerWrapper_nilReceiverAndNilInner(t *testing.T) {
	for _, w := range []*consumerWrapper{nil, {}} {
		w.Stop()
		w.ChangeMaxInFlight(1)
		if w.StopChan() != nil || w.IsStarved() || w.Stats() != nil {
			t.Fatal("expected zero values")
		}
	}
}
