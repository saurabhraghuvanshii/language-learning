/*
DAY 24: SELECT, TIMEOUTS, TIMERS AND TICKERS
RUN:  go run ./day24_select_timeouts

THEORY (plain words):
select is a switch statement for channel operations. Each case is a send or a
receive. select waits until at least one case can proceed, runs that one, and
moves on. Under the hood the runtime locks all the channels involved, scans the
cases in RANDOM order for a ready one, and if none is ready parks the goroutine on
every channel at once until one of them wakes it.
  - default case      makes select non-blocking: runs immediately if nothing is ready.
  - time.After(d)     returns a channel that receives one value after d. Perfect
    as a timeout case. (Since Go 1.23 unreferenced timers are
    garbage collected, so the old "leak" worry is mostly gone.)
  - time.NewTimer(d)  like After but you get a *Timer with Stop() and Reset().
  - time.NewTicker(d) delivers a value every d until you call Stop().

REAL WORLD: calling a slow upstream API with a deadline, health-check heartbeats,
periodic cache refresh, waiting for "work OR shutdown signal".
EXPLAIN-IT ANSWER: What does select do when several cases are ready?
It picks ONE of the ready cases uniformly at random (pseudo-random), runs it, and
ignores the others for that pass. This is deliberate: a fixed order (like top to
bottom) would starve the later cases. Do not rely on order. If you need priority,
check the high-priority channel first in its own select with a default.
PITFALLS:
 1. Assuming the first ready case wins (it does not; shown below).
 2. Forgetting ticker.Stop(): the ticker keeps the runtime busy and goroutines leak.
 3. A select with only a nil channel blocks forever; `select {}` blocks forever.
 4. Giving up on a call does NOT stop the callee goroutine. Use a buffered result
    channel (cap 1) so the abandoned goroutine can still finish its send and exit.

EXERCISES:
 1. TODO: change the timeout to 500ms and observe the call now succeeds.
 2. TODO: make the heartbeat interval 25ms and count how many beats you get.
 3. TODO: use time.NewTimer inside a loop and call Reset(d) after each message
    to build an "idle timeout" (give up after 100ms of silence).
*/
package main

import (
	"errors"
	"fmt"
	"time"
)

// demoRandomPick: both channels are ready; select picks randomly.
func demoRandomPick() {
	fmt.Println("--- select with several ready cases ---")
	a := make(chan string, 1000)
	b := make(chan string, 1000)
	for i := 0; i < 1000; i++ { // fill both so both are always ready
		a <- "a"
		b <- "b"
	}
	counts := map[string]int{}
	for i := 0; i < 1000; i++ {
		select {
		case v := <-a:
			counts[v]++
		case v := <-b:
			counts[v]++
		}
	}
	// Expect roughly 500/500, never 1000/0.
	fmt.Println("picked from a:", counts["a"], "picked from b:", counts["b"])
}

// demoDefault: default makes the select non-blocking.
func demoDefault() {
	fmt.Println("--- default: non-blocking receive ---")
	queue := make(chan string, 1)
	select {
	case msg := <-queue:
		fmt.Println("got", msg)
	default: // nothing ready right now, so do something else
		fmt.Println("queue empty, not waiting")
	}
}

// fetchPrice simulates a slow upstream pricing API.
func fetchPrice(sku string, latency time.Duration) (float64, error) {
	time.Sleep(latency)
	return 19.99, nil
}

// ErrTimeout is returned when the upstream takes too long.
var ErrTimeout = errors.New("upstream timed out")

// fetchWithTimeout runs fetchPrice in a goroutine and gives up after limit.
func fetchWithTimeout(sku string, latency, limit time.Duration) (float64, error) {
	type result struct {
		price float64
		err   error
	}
	// Buffer of 1: if we time out and leave, the goroutine can still send and
	// exit instead of blocking forever (that would be a goroutine leak).
	ch := make(chan result, 1)
	go func() {
		p, err := fetchPrice(sku, latency)
		ch <- result{p, err}
	}()

	timer := time.NewTimer(limit) // NewTimer so we can Stop it and release resources early
	defer timer.Stop()

	select {
	case r := <-ch:
		return r.price, r.err
	case <-timer.C: // fires once after limit
		return 0, ErrTimeout
	}
}

func demoTimeout() {
	fmt.Println("--- timeout: give up after 200ms ---")
	for _, latency := range []time.Duration{50 * time.Millisecond, 400 * time.Millisecond} {
		t0 := time.Now()
		price, err := fetchWithTimeout("SKU-1", latency, 200*time.Millisecond)
		fmt.Printf("upstream latency %v -> price=%.2f err=%v (waited %v)\n",
			latency, price, err, time.Since(t0).Round(10*time.Millisecond))
	}
	// time.After is the shortest form of the same idea:
	select {
	case <-time.After(30 * time.Millisecond):
		fmt.Println("time.After fired")
	}
}

// heartbeat prints a beat every interval until stop is closed, then signals done.
func heartbeat(interval time.Duration, stop <-chan struct{}, done chan<- struct{}) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop() // ALWAYS stop tickers
	defer close(done)   // tell the caller we've fully exited
	beats := 0
	for {
		select {
		case <-ticker.C:
			beats++
			fmt.Println("heartbeat", beats)
		case <-stop: // closing a channel wakes every receiver
			fmt.Println("heartbeat: stop requested, exiting after", beats, "beats")
			return
		}
	}
}

func demoHeartbeat() {
	fmt.Println("--- heartbeat ticker that stops cleanly ---")
	stop := make(chan struct{}) // struct{} carries no data, just a signal
	done := make(chan struct{})
	go heartbeat(60*time.Millisecond, stop, done)
	time.Sleep(200 * time.Millisecond)
	close(stop) // broadcast stop
	<-done      // wait until heartbeat has actually exited
	fmt.Println("heartbeat goroutine is gone")
}

func main() {
	demoRandomPick()
	demoDefault()
	demoTimeout()
	demoHeartbeat()
}
