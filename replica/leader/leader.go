package main

import (
	"context"
	"kvs/crdt"
	"kvs/crdt/delta"
	"kvs/crdt/generator"
	"kvs/crdt/op"
	"kvs/crdt/state"
	pb "kvs/proto"
	"kvs/replica/worker"
	"kvs/util"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type CRDTOption int

const (
	Delta CRDTOption = iota
	State
	Op
)

type leader[F crdt.Flavor] struct {
	pb.UnimplementedChiaveServer
	addr    string
	workers []worker.Worker[F]
}

type Leader interface {
	StartWorkers()
	pb.ChiaveServer
}

func NewLeader(opt CRDTOption) Leader {
	switch opt {
	case Delta:
		return LeaderWithFlavor[crdt.Delta](delta.Generator{})
	case State:
		return LeaderWithFlavor[crdt.State](state.Generator{})
	default:
		return LeaderWithFlavor[crdt.Op](op.Generator{})
	}
}

func LeaderWithFlavor[F crdt.Flavor](g generator.Generator[F]) *leader[F] {
	addr := "localhost:4747" // TODO: read from config
	workersPerReplica := 5   // TODO: read from config
	workers := make([]worker.Worker[F], workersPerReplica)
	for i := 0; i < workersPerReplica; i++ {
		r := util.NewReplica(addr, i)
		workers[i] = worker.New(r, worker.NewCache[F](), g)
	}
	return &leader[F]{
		addr:    addr,
		workers: workers,
	}
}

func (l *leader[_]) StartWorkers() {
	for _, w := range l.workers {
		go w.Start()
	}
}

func (l *leader[_]) Value(ctx context.Context, in *pb.Key) (*pb.ValueResponse, error) {
	// v, ok := l.workers[in.WorkerId].Get(in.Id)
	// c, ok := v.(int)
	return &pb.ValueResponse{}, nil
}

func (l *leader[_]) Get(ctx context.Context, in *pb.Key) (*pb.GetResponse, error) {
	req := worker.ClientRequest{
		Key:       in.Id,
		Operation: worker.Increment,
		Response:  make(chan worker.Response, 1),
	}
	l.workers[in.WorkerId].PutRequest(req)
	r := <-req.Response
	return &pb.GetResponse{
		Value:  r.Value,
		Exists: r.Exists,
	}, nil
}

func (l *leader[_]) Increment(ctx context.Context, in *pb.Key) (*emptypb.Empty, error) {
	req := worker.ClientRequest{
		Key:       in.Id,
		Operation: worker.Increment,
	}
	l.workers[in.WorkerId].PutRequest(req)
	return &emptypb.Empty{}, nil
}

func (l *leader[_]) Decrement(ctx context.Context, in *pb.Key) (*emptypb.Empty, error) {
	req := worker.ClientRequest{
		Key:       in.Id,
		Operation: worker.Decrement,
	}
	l.workers[in.WorkerId].PutRequest(req)
	return &emptypb.Empty{}, nil
}

func main() {
	leader := NewLeader(Op)
	leader.StartWorkers()
	addr := "localhost:4747"
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return
	}
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
