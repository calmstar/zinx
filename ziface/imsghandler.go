package ziface

type IMsgHandler interface {
	DoMsgHandler(r IRequest)
	AddRouter(msgId uint32, r IRouter) // 为消息添加具体的处理逻辑
}
