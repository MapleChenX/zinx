package znet

import (
	"fmt"
	"go_code/src/zinx/ziface"
	"net"
)

type Connection struct {
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
func NewConnection(conn *net.TCPConn, connID uint32, handler ziface.IMsgHandler) *Connection {
	c := &Connection{
		Conn:       conn,
		ConnID:     connID,
		isClosed:   false,
		ExitChan:   make(chan bool, 1),
		MsgChan:    make(chan []byte),
		MsgHandler: handler,
	}
	return c
}

// start connection
func (c *Connection) Start() {
	fmt.Println("Connection start()...ConnID = ", c.ConnID)
	go c.StartReader()
	go c.StartWriter()
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
			fmt.Println("unpack error: ", err)
			continue
		}

		data, err := pack.Pack(msg)
		if err != nil {
			fmt.Println("pack error: ", err)
			continue
		}

		c.MsgChan <- data
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

	c.Conn.Close()

	// 写channel关闭
	c.ExitChan <- true

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
