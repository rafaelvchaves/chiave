package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"kvs/client"
	"log"
	"os"
	"sync"
	"time"

	"github.com/spenczar/tdigest"
)

const (
	counterKey client.ChiaveCounter = "key1"
	setKey     client.ChiaveSet     = "key2"
)

func main() {
	nops := flag.Int("nops", 10, "number of operations per second")
	flag.Parse()
	f, err := os.OpenFile("results.csv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("could not create file: %v", err)
	}
	defer f.Close()
	w := csv.NewWriter(f)
	td := tdigest.New()
	proxy := client.NewProxy()
	defer proxy.Cleanup()

	n := *nops
	var wg sync.WaitGroup
	loopStart := time.Now()
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			proxy.Increment(counterKey)
			wg.Done()
		}()
	}
	serial := 1000
	for i := 0; i < serial; i++ {
		start := time.Now()
		if err := proxy.Increment(counterKey); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		elapsed := time.Since(start)
		td.Add(float64(elapsed), 1)
	}
	wg.Wait()
	<-time.After(time.Second - time.Since(loopStart))
	throughput := int(float64(n+serial) / time.Since(loopStart).Seconds())
	fmt.Println("true throughput:", throughput)
	latency := time.Duration(td.Quantile(0.95))
	fmt.Println("95% latency:", latency)
	w.Write([]string{fmt.Sprint(throughput), fmt.Sprint(latency.Nanoseconds())})
	w.Flush()
}