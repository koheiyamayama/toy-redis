package main

import (
	"bufio"
	"bytes"
	"fmt"
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

	r := bufio.NewReader(conn)
	b, err := r.ReadBytes('\n')
	b = b[:len(b)-1]

	var result []byte
	ver, command, payload := ParseQuery(b)
	slog.Info("query", map[string]any{
		"query":   string(b),
		"version": string(ver),
		"command": string(command),
		"payload": string(payload),
	})

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
		slog.Info(err.Error())
		result = []byte(err.Error())
	}

	if _, err := conn.Write(result); err != nil {
		slog.Info(err.Error())
	}
}
