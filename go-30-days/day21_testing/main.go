/*
DAY 21: TESTING
RUN:  go run ./day21_testing
TEST: go test -v ./day21_testing
BENCH: go test -bench=. ./day21_testing
RACE: go test -race ./day21_testing

THEORY:
Go has testing built in. A file ending in _test.go is compiled only by
"go test". A test is func TestXxx(t *testing.T); call t.Errorf to report a
failure and keep going, t.Fatalf to stop that test. Table-driven tests put
the cases in a slice of structs (name, input, expected) and loop over them,
so adding a case is one line. t.Run(name, func) makes a subtest that shows up
as TestX/name, can be run alone with -run, and reports separately.
Benchmarks are func BenchmarkXxx(b *testing.B) looping b.N times; the runner
picks N until timing is stable. Example functions (func ExampleXxx) have an
"// Output:" comment that go test compares with real stdout, so docs never lie.
-race builds with the race detector which reports unsynchronised memory access.
REAL WORLD: every Go project runs "go test ./..." in CI; table-driven tests are
the standard style in the stdlib; benchmarks guard hot paths.
EXPLAIN-IT ANSWER: How does a table-driven test make failures easy to read?
Each case has a name and its inputs and expected output side by side, and the
failure message prints name, input, got and want. You see exactly WHICH case
failed and why without reading test logic, all cases run in one pass (no
stopping at the first), and with t.Run you can rerun a single failing case.
PITFALLS:
  - Using t.Fatalf inside a table loop without subtests stops all remaining cases.
  - Failure messages that omit input or the got/want values.
  - Testing implementation details instead of behaviour.
  - Benchmark results optimised away: store the result in a package variable.

EXERCISES:
 1. Add a case to TestSlugify for "  Many   spaces  " and make it pass.
 2. Add an ApplyDiscount case where percent is exactly 100 and one that is -5.
 3. Write BenchmarkApplyDiscount and compare it with BenchmarkSlugify.
*/
package main

import (
	"errors"
	"fmt"
	"strings"
	"unicode"
)

// ErrBadPercent is returned when a discount percentage is outside 0-100.
var ErrBadPercent = errors.New("percent must be between 0 and 100")

// Slugify turns a title into a URL-safe slug: "Hello, World!" -> "hello-world".
func Slugify(title string) string {
	var b strings.Builder
	pendingDash := false // true when we saw a separator and still owe a dash
	for _, r := range strings.ToLower(title) {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			if pendingDash && b.Len() > 0 { // no leading dash
				b.WriteByte('-')
			}
			pendingDash = false
			b.WriteRune(r)
		} else {
			pendingDash = true // collapse runs of punctuation/spaces into one dash
		}
	}
	return b.String()
}

// ApplyDiscount reduces a price in cents by percent and rounds to the nearest cent.
// Money is kept in integer cents to avoid floating point errors in totals.
func ApplyDiscount(cents int, percent int) (int, error) {
	if percent < 0 || percent > 100 {
		return 0, fmt.Errorf("apply discount %d%%: %w", percent, ErrBadPercent)
	}
	discount := (cents*percent + 50) / 100 // +50 rounds half up using integer maths
	return cents - discount, nil
}

func main() {
	fmt.Println("== Slugify")
	for _, title := range []string{"Hello, World!", "  Go: 30 Days -- Day 21  ", "Crème Brûlée"} {
		fmt.Printf("%q -> %q\n", title, Slugify(title))
	}
	fmt.Println("== ApplyDiscount")
	for _, p := range []int{0, 15, 100, 120} {
		out, err := ApplyDiscount(1999, p)
		fmt.Printf("19.99 with %3d%% off -> %d cents, err=%v\n", p, out, err)
	}
	fmt.Println("Now run: go test -v ./day21_testing")
}
