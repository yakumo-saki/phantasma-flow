package main

import (
	"fmt"
	"net"
	"testing"
)

type Client struct {
	addr string
}

func NewClient(addr string) *Client {
	return &Client{addr: addr}
}

func (c *Client) Hello(b []byte) (string, error) {
	conn, err := net.Dial("tcp", c.addr)
	if err != nil {
		return "", err
	}
	defer conn.Close()

	_, err = conn.Write(b)
	if err != nil {
		return "", err
	}

	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		return "", err
	}

	return string(buf[:n]), nil
}

func TestConnect(t *testing.T) {
	// ここだけ変更
	testConn(":5000")

	t.Errorf("actual %v want %v", 1, 2)
}

func testConn(addr string) (string, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return "", err
	}
	defer conn.Close()

	_, err = send(conn, "COMMANDER")
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
