package main

import (
	"fmt"
	"strings"
)

func main() {
	p := NewParser(strings.NewReader("x = 7; x + 2; y = 8; x + y"))
	code, err := p.Parse()
	if err != nil {
		fmt.Println(err)
		return
	}
	i := NewInterpreter(code)
	if err := i.Exec(); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("stack:")
	for x := i.Pop(); x != nil; x = i.Pop() {
		fmt.Println(x)
	}
}
