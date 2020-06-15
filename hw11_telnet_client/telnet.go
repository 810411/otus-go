package main

import (
	"bufio"
	"errors"
	"io"
	"net"
	"time"
)

var (
	ErrConnectionClosed = errors.New("connection closed by peer")
	ErrEOF              = errors.New("...EOF")
)

type TelnetClient interface {
	Connect() error
	Close() error
	Send() error
	Receive() error
}

type telnet struct {
	address  string
	timeout  time.Duration
	in       io.ReadCloser
	out      io.Writer
	conn     net.Conn
	inScan   *bufio.Scanner
	connScan *bufio.Scanner
}

func (t *telnet) Connect() (err error) {
	t.conn, err = net.DialTimeout("tcp", t.address, t.timeout)
	t.connScan = bufio.NewScanner(t.conn)
	t.inScan = bufio.NewScanner(t.in)
	return
}

func (t *telnet) Close() (err error) {
	if t.conn != nil {
		err = t.conn.Close()
	}
	return
}

func (t *telnet) Send() (err error) {
	if t.conn == nil {
		return
	}
	if !t.inScan.Scan() {
		return ErrEOF
	}
	_, err = t.conn.Write(append(t.inScan.Bytes(), '\n'))
	return
}

func (t *telnet) Receive() (err error) {
	if t.conn == nil {
		return
	}
	if !t.connScan.Scan() {
		return ErrConnectionClosed
	}
	_, err = t.out.Write(append(t.connScan.Bytes(), '\n'))
	return
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &telnet{
		address: address,
		timeout: timeout,
		in:      in,
		out:     out,
	}
}
