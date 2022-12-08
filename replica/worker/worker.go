package worker

import (
	"kvs/crdt"
	"kvs/crdt/generator"
	pb "kvs/proto"
	"kvs/util"
	"time"
)

const (
	rpcTimeout = 60 * time.Second
)

type Worker[F crdt.Flavor] struct {
	replica     util.Replica
	generator   generator.Generator[F]
	kvs         Store[F]
	requests    chan LeaderRequest
	events      chan *pb.Event
	broadcaster *Broadcaster[F]
	logger      *util.Logger
	config      Config[F]
}

type LeaderRequest struct {
	Inner    *pb.Request
	Response chan Response
}

type Response struct {
	Context      *pb.Context
	CounterValue int64
	SetValue     []string
}

func New[F crdt.Flavor](replica util.Replica, generator generator.Generator[F], broadcaster *Broadcaster[F], config Config[F]) *Worker[F] {
	requestBufferSize := 1000000
	eventBufferSize := 1000000
	return &Worker[F]{
		generator:   generator,
		replica:     replica,
		kvs:         NewCache[F](),
		requests:    make(chan LeaderRequest, requestBufferSize),
		events:      make(chan *pb.Event, eventBufferSize),
		broadcaster: broadcaster,
		config:      config,
	}
}

func (w *Worker[F]) Start() {
	// set of keys modified in this epoch
	changeset := make(map[string]struct{})
	ticker := time.NewTicker(w.generator.BroadcastEpoch())
	for {
		select {
		case <-ticker.C:
			if len(w.requests) > 10 {
				newEpoch := len(w.requests) / 10
				ticker = time.NewTicker(time.Duration(newEpoch) * time.Millisecond)
			}
			for key := range changeset {
				v, ok := w.kvs.Get(key)
				if !ok {
					continue
				}
				e := v.PrepareEvent()
				e.Key = key
				w.broadcaster.Send(e)
			}
			changeset = make(map[string]struct{})
		case req := <-w.requests:
			w.process(req)
			if req.Inner.Operation != pb.OP_GETCOUNTER && req.Inner.Operation != pb.OP_GETSET {
				changeset[req.Inner.Key] = struct{}{}
			}
		case event := <-w.events:
			key := event.Key
			v := w.kvs.GetOrDefault(key, w.generator.New(event.Datatype, w.replica))
			v.PersistEvent(event)
			w.kvs.Put(key, v)
		}
	}
}

func (w *Worker[F]) process(req LeaderRequest) {
	r := req.Inner
	w.logRequestHandle(r.Key, r.Operation, r.Arg)
	switch r.Operation {
	case pb.OP_INCREMENT:
		v := w.kvs.GetOrDefault(r.Key, w.generator.New(pb.DT_Counter, w.replica))
		counter, ok := v.(crdt.Counter)
		if !ok {
			return
		}
		counter.Increment()
		w.kvs.Put(r.Key, v)
	case pb.OP_DECREMENT:
		v := w.kvs.GetOrDefault(r.Key, w.generator.New(pb.DT_Counter, w.replica))
		counter, ok := v.(crdt.Counter)
		if !ok {
			return
		}
		counter.Decrement()
		w.kvs.Put(r.Key, v)
	case pb.OP_GETCOUNTER:
		response := Response{}
		v, ok := w.kvs.Get(r.Key)
		if ok {
			if counter, ok := v.(crdt.Counter); ok {
				response.CounterValue = counter.Value()
			}
		}
		req.Response <- response
	case pb.OP_ADDSET:
		v := w.kvs.GetOrDefault(r.Key, w.generator.New(pb.DT_Set, w.replica))
		set, ok := v.(crdt.Set)
		if !ok {
			return
		}
		set.Add(r.Context, r.Arg)
		w.kvs.Put(r.Key, v)
		req.Response <- Response{
			Context: v.Context(),
		}
	case pb.OP_REMOVESET:
		v := w.kvs.GetOrDefault(r.Key, w.generator.New(pb.DT_Set, w.replica))
		set, ok := v.(crdt.Set)
		if !ok {
			return
		}
		set.Remove(r.Context, r.Arg)
		w.kvs.Put(r.Key, v)
		req.Response <- Response{
			Context: v.Context(),
		}
	case pb.OP_GETSET:
		response := Response{}
		v, ok := w.kvs.Get(r.Key)
		if ok {
			response.Context = v.Context()
			if set, ok := v.(crdt.Set); ok {
				response.SetValue = set.Value()
			}
		}
		req.Response <- response
	}
}

func (w *Worker[_]) logRequestHandle(key string, o pb.OP, arg string) {
	w.logger.Infof("worker %d handling %s(%s) on key %q", w.replica.WorkerID, o.String(), arg, key)
}

func (w *Worker[_]) PutRequest(r LeaderRequest) {
	w.requests <- r
}

func (w *Worker[F]) PutEvent(e *pb.Event) {
	w.events <- e
}
