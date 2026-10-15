package main

import (
	"errors"
	"fmt"
	"testing"
)

// TestSlugify is table-driven: one slice of cases, one loop, one subtest each.
func TestSlugify(t *testing.T) {
	tests := []struct {
		name string // shows up in output as TestSlugify/<name>
		in   string
		want string
	}{
		{"simple", "Hello World", "hello-world"},
		{"punctuation", "Hello, World!", "hello-world"},
		{"extra separators", "  Go: 30 Days -- Day 21  ", "go-30-days-day-21"},
		{"accents kept", "Crème Brûlée", "crème-brûlée"},
		{"empty", "", ""},
		{"only symbols", "!!!", ""},
	}
	for _, tc := range tests {
		// Since Go 1.22 tc is a fresh variable per iteration, so it is safe to capture.
		t.Run(tc.name, func(t *testing.T) {
			got := Slugify(tc.in)
			if got != tc.want { // the message has input, got and want: all you need to debug
				t.Errorf("Slugify(%q) = %q, want %q", tc.in, got, tc.want)
			}
		})
	}
}

// TestApplyDiscount checks both values and the error path with errors.Is.
func TestApplyDiscount(t *testing.T) {
	tests := []struct {
		name    string
		cents   int
		percent int
		want    int
		wantErr error
	}{
		{"no discount", 1999, 0, 1999, nil},
		{"fifteen percent rounds", 1999, 15, 1699, nil},
		{"full discount", 1999, 100, 0, nil},
		{"negative percent", 1000, -1, 0, ErrBadPercent},
		{"over 100", 1000, 101, 0, ErrBadPercent},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := ApplyDiscount(tc.cents, tc.percent)
			if !errors.Is(err, tc.wantErr) { // errors.Is(nil, nil) is true
				t.Fatalf("ApplyDiscount(%d, %d) error = %v, want %v", tc.cents, tc.percent, err, tc.wantErr)
			}
			if got != tc.want {
				t.Errorf("ApplyDiscount(%d, %d) = %d, want %d", tc.cents, tc.percent, got, tc.want)
			}
		})
	}
}

// sink stops the compiler from deleting the benchmarked call as unused.
var sink string

// BenchmarkSlugify measures one call; the runner chooses b.N.
func BenchmarkSlugify(b *testing.B) {
	for i := 0; i < b.N; i++ {
		sink = Slugify("  Go: 30 Days -- Day 21 Testing  ")
	}
}

// ExampleSlugify is documentation AND a test: go test checks the Output comment.
func ExampleSlugify() {
	fmt.Println(Slugify("Hello, World!"))
	// Output: hello-world
}
