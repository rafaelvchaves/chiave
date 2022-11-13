package benchmark_test

import (
	"fmt"
	"kvs/client"
	"os"
	"sync"
	"testing"
)

const (
	counterKey client.ChiaveCounter = "counter_key"
)

func BenchmarkIncrement(b *testing.B) {
	proxy := client.NewProxy()
	defer proxy.Cleanup()
	n := 100
	var wg sync.WaitGroup
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for j := 0; j < n; j++ {
			wg.Add(1)
			go func() {
				// time.Sleep(time.Second)
				if err := proxy.Increment(counterKey); err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
				wg.Done()
			}()
		}
	}
	wg.Wait()
}
