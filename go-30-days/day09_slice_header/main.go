/*
DAY 9: SLICE HEADER
RUN:  go run ./day09_slice_header

THEORY:
A slice does not hold data itself. It is a tiny struct (the "slice header") of three
words: a POINTER to an element of a backing array, a LENGTH (how many elements you
may read) and a CAPACITY (how many elements exist from the pointer to the end of
the backing array). Copying a slice copies just these 24 bytes (on 64-bit), so two
slices can point at the same backing array. append writes at index len; if len would
exceed cap, Go allocates a bigger array, copies the data over and returns a header
pointing at the new array. That is why you must always write s = append(s, x).
The unsafe package below is only a microscope to print the pointer; real code never needs it.
REAL WORLD: every []byte read from a socket or file, every list of rows from a
database, every slice of request headers. Understanding the header explains most
"why did my data change?" bugs.
EXPLAIN-IT ANSWER: What three things does a slice header hold?
(1) a pointer to the first element it can see in the backing array, (2) its length
len, and (3) its capacity cap. Nothing else; the elements live in the array.
PITFALLS:
 1. Forgetting to assign the result of append: append(s, x) alone is lost.
 2. A nil slice and an empty slice both have len 0, but only the nil one == nil.
 3. Passing a slice to a function copies the header, so append inside is invisible outside.

EXERCISES:
 1. TODO: change the loop to start with make([]int, 0, 10) and watch when growth happens.
 2. TODO: print len/cap of s[1:3] for a 5-element slice and explain cap.
 3. TODO: write func addTag(tags []string) that appends, then show the caller does not see it.
*/
package main

import (
	"fmt"
	"unsafe" // used ONLY as a microscope to print the backing array address
)

// addr returns the address of the first element of the backing array.
// unsafe.SliceData gives that pointer; we only print it, never compute with it.
func addr(s []int) unsafe.Pointer {
	return unsafe.Pointer(unsafe.SliceData(s))
}

// appendInside appends to its COPY of the header; the caller's len does not change.
func appendInside(s []int) {
	s = append(s, 100) // may or may not reallocate; either way the caller keeps its own header
	fmt.Println("  inside appendInside: len", len(s), "cap", cap(s))
}

func main() {
	fmt.Println("== growth: len, cap and backing array address on each append")
	var orderIDs []int // nil slice: pointer nil, len 0, cap 0
	fmt.Printf("start: len=%d cap=%d addr=%p\n", len(orderIDs), cap(orderIDs), addr(orderIDs))
	prev := addr(orderIDs) // remember the address to detect a reallocation
	for i := 1; i <= 17; i++ {
		orderIDs = append(orderIDs, 1000+i) // PITFALL 1 avoided: result assigned back
		now := addr(orderIDs)
		note := ""
		if now != prev { // a new address means Go allocated a new backing array
			note = "  <-- new backing array"
			prev = now
		}
		fmt.Printf("append #%2d: len=%2d cap=%2d addr=%p%s\n", i, len(orderIDs), cap(orderIDs), now, note)
	}

	fmt.Println("\n== PITFALL 1: forgetting to use append's result")
	cart := []string{"book"}
	_ = append(cart, "pen") // result thrown away on purpose (a bare append(...) statement will not even compile)
	fmt.Println("cart is unchanged:", cart)
	cart = append(cart, "pen")
	fmt.Println("cart after assigning:", cart)

	fmt.Println("\n== PITFALL 2: nil slice vs empty slice")
	var nilS []int
	emptyS := []int{}
	fmt.Println("nil  : len", len(nilS), "cap", cap(nilS), "== nil?", nilS == nil)
	fmt.Println("empty: len", len(emptyS), "cap", cap(emptyS), "== nil?", emptyS == nil)
	nilS = append(nilS, 1) // append works fine on a nil slice
	fmt.Println("append to nil slice works:", nilS)

	fmt.Println("\n== PITFALL 3: header is copied into functions")
	prices := make([]int, 3, 10) // len 3, cap 10: room to grow without reallocating
	appendInside(prices)
	fmt.Println("  caller still sees len", len(prices), "->", prices)
	// The element WAS written into the shared backing array, we can see it by reslicing:
	fmt.Println("  but backing array holds it:", prices[:4])

	fmt.Println("\n== reslicing shares the pointer")
	window := prices[1:3] // pointer moves to element 1; len 2; cap = 10-1 = 9
	fmt.Printf("window len=%d cap=%d addr=%p (base addr=%p)\n", len(window), cap(window), addr(window), addr(prices))
}
