package main

import (
	"fmt"
	"net"
	"time"
	"zinx/znet"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("net.Dial err:", err)
		return
	}

	i := 1
	for {
		// msg
		data := fmt.Sprintf("msg:%d", i)
		msg := znet.NewMsgPacket(uint32(i), []byte(data))

		pack := znet.NewDataPack()
		msgBytes, _ := pack.Pack(msg)
		_, err := conn.Write(msgBytes)
		if err != nil {
			fmt.Println("write err:", err)
			continue
		}

		//buf := make([]byte, 512)
		//cnt, err := conn.Read(buf)
		//if err != nil {
		//	fmt.Println("Read err:", err)
		//	continue
		//}
		//fmt.Println("read data:", string(buf), ", cnt:", cnt)

		time.Sleep(time.Second)
		i++
	}
}
