package util

import (
	"fmt"

	"github.com/buraksezer/consistent"
	farmhash "github.com/leemcloughlin/gofarmhash"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Replica struct {
	Addr     string
	WorkerID int
}

func NewReplica(addr string, workerID int) Replica {
	return Replica{
		Addr:     addr,
		WorkerID: workerID,
	}
}

func (r Replica) String() string {
	return r.Addr + "," + fmt.Sprintf("%d", r.WorkerID)
}

type hasher struct{}

func (hasher) Sum64(data []byte) uint64 {
	return farmhash.Hash64(data)
}

func GetConnections() map[string]*grpc.ClientConn {
	addrs := []string{
		"localhost:4747",
	}
	connections := make(map[string]*grpc.ClientConn)
	for _, addr := range addrs {
		// TODO: change from insecure credentials
		conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			continue
		}
		connections[addr] = conn
	}
	return connections
}

func GetHashRing() *consistent.Consistent {
	addrs := []string{
		"localhost:4747",
	}
	repFactor := 3
	workersPerReplica := 5
	cfg := consistent.Config{
		PartitionCount:    5, // TODO: change?
		ReplicationFactor: repFactor,
		Load:              1.25, // TODO: change?
		Hasher:            hasher{},
	}
	hashRing := consistent.New(nil, cfg)
	for _, addr := range addrs {
		for k := 0; k < workersPerReplica; k++ {
			hashRing.Add(NewReplica(addr, k))
		}
	}
	return hashRing
}
