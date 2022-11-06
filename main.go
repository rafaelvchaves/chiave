package main

import (
	"fmt"
	"kvs/client"
	"sync"
	"time"
)

const (
	counterKey client.ChiaveCounter = "key1"
	setKey     client.ChiaveSet     = "key2"
)

func main2() {
	n := 10000
	var wg sync.WaitGroup
	proxy := client.NewProxy()
	defer proxy.Cleanup()
	loopStart := time.Now()
	var latencies []time.Duration
	var mutex sync.Mutex
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			mutex.Lock()
			start := time.Now()
			if err := proxy.Increment(counterKey); err != nil {
				fmt.Println(err)
			}
			time := time.Since(start)
			latencies = append(latencies, time)
			mutex.Unlock()
			wg.Done()
		}()
		// key := client.ChiaveCounter(fmt.Sprint(i))
		// latencies = append(latencies, time.Since(start))
	}
	wg.Wait()
	time.Sleep(1 * time.Second)
	fmt.Println(time.Since(loopStart))
	// i := 1
	avgLatency := mean(latencies)
	fmt.Printf("Average write latency: %f ms\n", avgLatency)
	// fmt.Println("latencies: ", latencies)
}

func main() {
	// var latencies []time.Duration
	main2()
	// d1 := &pb.DVV{
	// 	Clock:   make(map[string]int32),
	// 	Sibling: "REMOVE_A",
	// }
	// d2 := &pb.DVV{
	// 	Clock:   make(map[string]int32),
	// 	Sibling: "ADD_A",
	// }
	// d1 = util.Update(nil, []*pb.DVV{d1}, "0")
	// d2 = util.Update(nil, []*pb.DVV{d2}, "1")
	// fmt.Println(util.String([]*pb.DVV{d1, d2}...))
	// // fmt.Println(util.Join([]*pb.DVV{d1, d2}))
	// reduce := func([]string) string {
	// 	return "ADD_A"
	// }
	// fmt.Println(util.String(util.Reconcile(reduce, []*pb.DVV{d1, d2}, "0")))
	// fmt.Println(util.String(util.Reconcile(reduce, []*pb.DVV{d1, d2}, "1")))

	// n := 4
	// var totalLatency time.Duration
	// proxy := client.NewProxy()
	// defer proxy.Cleanup()
	// start := time.Now()
	// for i := 0; i < n; i++ {
	// 	start := time.Now()
	// 	if err := proxy.AddSet(setKey, fmt.Sprint(i)); err != nil {
	// 		fmt.Println(err)
	// 	}
	// 	totalLatency += time.Since(start)
	// }
	// time.Sleep(1 * time.Second)
	// fmt.Println(time.Since(start))
	// start = time.Now()
	// _, err := proxy.Get(setKey)
	// fmt.Println(time.Since(start))
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// fmt.Println(v)
	// avgLatency := float64(totalLatency.Milliseconds()) / float64(n)
	// fmt.Printf("Average write latency: %f ms\n", avgLatency)
	// <-done
}

func mean(times []time.Duration) float64 {
	mean := int64(0)
	for _, t := range times {
		mean += t.Milliseconds()
	}
	return float64(mean) / float64(len(times))
}
