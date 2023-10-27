package main

import (
	"bytes"
	"log/slog"
)

func ParseQuery(b []byte) (version []byte, command []byte, value []byte) {
	if len(b) >= 12 {

		version = b[0:4]
		command = b[4:12]
		value = b[12:]
		return version, command, value
	} else {
		slog.Debug("query is not valid for protocol spec", map[string]any{"byte": string(b)})
		return version, command, value
	}
}

func ParseSet(b []byte) (key, value []byte) {
	splitBytes := bytes.Split(b, []byte("\r"))
	if len(splitBytes) != 2 {
		slog.Error("SET command handles two arguments: %s", string(b))
		return key, value
	}
	return splitBytes[0], splitBytes[1]
}
