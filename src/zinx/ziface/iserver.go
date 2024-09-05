package ziface

type IServer interface {
	// start server
	Start()
	// stop server
	Stop()
	// run server
	Serve()

	// add router
	AddRouter(id uint32, router IRouter)

	// get connection manager
	GetConnMgr() IConnManager

	// set hook
	SetOnConnStart(func(connection IConnection))
	SetOnConnStop(func(connection IConnection))
	CallOnConnStart(connection IConnection)
	CallOnConnStop(connection IConnection)
}
