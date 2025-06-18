// Practice 8.1: clockwall
package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"sync" // Need sync for concurrent map access if not using channel for display
	"time"
)

type clockInfo struct {
	name string
	addr string
	// Store the current time for this clock
	currentTime string
}

// updateMessage struct for passing updates to the display goroutine
type updateMessage struct {
	index int
	time  string
}

const (
	displayWidth = 20 // Each clock column width
	startRow     = 3  // Starting row for time display
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "usage: clockwall Name=host:port ...\n")
		os.Exit(1)
	}

	clocks := make([]*clockInfo, len(os.Args)-1)
	for i, arg := range os.Args[1:] {
		parts := strings.SplitN(arg, "=", 2)
		if len(parts) != 2 {
			log.Fatalf("invalid argument: %s", arg)
		}
		clocks[i] = &clockInfo{name: parts[0], addr: parts[1], currentTime: "N/A"}
	}

	// Channel to send time updates to the display goroutine
	updates := make(chan updateMessage)
	var wg sync.WaitGroup

	// Goroutine for each clock to connect and read time
	for i, clock := range clocks {
		wg.Add(1)
		go func(index int, c *clockInfo) {
			defer wg.Done()
			for {
				conn, err := net.Dial("tcp", c.addr)
				if err != nil {
					log.Printf("Failed to connect to %s (%s): %v", c.name, c.addr, err)
					updates <- updateMessage{index, "DISCONNECTED"}
					time.Sleep(3 * time.Second)
					continue
				}
				log.Printf("Connected to %s (%s)", c.name, c.addr)
				func() {
					defer conn.Close()
					reader := bufio.NewReader(conn)
					for {
						line, err := reader.ReadString('\n')
						if err != nil {
							if err != io.EOF {
								log.Printf("Failed to read from %s (%s): %v", c.name, c.addr, err)
							} else {
								log.Printf("Connection to %s (%s) closed by server.", c.name, c.addr)
							}
							break
						}
						updates <- updateMessage{index, strings.TrimSpace(line)} // Send update
					}
				}()
				time.Sleep(1 * time.Second)
			}
		}(i, clock)
	}

	// Goroutine to manage display updates
	go func() {
		// Initial setup: print headers
		fmt.Print("\033[2J\033[H")
		for _, clock := range clocks {
			fmt.Printf("%-*s", displayWidth, clock.name)
		}
		fmt.Println()
		for range clocks {
			fmt.Printf("%-*s", displayWidth, strings.Repeat("-", displayWidth-1))
		}
		fmt.Println()

		currentTimes := make([]string, len(clocks))
		for i := range currentTimes {
			currentTimes[i] = "N/A"
		}

		for update := range updates {
			currentTimes[update.index] = update.time
			fmt.Printf("\033[%d;%dH", startRow, 1)

			for _, t := range currentTimes {
				fmt.Printf("%-*s", displayWidth, t)
			}
			fmt.Print("\033[K")
		}
	}()

	wg.Wait() // You might use select {} or a signal handler to keep it running
}
