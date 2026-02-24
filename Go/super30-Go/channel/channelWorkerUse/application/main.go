package main

import "fmt"

func main() {
    go count()
}

func count() {
    for k := 0; k < 10; k++ {
        fmt.Println(k)
    }
}