package main

import (
	"flag"
	"fmt"
	"kvs/client"
	"math"
	"os"
	"sort"
	"sync"
	"time"
)

func measureLatency(proxy *client.Proxy, n int) {
	var latencies []time.Duration
	for i := 0; i < n; i++ {
		now := time.Now()
		if err := proxy.RemoveSet(client.ChiaveSet(fmt.Sprint(i)), fmt.Sprint(i)); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		elapsed := time.Since(now)
		latencies = append(latencies, elapsed)
	}
	sort.Slice(latencies, func(i, j int) bool {
		return latencies[i] < latencies[j]
	})
	fmt.Println(latencies[int(.95*float64(len(latencies)))])
}

func measureThroughput(proxy *client.Proxy, n int, nSeconds int) {
	for i := 0; i < nSeconds; i++ {
		var wg sync.WaitGroup
		loopStart := time.Now()
		for i := 0; i < n; i++ {
			wg.Add(1)
			go func(i int) {
				if err := proxy.AddSet(client.ChiaveSet(fmt.Sprint(i)), fmt.Sprint(i)); err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
				wg.Done()
			}(i)
		}
		wg.Wait()
		t := time.Since(loopStart)
		<-time.After(time.Second - t)
		throughput := math.Ceil(float64(n) / time.Since(loopStart).Seconds())
		fmt.Printf("true throughput = %d\n", int(throughput))
	}
}

func main() {
	nops := flag.Int("nops", 10, "number of operations per second")
	nSeconds := flag.Int("nsec", 20, "number of seconds to run (for throughput)")
	mode := flag.String("mode", "t", "t/l")
	flag.Parse()
	proxy := client.NewProxy()
	defer proxy.Cleanup()
	switch *mode {
	case "l":
		measureLatency(proxy, *nops)
	case "t":
		measureThroughput(proxy, *nops, *nSeconds)
	default:
		fmt.Printf("unknown mode %s\n", *mode)
	}
}
