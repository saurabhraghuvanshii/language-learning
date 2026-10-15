/*
DAY 27: CONCURRENCY PATTERNS
RUN:  go run ./day27_concurrency_patterns

THEORY (plain words):
Almost all real concurrent Go code is built from a few shapes:
  - WORKER POOL: N goroutines read jobs from one channel. Bounds concurrency
    (e.g. max 4 DB calls at once) instead of one goroutine per item.
  - PIPELINE: stages connected by channels; each stage is a goroutine (or few) that
    reads from an input channel, transforms, writes to an output channel and CLOSES
    its output when its input is drained. generator -> transform -> sink.
  - FAN-OUT: several goroutines read from the SAME channel to share the work.
  - FAN-IN (merge): several channels are combined into one with a goroutine per
    input plus a WaitGroup that closes the merged channel when all are done.
  - ERRGROUP BY HAND: run tasks, remember the FIRST error, cancel the others via a
    context. (golang.org/x/sync/errgroup does exactly this; we build it with
    sync.WaitGroup + sync.Once + context.WithCancel.)

The golden rule: whoever SENDS on a channel is the one who closes it, and every
goroutine you start must have a clear way to stop.
REAL WORLD: log/ETL processing, image thumbnailing, calling 10 microservices in
parallel for one page, crawling URLs, batch-importing CSV rows.
EXPLAIN-IT ANSWER: How do you stop all goroutines cleanly?
(1) Give every goroutine a stop signal: a context (ctx.Done()) or a closed "done"
channel; every blocking send/receive sits in a select with that signal.
(2) Let the producer close its output channel when finished, so downstream
`range` loops end naturally. (3) Wait with a WaitGroup so you KNOW they exited
before returning. (4) Cancel from the top with defer cancel(). Stopping = signal +
wait; never just "return and hope".
PITFALLS:
 1. Closing a channel from the receiver side or from several senders (panic).
 2. A stage that doesn't select on ctx.Done() when sending blocks forever if the
    consumer quit early (leak, see day 28).
 3. wg.Add inside the goroutine (race with Wait). Add before `go`.
 4. Capturing the loop variable: fine since Go 1.22 (each iteration has its own),
    but older code passed it as an argument.
 5. Results arrive in random order with fan-out; sort or index them if order matters.

EXERCISES:
 1. TODO: change the number of parsers (fan-out) from 3 to 1 and 8; the counts must not change.
 2. TODO: add a fourth stage that filters out "DEBUG" lines before counting.
 3. TODO: make runAll return the first error from a task that fails and see the others get cancelled.
*/
package main

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"
)

// rawLog is an in-memory log; real code would read a file or stdin.
const rawLog = `2025-01-10T10:00:01Z INFO  user=alice action=login
2025-01-10T10:00:02Z WARN  disk usage at 91%
2025-01-10T10:00:03Z ERROR payment failed order=1001
2025-01-10T10:00:04Z INFO  user=bob action=checkout
this line is garbage
2025-01-10T10:00:05Z ERROR db timeout
2025-01-10T10:00:06Z DEBUG cache hit sku=7
2025-01-10T10:00:07Z INFO  user=alice action=logout`

// Entry is a parsed log line.
type Entry struct {
	Level string
	Msg   string
}

// ---- Stage 1: generator. Emits lines, closes output when done or cancelled. ----
func readLines(ctx context.Context, text string) <-chan string {
	out := make(chan string)
	go func() {
		defer close(out) // the sender closes
		for _, line := range strings.Split(text, "\n") {
			select {
			case out <- line:
			case <-ctx.Done(): // consumer gone: stop instead of blocking forever
				return
			}
		}
	}()
	return out
}

// ---- Stage 2: parse. Run several copies (fan-out) on the same input. ----
func parse(ctx context.Context, in <-chan string) <-chan Entry {
	out := make(chan Entry)
	go func() {
		defer close(out)
		for line := range in {
			f := strings.Fields(line)
			// Malformed line: skip it. A valid line starts with a timestamp ending in Z.
			if len(f) < 3 || !strings.HasSuffix(f[0], "Z") {
				continue
			}
			e := Entry{Level: f[1], Msg: strings.Join(f[2:], " ")}
			select {
			case out <- e:
			case <-ctx.Done():
				return
			}
		}
	}()
	return out
}

// merge is FAN-IN: combine many channels into one.
func merge[T any](ins ...<-chan T) <-chan T {
	out := make(chan T)
	var wg sync.WaitGroup
	wg.Add(len(ins))
	for _, in := range ins {
		go func() { // one forwarder per input
			defer wg.Done()
			for v := range in {
				out <- v
			}
		}()
	}
	go func() { // closer: waits for ALL forwarders, then closes
		wg.Wait()
		close(out)
	}()
	return out
}

// ---- Stage 3: sink. Counts by level; runs in the caller's goroutine. ----
func countByLevel(in <-chan Entry) map[string]int {
	counts := map[string]int{}
	for e := range in { // only this goroutine touches the map: no mutex needed
		counts[e.Level]++
	}
	return counts
}

func demoPipeline() {
	fmt.Println("--- 3-stage pipeline: read -> parse (fan-out x3, fan-in) -> count ---")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	lines := readLines(ctx, rawLog)
	parsers := make([]<-chan Entry, 3)
	for i := range parsers {
		parsers[i] = parse(ctx, lines) // all three read the SAME channel (fan-out)
	}
	counts := countByLevel(merge(parsers...)) // fan-in then sink

	levels := make([]string, 0, len(counts))
	for l := range counts {
		levels = append(levels, l)
	}
	sort.Strings(levels) // map order is random; sort for stable output
	for _, l := range levels {
		fmt.Printf("  %-5s %d\n", l, counts[l])
	}
}

// ---- Worker pool: fixed number of workers, results collected ----
type job struct{ id, cents int }

func demoWorkerPool() {
	fmt.Println("--- worker pool: 4 workers, 10 invoices ---")
	jobs := make(chan job)
	results := make(chan int)
	var wg sync.WaitGroup
	for w := 0; w < 4; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := range jobs {
				results <- j.cents * 108 / 100 // add 8% tax
			}
		}()
	}
	go func() { // producer
		for i := 1; i <= 10; i++ {
			jobs <- job{i, i * 1000}
		}
		close(jobs)
	}()
	go func() { wg.Wait(); close(results) }() // close results once workers finish
	total := 0
	for r := range results {
		total += r
	}
	fmt.Println("  total with tax (cents):", total)
}

// ---- errgroup by hand ----
type group struct {
	wg     sync.WaitGroup
	once   sync.Once
	err    error
	cancel context.CancelFunc
}

// newGroup returns a group and a derived ctx that is cancelled on first error.
func newGroup(parent context.Context) (*group, context.Context) {
	ctx, cancel := context.WithCancel(parent)
	return &group{cancel: cancel}, ctx
}

// Go runs f in a goroutine; the first non-nil error is kept and cancels the rest.
func (g *group) Go(f func() error) {
	g.wg.Add(1)
	go func() {
		defer g.wg.Done()
		if err := f(); err != nil {
			g.once.Do(func() { // exactly once: only the FIRST error is stored
				g.err = err
				g.cancel()
			})
		}
	}()
}

// Wait blocks for all tasks, releases the context, and returns the first error.
func (g *group) Wait() error {
	g.wg.Wait()
	g.cancel()
	return g.err
}

func demoErrGroup() {
	fmt.Println("--- errgroup built by hand: first error cancels the rest ---")
	g, ctx := newGroup(context.Background())
	for _, svc := range []string{"inventory", "pricing", "shipping"} {
		g.Go(func() error {
			if svc == "pricing" {
				time.Sleep(20 * time.Millisecond)
				return errors.New("pricing service returned 500")
			}
			select {
			case <-time.After(time.Second): // would be slow
				return nil
			case <-ctx.Done():
				fmt.Printf("  %s: cancelled early (%v)\n", svc, ctx.Err())
				return ctx.Err()
			}
		})
	}
	fmt.Println("  Wait returned:", g.Wait())
}

func main() {
	demoPipeline()
	demoWorkerPool()
	demoErrGroup()
}
