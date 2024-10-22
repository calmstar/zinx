package znet

import (
	"fmt"
	"net"
	"zinx/ziface"
)

type Server struct {
	Name      string
	IPVersion string
	IP        string
	Port      int
}

func (s *Server) Start() {
	fmt.Printf("[START] zinx server start at IP: %s, port:%d \n", s.IP, s.Port)
	go func() {
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Printf("ResolveTCPAddr err:%v \n", err)
			return
		}
		list, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			fmt.Printf("ListenTCP err:%v \n", err)
			return
		}
		// 监听端口成功
		fmt.Println("Zinx listening success, name:", s.Name)

		// 循环处理来自n个用户的链接请求
		for {
			conn, err := list.AcceptTCP()
			if err != nil {
				fmt.Printf("accept fail, err: %s \n", err)
				continue
			}
			// 处理某个用户的请求
			go func() {
				// 循环处理监听该用户发送过来的消息，并发送回去
				for {
					buf := make([]byte, 512)
					cnt, err := conn.Read(buf) // 该协程阻塞监听用户消息
					if err != nil {
						fmt.Println("recv data err:", err)
						continue
					}
					fmt.Println("read data: ", string(buf))
					_, err = conn.Write(buf[:cnt])
					if err != nil {
						fmt.Println("send data err:", err)
						continue
					}
				}
			}()
		}
	}()

}

func (s *Server) Stop() {
	fmt.Println("[STOP] zinx server, name:", s.Name)
}

func (s *Server) Serve() {
	s.Start()
}

func NewServer(name string) ziface.IServer {
	return &Server{
		Name:      name,
		IPVersion: "tcp4",
		IP:        "0.0.0.0",
		Port:      8080,
	}
}
