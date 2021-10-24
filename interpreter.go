package main

import (
	"container/list"
)

type Interpreter struct {
	stack []interface{}
	sp    int
	ops   *list.List
}

func NewInterpreter(stackSize int, ops *list.List) *Interpreter {
	stack := make([]interface{}, stackSize)
	return &Interpreter{stack, 0, ops}
}

func (i *Interpreter) Push(val interface{}) {
	i.stack[i.sp] = val
	i.sp++
}

func (i *Interpreter) Pop() interface{} {
	i.sp--
	return i.stack[i.sp]
}

func (i *Interpreter) Peek() interface{} {
	return i.stack[i.sp-1]
}

func (i *Interpreter) Exec() (err error) {
	el := i.ops.Front()
	for el != nil {
		if el == nil {
			return
		}
		switch el.Value {
		case OpPush:
			el = el.Next()
			i.Push(el.Value)
		case OpPop:
			i.Pop()
		case OpAdd:
			y, x := i.Pop().(int), i.Pop().(int)
			i.Push(x + y)
		case OpSub:
			y, x := i.Pop().(int), i.Pop().(int)
			i.Push(x - y)
		case OpMul:
			y, x := i.Pop().(int), i.Pop().(int)
			i.Push(x * y)
		case OpDiv:
			y, x := i.Pop().(int), i.Pop().(int)
			i.Push(x / y)
		}
		el = el.Next()
	}
	return
}
