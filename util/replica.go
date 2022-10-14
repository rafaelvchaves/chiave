package util

import "fmt"

type Replica struct {
	Addr     string
	WorkerID int
}

func NewReplica(addr string, workerID int) Replica {
	return Replica{
		Addr:     addr,
		WorkerID: workerID,
	}
}

func (r Replica) String() string {
	return r.Addr + "," + fmt.Sprintf("%d", r.WorkerID)
}
