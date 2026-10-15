/*
DAY 5: FUNCTIONS AND ERRORS
RUN:  go run ./day05_functions_errors

THEORY:
A function is `func name(params) results`. It can return several values, which
is how Go reports failure: the convention is (value, error). `error` is a tiny
built-in interface with one method, Error() string; nil means success. Results
can be named, which declares them as variables and allows a bare `return`.
Variadic parameters (`prices ...float64`) arrive as a slice, and you pass an
existing slice with `prices...`. Functions are values: you can store and pass
them. Errors are wrapped with fmt.Errorf("...: %w", err) and inspected with
errors.Is and errors.As.
REAL WORLD: every call to a database, file system or HTTP API returns an error
you must handle; invoice and pricing code uses variadics and money rounding.
EXPLAIN-IT ANSWER: Go treats failure as a normal, expected result, not a
hidden control-flow jump. Returning an error makes every failure point visible
in the signature and at the call site (`if err != nil`), so reviewers see exactly
what can fail and what happens next; there are no invisible exits from the
middle of a function like exceptions, and performance stays predictable.
PITFALLS: (1) ignoring the error with _ ; (2) using the result when err != nil;
(3) a bare `return` with named results is easy to misread in long functions;
(4) float64 for money has rounding error (shown below; use integer cents).
EXERCISES:
1. TODO: add a "mod" operation to calculate and return an error for mod by 0.
2. TODO: wrap the divide error with %w in a caller and test with errors.Is.
3. TODO: change invoice to use int64 cents and compare the totals.
*/

package main

import (
	"errors" // errors.New, errors.Is
	"fmt"    // printing, Errorf
)

// ErrDivideByZero is a sentinel error callers can compare against with errors.Is.
var ErrDivideByZero = errors.New("division by zero")

// calculate applies op to a and b and returns (result, error).
func calculate(op string, a, b float64) (float64, error) {
	switch op {
	case "+":
		return a + b, nil // nil error = success
	case "-":
		return a - b, nil
	case "*":
		return a * b, nil
	case "/":
		if b == 0 {
			// Return a zero result PLUS an error; callers must check the error.
			return 0, ErrDivideByZero
		}
		return a / b, nil
	}
	// %q quotes the unknown operator so the message is unambiguous.
	return 0, fmt.Errorf("unknown operator %q", op)
}

// minMax uses NAMED results: lo and hi are pre-declared and zero-valued.
func minMax(nums ...int) (lo, hi int, err error) {
	if len(nums) == 0 {
		err = errors.New("no numbers given")
		return // bare return sends back lo, hi, err as they are now
	}
	lo, hi = nums[0], nums[0] // start with the first element
	for _, n := range nums[1:] {
		if n < lo {
			lo = n
		}
		if n > hi {
			hi = n
		}
	}
	return // bare return again: lo, hi and nil err
}

// invoiceTotal sums item prices (variadic), adds tax and subtracts a discount.
// taxRate 0.18 means 18 percent. Returns an error on invalid input.
func invoiceTotal(taxRate, discount float64, prices ...float64) (float64, error) {
	if taxRate < 0 || discount < 0 {
		return 0, errors.New("tax rate and discount must not be negative")
	}
	subtotal := 0.0
	for i, p := range prices { // prices is just a []float64 inside the function
		if p < 0 {
			// %w would wrap an error; here we build a fresh message with the index.
			return 0, fmt.Errorf("item %d has negative price %.2f", i, p)
		}
		subtotal += p
	}
	if discount > subtotal {
		return 0, fmt.Errorf("discount %.2f exceeds subtotal %.2f", discount, subtotal)
	}
	taxed := (subtotal - discount) * (1 + taxRate)
	return taxed, nil
}

// main shows the happy path, the error path and the pitfalls.
func main() {
	fmt.Println("== Calculator ==")
	// A table of test cases; each is an anonymous struct value.
	cases := []struct {
		op   string
		a, b float64
	}{{"+", 7, 3}, {"/", 10, 4}, {"/", 5, 0}, {"^", 2, 3}}
	for _, c := range cases {
		res, err := calculate(c.op, c.a, c.b)
		if err != nil { // ALWAYS check err before using res
			// errors.Is matches the sentinel even if it was wrapped.
			if errors.Is(err, ErrDivideByZero) {
				fmt.Printf("[calc] %v %s %v -> handled: %v\n", c.a, c.op, c.b, err)
			} else {
				fmt.Printf("[calc] %v %s %v -> error: %v\n", c.a, c.op, c.b, err)
			}
			continue
		}
		fmt.Printf("[calc] %v %s %v = %v\n", c.a, c.op, c.b, res)
	}

	fmt.Println("== Named returns ==")
	lo, hi, err := minMax(4, 9, -2, 7)
	fmt.Println("[minmax]", lo, hi, err)
	_, _, err = minMax() // no args -> error; _ discards values we do not need
	fmt.Println("[minmax] empty:", err)

	fmt.Println("== Invoice (variadic) ==")
	total, err := invoiceTotal(0.18, 10, 499.99, 25.50, 120)
	fmt.Printf("[invoice] total=%.2f err=%v\n", total, err)
	// Passing an existing slice requires the trailing "...".
	items := []float64{10, -5}
	_, err = invoiceTotal(0.18, 0, items...)
	fmt.Println("[invoice] bad item:", err)
	total, _ = invoiceTotal(0.0, 0) // zero variadic args is legal; total is 0
	fmt.Println("[invoice] empty cart:", total)

	fmt.Println("== Pitfalls ==")
	// PITFALL: discarding the error with _ hides the failure; result is just 0.
	bad, _ := calculate("/", 1, 0)
	fmt.Println("[pitfall] ignored error, silently got:", bad)
	// PITFALL: float money drifts; 0.1+0.2 is not exactly 0.3.
	fmt.Println("[pitfall] 0.1+0.2 == 0.3 ?", 0.1+0.2 == 0.3, "(untyped consts, exact)")
	x, y := 0.1, 0.2
	fmt.Println("[pitfall] float64 vars 0.1+0.2 == 0.3 ?", x+y == 0.3, "- use integer cents for money")
}
