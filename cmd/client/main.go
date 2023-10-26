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

	// key := []byte("hoge")
	// value := []byte("fuga")
	// setPayload := fmt.Sprintf("%d%s%d%s", len(key), key, len(value), value)
	// setQuery := fmt.Sprintf("000100000SET%s\n", setPayload)
	// if n, err := conn.Write([]byte(setQuery)); err == nil {
	// 	fmt.Printf("write %d bytes\n", n)
	// } else {
	// 	fmt.Println(err.Error())
	// }

	if n, err := conn.Write([]byte("000100000GEThoge\n")); err == nil {
		fmt.Printf("write %d bytes\n", n)
	} else {
		fmt.Println(err.Error())
	}

	res := make([]byte, 1024)
	if _, err := conn.Read(res); err == nil {
		fmt.Println(string(res))
	} else {
		fmt.Println(err.Error())
	}
	conn.Close()
}
