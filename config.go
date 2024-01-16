package main

import (
	"io"
	"log/slog"
	"os"
	"strconv"
)

func GetLogLevel() slog.Level {
	l := os.Getenv("TOY_REDIS_LOG_LEVEL")
	i, err := strconv.Atoi(l)
	if err != nil {
		return slog.LevelInfo
	}

	return slog.Level(i)
}

func GetLogFilePath() io.Writer {
	p := os.Getenv("TOY_REDIS_LOG_FILE_PATH")
	if p == "" {
		return os.Stdout
	}

	if f, err := os.OpenFile(p, os.O_RDWR, os.ModeAppend); err == nil {
		return f
	}

	return os.Stdout
}

func GetPyroscopeServerAddress() string {
	a := os.Getenv("PYROSCOPE_SERVER_ADDRESS")
	if a == "" {
		return "http://192.168.1.13:4040"
	} else {
		return a
	}
}

func GetHostName() string {
	return os.Getenv("HOSTNAME")
}
