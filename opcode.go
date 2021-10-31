package main

import "fmt"

const (
	OpPushI = 1 + iota
	OpPushF
	OpSwap
	OpDup
	OpOver
	OpRot
	OpDrop
	OpGet
	OpSetI
	OpSetF
	OpAddI
	OpSubI
	OpMulI
	OpDivI
	OpAddF
	OpSubF
	OpMulF
	OpDivF
	OpRet
	OpCall
)

type OpCode byte

func (op OpCode) String() string {
	switch op {
	case OpPushI:
		return "pushi"
	case OpPushF:
		return "pushf"
	case OpSwap:
		return "swap"
	case OpDup:
		return "dup"
	case OpOver:
		return "over"
	case OpRot:
		return "rot"
	case OpDrop:
		return "drop"
	case OpGet:
		return "get"
	case OpSetI:
		return "seti"
	case OpSetF:
		return "setf"
	case OpAddI:
		return "addi"
	case OpSubI:
		return "subi"
	case OpMulI:
		return "muli"
	case OpDivI:
		return "divi"
	case OpAddF:
		return "addf"
	case OpSubF:
		return "subf"
	case OpMulF:
		return "mulf"
	case OpDivF:
		return "divf"
	case OpRet:
		return "ret"
	case OpCall:
		return "call"
	default:
		return fmt.Sprintf("%d", op)
	}
}
