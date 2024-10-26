package znet

import (
	"fmt"
	"net"
	"zinx/ziface"
)

type Server struct {
	Name      string
	IPVersion string
	// 服务器绑定的ip
	IP string
	// 服务器绑定的端口
	Port   int
	Router ziface.IRouter
}

func NewServer(name string) ziface.IServer {
	return &Server{
		Name:      name,
		IPVersion: "tcp4",
		IP:        "0.0.0.0",
		Port:      8080,
		Router:    nil,
	}
}

//func callBackToClient(conn *net.TCPConn, buf []byte, cnt int) error {
//	// 处理某个用户的请求
//	fmt.Println("read data: ", string(buf))
//	_, err := conn.Write(buf[:cnt])
//	if err != nil {
//		fmt.Println("send data err:", err)
//		return fmt.Errorf("callBackToClient err: %s", err)
//	}
//	return nil
//}

func (s *Server) AddRouter(r ziface.IRouter) {
	s.Router = r
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

		cid := 0
		// 循环处理来自n个用户的链接请求
		for {
			conn, err := list.AcceptTCP()
			if err != nil {
				fmt.Printf("accept fail, err: %s \n", err)
				continue
			}
			dealConn := NewConnection(conn, uint32(cid), s.Router)
			cid++

			// 起一个协程单独出来该用户链接
			go dealConn.Start()
		}
	}()

}

func (s *Server) Stop() {
	fmt.Println("[STOP] zinx server, name:", s.Name)
}

func (s *Server) Serve() {
	s.Start()
}
