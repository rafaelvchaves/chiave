package main

import (
	"fmt"
	"kvs/client"
)

func main() {
	chiaveProxy := client.NewProxy(2)
	defer chiaveProxy.Cleanup()
	chiaveProxy.Increment("key1")
	chiaveProxy.Increment("key1")
	chiaveProxy.Increment("key3")
	v, err := chiaveProxy.Get("key1")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(v)
}
