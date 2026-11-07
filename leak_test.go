package nsq

import (
	"testing"

	"go.uber.org/goleak"
)

func TestMain(m *testing.M) {
	goleak.VerifyTestMain(m,
		goleak.IgnoreTopFunction("github.com/nsqio/go-nsq.(*Consumer).rdyLoop"),
		goleak.IgnoreTopFunction("github.com/nsqio/go-nsq.(*Consumer).handlerLoop"),
	)
}
