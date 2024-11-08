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
	// 发送消息给客户端（无缓冲）
	SendMsg(msgId uint32, data []byte) error
	// 发送消息给客户端（有缓冲）
	SendBuffMsg(msgId uint32, data []byte) error

	SetProperty(key string, value interface{})
	GetProperty(key string) (interface{}, error)
	RemoveProperty(key string)
}

type HandleFunc func(conn *net.TCPConn, data []byte, len int) error
