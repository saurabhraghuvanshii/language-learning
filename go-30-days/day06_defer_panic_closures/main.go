/*
DAY 6: DEFER, PANIC, RECOVER AND CLOSURES
RUN:  go run ./day06_defer_panic_closures

THEORY:
`defer f(x)` schedules f to run when the surrounding function returns (normally,
or while panicking). The arguments of a deferred call are evaluated IMMEDIATELY
when the defer statement executes; only the call itself is delayed. Multiple
defers run last-in-first-out (a stack), like unwinding nested setup. `panic`
stops normal execution, runs deferred calls while unwinding the stack, and
crashes the program unless a deferred function calls `recover()`, which returns
the panic value and stops the unwinding. recover only works directly inside a
deferred function. A closure is a function literal that captures variables from
its surrounding scope; the captured variable lives as long as the closure does.
REAL WORLD: defer file.Close(), mu.Unlock(), tx.Rollback(); HTTP servers recover
from handler panics so one bad request does not kill the process; closures make
counters, middleware and callbacks.
EXPLAIN-IT ANSWER: A deferred call runs when its enclosing function exits
(return, end of body, or panic), never earlier. Several defers in one function
run in reverse order of their defer statements (LIFO). Its arguments, however,
were already evaluated at the moment of the defer statement.
PITFALLS: (1) defer args evaluated immediately; (2) defer inside a loop piles
up until the function ends (can exhaust file handles); (3) recover() outside a
deferred function returns nil; (4) panic is not for ordinary errors.
EXERCISES:
1. TODO: change the deferred Println to a closure `defer func(){...}()` so it
   sees the final value, and observe the difference.
2. TODO: make safeDivide recover a custom error type and print it with errors.As.
3. TODO: write a closure generator makeAdder(n) returning func(int) int.
*/

package main

import (
	"errors"  // errors.New for the recovered error
	"fmt"     // printing
	"io"      // io.ReadAll
	"os"      // temp files
	"strings" // strings.Builder not needed; used for TrimSpace
)

// deferOrder shows LIFO order.
func deferOrder() {
	fmt.Println("[defer] start")
	for i := 1; i <= 3; i++ {
		// Each iteration pushes one deferred call; i is evaluated NOW.
		defer fmt.Println("[defer] deferred call with i =", i)
	}
	fmt.Println("[defer] end of function body (defers run after this line)")
}

// deferArgs demonstrates immediate evaluation of arguments.
func deferArgs() {
	status := "pending"
	// PITFALL: status is copied NOW, so this prints "pending".
	defer fmt.Println("[pitfall] defer arg captured early:", status)
	// A closure reads the variable when it RUNS, so it prints the final value.
	defer func() { fmt.Println("[fix] closure sees final value:", status) }()
	status = "shipped" // changed after both defers were registered
}

// readBack writes a temp file, then reads it, always closing and removing it.
func readBack() (string, error) {
	// CreateTemp makes a unique file in the OS temp dir ("" = default dir).
	f, err := os.CreateTemp("", "day06-*.txt")
	if err != nil {
		return "", err // nothing opened yet, so nothing to close
	}
	// Remove the file when we are done; registered first so it runs LAST.
	defer os.Remove(f.Name())
	// Write some content using the open handle.
	if _, err := f.WriteString("order #1001 shipped\n"); err != nil {
		f.Close() // we have not yet deferred Close, so do it by hand
		return "", err
	}
	// Close right now so data is flushed; we reopen below for reading.
	if err := f.Close(); err != nil {
		return "", err
	}
	// Reopen for reading.
	r, err := os.Open(f.Name())
	if err != nil {
		return "", err
	}
	// The key idiom: defer Close immediately after a successful open.
	// It runs on EVERY return path below, including errors.
	defer r.Close()
	data, err := io.ReadAll(r) // read the whole file into memory
	if err != nil {
		return "", err // r.Close still runs
	}
	return strings.TrimSpace(string(data)), nil
}

// safeDivide turns a runtime panic (integer divide by zero) into an error.
func safeDivide(a, b int) (result int, err error) {
	// Named result err lets the deferred function change what we return.
	defer func() {
		// recover() returns nil if no panic is in flight.
		if r := recover(); r != nil {
			err = fmt.Errorf("recovered from panic: %v", r)
		}
	}()
	return a / b, nil // b == 0 panics here at runtime
}

// mustPositive panics with a custom message; wrapped by protect below.
func mustPositive(n int) int {
	if n <= 0 {
		panic(errors.New("n must be positive"))
	}
	return n
}

// protect runs fn and converts any panic into an error (a generic wrapper).
func protect(fn func()) (err error) {
	defer func() {
		if r := recover(); r != nil {
			// r is an any; format it into an error value.
			err = fmt.Errorf("panic caught: %v", r)
		}
	}()
	fn() // run the risky function
	return nil
}

// newCounter returns a closure; count lives on after newCounter returns.
func newCounter(start int) func() int {
	count := start // captured variable, private to this counter
	return func() int {
		count++ // modifies the captured variable, not a copy
		return count
	}
}

// main runs the demos.
func main() {
	fmt.Println("== LIFO order ==")
	deferOrder()
	fmt.Println("== Argument evaluation time ==")
	deferArgs()
	fmt.Println("== File read with guaranteed close ==")
	text, err := readBack()
	fmt.Printf("[file] text=%q err=%v\n", text, err)
	fmt.Println("== Panic and recover ==")
	q, err := safeDivide(10, 2)
	fmt.Println("[recover] 10/2 =", q, "err:", err)
	q, err = safeDivide(10, 0)
	fmt.Println("[recover] 10/0 =", q, "err:", err)
	err = protect(func() { mustPositive(-5) })
	fmt.Println("[recover] protect:", err)
	err = protect(func() { mustPositive(5) })
	fmt.Println("[recover] protect ok, err:", err)
	fmt.Println("== Closures ==")
	a, b := newCounter(0), newCounter(100) // two independent counters
	fmt.Println("[closure] a:", a(), a(), a(), " b:", b(), b())
}
