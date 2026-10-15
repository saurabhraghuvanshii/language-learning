/*
DAY 25: MUTEX, DATA RACES, ATOMIC, RWMUTEX AND CHANNEL ALTERNATIVES
RUN:  go run ./day25_mutex_data_races
RACE: go run -race ./day25_mutex_data_races   (the report is EXPECTED for section 1 only)

THEORY (plain words):
A DATA RACE is two goroutines touching the same memory at the same time where at
least one writes, with no synchronisation. The result is undefined: lost updates,
torn values, rare crashes. Go gives you four tools, all shown below:
 1. sync.Mutex      only one goroutine at a time inside Lock()...Unlock().
 2. sync/atomic     the CPU does a single indivisible read-modify-write.
 3. sync.RWMutex    many readers OR one writer.
 4. channels        one goroutine OWNS the data; others send it requests.

Under the hood a Mutex is an integer flag plus a wait queue. Lock() tries to flip
the flag with an atomic compare-and-swap; if it is taken, the goroutine is parked by
the scheduler (not spinning forever) and woken by Unlock(). Unlock also creates a
"happens-before" edge: everything done before Unlock is visible to the next Locker.

WHEN TO USE WHAT (decision table)

	Situation                                        Use
	-----------------------------------------------  ---------------------------------
	One number: counter, flag, last-seen timestamp   sync/atomic (atomic.Int64, Bool)
	Several fields that must change together (a      sync.Mutex (invariants span >1
	  balance AND a transaction count), or a map      variable; maps cannot be atomic)
	Read-heavy shared data, writes rare (cache of    sync.RWMutex (readers run in
	  product prices, config, routing tables)         parallel). Measure first; for
	                                                  short critical sections a plain
	                                                  Mutex is often as fast.
	Passing OWNERSHIP of data / a pipeline of work,  channel (send the value, stop
	  a worker pool, events, shutdown signals         sharing it)
	Write-once-then-read-only data                   no lock at all (build before
	                                                  starting goroutines)
	Rule of thumb: "Share memory by communicating" for flows of work; use a mutex to
	protect a small piece of state ("a cache, a counter, a map"). Keep the locked
	section SHORT: no network calls, no sleeps, no callbacks while holding a lock.

WHEN NOT TO USE CONCURRENCY
  - The work is tiny (adding 1000 numbers): starting a goroutine (~2KB stack,
    scheduling, channel handoff) costs more than the work itself.
  - The steps are naturally sequential (each needs the previous result).
  - You'd be sharing mutable state everywhere: if every goroutine needs the same map,
    you have built a single-lane road with extra steps. Locks serialise the work.
  - The program is I/O-trivial and single-user (a small CLI). Simpler is faster to
    write, read and debug.
  - You can't answer "who stops this goroutine and when?" (see day 28: leaks).
    USE it when: independent I/O waits (many HTTP calls), CPU work that splits cleanly
    into independent chunks, a server handling many clients (net/http already does).

REAL WORLD: request counters and metrics, in-memory caches, rate limiters, session
stores, connection registries. Every net/http handler runs concurrently (day 29).
EXPLAIN-IT ANSWER: Why is counter++ not safe across goroutines?
counter++ looks like one step but compiles to THREE: (1) READ counter into a
register, (2) ADD 1 in the register, (3) WRITE the register back. Two goroutines can
both READ 41, both compute 42, both WRITE 42: two increments, one result. That is a
"lost update". Mutex makes the three steps one critical section; atomic.Add makes
them one hardware instruction.
PITFALLS:
 1. Lost updates from unsynchronised ++ (section 1).
 2. COPYING A MUTEX. A Mutex holds state; copying it copies that state and gives
    two unrelated locks. Example that `go vet` reports ("passes lock by value"):
    func (c Counter) Inc() { c.mu.Lock(); ... }   // value receiver copies c.mu!
    b := a                                         // copies a.mu
    Always use pointer receivers on types containing a Mutex and never copy them.
    (Not done in the code below so vet stays clean.)
 3. Forgetting Unlock (use defer mu.Unlock() right after Lock()).
 4. Locking twice in the same goroutine (Go mutexes are NOT reentrant: deadlock).
    A method that holds the lock must not call another method that also locks.
 5. Holding a RLock then trying to Lock (upgrade) deadlocks. Release first.
 6. Reading a map while another goroutine writes it is a fatal error, not just a
    wrong answer: "fatal error: concurrent map read and map write".

EXERCISES:
 1. TODO: run `go run -race ./day25_mutex_data_races` and read the report: find the
    two stack traces (read at line X, previous write at line Y).
 2. TODO: add a method Reset() to SafeCounter and call it from main; make sure it locks.
 3. TODO: add a PriceCache.Delete(sku) with the right kind of lock (RLock or Lock?).
*/
package main

import (
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
)

const workers = 100 // goroutines
const perWorker = 1000

// ---------- 1. BROKEN: data race ----------

// brokenCounter has no protection.
type brokenCounter struct{ n int }

// Inc is racy: n++ is read, add, write.
// We spell the three steps out and call runtime.Gosched() between the read and
// the write. Gosched politely lets another goroutine run. Real code gets
// interrupted at random moments too, but a tiny loop can finish before the
// next goroutine even starts, so on a fast machine the bug may stay hidden.
// Gosched makes the interleaving happen EVERY time, so you can see the bug.
func (c *brokenCounter) Inc() {
	tmp := c.n        // step 1: READ the current value
	runtime.Gosched() // another goroutine may now read the SAME old value
	c.n = tmp + 1     // step 2+3: ADD one and WRITE back, overwriting their work
}

func runBroken() {
	fmt.Println("--- 1. broken counter (data race, lost updates) ---")
	var c brokenCounter
	var wg sync.WaitGroup // waits for all goroutines to finish
	for i := 0; i < workers; i++ {
		wg.Add(1) // add BEFORE starting the goroutine
		go func() {
			defer wg.Done()
			for j := 0; j < perWorker; j++ {
				c.Inc() // RACE: concurrent unsynchronised read+write of c.n
			}
		}()
	}
	wg.Wait()
	fmt.Printf("expected %d, got %d (usually less: lost updates)\n", workers*perWorker, c.n)
}

// ---------- 2. FIXED with Mutex ----------

// SafeCounter bundles the data with the mutex that guards it. Put the mutex
// right above the fields it protects and say so in a comment.
type SafeCounter struct {
	mu sync.Mutex // guards n
	n  int
}

// Inc locks around the whole read-modify-write.
func (c *SafeCounter) Inc() {
	c.mu.Lock()
	defer c.mu.Unlock() // runs even if the code below panics
	c.n++
}

// Value must lock too: reads race with writes as well.
func (c *SafeCounter) Value() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.n
}

func runMutex() {
	fmt.Println("--- 2. fixed with sync.Mutex ---")
	var c SafeCounter
	var wg sync.WaitGroup
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < perWorker; j++ {
				c.Inc()
			}
		}()
	}
	wg.Wait()
	fmt.Printf("expected %d, got %d\n", workers*perWorker, c.Value())
}

// ---------- 3. FIXED with atomic ----------

func runAtomic() {
	fmt.Println("--- 3. fixed with sync/atomic ---")
	var n atomic.Int64 // typed atomic (Go 1.19+); safer than atomic.AddInt64(&x, 1)
	var wg sync.WaitGroup
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < perWorker; j++ {
				n.Add(1) // one indivisible CPU instruction
			}
		}()
	}
	wg.Wait()
	fmt.Printf("expected %d, got %d\n", workers*perWorker, n.Load())
}

// ---------- 4. RWMutex: read-heavy cache of product prices ----------

// PriceCache maps SKU -> price in cents.
type PriceCache struct {
	mu     sync.RWMutex // guards prices
	prices map[string]int
}

// NewPriceCache builds the map up front: a nil map would panic on write.
func NewPriceCache() *PriceCache {
	return &PriceCache{prices: make(map[string]int)}
}

// Get takes a READ lock: any number of goroutines can be in here together.
func (p *PriceCache) Get(sku string) (int, bool) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	v, ok := p.prices[sku]
	return v, ok
}

// Set takes the WRITE lock: waits for readers to leave, blocks new ones.
func (p *PriceCache) Set(sku string, cents int) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.prices[sku] = cents
}

func runRWMutex() {
	fmt.Println("--- 4. sync.RWMutex price cache (many readers, rare writer) ---")
	cache := NewPriceCache()
	cache.Set("SKU-1", 1999)
	cache.Set("SKU-2", 4999)
	var wg sync.WaitGroup
	var hits atomic.Int64
	for i := 0; i < 50; i++ { // 50 readers (like 50 concurrent HTTP requests)
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 200; j++ {
				if _, ok := cache.Get("SKU-1"); ok {
					hits.Add(1)
				}
			}
		}()
	}
	wg.Add(1)
	go func() { // one writer: a price update arrives
		defer wg.Done()
		for j := 0; j < 20; j++ {
			cache.Set("SKU-2", 4999+j)
		}
	}()
	wg.Wait()
	p, _ := cache.Get("SKU-2")
	fmt.Printf("reader hits=%d (expected 10000), final SKU-2 price=%d\n", hits.Load(), p)
}

// ---------- 5. Channel alternative: one goroutine owns the data ----------

// runChannel keeps the counter inside ONE goroutine. Others send "+1" messages.
// No lock needed because only the owner ever touches n.
func runChannel() {
	fmt.Println("--- 5. channel alternative (owner goroutine) ---")
	inc := make(chan struct{}, 64) // increment requests
	result := make(chan int)       // owner reports final value here
	go func() {                    // the OWNER
		n := 0
		for range inc { // runs until inc is closed
			n++
		}
		result <- n
	}()
	var wg sync.WaitGroup
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < perWorker; j++ {
				inc <- struct{}{}
			}
		}()
	}
	wg.Wait()
	close(inc) // all senders finished, so closing is safe
	fmt.Printf("expected %d, got %d\n", workers*perWorker, <-result)
	// Slower than atomic/mutex for a plain counter. Channels shine when you pass
	// ownership of real work, not when you protect one integer.
}

func main() {
	runBroken()
	runMutex()
	runAtomic()
	runRWMutex()
	runChannel()
	fmt.Println("rule: one number -> atomic; several fields/map -> mutex; read-heavy -> RWMutex; flow of work -> channel")
}
