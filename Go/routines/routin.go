package main

import (
	"fmt"
	"sync"
)

var count int = 0;


func main () {

	var wg sync.WaitGroup
	var mu sync.Mutex

	wg.Add(2);
	
	go func() {
		defer wg.Done();

		mu.Lock()
		count++;
		mu.Unlock()
	}()

	go func() {
		defer wg.Done();

		mu.Lock()
		count++;
		mu.Unlock()
	}()

	wg.Wait();
    fmt.Println(count);
}
