/*
DAY 15: INTERFACES
RUN:  go run ./day15_interfaces

THEORY:
An interface is a named set of method signatures. Any type that has those
methods IS that interface, automatically. You never write "implements".
This is called implicit (structural) satisfaction. The compiler checks it at
the place where you assign a value to an interface variable or pass it to a
function that wants the interface.
Under the hood an interface value holds two words: which concrete type it
contains, and the value itself (details on Day 16). Calling a method on an
interface looks the method up in a small table (itab) and jumps to it.
Go proverb: "the bigger the interface, the weaker the abstraction". The stdlib
is full of one-method interfaces: io.Reader (Read), io.Writer (Write),
fmt.Stringer (String), error (Error). Small interfaces are easy to satisfy
and easy to fake in tests.
REAL WORLD: http.ResponseWriter is an io.Writer, so fmt.Fprintf(w, ...) writes
an HTTP response; os.File, bytes.Buffer, net.Conn, gzip.Writer all satisfy
io.Writer, so logging, report and encoding code is written once.
EXPLAIN-IT ANSWER: Why does Go need no implements keyword?
Because satisfaction is decided by the method set, not by a declaration. The
type that defines Area() does not need to know that a Shape interface exists,
and the package that defines Shape does not need to import the type. This
lets you add an interface AFTER the concrete types exist (even types from the
stdlib or third parties), keeps packages decoupled, and lets the consumer
define just the tiny interface it needs.
PITFALLS:
  - Pointer vs value receivers: if Area has a pointer receiver, only *T has it
    in its method set, so T (value) does NOT satisfy the interface.
  - Accepting a huge interface when you only call one method.
  - Using interface{}/any when a small named interface would document intent.

EXERCISES:
 1. Add a Triangle type (base, height) and put it in the shapes slice. Do not
    touch the Shape interface; notice no "implements" is needed.
 2. Make Rectangle.Area use a pointer receiver and see the compile error in
    main; then fix it by using &Rectangle{...} in the slice.
 3. Write a type upperWriter that wraps an io.Writer and upper-cases data
    before passing it on; send the report through it.
*/
package main

import (
	"bytes"
	"fmt"
	"io"
	"math"
	"os"
	"strings"
)

// Shape is a small interface: two methods. Anything with both is a Shape.
type Shape interface {
	Area() float64
	Perimeter() float64
}

// Circle is a plain struct. Notice there is no mention of Shape anywhere on it.
type Circle struct {
	Radius float64
}

// Area uses a value receiver, so both Circle and *Circle have this method.
func (c Circle) Area() float64 { return math.Pi * c.Radius * c.Radius }

// Perimeter (circumference) of the circle.
func (c Circle) Perimeter() float64 { return 2 * math.Pi * c.Radius }

// String makes Circle a fmt.Stringer. fmt.Println calls it automatically.
func (c Circle) String() string { return fmt.Sprintf("Circle(r=%.1f)", c.Radius) }

// Rectangle is another unrelated struct that happens to have the same methods.
type Rectangle struct {
	Width, Height float64
}

// Area of the rectangle.
func (r Rectangle) Area() float64 { return r.Width * r.Height }

// Perimeter of the rectangle.
func (r Rectangle) Perimeter() float64 { return 2 * (r.Width + r.Height) }

// String makes Rectangle a fmt.Stringer too.
func (r Rectangle) String() string { return fmt.Sprintf("Rect(%.1fx%.1f)", r.Width, r.Height) }

// Compile-time check idiom: assigning to the blank identifier fails to compile
// if Circle ever stops satisfying Shape. It costs nothing at runtime.
var _ Shape = Circle{}
var _ Shape = Rectangle{}

// totalArea accepts the interface, so it works for every present and future shape.
func totalArea(shapes []Shape) float64 {
	total := 0.0
	for _, s := range shapes { // s is an interface value; the call below is dynamic
		total += s.Area()
	}
	return total
}

// LineItem is one row of an invoice for the report example.
type LineItem struct {
	Name  string
	Qty   int
	Price float64
}

// writeReport takes an io.Writer: it has no idea whether it writes to the
// terminal, memory or a file. fmt.Fprintf accepts any io.Writer.
func writeReport(w io.Writer, items []LineItem) error {
	total := 0.0
	for _, it := range items {
		line := float64(it.Qty) * it.Price
		total += line
		// Fprintf returns (n, err); we must check err because writers can fail (disk full, closed pipe).
		if _, err := fmt.Fprintf(w, "%-10s %2d x %7.2f = %8.2f\n", it.Name, it.Qty, it.Price, line); err != nil {
			return err
		}
	}
	_, err := fmt.Fprintf(w, "%-10s %27.2f\n", "TOTAL", total)
	return err
}

// countBytes shows io.Reader: it consumes any source of bytes.
func countBytes(r io.Reader) (int, error) {
	// io.Copy moves bytes from a Reader to a Writer; io.Discard is a Writer that drops everything.
	n, err := io.Copy(io.Discard, r)
	return int(n), err
}

func main() {
	// A slice of the interface type can hold different concrete types together.
	shapes := []Shape{Circle{Radius: 2}, Rectangle{Width: 3, Height: 4}}
	fmt.Println("== Shapes")
	for _, s := range shapes {
		// %v on a value with String() calls String() (fmt.Stringer in action).
		fmt.Printf("%v area=%.2f perimeter=%.2f\n", s, s.Area(), s.Perimeter())
	}
	fmt.Printf("total area = %.2f\n", totalArea(shapes))

	items := []LineItem{
		{"Keyboard", 2, 49.90},
		{"Mouse", 3, 19.50},
		{"Cable", 5, 4.25},
	}

	// One function, three destinations, because all three satisfy io.Writer.
	fmt.Println("== Report to os.Stdout")
	if err := writeReport(os.Stdout, items); err != nil { // *os.File has Write
		fmt.Println("error:", err)
	}

	fmt.Println("== Report to strings.Builder")
	var sb strings.Builder // zero value is ready to use; pointer receiver Write, so pass &sb
	if err := writeReport(&sb, items); err != nil {
		fmt.Println("error:", err)
	}
	fmt.Printf("builder captured %d bytes\n", sb.Len())

	fmt.Println("== Report to bytes.Buffer")
	var buf bytes.Buffer // also an io.Writer (and an io.Reader) via pointer
	if err := writeReport(&buf, items); err != nil {
		fmt.Println("error:", err)
	}
	fmt.Printf("buffer captured %d bytes, same as builder: %t\n", buf.Len(), buf.Len() == sb.Len())

	// A Buffer is also a Reader, so we can feed it to a function that takes io.Reader.
	n, _ := countBytes(&buf)
	fmt.Println("== io.Reader: counted", n, "bytes by reading the buffer")

	// Pitfall demo (commented because it does not compile): a type whose method has a
	// pointer receiver does not satisfy the interface as a value.
	//   func (c *Counter) Write(p []byte) (int, error) {...}
	//   var w io.Writer = Counter{}   // compile error: Write has pointer receiver
	//   var w io.Writer = &Counter{}  // ok
}
