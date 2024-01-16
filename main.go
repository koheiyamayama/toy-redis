package main

import (
	"bufio"
	"bytes"
	"context"
	"log/slog"
	"net"
	"net/http"
	"runtime"
	"time"

	"github.com/grafana/pyroscope-go"
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

	totalEntries = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "toy_redis_total_entries",
		Help: "total of entries in store",
	})
)

func main() {
	level := new(slog.LevelVar)
	level.Set(GetLogLevel())

	jHandler := slog.NewJSONHandler(
		GetLogFilePath(),
		&slog.HandlerOptions{Level: level},
	)
	logger := slog.New(jHandler)
	slog.SetDefault(logger)

	l, err := net.Listen("tcp", "0.0.0.0:9999")
	if err != nil {
		slog.Error(err.Error())
		return
	}

	reg := prometheus.NewRegistry()
	reg.Register(goRuntimeCollector)
	reg.Register(cmdProcessed)
	reg.Register(totalEntries)

	runtime.SetMutexProfileFraction(5)
	runtime.SetBlockProfileRate(5)

	pyroscope.Start(pyroscope.Config{
		ApplicationName: "toy-redis",

		// replace this with the address of pyroscope server
		ServerAddress: GetPyroscopeServerAddress(),

		// you can disable logging by setting this to nil
		Logger: pyroscope.StandardLogger,

		// you can provide static tags via a map:
		Tags: map[string]string{"hostname": GetHostName()},

		ProfileTypes: []pyroscope.ProfileType{
			// these profile types are enabled by default:
			pyroscope.ProfileCPU,
			pyroscope.ProfileAllocObjects,
			pyroscope.ProfileAllocSpace,
			pyroscope.ProfileInuseObjects,
			pyroscope.ProfileInuseSpace,

			// these profile types are optional:
			pyroscope.ProfileGoroutines,
			pyroscope.ProfileMutexCount,
			pyroscope.ProfileMutexDuration,
			pyroscope.ProfileBlockCount,
			pyroscope.ProfileBlockDuration,
		},
	})

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
	// 3秒ごとにGaugeを更新する
	// UpdateGaugeみたいな関数を作って呼び出す感じでも良いかもしれない
	go func() {
		for {
			time.Sleep(3 * time.Second)
			totalEntries.Set(float64(kv.Total()))
		}
	}()

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
	if err == nil {
		b = b[:len(b)-1]
	}

	var result []byte
	ver, command, payload := ParseQuery(b)
	logger.InfoCtx(ctx, "command",
		slog.String("request", string(b)),
		slog.String("version", string(ver)),
		slog.String("command", string(command)),
		slog.String("payload", string(payload)),
	)

	go func() {
		cmdProcessed.With(prometheus.Labels{"command": string(command)}).Inc()
	}()

	pyroscope.TagWrapper(context.Background(), pyroscope.Labels("handling_command", string(command)), func(ctx context.Context) {
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
	})

	if err != nil {
		slog.Info(err.Error())
		result = []byte("+" + err.Error())
	}

	if _, err := conn.Write(result); err != nil {
		slog.Info(err.Error())
		conn.Close()
	}
}
