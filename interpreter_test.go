package main

import (
	"fmt"
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
	n := i.Pop()
	if n != 5 {
		t.Fatal(fmt.Sprintf("expected %d, actual %d", 5, n))
	}
	defer func() {
		if err := recover(); err == nil {
			t.Fatal("expected runtime error")
		}
	}()
	// test empty stack
	i.Pop()
}

func TestIntExprList(t *testing.T) {
	p := NewParser(strings.NewReader("1 + 2 * 3 - 4 / 2; 4 * 5"))
	if err := p.Parse(); err != nil {
		t.Fatal(err)
	}
	i := NewInterpreter(256, p.Ops())
	if err := i.Exec(); err != nil {
		t.Fatal(err)
	}
	n := i.Pop()
	if n != 20 {
		t.Fatal(fmt.Sprintf("expected %d, actual %d", 20, n))
	}
	if n != 20 {
		t.Fatal(fmt.Sprintf("expected %d, actual %d", 20, n))
	}
	defer func() {
		if err := recover(); err == nil {
			t.Fatal("expected runtime error")
		}
	}()
	// test empty stack
	i.Pop()
}
