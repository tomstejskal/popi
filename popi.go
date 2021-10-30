package main

import (
	"fmt"
	"strings"
)

func main() {
	p := NewParser(strings.NewReader("x = 7; y = 8; x + y"))
	if err := p.Parse(); err != nil {
		fmt.Println(err)
		return
	}
	i := NewInterpreter(p.Ops())
	if err := i.Exec(); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("stack:")
	for x := i.Pop(); x != nil; x = i.Pop() {
		fmt.Println(x)
	}
}
