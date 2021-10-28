package main

import (
	"container/list"
	"fmt"
	"io"
)

const (
	OpPush = 1 + iota
	OpSwap
	OpDup
	OpOver
	OpRot
	OpDrop
	OpGet
	OpSet
	OpAddI
	OpSubI
	OpMulI
	OpDivI
	OpAddF
	OpSubF
	OpMulF
	OpDivF
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
	lex    *Lexer
	ops    *list.List
	tok    Token
	line   int
	pos    int
	scopes *list.List
}

type Scope struct {
	items     *list.List
	stackSize int
}

type Item struct {
	ident string
	pos   int
}

func NewParser(rs io.RuneScanner) *Parser {
	scopes := list.New()
	scopes.PushFront(&Scope{list.New(), 0})
	return &Parser{NewLexer(rs), list.New(), Token{}, 1, 1, scopes}
}

func (p *Parser) Parse() (err error) {
	return p.readExprList()
}

func (p *Parser) Ops() *list.List {
	return p.ops
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
			if err == io.EOF {
				err = nil
			}
			return
		}
		if p.tok.id != TokSColon {
			return p.unexpectedToken(";")
		}
		for {
			if err = p.readToken(); err != nil {
				if err == io.EOF {
					err = nil
				}
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
			p.addOp(OpDrop)
		}
		if err = p.readExpr(); err != nil {
			if err == io.EOF {
				err = nil
			}
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
			if err == io.EOF {
				err = nil
			}
			return
		}
		var op int
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
		p.addOp(op)
	}
}

func (p *Parser) readFactor() (err error) {
	if err = p.readVal(); err != nil {
		return
	}
	for {
		if err = p.readToken(); err != nil {
			if err == io.EOF {
				err = nil
			}
			return
		}
		var op int
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
		p.addOp(op)
	}
}

func (p *Parser) readVal() (err error) {
	if err = p.readToken(); err != nil {
		return
	}
	switch p.tok.id {
	case TokInt:
		p.addOp(OpPush)
		p.addOp(p.tok.val)
		p.scope().stackSize++
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
			p.scope().items.PushBack(&Item{ident, p.scope().stackSize - 1})
			p.addOp(OpDup)
		} else {
			if err = p.unreadToken(); err != nil {
				return
			}
			// variable evaluation
			sEl := p.scopes.Front()
			for sEl != nil {
				s := sEl.Value.(*Scope)
				iEl := s.items.Front()
				for iEl != nil {
					i := iEl.Value.(*Item)
					if i.ident == ident {
						p.addOp(OpGet)
						p.addOp(i.pos)
						return
					}
					iEl = iEl.Next()
				}
				sEl = sEl.Next()
			}
			return fmt.Errorf("Unknown variable: %s", ident)
		}
	default:
		return p.unexpectedToken("value")
	}
	return
}

func (p *Parser) readFunc() (err error) {
	if err = p.readToken(); err != nil {
		return
	}
	if p.tok.id != TokLParen {
		return p.unexpectedToken("(")
	}
	var params []string
	if params, err = p.readFuncParams(); err != nil {
		return
	}
	fmt.Println(params)
	if err = p.readToken(); err != nil {
		return
	}
	if p.tok.id != TokRParen {
		return p.unexpectedToken(")")
	}
	var body func()
	if body, err = p.readFuncBody(); err != nil {
		return
	}
	p.addOp(body)
	return
}

func (p *Parser) readFuncParams() (params []string, err error) {
	if err = p.readToken(); err != nil {
		return
	}
	params = append(params, p.tok.val.(string))
	for {
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
			err = p.unexpectedToken("function parameter")
			return
		}
		params = append(params, p.tok.val.(string))
	}
}

func (p *Parser) readFuncBody() (body func(), err error) {
	if err = p.readToken(); err != nil {
		return
	}
	if p.tok.id != TokLBrace {
		err = p.unexpectedToken("{")
		return
	}
	return
}

func (p *Parser) addOp(op interface{}) {
	p.ops.PushBack(op)
}

func (p *Parser) unexpectedToken(expected string) (err error) {
	return &ParserError{p.lex.line, p.lex.pos,
		fmt.Sprintf("Unexpected token: %s, %s expected", p.tok, expected)}
}

func (p *Parser) scope() (s *Scope) {
	return p.scopes.Front().Value.(*Scope)
}
