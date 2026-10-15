/*
DAY 16: INTERFACE INTERNALS
RUN:  go run ./day16_interface_internals

THEORY:
An interface value is a pair: (dynamic type, dynamic value). The first word
says which concrete type is stored, the second points to / holds the value.
An interface is nil ONLY when BOTH parts are empty. If you store a nil
*MyErr inside an error variable, the pair is (*MyErr, nil): the type part is
set, so the interface is not nil.
A type assertion x.(T) extracts the concrete value; the "comma ok" form
v, ok := x.(T) does not panic on mismatch. A type switch is a chain of such
assertions. Asserting to another interface (x.(fmt.Stringer)) asks "does the
stored type have these methods?".
REAL WORLD: functions that return a custom error type (validation, database
errors) are where this bug bites; also decoding JSON into any gives you
map[string]any and float64 values that you must type-switch on.
EXPLAIN-IT ANSWER: Why is err != nil true when a nil *MyErr is returned as error?
Comparing an interface to nil checks both the type word and the value word.
Returning a nil *MyErr converts it to error, producing (*MyErr, nil). The
type word is non-nil, so the interface differs from the nil interface (nil,
nil). The fix: return the literal nil when there is no error, and make the
function's return type error, not *MyErr.
PITFALLS:
  - Function returns a typed nil pointer as an interface (shown below).
  - Plain assertion x.(T) panics when the type is wrong; use comma-ok.
  - Printing with %v hides the problem; use %T to see the dynamic type.

EXERCISES:
 1. Add a case to describe() for []int and for map[string]int.
 2. Change brokenValidate so the bug disappears while still declaring a
    local variable "var e *ValidationError" (hint: only assign it when failing).
 3. Write a function that takes any and returns its length if it is a
    string, slice of int or map; otherwise -1, using a type switch.
*/
package main

import (
	"fmt"
	"strings"
)

// ValidationError is a custom error type with a pointer receiver.
type ValidationError struct {
	Field string
}

// Error makes *ValidationError satisfy the error interface.
func (e *ValidationError) Error() string {
	// e may be nil here; this line would panic on a nil receiver. Keep it in mind.
	return "invalid field: " + e.Field
}

// brokenValidate is the BUG. The return type is error (an interface), but
// inside we use a concrete pointer variable and return it unconditionally.
func brokenValidate(email string) error {
	var e *ValidationError // nil pointer, type *ValidationError
	if !strings.Contains(email, "@") {
		e = &ValidationError{Field: "email"}
	}
	// When email is valid, e is still nil, but converting it to error gives
	// the pair (*ValidationError, nil), which is NOT a nil interface.
	return e
}

// fixedValidate returns the literal nil on success. The nil literal converts
// to the interface as (nil, nil), a real nil interface.
func fixedValidate(email string) error {
	if !strings.Contains(email, "@") {
		return &ValidationError{Field: "email"}
	}
	return nil // never return a typed nil pointer through an interface
}

// describe shows a type switch over the dynamic type.
func describe(v any) string {
	switch x := v.(type) { // x has the concrete type inside each case
	case nil:
		return "nil interface"
	case int:
		return fmt.Sprintf("int %d", x)
	case string:
		return fmt.Sprintf("string %q (len %d)", x, len(x))
	case fmt.Stringer: // interface case: matches any type with String()
		return "Stringer " + x.String()
	case error:
		return "error " + x.Error()
	default: // x keeps type any here
		return fmt.Sprintf("other %T", x)
	}
}

// Money has a String method, so describe's Stringer case catches it.
type Money int

// String renders cents as a dollar amount.
func (m Money) String() string { return fmt.Sprintf("$%d.%02d", int(m)/100, int(m)%100) }

func main() {
	fmt.Println("== The nil-interface bug")
	err := brokenValidate("alice@example.com") // a VALID email, so we expect no error
	fmt.Println("broken: err != nil ->", err != nil)
	fmt.Printf("broken: dynamic type %T, value is nil pointer? %t\n", err, err.(*ValidationError) == nil)
	// err.(*ValidationError) is a type assertion extracting the pointer inside the pair.

	err = fixedValidate("alice@example.com")
	fmt.Println("fixed:  err != nil ->", err != nil)
	fmt.Printf("fixed:  dynamic type %T\n", err) // <nil>: both words empty

	err = fixedValidate("bad-email")
	fmt.Println("fixed with bad input:", err)

	fmt.Println("== Type assertion")
	var v any = "hello" // any is the empty interface: every type satisfies it
	if s, ok := v.(string); ok {
		fmt.Println("comma-ok: it is a string of length", len(s))
	}
	n, ok := v.(int) // wrong type with comma-ok: no panic, n is the zero value
	fmt.Println("comma-ok wrong type: n =", n, "ok =", ok)

	// A plain assertion panics when wrong. We recover to show it safely.
	func() {
		defer func() { fmt.Println("recovered:", recover()) }()
		_ = v.(int) // panics: interface conversion
	}()

	fmt.Println("== Type switch")
	values := []any{nil, 42, "go", Money(1999), &ValidationError{Field: "age"}, 3.14}
	for _, val := range values {
		fmt.Println(describe(val))
	}

	fmt.Println("== Interface to interface assertion")
	var e error = &ValidationError{Field: "name"}
	if st, ok := e.(fmt.Stringer); ok { // does the stored type also have String()?
		fmt.Println("also a Stringer:", st)
	} else {
		fmt.Printf("%T is not a Stringer\n", e)
	}
}
