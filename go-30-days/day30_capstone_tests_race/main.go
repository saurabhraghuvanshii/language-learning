/*
DAY 30: CAPSTONE - TESTING THE API WITH -race
RUN:  go run ./day30_capstone_tests_race          (same self-test demo as day 29)
TEST: go test -race ./day30_capstone_tests_race
BENCH: go test -bench . -benchmem ./day30_capstone_tests_race

THEORY (plain words):
This folder is a self-contained COPY of the day 29 inventory API (Store, MemStore,
handlers, logging middleware, graceful shutdown) plus main_test.go. Tests prove the
code is right; the race detector (-race) proves it is safe under concurrency.
  - httptest.NewRecorder: a fake ResponseWriter. You call handler.ServeHTTP(rec, req)
    directly: no network, no ports, microsecond fast. Best for unit-testing handlers.
  - httptest.NewServer: a real server on a random localhost port; you use a real
    http.Client. Slower but exercises routing, middleware, keep-alive, JSON on the wire.
  - Table-driven tests: a slice of cases {name, request, wantStatus, wantBody}, looped
    with t.Run. Adding a case is one line.
  - go test -race builds with the race detector: it instruments memory accesses and
    reports two goroutines touching the same memory without synchronisation.
  - Benchmarks (func BenchmarkX(b *testing.B)) run the body b.N times and report ns/op.

REAL WORLD: CI runs `go test -race ./...` on every pull request; handlers are tested
with httptest; stores with concurrent tests; hot paths with benchmarks.
EXPLAIN-IT ANSWER: teach the project in 10 minutes (TEACHING SCRIPT)
 1. (1 min) Goal: a tiny inventory REST API: GET/POST/DELETE on /items, JSON in/out.
 2. (1 min) Data: Item struct with json tags; Store interface = what handlers need.
 3. (1 min) MemStore: a map plus sync.RWMutex. Why? Every request is its own goroutine.
 4. (1 min) Handlers: decode -> validate -> store -> status code -> JSON. 201/204/404/422.
 5. (1 min) ServeMux patterns "GET /items/{id}" and r.PathValue("id").
 6. (1 min) Middleware: func(Handler) Handler; logging captures the status via a wrapper.
 7. (1 min) Server lifecycle: Listen, Serve in a goroutine, Shutdown(ctx) to drain.
 8. (1 min) Request path: TCP accept -> goroutine per conn -> mux -> middleware -> handler.
 9. (1 min) Tests: httptest.NewRecorder (table-driven) and NewServer (end to end).
 10. (1 min) Concurrency proof: 50 goroutines hammer the store under -race; then
    remove the mutex and watch -race fail. That is the whole point of day 25.

PITFALLS: tests that share one store across table cases (order dependent), t.Fatal

	called from a non-test goroutine (use t.Error there), forgetting resp.Body.Close.

EXERCISES:
 1. TODO: remove the RLock/Lock calls in MemStore and run go test -race: read the report.
 2. TODO: add a table case for POST with an unknown JSON field (expect 400).
 3. TODO: add a benchmark for GET /items with 1000 items stored and compare to 10.
*/
package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sort"
	"strings"
	"sync"
	"syscall"
	"time"
)

// ---------------- model + store ----------------

// Item is one inventory line. The `json:"..."` tags name the JSON fields.
type Item struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Qty        int    `json:"qty"`
	PriceCents int    `json:"price_cents"`
}

// Store is what handlers need. Depending on an interface lets us swap
// implementations (memory, SQL) and fake it in tests.
type Store interface {
	List() []Item
	Get(id int) (Item, bool)
	Create(it Item) Item // assigns the ID, returns the stored item
	Delete(id int) bool  // reports whether something was deleted
}

// MemStore is a Store backed by a map. It is shared by every request goroutine.
type MemStore struct {
	mu     sync.RWMutex // guards items and nextID
	items  map[int]Item
	nextID int
}

// NewMemStore returns a ready-to-use store (map initialised, IDs start at 1).
func NewMemStore() *MemStore {
	return &MemStore{items: make(map[int]Item), nextID: 1}
}

// List returns all items sorted by ID (map order is random, APIs should be stable).
func (s *MemStore) List() []Item {
	s.mu.RLock() // many GET requests may list at once
	defer s.mu.RUnlock()
	out := make([]Item, 0, len(s.items)) // never nil, so JSON encodes [] not null
	for _, it := range s.items {
		out = append(out, it)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out
}

// Get returns a copy of the item.
func (s *MemStore) Get(id int) (Item, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	it, ok := s.items[id]
	return it, ok
}

// Create needs the write lock: it mutates the map AND nextID together.
func (s *MemStore) Create(it Item) Item {
	s.mu.Lock()
	defer s.mu.Unlock()
	it.ID = s.nextID
	s.nextID++
	s.items[it.ID] = it
	return it
}

// Delete removes an item; reports whether it existed.
func (s *MemStore) Delete(id int) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, ok := s.items[id]
	delete(s.items, id)
	return ok
}

// ---------------- handlers ----------------

// API groups handlers and their dependency (the Store).
type API struct {
	store Store
}

// writeJSON sends a JSON body with the given status.
func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json") // headers BEFORE WriteHeader/Write
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v) // error here means the client went away; nothing to do
}

// writeError sends {"error": "..."} so clients get machine-readable failures.
func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

// listItems handles GET /items.
func (a *API) listItems(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, a.store.List())
}

// getItem handles GET /items/{id}.
func (a *API) getItem(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r.PathValue("id")) // value of {id} from the matched pattern
	if err != nil {
		writeError(w, http.StatusBadRequest, "id must be a positive integer")
		return
	}
	it, ok := a.store.Get(id)
	if !ok {
		writeError(w, http.StatusNotFound, "item not found")
		return
	}
	writeJSON(w, http.StatusOK, it)
}

// createItem handles POST /items.
func (a *API) createItem(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20) // cap body at 1 MiB
	var in Item
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields() // typos in field names become 400s, not silent data loss
	if err := dec.Decode(&in); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON: "+err.Error())
		return
	}
	if strings.TrimSpace(in.Name) == "" || in.Qty < 0 || in.PriceCents < 0 {
		writeError(w, http.StatusUnprocessableEntity, "name required; qty and price_cents must be >= 0")
		return
	}
	created := a.store.Create(in)
	w.Header().Set("Location", fmt.Sprintf("/items/%d", created.ID)) // where the new thing lives
	writeJSON(w, http.StatusCreated, created)                        // 201, not 200
}

// deleteItem handles DELETE /items/{id}.
func (a *API) deleteItem(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "id must be a positive integer")
		return
	}
	if !a.store.Delete(id) {
		writeError(w, http.StatusNotFound, "item not found")
		return
	}
	w.WriteHeader(http.StatusNoContent) // 204: success, no body
}

// parseID converts a path segment to a positive int.
func parseID(s string) (int, error) {
	var id int
	if _, err := fmt.Sscanf(s, "%d", &id); err != nil || id <= 0 || fmt.Sprint(id) != s {
		return 0, errors.New("bad id")
	}
	return id, nil
}

// Routes wires method+path patterns to handlers and wraps them in middleware.
func (a *API) Routes(logger *log.Logger) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /items", a.listItems)
	mux.HandleFunc("GET /items/{id}", a.getItem)
	mux.HandleFunc("POST /items", a.createItem)
	mux.HandleFunc("DELETE /items/{id}", a.deleteItem)
	return logging(logger, mux) // outermost wrapper runs first
}

// ---------------- middleware ----------------

// statusRecorder remembers the status code a handler wrote, because
// http.ResponseWriter has no getter for it.
type statusRecorder struct {
	http.ResponseWriter // embedded: all other methods pass through
	status              int
}

// WriteHeader records the code then forwards it.
func (s *statusRecorder) WriteHeader(code int) {
	s.status = code
	s.ResponseWriter.WriteHeader(code)
}

// logging logs method, path, status and duration of every request.
func logging(logger *log.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK} // 200 is the implicit default
		next.ServeHTTP(rec, r)                                           // run the rest of the chain
		logger.Printf("server: %s %s -> %d (%v)", r.Method, r.URL.Path, rec.status, time.Since(start).Round(time.Microsecond))
	})
}

// ---------------- server wiring ----------------

// newServer builds an http.Server with safe timeouts.
func newServer(h http.Handler, addr string) *http.Server {
	return &http.Server{
		Addr:              addr,
		Handler:           h,
		ReadHeaderTimeout: 5 * time.Second, // slowloris protection
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}
}

// run serves on ln until ctx is cancelled, then shuts down gracefully:
// stop accepting, let in-flight requests finish (up to 5s), then return.
func run(ctx context.Context, srv *http.Server, ln net.Listener) error {
	errCh := make(chan error, 1)
	go func() {
		// Serve blocks; returns http.ErrServerClosed after Shutdown is called.
		errCh <- srv.Serve(ln)
	}()
	select {
	case err := <-errCh: // server failed on its own
		return err
	case <-ctx.Done(): // told to stop
	}
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		return err
	}
	if err := <-errCh; !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

// ---------------- the self-test client ----------------

// call performs one request and prints both sides.
func call(client *http.Client, method, url, body string) {
	var rdr io.Reader
	if body != "" {
		rdr = strings.NewReader(body)
	}
	req, err := http.NewRequest(method, url, rdr)
	if err != nil {
		fmt.Println("client: build request:", err)
		return
	}
	if body != "" {
		req.Header.Set("Content-Type", "application/json")
	}
	fmt.Printf("client: >>> %s %s %s\n", method, req.URL.Path, body)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("client: error:", err)
		return
	}
	defer resp.Body.Close() // return the connection to the pool
	b, _ := io.ReadAll(resp.Body)
	fmt.Printf("client: <<< %d %s\n", resp.StatusCode, strings.TrimSpace(string(b)))
}

func main() {
	serve := flag.Bool("serve", false, "keep serving on :8080 instead of running the self-test")
	flag.Parse()

	logger := log.New(os.Stdout, "", 0) // no timestamps so output is stable
	api := &API{store: NewMemStore()}
	handler := api.Routes(logger)

	if *serve {
		// NotifyContext cancels ctx on Ctrl-C or SIGTERM (what Kubernetes sends).
		ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
		defer stop()
		ln, err := net.Listen("tcp", ":8080")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("serving on http://localhost:8080 (Ctrl-C to stop)")
		if err := run(ctx, newServer(handler, ":8080"), ln); err != nil {
			log.Fatal(err)
		}
		fmt.Println("shut down cleanly")
		return
	}

	// Self-test mode: port 0 asks the OS for any free port.
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		log.Fatal(err)
	}
	base := "http://" + ln.Addr().String()
	fmt.Println("self-test server on", base)

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() { done <- run(ctx, newServer(handler, ln.Addr().String()), ln) }()

	client := &http.Client{Timeout: 2 * time.Second} // never use a client with no timeout
	call(client, "POST", base+"/items", `{"name":"USB-C cable","qty":120,"price_cents":899}`)
	call(client, "POST", base+"/items", `{"name":"Laptop stand","qty":15,"price_cents":3499}`)
	call(client, "GET", base+"/items", "")
	call(client, "GET", base+"/items/1", "")
	call(client, "GET", base+"/items/99", "") // expect 404
	call(client, "POST", base+"/items", `{"name":""}`)
	call(client, "DELETE", base+"/items/1", "")
	call(client, "GET", base+"/items", "")

	cancel() // trigger graceful shutdown
	if err := <-done; err != nil {
		log.Fatal(err)
	}
	fmt.Println("server shut down gracefully")
}
