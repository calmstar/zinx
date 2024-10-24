package ziface

import "net"

type IConnection interface {
	//  连接开始工作
	Start()
	//  连接结束工作
	Stop()
	// 获取链接id
	GetConnId() uint32
	// 获取远程客户端的地址
	GetRemoteAddr() net.Addr
	// 获取tcp socket链接
	GetTCPConnection() *net.TCPConn
}

type HandleFunc func(conn *net.TCPConn, data []byte, len int) error
