package znet

import (
	"fmt"
	"net"
	"zinx/utils"
	"zinx/ziface"
)

type Server struct {
	Name      string
	IPVersion string
	// 服务器绑定的ip
	IP string
	// 服务器绑定的端口
	Port int
	// 消息处理模块：消息id路由处理功能；协程池worker处理业务逻辑功能
	msgHandler ziface.IMsgHandler
	// 连接管理
	connMgr ziface.IConnManager
}

func NewServer() ziface.IServer {
	return &Server{
		Name:       utils.GlobalObject.Name,
		IPVersion:  "tcp4",
		IP:         utils.GlobalObject.Host,
		Port:       utils.GlobalObject.TcpPort,
		msgHandler: NewMsgHandler(),
		connMgr:    NewConnManager(),
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

func (s *Server) AddRouter(msgId uint32, r ziface.IRouter) {
	s.msgHandler.AddRouter(msgId, r)
}

func (s *Server) Start() {
	fmt.Printf("[START] zinx server start at IP: %s, port:%d \n", s.IP, s.Port)
	go func() {
		// 开启业务处理协程池
		s.msgHandler.StartWorkerPool()

		// 监听服务启动开始
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
			// 判断连接是否超出限制
			if s.connMgr.ConnLen() > utils.GlobalObject.MaxConn {
				fmt.Printf("超出连接数量限制，当前链接数量：%v, 限制最大数量：%v", s.connMgr.ConnLen(), utils.GlobalObject.MaxConn)
				conn.Write([]byte("目前服务器链接数已经超过，请稍后再试"))
				conn.Close()
				continue
			}

			dealConn := NewConnection(s, conn, uint32(cid), s.msgHandler)
			cid++

			// 起一个协程单独处理该用户链接
			go dealConn.Start()
		}
	}()

}

func (s *Server) Stop() {
	fmt.Println("[STOP] zinx server, name:", s.Name)
	s.connMgr.ClearConn()
}

func (s *Server) Serve() {
	s.Start()
}

func (s *Server) GetConnManager() ziface.IConnManager {
	return s.connMgr
}
