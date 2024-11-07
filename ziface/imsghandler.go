package ziface

// 消息处理模块
type IMsgHandler interface {
	DoMsgHandler(r IRequest)
	AddRouter(msgId uint32, r IRouter) // 为消息添加具体的处理逻辑

	StartWorkerPool()          // 开始协程池
	SendMsgToQueue(i IRequest) // 将消息交给taskQueue，由worker处理
}
