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
	i := NewInterpreter(p.Ops())
	if err := i.Exec(); err != nil {
		t.Fatal(err)
	}
	assertEqualInt(t, 5, i.Pop().(int))
	defer func() {
		if err := recover(); err == nil {
			t.Fatal("expected runtime error")
		}
	}()
	// test empty stack
	i.Pop()
}

func TestIntExprListSemiColon(t *testing.T) {
	p := NewParser(strings.NewReader("1 + 2 * 3 - 4 / 2 ; 4 * 5"))
	if err := p.Parse(); err != nil {
		t.Fatal(err)
	}
	i := NewInterpreter(p.Ops())
	if err := i.Exec(); err != nil {
		t.Fatal(err)
	}
	assertEqualInt(t, 20, i.Pop().(int))
	defer func() {
		if err := recover(); err == nil {
			t.Fatal("expected runtime error")
		}
	}()
	// test empty stack
	i.Pop()
}

func TestIntExprListNewLine(t *testing.T) {
	p := NewParser(strings.NewReader("1 + 2 * 3 - 4 / 2 \n  4 * 5"))
	if err := p.Parse(); err != nil {
		t.Fatal(err)
	}
	i := NewInterpreter(p.Ops())
	if err := i.Exec(); err != nil {
		t.Fatal(err)
	}
	assertEqualInt(t, 20, i.Pop().(int))
	defer func() {
		if err := recover(); err == nil {
			t.Fatal("expected runtime error")
		}
	}()
	// test empty stack
	i.Pop()
}

func TestVarAssign(t *testing.T) {
	p := NewParser(strings.NewReader("x = 7"))
	if err := p.Parse(); err != nil {
		t.Fatal(err)
	}
	i := NewInterpreter(p.Ops())
	if err := i.Exec(); err != nil {
		t.Fatal(err)
	}
	assertEqualInt(t, 7, i.Pop().(int))
	assertEqualInt(t, 7, i.Pop().(int))
	defer func() {
		if err := recover(); err == nil {
			t.Fatal("expected runtime error")
		}
	}()
	// test empty stack
	i.Pop()
}

func TestVarEval(t *testing.T) {
	p := NewParser(strings.NewReader("x = 7; x + 2; y = 8; x + y"))
	if err := p.Parse(); err != nil {
		t.Fatal(err)
	}
	i := NewInterpreter(p.Ops())
	if err := i.Exec(); err != nil {
		t.Fatal(err)
	}
	assertEqualInt(t, 7, i.Pop().(int))
	assertEqualInt(t, 8, i.Pop().(int))
	assertEqualInt(t, 7, i.Pop().(int))
	defer func() {
		if err := recover(); err == nil {
			t.Fatal("expected runtime error")
		}
	}()
	// test empty stack
	i.Pop()
}

func assertEqualInt(t *testing.T, expected int, actual int) {
	if actual != expected {
		t.Fatal(fmt.Sprintf("expected %d, actual %d", expected, actual))
	}
}
