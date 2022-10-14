package replica

import (
	"fmt"
	"kvs/crdt"
	"time"
)

type Worker struct {
	addr       string
	workerID int
	kvs      Store
	requests chan ClientRequest
	events   chan crdt.Event
}

type Operation int

const (
	Increment Operation = iota
	Decrement
	Value
)

type ClientRequest struct {
	Key       string
	Operation Operation
	Params    any
}

func NewWorker(addr string, workerID int, kvs Store) Worker {
	return Worker{
		addr:       addr,
		workerID: workerID,
		kvs:      kvs,
		requests: make(chan ClientRequest),
		events:   make(chan crdt.Event),
	}
}

func (w *Worker) Start() {
	requestDeadline := 100 * time.Millisecond
	for {
		// events := make(map[string]crdt.Event)

		// phase 1: receive client requests and convert to events
		select {
		case req := <-w.requests:
			fmt.Println(req)
		// event := parse(req)
		// Events = // add to events map
		case <-time.After(requestDeadline):
			break
		}
		// for key, e := range events {
		// 	crdt, ok := w.kvs.Get(key)
		// 	if !ok {

		// 	}
		// 	// crdt.PersistEvents()
		// 	// w.broadcast(key, e)
		// }
	}
}

func (w *Worker) Act(r ClientRequest) {
	v, exists := w.kvs.Get(r.Key)
	switch r.Operation {
	case Increment:
		if !exists {
			// v = op.NewCounter()
		}
		counter, ok := v.(crdt.Counter)
		if !ok {

		}
		counter.Increment()
	}
}

func (w *Worker) Broadcast(key string, event crdt.Event) {
	// // map[key]ips
	// 	ips := w.localMap[key]
	// 	For ip := range ips {
	// 		If ip == LOCAL {
	// }
	// // send RPC to machine
	// }
}

func (w *Worker) AddRequest(r ClientRequest) {
	w.requests <- r
}

func (w *Worker) AddEvent(e crdt.Event) {
	w.events <- e
}
