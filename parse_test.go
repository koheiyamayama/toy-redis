package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseQuery(t *testing.T) {
	t.Parallel()

	type ret struct {
		version []byte
		command []byte
		value   []byte
	}

	type test struct {
		name string
		arg  []byte
		ret  ret
	}

	tests := []test{
		{
			name: "GET",
			arg:  []byte("000100000GET"),
			ret: ret{
				version: []byte("0001"),
				command: []byte("00000GET"),
				value:   []byte{},
			},
		},
		{
			name: "SET",
			arg:  []byte("000100000SETVALUE"),
			ret: ret{
				version: []byte("0001"),
				command: []byte("00000SET"),
				value:   []byte("VALUE"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			version, command, value := ParseQuery(tt.arg)
			assert.Equal(t, tt.ret.version, version)
			assert.Equal(t, tt.ret.command, command)
			assert.Equal(t, tt.ret.value, value)
		})
	}
}
