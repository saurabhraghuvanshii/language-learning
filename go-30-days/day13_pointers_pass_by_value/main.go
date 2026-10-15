/*
DAY 13: POINTERS AND PASS-BY-VALUE
RUN:  go run ./day13_pointers_pass_by_value

THEORY:
A pointer holds the memory address of a value. &x gives the address of x, *p reads or
writes the value at the address p holds. The zero value of a pointer is nil (points
nowhere); dereferencing nil panics. new(T) allocates a zeroed T and returns *T.
Go is ALWAYS pass-by-value: every argument is copied into the parameter. For an int the
copy is the number. For a pointer the copy is the address, so both copies refer to
the same target. A slice is a header (pointer, len, cap), a map is a pointer to a runtime
table, and a channel likewise: copying them copies the small handle, and the underlying
data is shared. That is why behaviour looks like "reference passing" without Go having it.

STACK vs HEAP: a function's local variables normally live on its stack, which is freed
automatically when the function returns (very cheap). If the compiler's escape analysis
sees that a variable's address outlives the function (returned, stored in a global,
captured by a goroutine), it allocates the variable on the HEAP instead, and the garbage
collector frees it later (more expensive). You do not choose; the compiler does.
See the decisions with:   go build -gcflags=-m ./day13_pointers_pass_by_value
REAL WORLD: methods with pointer receivers, json.Unmarshal(data, &cfg), optional fields
(*string, nil meaning absent), linked structures, avoiding big struct copies.
EXPLAIN-IT ANSWER: Go is pass-by-value, so why can a function change a map but not an int?
Because the map variable holds a pointer to the hash table. The function receives a
copy of that pointer, and both copies lead to the same table, so writes through it are
visible to the caller. An int variable holds the number itself, so the function
changes only its own copy. Pass &n if you want a function to change an int.
PITFALLS:
 1. Dereferencing a nil pointer panics.
 2. Appending inside a function is not visible to the caller (the slice header was copied).
 3. Reassigning a pointer parameter (p = &other) does not change the caller's pointer.

EXERCISES:
 1. TODO: write func increment(counter *int) and call it three times on one variable.
 2. TODO: run `go build -gcflags=-m ./day13_pointers_pass_by_value` and find the line saying "moved to heap".
 3. TODO: make appendItem return the new slice and use its result in main to fix PITFALL 2.
*/
package main

import "fmt"

// Invoice is a small struct to modify through a pointer.
type Invoice struct {
	ID    int
	Total float64
	Paid  bool
}

// swap exchanges two ints through their addresses.
func swap(a, b *int) {
	*a, *b = *b, *a // *a is the value at the address; assign both at once
}

// markPaid changes the caller's struct because it holds the struct's address.
func markPaid(inv *Invoice) {
	inv.Paid = true // inv.Paid is shorthand for (*inv).Paid
}

// newInvoice returns the address of a local variable. This is safe in Go:
// escape analysis moves the struct to the heap so it survives the return.
func newInvoice(id int, total float64) *Invoice {
	inv := Invoice{ID: id, Total: total} // "escapes to heap" because &inv leaves the function
	return &inv
}

// The four helpers below power the pass-by-value table.

func changeInt(n int)            { n = 99 }                 // changes the copy only
func changeSliceElem(s []int)    { s[0] = 99 }              // shared array: visible
func appendToSlice(s []int)      { s = append(s, 99) }      // header copy: invisible
func changeMap(m map[string]int) { m["x"] = 99 }            // shared table: visible
func changePointee(p *int)       { *p = 99 }                // shared target: visible
func repointParam(p *int)        { other := 1; p = &other } // only the local copy of p moves

func main() {
	fmt.Println("== new, nil and basic pointers")
	n := 10
	p := &n // p holds the address of n
	*p = 20 // write through the pointer
	fmt.Println("n via p:", n, "| *p:", *p)
	q := new(int) // allocates a zeroed int, returns *int
	*q = 5
	fmt.Println("new(int) holds:", *q)
	var nothing *Invoice // nil pointer
	fmt.Println("nothing == nil:", nothing == nil)

	fmt.Println("\n== swap via pointers")
	x, y := 3, 7
	swap(&x, &y)
	fmt.Println("x, y after swap:", x, y)

	fmt.Println("\n== modify a struct through a pointer")
	inv := Invoice{ID: 1, Total: 250.0}
	markPaid(&inv)
	fmt.Printf("%+v\n", inv)
	heap := newInvoice(2, 99.5)
	fmt.Printf("from constructor: %+v\n", *heap)

	fmt.Println("\n== PITFALL 1: nil dereference")
	func() {
		defer func() { fmt.Println("  recovered:", recover()) }()
		fmt.Println(nothing.Total) // panics: invalid memory address or nil pointer dereference
	}()

	fmt.Println("\n== PROOF: Go is always pass-by-value")
	intV := 1
	sliceV := []int{1, 2, 3}
	sliceApp := []int{1, 2, 3}
	mapV := map[string]int{"x": 1}
	ptrV := 1
	repoint := 1
	rp := &repoint
	rpBefore := rp // remember where it pointed

	changeInt(intV)
	changeSliceElem(sliceV)
	appendToSlice(sliceApp)
	changeMap(mapV)
	changePointee(&ptrV)
	repointParam(rp)

	fmt.Printf("%-26s %-12s %s\n", "call", "caller sees", "why")
	fmt.Printf("%-26s %-12v %s\n", "changeInt(n)", intV, "int copied; unchanged")
	fmt.Printf("%-26s %-12s %s\n", "changeSliceElem(s)", fmt.Sprint(sliceV), "header copied, same array; CHANGED")
	fmt.Printf("%-26s %-12s %s\n", "appendToSlice(s) PITFALL 2", fmt.Sprint(sliceApp), "caller's len still 3; unchanged")
	fmt.Printf("%-26s %-12s %s\n", "changeMap(m)", fmt.Sprint(mapV), "map pointer copied; CHANGED")
	fmt.Printf("%-26s %-12v %s\n", "changePointee(&n)", ptrV, "address copied, same target; CHANGED")
	fmt.Printf("%-26s %-12v %s\n", "repointParam(p) PITFALL 3", rp == rpBefore, "caller's pointer unchanged (true)")
}
