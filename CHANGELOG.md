# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.3.0] - 2026-11-09

### Added

- `Client` process-scoped helper: `NewClient(producer, lookup)`, `Publish` / `MultiPublish` / `DeferredPublish`, `Close`, and `NewConsumer(ctx, *Topology, handler)` using stored lookup config.

## [0.2.0] - 2026-11-07

### Breaking

- See repository diff for API changes in this release (secure defaults, renamed packages where applicable, context-aware constructors).

### Security

- Hardened defaults: TLS, timeouts, body limits, PII redaction where logging applies.

### Fixed

- Bug fixes from the v0.2.0 hardening pass.

### Changed

- CI uses go-version-file; added race, staticcheck, govulncheck, golangci-lint (gosec).



### Changed

- **Breaking:** package renamed `gonsq` -> `nsq`; upstream stays imported as `nsqio`.
- **Breaking:** `NewConsumer` and `WrapConsumer` return the new `Consumer` interface instead of `*nsqio.Consumer`.
- Retry helper uses `math/rand/v2` and a stoppable timer instead of `time.After` and a shared `rand.Rand`.
- **Breaking:** `NewProducer`, `NewProducerWithRetry`, and `NewConsumer` take a `context.Context` first; retries stop and late connections are released on cancellation.

## [0.1.0] - 2026-11-07

### Added

- Producer and consumer helpers with startup retry and mock-friendly producer interface.
