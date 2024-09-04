package znet

import (
	"fmt"
	"go_code/src/zinx/utils"
	"go_code/src/zinx/ziface"
	"net"
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
}

func NewServer() *Server {
	s := &Server{
		Name:       utils.GlobalVar.Name,
		IPVersion:  "tcp4",
		IP:         utils.GlobalVar.Host,
		Port:       utils.GlobalVar.TcpPort,
		MsgHandler: NewMsgHandler(),
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

		// 3 连接处理
		dealConn := NewConnection(conn.(*net.TCPConn), cid, this.MsgHandler)
		cid++

		go dealConn.Start()
	}
}

// stop server
func (s *Server) Stop() {

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
