/*
DAY 7: MINI PROJECT - IN-MEMORY TO-DO CLI
RUN:  go run ./day07_mini_todo_cli

THEORY:
This project reviews days 1-6. A struct groups related fields into one type
(Task: ID, Title, Done). A slice ([]Task) is a growable list backed by an array;
append may allocate a new array. Methods are functions with a receiver:
`func (l *TodoList) Add(...)`. A POINTER receiver (*TodoList) lets the method
modify the original list; a value receiver would modify a copy. Commands arrive
as strings from a scripted slice, are split with strings.Fields, and dispatched
with a switch. Failures are returned as errors, never hidden.
REAL WORLD: the same shape powers real CLIs and REST handlers: parse a command,
dispatch, mutate state, print a result or an error. Swap the slice for a
database and the script for os.Stdin and you have a real tool.
EXPLAIN-IT ANSWER: see the line-by-line comments below; read them aloud in
order, top to bottom. Each line says what it does and why it is needed.
PITFALLS: deleting from a slice while ranging over it by index; the slice
trick `append(s[:i], s[i+1:]...)` modifies the shared backing array; and a
value receiver silently changing only a copy (shown in brokenAdd).
EXERCISES:
1. TODO: add an "edit <id> <new title>" command.
2. TODO: add a "list-done" command that prints only completed tasks.
3. TODO: make Add reject duplicate titles by returning an error.
*/

// Declare that this file builds an executable.
package main

// Import the standard-library packages we use.
import (
	"errors"  // errors.New creates simple error values
	"fmt"     // fmt prints text and builds formatted errors
	"strconv" // strconv.Atoi converts "3" into the number 3
	"strings" // strings.Fields splits a command, strings.Join rebuilds titles
)

// Task is one to-do item: a struct groups its three fields together.
type Task struct {
	ID    int    // unique number the user types to refer to the task
	Title string // human-readable description
	Done  bool   // false (the zero value) until completed
}

// TodoList holds all tasks plus the counter that makes the next ID.
type TodoList struct {
	tasks  []Task // the slice of tasks, in the order they were added
	nextID int    // next ID to hand out; starts at 0, we bump it before use
}

// ErrNotFound is a sentinel error returned when an ID does not exist.
var ErrNotFound = errors.New("task not found")

// Add appends a new task and returns it. Pointer receiver so we change the real list.
func (l *TodoList) Add(title string) (Task, error) {
	// Trim spaces so "   " counts as empty.
	title = strings.TrimSpace(title)
	// Refuse empty titles by returning an error with a zero Task.
	if title == "" {
		return Task{}, errors.New("title must not be empty")
	}
	// Move to the next ID first so IDs start at 1.
	l.nextID++
	// Build the task; Done is left out so it takes its zero value, false.
	t := Task{ID: l.nextID, Title: title}
	// Append may allocate a bigger array, so we must store the result back.
	l.tasks = append(l.tasks, t)
	// Return the new task and a nil error to signal success.
	return t, nil
}

// find returns the slice index of the task with the given ID, or -1.
func (l *TodoList) find(id int) int {
	// Loop over every position; i is the index, we ignore the copy of the task.
	for i := range l.tasks {
		// Compare the stored ID with the wanted ID.
		if l.tasks[i].ID == id {
			return i // found: report where
		}
	}
	return -1 // not found sentinel
}

// Complete marks the task with this ID as done.
func (l *TodoList) Complete(id int) error {
	// Locate the task.
	i := l.find(id)
	// -1 means no such ID, so wrap the sentinel with details using %w.
	if i < 0 {
		return fmt.Errorf("complete %d: %w", id, ErrNotFound)
	}
	// Index into the slice (not a range copy) so the change sticks.
	l.tasks[i].Done = true
	return nil // success
}

// Delete removes the task with this ID, keeping the others in order.
func (l *TodoList) Delete(id int) error {
	// Locate the task.
	i := l.find(id)
	// Same not-found handling as Complete.
	if i < 0 {
		return fmt.Errorf("delete %d: %w", id, ErrNotFound)
	}
	// Join everything before i with everything after i, skipping element i.
	l.tasks = append(l.tasks[:i], l.tasks[i+1:]...)
	return nil // success
}

// List prints every task with a [x] or [ ] checkbox.
func (l *TodoList) List() {
	// Nothing to show? Say so and stop early.
	if len(l.tasks) == 0 {
		fmt.Println("  (no tasks)")
		return
	}
	// Visit each task; t is a COPY, which is fine because we only read it.
	for _, t := range l.tasks {
		// Start with an empty box and fill it when the task is done.
		box := "[ ]"
		if t.Done {
			box = "[x]"
		}
		// Print aligned: checkbox, ID, title.
		fmt.Printf("  %s %d. %s\n", box, t.ID, t.Title)
	}
}

// run executes one text command against the list and returns any error.
func run(l *TodoList, line string) error {
	// Split on whitespace: "add buy milk" becomes ["add","buy","milk"].
	parts := strings.Fields(line)
	// An empty line has no command at all.
	if len(parts) == 0 {
		return errors.New("empty command")
	}
	// The first word is the command, the rest are its arguments.
	cmd, args := parts[0], parts[1:]
	// Dispatch on the command word.
	switch cmd {
	case "add":
		// Re-join the remaining words into the title.
		t, err := l.Add(strings.Join(args, " "))
		if err != nil {
			return err // pass the failure up to main to print
		}
		fmt.Printf("  added #%d %q\n", t.ID, t.Title)
		return nil
	case "list":
		l.List() // printing is the whole job
		return nil
	case "done", "delete":
		// Both commands need exactly one numeric argument.
		if len(args) != 1 {
			return fmt.Errorf("%s needs exactly one id", cmd)
		}
		// Convert text to int; Atoi returns an error for non-numbers.
		id, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("bad id %q: %w", args[0], err)
		}
		if cmd == "done" {
			return l.Complete(id) // returns nil or an error
		}
		return l.Delete(id)
	}
	// No case matched: tell the user which word was unknown.
	return fmt.Errorf("unknown command %q", cmd)
}

// brokenAdd demonstrates the value-receiver PITFALL.
func (l TodoList) brokenAdd(title string) {
	// l is a COPY of the list, so this append never reaches the caller's list.
	l.tasks = append(l.tasks, Task{ID: 999, Title: title})
}

// main drives the program with a scripted list of commands so it terminates.
func main() {
	// Start with an empty list: zero value is ready to use (nil slice, nextID 0).
	var list TodoList
	// The scripted "user input"; includes deliberate mistakes to show error output.
	script := []string{
		"add buy milk",
		"add write Go code",
		"add call the bank",
		"list",
		"done 2",
		"done 42", // no such task
		"delete 1",
		"add   ",   // empty title
		"done abc", // not a number
		"dance",    // unknown command
		"list",
	}
	// Run each command in order.
	for _, line := range script {
		// Echo the command so the output reads like a session.
		fmt.Printf("> %s\n", line)
		// Execute it; an error is printed, not fatal, like a real CLI would do.
		if err := run(&list, line); err != nil {
			// errors.Is would let us branch on ErrNotFound; here we just report.
			fmt.Println("  error:", err)
		}
	}
	// Show the value-receiver pitfall: count must not change.
	before := len(list.tasks)
	list.brokenAdd("ghost task")
	fmt.Printf("[pitfall] value receiver: tasks before=%d after=%d (unchanged)\n", before, len(list.tasks))
}
