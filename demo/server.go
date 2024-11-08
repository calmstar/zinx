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
	s.AddRouter(1, &PingRouter{}) // 定义 id 映射 对应的路由信息
	s.AddRouter(2, &ZinxRouter{})

	s.SetOnConnStart(DoConnectionStart)
	s.SetOnConnStop(DoConnectionStop)

	s.Serve()
	select {}
}

type ZinxRouter struct {
	*znet.BaseRouter
}

func (z *ZinxRouter) Handle(r ziface.IRequest) {
	fmt.Printf("Serve ZinxRouter recv data, id：%v, data:%s \n", r.GetMsgId(), string(r.GetData()))
	// 这里的msgID，打包msg要用到，跟zinx的路由标识统一，也是为了是给客户端看的表示
	r.GetConnection().SendMsg(2, []byte("zinx router"))
}

type PingRouter struct {
	*znet.BaseRouter
}

func (p *PingRouter) Handle(r ziface.IRequest) {
	fmt.Printf("Serve PingRouter recv data, id：%v, data:%s \n", r.GetMsgId(), string(r.GetData()))
	// 发送数据
	//r.GetConnection().GetTCPConnection().Write([]byte("hello"))
	err := r.GetConnection().SendMsg(1, []byte("ping ping"))
	if err != nil {
		fmt.Println("handle err: ", err)
		return
	}
}

func DoConnectionStart(c ziface.IConnection) {
	fmt.Println("conn start do something")
	c.SetProperty("name", "ermao")
	c.SetProperty("home", "guangzhou")
	c.SendMsg(1, []byte("DoConnectionStart"))
}

func DoConnectionStop(c ziface.IConnection) {
	fmt.Println("conn stop do something")
	name, ok1 := c.GetProperty("name")
	home, ok2 := c.GetProperty("home")
	if ok1 != nil || ok2 != nil {

	}
	fmt.Printf("lost conn, name:%s, home:%s \n", name, home)
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
