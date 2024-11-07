package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"strconv"
	"time"
	"zinx/znet"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:8080")
	if err != nil {
		fmt.Println("net.Dial err:", err)
		return
	}
	flag.Parse()
	args := flag.Args()
	fmt.Println("参数个数：", flag.NArg(), " ,参数值：", args)

	msgId, _ := strconv.Atoi(args[0])
	i := 1
	for {
		// 写入msg到conn中
		data := fmt.Sprintf("msg:%d", i)
		i++
		if err := SendMsgBytes(conn, uint32(msgId), data); err != nil {
			continue
		}
		// 读取出msg
		ReadMsgBytes(conn)
		// 休眠
		time.Sleep(time.Second)
	}
}

func ReadMsgBytes(conn net.Conn) {
	pack := znet.NewDataPack()
	headData := make([]byte, pack.GetHeadLen())
	io.ReadFull(conn, headData)
	msgHead, _ := pack.UnPack(headData)
	if msgHead.GetDataLen() > 0 {
		writeMsg := msgHead.(*znet.Message)
		writeMsg.Data = make([]byte, writeMsg.GetDataLen())
		io.ReadFull(conn, writeMsg.Data)
		fmt.Printf("client recv msg, id:%d,  data:%s \n", writeMsg.GetMsgId(), writeMsg.GetData())
	}
}

func SendMsgBytes(conn net.Conn, msgId uint32, data string) error {
	msg := znet.NewMsg(msgId, []byte(data))

	pack := znet.NewDataPack()
	msgBytes, _ := pack.Pack(msg)
	_, err := conn.Write(msgBytes)
	if err != nil {
		fmt.Println("write err:", err)
		return err
	}
	return nil
}
