package main

import (
	"fmt"
	"net"
	"time"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:8080")
	if err != nil {
		fmt.Println("net.Dial err:", err)
		return
	}

	for {
		_, err := conn.Write([]byte("hello zinx"))
		if err != nil {
			fmt.Println("write err:", err)
			continue
		}

		buf := make([]byte, 512)
		cnt, err := conn.Read(buf)
		if err != nil {
			fmt.Println("Read err:", err)
			continue
		}
		fmt.Println("read data:", string(buf), ", cnt:", cnt)

		time.Sleep(time.Second)
	}
}
