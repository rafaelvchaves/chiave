package util_test

import (
	"fmt"
	pb "kvs/proto"
	"kvs/util"
	"testing"
)

func TestDVV(t *testing.T) {
	d1 := &pb.DVV{
		Clock:   make(map[string]int32),
		Sibling: "REMOVE_A",
	}
	d2 := &pb.DVV{
		Clock:   make(map[string]int32),
		Sibling: "ADD_A",
	}
	d1 = util.Update(nil, []*pb.DVV{d1}, "0")
	d2 = util.Update(nil, []*pb.DVV{d2}, "1")
	fmt.Println(util.String([]*pb.DVV{d1, d2}...))
	fmt.Println(util.Join([]*pb.DVV{d1, d2}))
	reduce := func([]string) string {
		return "ADD_A"
	}
	fmt.Println(util.String(util.Reconcile(reduce, []*pb.DVV{d1, d2}, "0")))
	fmt.Println(util.String(util.Reconcile(reduce, []*pb.DVV{d1, d2}, "1")))
}
