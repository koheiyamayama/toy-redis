package main

import (
	"fmt"
	"net"
)

func newConn() net.Conn {
	conn, err := net.Dial("tcp", "localhost:9999")
	if err != nil {
		fmt.Println(err.Error())
	}

	return conn
}

func main() {
	conn := newConn()
	key := "key"
	value := "value"
	exp := "3"
	setPayload := fmt.Sprintf("%s\r%s\r%s", key, value, exp)
	setQuery := fmt.Sprintf("000100000SET%s\n", setPayload)
	fmt.Println("set: ", setQuery)
	if n, err := conn.Write([]byte(setQuery)); err == nil {
		fmt.Printf("write %d bytes\n", n)
	} else {
		fmt.Println(err.Error())
	}
	conn.Close()

	conn = newConn()
	key = "key"
	value = "value"
	exp = "5"
	setPayload = fmt.Sprintf("%s\r%s\r%s", key, value, exp)
	setQuery = fmt.Sprintf("000100000SET%s\n", setPayload)
	fmt.Println("set: ", setQuery)
	if n, err := conn.Write([]byte(setQuery)); err == nil {
		fmt.Printf("write %d bytes\n", n)
	} else {
		fmt.Println(err.Error())
	}
	conn.Close()

	conn = newConn()
	if n, err := conn.Write([]byte("000100000GETkey\n")); err == nil {
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
