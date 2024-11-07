package znet

import (
	"errors"
	"fmt"
	"net"
	"zinx/utils"
	"zinx/ziface"
)

type Connection struct {
	// 获取当前链接的socket tcp套接字
	Conn *net.TCPConn
	// 获取当前的链接id 全局唯一 可以作为sessionId
	ConnId uint32
	//当前链接的关闭状态
	IsClosed bool

	// 该链接的方法处理api
	//HandleApi ziface.HandleFunc
	// 消息路由处理功能
	msgHandler ziface.IMsgHandler

	// 消息管道，用来解耦对于客户端的读协程和写协程
	msgChan chan []byte

	//告知该链接已经退出/停止的 channel
	ExitedBuffChan chan struct{}

	// 当前conn属于哪个server
	TcpServer ziface.IServer
}

func NewConnection(server ziface.IServer, conn *net.TCPConn, connId uint32, msgHandler ziface.IMsgHandler) ziface.IConnection {
	c := &Connection{
		TcpServer: server,
		Conn:      conn,
		ConnId:    connId,
		IsClosed:  false,
		//HandleApi:      callBack,
		ExitedBuffChan: make(chan struct{}, 0),
		msgHandler:     msgHandler,
		msgChan:        make(chan []byte),
	}
	// 将连接添加到管理模块
	c.TcpServer.GetConnManager().Add(c)
	return c
}

// 写入客户端消息的协程
func (c *Connection) startWriter() {
	fmt.Println("write goroutine start...")
	defer fmt.Println("write goroutine end...")

	for {
		select {
		case <-c.ExitedBuffChan: //读goroutine，负责掌管是否关闭链接
			return
		case data := <-c.msgChan:
			_, err := c.Conn.Write(data)
			if err != nil {
				fmt.Printf("conn write err:%v \n", err)
				return // 发生报错会导致有问题
			}
		}
	}
}

// 读取客户端消息的协程
func (c *Connection) startReader() {
	fmt.Println("reader goroutine start......")
	defer fmt.Printf("reader goroutine end...\n")
	defer c.Stop() //该方法退出，意味着链接请求完毕，该用户退出

	for {
		// 读取链接中的数据
		//buf := make([]byte, utils.GlobalObject.MaxPacketSize)
		//_, err := c.Conn.Read(buf)
		//if err != nil {
		//	fmt.Println("c.Conn.Read err: ", err)
		//	continue
		//}
		// 业务处理
		//err = c.HandleApi(c.Conn, buf, cnt)
		//if err != nil {
		//	fmt.Println("HandleApi err: ", err)
		//	return
		//}

		// 创建拆包解包对象
		dp := NewDataPack()
		headData := make([]byte, dp.GetHeadLen())
		// 循环阻塞读协程，不断从链接读取数据出来
		// 先读头数据，解包，知道后续数据包长度；再读数据包，解包从消息
		if _, err := c.Conn.Read(headData); err != nil {
			fmt.Println("conn read data err: ", err)
			return
		}
		msg, err := dp.UnPack(headData) // int类型的数据，要单独根据大小字节序列进行解包
		if err != nil {
			fmt.Println("unpack err: ", err)
			return
		}
		data := make([]byte, msg.GetDataLen()) // string类型的数据，直接获得
		if _, err = c.Conn.Read(data); err != nil {
			fmt.Println("read data err: ", err)
			return
		}
		msg.SetData(data)
		req := &Request{
			conn: c,
			msg:  msg,
		}

		// 业务操作，再开协程 【可能会导致过多业务协程，所以下一节会使用业务协程池来处理workerPool】
		if utils.GlobalObject.WorkerPoolSize > 0 {
			c.msgHandler.SendMsgToQueue(req)
		} else {
			go func(r ziface.IRequest) {
				c.msgHandler.DoMsgHandler(r)
			}(req)
		}

	}
}

func (c *Connection) Start() {
	// 针对该链接的读取和操作
	go c.startReader()
	go c.startWriter()

	// 监控该链接是否操作完毕
	select {
	case <-c.ExitedBuffChan: // 读goroutine，负责掌管是否关闭链接
		return
	}
}

func (c *Connection) Stop() {
	if c.IsClosed {
		return
	}
	c.IsClosed = true
	close(c.ExitedBuffChan)
	err := c.Conn.Close()
	if err != nil {
		fmt.Println("Close err: ", err)
		return
	}
	c.TcpServer.GetConnManager().Remove(c.GetConnId())
}

func (c *Connection) SendMsg(msgId uint32, data []byte) error {
	if c.IsClosed == true {
		return errors.New("connection closed when send msg")
	}
	dp := NewDataPack()
	msg := NewMsg(msgId, data)
	pkData, err := dp.Pack(msg)
	if err != nil {
		fmt.Println("pack err: ", err)
		return err
	}
	// 解耦, 发送给写链接的协程
	c.msgChan <- pkData
	//_, err = c.Conn.Write(pkData)
	//if err != nil {
	//	fmt.Printf("conn write err:%v \n", err)
	//	c.ExitedBuffChan <- struct{}{}
	//	return fmt.Errorf("conn write err:%v", err)
	//}
	return nil
}

func (c *Connection) GetConnId() uint32 {
	return c.ConnId
}

func (c *Connection) GetRemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}
