/*
DAY 12: STRUCTS AND METHODS
RUN:  go run ./day12_structs_methods

THEORY:
A struct groups named fields into one value; fields sit next to each other in memory.
A method is a function with a receiver: func (a *Account) Deposit(...). The receiver is
just a hidden first parameter, and like every parameter it is COPIED. With a value
receiver (a Account) the method gets a copy, so it cannot change the caller's struct.
With a pointer receiver (a *Account) it gets a copy of the address, so it can. Go
has no constructors; by convention you write NewX functions that validate and return
a ready value (or pointer). Go has no inheritance; instead struct EMBEDDING puts one
type inside another with no field name, and the outer type gains ("promotes") the inner
fields and methods. It is composition with shortcut syntax, not subclassing.
REAL WORLD: domain types (Order, User, Account), config structs, http.Server,
sql.DB wrappers, and embedding sync.Mutex or a logger into a service type.
EXPLAIN-IT ANSWER: When use a value receiver vs a pointer receiver?
Use a pointer receiver when the method must mutate the receiver, when the struct is
large (avoid copying), or when it contains something that must not be copied (a
sync.Mutex). Use a value receiver for small immutable-like types (Point, Money, time.Time)
where copying is cheap and the method only reads. Do not mix both on one type: choose
pointer if any method needs it, for consistency.
PITFALLS:
 1. A value-receiver method that "mutates" only changes its private copy.
 2. Using a struct literal directly skips validation; use a constructor.
 3. Embedding is not inheritance: the inner type does not know about the outer one.

EXERCISES:
 1. TODO: add a Transfer(to *BankAccount, amount float64) error method.
 2. TODO: add Title() to Person and to Employee (shadowing it), and call the embedded one explicitly e.Person.Title() explicitly.
 3. TODO: change DepositBroken to a pointer receiver and watch the balance update.
*/
package main

import (
	"errors"
	"fmt"
)

// Sentinel errors let callers test with errors.Is.
var (
	ErrInvalidAmount     = errors.New("amount must be positive")
	ErrInsufficientFunds = errors.New("insufficient funds")
)

// BankAccount keeps lowercase (unexported) fields so other packages must use methods.
type BankAccount struct {
	owner   string
	balance float64 // real systems use integer cents; float keeps the demo short
}

// NewBankAccount is the constructor by convention: validates and returns a pointer.
func NewBankAccount(owner string, opening float64) (*BankAccount, error) {
	if opening < 0 {
		return nil, ErrInvalidAmount
	}
	return &BankAccount{owner: owner, balance: opening}, nil // & on a literal allocates and returns its address
}

// Deposit has a POINTER receiver so the caller's balance really changes.
func (a *BankAccount) Deposit(amount float64) error {
	if amount <= 0 {
		return ErrInvalidAmount
	}
	a.balance += amount // a.balance is shorthand for (*a).balance
	return nil
}

// Withdraw refuses to overdraw and wraps the sentinel with context using %w.
func (a *BankAccount) Withdraw(amount float64) error {
	if amount <= 0 {
		return ErrInvalidAmount
	}
	if amount > a.balance {
		return fmt.Errorf("withdraw %.2f from %s (balance %.2f): %w", amount, a.owner, a.balance, ErrInsufficientFunds)
	}
	a.balance -= amount
	return nil
}

// Balance only reads, and the struct is tiny: a value receiver is fine here.
func (a BankAccount) Balance() float64 { return a.balance }

// DepositBroken has a VALUE receiver: the bug (PITFALL 1). It edits a copy.
func (a BankAccount) DepositBroken(amount float64) {
	a.balance += amount // changes the copy that dies when the method returns
}

// Person is the type that will be embedded.
type Person struct {
	Name string
	Age  int
}

// Greet is promoted to anything that embeds Person.
func (p Person) Greet() string { return "Hi, I'm " + p.Name }

// Employee embeds Person (a field with only a type, no name).
type Employee struct {
	Person        // embedded: Employee.Name and Employee.Greet() now work
	Role   string // ordinary field
	Salary float64
}

func main() {
	fmt.Println("== constructor and pointer-receiver methods")
	acct, err := NewBankAccount("Priya", 100)
	if err != nil {
		fmt.Println("create failed:", err)
		return
	}
	_ = acct.Deposit(50)
	fmt.Printf("after deposit 50 : %.2f\n", acct.Balance())

	fmt.Println("\n== errors from Deposit/Withdraw")
	fmt.Println("deposit -5       :", acct.Deposit(-5))
	err = acct.Withdraw(500)
	fmt.Println("withdraw 500     :", err)
	fmt.Println("is insufficient? :", errors.Is(err, ErrInsufficientFunds)) // errors.Is sees through %w wrapping
	fmt.Println("withdraw 30 ok   :", acct.Withdraw(30), "balance", acct.Balance())
	_, err = NewBankAccount("Bad", -1)
	fmt.Println("negative opening :", err)

	fmt.Println("\n== PITFALL 1: value receiver cannot mutate")
	before := acct.Balance()
	acct.DepositBroken(1000)
	fmt.Printf("balance before %.2f, after DepositBroken(1000): %.2f (unchanged)\n", before, acct.Balance())

	fmt.Println("\n== PITFALL 2: literal skips validation")
	rogue := BankAccount{owner: "Rogue", balance: -9999} // compiles; no check ran
	fmt.Printf("rogue balance %.2f (use NewBankAccount to prevent this)\n", rogue.Balance())

	fmt.Println("\n== embedding")
	emp := Employee{
		Person: Person{Name: "Marta", Age: 34},
		Role:   "Engineer",
		Salary: 85000,
	}
	fmt.Println("emp.Name (promoted):", emp.Name, "| emp.Person.Name:", emp.Person.Name)
	fmt.Println("emp.Greet()        :", emp.Greet())
	fmt.Printf("full value         : %+v\n", emp) // %+v prints field names

	fmt.Println("\n== PITFALL 3: embedding is not inheritance")
	// var p Person = emp   // compile error: Employee is not a Person
	var p Person = emp.Person // you must pick the inner value explicitly
	fmt.Println("extracted Person:", p.Name)
}
