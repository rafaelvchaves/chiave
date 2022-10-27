package worker

import (
	"kvs/crdt"
	"kvs/crdt/generator"
	"kvs/util"
	"time"
)

type Worker[F crdt.Flavor] struct {
	replica   util.Replica
	generator generator.Generator[F]
	kvs       Store[F]
	requests  chan ClientRequest
	events    chan crdt.Event
}

type Operation int

const (
	Increment Operation = iota
	Decrement
	Get
)

type ClientRequest struct {
	Key       string
	Operation Operation
	Params    any
	Response  chan Response
}

type Response = struct {
	Value  string
	Exists bool
}

func New[F crdt.Flavor](replica util.Replica, kvs Store[F], generator generator.Generator[F]) Worker[F] {
	return Worker[F]{
		generator: generator,
		replica:   replica,
		kvs:       kvs,
		requests:  make(chan ClientRequest),
		events:    make(chan crdt.Event),
	}
}

func (w *Worker[F]) Get(key string) (string, bool) {
	v, ok := w.kvs.Get(key)
	if ok {
		return v.String(), true
	}
	return "", false
}

func (w *Worker[F]) Start() {
	requestDeadline := 100 * time.Millisecond
	for {
		// set of keys modified in this epoch
		var changeset util.Set[string]
	reqLoop:
		for {
			// phase 1: receive client requests and convert to events
			select {
			case req := <-w.requests:
				w.process(req)
				changeset.Add(req.Key)
			case <-time.After(requestDeadline):
				break reqLoop
			}
		}
		// phase 2: go through all affected keys and broadcast to other owners
		changeset.Range(func(key string) bool {
			v, ok := w.kvs.Get(key)
			if !ok {
				return true
			}
			e := v.GetEvent()
			e.Key = key
			w.broadcast(e)
			return true
		})

		// phase 3: drain event queue and persist all events
		for event := range w.events {
			v := w.kvs.GetOrDefault(event.Key, w.generator.New(event.Type, w.replica))
			v.PersistEvent(event)
		}
	}
}

func (w *Worker[F]) process(r ClientRequest) {
	switch r.Operation {
	case Increment:
		v := w.kvs.GetOrDefault(r.Key, w.generator.New(crdt.CType, w.replica))
		counter, ok := v.(crdt.Counter)
		if !ok {
			return
		}
		counter.Increment()
		// w.kvs.Put(r.Key, v)
	case Decrement:
		v := w.kvs.GetOrDefault(r.Key, w.generator.New(crdt.CType, w.replica))
		counter, ok := v.(crdt.Counter)
		if !ok {
			return
		}
		counter.Decrement()
		// w.kvs.Put(r.Key, v)
	case Get:
		v, ok := w.kvs.Get(r.Key)
		r.Response <- Response{
			Value:  v.String(),
			Exists: ok,
		}
	}
}

func (w *Worker[F]) broadcast(event crdt.Event) {
	// TODO:
	// (1): hash event.Key to get address of leader(s)
	// (2): for each leader, send an RPC (ProcessEvent?)
	// (3): on leader side, processEvent implementation should simply
	// invoke the PutEvent() method on the proper worker.

	// should broadcast work differently depending on flavor?
}

func (w *Worker[_]) PutRequest(r ClientRequest) {
	w.requests <- r
}

func (w *Worker[F]) PutEvent(e crdt.Event) {
	w.events <- e
}
