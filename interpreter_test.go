package main

import (
	"fmt"
	"strings"
	"testing"
)

func TestInterpreter(t *testing.T) {
	p := NewParser(strings.NewReader("1 + 2 * 3 - 4 / 2"))
	if err := p.Parse(); err != nil {
		t.Fatal(err)
	}
	i := NewInterpreter(256, p.Ops())
	if err := i.Exec(); err != nil {
		t.Fatal(err)
	}
	fmt.Println(i.Pop())
}
