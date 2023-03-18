package client

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net"
	"net/textproto"
)

const MSG_SEPARATOR = "\x00\x00\x00"

type Client struct {
	addr string
	conn net.Conn
}

func NewClient(addr string) *Client {
	return &Client{addr: addr}
}

func (c *Client) Send(msg string) (int, error) {
	io.Copy(c.conn, bytes.NewBufferString(msg+"\n"))
	return 1, nil
}

func (c *Client) Ping(addr string) (string, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return "", err
	}
	defer conn.Close()
	c.conn = conn

	_, err = c.Send("PING")
	if err != nil {
		return "", err
	}

	reader := textproto.NewReader(bufio.NewReader(conn))

	ret := ""
	for {
		eof := false
		line, err := reader.ReadLine()
		if err == io.EOF {
			eof = true
		} else if err != nil {
			fmt.Println(err)
			break
		}
		if line == MSG_SEPARATOR {
			fmt.Println("got separator")
			break
		}

		ret = ret + line
		if eof {
			break
		}
	}

	fmt.Println(ret)
	return "", err
}
