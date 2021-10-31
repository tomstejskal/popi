package main

import "fmt"

const (
	TokNone = iota
	TokInt
	TokFloat
	TokAdd
	TokSub
	TokMul
	TokDiv
	TokLParen
	TokRParen
	TokLBrace
	TokRBrace
	TokComma
	TokSColon
	TokFunc
	TokIdent
	TokAssign
	TokEqual
	TokEOF
)

type Token struct {
	id  int
	val interface{}
}

func (t Token) String() string {
	switch t.id {
	case TokNone:
		return "None"
	case TokInt:
		return fmt.Sprintf("%d", t.val.(int))
	case TokFloat:
		return fmt.Sprintf("%f", t.val.(float64))
	case TokAdd:
		return "+"
	case TokSub:
		return "-"
	case TokMul:
		return "*"
	case TokDiv:
		return "/"
	case TokLParen:
		return "("
	case TokRParen:
		return ")"
	case TokLBrace:
		return "{"
	case TokRBrace:
		return "}"
	case TokComma:
		return ","
	case TokSColon:
		return ";"
	case TokFunc:
		return "func"
	case TokIdent:
		return t.val.(string)
	case TokAssign:
		return "="
	case TokEqual:
		return "=="
	case TokEOF:
		return "EOF"
	default:
		panic(fmt.Errorf("Unknown token: %d", t.id))
	}
}
