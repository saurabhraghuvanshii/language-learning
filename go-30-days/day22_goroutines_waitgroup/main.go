/*
DAY 22: GOROUTINES AND WAITGROUP
RUN:  go run ./day22_goroutines_waitgroup

THEORY:
A goroutine is a function running concurrently, started with the go keyword.
It is NOT an OS thread. It is a small task managed by the Go runtime in user
space, starting with a stack of only a few KB that grows as needed, so you can
have hundreds of thousands of them. The scheduler (in plain words) has three
actors, usually written G, M and P:

	G = a goroutine (the task to run),
	M = an OS thread (the worker that really executes code),
	P = a processor slot (a token holding a queue of runnable goroutines; there
	    are GOMAXPROCS of them, by default the number of CPU cores).

An M must hold a P to run Go code. When a goroutine blocks (network, sleep,
channel), the runtime parks it and runs another G on the same thread, so a few
threads serve many goroutines. main is also a goroutine; when main returns the
whole program exits, even if other goroutines are still running.
sync.WaitGroup waits for a group: wg.Add(n) before starting, wg.Done() when
each finishes (use defer), wg.Wait() blocks until the counter hits zero.
REAL WORLD: servers run one goroutine per request; fan-out calls to several
backends (like below) and wait for all to finish; background workers.
EXPLAIN-IT ANSWER: What is a goroutine and how is it different from an OS thread?
A goroutine is a lightweight function execution scheduled by the Go runtime;
an OS thread is scheduled by the kernel. Goroutines start with ~2-8 KB
stacks (threads reserve about 1-8 MB), switch cheaply without a kernel call,
and many are multiplexed (M:N) onto a few threads. Creating one costs
microseconds or less; blocking one does not block the thread.
LOOP-VARIABLE NOTE: before Go 1.22, "for i, u := range urls { go func(){ use(u) }() }"
shared ONE variable across iterations, so goroutines often saw the last
value. Since Go 1.22 (our go.mod says go 1.24) each iteration has its own i
and u, so capturing is safe. Older code passed them as arguments: go func(u string){...}(u).
PITFALLS:
  - Forgetting wg.Wait(): main exits and the goroutines are killed.
  - Calling wg.Add inside the goroutine (race with Wait); call it before "go".
  - Writing to a shared slice index/map from many goroutines: a slice with one
    distinct index per goroutine is safe; appending or writing a map is NOT.
  - Copying a WaitGroup by value; pass a pointer.

EXERCISES:
 1. Change the handler delays and watch the sequential vs concurrent times.
 2. Add a 6th endpoint that returns HTTP 500 and show it in the results
    (record the status code and error in the Result).
 3. Try appending results to a shared slice with append inside goroutines
    and run "go run -race ./day22_goroutines_waitgroup" to see the race.
*/
package main

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"sync"
	"time"
)

// Result holds the outcome of one fetch.
type Result struct {
	URL      string
	Status   int
	Bytes    int
	Duration time.Duration
	Err      error
}

// fetch downloads one URL and times it.
func fetch(client *http.Client, url string) Result {
	start := time.Now()
	res := Result{URL: url}
	resp, err := client.Get(url)
	if err != nil {
		res.Err = err
		res.Duration = time.Since(start)
		return res
	}
	defer resp.Body.Close() // always close the body or connections leak
	body, err := io.ReadAll(resp.Body)
	res.Status = resp.StatusCode
	res.Bytes = len(body)
	res.Err = err
	res.Duration = time.Since(start)
	return res
}

// fetchSequential does one request after another.
func fetchSequential(client *http.Client, urls []string) []Result {
	results := make([]Result, len(urls))
	for i, u := range urls {
		results[i] = fetch(client, u)
	}
	return results
}

// fetchConcurrent starts one goroutine per URL and waits for all of them.
func fetchConcurrent(client *http.Client, urls []string) []Result {
	results := make([]Result, len(urls)) // allocated up front with the final length
	var wg sync.WaitGroup
	for i, u := range urls {
		wg.Add(1) // count BEFORE starting the goroutine
		go func() {
			defer wg.Done() // runs even if fetch panics
			// i and u are per-iteration variables (Go 1.22+), so this is safe.
			// Each goroutine writes only results[i]: different memory, no data race,
			// and no mutex is needed.
			results[i] = fetch(client, u)
		}()
	}
	wg.Wait() // blocks until every Done has been called
	return results
}

func main() {
	// Five in-process servers, each simulating a backend with a different latency.
	delays := []time.Duration{
		120 * time.Millisecond, 80 * time.Millisecond, 200 * time.Millisecond,
		50 * time.Millisecond, 150 * time.Millisecond,
	}
	names := []string{"users", "orders", "inventory", "pricing", "reviews"}
	urls := make([]string, len(delays))
	for i, d := range delays {
		name := names[i]
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(d) // pretend to do slow work
			fmt.Fprintf(w, `{"service":%q}`, name)
		}))
		defer srv.Close() // shut each server down when main returns
		urls[i] = srv.URL + "/" + name
	}

	client := &http.Client{Timeout: 5 * time.Second} // never use a client without a timeout

	fmt.Println("== Sequential")
	t0 := time.Now()
	seq := fetchSequential(client, urls)
	seqTotal := time.Since(t0)
	for _, r := range seq {
		fmt.Printf("%-40s status=%d bytes=%d took=%v err=%v\n", r.URL, r.Status, r.Bytes, r.Duration.Round(time.Millisecond), r.Err)
	}

	fmt.Println("== Concurrent")
	t1 := time.Now()
	conc := fetchConcurrent(client, urls)
	concTotal := time.Since(t1)
	for _, r := range conc {
		fmt.Printf("%-40s status=%d bytes=%d took=%v err=%v\n", r.URL, r.Status, r.Bytes, r.Duration.Round(time.Millisecond), r.Err)
	}

	fmt.Println("== Comparison")
	fmt.Printf("sequential total: %v (about the SUM of delays, 600ms)\n", seqTotal.Round(10*time.Millisecond))
	fmt.Printf("concurrent total: %v (about the SLOWEST delay, 200ms)\n", concTotal.Round(10*time.Millisecond))
	fmt.Printf("faster concurrent: %t\n", concTotal < seqTotal)
}
