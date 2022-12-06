package worker

import (
	"context"
	"fmt"
	"kvs/crdt"
	pb "kvs/proto"
	"kvs/util"
	"os"

	"github.com/buraksezer/consistent"
	"google.golang.org/grpc"
)

type Config[F crdt.Flavor] struct {
	Workers     []*Worker[F]
	Connections map[string]*grpc.ClientConn
	HashRing    *consistent.Consistent
	RepFactor   int
}

type Broadcaster[F crdt.Flavor] struct {
	replica util.Replica
	events  chan *pb.Event
	config  Config[F]
}

func NewBroadcaster[F crdt.Flavor](replica util.Replica, config Config[F]) *Broadcaster[F] {
	eventBufferSize := 100000
	return &Broadcaster[F]{
		replica: replica,
		events:  make(chan *pb.Event, eventBufferSize),
		config:  config,
	}
}

func (b *Broadcaster[F]) broadcast(event *pb.Event) {
	owners, err := b.config.HashRing.GetClosestN([]byte(event.Key), b.config.RepFactor)
	if err != nil {
		os.Exit(1)
	}
	for _, o := range owners {
		v := o.(util.Replica)
		if v == b.replica {
			continue
		}
		if v.Addr == b.replica.Addr {
			b.config.Workers[v.WorkerID].PutEvent(event)
			continue
		}
		event.Dest = int32(v.WorkerID)
		client := pb.NewChiaveClient(b.config.Connections[v.Addr])
		ctx, cancel := context.WithTimeout(context.Background(), rpcTimeout)
		defer cancel()
		_, err := client.ProcessEvent(ctx, event)
		if err != nil {
			fmt.Printf("ProcessEvent from %s to %s: %v\n", b.replica.String(), v.String(), err)
		}
	}
}

func (b Broadcaster[_]) Start() {
	for event := range b.events {
		go b.broadcast(event)
	}
}

func (b Broadcaster[_]) Send(event *pb.Event) {
	b.events <- event
}
