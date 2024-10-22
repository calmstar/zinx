package ziface

import "net"

type IConnection interface {
	Start()
	Stop()
	GetConnId() uint32
	RemoteAddr() net.Addr
	GetTCPConnection() *net.TCPConn
}

type HandleFunc func(conn *net.TCPConn, data []byte, len int) error
