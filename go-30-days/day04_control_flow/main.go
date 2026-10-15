/*
DAY 4: CONTROL FLOW
RUN:  go run ./day04_control_flow

THEORY:
`if` needs no parentheses but always braces, and may start with an init
statement: `if v, err := f(); err != nil {}` where v and err exist only inside
the if/else chain. `switch` does not fall through: each case ends by itself
(use `fallthrough` to opt in), cases can be expressions, and a bare `switch {}`
replaces long if-else chains. Go has ONE loop keyword, `for`, in these forms:
three-part (init; cond; post), condition-only (like while), infinite (`for {}`)
and `for range` over slices, maps, strings, channels and (Go 1.22+) integers.
`break` leaves the loop, `continue` skips to the next iteration, and labels let
you break/continue an OUTER loop.
REAL WORLD: retry loops around HTTP calls, request routers using switch,
processing log lines with range, bounded polling.
EXPLAIN-IT ANSWER: Go's designers prefer fewer, orthogonal constructs. while,
do-while and for-each are all just `for` with different clauses, so there is one
thing to learn, read and tool (gofmt, vet), and no style wars over loop kinds.
PITFALLS: (1) switch cases do NOT fall through unlike C; (2) `break` inside a
switch inside a for only exits the switch, not the loop (use a label);
(3) a variable from an if init statement is scoped to that if/else only.
EXERCISES:
1. TODO: extend FizzBuzz so multiples of 7 print "Bazz" (add to the switch).
2. TODO: change guessing game seed (secret) and guesses; make it hit the cap.
3. TODO: write a labelled loop that finds the first pair (a,b), a*b==48, a<b.
*/

package main

import (
	"fmt"     // printing
	"strconv" // Atoi for the init-statement demo
)

// fizzBuzz prints 1..n replacing multiples of 3, 5 and both.
func fizzBuzz(n int) {
	// Classic three-part for: init; condition; post.
	for i := 1; i <= n; i++ {
		// A tagless switch: the first true case wins; no fallthrough.
		switch {
		case i%15 == 0: // check 15 first, otherwise 3 or 5 would match
			fmt.Print("FizzBuzz ")
		case i%3 == 0:
			fmt.Print("Fizz ")
		case i%5 == 0:
			fmt.Print("Buzz ")
		default:
			fmt.Print(i, " ")
		}
	}
	fmt.Println() // end the line
}

// guessGame plays a deterministic game: secret is fixed, guesses are scripted.
func guessGame(secret int, guesses []int) {
	// Condition-only for (Go's "while"): runs while attempts remain.
	attempts := 0
	for attempts < len(guesses) {
		g := guesses[attempts] // the scripted guess for this round
		attempts++             // count it
		switch {
		case g < secret:
			fmt.Printf("[guess %d] %d is too low\n", attempts, g)
		case g > secret:
			fmt.Printf("[guess %d] %d is too high\n", attempts, g)
		default:
			fmt.Printf("[guess %d] %d is correct! Won in %d tries\n", attempts, g, attempts)
			return // leave the function: game over
		}
	}
	fmt.Println("[game] out of scripted guesses, secret was", secret)
}

// ifInit shows an init statement scoped to the if.
func ifInit() {
	for _, in := range []string{"42", "abc"} {
		// n and err exist only within this if/else chain.
		if n, err := strconv.Atoi(in); err != nil {
			fmt.Printf("[if-init] %q is not a number\n", in)
		} else if n > 40 {
			fmt.Printf("[if-init] %q parsed to %d (big)\n", in, n)
		} else {
			fmt.Printf("[if-init] %q parsed to %d\n", in, n)
		}
		// Using n here would not compile: it is out of scope.
	}
}

// loops shows the other for forms plus break/continue/labels.
func loops() {
	// range over a slice: index and value.
	for i, host := range []string{"api", "db", "cache"} {
		fmt.Printf("[range] %d=%s\n", i, host)
	}
	// range over an integer (Go 1.22+): 0,1,2.
	for i := range 3 {
		fmt.Print("[range-int] ", i, "\n")
	}
	// Infinite loop with an explicit exit condition.
	retries := 0
	for {
		retries++
		if retries%2 == 1 {
			continue // skip odd attempts (pretend they failed fast)
		}
		fmt.Println("[infinite] succeeded on attempt", retries)
		break // exit the loop
	}
	// switch with explicit fallthrough, the opt-in behaviour.
	switch 1 {
	case 1:
		fmt.Println("[switch] case 1")
		fallthrough // runs the next case body regardless of its condition
	case 2:
		fmt.Println("[switch] case 2 (reached by fallthrough)")
	case 3:
		fmt.Println("[switch] case 3 (never)")
	}
	// PITFALL: plain break inside switch only exits the switch.
	for i := 0; i < 3; i++ {
		switch i {
		case 1:
			break // leaves the switch only, loop keeps going
		}
		fmt.Println("[pitfall] loop still running, i =", i)
	}
	// Fix: label the loop and break to it.
outer: // label names the statement that follows
	for a := 1; a < 10; a++ {
		for b := a + 1; b < 10; b++ {
			if a*b == 12 {
				fmt.Printf("[label] found %d*%d=12\n", a, b)
				break outer // exits BOTH loops
			}
		}
	}
}

// main runs everything.
func main() {
	fmt.Println("== FizzBuzz ==")
	fizzBuzz(15)
	fmt.Println("== Guessing game (secret 37) ==")
	guessGame(37, []int{50, 25, 40, 37, 99})
	fmt.Println("== if with init ==")
	ifInit()
	fmt.Println("== loops ==")
	loops()
}
