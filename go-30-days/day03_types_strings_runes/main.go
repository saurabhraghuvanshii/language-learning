/*
DAY 3: TYPES, STRINGS, BYTES AND RUNES
RUN:  go run ./day03_types_strings_runes

THEORY:
Go has sized ints (int8..int64, uint8..uint64), plain int/uint (platform size),
float32/64, bool, string, and byte (alias of uint8) and rune (alias of int32).
There are NO implicit conversions: you write float64(n). A string is an
immutable sequence of BYTES, normally holding UTF-8 text. UTF-8 encodes each
Unicode code point (a "rune") in 1 to 4 bytes: ASCII is 1 byte, Devanagari is 3.
len(s) counts bytes. s[i] is a byte. `for i, r := range s` decodes UTF-8 and
yields the byte offset and the rune. []rune(s) copies into code points.
Even runes are not always "letters": a visible character can be several runes
(base letter + combining mark).
REAL WORLD: validating usernames, truncating text for SMS or tweets, parsing
protocols, storing names from any language in a database without corruption.
EXPLAIN-IT ANSWER: len("नमस्ते") returns the number of BYTES (18), because the
string is UTF-8 and each of its 6 code points takes 3 bytes. Those 6 runes in
turn form only 3 user-visible clusters (न, म, स्ते), since the virama and
vowel sign are combining marks. len measures storage, not letters.
PITFALLS: (1) len() on non-ASCII text; (2) indexing s[i] gives a byte, slicing
in the middle of a rune produces invalid UTF-8; (3) strings are immutable so
s[0] = 'x' does not compile; (4) converting a number to string with string(65)
gives "A", not "65" (use strconv.Itoa).
EXERCISES:
1. TODO: add the string "héllo wörld" and print bytes vs runes.
2. TODO: write reverseRunes(s string) string that reverses by rune; test on Hindi.
3. TODO: use unicode.IsUpper to count capital letters in "Go Is Fun In India".
*/

package main

import (
	"fmt"          // printing
	"strconv"      // number <-> string conversions
	"unicode/utf8" // utf8.RuneCountInString etc.
)

// analyse prints bytes versus runes for one string.
func analyse(label, s string) {
	// len(s) = number of bytes stored.
	bytes := len(s)
	// RuneCountInString decodes UTF-8 and counts code points.
	runes := utf8.RuneCountInString(s)
	fmt.Printf("[%s] %q bytes=%d runes=%d\n", label, s, bytes, runes)
}

// walk iterates with range to show offsets and code points.
func walk(s string) {
	// range over a string yields (byte offset, rune); offsets skip by rune width.
	for i, r := range s {
		// %c prints the character, %U the Unicode code point, utf8.RuneLen the byte width.
		fmt.Printf("  offset=%2d rune=%c %U width=%d\n", i, r, r, utf8.RuneLen(r))
	}
}

// main runs the demonstrations.
func main() {
	fmt.Println("== Bytes vs runes ==")
	english := "Hello"
	hindi := "नमस्ते" // 6 code points, 18 bytes
	analyse("english", english)
	analyse("hindi", hindi)

	fmt.Println("== range over Hindi (decodes runes) ==")
	walk(hindi)

	fmt.Println("== Plain index gives bytes ==")
	// hindi[0] is just the first byte of the 3-byte encoding of न.
	fmt.Printf("[index] hindi[0]=%d (0x%x) is a byte, not a letter\n", hindi[0], hindi[0])
	// Convert to []rune to index by code point.
	rs := []rune(hindi)
	fmt.Printf("[index] []rune(hindi)[0]=%c len=%d\n", rs[0], len(rs))
	// PITFALL: slicing bytes mid-rune creates invalid UTF-8.
	broken := hindi[:2]
	fmt.Printf("[pitfall] hindi[:2] valid UTF-8? %t\n", utf8.ValidString(broken))
	// Safe truncation: slice the rune slice, then convert back.
	fmt.Println("[safe] first 3 runes:", string(rs[:3]))

	fmt.Println("== Type conversions ==")
	var count int = 7
	var price float64 = 19.99
	// No implicit mixing: convert one side explicitly.
	total := float64(count) * price
	fmt.Printf("[convert] %d x %.2f = %.2f\n", count, price, total)
	// Float -> int truncates toward zero.
	fmt.Println("[convert] int(19.99) =", int(price))
	// PITFALL: string(int) is a rune conversion, not number formatting.
	fmt.Printf("[pitfall] string(rune(65)) = %q, strconv.Itoa(65) = %q\n", string(rune(65)), strconv.Itoa(65))
	// Parsing returns an error for bad input instead of crashing.
	n, err := strconv.Atoi("42x")
	fmt.Println("[parse] Atoi(\"42x\") =", n, "err:", err)
	// byte and rune are just integer aliases.
	var b byte = 'A'
	var r rune = 'अ'
	fmt.Println("[alias] byte 'A' =", b, " rune 'अ' =", r)
	// Strings are immutable: build changes via []byte or []rune copies.
	bs := []byte("go")
	bs[0] = 'G'
	fmt.Println("[immutable] modified copy:", string(bs))
}
