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

type Proxy struct {
	connections map[string]*grpc.ClientConn
	hashRing    *consistent.Consistent
	repFactor   int
}

func NewProxy(repFactor, threadsPerServer int) *Proxy {
	p := &Proxy{
		connections: util.GetConnections(),
		hashRing:    util.GetHashRing(),
		repFactor:   repFactor,
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
