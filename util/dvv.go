package util

import (
	pb "kvs/proto"

	"golang.org/x/exp/constraints"
)

type dot struct {
	r string
	n int32
}

type DVV struct {
	d  dot
	vv map[string]int32
}

func Lt(d1, d2 *pb.DVV) bool {
	dot := d1.Dot
	return dot == nil || dot.N <= d2.Clock[dot.Replica]
}

func Sync(D1, D2 []*pb.DVV) []*pb.DVV {
	var result []*pb.DVV
	for _, x := range D1 {
		include := true
		for _, y := range D2 {
			if Lt(x, y) {
				include = false
				break
			}
		}
		if include {
			result = append(result, x)
		}
	}
	for _, x := range D2 {
		include := true
		for _, y := range D1 {
			if Lt(x, y) {
				include = false
				break
			}
		}
		if include {
			result = append(result, x)
		}
	}
	return result
}

func dvvIDs(dvv *pb.DVV) []string {
	var result []string
	result = append(result, dvv.Dot.Replica)
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

func dvvCeil(dvv *pb.DVV, r string) int32 {
	dot := dvv.Dot
	if dot != nil && dot.Replica == r {
		return max(dot.N, dvv.Clock[r])
	}
	return dvv.Clock[r]
}

func ceil(dvvs []*pb.DVV, r string) int32 {
	m := int32(0)
	for _, dvv := range dvvs {
		m = max(m, dvvCeil(dvv, r))
	}
	return m
}

func Update(S []*pb.DVV, S_r []*pb.DVV, r string) *pb.DVV {
	result := &pb.DVV{
		Dot: &pb.Dot{
			Replica: r,
			N: ceil(S_r, r) + 1,
		},
		Clock: make(map[string]int32),
	}
	for _, i := range ids(S) {
		result.Clock[i] = ceil(S, i)
	}
	return result
}

func max[T constraints.Ordered](t1, t2 T) T {
	if t1 > t2 {
		return t1
	}
	return t2
}
