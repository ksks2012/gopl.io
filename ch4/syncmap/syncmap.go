package main

import (
	"fmt"
	"sync"
)

func main() {
	var sm sync.Map // Zero value is usable, no need to make()

	wg := sync.WaitGroup{}

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			key := fmt.Sprintf("user:%d", id)
			val := id * 10
			sm.Store(key, val) // Store

			loadedVal, ok := sm.Load(key) // Load
			if ok {
				fmt.Printf("Goroutine %d: loaded %s = %v\n", id, key, loadedVal)
			}
		}(i)
	}

	wg.Wait()

	// Iterate over sync.Map
	sm.Range(func(key, value interface{}) bool {
		fmt.Printf("Key: %v, Value: %v\n", key, value)
		return true // Return true to continue iteration
	})
}
