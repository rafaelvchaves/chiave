package main

import (
	"fmt"
	"kvs/client"
	"os"
)

const key client.ChiaveSet = "key1"

func main() {
	proxy := client.NewProxy()
	defer proxy.Cleanup()
	// for i := 0; i < 10000; i++ {
	// 	if err := proxy.AddSet(key, fmt.Sprint(i)); err != nil {
	// 		fmt.Println(err)
	// 		os.Exit(1)
	// 	}
	// }
	// for i := 0; i < 10000; i++ {
	// 	if err := proxy.RemoveSet(key, fmt.Sprint(i)); err != nil {
	// 		fmt.Println(err)
	// 		os.Exit(1)
	// 	}
	// }
	// time.Sleep(6 * time.Second)
	s, err := proxy.Get(key)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Printf("set: %s\n", s)
}
