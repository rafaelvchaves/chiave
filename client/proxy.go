package client

import (
	"context"
	pb "kvs/proto"
	"math/rand"

	"github.com/buraksezer/consistent"
	"google.golang.org/grpc"
)

type Proxy struct {
	connections    map[string]*grpc.ClientConn
	consistentHash *consistent.Consistent
	repFactor      int
}

func NewProxy(serverAddrs []string, threadsPerServer int, repFactor int) *Proxy {
	cfg := consistent.Config{
		PartitionCount:    5,
		ReplicationFactor: repFactor,
		Load:              1.25,
		Hasher:            hasher{},
	}
	p := &Proxy{
		consistentHash: consistent.New(nil, cfg),
		repFactor:      repFactor,
	}
	for _, addr := range serverAddrs {
		conn, err := grpc.Dial(addr)
		if err != nil {
			return nil
		}
		p.connections[addr] = conn
		for k := 0; k < threadsPerServer; k++ {
			p.consistentHash.Add(NewReplica(addr, k))
		}
	}
	return p
}

func (p *Proxy) Cleanup() error {
	for _, conn := range p.connections {
		if err := conn.Close(); err != nil {
			return err
		}
	}
	return nil
}

func (p *Proxy) clientOfOwner(key string) (pb.ChiaveClient, error) {
	owners, err := p.consistentHash.GetClosestN([]byte(key), p.repFactor)
	if err != nil {
		return nil, err
	}
	i := rand.Intn(p.repFactor)
	return pb.NewChiaveClient(p.connections[owners[i].String()]), nil
}

func (p *Proxy) Increment(key string) error {
	client, err := p.clientOfOwner(key)
	_, err = client.Increment(context.TODO(), &pb.Key{Id: key})
	if err != nil {
		return err
	}
	return nil
}

func (p *Proxy) Decrement(key string) error {
	client, err := p.clientOfOwner(key)
	_, err = client.Decrement(context.TODO(), &pb.Key{Id: key})
	if err != nil {
		return err
	}
	return nil
}

func (p *Proxy) Value(key string) (int64, error) {
	client, err := p.clientOfOwner(key)
	r, err := client.Value(context.TODO(), &pb.Key{Id: key})
	if err != nil {
		return 0, err
	}
	return r.Value, nil
}
