package main

import (
	"fmt"
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
	fmt.Println("handling, ", string(r.GetData()))
	r.GetConnection().GetTCPConnection().Write([]byte("hello"))
}

func (p *PingRouter) PostHandle(r ziface.IRequest) {
	fmt.Println("postHandle")
}

func (p *PingRouter) PreHandle(r ziface.IRequest) {
	fmt.Println("preHandle")
}
