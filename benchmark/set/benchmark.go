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
		if err := proxy.RemoveSet(client.ChiaveSet(fmt.Sprint(i%1000)), fmt.Sprint(i)); err != nil {
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

func measureThroughput(proxy *client.Proxy, n int, nSeconds int, stop chan bool, done chan bool) {
	var wg sync.WaitGroup
	stop <- false
	for !<-stop {
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
		stop <- false
	}
	done <- true
}

func main() {
	nops := flag.Int("nops", 10, "number of operations per second")
	nSeconds := flag.Int("nsec", 20, "number of seconds to run (for throughput)")
	mode := flag.String("mode", "t", "t/l")
	flag.Parse()
	proxy := client.NewProxy()
	defer proxy.Cleanup()
	var latencies []float64
	switch *mode {
	case "l":
		nSamples := 15
		for i := 0; i < nSamples; i++ {
			latency := measureLatency(proxy, 50)
			latencies = append(latencies, latency)
		}
		mu := mean(latencies)
		fmt.Println(int64(mu))
	case "t":
		// measureThroughput(proxy, *nops, *nSeconds)
	default:
		stop := make(chan bool, 2)
		// doneMap := sync.Map{}
		// doneMap.Store("done", false)
		done := make(chan bool, 1)
		go measureThroughput(proxy, *nops, *nSeconds, stop, done)
		nSamples := 5
		time.Sleep(2 * time.Second)
		p2 := client.NewProxy()
		defer p2.Cleanup()
		for i := 0; i < nSamples; i++ {
			latency := measureLatency(p2, 500)
			latencies = append(latencies, latency)
			fmt.Printf("latency: %f\n", latency)
		}
		mu := median(latencies)
		fmt.Println(latencies)
		fmt.Printf("%d,%d\n", int64(mu), *nops)
		// doneMap.Store("done", true)
		// fmt.Println(doneMap.Load("done"))
		stop <- true
		<-done
	}
}

func median(lst []float64) float64 {
	sort.Slice(lst, func(i, j int) bool {
		return lst[i] < lst[j]
	})
	return lst[int(.5*float64(len(lst)))]
}

func mean(lst []float64) float64 {
	sum := float64(0)
	for _, l := range lst {
		sum += l
	}
	return sum / float64(len(lst))
}
