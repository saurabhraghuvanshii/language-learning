/*
DAY 11: MAPS
RUN:  go run ./day11_maps

THEORY:
A map is a hash table: map[K]V stores values by key with average O(1) lookup. Keys must
be comparable (==): ints, strings, structs of comparable fields, arrays; not slices,
maps or funcs. A map value is a small pointer to a runtime structure (hmap) that holds
buckets; make(map[K]V) allocates it. A nil map has no structure behind it: reading
returns the zero value, but a write has nowhere to put the entry. Looking up a missing
key returns the zero value of V, so use the comma-ok form v, ok := m[k] to tell "absent"
from "present with zero value". delete(m, k) removes a key (no error if absent).
Iteration order is deliberately randomised by the runtime. Maps are not safe for
concurrent writes (use a mutex or sync.Map).
REAL WORLD: counting events from logs, caching users by ID, grouping orders by
customer, de-duplicating, HTTP header sets (map[string][]string).
EXPLAIN-IT ANSWER: Why does writing to a nil map panic and why is iteration order random?
A nil map is just a nil pointer with no hash table allocated, so an insert would need
to create storage inside a map you only hold by value, which Go refuses to do implicitly:
it panics "assignment to entry in nil map". Reads are safe because they just find nothing.
Iteration order is random because the runtime picks a random start bucket and offset on
every range, on purpose, so programs never depend on an order that the implementation
can change when the table grows. Sort the keys if you need stable output.
PITFALLS:
 1. Writing to a nil map panics.
 2. m[k] on a missing key gives the zero value; use comma-ok to distinguish.
 3. Relying on range order; it differs between runs.

EXERCISES:
 1. TODO: add a "longest word" feature using a second map[string]int of lengths.
 2. TODO: change topN(counts, 3) to topN(counts, 5) and print the share of each word as a percentage.
 3. TODO: write func invert(m map[string]int) map[int][]string grouping words by count.
*/
package main

import (
	"fmt"
	"sort"
	"strings"
)

// paragraph is real text (the opening of Moby-Dick, public domain).
const paragraph = `Call me Ishmael. Some years ago, never mind how long precisely, having
little or no money in my purse, and nothing particular to interest me on shore, I
thought I would sail about a little and see the watery part of the world. It is a way
I have of driving off the spleen and regulating the circulation. Whenever I find myself
growing grim about the mouth; whenever it is a damp, drizzly November in my soul;
whenever I find myself involuntarily pausing before coffin warehouses, and bringing up
the rear of every funeral I meet; then, I account it high time to get to sea as soon
as I can.`

// wordCount lowercases the text, splits on anything that is not a letter, and counts.
func wordCount(text string) map[string]int {
	counts := make(map[string]int) // make allocates the hash table; without it we would have a nil map
	// FieldsFunc splits where the function returns true, i.e. on every non-letter rune.
	words := strings.FieldsFunc(strings.ToLower(text), func(r rune) bool {
		return r < 'a' || r > 'z'
	})
	for _, w := range words {
		counts[w]++ // missing key reads as 0, so ++ works as "start at 0 then add 1"
	}
	return counts
}

// topN returns the n most frequent words, ties broken alphabetically for stable output.
func topN(counts map[string]int, n int) []string {
	keys := make([]string, 0, len(counts)) // preallocate: we know the final size
	for k := range counts {
		keys = append(keys, k) // collect keys: map order is random, so we sort next
	}
	sort.Slice(keys, func(i, j int) bool {
		if counts[keys[i]] != counts[keys[j]] {
			return counts[keys[i]] > counts[keys[j]] // higher count first
		}
		return keys[i] < keys[j] // tie: alphabetical
	})
	if n > len(keys) {
		n = len(keys) // guard against slicing past the end
	}
	return keys[:n]
}

func main() {
	fmt.Println("== word frequency")
	counts := wordCount(paragraph)
	fmt.Println("distinct words:", len(counts))
	for i, w := range topN(counts, 3) {
		fmt.Printf("  #%d %-5s %d\n", i+1, w, counts[w])
	}

	fmt.Println("\n== lookup with comma-ok (PITFALL 2)")
	fmt.Println("counts[\"whale\"] =", counts["whale"], "(0 means absent OR zero: ambiguous)")
	if n, ok := counts["sea"]; ok { // ok is true only when the key exists
		fmt.Println("sea is present, count", n)
	}
	if _, ok := counts["whale"]; !ok {
		fmt.Println("whale is really absent")
	}

	fmt.Println("\n== delete")
	delete(counts, "sea")
	delete(counts, "no-such-key") // deleting a missing key is a no-op, not an error
	_, stillThere := counts["sea"]
	fmt.Println("sea present after delete?", stillThere)

	fmt.Println("\n== PITFALL 3: random iteration order (run twice to see it change)")
	stock := map[string]int{"apple": 3, "pear": 5, "plum": 7, "fig": 1, "kiwi": 9}
	fmt.Print("  range order :")
	for k := range stock {
		fmt.Print(" ", k) // may differ on every run
	}
	fmt.Println()
	names := make([]string, 0, len(stock))
	for k := range stock {
		names = append(names, k)
	}
	sort.Strings(names) // stable, predictable order
	fmt.Print("  sorted order:")
	for _, k := range names {
		fmt.Print(" ", k)
	}
	fmt.Println()

	fmt.Println("\n== PITFALL 1: nil map")
	var cache map[string]string // nil: declared but never made
	fmt.Println("read from nil map ok:", cache["user:1"] == "", "len:", len(cache))
	func() {
		// recover turns the panic into a value so the demo can continue.
		defer func() { fmt.Println("  recovered from panic:", recover()) }()
		cache["user:1"] = "alice" // panics: assignment to entry in nil map
	}()
	cache = make(map[string]string) // the fix: allocate before writing
	cache["user:1"] = "alice"
	fmt.Println("after make:", cache)
}
