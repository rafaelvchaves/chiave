package crdt

import "testing"

func TestAddSet(t *testing.T) {
	s := Init(NewORSet(), setHandlers, setQueries)
	s.Process("ADD", "a")
	got, err := s.Query("EXISTS", "a")
	if err != nil {
		t.Fatalf("unexpected Query() error: %q", err)
	}
	if got != "true" {
		t.Errorf("Query(EXISTS a): expected true, got false")
	}
	got, err = s.Query("EXISTS", "b")
	if err != nil {
		t.Fatalf("unexpected Query() error: %q", err)
	}
	if got != "false" {
		t.Errorf("Query(EXISTS b): expected false, got true")
	}
}

func TestRemoveSet(t *testing.T) {
	s := Init(NewORSet(), setHandlers, setQueries)
	s.Process("ADD", "a")
	s.Process("REM", "a")
	got, err := s.Query("EXISTS", "a")
	if err != nil {
		t.Fatalf("unexpected Query() error: %q", err)
	}
	if got != "false" {
		t.Errorf("Query(EXISTS a): expected false, got true")
	}
}

func TestDuplicateInSet(t *testing.T) {
	s := Init(NewORSet(), setHandlers, setQueries)
	s.Process("ADD", "a")
	s.Process("ADD", "a")
	s.Process("ADD", "a")
	s.Process("REM", "a")
	got, err := s.Query("EXISTS", "a")
	if err != nil {
		t.Fatalf("unexpected Query() error: %q", err)
	}
	if got != "false" {
		t.Errorf("Query(EXISTS a): expected false, got true")
	}
}

func TestRemoveBeforeAdd(t *testing.T) {
	s := Init(NewORSet(), setHandlers, setQueries)
	s.Process("REM", "a")
	s.Process("ADD", "a")
	got, err := s.Query("EXISTS", "a")
	if err != nil {
		t.Fatalf("unexpected Query() error: %q", err)
	}
	if got != "true" {
		t.Errorf("Query(EXISTS a): expected true, got false")
	}
}
