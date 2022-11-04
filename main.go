package main

import (
	"fmt"
	"kvs/client"
)

const (
	counterKey client.ChiaveCounter = "key1"
	setKey     client.ChiaveSet     = "key2"
)

func main() {
	proxy := client.NewProxy()
	defer proxy.Cleanup()
	proxy.AddSet(setKey, "a")
	proxy.AddSet(setKey, "b")
	v, err := proxy.Get(setKey)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(v)
}
