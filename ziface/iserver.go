package ziface

type IServer interface {
	Start()
	Stop()
	Serve()
	AddRouter(msgId uint32, r IRouter)
	GetConnManager() IConnManager
	SetOnConnStart(func(connection IConnection))
	SetOnConnStop(func(connection IConnection))
	CallOnConnStart(connection IConnection)
	CallOnConnStop(connection IConnection)
}
