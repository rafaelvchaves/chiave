package crdt

import (
	"fmt"
)

// event comes in: operation name, payload.
// lookup handler in map[string]handler
// do operation, ok = handler.prepare(state, payload)
// if not ok: push event on back of queue
// otherwise, invoke handler.effect(state, operation) and
// send event to all other nodes.

type CRDT[S any] struct {
	state    S
	handlers map[string]Handler[S]
	opQueue  []any
}

type handler[S any, D any] interface {
	Prepare(S, any) (D, bool)
	Effect(*S, D)
}

type Handler[S any] interface {
	handler[S, any]
}

func Init[S any](handlers map[string]Handler[S]) CRDT[S] {
	return CRDT[S]{
		handlers: handlers,
	}
}

func test() {
	handlers := map[string]Handler[Graph]{
		"AddVertex": AddVertexHandler(struct{}{}),
	}
	g := Init(handlers)
	g.Process("AddVertex", "a")
}

func (c *CRDT[S]) Process(opName string, payload any) error {
	handler, ok := c.handlers[opName]
	if !ok {
		return fmt.Errorf("Could not find handler for operation %q", opName)
	}
	operation, ok := handler.Prepare(c.state, payload)
	if !ok {
		c.opQueue = append(c.opQueue, operation)
		return nil
	}
	handler.Effect(&c.state, operation)
	// send update to others
	return nil
}
