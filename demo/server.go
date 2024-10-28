package main

import (
	"fmt"
	"io"
	"net"
	"zinx/ziface"
	"zinx/znet"
)

func main() {
	main2()
	//s := znet.NewServer()
	//s.AddRouter(&PingRouter{})
	//s.Serve()
	//
	//select {}
}

type PingRouter struct {
	*znet.BaseRouter
}

func (p *PingRouter) Handle(r ziface.IRequest) {
	fmt.Println("handling, ", string(r.GetData()))
	r.GetConnection().GetTCPConnection().Write([]byte("hello"))
}

func (p *PingRouter) PostHandle(r ziface.IRequest) {
	fmt.Println("postHandle")
}

func (p *PingRouter) PreHandle(r ziface.IRequest) {
	fmt.Println("preHandle")
}

func main2() {
	l, _ := net.Listen("tcp", "127.0.0.1:7777")

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
