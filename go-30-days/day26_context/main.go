/*
DAY 26: CONTEXT
RUN:  go run ./day26_context

THEORY (plain words):
A context.Context carries three things across API and goroutine boundaries: a
CANCELLATION signal, an optional DEADLINE, and a few request-scoped VALUES. It is an
interface: Done() <-chan struct{} (closed when cancelled), Err() (context.Canceled
or context.DeadlineExceeded), Deadline(), Value(key). Contexts form a TREE: you
derive a child from a parent; cancelling the parent cancels every child, but never
the other way round.
  - context.Background()          the empty root; use in main, tests, top level.
  - context.TODO()                placeholder when you haven't decided yet.
  - WithCancel(parent)            child + cancel func you call to stop it.
  - WithTimeout(parent, d)        cancels itself after d.
  - WithDeadline(parent, t)       cancels itself at time t (Timeout = Deadline(now+d)).
  - WithValue(parent, key, val)   attaches a value. Use SPARINGLY: request IDs,
    trace IDs, auth identity. NOT for optional function
    parameters or dependencies (DB handles, loggers).
    Use an unexported custom key TYPE to avoid clashes.

REAL WORLD: every incoming HTTP request has r.Context(), cancelled when the client
disconnects; you pass it to database calls (db.QueryContext) and outgoing HTTP calls
so work stops when nobody is waiting for the answer. Servers use it for graceful shutdown.
EXPLAIN-IT ANSWER: Why is ctx always the first parameter?
Pure convention, but a universal one (go vet / linters check it): func Do(ctx
context.Context, ...). Because nearly every function that does I/O or waits needs
it, putting it in the same place makes call chains easy to read and easy to thread
through, makes it obvious in a signature that a function can be cancelled, and means
you never store it in a struct (a context is per-call, like a request, not per-object).
PITFALLS:
 1. Not calling cancel(): the timer and the child registration leak until the parent
    ends. ALWAYS `defer cancel()` right after creating the context. go vet checks this.
 2. Storing ctx in a struct instead of passing it.
 3. Passing a nil context. Pass context.TODO() if unsure.
 4. A function that ignores ctx.Done() cannot be cancelled. Cancellation is
    cooperative: your code must select on ctx.Done() or use ctx-aware APIs.
 5. Putting business data in WithValue: it is untyped and invisible in signatures.

EXERCISES:
 1. TODO: change the WithTimeout in demoHTTP from 100ms to 500ms and watch the request succeed.
 2. TODO: use context.WithDeadline(ctx, time.Now().Add(...)) in place of WithTimeout.
 3. TODO: add a second value (user ID) with its own key type and print both.
*/
package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"sync"
	"time"
)

// ctxKey is an unexported type so other packages' keys can never collide with ours.
type ctxKey string

const requestIDKey ctxKey = "requestID"

// requestID reads the request ID back out of a context.
func requestID(ctx context.Context) string {
	id, ok := ctx.Value(requestIDKey).(string) // type assertion: Value returns any
	if !ok {
		return "none"
	}
	return id
}

// handleOrder shows ctx as FIRST parameter and reading a value from it.
func handleOrder(ctx context.Context, orderID int) {
	fmt.Printf("  [req %s] processing order %d\n", requestID(ctx), orderID)
}

func demoValue() {
	fmt.Println("--- WithValue: request ID (use sparingly) ---")
	ctx := context.WithValue(context.Background(), requestIDKey, "req-42")
	handleOrder(ctx, 1001)
}

// worker processes jobs until ctx is cancelled. ctx first, as always.
func worker(ctx context.Context, id int, jobs <-chan int, wg *sync.WaitGroup, processed chan<- string) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done(): // cancelled or timed out: stop promptly
			processed <- fmt.Sprintf("worker %d stopping: %v", id, ctx.Err())
			return
		case j, ok := <-jobs:
			if !ok {
				return
			}
			time.Sleep(10 * time.Millisecond) // pretend to resize an image
			_ = j
		}
	}
}

func demoWorkerPool() {
	fmt.Println("--- WithCancel: worker pool stops on cancel ---")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // safety net even if we cancel explicitly below

	jobs := make(chan int)
	msgs := make(chan string, 3) // buffered so workers never block on reporting
	var wg sync.WaitGroup
	for i := 1; i <= 3; i++ {
		wg.Add(1)
		go worker(ctx, i, jobs, &wg, msgs)
	}
	for j := 0; j < 5; j++ {
		jobs <- j
	}
	cancel()  // tell everyone to stop
	wg.Wait() // wait until they all have
	close(msgs)
	n := 0
	for range msgs {
		n++
	}
	fmt.Println("all workers stopped, stop messages received:", n)
}

func demoTimeoutDeadline() {
	fmt.Println("--- WithTimeout / WithDeadline ---")
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()
	select {
	case <-ctx.Done():
		fmt.Println("timeout ctx:", ctx.Err(), "| is DeadlineExceeded:", errors.Is(ctx.Err(), context.DeadlineExceeded))
	case <-time.After(time.Second):
		fmt.Println("unexpected")
	}
	dctx, dcancel := context.WithDeadline(context.Background(), time.Now().Add(20*time.Millisecond))
	defer dcancel()
	dl, _ := dctx.Deadline()
	fmt.Println("deadline is in the future:", time.Until(dl) > 0)
	<-dctx.Done()
	fmt.Println("deadline ctx:", dctx.Err())
}

// fetch performs a GET bound to ctx. When ctx expires, the in-flight request is aborted.
func fetch(ctx context.Context, client *http.Client, url string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", err // includes context deadline exceeded
	}
	defer resp.Body.Close() // always close the body or the connection leaks
	b, err := io.ReadAll(resp.Body)
	return string(b), err
}

func demoHTTP() {
	fmt.Println("--- HTTP call that times out (httptest server sleeps 300ms) ---")
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		select {
		case <-time.After(300 * time.Millisecond): // slow upstream
			fmt.Fprint(w, "price=19.99")
		case <-r.Context().Done(): // server notices the client went away
		}
	}))
	defer srv.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	t0 := time.Now()
	_, err := fetch(ctx, srv.Client(), srv.URL)
	fmt.Printf("error: %v\n", err)
	fmt.Println("is deadline exceeded:", errors.Is(err, context.DeadlineExceeded), "| took ~", time.Since(t0).Round(50*time.Millisecond))

	body, err := fetch(context.Background(), srv.Client(), srv.URL) // no deadline: waits full 300ms
	fmt.Printf("without deadline: body=%q err=%v\n", body, err)
}

func main() {
	demoValue()
	demoWorkerPool()
	demoTimeoutDeadline()
	demoHTTP()
}
