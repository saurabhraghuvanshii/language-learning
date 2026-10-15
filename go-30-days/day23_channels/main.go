/*
DAY 23: CHANNELS
RUN:  go run ./day23_channels

THEORY (plain words):
A channel is a typed pipe between goroutines: one side puts values in, the other
takes them out. Under the hood it is a small struct (hchan) holding a queue of
waiting values, a lock, and two lists of parked goroutines (blocked senders and
blocked receivers). Channels let goroutines COMMUNICATE instead of sharing memory.
  - make(chan T)      UNBUFFERED: capacity 0. A send and a receive must meet
    ("rendezvous"). The sender sleeps until a receiver arrives.
  - make(chan T, N)   BUFFERED: a queue of N slots. Send only blocks when the
    queue is full; receive only blocks when it is empty.
  - close(ch)         Says "no more values will ever be sent". Receivers still
    drain what is buffered, then get the zero value and ok=false.
  - for v := range ch Receives until the channel is closed AND drained.
  - chan<- T          send-only type, <-chan T receive-only type. The compiler
    enforces who may send or receive: documentation + safety.
  - v, ok := <-ch     comma-ok receive. ok==false means closed and empty.

REAL WORLD: job queues (HTTP handler enqueues emails, worker sends them), log
shippers, pipelines that read lines -> parse -> store, signalling "shut down now".
EXPLAIN-IT ANSWER: Why does an unbuffered send block until someone receives?
Because an unbuffered channel has NO storage. The value has nowhere to live, so the
only way to deliver it is to hand it directly to a receiver. The runtime parks the
sending goroutine on the channel's "waiting senders" list; when a receiver arrives
it copies the value straight from the sender and wakes the sender. That gives a
free guarantee: when the send returns, a receiver HAS the value (synchronisation).
PITFALLS:
 1. Sending on a closed channel panics. Only the SENDER side should close.
 2. Receiving from a nil channel or a channel nobody sends to blocks forever
    (deadlock: "all goroutines are asleep").
 3. Closing twice panics. Closing is NOT required for GC, only to tell receivers
    "that's all" (needed for range).
 4. A buffered channel hides ordering bugs: code "works" until the buffer fills.

EXERCISES:
 1. TODO: change the buffered channel capacity in demoBuffered to 1 and predict
    which send now blocks. Run it to check.
 2. TODO: add a second consumer goroutine and make both sum orders; combine the
    two partial sums at the end (hint: a results channel).
 3. TODO: try sending on the closed channel in demoClose (uncomment) and read the panic.
*/
package main

import (
	"fmt"
	"time"
)

// Order is a tiny invoice-like record.
type Order struct {
	ID    int
	Total float64 // in dollars; real code would use integer cents
}

// start is used to print elapsed time since program start in every log line.
var start = time.Now()

// logf prints a message with elapsed milliseconds so ordering is visible.
func logf(format string, args ...any) {
	fmt.Printf("[%4dms] %s\n", time.Since(start).Milliseconds(), fmt.Sprintf(format, args...))
}

// producer takes a SEND-ONLY channel (chan<- int): inside this function it is a
// compile error to receive from or range over out.
func producer(ids []int, out chan<- int) {
	for _, id := range ids {
		logf("producer: sending order id %d ...", id)
		out <- id // blocks (unbuffered) until the consumer receives
		logf("producer: order id %d was received", id)
	}
	close(out) // sender closes; tells the consumer's range loop to finish
}

// lookupOrder simulates a database lookup by id.
func lookupOrder(id int) Order {
	return Order{ID: id, Total: float64(id) * 12.50}
}

// consumer takes a RECEIVE-ONLY channel (<-chan int) and returns the total.
func consumer(in <-chan int) float64 {
	sum := 0.0
	for id := range in { // loop ends when in is closed and drained
		time.Sleep(50 * time.Millisecond) // pretend work; producer waits meanwhile
		o := lookupOrder(id)
		sum += o.Total
		logf("consumer: got order %d total %.2f (running %.2f)", o.ID, o.Total, sum)
	}
	return sum
}

// demoUnbuffered shows the producer blocking until the consumer takes each value.
func demoUnbuffered() {
	fmt.Println("--- unbuffered: send blocks until receive ---")
	ids := make(chan int) // capacity 0
	go producer([]int{101, 102, 103}, ids)
	total := consumer(ids) // runs in main goroutine
	logf("sum of order totals = %.2f", total)
}

// demoBuffered shows sends succeeding immediately while there is room.
func demoBuffered() {
	fmt.Println("--- buffered: send only blocks when full ---")
	ch := make(chan int, 2) // two free slots
	ch <- 1
	ch <- 2
	// ch <- 3 would block forever here (full, no receiver): deadlock panic.
	fmt.Println("len =", len(ch), "cap =", cap(ch)) // len = values waiting
	fmt.Println("received", <-ch, "then", <-ch)
}

// demoClose shows close, comma-ok and range on a buffered channel.
func demoClose() {
	fmt.Println("--- close and comma-ok ---")
	ch := make(chan string, 2)
	ch <- "paid"
	ch <- "shipped"
	close(ch) // buffered values are still deliverable
	// ch <- "x" // would panic: send on closed channel (TODO exercise 3)
	for i := 0; i < 3; i++ {
		v, ok := <-ch // ok is false once closed AND empty
		fmt.Printf("receive %d: value=%q ok=%v\n", i+1, v, ok)
	}
}

func main() {
	demoUnbuffered()
	demoBuffered()
	demoClose()
}
