package main

import "log/slog"

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
	return key, value
}
