package main

import (
	"context"
	"fmt"
	pb "kvs/proto"
	"kvs/replica"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Leader struct {
	pb.UnimplementedChiaveServer
	addr string
	workers []replica.Worker
}

func NewLeader() *Leader {
	addr := "localhost:4747"
	workersPerReplica := 5
	workers := make([]replica.Worker, workersPerReplica)
	for i := 0; i < workersPerReplica; i++ {
		workers[i] = replica.NewWorker(addr, i, replica.NewCache())
	}
	return &Leader{
		addr: addr,
		workers: workers,
	}
}

func (l *Leader) Value(ctx context.Context, in *pb.Key) (*pb.ValueResponse, error) {
	fmt.Println("Value")
	return nil, nil
}

func (l *Leader) Increment(ctx context.Context, in *pb.Key) (*emptypb.Empty, error) {
	fmt.Println("Increment")
	return nil, nil
}

func (l *Leader) Decrement(ctx context.Context, in *pb.Key) (*emptypb.Empty, error) {
	fmt.Println("Decrement")
	return nil, nil
}

func main() {
	// addr, err := getAddress()
	// if err != nil {
	// 	fmt.Println("cannot get IP address")
	// 	return
	// }
	addr := "localhost:4747"
	fmt.Println(fmt.Sprintf("listening at address %s...", addr))
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		fmt.Println(err)
		return
	}
	leader := NewLeader()
	grpcServer := grpc.NewServer()
	pb.RegisterChiaveServer(grpcServer, leader)
	grpcServer.Serve(listener)
}

func getAddress() (string, error) {
	ip, err := getHost()
	port := "4747"
	if err != nil {
		return "", err
	}
	addr := net.JoinHostPort(ip, port)
	return addr, nil
}

func getHost() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", fmt.Errorf("could not get interface addresses")
	}
	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String(), nil
			}
		}
	}
	return "", fmt.Errorf("cannot find IP")
}
