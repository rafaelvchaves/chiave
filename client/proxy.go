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

type Key interface {
	string() string
}

type ChiaveCounter string

func (c ChiaveCounter) string() string {
	return string(c)
}

type ChiaveSet string

func (c ChiaveSet) string() string {
	return string(c)
}

type ChiaveRegister string

type Proxy struct {
	connections map[string]*grpc.ClientConn
	hashRing    *consistent.Consistent
	repFactor   int
	context     []*pb.DVV
	// id uuid
}

func NewProxy() *Proxy {
	p := &Proxy{
		connections: util.GetConnections(),
		hashRing:    util.GetHashRing(),
		repFactor:   util.LoadConfig().RepFactor,
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

func (p *Proxy) Increment(key ChiaveCounter) error {
	owners, err := p.ownersOf(key.string())
	if err != nil {
		return err
	}
	for _, owner := range owners {
		r := owner.(util.Replica)
		client := pb.NewChiaveClient(p.connections[r.Addr])
		ctx, cancel := context.WithTimeout(context.Background(), RPCTimeout)
		defer cancel()
		_, err := client.Increment(ctx, &pb.Request{
			Key:      key.string(),
			WorkerId: int32(r.WorkerID),
			Context:  p.context,
		})
		if err == nil {
			return nil
		}
	}
	return fmt.Errorf("failed to reach owners of key %q", key)
}

func (p *Proxy) Decrement(key ChiaveCounter) error {
	owners, err := p.ownersOf(key.string())
	if err != nil {
		return err
	}
	for _, owner := range owners {
		r := owner.(util.Replica)
		client := pb.NewChiaveClient(p.connections[r.Addr])
		ctx, cancel := context.WithTimeout(context.Background(), RPCTimeout)
		defer cancel()
		_, err := client.Decrement(ctx, &pb.Request{
			Key:      key.string(),
			WorkerId: int32(r.WorkerID),
			Context:  p.context,
		})
		if err == nil {
			return nil
		}
	}
	return fmt.Errorf("failed to reach owners of key %q", key)
}

func (p *Proxy) Get(key Key) (string, error) {
	k := key.string()
	owners, err := p.ownersOf(k)
	if err != nil {
		return "", err
	}
	for _, owner := range owners {
		r := owner.(util.Replica)
		client := pb.NewChiaveClient(p.connections[r.Addr])
		ctx, cancel := context.WithTimeout(context.Background(), RPCTimeout)
		defer cancel()
		res, err := client.Get(ctx, &pb.Request{
			Key:      k,
			WorkerId: int32(r.WorkerID),
		})
		if err == nil {
			p.context = util.Sync(p.context, res.Context)
			return res.Value, nil
		}
	}
	return "", fmt.Errorf("failed to reach owners of key %q", key)
}

func (p *Proxy) AddSet(key ChiaveSet, element string) error {
	owners, err := p.ownersOf(key.string())
	if err != nil {
		return err
	}
	for _, owner := range owners {
		r := owner.(util.Replica)
		client := pb.NewChiaveClient(p.connections[r.Addr])
		ctx, cancel := context.WithTimeout(context.Background(), RPCTimeout)
		defer cancel()
		res, err := client.AddSet(ctx, &pb.Request{
			Key:      key.string(),
			WorkerId: int32(r.WorkerID),
			Context:  p.context,
			Args:     []string{element},
		})
		if err == nil {
			p.context = util.Sync(p.context, res.Context)
			return nil
		}
	}
	return fmt.Errorf("failed to reach owners of key %q", key)
}

func (p *Proxy) RemoveSet(key ChiaveSet, element string) error {
	owners, err := p.ownersOf(key.string())
	if err != nil {
		return err
	}
	for _, owner := range owners {
		r := owner.(util.Replica)
		client := pb.NewChiaveClient(p.connections[r.Addr])
		ctx, cancel := context.WithTimeout(context.Background(), RPCTimeout)
		defer cancel()
		res, err := client.RemoveSet(ctx, &pb.Request{
			Key:      key.string(),
			WorkerId: int32(r.WorkerID),
			Context:  p.context,
			Args:     []string{element},
		})
		if err == nil {
			p.context = util.Sync(p.context, res.Context)
			return nil
		}
	}
	return fmt.Errorf("failed to reach owners of key %q", key)
}

func (p *Proxy) Cleanup() error {
	for _, conn := range p.connections {
		if err := conn.Close(); err != nil {
			return err
		}
	}
	return nil
}
