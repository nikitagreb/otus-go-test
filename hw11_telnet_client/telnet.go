package main

import (
	"io"
	"net"
	"time"
)

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

type client struct {
	connect net.Conn
	address string
	in      io.ReadCloser
	out     io.Writer
	timeout time.Duration
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &client{
		address: address,
		timeout: timeout,
		in:      in,
		out:     out,
	}
}

func (c *client) Connect() error {
	connect, err := net.DialTimeout("tcp", c.address, c.timeout)
	if err != nil {
		return err
	}
	c.connect = connect
	return nil
}

func (c *client) Send() error {
	if _, err := io.Copy(c.connect, c.in); err != nil {
		return err
	}
	return nil
}

func (c *client) Receive() error {
	if _, err := io.Copy(c.out, c.connect); err != nil {
		return err
	}
	return nil
}

func (c *client) Close() error {
	if err := c.connect.Close(); err != nil {
		return err
	}
	return nil
}
