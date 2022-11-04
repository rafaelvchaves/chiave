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
	addr    string
	workers []*worker.Worker[F]
	logger  *util.Logger
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
	l, err := util.NewLogger("log.txt")
	if err != nil {
		panic(fmt.Sprintf("error creating logger: %v", err))
	}
	cfg := util.LoadConfig()
	workers := make([]*worker.Worker[F], cfg.WorkersPerServer)
	for i := 0; i < cfg.WorkersPerServer; i++ {
		r := util.NewReplica(addr, i)
		workers[i] = worker.New(r, g, l)
	}
	return &leader[F]{
		addr:    addr,
		workers: workers,
		logger:  l,
	}
}

func (l *leader[_]) StartWorkers() {
	for _, w := range l.workers {
		go w.Start()
	}
}

func (l *leader[_]) ProcessEvent(ctx context.Context, in *pb.Event) (*emptypb.Empty, error) {
	l.workers[in.Dest].PutEvent(in)
	return &emptypb.Empty{}, nil
}

func (l *leader[_]) Get(ctx context.Context, in *pb.Request) (*pb.GetResponse, error) {
	req := worker.ClientRequest{
		Key:       in.Key,
		Operation: worker.Get,
		Response:  make(chan worker.Response, 1),
	}
	l.workers[in.WorkerId].PutRequest(req)
	r := <-req.Response
	return &pb.GetResponse{
		Value:   r.Value,
		Exists:  r.Exists,
		Context: r.Context,
	}, nil
}

func (l *leader[_]) Increment(ctx context.Context, in *pb.Request) (*pb.Response, error) {
	req := worker.ClientRequest{
		Key:       in.Key,
		Operation: worker.Increment,
		Context:   in.Context,
		Response:  make(chan worker.Response, 1),
	}
	l.workers[in.WorkerId].PutRequest(req)
	r := <-req.Response
	return &pb.Response{Context: r.Context}, nil
}

func (l *leader[_]) Decrement(ctx context.Context, in *pb.Request) (*pb.Response, error) {
	req := worker.ClientRequest{
		Key:       in.Key,
		Operation: worker.Decrement,
		Context:   in.Context,
		Response:  make(chan worker.Response, 1),
	}
	l.workers[in.WorkerId].PutRequest(req)
	r := <-req.Response
	return &pb.Response{Context: r.Context}, nil
}

func (l *leader[_]) AddSet(ctx context.Context, in *pb.Request) (*pb.Response, error) {
	req := worker.ClientRequest{
		Key:       in.Key,
		Operation: worker.AddSet,
		Context:   in.Context,
		Response:  make(chan worker.Response, 1),
		Params:    in.Args,
	}
	l.workers[in.WorkerId].PutRequest(req)
	r := <-req.Response
	return &pb.Response{Context: r.Context}, nil
}

func main() {
	addr := flag.String("a", "localhost:4747", "ip address to start leader at")
	flavor := flag.String("crdt", "op", "CRDT flavor (op, state, delta)")
	flag.Parse()
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

// func getAddress() (string, error) {
// 	ip, err := getHost()
// 	port := "4747"
// 	if err != nil {
// 		return "", err
// 	}
// 	addr := net.JoinHostPort(ip, port)
// 	return addr, nil
// }

// func getHost() (string, error) {
// 	addrs, err := net.InterfaceAddrs()
// 	if err != nil {
// 		return "", fmt.Errorf("could not get interface addresses")
// 	}
// 	for _, address := range addrs {
// 		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
// 			if ipnet.IP.To4() != nil {
// 				return ipnet.IP.String(), nil
// 			}
// 		}
// 	}
// 	return "", fmt.Errorf("cannot find IP")
// }
