package znet

import (
	"fmt"
	"go_code/src/zinx/utils"
	"go_code/src/zinx/ziface"
	"net"
	"time"
)

type Server struct {
	// server name
	Name string
	// server IP version
	IPVersion string
	// server IP
	IP string
	// server port
	Port int
	// cmd handler
	MsgHandler ziface.IMsgHandler

	// connection manager
	ConnMgr ziface.IConnManager

	// hook
	OnConnStart func(conn ziface.IConnection)
	OnConnStop  func(conn ziface.IConnection)
}

func NewServer() *Server {
	s := &Server{
		Name:       utils.GlobalVar.Name,
		IPVersion:  "tcp4",
		IP:         utils.GlobalVar.Host,
		Port:       utils.GlobalVar.TcpPort,
		MsgHandler: NewMsgHandler(),
		ConnMgr:    NewConnManager(),
	}
	return s
}

// start server
func (this *Server) Start() {
	// 启动worker工作池
	this.MsgHandler.StartWorkerPool()

	// 1 连接
	listener, err := net.Listen(this.IPVersion, fmt.Sprintf("%s:%d", this.IP, this.Port))
	if err != nil {
		fmt.Println("net.Listen err:", err)
		return
	}
	defer listener.Close()

	var cid uint32
	cid = 0

	// 2 处理客户端请求
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("listener.Accept err:", err)
			continue
		}

		// 2.1 设置最大连接数
		if this.ConnMgr.Len() >= utils.GlobalVar.MaxConn {
			message := &Message{
				Id:   100,
				Data: []byte("Too many connections, server busy, please try again later\n"),
			}
			message.DataLen = uint32(len(message.Data))

			dp := NewDataPack()
			binaryMsg, err := dp.Pack(message)
			if err != nil {
				fmt.Println("Pack error:", err)
				return
			}

			conn.SetWriteDeadline(time.Now().Add(1 * time.Second))
			conn.Write(binaryMsg)
			conn.Close()
			continue
		}

		// 3 连接处理
		dealConn := NewConnection(this, conn.(*net.TCPConn), cid, this.MsgHandler)
		cid++

		go dealConn.Start()
	}
}

// stop server
func (s *Server) Stop() {
	fmt.Println("[STOP] Zinx server name ", s.Name)
	s.ConnMgr.ClearConn()
}

// run server
func (s *Server) Serve() {
	go s.Start()
	select {}
}

// add router
func (s *Server) AddRouter(id uint32, router ziface.IRouter) {
	s.MsgHandler.AddRouter(id, router)
	fmt.Println("Add Router success!")
}

// get connection manager
func (s *Server) GetConnMgr() ziface.IConnManager {
	return s.ConnMgr
}

// set hook
func (s *Server) SetOnConnStart(hook func(connection ziface.IConnection)) {
	s.OnConnStart = hook
}

func (s *Server) SetOnConnStop(hook func(connection ziface.IConnection)) {
	s.OnConnStop = hook
}

func (s *Server) CallOnConnStart(connection ziface.IConnection) {
	if s.OnConnStart != nil {
		fmt.Println("连接建立时的回调函数启动！")
		s.OnConnStart(connection)
	}
}

func (s *Server) CallOnConnStop(connection ziface.IConnection) {
	if s.OnConnStop != nil {
		fmt.Println("连接关闭时的回调函数启动！")
		s.OnConnStop(connection)
	}
}
