package znet

import (
	"fmt"
	"go_code/src/zinx/utils"
	"go_code/src/zinx/ziface"
	"net"
)

type Connection struct {
	// 连接属于哪个server
	TcpServer ziface.IServer

	// current connection
	Conn *net.TCPConn

	// connection ID
	ConnID uint32

	// connection status
	isClosed bool

	// channel to notify connection close
	ExitChan chan bool

	// 读写消息的通道
	MsgChan chan []byte

	// handler
	MsgHandler ziface.IMsgHandler
}

// initialize connection
func NewConnection(server ziface.IServer, conn *net.TCPConn, connID uint32, handler ziface.IMsgHandler) *Connection {
	c := &Connection{
		TcpServer:  server,
		Conn:       conn,
		ConnID:     connID,
		isClosed:   false,
		ExitChan:   make(chan bool, 1),
		MsgChan:    make(chan []byte),
		MsgHandler: handler,
	}

	c.TcpServer.GetConnMgr().Add(c)

	return c
}

// start connection
func (c *Connection) Start() {
	fmt.Println("Connection start()...ConnID = ", c.ConnID)
	go c.StartReader()
	go c.StartWriter()

	// 开始hook
	c.TcpServer.CallOnConnStart(c)
}

func (c *Connection) StartReader() {
	fmt.Println("Reader Goroutine is running...")
	defer fmt.Println("connID = ", c.ConnID, "Reader is exit, remote addr is ", c.RemoteAddr().String())
	defer c.Stop()

	for {
		pack := &DataPack{}

		// 阻塞读取客户端数据
		msg, err := pack.GetMsgFromConn(c.Conn)
		if err != nil {
			fmt.Println("GetMsgFromConn error: ", err)
			break
		}

		req := &Request{
			conn: c,
			msg:  msg,
		}

		if utils.GlobalVar.WorkerPoolSize > 0 {
			c.MsgHandler.SendMsgToTaskQueue(req)
		} else {
			go c.MsgHandler.DoMsgHandler(req)
		}
	}
}

// start write goroutine
func (c *Connection) StartWriter() {
	fmt.Println("[Writer Goroutine is running]")
	defer fmt.Println(c.RemoteAddr().String(), "[conn Writer exit!]")

	// 阻塞等待channel的消息
	for {
		select {
		case data := <-c.MsgChan:
			// 有数据要写给客户端
			if _, err := c.Conn.Write(data); err != nil {
				fmt.Println("Send data error: ", err)
				return
			}

		case <-c.ExitChan:
			// conn已经关闭
			return
		}
	}
}

// stop connection
func (c *Connection) Stop() {
	fmt.Println("Connection stop()...ConnID = ", c.ConnID)
	if c.isClosed == true {
		return
	}
	c.isClosed = true

	// 关闭hook
	c.TcpServer.CallOnConnStop(c)

	c.Conn.Close()

	// 写channel关闭
	c.ExitChan <- true

	// 将当前连接从连接管理器中删除
	c.TcpServer.GetConnMgr().Remove(c)

	close(c.ExitChan)
	close(c.MsgChan)
}

// get connection ID
func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

// get connection
func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

// get remote client address
func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

// send data to client
func (c *Connection) Send(data []byte) error {
	return nil
}

func (c *Connection) SendMsg(msg ziface.IMessage) error {
	if c.isClosed == true {
		return nil
	}

	// 将Message封装成二进制数据
	dp := &DataPack{}
	binaryMsg, err := dp.Pack(msg)
	if err != nil {
		fmt.Println("Pack error msg id = ", msg.GetMsgId())
		return err
	}

	// 发送数据给客户端
	if _, err := c.Conn.Write(binaryMsg); err != nil {
		fmt.Println("Write error msg id = ", msg.GetMsgId())
		return err
	}

	return nil
}

func (c *Connection) SendData(id uint32, data []byte) error {
	msg := &Message{
		Id:      id,
		DataLen: uint32(len(data)),
		Data:    data,
	}

	return c.SendMsg(msg)
}
