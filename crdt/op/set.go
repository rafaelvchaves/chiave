package op

import (
	"fmt"
	pb "kvs/proto"
	"kvs/util"

	"github.com/google/uuid"
)

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

func (s *Set) Add(e string) {
	u := uuid.New().String()
	s.elements[e] = append(s.elements[e], u)
	eventData := s.current.GetOpSet()
	eventData.Operations = append(eventData.Operations, &pb.SetOperation{
		Op:      pb.SET_OP_ADD,
		Element: e,
		Tag:     u,
	})
}
func (s *Set) Remove(e string) {
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
	str := "{"
	i := 0
	for e := range s.elements {
		i++
		str += e
		if i < len(s.elements) {
			str += ","
		}
	}
	return str + "}"
}

func (s *Set) GetEvent() *pb.Event {
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
			filter(func(u string) bool { return !contains(u, op.RemoveTags) }, &tags)
		}
		s.elements[op.Element] = tags
	}
}

func filter[T any](p func(T) bool, lst *[]T) {
	i := 0
	for _, x := range *lst {
		if p(x) {
			(*lst)[i] = x
			i++
		}
	}
	*lst = (*lst)[:i]
}

func contains[T comparable](target T, lst []T) bool {
	for _, x := range lst {
		if x == target {
			return true
		}
	}
	return false
}
