package nsq

import (
	"testing"
	"time"
)

type fakeProducer struct {
	published []string
}

func (f *fakeProducer) Publish(topic string, body []byte) error {
	f.published = append(f.published, topic)
	return nil
}
func (f *fakeProducer) DeferredPublish(topic string, _ time.Duration, body []byte) error {
	return f.Publish(topic, body)
}
func (f *fakeProducer) MultiPublish(topic string, body [][]byte) error {
	for _, b := range body {
		if err := f.Publish(topic, b); err != nil {
			return err
		}
	}
	return nil
}
func (f *fakeProducer) Stop() {}

var _ Producer = (*fakeProducer)(nil)

func TestFakeProducer_Publish(t *testing.T) {
	f := &fakeProducer{}
	if err := f.Publish("t", []byte("x")); err != nil {
		t.Fatal(err)
	}
	if len(f.published) != 1 {
		t.Fatal(f.published)
	}
}
