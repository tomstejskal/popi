package main

import (
	"fmt"
	"strings"
	"testing"
)

func TestIntExpr(t *testing.T) {
	i := exec(t, "1 + 2 * 3 - 4 / 2")
	checkEqualInt(t, 5, i.Pop().(int))
	checkNil(t, i.Pop())
}

func TestIntExprListSemiColon(t *testing.T) {
	i := exec(t, "1 + 2 * 3 - 4 / 2 ; 4 * 5")
	checkEqualInt(t, 20, i.Pop().(int))
	checkNil(t, i.Pop())
}

func TestIntExprListNewLine(t *testing.T) {
	i := exec(t, "1 + 2 * 3 - 4 / 2 \n  4 * 5")
	checkEqualInt(t, 20, i.Pop().(int))
	checkNil(t, i.Pop())
}

func TestVarAssign(t *testing.T) {
	i := exec(t, "x = 7")
	checkEqualInt(t, 7, i.Pop().(int))
	checkEqualInt(t, 7, i.Pop().(int))
	checkNil(t, i.Pop())
}

func TestVarEval(t *testing.T) {
	i := exec(t, "x = 7; x + 2; y = 8; x + y")
	checkEqualInt(t, 15, i.Pop().(int))
	checkEqualInt(t, 8, i.Pop().(int))
	checkEqualInt(t, 7, i.Pop().(int))
	checkNil(t, i.Pop())
}

func checkEqualInt(t *testing.T, expected int, actual int) {
	if actual != expected {
		t.Fatal(fmt.Sprintf("%s: expected %d, actual %d", t.Name(), expected, actual))
	}
}

func checkNil(t *testing.T, val interface{}) {
	if val != nil {
		t.Fatal(fmt.Sprintf("%s: expected nil, actual %v", t.Name(), val))
	}
}

func exec(t *testing.T, s string) (i *Interpreter) {
	p := NewParser(strings.NewReader(s))
	code, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}
	i = NewInterpreter(code)
	if err := i.Exec(); err != nil {
		t.Fatal(err)
	}
	return
}
