package main

import (
	"fmt"
	"kvs/crdt"
	"kvs/crdt/op"
)

func main() {
	m := map[string]crdt.CRDT{
		"key1": op.NewOCounter(),
	}

	val := m["key1"]
	c, ok := val.(crdt.Counter)
	if !ok {
		fmt.Println("key1 is not a counter value")
		return
	}
	fmt.Println(fmt.Sprintf("Counter value: %d", c.Value()))
	c.Increment()
	c.Increment()
	m["key1"] = val
	fmt.Println(fmt.Sprintf("Counter value: %d", m["key1"].(crdt.Counter).Value()))
}
