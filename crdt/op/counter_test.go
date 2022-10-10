package crdt

import "testing"

func TestIncrement(t *testing.T) {
	c := Init(NewCounter(), counterHandlers, counterQueries)
	c.Process("INC", 1)
	c.Process("INC", 1)
	got, err := c.Query("VALUE", "")
	if err != nil {
		t.Fatalf("unexpected Query() error: %q", err)
	}
	if got != "2" {
		t.Errorf("Query(VALUE): expected 2, got %s", got)
	}
}

func TestDecrement(t *testing.T) {
	c := Init(NewCounter(), counterHandlers, counterQueries)
	c.Process("INC", 1)
	c.Process("INC", 1)
	c.Process("INC", 1)
	c.Process("DEC", 1)
	c.Process("DEC", 1)
	got, err := c.Query("VALUE", "")
	if err != nil {
		t.Fatalf("unexpected Query() error: %q", err)
	}
	if got != "1" {
		t.Errorf("Query(VALUE): expected 1, got %s", got)
	}
}
