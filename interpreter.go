package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
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
	code      *bytes.Buffer
}

func NewInterpreter(code *bytes.Buffer) *Interpreter {
	dataStack := make([]interface{}, 1<<10)
	callStack := make([]*StackFrame, 1<<10)
	dp := -1
	cp := 0
	callStack[cp] = &StackFrame{"top", dp}
	return &Interpreter{dataStack, callStack, dp, cp, code}
}

func (i *Interpreter) Exec() (err error) {
	var op byte
	for op, err = i.code.ReadByte(); err == nil; op, err = i.code.ReadByte() {
		switch op {
		case OpPushI:
			if err = i.pushI(); err != nil {
				return
			}
		case OpPushF:
			if err = i.pushF(); err != nil {
				return
			}
		case OpSwap:
			if err = i.swap(); err != nil {
				return
			}
		case OpDup:
			if err = i.dup(); err != nil {
				return
			}
		case OpOver:
			if err = i.over(); err != nil {
				return
			}
		case OpRot:
			if err = i.rot(); err != nil {
				return
			}
		case OpDrop:
			if err = i.drop(); err != nil {
				return
			}
		case OpGet:
			if err = i.get(); err != nil {
				return
			}
		case OpSetI:
			if err = i.setI(); err != nil {
				return
			}
		case OpSetF:
			if err = i.setF(); err != nil {
				return
			}
		case OpAddI:
			if err = i.addI(); err != nil {
				return
			}
		case OpSubI:
			if err = i.subI(); err != nil {
				return
			}
		case OpMulI:
			if err = i.mulI(); err != nil {
				return
			}
		case OpDivI:
			if err = i.divI(); err != nil {
				return
			}
		case OpAddF:
			if err = i.addF(); err != nil {
				return
			}
		case OpSubF:
			if err = i.subF(); err != nil {
				return
			}
		case OpMulF:
			if err = i.mulF(); err != nil {
				return
			}
		case OpDivF:
			if err = i.divF(); err != nil {
				return
			}
		case OpCall:
			if err = i.call(); err != nil {
				return
			}
		default:
			panic(fmt.Errorf("Unexpected op: %d", op))
		}
	}
	if err == io.EOF {
		err = nil
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

func (i *Interpreter) pushI() (err error) {
	var n int
	if n, err = i.readInt(); err != nil {
		return err
	}
	i.Push(n)
	return
}

func (i *Interpreter) pushF() (err error) {
	var n float64
	if n, err = i.readFloat(); err != nil {
		return err
	}
	i.Push(n)
	return
}

func (i *Interpreter) drop() (err error) {
	i.Pop()
	return
}

func (i *Interpreter) get() (err error) {
	var offset int
	if offset, err = i.readInt(); err != nil {
		return
	}
	addr := i.offsetToAddr(offset)
	i.Push(i.dataStack[addr])
	return
}

func (i *Interpreter) setI() (err error) {
	var (
		offset int
		n      int
	)
	if offset, err = i.readInt(); err != nil {
		return
	}
	if n, err = i.readInt(); err != nil {
		return
	}
	addr := i.offsetToAddr(offset)
	i.dataStack[addr] = n
	return
}

func (i *Interpreter) setF() (err error) {
	var (
		offset int
		n      float64
	)
	if offset, err = i.readInt(); err != nil {
		return
	}
	if n, err = i.readFloat(); err != nil {
		return
	}
	addr := i.offsetToAddr(offset)
	i.dataStack[addr] = n
	return
}

func (i *Interpreter) swap() (err error) {
	i.dataStack[i.dp], i.dataStack[i.dp-1] = i.dataStack[i.dp-1], i.dataStack[i.dp]
	return
}

func (i *Interpreter) dup() (err error) {
	i.Push(i.dataStack[i.dp])
	return
}

func (i *Interpreter) over() (err error) {
	i.Push(i.dataStack[i.dp-1])
	return
}

func (i *Interpreter) rot() (err error) {
	i.dataStack[i.dp], i.dataStack[i.dp-1], i.dataStack[i.dp-2] =
		i.dataStack[i.dp-1], i.dataStack[i.dp-2], i.dataStack[i.dp]
	return
}

func (i *Interpreter) addI() (err error) {
	y := i.Pop().(int)
	x := i.Pop().(int)
	i.Push(x + y)
	return
}

func (i *Interpreter) subI() (err error) {
	y := i.Pop().(int)
	x := i.Pop().(int)
	i.Push(x - y)
	return
}

func (i *Interpreter) mulI() (err error) {
	y := i.Pop().(int)
	x := i.Pop().(int)
	i.Push(x * y)
	return
}

func (i *Interpreter) divI() (err error) {
	y := i.Pop().(int)
	x := i.Pop().(int)
	i.Push(x / y)
	return
}

func (i *Interpreter) addF() (err error) {
	y := i.Pop().(float32)
	x := i.Pop().(float32)
	i.Push(x + y)
	return
}

func (i *Interpreter) subF() (err error) {
	y := i.Pop().(float32)
	x := i.Pop().(float32)
	i.Push(x - y)
	return
}

func (i *Interpreter) mulF() (err error) {
	y := i.Pop().(float32)
	x := i.Pop().(float32)
	i.Push(x * y)
	return
}

func (i *Interpreter) divF() (err error) {
	y := i.Pop().(float32)
	x := i.Pop().(float32)
	i.Push(x / y)
	return
}

func (i *Interpreter) call() (err error) {
	return fmt.Errorf("Not implemented")
}

func (i *Interpreter) readInt() (n int, err error) {
	var x int64
	err = binary.Read(i.code, binary.LittleEndian, &x)
	n = int(x)
	return
}

func (i *Interpreter) readFloat() (n float64, err error) {
	err = binary.Read(i.code, binary.LittleEndian, n)
	return
}
