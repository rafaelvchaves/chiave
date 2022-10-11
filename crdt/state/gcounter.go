package state

import "kvs/util"

type Ord int

const (
	LT Ord = -1
	EQ Ord = 0
	GT Ord = 1
	CC Ord = 2
)

type GCounter struct {
	id  string
	vec map[string]int
}

func NewGCounter(id string) GCounter {
	vec := make(map[string]int)
	vec[id] = 0
	return GCounter{
		id:  id,
		vec: vec,
	}
}

func (g *GCounter) Increment() {
	v, ok := g.vec[g.id]
	if !ok {
		g.vec[g.id] = 1
		return
	}
	g.vec[g.id] = v + 1
}

func (g *GCounter) Value() int {
	sum := 0
	for _, count := range g.vec {
		sum += count
	}
	return sum
}

func (g *GCounter) SafeGet(r string) int {
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

func (g *GCounter) Merge(o GCounter) {
	for k, vo := range o.vec {
		v := g.SafeGet(k)
		g.vec[k] = util.Max(v, vo)
	}
}
