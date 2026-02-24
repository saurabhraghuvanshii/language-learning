package main

import "time"

// Unbuffered channels are used to synchronize two goroutines
func main() {
    syncCh := make(chan bool)

    // launch a second goroutine
	go func() {
		longTask2()
        // finished
		syncCh <- true
	}()

	longTask()
    // blocks until the second goroutine has finished
	<-syncCh
}

func longTask() {
	time.Sleep(1 * time.Second)
}

func longTask2() {
	time.Sleep(3 * time.Second)
}