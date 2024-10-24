package znet

import (
	"fmt"
	"net"
	"zinx/ziface"
)

type Connection struct {
	// 获取当前链接的socket tcp套接字
	Conn *net.TCPConn
	// 获取当前的链接id 全局唯一 可以作为sessionId
	ConnId uint32
	//当前链接的关闭状态
	IsClosed bool

	// 该链接的方法处理api
	HandleApi ziface.HandleFunc
	//告知该链接已经退出/停止的 channel
	ExitedBuffChan chan struct{}
}

func NewConnection(conn *net.TCPConn, connId uint32, callBack ziface.HandleFunc) ziface.IConnection {
	return &Connection{
		Conn:           conn,
		ConnId:         connId,
		IsClosed:       false,
		HandleApi:      callBack,
		ExitedBuffChan: make(chan struct{}, 0),
	}
}

func (c *Connection) StartReader() {
	fmt.Println("reader goroutine start......")
	defer fmt.Printf("reader goroutine end...")
	defer c.Stop()

	for {
		// 读取链接中的数据
		buf := make([]byte, 512)
		cnt, err := c.Conn.Read(buf)
		if err != nil {
			fmt.Println("c.Conn.Read err: ", err)
			continue
		}

		// 业务处理
		err = c.HandleApi(c.Conn, buf, cnt)
		if err != nil {
			fmt.Println("HandleApi err: ", err)
			return
		}
	}
}

func (c *Connection) Start() {
	// 针对该链接的读取和操作
	go c.StartReader()

	select {
	case <-c.ExitedBuffChan:
		return
	}
}

func (c *Connection) Stop() {
	if c.IsClosed {
		return
	}
	c.IsClosed = true
	close(c.ExitedBuffChan)
	err := c.Conn.Close()
	if err != nil {
		fmt.Println("Close err: ", err)
		return
	}
}

func (c *Connection) GetConnId() uint32 {
	return c.ConnId
}

func (c *Connection) GetRemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}
