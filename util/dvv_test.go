package util_test

import (
	"fmt"
	pb "kvs/proto"
	"kvs/util"
	"testing"
)

func TestDVV(t *testing.T) {
	// d1 := &pb.DVV{
	// 	Clock:   make(map[string]int32),
	// 	Sibling: "REMOVE_A",
	// }
	// d2 := &pb.DVV{
	// 	Clock:   make(map[string]int32),
	// 	Sibling: "ADD_A",
	// }
	// //
	// d1 = util.Update(nil, []*pb.DVV{d1}, "0")
	// d2 = util.Update(nil, []*pb.DVV{d2}, "1")
	// fmt.Println(util.String([]*pb.DVV{d1, d2}...))
	// reduce := func([]string) string {
	// 	return "ADD_A"
	// }
	// new1 := util.Reconcile(reduce, []*pb.DVV{d1, d2}, "0")
	// new2 := util.Reconcile(reduce, []*pb.DVV{d1, d2}, "1")
	var client1Ctx *pb.DVV
	var client2Ctx *pb.DVV

	replicaA := replica{id: "a", ctx: &pb.DVV{}}
	replicaB := replica{id: "b", ctx: &pb.DVV{}}

	// Client 1 performs an operation and goes to replica a
	client1Ctx = replicaA.Update(client1Ctx)

	// Client 1 performs an operation and goes to replica b
	client1Ctx = replicaB.Update(client1Ctx)
	client1Ctx = replicaB.Update(client1Ctx)
	client1Ctx = replicaB.Update(client1Ctx)
	client1Ctx = replicaB.Update(client1Ctx)
	client1Ctx = replicaA.Update(client1Ctx)
	client2Ctx = replicaA.Update(client2Ctx)
	// client2Ctx = replicaA.Update(client2Ctx)

	// replicaA.Sync(replicaB)
	// replicaB.Sync(replicaA)
	fmt.Println(util.String(client1Ctx))
	fmt.Println(util.String(client2Ctx))
	// fmt.Println(util.String(replicaA.ctx))
	// fmt.Println(util.String(replicaB.ctx))
	// fmt.Println(util.String(util.Join([]*pb.DVV{replicaA.ctx, replicaB.ctx})))

}

type replica struct {
	ctx *pb.DVV
	id  string
}

func (r *replica) Update(clientCtx *pb.DVV) *pb.DVV {
	u := util.UpdateSingle(clientCtx, r.ctx, r.id)
	fmt.Println(util.String(u))
	if util.Lt(r.ctx, u) {
		r.ctx = u
	} else if util.Lt(u, r.ctx) {
		fmt.Println("client context greater than update?")
	} else {
		fmt.Printf("concurrent: %s and %s\n", util.String(u), util.String(r.ctx))
	}
	return r.ctx
}

func (r *replica) Sync(o replica) {
	if util.Lt(r.ctx, o.ctx) {
		// r.ctx = o.ctx
		fmt.Printf("%s lt %s \n", r.id, o.id)

	} else if util.Lt(o.ctx, r.ctx) {
		fmt.Printf("%s lt %s \n", o.id, r.id)
		// fmt.Println("client context greater than update?")
	} else {
		fmt.Printf("%s cc %s \n", o.id, r.id)
	}
}
