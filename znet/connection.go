package znet

import (
	"errors"
	"fmt"
	"net"
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
	Router ziface.IRouter

	//告知该链接已经退出/停止的 channel
	ExitedBuffChan chan struct{}
}

func NewConnection(conn *net.TCPConn, connId uint32, r ziface.IRouter) ziface.IConnection {
	return &Connection{
		Conn:     conn,
		ConnId:   connId,
		IsClosed: false,
		//HandleApi:      callBack,
		ExitedBuffChan: make(chan struct{}, 0),
		Router:         r,
	}
}

func (c *Connection) StartReader() {
	fmt.Println("reader goroutine start......")
	defer fmt.Printf("reader goroutine end...")
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
		if _, err := c.Conn.Read(headData); err != nil {
			fmt.Println("conn read data err: ", err)
			return
		}

		msg, err := dp.UnPack(headData)
		if err != nil {
			fmt.Println("unpack err: ", err)
			return
		}
		data := make([]byte, msg.GetDataLen())
		if _, err = c.Conn.Read(data); err != nil {
			fmt.Println("read data err: ", err)
			return
		}
		msg.SetData(data)

		req := Request{
			conn: c,
			msg:  msg,
		}
		go func(r ziface.IRequest) {
			c.Router.PreHandle(r)
			c.Router.Handle(r)
			c.Router.PostHandle(r)
		}(&req)
	}
}

func (c *Connection) Start() {
	// 针对该链接的读取和操作
	go c.StartReader()

	// 监控该链接是否操作完毕
	select {
	case <-c.ExitedBuffChan:
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
}

func (c *Connection) SendMsg(msgId uint32, data []byte) error {
	if c.IsClosed == true {
		return errors.New("connection closed when send msg")
	}
	dp := NewDataPack()
	msg := NewMsgPacket(msgId, data)
	pkData, err := dp.Pack(msg)
	if err != nil {
		fmt.Println("pack err: ", err)
		return err
	}
	_, err = c.Conn.Write(pkData)
	if err != nil {
		fmt.Printf("conn write err:%v \n", err)
		c.ExitedBuffChan <- struct{}{}
		return fmt.Errorf("conn write err:%v", err)
	}
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
