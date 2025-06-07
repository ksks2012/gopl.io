package main

import (
	"fmt"
	"sync"
	"time"
)

type ConcurrentMap struct {
	sync.RWMutex
	data map[string]int
}

func NewConcurrentMap() *ConcurrentMap {
	return &ConcurrentMap{
		data: make(map[string]int),
	}
}

func (cm *ConcurrentMap) Store(key string, value int) {
	cm.Lock()
	defer cm.Unlock()
	cm.data[key] = value
}

func (cm *ConcurrentMap) Load(key string) (int, bool) {
	cm.RLock()
	defer cm.RUnlock()
	val, ok := cm.data[key]
	return val, ok
}

func (cm *ConcurrentMap) Delete(key string) {
	cm.Lock()
	defer cm.Unlock()
	delete(cm.data, key)
}

func main() {
	cm := NewConcurrentMap()
	wg := sync.WaitGroup{}

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			cm.Store(fmt.Sprintf("key-%d", id), id)
			time.Sleep(time.Millisecond)
			val, ok := cm.Load(fmt.Sprintf("key-%d", id/2)) // read some data
			if ok {
				fmt.Printf("Goroutine %d loaded key-%d: %d\n", id, id/2, val)
			}
		}(i)
	}

	wg.Wait()
	fmt.Println("ConcurrentMap operations finished.")
	// fmt.Println("Final map size:", len(cm.data)) // This needs to be read with a lock
}
