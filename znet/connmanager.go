package znet

import (
	"errors"
	"fmt"
	"sync"
	"zinx/ziface"
)

type ConnManager struct {
	sync.RWMutex
	connections map[uint32]ziface.IConnection
}

func NewConnManager() *ConnManager {
	return &ConnManager{
		connections: make(map[uint32]ziface.IConnection),
	}
}

func (c *ConnManager) Add(conn ziface.IConnection) {
	c.Lock()
	defer c.Unlock()
	c.connections[conn.GetConnId()] = conn
}

func (c *ConnManager) Remove(id uint32) {
	c.Lock()
	defer c.Unlock()

	delete(c.connections, id)
}

func (c *ConnManager) ConnLen() uint32 {
	return uint32(len(c.connections))
}

func (c *ConnManager) Get(id uint32) (ziface.IConnection, error) {
	c.RLock()
	defer c.RUnlock()

	if res, ok := c.connections[id]; ok {
		return res, nil
	} else {
		return nil, errors.New("conn id not exist")
	}
}

func (c *ConnManager) ClearConn() {
	c.Lock()
	defer c.Unlock()

	for id, conn := range c.connections {
		fmt.Println("connID clear: ", id)
		conn.Stop()
	}

	c.connections = map[uint32]ziface.IConnection{}
}
