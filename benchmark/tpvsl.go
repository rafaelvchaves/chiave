package main

import (
	// "encoding/csv"
	// "encoding/csv"
	"flag"
	"fmt"
	"kvs/client"
	"log"
	"os"
	"sync"
	"time"
	// "github.com/spenczar/tdigest"
)

const (
	counterKey client.ChiaveCounter = "key1"
	setKey     client.ChiaveSet     = "key2"
)

func main() {
	nops := flag.Int("nops", 10, "number of operations per second")
	flag.Parse()
	f, err := os.OpenFile("results2.csv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("could not create file: %v", err)
	}
	defer f.Close()
	// w := csv.NewWriter(f)
	// td := tdigest.New()
	proxy := client.NewProxy()
	defer proxy.Cleanup()
	// var mutex sync.Mutex

	n := *nops
	// latencies := make(chan float64, n)
	var wg sync.WaitGroup
	loopStart := time.Now()
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(i int) {
			// start := time.Now()
			if err := proxy.Increment(client.ChiaveCounter(fmt.Sprint(i % 10000))); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			// elapsed := time.Since(start)
			// fmt.Println(elapsed)
			// latencies <- float64(elapsed)
			// mutex.Lock()
			// td.Add(float64(elapsed), 1)
			// mutex.Unlock()
			wg.Done()
		}(i)
	}
	wg.Wait()
	// serial := 5
	// for i := 0; i < serial; i++ {
	// 	start := time.Now()
	// 	if err := proxy.Increment(counterKey); err != nil {
	// 		fmt.Println(err)
	// 		os.Exit(1)
	// 	}
	// 	elapsed := time.Since(start)
	// 	fmt.Println(elapsed)
	// 	// td.Add(float64(elapsed), 1)
	// }
	// fmt.Println(time.Since(loopStart))
	<-time.After(time.Second - time.Since(loopStart))
	t := time.Since(loopStart)
	fmt.Println(time.Duration(t.Nanoseconds() / int64(n)))
	throughput := int(float64(n) / t.Seconds())
	fmt.Println("true throughput:", throughput)
	// i := 0
	// for l := range latencies {
	// 	td.Add(l, 1)
	// 	i++
	// 	if i == n {
	// 		break
	// 	}
	// }
	// latency := time.Duration(td.Quantile(0.90))
	// fmt.Println("90% latency:", latency)
	// w.Write([]string{fmt.Sprint(throughput), fmt.Sprint(latency.Nanoseconds())})
	// w.Flush()
}
