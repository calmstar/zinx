package znet

import (
	"log"
	"zinx/ziface"
)

// 消息管理模块, 实现路由功能，通过msgID来指向不同router处理方法

type MsgHandler struct {
	Apis map[uint32]ziface.IRouter // 存放每个msgId所对应处理方法的map属性
}

func NewMsgHandler() *MsgHandler {
	return &MsgHandler{Apis: map[uint32]ziface.IRouter{}}
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
