package main

import (
	"kvs/client"
)

func main() {
	chiaveProxy := client.NewProxy(3, 5)
	chiaveProxy.Increment("key1")
	chiaveProxy.Increment("key1")
	chiaveProxy.Increment("key1")
}
