/*
DAY 8: ARRAYS
RUN:  go run ./day08_arrays

THEORY:
An array is a fixed-length sequence of values of ONE type, stored side by side in
memory. The length is part of the type: [3]int and [4]int are two different,
incompatible types. The length must be a compile-time constant. An array variable
IS the data (it is not a pointer to the data like in C). So assigning an array to
another variable, or passing it to a function, copies every element. Go chose this
because it keeps the rule simple: "assignment copies the value". Arrays are the
building block under slices (Day 9), which are the type you use day to day.
The zero value of an array is every element set to its zero value.
REAL WORLD: fixed-size things such as a SHA-256 digest ([32]byte), an IPv4 address
([4]byte), board games, lookup tables, and as map keys (arrays are comparable
with ==, slices are not).
EXPLAIN-IT ANSWER: Why does passing an array to a function copy it?
Because in Go an array variable is the whole block of values, not a reference to it,
and every function argument is evaluated and copied into the parameter. The function
gets its own independent block of memory, so changes inside do not affect the caller.
To share the original you pass a pointer (&arr) or use a slice (arr[:]).
PITFALLS:
 1. Modifying an array parameter does not change the caller's array (it was copied).
 2. [3]int and [4]int are different types; you cannot assign or pass one as the other.
 3. Large arrays passed by value are copied every call, which is slow.

EXERCISES:
 1. TODO: write sumByPointer(p *[5]float64) float64 and call it with &prices.
 2. TODO: declare [2][3]int (a 2x3 grid), fill it with nested loops and print it.
 3. TODO: use an array as a map key: map[[2]int]string for grid coordinates.
*/
package main

import "fmt"

// reverse swaps elements in place. It takes a POINTER to an array of exactly 5
// ints, so the caller's array really is modified.
func reverse(a *[5]int) {
	// i walks from the left, j from the right; stop when they meet in the middle.
	for i, j := 0, len(a)-1; i < j; i, j = i+1, j-1 {
		// a is a pointer, but Go lets us index through it directly: a[i] means (*a)[i].
		a[i], a[j] = a[j], a[i] // tuple assignment swaps without a temp variable
	}
}

// reversedCopy takes the array BY VALUE: it works on its own private copy and
// returns that copy. The caller's array stays untouched.
func reversedCopy(a [5]int) [5]int {
	for i, j := 0, len(a)-1; i < j; i, j = i+1, j-1 {
		a[i], a[j] = a[j], a[i]
	}
	return a // returning an array copies it again to the caller
}

// byValue shows the address of its OWN copy and tries to mutate it.
func byValue(a [5]int) {
	fmt.Printf("  inside byValue : &a[0] = %p  (a different address: a copy)\n", &a[0])
	a[0] = 999 // only changes the copy
}

// byPointer receives the address of the caller's array, so &a[0] is the same.
func byPointer(a *[5]int) {
	fmt.Printf("  inside byPointer: &a[0] = %p  (same address: no copy)\n", &a[0])
	a[0] = 999 // changes the caller's array
}

func main() {
	// Declaring with a literal: the length 5 is part of the type [5]int.
	monthlyOrders := [5]int{10, 20, 30, 40, 50}
	// [...] lets the compiler count the elements; the type is still [3]string.
	regions := [...]string{"EU", "US", "APAC"}
	// An array declared without values is zeroed.
	var empty [3]int
	fmt.Println("== basics")
	fmt.Println("orders:", monthlyOrders, "len:", len(monthlyOrders))
	fmt.Println("regions:", regions)
	fmt.Println("zero array:", empty)
	// %T prints the type, which includes the length.
	fmt.Printf("types: %T and %T\n", monthlyOrders, regions)

	fmt.Println("\n== assignment copies")
	backup := monthlyOrders // copies all 5 ints
	backup[0] = -1          // changes only backup
	fmt.Println("original:", monthlyOrders, "backup:", backup)
	// Arrays of the same type are comparable with ==.
	fmt.Println("equal after change?", monthlyOrders == backup)

	// PITFALL 2 (kept as a comment because it will not compile):
	//   var other [4]int = monthlyOrders   // error: cannot use [5]int as [4]int

	fmt.Println("\n== reverse in place via pointer")
	reverse(&monthlyOrders) // &x takes the address of x
	fmt.Println("reversed:", monthlyOrders)
	fmt.Println("reversedCopy of it:", reversedCopy(monthlyOrders), "(original still", monthlyOrders, ")")

	fmt.Println("\n== PITFALL 1: addresses before/after, by value vs by pointer")
	fmt.Printf("  in main        : &arr[0] = %p\n", &monthlyOrders[0])
	byValue(monthlyOrders)
	fmt.Println("  after byValue  :", monthlyOrders, "(unchanged)")
	byPointer(&monthlyOrders)
	fmt.Println("  after byPointer:", monthlyOrders, "(first element changed)")
	fmt.Printf("  in main again  : &arr[0] = %p\n", &monthlyOrders[0])

	fmt.Println("\n== PITFALL 3: cost of copying")
	var big [1000]int // 8000 bytes copied on every by-value call
	fmt.Println("a [1000]int is", len(big)*8, "bytes; prefer *[1000]int or a slice")
}
