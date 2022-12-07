package main

import (
	"context"
	"flag"
	"fmt"
	"kvs/crdt"
	"kvs/crdt/generator"
	pb "kvs/proto"
	"kvs/replica/worker"
	"kvs/util"
	"net"

	"github.com/pkg/profile"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type CRDTOption int

const (
	Op CRDTOption = iota
	State
	Delta
)

var fromString = map[string]CRDTOption{
	"op":    Op,
	"state": State,
	"delta": Delta,
}

var toString = map[CRDTOption]string{
	Op:    "op",
	State: "state",
	Delta: "delta",
}

type leader[F crdt.Flavor] struct {
	pb.UnimplementedChiaveServer
	addr         string
	workers      []*worker.Worker[F]
	broadcasters []*worker.Broadcaster[F]
	logger       *util.Logger
}

type Leader interface {
	StartWorkers()
	pb.ChiaveServer
}

func NewLeader(addr string, opt CRDTOption) Leader {
	switch opt {
	case Delta:
		return leaderWithFlavor[crdt.Delta](addr, generator.Delta{})
	case State:
		return leaderWithFlavor[crdt.State](addr, generator.State{})
	default:
		return leaderWithFlavor[crdt.Op](addr, generator.Op{})
	}
}

func leaderWithFlavor[F crdt.Flavor](addr string, g generator.Generator[F]) *leader[F] {
	logger, err := util.NewLogger("log.txt")
	if err != nil {
		panic(fmt.Sprintf("error creating logger: %v", err))
	}
	cfg := util.LoadConfig()
	workers := make([]*worker.Worker[F], cfg.WorkersPerServer)
	broadcasters := make([]*worker.Broadcaster[F], cfg.WorkersPerServer)
	config := worker.Config[F]{
		Workers:     workers,
		Connections: util.GetConnections(),
		HashRing:    util.GetHashRing(),
		RepFactor:   cfg.RepFactor,
	}
	for i := 0; i < cfg.WorkersPerServer; i++ {
		r := util.NewReplica(addr, i)
		broadcasters[i] = worker.NewBroadcaster(r, config)
		workers[i] = worker.New(r, g, broadcasters[i], config)
	}
	return &leader[F]{
		addr:         addr,
		workers:      workers,
		broadcasters: broadcasters,
		logger:       logger,
	}
}

func (l *leader[_]) StartWorkers() {
	for i := 0; i < len(l.workers); i++ {
		go l.workers[i].Start()
		go l.broadcasters[i].Start()
	}
}

func (l *leader[_]) ProcessEvent(ctx context.Context, in *pb.Event) (*emptypb.Empty, error) {
	l.workers[in.Dest].PutEvent(in)
	return &emptypb.Empty{}, nil
}

func (l *leader[_]) GetCounter(ctx context.Context, in *pb.Request) (*pb.GetCounterResponse, error) {
	req := worker.LeaderRequest{
		Inner:    in,
		Response: make(chan worker.Response, 1),
	}
	l.workers[in.WorkerId].PutRequest(req)
	r := <-req.Response
	return &pb.GetCounterResponse{
		Context: r.Context,
		Value:   r.CounterValue,
	}, nil
}

func (l *leader[_]) GetSet(ctx context.Context, in *pb.Request) (*pb.GetSetResponse, error) {
	req := worker.LeaderRequest{
		Inner:    in,
		Response: make(chan worker.Response, 1),
	}
	l.workers[in.WorkerId].PutRequest(req)
	r := <-req.Response
	return &pb.GetSetResponse{
		Context: r.Context,
		Value:   r.SetValue,
	}, nil
}

func isAsync(op pb.OP) bool {
	return op == pb.OP_INCREMENT || op == pb.OP_DECREMENT
}

func (l *leader[_]) Write(ctx context.Context, in *pb.Request) (*pb.Response, error) {
	// fmt.Printf("Write at time %v\n", time.Now())
	if isAsync(in.Operation) {
		req := worker.LeaderRequest{
			Inner: in,
		}
		l.workers[in.WorkerId].PutRequest(req)
		return &pb.Response{}, nil
	}
	req := worker.LeaderRequest{
		Inner:    in,
		Response: make(chan worker.Response, 1),
	}
	l.workers[in.WorkerId].PutRequest(req)
	r := <-req.Response
	return &pb.Response{Context: r.Context}, nil
}

func main() {
	prof := flag.Bool("prof", false, "whether to run pprof")
	addr := flag.String("ip", util.LoadConfig().Addresses[0], "ip address to start leader at")
	flavor := flag.String("crdt", "op", "CRDT flavor (op, state, delta)")
	flag.Parse()
	if *prof {
		defer profile.Start(profile.ProfilePath(".")).Stop()
	}
	f := fromString[*flavor]
	fmt.Printf("using %s CRDTs\n", toString[f])
	leader := NewLeader(*addr, f)
	leader.StartWorkers()
	listener, err := net.Listen("tcp", *addr)
	if err != nil {
		return
	}
	fmt.Printf("serving at %q...\n", *addr)
	server := grpc.NewServer()
	pb.RegisterChiaveServer(server, leader)
	server.Serve(listener)
}
