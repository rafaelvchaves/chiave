package state

import (
	"fmt"
	pb "kvs/proto"
	"kvs/util"
)

type Ord int

const (
	LT Ord = -1
	EQ Ord = 0
	GT Ord = 1
	CC Ord = 2
)

func NewGCounter(replica string) *pb.GCounter {
	vec := make(map[string]int64)
	vec[replica] = 0
	return &pb.GCounter{
		Replica: replica,
		Vec:     vec,
	}
}

func Increment(g *pb.GCounter) {
	id := g.Replica
	v, ok := g.Vec[id]
	if !ok {
		g.Vec[id] = 1
		return
	}
	g.Vec[id] = v + 1
}

func Value(g *pb.GCounter) int {
	sum := 0
	for _, count := range g.Vec {
		sum += int(count)
	}
	return sum
}

func SafeGet(g *pb.GCounter, r string) int64 {
	v, ok := g.Vec[r]
	if !ok {
		return 0
	}
	return v
}

func Compare(g1, g2 *pb.GCounter) Ord {
	ord := EQ
	for k, va := range g1.Vec {
		vb := SafeGet(g2, k)
		switch {
		case ord == EQ && va > vb:
			ord = GT
		case ord == EQ && va < vb:
			ord = LT
		case ord == LT && va > vb:
			ord = CC
		case ord == GT && va < vb:
			ord = CC
		}
	}
	return ord
}

func String(g *pb.GCounter) string {
	return fmt.Sprintf("%v", g.Vec)
}

func Merge(g1, g2 *pb.GCounter) {
	for k, vo := range g2.Vec {
		v := SafeGet(g1, k)
		g1.Vec[k] = util.Max(v, vo)
	}
}

func Copy(g *pb.GCounter) *pb.GCounter {
	cpy := NewGCounter(g.Replica)
	for k, v := range g.Vec {
		cpy.Vec[k] = v
	}
	return cpy
}
