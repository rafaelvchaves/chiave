package client

import (
	"context"
	"fmt"
	pb "kvs/proto"
	"kvs/util"
	"math/rand"
	"sync"
	"time"

	"github.com/buraksezer/consistent"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
)

const (
	RPCTimeout = 10 * time.Second
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
	context     *pb.Context
	mu          sync.Mutex
}

func NewProxy() *Proxy {
	p := &Proxy{
		connections: util.GetConnections(),
		hashRing:    util.GetHashRing(),
		repFactor:   util.LoadConfig().RepFactor,
		context:     &pb.Context{},
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

func (p *Proxy) writeSync(key string, op pb.OP, args ...string) error {
	owners, err := p.ownersOf(key)
	if err != nil {
		return err
	}
	var lastError error
	for _, owner := range owners {
		r := owner.(util.Replica)
		client := pb.NewChiaveClient(p.connections[r.Addr])
		ctx, cancel := context.WithTimeout(context.Background(), RPCTimeout)
		defer cancel()
		p.mu.Lock()
		context := proto.Clone(p.context)
		p.mu.Unlock()
		res, err := client.Write(ctx, &pb.Request{
			Key:       key,
			WorkerId:  int32(r.WorkerID),
			Operation: op,
			Args:      args,
			Context:   context.(*pb.Context),
		})
		if err == nil {
			p.mu.Lock()
			p.context.Dvv = util.Sync(p.context.Dvv, res.Context.Dvv)
			p.mu.Unlock()
			return nil
		}
		lastError = err
	}
	return fmt.Errorf("failed to reach owners of key %q, last error = %v", key, lastError)
}

func (p *Proxy) writeAsync(key string, op pb.OP, args ...string) error {
	owners, err := p.ownersOf(key)
	if err != nil {
		return err
	}
	var lastError error
	for _, owner := range owners {
		r := owner.(util.Replica)
		client := pb.NewChiaveClient(p.connections[r.Addr])
		ctx, cancel := context.WithTimeout(context.Background(), RPCTimeout)
		defer cancel()
		_, err := client.Write(ctx, &pb.Request{
			Key:       key,
			WorkerId:  int32(r.WorkerID),
			Operation: op,
			Args:      args,
		})
		if err == nil {
			return nil
		}
		lastError = err
	}
	return fmt.Errorf("failed to reach owners of key %q, last error = %v", key, lastError)
}

func (p *Proxy) Increment(key ChiaveCounter) error {
	return p.writeAsync(key.string(), pb.OP_INCREMENT)
}

func (p *Proxy) Decrement(key ChiaveCounter) error {
	return p.writeAsync(key.string(), pb.OP_DECREMENT)
}

func (p *Proxy) GetCounter(key ChiaveCounter) (int64, error) {
	k := key.string()
	owners, err := p.ownersOf(k)
	if err != nil {
		return 0, err
	}
	var lastError error
	for _, owner := range owners {
		r := owner.(util.Replica)
		client := pb.NewChiaveClient(p.connections[r.Addr])
		ctx, cancel := context.WithTimeout(context.Background(), RPCTimeout)
		defer cancel()
		res, err := client.GetCounter(ctx, &pb.Request{
			Key:       k,
			WorkerId:  int32(r.WorkerID),
			Operation: pb.OP_GETCOUNTER,
		})
		if err == nil {
			return res.Value, nil
		}
		lastError = err
	}
	return 0, fmt.Errorf("failed to reach owners of key %q, last error = %v", key, lastError)
}

func (p *Proxy) AddSet(key ChiaveSet, element string) error {
	return p.writeSync(key.string(), pb.OP_ADDSET, element)
}

func (p *Proxy) RemoveSet(key ChiaveSet, element string) error {
	return p.writeSync(key.string(), pb.OP_REMOVESET, element)
}

func (p *Proxy) GetSet(key ChiaveSet) ([]string, error) {
	k := key.string()
	owners, err := p.ownersOf(k)
	if err != nil {
		return nil, err
	}
	var lastError error
	for _, owner := range owners {
		r := owner.(util.Replica)
		client := pb.NewChiaveClient(p.connections[r.Addr])
		ctx, cancel := context.WithTimeout(context.Background(), RPCTimeout)
		defer cancel()
		p.mu.Lock()
		context := proto.Clone(p.context)
		p.mu.Unlock()
		res, err := client.GetSet(ctx, &pb.Request{
			Key:       k,
			WorkerId:  int32(r.WorkerID),
			Operation: pb.OP_GETSET,
			Context:   context.(*pb.Context),
		})
		if err == nil {
			p.mu.Lock()
			p.context.Dvv = util.Sync(p.context.Dvv, res.Context.Dvv)
			p.mu.Unlock()
			return res.Value, nil
		}
		lastError = err
	}
	return nil, fmt.Errorf("failed to reach owners of key %q, last error = %v", key, lastError)
}

func (p *Proxy) Cleanup() error {
	for _, conn := range p.connections {
		if err := conn.Close(); err != nil {
			return err
		}
	}
	return nil
}
