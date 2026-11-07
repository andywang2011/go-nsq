// Example NSQ producer: connect, publish once, stop.
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/andywang2011/go-nsq"
)

func main() {
	p, err := nsq.NewProducer(context.Background(), nsq.ProducerConfig{Address: "127.0.0.1:4150"})
	if err != nil {
		log.Fatal(err)
	}
	defer p.Stop()

	if err := p.Publish("events", []byte("hello")); err != nil {
		log.Fatal(err)
	}
	fmt.Println("published")
}
