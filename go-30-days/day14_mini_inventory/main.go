/*
DAY 14: MINI PROJECT, INVENTORY MANAGER
RUN:  go run ./day14_mini_inventory

THEORY:
This review project combines structs (Item), a map for fast lookup by SKU, and a slice
for ordered output. The map gives O(1) "find item by SKU"; because map order is random
we build a slice of keys, sort it, and print in that order. Items are stored as
*Item pointers in the map so Restock and Sell can update them in place; with
map[string]Item (values) you cannot write m[k].Qty++ because map elements are not
addressable. Errors are values: Sell returns a wrapped sentinel error that callers can
test with errors.Is. The script at the bottom runs a fixed list of operations, so the
program always terminates.
REAL WORLD: a shop backend, warehouse stock service, cart/order validation. In production
the map would be a database table and the mutex/transaction would guard concurrency.
EXPLAIN-IT ANSWER: draw the memory layout (ASCII diagram below).

	Inventory (struct, lives in main's stack or heap)
	+-------------------+
	| items  map ptr ---+--> hash table (heap)
	+-------------------+     +--------+-----------+
	                          | "A100" | *Item ----+--> Item{SKU:"A100", Name:"Keyboard", Qty:12, Min:5}
	                          | "B200" | *Item ----+--> Item{SKU:"B200", Name:"Mouse",    Qty: 3, Min:5}
	                          | "C300" | *Item ----+--> Item{...}
	                          +--------+-----------+
	skus  (slice header on the stack)      backing array (heap)
	+-----+-----+-----+                    +--------+--------+--------+
	| ptr | len | cap |  ---ptr---------> | "A100" | "B200" | "C300" |   sorted copy of the map keys
	+-----+-----+-----+                    +--------+--------+--------+
	Each string is itself a header (pointer to bytes, length). The map and the slice
	both reach the same Item values only through the SKU strings and the *Item pointers.

PITFALLS:
 1. Storing struct values in the map and trying to modify a field (does not compile).
 2. Selling more than stock without a check (quantity goes negative).
 3. Printing straight from range over the map (random order).

EXERCISES:
 1. TODO: add a Remove(sku) method that deletes the item and returns an error if absent.
 2. TODO: add a total stock value report (Qty * PriceCents) sorted by value descending.
 3. TODO: make the low-stock report accept a threshold parameter instead of using Item.Min.
*/
package main

import (
	"errors"
	"fmt"
	"sort"
)

// Sentinel errors for callers to test with errors.Is.
var (
	ErrUnknownSKU    = errors.New("unknown sku")
	ErrInsufficient  = errors.New("insufficient stock")
	ErrDuplicateItem = errors.New("item already exists")
	ErrBadQuantity   = errors.New("quantity must be positive")
)

// Item is one product line.
type Item struct {
	SKU        string
	Name       string
	Qty        int // units on hand
	Min        int // reorder threshold
	PriceCents int // money as integer cents to avoid float rounding
}

// Inventory owns the map. Fields are unexported so callers go through methods.
type Inventory struct {
	items map[string]*Item
}

// NewInventory is the constructor: it makes the map so writes never hit a nil map.
func NewInventory() *Inventory {
	return &Inventory{items: make(map[string]*Item)}
}

// Add inserts a new item.
func (inv *Inventory) Add(it Item) error {
	if _, exists := inv.items[it.SKU]; exists { // comma-ok: is this key already present?
		return fmt.Errorf("add %s: %w", it.SKU, ErrDuplicateItem)
	}
	copyIt := it // it is already a copy (value param); take the address of our own copy
	inv.items[it.SKU] = &copyIt
	return nil
}

// Restock increases quantity of an existing item.
func (inv *Inventory) Restock(sku string, qty int) error {
	if qty <= 0 {
		return fmt.Errorf("restock %s: %w", sku, ErrBadQuantity)
	}
	it, ok := inv.items[sku]
	if !ok {
		return fmt.Errorf("restock %s: %w", sku, ErrUnknownSKU)
	}
	it.Qty += qty // it is a pointer, so this updates the stored item
	return nil
}

// Sell decreases stock, refusing when there is not enough (PITFALL 2 avoided).
func (inv *Inventory) Sell(sku string, qty int) error {
	if qty <= 0 {
		return fmt.Errorf("sell %s: %w", sku, ErrBadQuantity)
	}
	it, ok := inv.items[sku]
	if !ok {
		return fmt.Errorf("sell %s: %w", sku, ErrUnknownSKU)
	}
	if it.Qty < qty {
		return fmt.Errorf("sell %d of %s (have %d): %w", qty, sku, it.Qty, ErrInsufficient)
	}
	it.Qty -= qty
	return nil
}

// sortedSKUs returns the keys in alphabetical order (PITFALL 3 avoided).
func (inv *Inventory) sortedSKUs() []string {
	skus := make([]string, 0, len(inv.items))
	for sku := range inv.items {
		skus = append(skus, sku)
	}
	sort.Strings(skus)
	return skus
}

// Report prints every item in SKU order.
func (inv *Inventory) Report() {
	fmt.Printf("  %-6s %-14s %4s %4s %10s\n", "SKU", "NAME", "QTY", "MIN", "PRICE")
	for _, sku := range inv.sortedSKUs() {
		it := inv.items[sku]
		fmt.Printf("  %-6s %-14s %4d %4d %7d.%02d\n", it.SKU, it.Name, it.Qty, it.Min, it.PriceCents/100, it.PriceCents%100)
	}
}

// LowStock returns items with Qty at or below Min, fewest units first.
func (inv *Inventory) LowStock() []Item {
	var low []Item
	for _, it := range inv.items {
		if it.Qty <= it.Min {
			low = append(low, *it) // copy the value out so callers cannot mutate our data
		}
	}
	sort.Slice(low, func(i, j int) bool {
		if low[i].Qty != low[j].Qty {
			return low[i].Qty < low[j].Qty
		}
		return low[i].SKU < low[j].SKU // tie-break for stable output
	})
	return low
}

func main() {
	inv := NewInventory()

	fmt.Println("== add items")
	seed := []Item{
		{SKU: "A100", Name: "Keyboard", Qty: 12, Min: 5, PriceCents: 4999},
		{SKU: "B200", Name: "Mouse", Qty: 3, Min: 5, PriceCents: 1999},
		{SKU: "C300", Name: "USB-C Cable", Qty: 40, Min: 10, PriceCents: 799},
		{SKU: "D400", Name: "Monitor", Qty: 2, Min: 2, PriceCents: 18900},
	}
	for _, it := range seed {
		if err := inv.Add(it); err != nil {
			fmt.Println("  error:", err)
		}
	}
	fmt.Println("  duplicate add ->", inv.Add(seed[0]))
	inv.Report()

	fmt.Println("\n== scripted operations")
	// A table of steps keeps the script short and guarantees termination.
	steps := []struct {
		op  string
		sku string
		qty int
	}{
		{"restock", "B200", 10},
		{"sell", "A100", 8},
		{"sell", "D400", 5}, // too many: must fail
		{"sell", "Z999", 1}, // unknown sku
		{"sell", "C300", 35},
		{"restock", "A100", 0}, // bad quantity
	}
	for _, s := range steps {
		var err error
		switch s.op {
		case "restock":
			err = inv.Restock(s.sku, s.qty)
		case "sell":
			err = inv.Sell(s.sku, s.qty)
		}
		switch {
		case errors.Is(err, ErrInsufficient):
			fmt.Printf("  %-8s %s x%d -> REFUSED (insufficient): %v\n", s.op, s.sku, s.qty, err)
		case err != nil:
			fmt.Printf("  %-8s %s x%d -> error: %v\n", s.op, s.sku, s.qty, err)
		default:
			fmt.Printf("  %-8s %s x%d -> ok\n", s.op, s.sku, s.qty)
		}
	}

	fmt.Println("\n== sorted inventory")
	inv.Report()

	fmt.Println("\n== low-stock report")
	for _, it := range inv.LowStock() {
		fmt.Printf("  REORDER %s %s: have %d, min %d\n", it.SKU, it.Name, it.Qty, it.Min)
	}

	fmt.Println("\n== PITFALL 1: map of struct values")
	byValue := map[string]Item{"A100": {SKU: "A100", Qty: 1}}
	// byValue["A100"].Qty++   // compile error: cannot assign to struct field in map
	tmp := byValue["A100"] // the workaround: copy out, modify, store back
	tmp.Qty++
	byValue["A100"] = tmp
	fmt.Println("  after copy-modify-store:", byValue["A100"].Qty, "(that is why we store *Item)")
}
