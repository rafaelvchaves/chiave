package worker

import (
	"context"
	"fmt"
	"kvs/crdt"
	"kvs/crdt/generator"
	pb "kvs/proto"
	"kvs/util"
	"os"
	"time"

	"github.com/buraksezer/consistent"
	"google.golang.org/grpc"
)

const (
	rpcTimeout   = 5 * time.Second
	requestEpoch = 100 * time.Millisecond
	eventEpoch   = 100 * time.Millisecond
)

type Worker[F crdt.Flavor] struct {
	replica     util.Replica
	generator   generator.Generator[F]
	kvs         Store[F]
	requests    chan LeaderRequest
	events      chan *pb.Event
	hashRing    *consistent.Consistent
	connections map[string]*grpc.ClientConn
	cfg         util.Config
	logger      *util.Logger
	workers     []*Worker[F]
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

func New[F crdt.Flavor](replica util.Replica, generator generator.Generator[F], logger *util.Logger, workers []*Worker[F]) *Worker[F] {
	cfg := util.LoadConfig()
	requestBufferSize := 100000
	eventBufferSize := 100000
	return &Worker[F]{
		generator:   generator,
		replica:     replica,
		kvs:         NewCache[F](),
		requests:    make(chan LeaderRequest, requestBufferSize),
		events:      make(chan *pb.Event, eventBufferSize),
		hashRing:    util.GetHashRing(),
		connections: util.GetConnections(),
		cfg:         cfg,
		logger:      logger,
		workers:     workers,
	}
}

func (w *Worker[F]) Start() {
	for {
		// set of keys modified in this epoch
		changeset := make(map[string]struct{})
		// phase 1: receive client requests and convert to events
		requestsProcessed := 0
		eventsProcessed := 0
		wid := 3
	reqLoop:
		for timeout := time.After(requestEpoch); ; {
			select {
			case <-timeout:
				break reqLoop
			case req := <-w.requests:
				requestsProcessed++
				if w.replica.WorkerID == wid && len(w.requests) > 0 {
					fmt.Printf("request buffer size: %d\n", len(w.requests))
				}
				w.process(req)
				if req.Inner.Operation != pb.OP_GETCOUNTER && req.Inner.Operation != pb.OP_GETSET {
					changeset[req.Inner.Key] = struct{}{}
				}
			}
		}

		// phase 2: go through all affected keys and broadcast to other replicas
		for key := range changeset {
			v, ok := w.kvs.Get(key)
			if !ok {
				continue
			}
			e := v.PrepareEvent()
			e.Key = key
			w.broadcast(e)
		}
		if w.replica.WorkerID == wid && len(changeset) > 0 {
			fmt.Printf("%d events sent\n", len(changeset)*(w.cfg.RepFactor-1))
		}

		// phase 3: receive events from other replicas
	eventLoop:
		for timeout := time.After(eventEpoch); ; {
			select {
			case <-timeout:
				break eventLoop
			case event := <-w.events:
				eventsProcessed++
				// if w.replica.WorkerID == wid {
				// 	fmt.Printf("event buffer size: %d\n", len(w.events))
				// }
				v := w.kvs.GetOrDefault(event.Key, w.generator.New(event.Datatype, w.replica))
				v.PersistEvent(event)
				w.kvs.Put(event.Key, v)
			default:
				break eventLoop
			}
		}
		if requestsProcessed > 0 && w.replica.WorkerID == wid {
			fmt.Printf("worker %d requests processed: %d\n", w.replica.WorkerID, requestsProcessed)
		}
		if eventsProcessed > 0 && w.replica.WorkerID == wid {
			fmt.Printf("worker %d events processed: %d\n", w.replica.WorkerID, eventsProcessed)
		}
	}
}

func (w *Worker[F]) process(req LeaderRequest) {
	r := req.Inner
	w.logRequestHandle(r.Key, r.Operation, r.Args)
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
		set.Add(r.Context, r.Args[0])
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
		set.Remove(r.Context, r.Args[0])
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

func (w *Worker[F]) broadcast(event *pb.Event) {
	owners, err := w.hashRing.GetClosestN([]byte(event.Key), w.cfg.RepFactor)
	if err != nil {
		os.Exit(1)
	}
	for _, o := range owners {
		v := o.(util.Replica)
		if v == w.replica {
			continue
		}
		if v.Addr == w.replica.Addr {
			w.workers[v.WorkerID].PutEvent(event)
			continue
		}
		event.Dest = int32(v.WorkerID)
		client := pb.NewChiaveClient(w.connections[v.Addr])
		ctx, cancel := context.WithTimeout(context.Background(), rpcTimeout)
		defer cancel()
		_, err := client.ProcessEvent(ctx, event)
		if err != nil {
			w.logger.Errorf("ProcessEvent from %s to %s: %v", w.replica.String(), v.String(), err)
		}
	}
}

func (w *Worker[F]) logRequestHandle(key string, o pb.OP, args []string) {
	w.logger.Infof("worker %d handling %s(%v) on key %q", w.replica.WorkerID, o.String(), args, key)
}

func (w *Worker[_]) PutRequest(r LeaderRequest) {
	w.requests <- r
}

func (w *Worker[F]) PutEvent(e *pb.Event) {
	w.events <- e
}
