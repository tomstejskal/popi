package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

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

type ParserError struct {
	line int
	pos  int
	msg  string
}

func (err *ParserError) Error() string {
	return fmt.Sprintf("%s at line %d and position %d", err.msg, err.line, err.pos)
}

func (err *ParserError) Line() int {
	return err.line
}

func (err *ParserError) Pos() int {
	return err.pos
}

func (err *ParserError) Msg() string {
	return err.msg
}

type Parser struct {
	lex   *Lexer
	code  *bytes.Buffer
	tok   Token
	line  int
	pos   int
	scope *Scope
}

type Scope struct {
	item      *Item
	stackSize int
	next      *Scope
}

const (
	ItemVar = iota
	ItemParam
)

type Item struct {
	typ   byte
	ident string
	val   int
	next  *Item
}

func NewParser(rs io.RuneScanner) *Parser {
	code := &bytes.Buffer{}
	return &Parser{NewLexer(rs), code, Token{}, 1, 1, &Scope{}}
}

func (p *Parser) Parse() (code *bytes.Buffer, err error) {
	err = p.readExprList()
	if err == nil {
		code = p.code
	}
	return
}

func (p *Parser) readToken() (err error) {
	p.line, p.pos = p.lex.line, p.lex.pos
	p.tok, err = p.lex.ReadToken()
	return err
}

func (p *Parser) unreadToken() (err error) {
	return p.lex.UnreadToken(p.tok)
}

func (p *Parser) readExprList() (err error) {
	if err = p.readExpr(); err != nil {
		return
	}
	for {
		if err = p.readToken(); err != nil {
			return
		}
		if p.tok.id == TokEOF {
			return
		}
		if p.tok.id != TokSColon {
			return p.unexpectedToken(";")
		}
		for {
			if err = p.readToken(); err != nil {
				return
			}
			if p.tok.id != TokSColon {
				if err = p.unreadToken(); err != nil {
					return
				}
				break
			}
		}
		if err = p.readToken(); err == nil {
			if err = p.unreadToken(); err != nil {
				return
			}
			p.writeOp(OpDrop)
		}
		if err = p.readExpr(); err != nil {
			return
		}
	}
}

func (p *Parser) readExpr() (err error) {
	return p.readTerm()
}

func (p *Parser) readTerm() (err error) {
	if err = p.readFactor(); err != nil {
		return
	}
	for {
		if err = p.readToken(); err != nil {
			return
		}
		var op byte
		switch p.tok.id {
		case TokAdd:
			op = OpAddI
		case TokSub:
			op = OpSubI
		default:
			return p.unreadToken()
		}
		if err = p.readFactor(); err != nil {
			return
		}
		p.writeOp(op)
	}
}

func (p *Parser) readFactor() (err error) {
	if err = p.readVal(); err != nil {
		return
	}
	for {
		if err = p.readToken(); err != nil {
			return
		}
		var op byte
		switch p.tok.id {
		case TokMul:
			op = OpMulI
		case TokDiv:
			op = OpDivI
		default:
			return p.unreadToken()
		}
		if err = p.readVal(); err != nil {
			return
		}
		p.writeOp(op)
	}
}

func (p *Parser) readVal() (err error) {
	if err = p.readToken(); err != nil {
		return
	}
	switch p.tok.id {
	case TokInt:
		p.writeOp(OpPushI)
		p.writeInt(p.tok.val.(int))
	case TokFunc:
		return p.readFunc()
	case TokIdent:
		ident := p.tok.val.(string)
		if err = p.readToken(); err != nil {
			return
		}
		if p.tok.id == TokAssign {
			// variable assignment
			if err = p.readExpr(); err != nil {
				return
			}
			p.pushItem(p.newVar(ident))
			p.writeOp(OpDup)
		} else {
			if err = p.unreadToken(); err != nil {
				return
			}
			// variable evaluation
			scope := p.scope
			for scope != nil {
				item := scope.item
				for item != nil {
					if item.ident == ident {
						p.writeOp(OpGet)
						p.writeInt(item.val)
						return
					}
					item = item.next
				}
				scope = scope.next
			}
			return fmt.Errorf("Unknown variable: %s", ident)
		}
	default:
		return p.unexpectedToken("value")
	}
	return
}

func (p *Parser) readFunc() (err error) {
	p.pushScope(&Scope{})
	defer p.popScope()
	p.writeOp(OpPushI)
	p.writeInt(p.code.Len() - 1) // function address
	if err = p.readToken(); err != nil {
		return
	}
	if p.tok.id != TokLParen {
		return p.unexpectedToken("(")
	}
	if err = p.readFuncParams(); err != nil {
		return
	}
	if err = p.readToken(); err != nil {
		return
	}
	if p.tok.id != TokRParen {
		return p.unexpectedToken(")")
	}
	if err = p.readFuncBody(); err != nil {
		return
	}
	p.writeOp(OpRet)
	return
}

func (p *Parser) readFuncParams() (err error) {
	if err = p.readToken(); err != nil {
		return
	}
	if p.tok.id != TokIdent {
		return p.unexpectedToken("ident")
	}
	ident := p.tok.val.(string)
	p.pushItem(p.newParam(ident, 0))
	for pos := 1; ; pos++ {
		if err = p.readToken(); err != nil {
			return
		}
		if p.tok.id != TokComma {
			return
		}
		if err = p.readToken(); err != nil {
			return
		}
		if p.tok.id != TokIdent {
			err = p.unexpectedToken("ident")
			return
		}
		ident = p.tok.val.(string)
		p.pushItem(p.newParam(ident, pos))
	}
}

func (p *Parser) readFuncBody() (err error) {
	if err = p.readToken(); err != nil {
		return
	}
	if p.tok.id != TokLBrace {
		err = p.unexpectedToken("{")
		return
	}
	if err = p.readExprList(); err != nil {
		return
	}
	if p.tok.id != TokRBrace {
		err = p.unexpectedToken("}")
		return
	}
	return
}

func (p *Parser) unexpectedToken(expected string) (err error) {
	return &ParserError{p.lex.line, p.lex.pos,
		fmt.Sprintf("Unexpected token: %s, %s expected", p.tok, expected)}
}

func (p *Parser) pushItem(item *Item) {
	item.next = p.scope.item
	p.scope.item = item
}

func (p *Parser) pushScope(scope *Scope) {
	scope.next = p.scope
	p.scope = scope
}

func (p *Parser) popScope() (scope *Scope) {
	scope = p.scope
	p.scope = scope.next
	scope.next = nil
	return
}

func (p *Parser) writeOp(op byte) (err error) {
	if err = p.code.WriteByte(op); err != nil {
		return
	}
	switch op {
	case OpPushI:
		p.scope.stackSize++
	case OpPushF:
		p.scope.stackSize++
	case OpSwap:
	case OpDup:
		p.scope.stackSize++
	case OpOver:
		p.scope.stackSize++
	case OpRot:
	case OpDrop:
		p.scope.stackSize--
	case OpGet:
		p.scope.stackSize++
	case OpSetI:
	case OpSetF:
	case OpAddI:
		p.scope.stackSize--
	case OpSubI:
		p.scope.stackSize--
	case OpMulI:
		p.scope.stackSize--
	case OpDivI:
		p.scope.stackSize--
	case OpAddF:
		p.scope.stackSize--
	case OpSubF:
		p.scope.stackSize--
	case OpMulF:
		p.scope.stackSize--
	case OpDivF:
		p.scope.stackSize--
	case OpRet:
	case OpCall:
	default:
		panic(fmt.Errorf("Unexpected op code: %d\n", op))
	}
	return
}

func (p *Parser) writeInt(n int) error {
	return binary.Write(p.code, binary.LittleEndian, int64(n))
}

func (p *Parser) writeFloat(n float64) error {
	return binary.Write(p.code, binary.LittleEndian, n)
}

func (p *Parser) newVar(ident string) *Item {
	return &Item{typ: ItemVar, ident: ident, val: p.scope.stackSize - 1}
}

func (p *Parser) newParam(ident string, pos int) *Item {
	return &Item{typ: ItemParam, ident: ident, val: pos}
}
