package client

import (
	"fmt"

	farmhash "github.com/leemcloughlin/gofarmhash"
)

type replica struct {
	addr string
	workerID int
}

func NewReplica(addr string, workerID int) replica {
	return replica{
		addr: addr,
		workerID: workerID,
	}
}

func (r replica) String() string {
	return r.addr + fmt.Sprintf("%d", r.workerID)
}

type hasher struct{}

func (hasher) Sum64(data []byte) uint64 {
	return farmhash.Hash64(data)
}
