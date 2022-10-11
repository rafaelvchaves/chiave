package crdt

type Event struct {
	Source string // id of source replica
	Data any // CvRDT: state, CmRDT: update, dCvRDT: delta
}

type CRDT interface {
	GetEvents() []Event
	PersistEvents([]Event)
}


// Alternative:

// type CRDT interface {
// 	GetEvents() []any
// 	PersistEvents([]any)
// }
