package znet

import (
	"fmt"
	"log"
	"zinx/utils"
	"zinx/ziface"
)

// 消息管理模块, 实现路由功能，通过msgID来指向不同router处理方法

type MsgHandler struct {
	Apis map[uint32]ziface.IRouter // 存放每个msgId所对应处理方法的map属性

	WorkerPoolSize uint32                 // worker数量
	TaskQueues     []chan ziface.IRequest // worker取任务的消息队列
}

func NewMsgHandler() *MsgHandler {
	return &MsgHandler{
		Apis:           map[uint32]ziface.IRouter{},
		WorkerPoolSize: utils.GlobalObject.WorkerPoolSize,
		TaskQueues:     make([]chan ziface.IRequest, utils.GlobalObject.WorkerPoolSize),
	}
}

func (m *MsgHandler) DoMsgHandler(r ziface.IRequest) {
	msgId := r.GetMsgId()
	if router, exist := m.Apis[msgId]; !exist {
		log.Fatalf("DoMsgHandler err not exist, msgId: %v", msgId)
		return
	} else {
		router.PreHandle(r)
		router.Handle(r)
		router.PostHandle(r)
	}
}

func (m *MsgHandler) AddRouter(msgId uint32, r ziface.IRouter) {
	if _, exist := m.Apis[msgId]; exist {
		return
	}

	m.Apis[msgId] = r
}

func (m *MsgHandler) StartWorkerPool() {
	fmt.Println("StartWorkerPool, workerNum: ", m.WorkerPoolSize)
	for i := 0; i < int(m.WorkerPoolSize); i++ {
		// 给每个taskQueue指定长度（最大任务数量）
		m.TaskQueues[i] = make(chan ziface.IRequest, utils.GlobalObject.MaxWorkerTaskLen)
		// 一个任务对应一个queue
		go m.startOneWorker(i, m.TaskQueues[i])
	}
}

func (m *MsgHandler) startOneWorker(workerId int, taskChan chan ziface.IRequest) {
	fmt.Println("workerID: ", workerId, " 启动，监听业务消息")
	for {
		select {
		case req := <-taskChan: // 从客户端读出消息的goroutine，会往此taskChan写数据，让worker来取
			m.DoMsgHandler(req)
		}
	}
}

func (m *MsgHandler) SendMsgToQueue(i ziface.IRequest) {
	// 负载算法，使得每个queue都能平均
	connID := i.GetConnection().GetConnId()
	workerID := connID % m.WorkerPoolSize // connID是自己分配的。msgID带有业务路由属性，不要用
	fmt.Printf("add connID:%v, msgID:%v, to workerID:%v \n ", connID, i.GetMsgId(), workerID)
	taskQueue := m.TaskQueues[workerID]
	taskQueue <- i
}
