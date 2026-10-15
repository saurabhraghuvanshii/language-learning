/*
DAY 10: SLICE SHARING AND APPEND
RUN:  go run ./day10_slice_sharing_append

THEORY:
Slices are windows onto a backing array (Day 9). Reslicing s[a:b] makes a new window on
the SAME array: len = b-a, cap = cap(s)-a. The full slice expression s[a:b:c] also
caps the window at c-a, so append cannot run past c into data you do not own.
append(s, x) has two outcomes: if len(s) < cap(s) it writes x into the existing array
at index len (visible to every slice sharing that array); if len == cap it allocates a
larger array (roughly double for small slices), copies, and writes there, leaving old
sharers behind. copy(dst, src) copies min(len(dst), len(src)) elements into a
separate array. make([]T, len, cap) allocates a zeroed array up front.
REAL WORLD: parsing a log file into lines, splitting a request body, batching DB
rows. Aliasing bugs show up when a helper returns a sub-slice and the caller appends to it.
EXPLAIN-IT ANSWER: When does append allocate a new array and when not?
It allocates a new array only when len(s) == cap(s) (no room left). If there is spare
capacity it writes in place in the existing backing array and returns a header with
len+1. So whether two slices still share data after an append depends on spare capacity.
PITFALLS:
 1. Aliasing: modifying one slice changes the other because they share one array.
 2. append into spare capacity silently overwrites a neighbour slice's elements.
 3. Keeping a tiny slice of a huge array keeps the whole array alive (memory leak).

EXERCISES:
 1. TODO: change `head := ...[:3]` to the 3-index form [:3:3] and observe that append no longer overwrites.
 2. TODO: write func lastN(s []int, n int) []int that returns an independent copy.
 3. TODO: print len and cap after each step of s = make([]int, 0, 4) appending 6 values.
*/
package main

import "fmt"

func main() {
	fmt.Println("== PITFALL 1: aliasing, two slices share one backing array")
	invoices := []int{100, 200, 300, 400, 500}
	firstThree := invoices[:3] // window on the same array: len 3, cap 5
	lastThree := invoices[2:]  // another window, overlaps at index 2
	firstThree[2] = 999        // writes the shared array at index 2
	fmt.Println("invoices  :", invoices)
	fmt.Println("firstThree:", firstThree, "lastThree:", lastThree, "(both see 999)")

	fmt.Println("\n== PITFALL 2: append into spare capacity overwrites the neighbour")
	data := []string{"a", "b", "c", "d", "e"}
	head := data[:2]           // len 2 cap 5: spare room that belongs to "c","d","e"
	head = append(head, "XXX") // len < cap, so it writes in place at index 2
	fmt.Println("head:", head, "data:", data, "(c was destroyed)")

	fmt.Println("\n== FIX A: three-index slice caps the capacity")
	data = []string{"a", "b", "c", "d", "e"}
	safeHead := data[:2:2] // len 2, cap 2 (the third number is the cap limit)
	fmt.Println("safeHead len", len(safeHead), "cap", cap(safeHead))
	safeHead = append(safeHead, "XXX") // len == cap, so append MUST allocate a new array
	fmt.Println("safeHead:", safeHead, "data:", data, "(data intact)")

	fmt.Println("\n== FIX B: copy into a fresh slice")
	data = []string{"a", "b", "c", "d", "e"}
	clone := make([]string, 2)   // new backing array, len 2
	n := copy(clone, data[:2])   // copy returns how many elements it copied
	clone = append(clone, "XXX") // cap 2 -> reallocates anyway; independent
	clone[0] = "changed"         // does not touch data
	fmt.Println("copied", n, "clone:", clone, "data:", data)

	fmt.Println("\n== make with len and cap")
	buf := make([]int, 0, 4) // len 0, cap 4: appends up to 4 do not allocate
	for i := 1; i <= 6; i++ {
		buf = append(buf, i*10)
		fmt.Printf("  after append %d: len=%d cap=%d %v\n", i, len(buf), cap(buf), buf)
	}

	fmt.Println("\n== when does append share? (spare capacity vs full)")
	a := make([]int, 3, 4) // one spare slot
	b := append(a, 1)      // fits: shares a's array
	c := append(a, 2)      // fits again: overwrites the same slot b used!
	fmt.Println("b:", b, "c:", c, "(b was clobbered by c)")
	full := []int{1, 2, 3} // len == cap
	d := append(full, 4)   // must reallocate
	d[0] = 77
	fmt.Println("full:", full, "d:", d, "(independent)")

	fmt.Println("\n== PITFALL 3: tiny slice of a huge array keeps it alive")
	huge := make([]byte, 10_000_000) // 10 MB backing array, e.g. a whole file read into memory
	copy(huge, "ERROR: disk full\nrest of file...")
	leaky := huge[:16] // only 16 bytes visible, but the whole 10 MB stays reachable
	fmt.Printf("leaky : len=%d cap=%d (the cap reveals the 10MB array still held)\n", len(leaky), cap(leaky))
	tight := make([]byte, 16) // copy just what we need
	copy(tight, huge[:16])
	huge = nil // now the 10MB array has no references and the GC can free it
	fmt.Printf("tight : len=%d cap=%d %q\n", len(tight), cap(tight), tight)
	_ = leaky // leaky is still holding the memory; the fix is to keep only tight
	// Go 1.21+ also offers slices.Clone(huge[:16]) for the same job.
}
