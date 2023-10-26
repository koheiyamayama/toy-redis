package main

import (
	"bytes"
	"fmt"
	"io"
	"log/slog"
	"net"
)

var (
	GREETING = []byte("Hello World")
	GET      = []byte("00000GET")
	SET      = []byte("00000SET")
)

func main() {
	l, err := net.Listen("tcp", "localhost:9999")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	kv := NewKV()

	slog.Info("start koheiyamayama/toy-redis")
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println(err.Error())
		}

		go handleConn(conn, kv)
	}
}

func handleConn(conn net.Conn, kv *KV) {
	slog.Info("start handling connection")
	by := make([]byte, 1024)
	conn.Read(by)
	slog.Info("connection payload: " + string(by))
	b, err := io.ReadAll(conn)
	if err != nil {
		slog.Error(err.Error())
		panic(err)
	}

	var result []byte
	ver, command, payload := ParseQuery(b)
	fmt.Printf("version: %s, command: %s, payload: %s\n", ver, command, payload)
	switch {
	case bytes.Equal(command, GET):
		result, err = kv.GET(payload)
	case bytes.Equal(command, SET):
		key, value := ParseSet(payload)
		kv.SET(key, value)
	default:
		_, err = conn.Write([]byte("NOP!"))
	}

	if err != nil {
		fmt.Println(err.Error())
		panic(err)
	}

	if _, err := conn.Write(result); err != nil {
		fmt.Println(err.Error())
		panic(err)
	}
}
