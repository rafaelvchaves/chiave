package client

import (
	"context"
	"fmt"
	pb "kvs/proto"
	"kvs/util"
	"math/rand"
	"time"

	"github.com/buraksezer/consistent"
	farmhash "github.com/leemcloughlin/gofarmhash"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type hasher struct{}

func (hasher) Sum64(data []byte) uint64 {
	return farmhash.Hash64(data)
}

type Proxy struct {
	connections map[string]*grpc.ClientConn
	hashRing    *consistent.Consistent
	repFactor   int
}

func NewProxy(repFactor, threadsPerServer int) *Proxy {
	serverAddrs := []string{
		"localhost:4747",
	}
	cfg := consistent.Config{
		PartitionCount:    5, // TODO: change?
		ReplicationFactor: repFactor,
		Load:              1.25, // TODO: change?
		Hasher:            hasher{},
	}
	p := &Proxy{
		connections: make(map[string]*grpc.ClientConn),
		hashRing:    consistent.New(nil, cfg),
		repFactor:   repFactor,
	}
	for _, addr := range serverAddrs {
		// TODO: change from insecure credentials
		conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			continue
		}
		p.connections[addr] = conn
		for k := 0; k < threadsPerServer; k++ {
			p.hashRing.Add(util.NewReplica(addr, k))
		}
	}
	return p
}

func (p *Proxy) chooseLeader(key string) (pb.ChiaveClient, error) {
	owners, err := p.hashRing.GetClosestN([]byte(key), p.repFactor)
	if err != nil {
		return nil, err
	}
	// pick owner randomly
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(owners), func(i, j int) { owners[i], owners[j] = owners[j], owners[i] })
	for _, owner := range owners {
		addr := owner.(util.Replica).Addr
		conn, ok := p.connections[addr]
		if !ok {
			continue
		}
		return pb.NewChiaveClient(conn), nil
	}
	return nil, fmt.Errorf("connection to key owners failed")
}

func (p *Proxy) Increment(key string) error {
	leader, err := p.chooseLeader(key)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err = leader.Increment(ctx, &pb.Key{Id: key})
	if err != nil {
		return err
	}
	return nil
}

func (p *Proxy) Decrement(key string) error {
	leader, err := p.chooseLeader(key)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err = leader.Decrement(ctx, &pb.Key{Id: key})
	if err != nil {
		return err
	}
	return nil
}

func (p *Proxy) Value(key string) (int64, error) {
	leader, err := p.chooseLeader(key)
	if err != nil {
		return 0, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	r, err := leader.Value(ctx, &pb.Key{Id: key})
	if err != nil {
		return 0, err
	}
	return r.Value, nil
}

func (p *Proxy) Cleanup() error {
	for _, conn := range p.connections {
		if err := conn.Close(); err != nil {
			return err
		}
	}
	return nil
}
