package main

import (
	"bytes"
	"fmt"
	"net"
)

var (
	GREETING = []byte("HelloWorld")
)

func main() {
	l, err := net.Listen("tcp", "localhost:9999")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println(err.Error())
		}

		go handleConn(conn)
	}
}

func handleConn(conn net.Conn) {
	defer conn.Close()

	m := 1024
	b := make([]byte, m)
	conn.Read(b)

	var err error
	// TODO: bytes.Containsからbytes.Equalに変えたい
	switch {
	case bytes.Contains(b, GREETING):
		_, err = conn.Write([]byte("Hey!"))
	default:
		_, err = conn.Write([]byte("NOP!"))
	}

	if err != nil {
		fmt.Println(err.Error())
		panic(err)
	}
}
