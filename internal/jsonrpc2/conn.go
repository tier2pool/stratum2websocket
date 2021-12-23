package jsonrpc2

import (
	"bufio"
	"encoding/json"
	"errors"
	"github.com/tier2pool/tier2pool/internal/pool"
	"io"
	"net"
	"time"
)

type Conn interface {
	net.Conn
	io.ReadWriteCloser

	WriteStruct(r interface{}) error
	ReadStruct(r interface{}) error
}

func NewConn(c net.Conn) Conn {
	return &conn{
		conn:   c,
		reader: bufio.NewReader(c),
	}
}

type conn struct {
	conn   net.Conn
	reader *bufio.Reader
}

func (c *conn) LocalAddr() net.Addr {
	return c.conn.LocalAddr()
}

func (c *conn) RemoteAddr() net.Addr {
	return c.conn.RemoteAddr()
}

func (c *conn) SetDeadline(t time.Time) error {
	return c.conn.SetDeadline(t)
}

func (c *conn) SetReadDeadline(t time.Time) error {
	return c.conn.SetReadDeadline(t)
}

func (c *conn) SetWriteDeadline(t time.Time) error {
	return c.conn.SetWriteDeadline(t)
}

func (c *conn) Read(p []byte) (n int, err error) {
	line, isPrefix, err := c.reader.ReadLine()
	if err != nil {
		return 0, err
	}
	if isPrefix {
		return copy(p, line), errors.New("data too long")
	}
	return copy(p, line), nil
}

func (c *conn) Write(p []byte) (n int, err error) {
	return c.conn.Write(append(p, []byte{'\n'}...))
}

func (c *conn) Close() error {
	return c.conn.Close()
}

func (c *conn) WriteStruct(r interface{}) error {
	data, err := json.Marshal(r)
	if err != nil {
		return err
	}
	_, err = c.Write(data)
	return err
}

func (c *conn) ReadStruct(r interface{}) error {
	data := pool.Get()
	defer pool.Put(data)
	n, err := c.Read(data)
	if err != nil {
		return err
	}
	return json.Unmarshal(data[:n], r)
}
