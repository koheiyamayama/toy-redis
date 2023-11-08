package main

import (
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
