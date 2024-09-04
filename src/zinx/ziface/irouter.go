package ziface

/*
	路由就是业务处理
*/

type IRouter interface {
	// 处理conn业务之前的钩子方法Hook
	PreHandle(request IRequest)

	// 处理conn业务的主方法Hook
	Handle(request IRequest)

	// 处理conn业务之后的钩子方法Hook
	PostHandle(request IRequest)
}
