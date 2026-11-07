package nsq

import (
	nsqio "github.com/nsqio/go-nsq"
)

func applyProducerConfig(cfg ProducerConfig) *nsqio.Config {
	if cfg.NSQConfig != nil {
		return cfg.NSQConfig
	}
	c := nsqio.NewConfig()
	if cfg.DialTimeout > 0 {
		c.DialTimeout = cfg.DialTimeout
	} else {
		c.DialTimeout = defaultDialTimeout
	}
	if cfg.ReadTimeout > 0 {
		c.ReadTimeout = cfg.ReadTimeout
	} else {
		c.ReadTimeout = defaultReadTimeout
	}
	if cfg.WriteTimeout > 0 {
		c.WriteTimeout = cfg.WriteTimeout
	} else {
		c.WriteTimeout = defaultWriteTimeout
	}
	c.Snappy = cfg.Snappy
	c.Deflate = cfg.Deflate
	return c
}

func applyConsumerConfig(topo *Topology, cfg ConsumerConfig) *nsqio.Config {
	var c *nsqio.Config
	if cfg.NSQConfig != nil {
		c = cfg.NSQConfig
	} else {
		c = nsqio.NewConfig()
	}
	c.MaxAttempts = topo.MaxAttempts
	c.MaxInFlight = topo.Concurrency
	c.MaxBackoffDuration = 0
	c.DefaultRequeueDelay = topo.RequeueDelay
	if cfg.LookupdPollInterval > 0 {
		c.LookupdPollInterval = cfg.LookupdPollInterval
	}
	return c
}

func lookupAddresses(cfg ConsumerConfig) []string {
	if len(cfg.LookupAddresses) > 0 {
		return cfg.LookupAddresses
	}
	if cfg.LookupAddress != "" {
		return []string{cfg.LookupAddress}
	}
	return nil
}
