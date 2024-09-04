package ziface

type IServer interface {
	// start server
	Start()
	// stop server
	Stop()
	// run server
	Serve()

	// add router
	AddRouter(router IRouter)
}
