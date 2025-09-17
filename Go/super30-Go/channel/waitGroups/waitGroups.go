package main

import (
	"fmt"
	"sync"
	"time"
)


func main() {
	fmt.Println("program start")

	var waitGroup sync.WaitGroup
	waitGroup.Add(10)
	for i := 0; i < 10; i++ {
		go concurrentTask(i, &waitGroup)
	}
	waitGroup.Wait()

	finishTask()
	fmt.Println("Program end")
}


func finishTask() {
	fmt.Println("Execcuting finish task")
}

func concurrentTask(taskNumber int, waitGroup *sync.WaitGroup) {
	fmt.Println("Begin Execute task number ", taskNumber)
	time.Sleep(100 * time.Millisecond)
	fmt.Println("end execute tak ", taskNumber)
	waitGroup.Done()
}