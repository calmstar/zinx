package znet

import (
	"net"
	"zinx/ziface"
)

type Connection struct {
	Conn     *net.TCPConn
	ConnId   uint32
	IsClosed bool

	HandleApi      ziface.HandleFunc
	ExitedBuffChan chan bool
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
