package state

import (
	"fmt"
	"kvs/util"
)

type Ord int

const (
	LT Ord = -1
	EQ Ord = 0
	GT Ord = 1
	CC Ord = 2
)

type GCounter struct {
	replica util.Replica
	vec     map[string]int64
}

func NewGCounter(replica util.Replica) GCounter {
	vec := make(map[string]int64)
	vec[replica.String()] = 0
	return GCounter{
		replica: replica,
		vec:     vec,
	}
}

func (g *GCounter) Increment() {
	id := g.replica.String()
	v, ok := g.vec[id]
	if !ok {
		g.vec[id] = 1
		return
	}
	g.vec[id] = v + 1
}

func (g *GCounter) Value() int {
	sum := 0
	for _, count := range g.vec {
		sum += int(count)
	}
	return sum
}

func (g *GCounter) SafeGet(r string) int64 {
	v, ok := g.vec[r]
	if !ok {
		return 0
	}
	return v
}

func (g *GCounter) Compare(o GCounter) Ord {
	ord := EQ
	for k, va := range g.vec {
		vb := o.SafeGet(k)
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

func (g GCounter) String() string {
	return fmt.Sprintf("%v", g.vec)
}

func (g *GCounter) Merge(ovec map[string]int64) {
	for k, vo := range ovec {
		v := g.SafeGet(k)
		g.vec[k] = util.Max(v, vo)
	}
}

func (g *GCounter) Copy() GCounter {
	cpy := NewGCounter(g.replica)
	for k, v := range g.vec {
		cpy.vec[k] = v
	}
	return cpy
}
