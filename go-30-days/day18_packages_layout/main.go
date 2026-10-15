/*
DAY 18: PACKAGES LAYOUT
RUN:  go run ./day18_packages_layout

LAYOUT (import path prefix is go30/day18_packages_layout):

	main.go                     package main, wires everything together
	inventory/inventory.go      business logic, defines the Store interface
	storage/storage.go          concrete in-memory storage, imports nothing of ours
	internal/idgen/idgen.go     private helper, importable only from inside this tree

THEORY:
A package is a directory; every .go file in it starts with the same
"package name" line. You import by module path + directory. Visibility is by
the first letter: Upper case = exported (usable by other packages), lower
case = unexported (private to the package). No public/private keywords.
A directory named internal/ is special: the go tool only allows packages
rooted at the PARENT of internal/ to import it. Anyone else gets the compile
error "use of internal package ... not allowed". It is how you keep helpers
private to your module or project.
Go forbids import cycles: if A imports B, B must not (even indirectly) import A.
REAL WORLD: big services split into packages like user, billing, storage; a
cmd/ folder holds main packages, internal/ holds code not meant to be
imported by other modules.

CIRCULAR IMPORT EXAMPLE:
Suppose inventory calls storage.Memory directly, so inventory imports storage.
Later storage wants to use inventory.Item (e.g. Put(it inventory.Item)), so
storage imports inventory. Compile error:

	package go30/day18_packages_layout
	    imports go30/day18_packages_layout/inventory
	    imports go30/day18_packages_layout/storage
	    imports go30/day18_packages_layout/inventory: import cycle not allowed

THE FIX: inventory declares a small interface Store (Put/Get with built-in types)
and depends on THAT, not on package storage. storage.Memory simply has those
methods, so it satisfies inventory.Store without importing inventory. Only
main (the composition root) imports both and plugs one into the other. The
arrows now point: main -> inventory, main -> storage, and nothing points back.

EXPLAIN-IT ANSWER: How does an interface defined in package A let B and C avoid
importing each other?
The interface lists behaviour, not types. C (the consumer) says "I need
something with Put and Get" in its own package. B (the provider) happens to
have those methods; satisfaction is implicit so B needs no import of C to
declare it. A third place (main) imports both and connects them. Dependencies
flow to the interface, and the interface lives with the consumer.
PITFALLS:
  - Cycles from putting shared types in the wrong package.
  - Exporting too much: every Upper-case name becomes API you must keep.
  - Putting interfaces next to the implementation instead of the consumer.
  - Fields of structs must be Upper case for encoding/json to see them.

EXERCISES:
 1. Add a method Keys() []string to the inventory.Store interface and see
    what the compiler says about other implementations (it still works for Memory).
 2. Try adding `import "go30/day18_packages_layout/internal/idgen"` to
    a package outside this tree (e.g. a copy of day17) and read the error.
 3. Write a second Store implementation (a fake that always fails) in main.go
    and pass it to inventory.New to see Add return an error.
*/
package main

import (
	"errors"
	"fmt"

	"go30/day18_packages_layout/inventory" // exported names are used as inventory.X
	"go30/day18_packages_layout/storage"
)

func main() {
	// main is the composition root: the only place that knows both concrete pieces.
	mem := storage.NewMemory()
	svc := inventory.New(mem) // *storage.Memory is accepted as inventory.Store

	fmt.Println("== Add items")
	var ids []string
	for _, p := range []struct {
		name  string
		stock int
		price float64
	}{{"Widget", 10, 2.50}, {"Gadget", 3, 19.99}, {"Gizmo", 7, 5.00}} {
		it, err := svc.Add(p.name, p.stock, p.price)
		if err != nil {
			fmt.Println("add failed:", err)
			return
		}
		ids = append(ids, it.ID)
		fmt.Printf("added %+v\n", it)
	}

	fmt.Println("== Find and total")
	it, err := svc.Find(ids[1])
	fmt.Printf("find %s -> %+v err=%v\n", ids[1], it, err)
	total, _ := svc.TotalValue(ids)
	fmt.Printf("total stock value: %.2f\n", total)

	fmt.Println("== Missing item keeps a matchable error through the wrap")
	_, err = svc.Find("SKU-9999")
	fmt.Println("error:", err)
	fmt.Println("is storage.ErrNotFound:", errors.Is(err, storage.ErrNotFound))

	fmt.Println("== Storage keys (exported method of the concrete type)")
	fmt.Println(mem.Keys())

	// Unexported names are invisible here. These lines would not compile:
	//   storage.helper()        // undefined (cannot refer to unexported name)
	//   svc.store               // svc.store undefined (unexported field)
	// main sits inside the project tree, so it COULD import internal/idgen;
	// code in another module could not.
}
