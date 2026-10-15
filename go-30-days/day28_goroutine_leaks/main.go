/*
DAY 28: GOROUTINE LEAKS AND THE SCHEDULER
RUN:  go run ./day28_goroutine_leaks

THEORY (plain words):
A goroutine is a function running on its own tiny stack (starts ~2KB, grows as
needed). The runtime schedules many goroutines onto few OS threads. Vocabulary:

	G = a goroutine (stack + the function's state).
	M = a "machine": a real OS thread. Only an M can actually run code.
	P = a "processor": a scheduling token holding a local run queue of Gs. There are
	    GOMAXPROCS of them (default = number of CPU cores). An M must hold a P to run Gs.

Flow: new goroutines go onto a P's local queue; its M runs them one by one. If a G
blocks on a channel/mutex/sleep, it is parked (set aside) and the M immediately runs
the next G: blocking is cheap. If a G makes a blocking system call, the P is handed
to another M so other work continues. WORK STEALING: when a P's queue is empty it
steals half of another P's queue (or checks the global queue / network poller), so
cores don't sit idle. PREEMPTION: since Go 1.14 the runtime can interrupt a G that
runs too long (about 10ms) via a signal, so one busy loop cannot hog a core.
Unlike threads, nobody ever kills a goroutine from outside: it ends only when its
function returns (or the whole program exits).
REAL WORLD: leaked goroutines in servers show up as steadily growing memory and
goroutine counts (pprof /debug/pprof/goroutine), eventually OOM kills. Typical
sources: a handler spawning a goroutine to call a slow API, then returning on timeout.
EXPLAIN-IT ANSWER: How does a forgotten goroutine stay alive forever?
A goroutine only ends when its function returns. If it is blocked on a channel
send/receive (or a lock, or a select) and no other goroutine will ever complete that
operation, it never returns. The runtime does not garbage-collect it: the channel is
reachable FROM the goroutine's own stack, and the goroutine is reachable from the
scheduler's bookkeeping, so GC sees it as live. Its stack and everything it points
to stay in memory until the process exits.
THREE CLASSIC LEAK CAUSES:
 1. Sending on an unbuffered channel nobody reads (the receiver left after a timeout
    or an early return). Fix: buffered result channel of 1, or select with ctx.Done().
 2. Receiving from / ranging over a channel nobody ever closes. Fix: the sender
    closes it, or the loop also selects on a done signal.
 3. Forgotten resources: time.Ticker not Stopped, a goroutine waiting on a
    context that never cancels (missing defer cancel()), or wg.Wait() / mutex
    Lock() that can never complete.

PITFALLS: every `go` statement needs an answer to "when and how does this stop?".

	Counting with runtime.NumGoroutine() is a debugging aid, not production logic.

EXERCISES:
 1. TODO: print runtime.GOMAXPROCS(0) and runtime.NumCPU() at startup.
 2. TODO: add a third leak variant (range over a never-closed channel), measure it, then fix it.
 3. TODO: replace the done channel in the fixed version with context.WithCancel.
*/
package main

import (
	"fmt"
	"runtime"
	"time"
)

// leakyFetch simulates an upstream call whose result nobody wants anymore.
// LEAKY: if the caller leaves, the send blocks forever.
func leakyFetch(timeout time.Duration) (string, bool) {
	ch := make(chan string) // unbuffered and nobody will read after timeout
	go func() {
		time.Sleep(50 * time.Millisecond) // slow upstream
		ch <- "price=19.99"               // blocks forever once the caller is gone
	}()
	select {
	case v := <-ch:
		return v, true
	case <-time.After(timeout): // caller gives up after 10ms
		return "", false
	}
}

// fixedFetch uses a done channel so the worker can always exit.
func fixedFetch(timeout time.Duration) (string, bool) {
	ch := make(chan string)
	done := make(chan struct{})
	defer close(done) // when fixedFetch returns for ANY reason, release the worker
	go func() {
		time.Sleep(50 * time.Millisecond)
		select {
		case ch <- "price=19.99":
		case <-done: // caller is gone: just exit
		}
	}()
	select {
	case v := <-ch:
		return v, true
	case <-time.After(timeout):
		return "", false
	}
}

func main() {
	fmt.Println("goroutines at start:", runtime.NumGoroutine()) // main only: 1

	fmt.Println("--- leaky version: 100 calls that time out ---")
	for i := 0; i < 100; i++ {
		leakyFetch(10 * time.Millisecond)
	}
	time.Sleep(200 * time.Millisecond) // give workers time to reach their blocked send
	fmt.Println("goroutines after leaky calls:", runtime.NumGoroutine(), "(100 stuck forever)")

	before := runtime.NumGoroutine()
	fmt.Println("--- fixed version: 100 calls that time out ---")
	for i := 0; i < 100; i++ {
		fixedFetch(10 * time.Millisecond)
	}
	time.Sleep(200 * time.Millisecond) // workers wake, see done, exit
	after := runtime.NumGoroutine()
	fmt.Println("goroutines before fixed calls:", before, "after:", after, "(no growth)")
	fmt.Println("scheduler: GOMAXPROCS (P count) =", runtime.GOMAXPROCS(0), "| NumCPU =", runtime.NumCPU())
}
