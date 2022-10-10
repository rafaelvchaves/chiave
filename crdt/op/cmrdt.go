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

type CmRDT[S any] struct {
	state    S
	handlers map[string]Handler[S]
	queries map[string]Query[S]
	opQueue  []any
}

type handler[S any, D any] interface {
	Prepare(S, any) (D, bool)
	Effect(*S, D)
}

type Query[S any] interface {
	Query(S, any) string
}

type Handler[S any] interface {
	handler[S, any]
}

func Init[S any](state S, handlers map[string]Handler[S], queries map[string]Query[S]) CmRDT[S] {
	return CmRDT[S]{
		state: state,
		handlers: handlers,
		queries: queries,
	}
}

func (c *CmRDT[S]) Process(opName string, payload any) error {
	handler, ok := c.handlers[opName]
	if !ok {
		return fmt.Errorf("could not find handler for operation %q", opName)
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

func (c *CmRDT[S]) Query(queryType string, args any) (string, error) {
	query, ok := c.queries[queryType]
	if !ok {
		return "", fmt.Errorf("unknown query %q", queryType)
	}
	return query.Query(c.state, args), nil
}
