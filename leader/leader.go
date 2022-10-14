package main

import (
	"context"
	pb "kvs/proto"
	"kvs/replica"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Leader struct {
	pb.UnimplementedChiaveServer
	addr    string
	workers []replica.Worker
}

func NewLeader() *Leader {
	addr := "localhost:4747" // TODO: read from config
	workersPerReplica := 5 // TODO: read from config
	workers := make([]replica.Worker, workersPerReplica)
	for i := 0; i < workersPerReplica; i++ {
		workers[i] = replica.NewWorker(addr, i, replica.NewCache())
	}
	return &Leader{
		addr:    addr,
		workers: workers,
	}
}

func (l *Leader) StartWorkers() {
	for _, w := range l.workers {
		go w.Start()
	}
}

func (l *Leader) Value(ctx context.Context, in *pb.Key) (*pb.ValueResponse, error) {
	return &pb.ValueResponse{}, nil
}

func (l *Leader) Increment(ctx context.Context, in *pb.Key) (*emptypb.Empty, error) {
	req := replica.ClientRequest{
		Key:       in.Id,
		Operation: replica.Increment,
	}
	l.workers[in.WorkerId].AddRequest(req)
	return &emptypb.Empty{}, nil
}

func (l *Leader) Decrement(ctx context.Context, in *pb.Key) (*emptypb.Empty, error) {
	req := replica.ClientRequest{
		Key:       in.Id,
		Operation: replica.Decrement,
	}
	l.workers[in.WorkerId].AddRequest(req)
	return &emptypb.Empty{}, nil
}

func main() {
	leader := NewLeader()
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
