package main

import (
	"container/list"
	"fmt"
)

type StackFrame struct {
	ident string
	dp    int
}

type Interpreter struct {
	dataStack []interface{}
	callStack []*StackFrame
	dp        int
	cp        int
	ops       *list.List
	op        *list.Element
}

func NewInterpreter(ops *list.List) *Interpreter {
	dataStack := make([]interface{}, 1<<10)
	callStack := make([]*StackFrame, 1<<10)
	dp := -1
	cp := 0
	callStack[cp] = &StackFrame{"top", dp}
	return &Interpreter{dataStack, callStack, dp, cp, ops, ops.Front()}
}

func (i *Interpreter) Exec() (err error) {
	for i.op != nil {
		switch i.op.Value {
		case OpPush:
			i.push()
		case OpSwap:
			i.swap()
		case OpDup:
			i.dup()
		case OpOver:
			i.over()
		case OpRot:
			i.rot()
		case OpDrop:
			i.drop()
		case OpGet:
			i.get()
		case OpSet:
			i.set()
		case OpAddI:
			i.addI()
		case OpSubI:
			i.subI()
		case OpMulI:
			i.mulI()
		case OpDivI:
			i.divI()
		case OpAddF:
			i.addF()
		case OpSubF:
			i.subF()
		case OpMulF:
			i.mulF()
		case OpDivF:
			i.divF()
		case OpCall:
			i.call()
		default:
			panic(fmt.Errorf("Unexpected op: %d", i.op.Value.(int)))
		}
		i.nextOp()
	}
	return
}

func (i *Interpreter) Push(val interface{}) {
	i.dp++
	if i.dp >= len(i.dataStack) {
		i.growDataStack()
	}
	i.dataStack[i.dp] = val
}

func (i *Interpreter) Pop() (val interface{}) {
	if i.dp < 0 {
		return nil
	}
	val = i.dataStack[i.dp]
	i.dp--
	return
}

func (i *Interpreter) growDataStack() {
	o := len(i.dataStack)
	n := o * 2
	tmp := i.dataStack
	i.dataStack = make([]interface{}, n)
	copy(i.dataStack, tmp)
}

func (i *Interpreter) growCallStack() {
	o := len(i.callStack)
	n := o * 2
	tmp := i.callStack
	i.callStack = make([]*StackFrame, n)
	copy(i.callStack, tmp)
}

func (i *Interpreter) offsetToAddr(offset int) (addr int) {
	if offset >= 0 {
		addr = i.callStack[i.cp].dp + offset + 1
	} else {
		addr = i.dp + offset
	}
	return
}

func (i *Interpreter) stackFrame() (frame *StackFrame) {
	return i.callStack[i.cp]
}

func (i *Interpreter) push() {
	i.Push(i.nextOp())
}

func (i *Interpreter) drop() {
	i.Pop()
}

func (i *Interpreter) get() {
	offset := i.nextOp().(int)
	addr := i.offsetToAddr(offset)
	i.Push(i.dataStack[addr])
}

func (i *Interpreter) set() {
	offset := i.nextOp().(int)
	val := i.nextOp()
	addr := i.offsetToAddr(offset)
	i.dataStack[addr] = val
}

func (i *Interpreter) swap() {
	i.dataStack[i.dp], i.dataStack[i.dp-1] = i.dataStack[i.dp-1], i.dataStack[i.dp]
}

func (i *Interpreter) dup() {
	i.Push(i.dataStack[i.dp])
}

func (i *Interpreter) over() {
	i.Push(i.dataStack[i.dp-1])
}

func (i *Interpreter) rot() {
	i.dataStack[i.dp], i.dataStack[i.dp-1], i.dataStack[i.dp-2] =
		i.dataStack[i.dp-1], i.dataStack[i.dp-2], i.dataStack[i.dp]
}

func (i *Interpreter) addI() {
	x := i.dataStack[i.dp-1].(int) + i.dataStack[i.dp].(int)
	i.dp--
	i.dataStack[i.dp] = x
}

func (i *Interpreter) subI() {
	x := i.dataStack[i.dp-1].(int) - i.dataStack[i.dp].(int)
	i.dp--
	i.dataStack[i.dp] = x
}

func (i *Interpreter) mulI() {
	x := i.dataStack[i.dp-1].(int) * i.dataStack[i.dp].(int)
	i.dp--
	i.dataStack[i.dp] = x
}

func (i *Interpreter) divI() {
	x := i.dataStack[i.dp-1].(int) / i.dataStack[i.dp].(int)
	i.dp--
	i.dataStack[i.dp] = x
}

func (i *Interpreter) addF() {
	x := i.dataStack[i.dp-1].(float64) + i.dataStack[i.dp].(float64)
	i.dp--
	i.dataStack[i.dp] = x
}

func (i *Interpreter) subF() {
	x := i.dataStack[i.dp-1].(float64) - i.dataStack[i.dp].(float64)
	i.dp--
	i.dataStack[i.dp] = x
}

func (i *Interpreter) mulF() {
	x := i.dataStack[i.dp-1].(float64) * i.dataStack[i.dp].(float64)
	i.dp--
	i.dataStack[i.dp] = x
}

func (i *Interpreter) divF() {
	x := i.dataStack[i.dp-1].(float64) / i.dataStack[i.dp].(float64)
	i.dp--
	i.dataStack[i.dp] = x
}

func (i *Interpreter) call() {
	argc := i.dataStack[i.dp].(int)
	closure := i.dataStack[i.dp-argc-1].(func())
	ident := i.dataStack[i.dp-argc-2].(string)
	frame := &StackFrame{ident, i.dp}
	i.cp++
	if i.cp >= len(i.callStack) {
		i.growCallStack()
	}
	i.callStack[i.cp] = frame
	closure()
	i.cp--
}

func (i *Interpreter) nextOp() (op interface{}) {
	i.op = i.op.Next()
	if i.op != nil {
		op = i.op.Value
	}
	return
}
