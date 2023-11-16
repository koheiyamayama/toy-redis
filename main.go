package main

import (
	"bufio"
	"bytes"
	"context"
	"log/slog"
	"net"
	"net/http"
	"os"

	"github.com/koheiyamayama/toy-redis/logger"
	"github.com/prometheus/client_golang/prometheus"
	promCollectors "github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	GREETING = []byte("Hello World")
	GET      = []byte("00000GET")
	SET      = []byte("00000SET")
	EXPIRE   = []byte("00EXPIRE")
)

var (
	cmdProcessed = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "toy_redis_command_counter",
		Help: "total of processed commands",
	}, []string{"command"})

	goRuntimeCollector = promCollectors.NewGoCollector(promCollectors.WithGoCollectorRuntimeMetrics(promCollectors.MetricsAll))
)

func main() {
	level := new(slog.LevelVar)
	level.Set(GetLogLevel())
	jHandler := slog.NewJSONHandler(
		os.Stdout,
		&slog.HandlerOptions{Level: level},
	)
	logger := slog.New(jHandler)
	slog.SetDefault(logger)

	l, err := net.Listen("tcp", "localhost:9999")
	if err != nil {
		slog.Error(err.Error())
		return
	}

	reg := prometheus.NewRegistry()
	reg.Register(goRuntimeCollector)
	reg.Register(cmdProcessed)

	kv := NewKV()

	slog.Info("start exposing metrics for Prometheus")
	http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))
	go func() {
		err := http.ListenAndServe(":8080", nil)
		if err != nil {
			slog.Error(err.Error())
		}
	}()

	slog.Info("start koheiyamayama/toy-redis")
	for {
		conn, err := l.Accept()
		if err != nil {
			slog.Error(err.Error())
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
	logger.InfoCtx(ctx, "command",
		slog.String("request", string(b)),
		slog.String("version", string(ver)),
		slog.String("command", string(command)),
		slog.String("payload", string(payload)),
	)

	switch {
	case bytes.Equal(command, GET):
		result, err = kv.Get(payload)
	case bytes.Equal(command, SET):
		key, value, exp := ParseSet(payload)
		kv.Set(key, value, exp)
	case bytes.Equal(command, EXPIRE):
		key, exp := ParseExpire(b)
		ok, eErr := kv.Expire(key, exp)
		if !ok {
			result = nil
			err = eErr
		}
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
