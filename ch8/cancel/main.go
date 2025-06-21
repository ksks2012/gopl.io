// Example of using context.Context to cancel a goroutine
package main

import (
	"context"
	"fmt"
	"time"
)

// download simulates a file download process.
// It takes a context.Context to listen for cancellation signals.
func download(ctx context.Context, fileName string) {
	fmt.Printf("Starting download of %s...\n", fileName)

	// Simulate multiple steps in the download process
	for i := 1; i <= 5; i++ {
		select {
		case <-time.After(1 * time.Second): // Simulate downloading a part every second
			fmt.Printf("Downloading %s: completed %d/5 parts...\n", fileName, i)
		case <-ctx.Done(): // Listen for cancellation signal
			fmt.Printf("Download of %s was cancelled: %v\n", fileName, ctx.Err())
			return // Exit goroutine immediately upon cancellation
		}
	}

	fmt.Printf("Download of %s completed!\n", fileName)
}

func main() {
	fmt.Println("Main program starting download task.")

	// 1. Create a cancellable Context
	// context.WithCancel returns a new Context and a cancel function
	ctx, cancel := context.WithCancel(context.Background())

	// 2. Start a goroutine to perform the download task, passing the Context to it
	go download(ctx, "big_file.zip")

	// 3. Main program waits for a while, simulating other operations
	time.Sleep(3 * time.Second)

	// 4. At some point (e.g., user clicks cancel button), call the cancel function
	fmt.Println("\nMain program: preparing to cancel download task...")
	cancel() // Calling cancel closes the ctx.Done() channel

	// 5. Wait for the goroutine to respond to the cancellation signal and exit
	// To observe the cancellation effect, ensure the main program doesn't exit immediately
	time.Sleep(2 * time.Second) // Give the goroutine some time to respond to cancellation
	fmt.Println("Main program ends.")
}
