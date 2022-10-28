package client

import (
	"context"
	"fmt"
	pb "kvs/proto"
	"kvs/util"
	"math/rand"
	"time"

	"github.com/buraksezer/consistent"
	"google.golang.org/grpc"
)

const (
	RPCTimeout = 5 * time.Second
)

type Proxy struct {
	connections map[string]*grpc.ClientConn
	hashRing    *consistent.Consistent
	repFactor   int
}

func NewProxy(repFactor int) *Proxy {
	p := &Proxy{
		connections: util.GetConnections(),
		hashRing:    util.GetHashRing(),
		repFactor:   repFactor,
	}
	return p
}

func (p *Proxy) ownersOf(key string) ([]consistent.Member, error) {
	owners, err := p.hashRing.GetClosestN([]byte(key), p.repFactor)
	if err != nil {
		return nil, err
	}
	// pick owner randomly
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(owners), func(i, j int) { owners[i], owners[j] = owners[j], owners[i] })
	return owners, nil
}

func (p *Proxy) Increment(key string) error {
	owners, err := p.ownersOf(key)
	if err != nil {
		return err
	}
	for _, owner := range owners {
		r := owner.(util.Replica)
		client := pb.NewChiaveClient(p.connections[r.Addr])
		ctx, cancel := context.WithTimeout(context.Background(), RPCTimeout)
		defer cancel()
		_, err := client.Increment(ctx, &pb.Key{
			Id:       key,
			WorkerId: int32(r.WorkerID),
		})
		if err == nil {
			// if no error occured during RPC, then return. Otherwise, try the next
			// owner.
			return nil
		}
	}
	return fmt.Errorf("failed to reach owners of key %q", key)
}

func (p *Proxy) Decrement(key string) error {
	owners, err := p.ownersOf(key)
	if err != nil {
		return err
	}
	for _, owner := range owners {
		r := owner.(util.Replica)
		client := pb.NewChiaveClient(p.connections[r.Addr])
		ctx, cancel := context.WithTimeout(context.Background(), RPCTimeout)
		defer cancel()
		_, err := client.Decrement(ctx, &pb.Key{
			Id:       key,
			WorkerId: int32(r.WorkerID),
		})
		if err == nil {
			return nil
		}
	}
	return fmt.Errorf("failed to reach owners of key %q", key)
}

func (p *Proxy) Get(key string) (string, error) {
	owners, err := p.ownersOf(key)
	if err != nil {
		return "", err
	}
	for _, owner := range owners {
		r := owner.(util.Replica)
		client := pb.NewChiaveClient(p.connections[r.Addr])
		ctx, cancel := context.WithTimeout(context.Background(), RPCTimeout)
		defer cancel()
		res, err := client.Get(ctx, &pb.Key{
			Id:       key,
			WorkerId: int32(r.WorkerID),
		})
		if err == nil {
			return res.Value, nil
		}
	}
	return "", fmt.Errorf("failed to reach owners of key %q", key)
}

func (p *Proxy) Cleanup() error {
	for _, conn := range p.connections {
		if err := conn.Close(); err != nil {
			return err
		}
	}
	return nil
}
