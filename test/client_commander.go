package main

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"testing"
)

type Client struct {
	addr string
	conn net.Conn
}

func NewClient(addr string) *Client {
	return &Client{addr: addr}
}
func TestConnect(t *testing.T) {
	var client Client
	// ここだけ変更
	client.TestConn(":5000")

	t.Errorf("actual %v want %v", 1, 2)
}

func (c *Client) send(msg string) (int, error) {
	io.Copy(c.conn, bytes.NewBufferString(msg+"\n"))
	return 1, nil
}

func (c *Client) TestConn(addr string) (string, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return "", err
	}
	defer conn.Close()
	c.conn = conn

	_, err = c.send("COMMANDER")
	if err != nil {
		return "", err
	}

	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		return "", err
	}

	fmt.Sprintln(n)
	return "", err
}
