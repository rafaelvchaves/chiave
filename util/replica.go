package util

import (
	"fmt"

	"github.com/buraksezer/consistent"
	farmhash "github.com/leemcloughlin/gofarmhash"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var cfg = LoadConfig()

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
	connections := make(map[string]*grpc.ClientConn)
	for _, addr := range cfg.Addresses {
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
	hashCfg := consistent.Config{
		PartitionCount:    cfg.PartitionCount,
		ReplicationFactor: cfg.RepFactor,
		Load:              cfg.Load,
		Hasher:            hasher{},
	}
	hashRing := consistent.New(nil, hashCfg)
	for _, addr := range cfg.Addresses {
		for k := 0; k < cfg.WorkersPerServer; k++ {
			hashRing.Add(NewReplica(addr, k))
		}
	}
	return hashRing
}
