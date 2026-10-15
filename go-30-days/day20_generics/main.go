/*
DAY 20: GENERICS
RUN:  go run ./day20_generics

THEORY:
Generics (Go 1.18+) let one function or type work for many types while staying
statically type checked. You declare type parameters in square brackets:
func Map[T, U any](...). A constraint says what operations the type allows:
  - any        : no operations (just store, pass around)
  - comparable : supports == and != (needed for map keys, Contains)
  - cmp.Ordered: supports < > (ints, floats, strings)
  - custom     : an interface listing types, e.g. interface{ ~int | ~float64 }

The "~" means "any type whose underlying type is int", so type Cents int fits.
The compiler usually infers type arguments from the call, so you rarely write
Map[int,string](...). Under the hood Go generates code per "shape" of type
(GC shape stenciling with dictionaries); you do not need to care.
Generic types: type Stack[T any] struct{ items []T }, methods repeat [T].
REAL WORLD: slices.Contains, slices.SortFunc, maps.Keys, sync.Map wrappers,
typed caches, Result/Optional helpers, generic repositories.
EXPLAIN-IT ANSWER: When is an interface better than a generic function?
Use an interface when behaviour differs per type and you want to call
methods dynamically (a []Shape with different Area implementations, or an
io.Writer). Use generics when the algorithm is the same and only the data
type changes (Map, Filter, Stack). If you would have to call a method on T,
an interface parameter is simpler and also lets one slice hold mixed types,
which generics cannot do ([]T has ONE T).
PITFALLS:
  - Over-generalizing: if only one type is used, write the plain function.
  - Methods cannot have their own type parameters (only the type can).
  - Forgetting the constraint: a > b on T any does not compile.
  - Generic slices cannot hold mixed concrete types; interfaces can.

EXERCISES:
 1. Write Keys[K comparable, V any](m map[K]V) []K and print it sorted.
 2. Add Peek() (T, bool) to Stack[T] that returns the top without removing it.
 3. Write Sum[T Number](xs []T) T and call it with []Cents.
*/
package main

import (
	"cmp"
	"fmt"
	"strings"
)

// Map applies f to every element and returns a new slice of a possibly different type.
func Map[T, U any](xs []T, f func(T) U) []U {
	out := make([]U, 0, len(xs)) // preallocate: we know the final length
	for _, x := range xs {
		out = append(out, f(x))
	}
	return out
}

// Filter keeps the elements for which keep returns true.
func Filter[T any](xs []T, keep func(T) bool) []T {
	var out []T
	for _, x := range xs {
		if keep(x) {
			out = append(out, x)
		}
	}
	return out
}

// Reduce folds the slice into one value, starting from initial.
func Reduce[T, A any](xs []T, initial A, f func(A, T) A) A {
	acc := initial
	for _, x := range xs {
		acc = f(acc, x)
	}
	return acc
}

// Contains needs ==, so T must be comparable.
func Contains[T comparable](xs []T, target T) bool {
	for _, x := range xs {
		if x == target {
			return true
		}
	}
	return false
}

// MaxOf needs <, so T must be cmp.Ordered. cmp.Compare gives -1, 0 or 1.
func MaxOf[T cmp.Ordered](xs []T) (T, bool) {
	var zero T // the zero value of whatever T is
	if len(xs) == 0 {
		return zero, false
	}
	best := xs[0]
	for _, x := range xs[1:] {
		if cmp.Compare(x, best) > 0 {
			best = x
		}
	}
	return best, true
}

// Number is a custom constraint: any type whose underlying type is one of these.
type Number interface {
	~int | ~int64 | ~float64
}

// Sum works for ints, floats and named types such as Cents.
func Sum[T Number](xs []T) T {
	var total T
	for _, x := range xs {
		total += x // + is allowed because every type in Number supports it
	}
	return total
}

// Cents is a named type; ~int in Number lets it through.
type Cents int

// Stack is a generic type. The type parameter belongs to the type; methods repeat it.
type Stack[T any] struct {
	items []T
}

// Push adds an element on top.
func (s *Stack[T]) Push(v T) { s.items = append(s.items, v) }

// Pop removes and returns the top. The bool says whether there was anything.
func (s *Stack[T]) Pop() (T, bool) {
	var zero T
	if len(s.items) == 0 {
		return zero, false
	}
	last := len(s.items) - 1
	v := s.items[last]
	s.items = s.items[:last] // shrink the slice by one
	return v, true
}

// Len reports how many elements are stored.
func (s *Stack[T]) Len() int { return len(s.items) }

// ---- Case where a plain interface beats generics ----

// Notifier differs by behaviour per type, so an interface is the right tool.
type Notifier interface {
	Notify(msg string) string
}

// EmailNotifier sends by email (simulated).
type EmailNotifier struct{ To string }

// Notify formats an email.
func (e EmailNotifier) Notify(msg string) string { return "email to " + e.To + ": " + msg }

// SMSNotifier sends by SMS (simulated).
type SMSNotifier struct{ Phone string }

// Notify formats an SMS.
func (s SMSNotifier) Notify(msg string) string { return "sms to " + s.Phone + ": " + msg }

// notifyAll takes a slice of the interface so DIFFERENT types live together.
// A generic version func notifyAll[T Notifier](ns []T) would force ONE concrete
// type per call, so you could not mix email and sms in the same slice.
func notifyAll(ns []Notifier, msg string) []string {
	return Map(ns, func(n Notifier) string { return n.Notify(msg) })
}

// Order is used in the invoice example below.
type Order struct {
	Customer string
	Total    Cents
}

func main() {
	fmt.Println("== Map / Filter / Reduce")
	orders := []Order{{"ana", 12000}, {"bo", 4550}, {"cy", 99900}, {"di", 700}}
	big := Filter(orders, func(o Order) bool { return o.Total >= 5000 })
	names := Map(big, func(o Order) string { return strings.ToUpper(o.Customer) })
	fmt.Println("big order customers:", names)
	totalCents := Reduce(orders, Cents(0), func(acc Cents, o Order) Cents { return acc + o.Total })
	fmt.Printf("all orders total: %d cents ($%.2f)\n", totalCents, float64(totalCents)/100)

	fmt.Println("== Contains / MaxOf / Sum")
	fmt.Println("contains \"bo\":", Contains(names, "BO"), "| contains 3:", Contains([]int{1, 2, 3}, 3))
	maxPrice, _ := MaxOf([]float64{2.5, 9.99, 4})
	maxName, _ := MaxOf([]string{"pear", "apple", "zebra"})
	_, ok := MaxOf([]int{})
	fmt.Println("max price:", maxPrice, "| max name:", maxName, "| empty ok:", ok)
	fmt.Println("sum ints:", Sum([]int{1, 2, 3}), "| sum cents:", Sum([]Cents{100, 250}))

	fmt.Println("== Generic Stack")
	var undo Stack[string] // zero value is a usable empty stack
	undo.Push("type 'a'")
	undo.Push("delete line")
	fmt.Println("len:", undo.Len())
	for undo.Len() > 0 {
		op, _ := undo.Pop()
		fmt.Println("undo:", op)
	}
	_, ok = undo.Pop()
	fmt.Println("pop on empty ok:", ok)
	nums := &Stack[int]{} // same code, different T
	nums.Push(7)
	fmt.Println("int stack len:", nums.Len())

	fmt.Println("== Interface is better: mixed notifiers")
	ns := []Notifier{EmailNotifier{"a@x.com"}, SMSNotifier{"+15550100"}}
	for _, line := range notifyAll(ns, "invoice paid") {
		fmt.Println(line)
	}
}
