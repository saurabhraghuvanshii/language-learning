# Go in 30 Days: runnable practice

One folder per day. Every `main.go` starts with a comment block: THEORY, REAL WORLD, the answer to that day's EXPLAIN-IT question, PITFALLS, and 3 EXERCISES. The code is commented line by line.

## Setup (once)

1. Install Go 1.22 or newer from https://go.dev/dl/ (built and tested on Go 1.24).
2. Unzip this folder and open a terminal inside `go-30-days`.
3. Check: `go version`

## Daily loop (about 90 minutes)

1. Read the comment block at the top of the day's `main.go`.
2. Run it: `go run ./day01_hello_toolchain` (replace with the day's folder).
3. Read the code and the output side by side, line by line.
4. Do the 3 EXERCISES by editing the file. Break things on purpose and read the error.
5. Explain every line out loud as if teaching a friend. If you cannot say why a line is needed, read it again.

## Index

| Day | Folder | Topic |
| --- | --- | --- |
| 1 | day01_hello_toolchain | toolchain, packages, `main` |
| 2 | day02_variables_zero_values | variables, constants, zero values |
| 3 | day03_types_strings_runes | types, strings, bytes, runes |
| 4 | day04_control_flow | if, switch, the one `for` |
| 5 | day05_functions_errors | functions, multiple returns, errors |
| 6 | day06_defer_panic_closures | defer, panic, recover, closures |
| 7 | day07_mini_todo_cli | mini project |
| 8 | day08_arrays | arrays are values |
| 9 | day09_slice_header | slice = pointer + len + cap |
| 10 | day10_slice_sharing_append | append, sharing, copy |
| 11 | day11_maps | maps |
| 12 | day12_structs_methods | structs, receivers, embedding |
| 13 | day13_pointers_pass_by_value | pointers, pass by value |
| 14 | day14_mini_inventory | mini project |
| 15 | day15_interfaces | interfaces, `io.Writer` |
| 16 | day16_interface_internals | nil interface bug |
| 17 | day17_errors | wrapping, `Is`, `As` |
| 18 | day18_packages_layout | packages, `internal/`, no circular imports |
| 19 | day19_decoupling_interfaces | consumer-side interfaces |
| 20 | day20_generics | generics |
| 21 | day21_testing | tests, benchmarks |
| 22 | day22_goroutines_waitgroup | goroutines, concurrent HTTP calls |
| 23 | day23_channels | channels |
| 24 | day24_select_timeouts | select, timeouts, tickers |
| 25 | day25_mutex_data_races | **when to use a mutex** |
| 26 | day26_context | context, cancellation |
| 27 | day27_concurrency_patterns | worker pool, pipeline |
| 28 | day28_goroutine_leaks | goroutine leaks, scheduler |
| 29 | day29_capstone_http_api | **how an HTTP request works in Go** |
| 30 | day30_capstone_tests_race | tests with `-race` |

## Useful commands

- Run a day: `go run ./dayNN_name`
- Check for mistakes: `go vet ./...`
- Format code: `gofmt -w .`
- Find data races: `go run -race ./day25_mutex_data_races` (the race report on section 1 is intentional)
- Run tests: `go test -race ./day21_testing ./day30_capstone_tests_race`
- Keep the Day 29 API running: `go run ./day29_capstone_http_api -serve`, then `curl http://localhost:8080/items`

## Your two main questions

- **When to use a mutex, and when concurrency at all:** Day 25 has a decision table in its header comment, plus a broken counter that really loses updates so you can see the problem.
- **How an HTTP request is created and served:** Day 22 (calling many endpoints concurrently), Day 26 (timeouts with `context`) and Day 29 (a full REST API plus a client, with the request's path traced step by step).
