package main

import (
	"fmt"
	"log"
	"strings"
)

func main() {
	p := NewParser(strings.NewReader("x = 7; x + 1"))
	if err := p.Parse(); err != nil {
		log.Fatal(err)
	}
	i := NewInterpreter(p.Ops())
	if err := i.Exec(); err != nil {
		log.Fatal(err)
	}
	fmt.Println(i.Pop())
}
