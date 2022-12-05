package util_test

import (
	pb "kvs/proto"
	"kvs/util"
	"testing"
)

func TestList(t *testing.T) {
	lst := []*pb.Dot{
		{
			N:       1,
			Replica: "a",
		},
		{
			N:       2,
			Replica: "b",
		},
		{
			N:       3,
			Replica: "c",
		},
		{
			N:       2,
			Replica: "b",
		},
	}
	util.Filter(func(dot *pb.Dot) bool {
		return dot.N >= 2
	}, &lst)
	if len(lst) != 2 {
		t.Errorf("len(lst): expected %d, got %d", 2, len(lst))
	}
}
