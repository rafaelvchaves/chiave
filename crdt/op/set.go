package op

import (
	"fmt"
	"kvs/crdt"
	pb "kvs/proto"
	"kvs/util"

	"github.com/google/uuid"
)

var _ crdt.Set = &Set{}
var _ crdt.CRDT[crdt.Op] = &Set{}

type Set struct {
	replica  util.Replica
	elements map[string][]string
	current  *pb.Event
}

func newSetEvent(replica util.Replica) *pb.Event {
	return &pb.Event{
		Source:   replica.String(),
		Datatype: pb.DT_Set,
		Data: &pb.Event_OpSet{
			OpSet: &pb.OpSet{},
		},
	}
}

func NewSet(replica util.Replica) *Set {
	return &Set{
		replica:  replica,
		elements: make(map[string][]string),
		current:  newSetEvent(replica),
	}
}

func (s *Set) Add(ctx *pb.Context, e string) {
	u := uuid.New().String()
	s.elements[e] = append(s.elements[e], u)
	eventData := s.current.GetOpSet()
	eventData.Operations = append(eventData.Operations, &pb.SetOperation{
		Op:      pb.SET_OP_ADD,
		Element: e,
		Tag:     u,
	})
}

func (s *Set) Remove(ctx *pb.Context, e string) {
	removeTags := s.elements[e]
	delete(s.elements, e)
	eventData := s.current.GetOpSet()
	eventData.Operations = append(eventData.Operations, &pb.SetOperation{
		Op:         pb.SET_OP_REM,
		Element:    e,
		RemoveTags: removeTags,
	})
}

func (s *Set) Value() []string {
	var result []string
	for e, tags := range s.elements {
		if len(tags) == 0 {
			continue
		}
		result = append(result, e)
	}
	return result
}

func (s *Set) String() string {
	set := s.Value()
	str := "{"
	for i, e := range set {
		str += e
		if i < len(s.elements) {
			str += ","
		}
	}
	return str + "}"
}

func (s *Set) PrepareEvent() *pb.Event {
	current := s.current
	s.current = newSetEvent(s.replica)
	return current
}

func (s *Set) PersistEvent(event *pb.Event) {
	os := event.GetOpSet()
	if os == nil {
		fmt.Println("warning: nil opset encountered in PersistEvent")
		return
	}
	for _, op := range os.Operations {
		tags := s.elements[op.Element]
		switch op.Op {
		case pb.SET_OP_ADD:
			tags = append(tags, op.Tag)
		case pb.SET_OP_REM:
			util.Filter(func(u string) bool { return !util.Contains(u, op.RemoveTags) }, &tags)
		}
		s.elements[op.Element] = tags
	}
}

func (s *Set) Context() *pb.Context {
	return &pb.Context{Dvv: &pb.DVV{}}
}

//lint:ignore U1000 Ignore unused warning: only used for debugging
func (s *Set) printState(header string) {
	fmt.Println(header)
	var result []string
	for e, tags := range s.elements {
		if len(tags) == 0 {
			continue
		}
		for _, t := range tags {
			result = append(result, fmt.Sprintf("(%s, %s)", e, t[:5]))
		}
	}
	fmt.Println(util.ListToString(result))
}
