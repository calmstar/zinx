package main

import (
	"fmt"
	"io"
	"net"
	"time"
	"zinx/znet"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:8080")
	if err != nil {
		fmt.Println("net.Dial err:", err)
		return
	}

	i := 1
	for {
		// 写入msg到conn中
		data := fmt.Sprintf("msg:%d", i)
		msg := znet.NewMsg(uint32(i), []byte(data))

		pack := znet.NewDataPack()
		msgBytes, _ := pack.Pack(msg)
		_, err := conn.Write(msgBytes)
		if err != nil {
			fmt.Println("write err:", err)
			continue
		}

		// 读取出msg
		headData := make([]byte, pack.GetHeadLen())
		io.ReadFull(conn, headData)
		msgHead, _ := pack.UnPack(headData)
		if msgHead.GetDataLen() > 0 {
			writeMsg := msgHead.(*znet.Message)
			writeMsg.Data = make([]byte, writeMsg.GetDataLen())
			io.ReadFull(conn, writeMsg.Data)
			fmt.Printf("client recv msg, id:%d, dataLen:%d, data:%s \n", writeMsg.GetMsgId(), writeMsg.GetDataLen(), writeMsg.GetData())
		}

		time.Sleep(time.Second)
		i++
	}
}
