package main

import (
	"fmt"
	"time"
)

// using sync.WaitGroup
// func main() {
// 	var wg sync.WaitGroup
// 	wg.Add(1)
// 	go func() {
// 		count("hello")
// 		wg.Done()
// 	}()
// 	wg.Wait()
// }

func main() {
	c := make(chan string)
	go countCh("hello", c)
	for msg := range c {
		fmt.Println(msg)
	}
}

func countCh(thing string, c chan string) {
	defer close(c)
	for range 5 {
		c <- thing
		time.Sleep(time.Millisecond * 500)
	}
}

func count(thing string) {
	for range 5 {
		fmt.Println(thing)
		time.Sleep(time.Millisecond * 500)
	}
}
