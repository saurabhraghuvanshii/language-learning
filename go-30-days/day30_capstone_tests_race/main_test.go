package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"
)

// newTestHandler builds a fresh API with a fresh store and a silent logger.
func newTestHandler() (http.Handler, *MemStore) {
	store := NewMemStore()
	api := &API{store: store}
	return api.Routes(log.New(io.Discard, "", 0)), store
}

// TestHandlersTable is table-driven: each case gets its OWN fresh store, so cases
// cannot depend on each other. seed items are created first.
func TestHandlersTable(t *testing.T) {
	tests := []struct {
		name       string
		seed       []Item
		method     string
		path       string
		body       string
		wantStatus int
		wantInBody string // substring expected in the response body
	}{
		{"list empty", nil, "GET", "/items", "", 200, "[]"},
		{"list one", []Item{{Name: "Cable", Qty: 1, PriceCents: 100}}, "GET", "/items", "", 200, "Cable"},
		{"get existing", []Item{{Name: "Cable"}}, "GET", "/items/1", "", 200, `"id":1`},
		{"get missing", nil, "GET", "/items/42", "", 404, "not found"},
		{"get bad id", nil, "GET", "/items/abc", "", 400, "positive integer"},
		{"create ok", nil, "POST", "/items", `{"name":"Stand","qty":3,"price_cents":999}`, 201, `"name":"Stand"`},
		{"create bad json", nil, "POST", "/items", `{nope`, 400, "invalid JSON"},
		{"create unknown field", nil, "POST", "/items", `{"name":"x","colour":"red"}`, 400, "invalid JSON"},
		{"create invalid", nil, "POST", "/items", `{"name":" "}`, 422, "name required"},
		{"delete ok", []Item{{Name: "Cable"}}, "DELETE", "/items/1", "", 204, ""},
		{"delete missing", nil, "DELETE", "/items/7", "", 404, "not found"},
		{"wrong method", nil, "PUT", "/items/1", "", 405, ""},
		{"unknown path", nil, "GET", "/nope", "", 404, ""},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			h, store := newTestHandler()
			for _, it := range tc.seed {
				store.Create(it)
			}
			req := httptest.NewRequest(tc.method, tc.path, strings.NewReader(tc.body))
			rec := httptest.NewRecorder() // fake ResponseWriter that records everything
			h.ServeHTTP(rec, req)         // no network involved
			if rec.Code != tc.wantStatus {
				t.Fatalf("status = %d, want %d (body %q)", rec.Code, tc.wantStatus, rec.Body.String())
			}
			if !strings.Contains(rec.Body.String(), tc.wantInBody) {
				t.Errorf("body %q does not contain %q", rec.Body.String(), tc.wantInBody)
			}
		})
	}
}

// TestEndToEndServer uses a real server on a random port and a real client.
func TestEndToEndServer(t *testing.T) {
	h, _ := newTestHandler()
	srv := httptest.NewServer(h)
	defer srv.Close()
	client := &http.Client{Timeout: 2 * time.Second}

	resp, err := client.Post(srv.URL+"/items", "application/json", strings.NewReader(`{"name":"Monitor","qty":2,"price_cents":15000}`))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("POST status = %d, want 201", resp.StatusCode)
	}
	if loc := resp.Header.Get("Location"); loc != "/items/1" {
		t.Errorf("Location = %q, want /items/1", loc)
	}
	var created Item
	if err := json.NewDecoder(resp.Body).Decode(&created); err != nil {
		t.Fatal(err)
	}
	if created.ID != 1 || created.Name != "Monitor" {
		t.Errorf("created = %+v", created)
	}

	get, err := client.Get(srv.URL + "/items/1")
	if err != nil {
		t.Fatal(err)
	}
	defer get.Body.Close()
	if get.StatusCode != 200 {
		t.Errorf("GET status = %d, want 200", get.StatusCode)
	}
	if ct := get.Header.Get("Content-Type"); ct != "application/json" {
		t.Errorf("Content-Type = %q", ct)
	}
}

// TestStoreConcurrent hammers the store from 50 goroutines. Under -race this FAILS
// if the mutex is missing; without -race it still checks the final counts.
func TestStoreConcurrent(t *testing.T) {
	store := NewMemStore()
	const goroutines = 50
	const perG = 20
	var wg sync.WaitGroup
	for g := 0; g < goroutines; g++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < perG; i++ {
				it := store.Create(Item{Name: "widget", Qty: i})
				store.Get(it.ID)
				store.List()
				if i%2 == 0 {
					store.Delete(it.ID)
				}
			}
		}()
	}
	wg.Wait()
	want := goroutines * perG / 2 // half were deleted
	if got := len(store.List()); got != want {
		t.Fatalf("items left = %d, want %d", got, want)
	}
	// IDs must be unique: next ID is exactly one past all creates.
	next := store.Create(Item{Name: "last"})
	if next.ID != goroutines*perG+1 {
		t.Errorf("next ID = %d, want %d", next.ID, goroutines*perG+1)
	}
}

// BenchmarkGetItem measures one GET /items/1 through the whole handler chain.
func BenchmarkGetItem(b *testing.B) {
	h, store := newTestHandler()
	store.Create(Item{Name: "Cable", Qty: 1, PriceCents: 100})
	b.ReportAllocs()
	b.ResetTimer() // exclude setup from timing
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest("GET", "/items/1", nil)
		h.ServeHTTP(httptest.NewRecorder(), req)
	}
}

// BenchmarkStoreParallel measures the store under contention.
func BenchmarkStoreParallel(b *testing.B) {
	store := NewMemStore()
	store.Create(Item{Name: "Cable"})
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			store.Get(1)
		}
	})
}
