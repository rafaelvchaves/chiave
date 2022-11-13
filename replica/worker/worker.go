package worker

import (
	"context"
	"kvs/crdt"
	"kvs/crdt/generator"
	pb "kvs/proto"
	"kvs/util"
	"os"
	"time"

	"github.com/buraksezer/consistent"
	"google.golang.org/grpc"
)

const RPCTimeout = 10 * time.Second

type Worker[F crdt.Flavor] struct {
	replica     util.Replica
	generator   generator.Generator[F]
	kvs         Store[F]
	contexts    map[string][]*pb.DVV
	requests    chan ClientRequest
	events      chan *pb.Event
	hashRing    *consistent.Consistent
	connections map[string]*grpc.ClientConn
	cfg         util.Config
	logger      *util.Logger
}

type Operation int

const (
	Get Operation = iota
	Increment
	Decrement
	AddSet
	RemoveSet
)

func (o Operation) String() string {
	switch o {
	case Get:
		return "GET"
	case Increment:
		return "INC"
	case Decrement:
		return "DEC"
	case AddSet:
		return "ADD"
	case RemoveSet:
		return "REMOVE"
	}
	return "UNKNOWN"
}

type ClientRequest struct {
	Key       string
	Operation Operation
	Params    []string
	Response  chan Response
	Context   *pb.Context
}

type Response = struct {
	Value   string
	Exists  bool
	Context []*pb.DVV
}

func New[F crdt.Flavor](replica util.Replica, generator generator.Generator[F], logger *util.Logger) *Worker[F] {
	return &Worker[F]{
		generator:   generator,
		replica:     replica,
		kvs:         NewCache[F](),
		contexts:    make(map[string][]*pb.DVV),
		requests:    make(chan ClientRequest, 10000),
		events:      make(chan *pb.Event, 10000),
		hashRing:    util.GetHashRing(),
		connections: util.GetConnections(),
		cfg:         util.LoadConfig(),
		logger:      logger,
	}
}

func (w *Worker[F]) Start() {
	requestDeadline := 50 * time.Millisecond
	for {
		// set of keys modified in this epoch
		changeset := util.NewSet[string]()
		timeout := time.After(requestDeadline)
	reqLoop:
		for {
			// phase 1: receive client requests and convert to events
			select {
			case <-timeout:
				break reqLoop
			case req := <-w.requests:
				w.process(req)
				if req.Operation != Get {
					changeset.Add(req.Key)
				}
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

		timeout = time.After(requestDeadline)
	eventLoop:
		for {
			select {
			case event := <-w.events:
				v := w.kvs.GetOrDefault(event.Key, w.generator.New(event.Datatype, w.replica))
				v.PersistEvent(event)
				w.kvs.Put(event.Key, v)
			case <-timeout:
				break eventLoop
			}
		}
	}
}

func (w *Worker[F]) process(r ClientRequest) {
	w.logRequestHandle(r.Key, r.Operation, r.Params)
	switch r.Operation {
	case Increment:
		v := w.kvs.GetOrDefault(r.Key, w.generator.New(pb.DT_Counter, w.replica))
		counter, ok := v.(crdt.Counter)
		if !ok {
			return
		}
		counter.Increment()
		w.kvs.Put(r.Key, v)
	case Decrement:
		v := w.kvs.GetOrDefault(r.Key, w.generator.New(pb.DT_Counter, w.replica))
		counter, ok := v.(crdt.Counter)
		if !ok {
			return
		}
		counter.Decrement()
		w.kvs.Put(r.Key, v)
	case Get:
		var value string
		v, ok := w.kvs.Get(r.Key)
		if ok {
			value = v.String()
		}
		r.Response <- Response{
			Value:   value,
			Exists:  ok,
			Context: w.contexts[r.Key],
		}
	case AddSet:
		v := w.kvs.GetOrDefault(r.Key, w.generator.New(pb.DT_Set, w.replica))
		set, ok := v.(crdt.Set)
		if !ok {
			return
		}
		set.Add(r.Params[0])
		w.kvs.Put(r.Key, v)
	case RemoveSet:
		v := w.kvs.GetOrDefault(r.Key, w.generator.New(pb.DT_Set, w.replica))
		set, ok := v.(crdt.Set)
		if !ok {
			return
		}
		set.Remove(r.Params[0])
		w.kvs.Put(r.Key, v)
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
		event.Dest = int32(v.WorkerID)
		client := pb.NewChiaveClient(w.connections[v.Addr])
		ctx, cancel := context.WithTimeout(context.Background(), RPCTimeout)
		defer cancel()
		// w.logger.Infof("worker %d is sending %v to worker %d", w.replica.WorkerID, event.Data, v.WorkerID)
		_, err := client.ProcessEvent(ctx, event)
		if err != nil {
			w.logger.Errorf("ProcessEvent from %s to %s: %v", w.replica.String(), v.String(), err)
		}
	}
}

func (w *Worker[F]) logRequestHandle(key string, o Operation, args []string) {
	w.logger.Infof("worker %d handling %s(%v) on key %q", w.replica.WorkerID, o.String(), args, key)
}

func (w *Worker[_]) PutRequest(r ClientRequest) {
	w.requests <- r
}

func (w *Worker[F]) PutEvent(e *pb.Event) {
	w.events <- e
}
