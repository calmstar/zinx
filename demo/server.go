package main

import (
	"fmt"
	"io"
	"net"
	"zinx/ziface"
	"zinx/znet"
)

func main() {
	s := znet.NewServer()
	s.AddRouter(&PingRouter{})
	s.Serve()

	select {}
}

type PingRouter struct {
	*znet.BaseRouter
}

func (p *PingRouter) Handle(r ziface.IRequest) {
	fmt.Printf("serve recv data, id：%v, data:%s \n", r.GetMsgId(), string(r.GetData()))

	// 发送数据
	//r.GetConnection().GetTCPConnection().Write([]byte("hello"))
	err := r.GetConnection().SendMsg(1, []byte("ping ping"))
	if err != nil {
		fmt.Println("handle err: ", err)
		return
	}
}

func (p *PingRouter) PostHandle(r ziface.IRequest) {
	fmt.Println("postHandle")
}

func (p *PingRouter) PreHandle(r ziface.IRequest) {
	fmt.Println("preHandle")
}

func main2() {
	l, _ := net.Listen("tcp", "127.0.0.1:8080")

	for {
		conn, _ := l.Accept()
		go func(conn net.Conn) {
			for {
				pack := znet.NewDataPack()
				headData := make([]byte, pack.GetHeadLen())
				io.ReadFull(conn, headData)

				headMsg, _ := pack.UnPack(headData)
				data := make([]byte, headMsg.GetDataLen())
				io.ReadFull(conn, data)

				fmt.Printf("data:%v , dataLen:%v, id:%v \n", string(data), headMsg.GetDataLen(), headMsg.GetMsgId())
			}
		}(conn)
	}

}
