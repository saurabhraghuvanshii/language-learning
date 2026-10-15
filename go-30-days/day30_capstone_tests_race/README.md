# Day 30: testing the inventory API

`main.go` is a self-contained copy of the day 29 API (Store, MemStore, handlers, middleware, graceful shutdown). It is not imported by any other folder.

## How to run
- Demo (server + client, then exit): `go run ./day30_capstone_tests_race`
- Tests with the race detector: `go test -race ./day30_capstone_tests_race`
- Verbose: `go test -race -v ./day30_capstone_tests_race`
- Benchmarks: `go test -run xxx -bench . -benchmem ./day30_capstone_tests_race`

## What each test proves
| Test | Proves |
|------|--------|
| `TestHandlersTable` | Every route returns the right status and body (200, 201, 204, 400, 404, 405, 422), using `httptest.NewRecorder`. Each case has a fresh store. |
| `TestEndToEndServer` | A real HTTP round trip works over TCP: routing, JSON, `Location` and `Content-Type` headers (`httptest.NewServer`). |
| `TestStoreConcurrent` | 50 goroutines create/get/list/delete at once with no data race and consistent results. Delete the locks in `MemStore` and `-race` fails. |
| `BenchmarkGetItem` | Cost of one GET through middleware, mux and handler. |
| `BenchmarkStoreParallel` | Read throughput of the RWMutex store under parallel load. |
