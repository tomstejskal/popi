package main

import (
	"strings"
	"testing"
)

func TestIntExpr(t *testing.T) {
	p := NewParser(strings.NewReader("1 + 2 * 3 - 4 / 2"))
	if err := p.Parse(); err != nil {
		t.Fatal(err)
	}
	i := NewInterpreter(256, p.Ops())
	if err := i.Exec(); err != nil {
		t.Fatal(err)
	}
	if i.Pop() != 5 {
		t.Fatal("expected result 5")
	}
	defer func() {
		if err := recover(); err == nil {
			t.Fatal("expected runtime error")
		}
	}()
	// test empty stack
	i.Pop()
}
