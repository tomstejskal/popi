package main

import (
	"container/list"
	"fmt"
	"io"
)

const (
	OpPush = 1 + iota
	OpPop
	OpAdd
	OpSub
	OpMul
	OpDiv
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
	lex  *Lexer
	ops  *list.List
	tok  Token
	val  interface{}
	line int
	pos  int
}

func NewParser(rs io.RuneScanner) *Parser {
	return &Parser{NewLexer(rs), list.New(), 0, nil, 1, 1}
}

func (p *Parser) Parse() (err error) {
	return p.readExprList()
}

func (p *Parser) Ops() *list.List {
	return p.ops
}

func (p *Parser) readToken() (err error) {
	p.line, p.pos = p.lex.line, p.lex.pos
	p.tok, p.val, err = p.lex.ReadToken()
	return err
}

func (p *Parser) unreadToken() (err error) {
	return p.lex.UnreadToken(p.tok, p.val)
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
		if p.tok != TokSColon {
			return p.unexpectedToken("';'")
		}
		for {
			if err = p.readToken(); err != nil {
				if err == io.EOF {
					err = nil
				}
				return
			}
			if p.tok != TokSColon {
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
			p.addOp(OpPop)
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
		switch p.tok {
		case TokAdd:
			op = OpAdd
		case TokSub:
			op = OpSub
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
		switch p.tok {
		case TokMul:
			op = OpMul
		case TokDiv:
			op = OpDiv
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
	switch p.tok {
	case TokNum:
		p.addOp(OpPush)
		p.addOp(p.val)
	case TokFunc:
		return p.readFunc()
	default:
		return p.unexpectedToken("value")
	}
	return
}

func (p *Parser) readFunc() (err error) {
	if err = p.readToken(); err != nil {
		return
	}
	if p.tok != TokLParen {
		return p.unexpectedToken("'('")
	}
	var params []string
	if params, err = p.readFuncParams(); err != nil {
		fmt.Println(params)
		return
	}
	if err = p.readToken(); err != nil {
		return
	}
	if p.tok != TokRParen {
		return p.unexpectedToken("')'")
	}
	return
}

func (p *Parser) readFuncParams() (params []string, err error) {
	return
}

func (p *Parser) addOp(op interface{}) {
	p.ops.PushBack(op)
}

func (p *Parser) unexpectedToken(expected string) (err error) {
	return &ParserError{p.lex.line, p.lex.pos,
		fmt.Sprintf("Unexpected token: %d, %s expected", p.tok, expected)}
}
