package main

import "io"

type ByteCode struct {
	buf  []byte
	addr int
}

func NewByteCode(buf []byte) *ByteCode {
	return &ByteCode{buf, 0}
}

func (b *ByteCode) Addr() int {
	return b.addr
}

func (b *ByteCode) SetAddr(addr int) {
	b.addr = addr
}

func (b *ByteCode) Read(p []byte) (n int, err error) {
	n = len(b.buf) - b.addr
	if n == 0 {
		err = io.EOF
		return
	}
	if len(p) < n {
		n = len(p)
	}
	copy(p, b.buf[b.addr:b.addr+n])
	b.addr += n
	return
}

func (b *ByteCode) ReadByte() (c byte, err error) {
	n := len(b.buf) - b.addr
	if n == 0 {
		err = io.EOF
		return
	}
	c = b.buf[b.addr]
	b.addr++
	return
}
