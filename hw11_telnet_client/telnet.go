package main

import (
	"io"
	"log"
	"net"
	"time"
)

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &client{
		address: address,
		timeout: timeout,
		in:      in,
		out:     out,
	}
}

type client struct {
	address    string
	timeout    time.Duration
	in         io.ReadCloser
	out        io.Writer
	connection net.Conn
}

func (c *client) Connect() error {
	conn, err := net.DialTimeout("tcp", c.address, c.timeout)
	if err != nil {
		return err
	}

	c.connection = conn
	log.Printf("...Connected to %s", c.address)

	return nil
}

func (c *client) Send() error {
	_, err := io.Copy(c.connection, c.in)
	return err
}

func (c *client) Receive() error {
	_, err := io.Copy(c.out, c.connection)
	return err
}

func (c *client) Close() error {
	if c.connection != nil {
		err := c.connection.Close()
		if err != nil {
			return err
		}
		c.connection = nil
	}

	return nil
}
