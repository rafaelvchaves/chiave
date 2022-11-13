package main

import (
	"flag"
	"fmt"
	"kvs/client"
	"log"
	"os"
	"sort"
	"sync"
	"time"
)

func measureLatency(n int, proxy *client.Proxy) {
	var latencies []time.Duration
	// td := tdigest.New()
	for i := 0; i < n; i++ {
		now := time.Now()
		if err := proxy.Increment(client.ChiaveCounter("l" + fmt.Sprint(i))); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		elapsed := time.Since(now)
		latencies = append(latencies, elapsed)
	}
	sort.Slice(latencies, func(i, j int) bool {
		return latencies[i] < latencies[j]
	})
	// latency := time.Duration(td.Quantile(0.95))

	fmt.Println(latencies[int(.95*float64(len(latencies)))])
}

func main() {
	nops := flag.Int("nops", 10, "number of operations per second")
	mode := flag.String("mode", "t", "t/l")
	flag.Parse()
	f, err := os.OpenFile("results2.csv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("could not create file: %v", err)
	}
	defer f.Close()
	proxy := client.NewProxy()
	defer proxy.Cleanup()
	if *mode == "l" {
		measureLatency(*nops, proxy)
		return
	}
	// w := csv.NewWriter(f)
	// td := tdigest.New()
	n := *nops
	// latencies := make(chan float64, n)
	nSeconds := 40
	// latencyChan := make(chan time.Duration, 1)
	for i := 0; i < nSeconds; i++ {
		var wg sync.WaitGroup
		loopStart := time.Now()
		for i := 0; i < n; i++ {
			wg.Add(1)
			go func(i int) {
				if err := proxy.Increment(client.ChiaveCounter(fmt.Sprint(i % 10000))); err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
				wg.Done()
			}(i)
		}
		// wg.Add(1)
		// go measureLatency(serialOps, proxy, &wg, latencyChan)
		wg.Wait()
		t := time.Since(loopStart)
		<-time.After(time.Second - t)
		// latency := <-latencyChan
		throughput := int(float64(n) / time.Since(loopStart).Seconds())
		fmt.Printf("true throughput = %d\n", throughput)
		// fmt.Println("95% latency:", latency.String())
	}
	// fmt.Println("90% latency:", latency)
	// w.Write([]string{fmt.Sprint(throughput), fmt.Sprint(latency.Nanoseconds())})
	// w.Flush()
}
