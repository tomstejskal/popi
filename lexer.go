package main

import (
	"fmt"
	"io"
	"unicode"
)

type LexerError struct {
	line int
	pos  int
	msg  string
}

func (err *LexerError) Error() string {
	return fmt.Sprintf("%s at line %d and position %d", err.msg, err.line, err.pos)
}

func (err *LexerError) Line() int {
	return err.line
}

func (err *LexerError) Pos() int {
	return err.pos
}

func (err *LexerError) Msg() string {
	return err.msg
}

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

type Lexer struct {
	rs   io.RuneScanner
	line int
	pos  int
	tok  Token
}

func NewLexer(rs io.RuneScanner) *Lexer {
	return &Lexer{rs, 1, 1, Token{}}
}

func (l *Lexer) ReadToken() (tok Token, err error) {
	if l.tok.id > TokNone {
		tok, l.tok = l.tok, Token{}
		return
	}
	if err = l.skipSpace(); err != nil {
		if err == io.EOF {
			tok = Token{id: TokEOF}
			err = nil
		}
		return
	}
	r, err := l.readRune()
	if err != nil {
		if err == io.EOF {
			tok = Token{id: TokEOF}
			err = nil
		}
		return
	}
	switch r {
	case '+':
		tok = Token{id: TokAdd}
	case '-':
		tok = Token{id: TokSub}
	case '*':
		tok = Token{id: TokMul}
	case '/':
		tok = Token{id: TokDiv}
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		if err = l.rs.UnreadRune(); err != nil {
			return
		}
		tok = Token{id: TokInt}
		if tok.val, err = l.readNum(); err != nil {
			return
		}
	case '(':
		tok = Token{id: TokLParen}
	case ')':
		tok = Token{id: TokRParen}
	case '{':
		tok = Token{id: TokLBrace}
	case '}':
		tok = Token{id: TokRBrace}
	case ',':
		tok = Token{id: TokComma}
	case ';':
		tok = Token{id: TokSColon}
	case '\n':
		tok = Token{id: TokSColon} // implicit semicolon
	case '=':
		if r, err = l.readRune(); err != nil {
			return
		}
		if r == '=' {
			tok = Token{id: TokEqual}
		} else {
			tok = Token{id: TokAssign}
			if err = l.rs.UnreadRune(); err != nil {
				return
			}
		}
	default:
		if err = l.rs.UnreadRune(); err != nil {
			return
		}
		var val string
		if val, err = l.readIdent(); err != nil {
			return
		}
		if val == "func" {
			tok = Token{id: TokFunc, val: val}
		} else {
			tok = Token{id: TokIdent, val: val}
		}
	}

	return
}

func (l *Lexer) UnreadToken(tok Token) (err error) {
	if l.tok.id > TokNone {
		return l.makeError("Cannot unread token")
	}
	l.tok = tok
	return
}

func (l *Lexer) skipSpace() (err error) {
	var r rune
	for {
		r, err = l.readRune()
		if err != nil {
			return
		}
		if !unicode.IsSpace(r) || r == '\n' {
			return l.rs.UnreadRune()
		}
	}
}

func (l *Lexer) readRune() (r rune, err error) {
	if r, _, err = l.rs.ReadRune(); err != nil {
		return
	}
	if r == '\n' {
		l.line++
		l.pos = 1
	}
	return
}

func (l *Lexer) readNum() (val int, err error) {
	val = 0
	found := false
	for {
		var r rune
		r, err = l.readRune()
		if err != nil {
			if err == io.EOF && found {
				err = nil
			}
			return
		}
		if r < '0' || r > '9' {
			err = l.rs.UnreadRune()
			if !found {
				err = l.unexpectedChar(r)
			}
			return
		}
		val = val*10 + (int(r) - '0')
		found = true
	}
}

func (l *Lexer) readIdent() (val string, err error) {
	var r rune
	if r, err = l.readRune(); err != nil {
		return
	}
	if !unicode.IsLetter(r) {
		err = l.unexpectedChar(r)
		return
	}
	val = string(r)
	for {
		r, err = l.readRune()
		if err != nil {
			if err == io.EOF {
				err = nil
			}
			return
		}
		if !(unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_') {
			err = l.rs.UnreadRune()
			return
		}
		val += string(r)
	}
}

func (l *Lexer) makeError(format string, a ...interface{}) error {
	return &LexerError{l.line, l.pos, fmt.Sprintf(format, a...)}
}

func (l *Lexer) unexpectedChar(r rune) error {
	return l.makeError("Unexpected char %c at line %d and pos %d", r, l.line, l.pos)
}
