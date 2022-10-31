package main

import (
	"fmt"
	"kvs/client"

	// "sync"
	"time"
)

func main() {
	chiaveProxy := client.NewProxy()
	defer chiaveProxy.Cleanup()

	start := time.Now()
	// var wg sync.WaitGroup
	// for i := 0; i < 100; i++ {
		// wg.Add(2)
		// chiaveProxy.Increment("key1")
		// chiaveProxy.Decrement("key1")
		// go func() { defer wg.Done(); chiaveProxy.Increment("key1") }()
		// go func() { defer wg.Done(); chiaveProxy.Decrement("key1") }()
	// }
	// chiaveProxy.Increment("key1")
	// wg.Wait()
	v, err := chiaveProxy.Get("key1")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(time.Since(start))
	fmt.Println(v)
}
