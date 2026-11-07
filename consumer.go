package nsq

import (
	"context"
	"fmt"

	nsqio "github.com/nsqio/go-nsq"
)

var (
	newNSQConsumer     = nsqio.NewConsumer
	connectOneLookupd  = func(c *nsqio.Consumer, addr string) error { return c.ConnectToNSQLookupd(addr) }
	connectManyLookupd = func(c *nsqio.Consumer, addrs []string) error { return c.ConnectToNSQLookupds(addrs) }
	connectLookupd     = defaultConnectLookupd
)

func defaultConnectLookupd(consumer *nsqio.Consumer, cfg ConsumerConfig) error {
	addrs := lookupAddresses(cfg)
	if len(addrs) == 0 {
		return fmt.Errorf("nsq: no lookup addresses")
	}
	if len(addrs) == 1 {
		return connectOneLookupd(consumer, addrs[0])
	}
	return connectManyLookupd(consumer, addrs)
}

// NewConsumer creates a consumer from topology and lookup config and returns it as a Consumer interface.
// ctx bounds the lookupd connect; on cancellation the partially built consumer is stopped.
func NewConsumer(ctx context.Context, topo *Topology, cfg ConsumerConfig, handler nsqio.Handler) (Consumer, error) {
	ctx = ctxOrBackground(ctx)
	if topo == nil {
		return nil, fmt.Errorf("nsq: topology is nil")
	}
	if err := topo.Validate(); err != nil {
		return nil, err
	}
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	nsqCfg := applyConsumerConfig(topo, cfg)
	consumer, err := newNSQConsumer(topo.Topic, topo.Channel, nsqCfg)
	if err != nil {
		return nil, fmt.Errorf("nsq new consumer: %w", err)
	}
	consumer.AddConcurrentHandlers(handler, topo.Concurrency)

	_, err = buildWithCtx(ctx, func() (struct{}, error) {
		return struct{}{}, connectLookupd(consumer, cfg)
	}, func(struct{}) { consumer.Stop() })
	if err != nil {
		consumer.Stop()
		return nil, fmt.Errorf("nsq connect lookupd: %w", err)
	}
	return WrapConsumer(consumer), nil
}
