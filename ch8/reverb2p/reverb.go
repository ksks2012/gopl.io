// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 224.

// Reverb2 is a TCP server that simulates an echo.
// Practice 8.4: Modify the reverb2 program to close the write side of the connection
package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
	"time"
)

func echo(c net.Conn, shout string, delay time.Duration) {
	fmt.Fprintln(c, "\t", strings.ToUpper(shout))
	time.Sleep(delay)
	fmt.Fprintln(c, "\t", shout)
	time.Sleep(delay)
	fmt.Fprintln(c, "\t", strings.ToLower(shout))
}

func handleConn(c net.Conn) {
	var wg sync.WaitGroup
	input := bufio.NewScanner(c)

	for input.Scan() {
		wg.Add(1)
		go func(shout string) {
			defer wg.Done()
			echo(c, shout, 1*time.Second)
		}(input.Text())
	}
	wg.Wait()

	if tcpConn, ok := c.(*net.TCPConn); ok {
		log.Printf("Closing write half of connection for %s", c.RemoteAddr())
		tcpConn.CloseWrite()
	} else {
		log.Printf("Non-TCP connection detected for %s, closing full connection.", c.RemoteAddr())
		c.Close()
	}
}

func main() {
	l, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Reverb2 server listening on localhost:8000")

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Print(err) // e.g., connection aborted
			continue
		}
		log.Printf("Accepted connection from %s", conn.RemoteAddr())
		go handleConn(conn)
	}
}
