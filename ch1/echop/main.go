// Practice: Echo program
package main

import (
	"fmt"
	"os"
	"strings"
	"time"
)

// !+
func main() {
	// P1
	fmt.Println(os.Args[0])

	// P2
	for i, arg := range os.Args[1:] {
		fmt.Printf("Argument %d: %s\n", i+1, arg)
	}

	// P3
	var loopNumber = 10000000
	if len(os.Args) < 2 {
		os.Args = append(os.Args, "arg1", "arg2", "arg3", "arg4", "arg5")
	}

	// Method 1: String concatenation using +=
	// Each loop starts concatenation from scratch to avoid exponential string growth
	startTime := time.Now()
	for j := 0; j < loopNumber; j++ {
		var s_concat, sep string
		for i := 1; i < len(os.Args); i++ {
			s_concat += sep + os.Args[i]
			sep = " "
		}
	}
	endTime := time.Now()
	duration := endTime.Sub(startTime)
	fmt.Printf("Concatenation took %s\n", duration)

	// Method 2: Using strings.Join
	// Each loop performs a complete Join operation
	startTime = time.Now()
	for j := 0; j < loopNumber; j++ {
		_ = strings.Join(os.Args[1:], " ")
	}
	endTime = time.Now()
	duration = endTime.Sub(startTime)
	fmt.Printf("Join took %s\n", duration)

	// Method 3 (Recommended): Use strings.Builder (more efficient for large amounts of string concatenation)
	// Each loop starts building from scratch to avoid exponential string growth
	startTime = time.Now()
	for j := 0; j < loopNumber; j++ {
		var builder strings.Builder
		for i := 1; i < len(os.Args); i++ {
			if i > 1 {
				builder.WriteString(" ")
			}
			builder.WriteString(os.Args[i])
		}
		_ = builder.String()
	}
	endTime = time.Now()
	duration = endTime.Sub(startTime)
	fmt.Printf("Builder took %s\n", duration)
}

//!-
