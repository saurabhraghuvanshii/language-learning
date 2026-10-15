/*
DAY 2: VARIABLES, CONSTANTS AND ZERO VALUES
RUN:  go run ./day02_variables_zero_values

THEORY:
`var x int` declares a variable; `x := 5` is shorthand that declares AND infers
the type, but only works inside functions. `const` values are fixed at compile
time and may be "untyped" (they adapt to where they are used). `iota` is a
counter inside a const block that goes 0,1,2... and is how Go writes enums.
Every variable is always initialised: if you do not give a value, Go sets it to
the type's ZERO VALUE: 0 for numbers, false for bool, "" for string, nil for
pointers, slices, maps, funcs, interfaces and channels. Memory is zeroed by the
runtime, so there is no "uninitialised garbage" like in C.
REAL WORLD: config structs where an unset field means "default", enum-like
status codes for orders/HTTP states, unit conversions in billing and IoT code.
EXPLAIN-IT ANSWER: Because Go guarantees memory safety and predictability: a
variable can never hold random leftover bytes. Choosing a useful zero value
(empty string, 0, nil slice you can append to, sync.Mutex ready to lock) lets
types work without constructors, and removes a whole class of bugs.
PITFALLS: (1) `:=` inside an inner block creates a NEW shadowed variable
instead of assigning to the outer one; (2) integer division truncates
(5/2 == 2); (3) mixing int and float64 needs explicit conversion.
EXERCISES:
1. TODO: add a Status constant "Refunded" to the iota block and print it.
2. TODO: add a miles-to-km converter and check that round trip is ~identical.
3. TODO: make a struct with a time.Duration field, print its zero value.
*/

package main

import "fmt" // fmt for printing

// Status is our own named type built on int, used as an enum.
type Status int

// iota starts at 0 in each const block and increases by 1 per line.
const (
	Pending Status = iota // 0: the zero value, so "unset" means Pending
	Paid                  // 1
	Shipped               // 2
)

// String lets fmt print the name instead of the number (Stringer interface).
func (s Status) String() string {
	// switch picks the label that matches s.
	switch s {
	case Pending:
		return "Pending"
	case Paid:
		return "Paid"
	case Shipped:
		return "Shipped"
	}
	// Anything else is a bug; make it visible.
	return "Unknown"
}

// Conversion factors are untyped constants: fixed forever, checked at compile time.
const (
	kmToMiles      = 0.621371  // miles per kilometre
	celsiusToFMult = 9.0 / 5.0 // multiplier part of C -> F
	celsiusToFAdd  = 32.0      // offset part of C -> F
)

// kmToMilesConv converts kilometres to miles.
func kmToMilesConv(km float64) float64 {
	return km * kmToMiles // float * untyped const works without a cast
}

// celsiusToFahrenheit converts degrees Celsius to Fahrenheit.
func celsiusToFahrenheit(c float64) float64 {
	return c*celsiusToFMult + celsiusToFAdd // F = C*9/5 + 32
}

// zeroValues shows what each type holds when never assigned.
func zeroValues() {
	var i int            // declared, not assigned -> 0
	var f float64        // -> 0
	var b bool           // -> false
	var s string         // -> "" (empty, NOT nil)
	var p *int           // pointer -> nil
	var sl []string      // slice -> nil, but len() and append still work
	var m map[string]int // map -> nil; reading is fine, writing panics
	var st Status        // our enum -> 0 == Pending
	// %q quotes strings so empty is visible; %v prints any value.
	fmt.Printf("[zero] int=%d float=%v bool=%v string=%q ptr=%v\n", i, f, b, s, p)
	fmt.Printf("[zero] slice=%v (nil? %t, len %d) map=%v (nil? %t) status=%v\n",
		sl, sl == nil, len(sl), m, m == nil, st)
	// Reading a missing key from a nil map is safe and yields the zero value.
	fmt.Println("[zero] nil map read:", m["missing"])
}

// declarations shows var versus :=.
func declarations() {
	var name string = "meter" // full form: var, name, type, value
	var count = 3             // type inferred (int)
	var width, height = 4, 5  // several at once
	unit := "km"              // short form, only inside functions
	ratio := 2.5              // inferred float64
	fmt.Println("[decl]", name, count, width*height, unit, ratio)
	// Types differ, so this would NOT compile: count * ratio (int * float64).
	fmt.Println("[decl] explicit conversion:", float64(count)*ratio)
}

// pitfalls shows shadowing and integer division.
func pitfalls() {
	total := 10 // outer variable
	if true {
		// PITFALL: := here makes a brand-new total that only lives in this block.
		total := 99
		_ = total // silence "declared and not used" for the inner copy
	}
	// Still 10: the assignment above never touched the outer variable.
	fmt.Println("[pitfall] shadowing, outer total is still:", total)
	if true {
		total = 99 // plain = assigns to the outer variable (the fix)
	}
	fmt.Println("[pitfall] after plain assignment:", total)
	// PITFALL: both operands are ints, so the result is an int (truncated).
	fmt.Println("[pitfall] 5/2 =", 5/2, "but 5/2.0 =", 5/2.0)
}

// main runs the unit converter and the demonstrations.
func main() {
	fmt.Println("== Unit converter ==")
	// Table of sample inputs; a slice literal of float64.
	for _, km := range []float64{1, 5, 42.195} {
		// %.3f prints three decimal places.
		fmt.Printf("[convert] %.3f km = %.3f miles\n", km, kmToMilesConv(km))
	}
	for _, c := range []float64{-40, 0, 36.6, 100} {
		fmt.Printf("[convert] %.1f C = %.1f F\n", c, celsiusToFahrenheit(c))
	}
	fmt.Println("== Status enum ==")
	fmt.Println("[iota]", Pending, Paid, Shipped, "numeric:", int(Shipped))
	fmt.Println("== Declarations ==")
	declarations()
	fmt.Println("== Zero values ==")
	zeroValues()
	fmt.Println("== Pitfalls ==")
	pitfalls()
}
