package nsq_test

import (
	"context"
	"fmt"

	"github.com/andywang2011/go-nsq"
)

func ExampleNewProducer() {
	cfg := nsq.ConfigDefault()
	cfg.Address = "127.0.0.1:65530"

	p, err := nsq.NewProducer(context.Background(), cfg)
	if err != nil {
		fmt.Println("unavailable")
		return
	}
	defer p.Stop()
	fmt.Println("ok")
	// Output: unavailable
}
