package util_test

import (
	"fmt"
	pb "kvs/proto"
	"kvs/util"
	"testing"
)

func TestDVV(t *testing.T) {
	var client1Ctx *pb.DVV
	// var client2Ctx *pb.DVV

	replicaA := replica{id: "a", ctx: &pb.DVV{}}
	// replicaB := replica{id: "b", ctx: &pb.DVV{}}

	client1Ctx = replicaA.Update(client1Ctx)
	// client1Ctx = replicaB.Update(client1Ctx)
	client1Ctx = replicaA.Update(client1Ctx)
	_ = replicaA.Update(client1Ctx)

	// client1Ctx = replicaB.Update(client1Ctx)

	// fmt.Println(util.String(client1Ctx))
}

type replica struct {
	ctx *pb.DVV
	id  string
}

func (r *replica) Update(clientCtx *pb.DVV) *pb.DVV {
	u := util.UpdateSingle(clientCtx, r.ctx, r.id)
	switch util.Compare(r.ctx, u) {
	case util.LT:
		r.ctx = u
	case util.CC:
		r.ctx = util.Join(u, r.ctx)
	default:
		fmt.Println("client context greater than update?")
	}
	fmt.Println(util.String(r.ctx))
	return r.ctx
}
