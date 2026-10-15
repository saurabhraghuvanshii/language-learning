/*
DAY 1: HELLO, TOOLCHAIN
RUN:  go run ./day01_hello_toolchain

THEORY:
Go ships as one toolchain: `go run` compiles to a temp binary and runs it,
`go build` leaves a binary on disk, `go mod init <name>` creates go.mod which
names your module and pins the Go version and dependencies, `go vet` finds
suspicious code, and `gofmt` formats it. Every Go file starts with
`package <name>`. A package is a folder of files that are compiled together and
share names; other packages reach in through `import "path"`. A package named
`main` is special: the compiler turns it into an executable, and the runtime
starts by calling `func main()` inside it. Any other package name produces a
library that can only be imported. Go is compiled to a single static binary,
so deploying means copying one file.
REAL WORLD: every Go microservice, CLI (kubectl, docker) and tool starts at
package main / func main(); os.Args is how simple CLIs read their flags.
EXPLAIN-IT ANSWER: A package is a directory of .go files compiled as one unit
with a shared namespace; it is Go's unit of code reuse and visibility
(Capitalised names are exported). `main` must exist because the linker needs a
single entry point: package main says "build an executable", func main says
"start here". Without them you only have a library.
PITFALLS: (1) an unused import is a COMPILE ERROR in Go, not a warning; (2) an
unused local variable is also a compile error; (3) exported names need a
Capital first letter. Shown below in comments and by using os.Args[0] carefully.
EXERCISES:
1. TODO: run `go build -o hello ./day01_hello_toolchain` then run ./hello and
   compare os.Args[0] with `go run`'s temp path.
2. TODO: add a second function greet(name string) string and call it from main.
3. TODO: add `import "strings"` WITHOUT using it, read the compiler error,
   then use strings.ToUpper on the greeting to fix it.
*/

// Every file declares its package. "main" means: build an executable.
package main

// The import block lists packages this file uses from the standard library.
import (
	"fmt"     // fmt: formatted printing (Println, Printf, Sprintf)
	"os"      // os: operating system access, here the command-line arguments
	"runtime" // runtime: info about the Go toolchain and platform
)

// PITFALL (shown as a comment because it would not compile):
//
//	import "strings"   // imported but never used -> "imported and not used" error
//
//	func main() { x := 5 }   // x declared and not used -> compile error
//
// Go refuses to build these so codebases never accumulate dead imports.

// greeting is a package-level constant visible to every file in package main.
const greeting = "Hello, World!"

// helloWorld prints the classic first program output.
func helloWorld() {
	// Println writes its arguments followed by a newline to standard output.
	fmt.Println("[hello]", greeting)
}

// showArgs prints the program's own command line, as the OS passed it.
func showArgs() {
	// os.Args is a []string; element 0 is the path of the running program.
	fmt.Println("[args] program path (os.Args[0]):", os.Args[0] != "")
	// len gives the number of elements; with no extra flags it is 1.
	fmt.Println("[args] count:", len(os.Args))
	// range over a slice yields index and value; we only print the extra args.
	for i, a := range os.Args[1:] {
		// i is relative to the sliced view, so add 1 to match os.Args.
		fmt.Printf("[args] os.Args[%d] = %q\n", i+1, a)
	}
	// If the user passed nothing, say so, to make the output self-explaining.
	if len(os.Args) == 1 {
		fmt.Println("[args] no extra arguments; try: go run ./day01_hello_toolchain a b c")
	}
}

// showToolchain prints facts about the compiler that built this program.
func showToolchain() {
	// runtime.Version() is the Go version, e.g. "go1.24.0".
	fmt.Println("[toolchain] go version:", runtime.Version())
	// GOOS/GOARCH describe the target platform; cross-compiling just sets them.
	fmt.Println("[toolchain] platform:", runtime.GOOS+"/"+runtime.GOARCH)
}

// main is the entry point: the runtime calls it once after initialising packages.
func main() {
	helloWorld()    // print Hello World
	showArgs()      // print our own os.Args
	showToolchain() // print compiler info
}
