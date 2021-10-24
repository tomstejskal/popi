package main

import (
	"fmt"
	"log"
	"strings"
)

func main() {
	p := NewParser(strings.NewReader("1 + 2 * 3 - 4 / 2"))
	if err := p.Parse(); err != nil {
		log.Fatal(err)
	}
	i := NewInterpreter(256, p.Ops())
	if err := i.Exec(); err != nil {
		log.Fatal(err)
	}
	fmt.Println(i.Pop())
}
