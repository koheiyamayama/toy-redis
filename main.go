package main

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"net"
	"os"
)

var (
	GREETING = []byte("Hello World")
	GET      = []byte("00000GET")
	SET      = []byte("00000SET")
)

func main() {
	jHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{})
	logger := slog.New(jHandler)
	slog.SetDefault(logger)

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
	ctx := context.Background()
	slog.Info("start handling connection")
	defer func() {
		slog.Info("complete handling connection")
	}()

	r := bufio.NewReader(conn)
	b, err := r.ReadBytes('\n')
	b = b[:len(b)-1]

	var result []byte
	ver, command, payload := ParseQuery(b)
	slog.LogAttrs(ctx, slog.LevelInfo, "query",
		slog.String("query", string(b)),
		slog.String("version", string(ver)),
		slog.String("command", string(command)),
		slog.String("payload", string(payload)),
	)
	switch {
	case bytes.Equal(command, GET):
		result, err = kv.Get(payload)
	case bytes.Equal(command, SET):
		key, value := ParseSet(payload)
		kv.Set(key, value)
	default:
		result = []byte("NOP")
	}

	if err != nil {
		slog.Info(err.Error())
		result = []byte(err.Error())
	}

	if _, err := conn.Write(result); err != nil {
		slog.Info(err.Error())
		conn.Close()
	}
}
