# Realtime Bidding System Engine

Lightweight realtime bidding (RTB) engine used for learning and experiments.  
It simulates DSP-style bid handling and can be extended with Redis Streams (ingest), Kafka (async delivery) and Prometheus/Grafana (monitoring).

## Current features
- Priority-based bid queue (priority heap)
- Concurrent processing with a configurable concurrency limit
- Mock external bidder HTTP server (cmd/mockbidder) that simulates latency/failures
- Kafka integration example: write bid events to a Kafka topic
- Redis Streams reader helper for ingestion (utils)
- Prometheus metrics helpers (histogram + counters) and a small helper to expose `/metrics`
- Unit test harness for engine (pkg/engine/engine_test.go)

## Repo layout
- cmd/app — main application (runs engine)
- cmd/mockbidder — mock HTTP bidder (simulates latency/failure)
- pkg/engine — engine logic and orchestration
- pkg/utils — shared types, helpers (GenerateBids, Redis/Kafka helpers, Prometheus helpers)

## Quickstart (Windows)

Prereqs:
- Go 1.20+
- (optional) Docker for Kafka/Redis/Prometheus

1. Run the mock bidder (graceful shutdown via context + signals):
   ```sh
   go run ./cmd/mockbidder
   ```
   Default endpoint: http://localhost:8081/bid

2. Run the main app (notice new RunEngineWithBids signature that accepts a context and bidder URL):
   ```go
   // example in cmd/app/main.go
   utils.InitPrometheus()
   ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
   defer cancel()
   bids := utils.GenerateBids(10)
   engine.RunEngineWithBids(ctx, bids, 5, 20, "http://localhost:8081/bid")
   ```

3. Run tests (tests use httptest.Server and inject its URL into the engine):
   ```sh
   go test ./pkg/engine -v
   ```

## Mock bidder

- The mock bidder (cmd/mockbidder) now uses signal.NotifyContext for clean shutdown on SIGINT/SIGTERM.
- Use it for local integration testing or run tests which use httptest.Server instead.

## Kafka & Redis

- utils provides helpers:
  - ReadBidsFromRedisStream(ctx, rdb, stream, out) — Redis Streams ingestion helper.
  - SendBidResultToKafka(writer *kafka.Writer, bidID int, status string) — convenience wrapper.

Notes:
- SendBidResultToKafka now early-returns when writer == nil, so the code is safe when Kafka is not configured or unavailable.
- The engine orchestration (RunEngineWithBids) accepts a context and bidderURL, and creates/uses a kafka.Writer (configured in code). Adjust brokers/topic in engine if you need a different setup.
- For local testing, use the provided docker-compose to bring up Kafka / Redis.

## Prometheus & Metrics

- Call utils.InitPrometheus() early (e.g., in cmd/app/main.go) to expose /metrics (default :2112).
- RunEngineWithBids now aggregates results via a results channel carrying duration + success flag, and updates Prometheus metrics:
  - bid_latency_seconds (histogram)
  - bid_success_total (counter)
  - bid_timeout_total (counter)

Example call flow:
```go
// main.go
utils.InitPrometheus()
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()
engine.RunEngineWithBids(ctx, bids, 5, 20, "http://localhost:8081/bid")
```

## Testing / Integration

- Unit/integration tests in pkg/engine/engine_test.go spin up an httptest.Server that simulates latency/failures and pass server.URL+"/bid" to RunEngineWithBids. This keeps tests hermetic and fast.

## Troubleshooting

- "not enough arguments in call to kafka.NewWriter": call kafka.NewWriter(kafka.WriterConfig{...}) — include Brokers and Topic.
- SendBidResultToKafka will no-op (return nil) if the provided *kafka.Writer is nil.
- If Prometheus avg latency divides by zero, ensure there are successful counts before computing averages (utils.Metrics.Report does this).
- If Kafka/Redis connection issues occur, verify Docker compose or your broker endpoints.

## Next steps / Extensions

- Use ReadBidsFromRedisStream to ingest real requests into the engine.
- Send processed bid events to Kafka (SendBidResultToKafka) for downstream analytics.
- Import Prometheus metrics into Grafana and build dashboards for latency / success rates

## Contact / Contribution
Open issues or PRs with small incremental changes (fixes, metrics, test coverage). Keep code simple and platform-agnostic where possible.
