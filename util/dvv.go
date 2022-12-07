package util

import (
	"fmt"
	pb "kvs/proto"
)

type Ord int

const (
	LT Ord = iota
	GT
	CC
)

func displayMap(m map[string]int64) string {
	str := "{"
	i := 0
	for k, v := range m {
		i++
		str += fmt.Sprintf("%s:%d", k, v)
		if i < len(m) {
			str += ","
		}
	}
	return str + "}"
}

func String(dvvs ...*pb.DVV) string {
	str := ""
	for _, dvv := range dvvs {
		if dvv == nil {
			str += "(<nil>)"
		} else {
			str += fmt.Sprintf("(%s:%d, %s)", dvv.Dot.Replica, dvv.Dot.N, displayMap(dvv.Clock))
		}
	}
	return str
}

func Join(D ...*pb.DVV) *pb.DVV {
	var r string
	var N int64
	for _, dvv := range D {
		if dvv.Dot.N >= N {
			r = dvv.Dot.Replica
			N = dvv.Dot.N
		}
	}
	result := &pb.DVV{
		Clock: make(map[string]int64),
		Dot: &pb.Dot{
			Replica: r,
			N:       N,
		},
	}
	for _, i := range ids(D) {
		result.Clock[i] = ceil(D, i)
	}
	return result
}

func ContainedIn(dot *pb.Dot, dvv *pb.DVV) bool {
	if dot == nil {
		return true
	}
	return dot.N < dvv.Clock[dot.Replica] || (dvv.Dot.Replica == dot.Replica && dot.N <= dvv.Dot.N)
}

func Compare(d1, d2 *pb.DVV) Ord {
	if lt(d1, d2) {
		return LT
	} else if lt(d2, d1) {
		return GT
	}
	return CC
}

func lt(d1, d2 *pb.DVV) bool {
	if d1 == nil {
		return true
	}
	if d2 == nil {
		return false
	}
	dot := d1.Dot
	return dot == nil || dot.N <= d2.Clock[dot.Replica]
}

func Sync(d1, d2 *pb.DVV) *pb.DVV {
	switch Compare(d1, d2) {
	case GT:
		return d1
	case LT:
		return d2
	}
	return Join(d1, d2)
}

func dvvIDs(dvv *pb.DVV) []string {
	if dvv == nil {
		return nil
	}
	var result []string
	if dvv.Dot != nil {
		result = append(result, dvv.Dot.Replica)
	}
	for r := range dvv.Clock {
		result = append(result, r)
	}
	return result
}

func ids(dvvs []*pb.DVV) []string {
	var result []string
	for _, dvv := range dvvs {
		result = append(result, dvvIDs(dvv)...)
	}
	return result
}

func dvvCeil(dvv *pb.DVV, r string) int64 {
	dot := dvv.Dot
	if dot != nil && dot.Replica == r {
		return Max(dot.N, dvv.Clock[r])
	}
	return dvv.Clock[r]
}

func ceil(dvvs []*pb.DVV, r string) int64 {
	m := int64(0)
	for _, dvv := range dvvs {
		m = Max(m, dvvCeil(dvv, r))
	}
	return m
}

func UpdateSingle(d1 *pb.DVV, d2 *pb.DVV, r string) *pb.DVV {
	result := &pb.DVV{
		Dot: &pb.Dot{
			Replica: r,
			N:       dvvCeil(d2, r) + 1,
		},
		Clock: make(map[string]int64),
	}
	for _, i := range dvvIDs(d1) {
		result.Clock[i] = dvvCeil(d1, i)
	}
	return result
}
