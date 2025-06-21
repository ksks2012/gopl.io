// Example of a buffered channel in Go
package main

import (
	"fmt"
	"time"
)

func producer(ch chan int) {
	for i := 0; i < 5; i++ {
		ch <- i // Send data to the channel
		fmt.Printf("Produced: %d\n", i)
		time.Sleep(50 * time.Millisecond)
	}
	close(ch) // Close the channel after the producer is done
}

func consumer(ch chan int) {
	for num := range ch { // Use for-range loop to receive data from the channel until it is closed and empty
		fmt.Printf("Consumed: %d\n", num)
		time.Sleep(100 * time.Millisecond) // Consumer processes more slowly
	}
	fmt.Println("Consumer finished.")
}

func main() {
	// Create a buffered channel with size 3
	bufferedCh := make(chan int, 3)

	go producer(bufferedCh)
	go consumer(bufferedCh)

	// Main goroutine waits for a while to allow producer and consumer to finish
	time.Sleep(1 * time.Second)
	fmt.Println("Main goroutine ends.")
}
