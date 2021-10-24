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
	TokNum
	TokAdd
	TokSub
	TokMul
	TokDiv
	TokLParen
	TokRParen
)

type Token int

type Lexer struct {
	rs   io.RuneScanner
	line int
	pos  int
	tok  Token
	val  interface{}
}

func NewLexer(rs io.RuneScanner) *Lexer {
	return &Lexer{rs, 1, 1, TokNone, nil}
}

func (l *Lexer) ReadToken() (tok Token, val interface{}, err error) {
	if l.tok > TokNone {
		tok, l.tok, val, l.val = l.tok, TokNone, l.val, nil
		return
	}
	if err = l.skipSpace(); err != nil {
		return
	}
	r, err := l.readRune()
	if err != nil {
		return
	}
	switch r {
	case '+':
		tok = TokAdd
	case '-':
		tok = TokSub
	case '*':
		tok = TokMul
	case '/':
		tok = TokDiv
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		if err = l.rs.UnreadRune(); err != nil {
			return
		}
		tok = TokNum
		if val, err = l.readNum(); err != nil {
			return
		}
	case '(':
		tok = TokLParen
	case ')':
		tok = TokRParen
	default:
		err = l.unexpectedChar(r)
	}

	return
}

func (l *Lexer) UnreadToken(tok Token, val interface{}) (err error) {
	if l.tok > TokNone {
		return l.makeError("Cannot unread token")
	}
	l.tok, l.val = tok, val
	return
}

func (l *Lexer) skipSpace() (err error) {
	var r rune
	for {
		r, err = l.readRune()
		if err != nil {
			return
		}
		if !unicode.IsSpace(r) {
			err = l.rs.UnreadRune()
			return
		}
	}
}

func (l *Lexer) readRune() (r rune, err error) {
	for {
		if r, _, err = l.rs.ReadRune(); err != nil {
			return
		}
		switch r {
		case '\r':
			if r, _, err = l.rs.ReadRune(); err != nil {
				return
			}
			if r == '\n' {
				l.line++
				l.pos = 1
			} else {
				l.pos++
				return
			}
		case '\n':
			l.line++
			l.pos = 1
		default:
			return
		}
	}
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

func (l *Lexer) makeError(format string, a ...interface{}) error {
	return &LexerError{l.line, l.pos, fmt.Sprintf(format, a...)}
}

func (l *Lexer) unexpectedChar(r rune) error {
	return l.makeError("Unexpected char %c at line %d and pos %d", r, l.line, l.pos)
}
