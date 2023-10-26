package main

import (
	"fmt"
	"net"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:9999")
	if err != nil {
		fmt.Println(err.Error())
	}
	defer conn.Close()

	if n, err := conn.Write([]byte("000100000GEThoge")); err == nil {
		fmt.Printf("write %d bytes\n", n)
	} else {
		fmt.Println(err.Error())
	}

	res := make([]byte, 1024)
	if n, err := conn.Read(res); err == nil {
		fmt.Println(string(res), n)
	} else {
		fmt.Println(err.Error(), n)
	}
}
