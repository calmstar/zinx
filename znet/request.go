package znet

import "zinx/ziface"

type Request struct {
	conn ziface.IConnection
	msg  ziface.IMessage
}

func (r *Request) GetData() []byte {
	return r.msg.GetData()
}

func (r *Request) GetConnection() ziface.IConnection {
	return r.conn
}

func (r *Request) GetMsgId() uint32 {
	return r.msg.GetMsgId()
}
