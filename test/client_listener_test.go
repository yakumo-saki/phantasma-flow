package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net"
	"testing"
	"time"
)

type TestListenerClient struct {
	addr string
	conn net.Conn
}

func NewListnerClient(addr string) *TestListenerClient {
	lc := new(TestListenerClient)
	lc.addr = addr
	return lc
}

func (lc *TestListenerClient) send(msg string) (int, error) {
	io.Copy(lc.conn, bytes.NewBufferString(msg+"\n"))
	return 1, nil

	// length, err := lc.conn.Write([]byte(msg + "\n"))
	// if err != nil {
	// 	return 0, err
	// }

	// return length, err
}

func (lc *TestListenerClient) testListener(stop <-chan string) error {
	conn, err := net.Dial("tcp", lc.addr)
	if err != nil {
		fmt.Println("Error on dial")
		fmt.Println(err)
		return err
	}
	fmt.Println("Connect OK")
	lc.conn = conn
	defer lc.conn.Close()

	_, err = lc.send("LISTENER")
	if err != nil {
		fmt.Println("send LISTENER FAIL")
		return err
	}
	fmt.Println("send LISTENER OK")

	scanner := bufio.NewScanner(lc.conn)

	stopFlag := false
	for {
		select {
		case v := <-stop:
			fmt.Println("Signal received" + v)
			stopFlag = true
		default:
			if scanner.Scan() {
				line := scanner.Text() // スキャンした内容を文字列で取得
				fmt.Println(line)
			}
		}

		if stopFlag {
			fmt.Println("Exiting")
			break
		}
	}

	// for {
	// 	line, err := bufio.NewReader(lc.conn).ReadBytes('\n')
	// 	if err != nil {
	// 		fmt.Println("read FAIL")
	// 		return err
	// 	}
	// 	fmt.Printf("got: %s", line)
	// }
	fmt.Println("Exitted")
	return nil
}

func TestListener(t *testing.T) {
	// ここだけ変更
	lc := NewListnerClient(":5000")

	cha := make(chan string)
	go lc.testListener(cha)

	time.Sleep(5 * time.Second)
	fmt.Println("Signal")
	cha <- "STOP"

	t.Errorf("actual %v want %v", 1, 2)
}
