package main

import (
	"bytes"
	"strconv"
)

func ParseQuery(b []byte) (version []byte, command []byte, value []byte) {
	if len(b) >= 12 {

		version = b[0:4]
		command = b[4:12]
		value = b[12:]
		return version, command, value
	} else {
		return version, command, value
	}
}

func ParseSet(b []byte) (key, value []byte, exp uint32) {
	splitBytes := bytes.Split(b, []byte("\r"))
	if len(splitBytes) != 3 {
		return key, value, exp
	}
	// この変換どうにかしたい
	e, _ := strconv.Atoi(string(splitBytes[2]))

	return splitBytes[0], splitBytes[1], uint32(e)
}

func ParseExpire(b []byte) (key []byte, exp uint32) {
	splitBytes := bytes.Split(b, []byte("\r"))
	if len(splitBytes) == 2 {
		return key, exp
	}
	// この変換どうにかしたい
	e, _ := strconv.Atoi(string(splitBytes[2]))

	return splitBytes[0], uint32(e)
}
