// Package nsq provides a thin, testable wrapper around NSQ producers and consumers.
//
// Design constraints (for use in application frameworks):
//   - No viper, no global singletons, no log.Fatal inside the library.
//   - Producers return a mock-friendly Producer interface.
//   - Consumers return a mock-friendly Consumer interface (Stop, StopChan, Stats, ...); the upstream package is imported as nsqio.
//   - Start consumers from a dedicated cmd process (not inside HTTP servers); see README.
//
// High concurrency: tune Topology.Concurrency (MaxInFlight) and use AddConcurrentHandlers via NewConsumer.
// One Producer per nsqd address; reuse it across goroutines.
package nsq
