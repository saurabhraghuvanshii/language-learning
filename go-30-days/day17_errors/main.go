/*
DAY 17: ERRORS
RUN:  go run ./day17_errors

THEORY:
In Go an error is just a value of the built-in interface
type error interface { Error() string }. Functions return it as the last
result and callers check "if err != nil". There are no exceptions.
errors.New makes a simple error. fmt.Errorf with the %w verb WRAPS an error:
the new error adds context ("loading config: ...") and remembers the original,
reachable through Unwrap. errors.Is(err, target) walks the chain looking for
a specific value (a sentinel such as ErrMissingKey). errors.As(err, &target)
walks the chain looking for a specific TYPE (like *ParseError) and copies it
into target so you can read its fields.
Sentinel errors are package-level variables (var ErrX = errors.New(...)) that
callers compare with errors.Is. Custom error types carry structured data
(line number, key).
REAL WORLD: os.Open returns *fs.PathError wrapping syscall errors; sql.ErrNoRows
is a sentinel; HTTP handlers use errors.Is/As to choose a 404 vs a 500.
EXPLAIN-IT ANSWER: When to wrap an error vs return it as is?
Wrap (with %w) when you add useful context the caller does not have, such as
which file, which line, which operation, AND you are willing to make the
underlying error part of your API (callers may errors.Is it). Return as is
when you have nothing to add (the error already says it all) or when you do
not want to expose the internal cause (use %v instead of %w to hide it).
Do not wrap at every layer with the same words: "failed: failed: failed:".
PITFALLS:
  - Comparing with == breaks once an error is wrapped; use errors.Is.
  - Using %v instead of %w silently drops the chain.
  - errors.As needs a pointer to the target variable.
  - Ignoring the error returned by a function.

EXERCISES:
 1. Add a new sentinel ErrDuplicateKey and make Load fail when a key repeats.
 2. Add a "port" validation: it must be an int 1-65535; return a wrapped
    custom error carrying the bad value.
 3. In main, use errors.Join to combine two errors and test one with Is.
*/
package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// Sentinel errors: fixed values callers can test with errors.Is.
var (
	ErrMissingKey = errors.New("missing required key")
	ErrBadLine    = errors.New("malformed line")
)

// ParseError is a custom error type carrying structured data.
type ParseError struct {
	Line int    // 1-based line number in the input
	Text string // the offending text
	Err  error  // the underlying cause (a sentinel or strconv error)
}

// Error builds the message. Pointer receiver, so *ParseError is the error.
func (e *ParseError) Error() string {
	return fmt.Sprintf("line %d (%q): %v", e.Line, e.Text, e.Err)
}

// Unwrap lets errors.Is / errors.As look inside ParseError at its cause.
func (e *ParseError) Unwrap() error { return e.Err }

// Config is the result of loading.
type Config struct {
	Host    string
	Port    int
	Verbose bool
}

// Load parses "key=value" lines. Blank lines and lines starting with # are skipped.
func Load(text string) (*Config, error) {
	values := map[string]string{}
	for i, raw := range strings.Split(text, "\n") {
		line := strings.TrimSpace(raw)
		if line == "" || strings.HasPrefix(line, "#") {
			continue // nothing to parse
		}
		key, val, found := strings.Cut(line, "=") // splits at the first "="
		if !found || strings.TrimSpace(key) == "" {
			// Return a custom error type that wraps the sentinel.
			return nil, &ParseError{Line: i + 1, Text: line, Err: ErrBadLine}
		}
		values[strings.TrimSpace(key)] = strings.TrimSpace(val)
	}

	cfg := &Config{}
	// Required key: host.
	host, ok := values["host"]
	if !ok {
		// %w wraps the sentinel and adds the name of the key as context.
		return nil, fmt.Errorf("key %q: %w", "host", ErrMissingKey)
	}
	cfg.Host = host

	// Required key: port, must be an integer.
	portText, ok := values["port"]
	if !ok {
		return nil, fmt.Errorf("key %q: %w", "port", ErrMissingKey)
	}
	port, err := strconv.Atoi(portText)
	if err != nil {
		// err is *strconv.NumError; we wrap it so errors.As can still find it.
		return nil, fmt.Errorf("key %q: %w", "port", err)
	}
	cfg.Port = port

	// Optional key: verbose.
	if v, ok := values["verbose"]; ok {
		b, err := strconv.ParseBool(v)
		if err != nil {
			return nil, fmt.Errorf("key %q: %w", "verbose", err)
		}
		cfg.Verbose = b
	}
	return cfg, nil
}

// LoadNamed adds the source name as context. This is the "wrap on the way up" step.
func LoadNamed(name, text string) (*Config, error) {
	cfg, err := Load(text)
	if err != nil {
		return nil, fmt.Errorf("loading %s: %w", name, err)
	}
	return cfg, nil
}

// classify inspects an error using Is and As, as a caller would.
func classify(err error) string {
	var pe *ParseError
	var ne *strconv.NumError
	switch {
	case err == nil:
		return "ok"
	case errors.Is(err, ErrMissingKey):
		return "missing key (sentinel match)"
	case errors.As(err, &pe): // pe now points at the ParseError in the chain
		return fmt.Sprintf("parse error at line %d, is bad line: %t", pe.Line, errors.Is(err, ErrBadLine))
	case errors.As(err, &ne):
		return fmt.Sprintf("number problem: func=%s input=%q", ne.Func, ne.Num)
	default:
		return "unknown error"
	}
}

func main() {
	inputs := []struct{ name, text string }{
		{"good.conf", "# server\nhost = example.com\nport = 8080\nverbose = true\n"},
		{"nohost.conf", "port=80\n"},
		{"badline.conf", "host=a.com\nthis line is wrong\nport=1\n"},
		{"badport.conf", "host=a.com\nport=eighty\n"},
	}
	for _, in := range inputs {
		fmt.Println("==", in.name)
		cfg, err := LoadNamed(in.name, in.text)
		if err != nil {
			fmt.Println("error  :", err)
			fmt.Println("class  :", classify(err))
			fmt.Println("== compare (err == ErrMissingKey):", err == ErrMissingKey) // false once wrapped
			continue
		}
		fmt.Printf("config : %+v\n", *cfg)
	}

	// Hiding the cause with %v: Is can no longer see the sentinel.
	hidden := fmt.Errorf("opaque: %v", ErrMissingKey)
	fmt.Println("== verb v vs verb w")
	fmt.Println("hidden Is sentinel :", errors.Is(hidden, ErrMissingKey))
	fmt.Println("wrapped Is sentinel:", errors.Is(fmt.Errorf("ctx: %w", ErrMissingKey), ErrMissingKey))
}
