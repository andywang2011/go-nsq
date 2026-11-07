// Example NSQ consumer: lookupd discovery, concurrent handlers, graceful stop.
//
// In production, trap SIGTERM/SIGINT and call c.Stop(); wait on <-c.StopChan().
package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/andywang2011/go-nsq"
	nsqio "github.com/nsqio/go-nsq"
)

func main() {
	handler := nsqio.HandlerFunc(func(msg *nsqio.Message) error {
		log.Printf("body=%s", msg.Body)
		return nil
	})

	topo := nsq.DefaultTopology()
	topo.Topic = "events"
	topo.Channel = "worker"

	c, err := nsq.NewConsumer(context.Background(), &topo, nsq.ConsumerConfig{LookupAddress: "127.0.0.1:4161"}, handler)
	if err != nil {
		log.Fatal(err)
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig

	c.Stop()
	<-c.StopChan()
}
