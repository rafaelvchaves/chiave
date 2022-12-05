package main

import (
	"flag"
	"fmt"
	"kvs/client"
	"os"
	"sort"
	"strconv"
	"sync"
	"time"

	"golang.org/x/exp/constraints"
)

func measureConvergenceTime(proxy *client.Proxy) time.Duration {
	setKey := client.ChiaveSet("set_key")
	var expected []string
	for i := 0; i < 20; i++ {
		element := strconv.Itoa(i)
		proxy.AddSet(setKey, element)
		expected = append(expected, element)
	}
	t, err := proxy.GetConvergenceTime("set_key", expected)
	if err != nil {
		fmt.Println(err)
		return 0
	}
	return t
}

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

func measureThroughput(proxy *client.Proxy, n int, stop chan bool, done chan bool) {
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
	mode := flag.String("mode", "t", "t/l")
	flag.Parse()
	proxy := client.NewProxy()
	defer proxy.Cleanup()
	switch *mode {
	case "c":
		var cts []time.Duration
		stop := make(chan bool, 2)
		done := make(chan bool, 1)
		go measureThroughput(proxy, *nops, stop, done)
		nSamples := 5
		time.Sleep(2 * time.Second)
		p2 := client.NewProxy()
		defer p2.Cleanup()
		for i := 0; i < nSamples; i++ {
			ct := measureConvergenceTime(p2)
			cts = append(cts, ct)
			fmt.Printf("convergence time: %v\n", ct)
		}
		mu := median(cts)
		fmt.Printf("%d,%d\n", int64(mu), *nops)
		stop <- true
		<-done
	default:
		var latencies []float64
		stop := make(chan bool, 2)
		done := make(chan bool, 1)
		go measureThroughput(proxy, *nops, stop, done)
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
		stop <- true
		<-done
	}
}

func median[T constraints.Ordered](lst []T) T {
	sort.Slice(lst, func(i, j int) bool {
		return lst[i] < lst[j]
	})
	return lst[int(.5*float64(len(lst)))]
}