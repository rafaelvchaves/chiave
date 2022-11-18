package main

import (
	"flag"
	"fmt"
	"kvs/client"
	"os"
	"sort"
	"sync"
	"time"
)

func measureLatency(proxy *client.Proxy, n int) float64 {
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
	return float64(latencies[int(.95*float64(len(latencies)))].Nanoseconds())
}

func measureThroughput(proxy *client.Proxy, n int, nSeconds int) {
	for i := 0; i < nSeconds; i++ {
		var wg sync.WaitGroup
		loopStart := time.Now()
		for i := 0; i < n; i++ {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				if err := proxy.AddSet(client.ChiaveSet(fmt.Sprint(i%1000)), fmt.Sprint(i)); err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
			}(i)
		}
		wg.Wait()
		t := time.Since(loopStart)
		<-time.After(time.Second - t)
		// throughput := math.Ceil(float64(n) / time.Since(loopStart).Seconds())
		// fmt.Println(int(throughput))
	}
}

func main() {
	nops := flag.Int("nops", 10, "number of operations per second")
	nSeconds := flag.Int("nsec", 20, "number of seconds to run (for throughput)")
	mode := flag.String("mode", "t", "t/l")
	flag.Parse()
	proxy := client.NewProxy()
	var latencies []float64
	defer proxy.Cleanup()
	switch *mode {
	case "l":
		measureLatency(proxy, *nops)
	case "t":
		measureThroughput(proxy, *nops, *nSeconds)
	default:
		go measureThroughput(proxy, *nops, *nSeconds)
		nSamples := 30
		for i := 0; i < nSamples; i++ {
			latency := measureLatency(proxy, 200)
			latencies = append(latencies, latency)
			time.Sleep(time.Duration(*nSeconds / nSamples))
		}
		mu := mean(latencies)
		// variance := variance(latencies, mu)
		// width := 1.96 * math.Sqrt(variance/float64(len(latencies)))
		// lower := mu - width
		// upper := mu + width
		fmt.Printf("%d,%d\n", int64(mu), *nops)
	}
}

func mean(lst []float64) float64 {
	sum := float64(0)
	for _, l := range lst {
		sum += l
	}
	return sum / float64(len(lst))
}
