package main

import (
	"fmt"
	"net"
	"time"
)

func newConn() net.Conn {
	conn, err := net.Dial("tcp", "localhost:9999")
	if err != nil {
		fmt.Println(err.Error())
	}

	return conn
}

func main() {
	conn := newConn()

	key := []byte("hogehogehogehoge")
	value := []byte("fugafuga2times")
	exp := []byte("3")
	setPayload := fmt.Sprintf("%s\r%s\r%s", key, value, exp)
	setQuery := fmt.Sprintf("000100000SET%s\n", setPayload)
	fmt.Println(setQuery)
	if n, err := conn.Write([]byte(setQuery)); err == nil {
		fmt.Printf("write %d bytes\n", n)
	} else {
		fmt.Println(err.Error())
	}
	conn.Close()

	time.Sleep(4 * time.Second)
	conn = newConn()
	if n, err := conn.Write([]byte("000100000GEThogehogehogehoge\n")); err == nil {
		fmt.Printf("write %d bytes\n", n)
	} else {
		fmt.Println(err.Error())
	}

	res := make([]byte, 1024)
	if _, err := conn.Read(res); err == nil {
		fmt.Println(string(res))
	} else {
		fmt.Println(err.Error())
	}
	conn.Close()
}
