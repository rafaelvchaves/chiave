package delta

import (
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
	delta   map[string]int64
}

func NewGCounter(replica util.Replica) GCounter {
	vec := make(map[string]int64)
	vec[replica.String()] = 0
	return GCounter{
		replica: replica,
		vec:     vec,
		delta:   make(map[string]int64),
	}
}

func (g *GCounter) Increment() {
	id := g.replica.String()
	v, ok := g.vec[id]
	if !ok {
		g.vec[id] = 1
		g.delta[id] = 1
		return
	}
	g.vec[id] = v + 1
	g.delta[id] = v + 1
}

func (g *GCounter) Value() int {
	sum := 0
	for _, count := range g.vec {
		sum += int(count)
	}
	return sum
}

func safeGet(m map[string]int64, r string) int64 {
	v, ok := m[r]
	if !ok {
		return 0
	}
	return v
}

func (g *GCounter) Compare(o GCounter) Ord {
	ord := EQ
	for k, va := range g.vec {
		vb := safeGet(o.vec, k)
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

func (g *GCounter) Merge(delta map[string]int64) {
	for k, vo := range delta {
		v := safeGet(g.vec, k)
		g.vec[k] = util.Max(v, vo)
		vd := safeGet(g.delta, k)
		g.delta[k] = util.Max(vd, vo)
	}
}

func (g *GCounter) GetDelta() map[string]int64 {
	d := g.delta
	g.delta = make(map[string]int64)
	return d
}
